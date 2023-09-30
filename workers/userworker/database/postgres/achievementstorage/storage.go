package achievementstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	ORDER BY %s ASC, %s ASC
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
	sqlutils.Full(achievementsql.TypeId),
	sqlutils.Full(achievementsql.Level),
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

var insertUserAchievement = fmt.Sprintf(
	`
	INSERT INTO %s (%s,%s,%s) VALUES($1,$2,$3)
	`,
	userachievementsql.Table,

	userachievementsql.UserId,
	userachievementsql.AchievementId,
	userachievementsql.ReachDate,
)

func (s *Storage) InsertUserAchievements(ctx context.Context, userId int64, reachDate time.Time, achievementsIds []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return Error(err, exception.NewCause("begin tx", "InsertUserAchievements", _PROVIDER))
	}
	defer tx.Rollback()
	for _, id := range achievementsIds {
		_, err := tx.ExecContext(ctx, insertUserAchievement, userId, id, amidtime.Timestamp{Time: reachDate})
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
