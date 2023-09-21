package motivationstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/errorutils/motivationerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/motivationsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/lib/pq"
)

const _PROVIDER = "workers/userworker/database/postgres/motivationstorage.Storage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Constraint {
		case usersql.MotivationsForeignKey:
			return exception.Wrap(motivationerror.ExceptionMotivationNotExist(), cause)
		}
	}
	return exception.Wrap(err, cause)
}

type MotivationStorage interface {
	NextUserMotivation(ctx context.Context, userId int64) (model.Motivation, error)
	UpdateUserMotivation(ctx context.Context, userId int64, motivationId int) error
}

var nextUserMotivationQuery = fmt.Sprintf(`
	WITH user_motivation as (
		SELECT %s as id FROM %s WHERE %s = $1
	),
	min_motivation as (
		SELECT coalesce(min(%s), (SELECT id FROM user_motivation)) as id FROM %s WHERE %s > (SELECT id FROM user_motivation)
	)
	SELECT %s,%s FROM %s WHERE %s = (SELECT id FROM min_motivation)
`,
	usersql.MotivationId,
	usersql.Table,
	usersql.ID,

	motivationsql.ID,
	motivationsql.Table,
	motivationsql.ID,

	motivationsql.ID,
	motivationsql.Motivation,
	motivationsql.Table,
	motivationsql.ID,
)

func (s *Storage) NextUserMotivation(ctx context.Context, userId int64) (model.Motivation, error) {
	var motivation model.Motivation
	err := s.db.QueryRowContext(ctx, nextUserMotivationQuery, userId).Scan(&motivation.ID, &motivation.Motivation)
	if err != nil {
		return motivation, Error(err, exception.NewCause("get next user motivation", "NextUserMotivation", _PROVIDER))
	}
	return motivation, nil
}

var updateUserMotivationQuery = fmt.Sprintf(`
	UPDATE %s SET %s = $2 WHERE %s = $1
`,
	usersql.Table,
	usersql.MotivationId,
	usersql.ID,
)

func (s *Storage) UpdateUserMotivation(ctx context.Context, userId int64, motivationId int) error {
	r, err := s.db.ExecContext(ctx, updateUserMotivationQuery, userId, motivationId)
	if err != nil {
		return Error(err, exception.NewCause("update user motivation", "UpdateUserMotivation", _PROVIDER))
	}
	if rows, _ := r.RowsAffected(); rows == 0 {
		return exception.Wrap(usererror.ExceptionUserNotFound(), exception.NewCause("check rows", "UpdateUserMotivation", _PROVIDER))
	}
	return nil
}
