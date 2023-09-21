package achievementstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
)

const _PROVIDER = "achievements/database/postgres/achievementstorage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
	return exception.Wrap(err, cause)
}

var userAchievementQuery = fmt.Sprintf(
	`
	SELECT %s, coalesce(%s,NULL),coalesce(%s,NULL),coalesce(%s,TRUE) FROM %s
	INNER JOIN %s ON %s = %s 
	LEFT JOIN %s ON %s = %s AND %s = $1
	GROUP BY %s
	ORDER BY %s
	`,
	sqlutils.Full(
		achievementsql.ID,
		achievementtypesql.Type,
		achievementsql.Exp,
		achievementsql.Level,
		achievementsql.Description,
	),
	sqlutils.Full(userachievementsql.OpenDate),
	sqlutils.Full(userachievementsql.ReachDate),
	sqlutils.Full(userachievementsql.Shown),

	achievementsql.Table,

	// inner join
	achievementtypesql.Table,
	sqlutils.Full(achievementsql.TypeId),
	sqlutils.Full(achievementtypesql.ID),

	// left join
	userachievementsql.Table,
	sqlutils.Full(achievementsql.ID),
	sqlutils.Full(userachievementsql.AchievementId),

	// where user id eq $1
	sqlutils.Full(userachievementsql.UserId),

	// group by
	sqlutils.Full(
		achievementsql.ID,
		achievementtypesql.Type,
		achievementsql.Exp,
		achievementsql.Level,
		achievementsql.Description,
		userachievementsql.OpenDate,
		userachievementsql.ReachDate,
		userachievementsql.Shown,
	),

	// order by
	sqlutils.Full(
		achievementsql.TypeId,
		achievementsql.Level,
	),
)

func (s *Storage) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	rows, err := s.db.QueryContext(ctx, userAchievementQuery, userId)
	if err != nil {
		return nil, Error(err, exception.NewCause("user achievement query", "UserAchievements", _PROVIDER))
	}
	defer rows.Close()
	achievements := make([]*achievementmodel.Achievement, 0)
	for rows.Next() {
		var achievement achievementmodel.Achievement
		err := rows.Scan(
			&achievement.ID,
			&achievement.Type,
			&achievement.Exp,
			&achievement.Level,
			&achievement.Description,
			&achievement.OpenDate,
			&achievement.ReachDate,
			&achievement.Shown,
		)
		if err != nil {
			return nil, Error(err, exception.NewCause("scan user achievement", "UserAchievements", _PROVIDER))
		}
		achievements = append(achievements, &achievement)
	}
	return achievements, nil
}

var markShownQuery = fmt.Sprintf(
	`
	UPDATE %s SET %s = TRUE WHERE %s = $1 AND %s IS FALSE
	`,
	userachievementsql.Table,
	userachievementsql.Shown,
	userachievementsql.UserId,
	userachievementsql.Shown,
)

func (s *Storage) MarkShown(ctx context.Context, userId int64) error {
	_, err := s.db.ExecContext(ctx, markShownQuery, userId)
	if err != nil {
		return Error(err, exception.NewCause("mark user achievement shown", "MarkShown", _PROVIDER))
	}
	return nil
}

var insertUserAchievement = fmt.Sprintf(
	`
	INSERT INTO %s (%s,%s,%s) VALUES($1,$2,$3)
	`,
	userachievementsql.Table,

	userachievementsql.UserId,
	userachievementsql.AchievementId,
	userachievementsql.ReachDate,
)

func (s *Storage) InsertUserAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return Error(err, exception.NewCause("begin tx", "InsertUserAchievements", _PROVIDER))
	}
	defer tx.Rollback()
	for _, id := range achievementsIds {
		_, err := tx.ExecContext(ctx, insertUserAchievement, userId, id, reachDate)
		if err != nil {
			return Error(err, exception.NewCause("insert user achievement", "InsertUserAchievements", _PROVIDER))
		}
	}
	err = tx.Commit()
	if err != nil {
		return Error(err, exception.NewCause("commit tx", "InsertUserAchievements", _PROVIDER))
	}
	return nil
}

var openSingleQuery = fmt.Sprintf(
	`
	UPDATE %s SET %s = $1 WHERE %s = $2 AND %s = $3 AND %s IS NULL
	`,
	userachievementsql.Table,
	userachievementsql.OpenDate,
	userachievementsql.AchievementId,
	userachievementsql.UserId,
	userachievementsql.OpenDate,
)

