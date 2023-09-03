package achievementstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
)

const _PROVIDER = "user/database/postgres/achievementstorage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
	return exception.Wrap(err, cause)
}

var achievementPreviewQuery = fmt.Sprintf(
	`
	SELECT DISTINCT ON (%s) %s FROM %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s = %s
	WHERE %s = $1 AND %s IS NOT NULL

	ORDER BY %s DESC
	`,
	sqlutils.Full(achievementtypesql.Type),
	sqlutils.Full(
		achievementsql.Level,
		achievementtypesql.Type,
	),
	achievementsql.Table,

	// inner join achievement types
	achievementtypesql.Table,
	sqlutils.Full(achievementsql.TypeId),
	sqlutils.Full(achievementtypesql.ID),

	// inner join user achievements
	userachievementsql.Table,
	sqlutils.Full(achievementsql.ID),
	sqlutils.Full(userachievementsql.AchievementId),

	// where user id eq $1
	sqlutils.Full(userachievementsql.UserId),
	sqlutils.Full(userachievementsql.OpenDate),

	// order by achievement level
	sqlutils.Full(achievementtypesql.Type, achievementsql.Level),
)

func (s *Storage) AchievementPreview(ctx context.Context, userId int64) []*usermodel.Achievement {
	rows, err := s.db.QueryContext(ctx, achievementPreviewQuery, userId)
	if err != nil {
		return make([]*usermodel.Achievement, 0)
	}
	defer rows.Close()
	achievements := make([]*usermodel.Achievement, 0)
	for rows.Next() {
		var ach usermodel.Achievement
		err := rows.Scan(&ach.Level, &ach.Type)
		if err != nil {
			return achievements
		}
		achievements = append(achievements, &ach)
	}
	return achievements
}