func (s *Storage) OpenSingle(ctx context.Context, userId int64, ach model.OpenAchievement) error {
	tx, err := s.db.Begin()
	if err != nil {
		return Error(err, exception.NewCause("begin tx", "OpenSingle", _PROVIDER))
	}
	defer tx.Rollback()
	r, err := tx.ExecContext(ctx, openSingleQuery, ach.OpenTime, ach.AchievementId, userId)
	if err != nil {
		return Error(err, exception.NewCause("exec open single query", "OpenSingle", _PROVIDER))
	}
	rows, _ := r.RowsAffected()
	if rows == 0 || rows > 1 {
		return exception.Wrap(usererror.ExceptionUserNotFound(), exception.NewCause("no rows", "OpenSingle", _PROVIDER))
	}
	err = tx.Commit()
	if err != nil {
		return Error(err, exception.NewCause("commit tx", "OpenSingle", _PROVIDER))
	}
	return nil
}

var openTypeQuery = fmt.Sprintf(
	`
	UPDATE %s SET %s = $1 FROM %s 
	INNER JOIN %s ON %s = %s
	WHERE %s = %s AND %s = $2 AND %s = $3 AND %s IS NULL
	RETURNING %s
	`,

	userachievementsql.Table,
	userachievementsql.OpenDate,

	achievementsql.Table,

	achievementtypesql.Table,
	sqlutils.Full(achievementsql.TypeId),
	sqlutils.Full(achievementtypesql.ID),

	// whre userachievement achievement id eq achievement id
	sqlutils.Full(userachievementsql.AchievementId),
	sqlutils.Full(achievementsql.ID),

	sqlutils.Full(achievementtypesql.Type),
	sqlutils.Full(userachievementsql.UserId),

	sqlutils.Full(userachievementsql.OpenDate),

	sqlutils.Full(userachievementsql.AchievementId),
)

func (s *Storage) OpenType(ctx context.Context, userId int64, ach model.OpenAchievementType) ([]int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, Error(err, exception.NewCause("begin tx", "OpenType", _PROVIDER))
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, openTypeQuery, ach.OpenTime, ach.AchievementType, userId)
	if err != nil {
		return nil, Error(err, exception.NewCause("exec openTypeQuery", "OpenType", _PROVIDER))
	}
	defer rows.Close()
	achIds := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, Error(err, exception.NewCause("scan achievement id", "OpenType", _PROVIDER))
		}
		achIds = append(achIds, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("commit tx", "OpenType", _PROVIDER))
	}
	return achIds, nil
}

var openAllQuery = fmt.Sprintf(
	`
	UPDATE %s SET %s = $1 WHERE %s IS NULL AND %s = $2 RETURNING %s
	`,
	userachievementsql.Table,
	userachievementsql.OpenDate,
	userachievementsql.OpenDate,
	userachievementsql.UserId,
	userachievementsql.AchievementId,
)

func (s *Storage) OpenAll(ctx context.Context, userId int64, openTime amidtime.Timestamp) ([]int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, Error(err, exception.NewCause("begin tx", "OpenAll", _PROVIDER))
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, openAllQuery, openTime, userId)
	if err != nil {
		return nil, Error(err, exception.NewCause("exec open all query", "OpenAll", _PROVIDER))
	}
	defer rows.Close()
	achIds := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, Error(err, exception.NewCause("scan achievement id", "OpenAll", _PROVIDER))
		}
		achIds = append(achIds, id)
	}
	err = tx.Commit()
	if err != nil {
		return nil, Error(err, exception.NewCause("commit tx", "OpenAll", _PROVIDER))
	}
	return achIds, nil
}

var selectAchievementMotivationQuery = fmt.Sprintf(`
	SELECT %s FROM %s WHERE %s = $1
`,
	achievementsql.Motivation,
	achievementsql.Table,
	achievementsql.ID,
)

func (s *Storage) AchievementMotivation(ctx context.Context, achId int64) (string, error) {
	var motivation string
	err := s.db.QueryRowContext(ctx, selectAchievementMotivationQuery, achId).Scan(&motivation)
	if err != nil {
		return motivation, Error(err, exception.NewCause("get achievement motivation", "AchievementMotivation", _PROVIDER))
	}
	return motivation, nil
}
