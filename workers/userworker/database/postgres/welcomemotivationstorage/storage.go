package welcomemotivationstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/errorutils/welcomemotivationerror"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/welcomemotivationsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/lib/pq"
)

const _PROVIDER = "workers/userworker/database/postgres/welcomemotivationstorage.Storage"

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
		case usersql.WelcomeMotivationsForeignKey:
			return exception.Wrap(welcomemotivationerror.ExceptionMotivationNotExist(), cause)
		}
	}
	return exception.Wrap(err, cause)
}

var nextUserWelcomeMotivationQuery = fmt.Sprintf(
	`
	WITH user_motivation as (
		SELECT %s as id FROM %s WHERE %s = $1
	),
	min_motivation as (
		SELECT coalesce(min(%s), (SELECT min(%s) FROM %s)) as id FROM %s WHERE %s > (SELECT id FROM user_motivation)
	)
	SELECT %s,%s FROM %s WHERE %s = (SELECT id FROM min_motivation)
	`,
	usersql.WelcomeMotivationId,
	usersql.Table,
	usersql.ID,

	// select min welcome motivation id
	welcomemotivationsql.ID,

	// select min welcome motivation in coalesce
	welcomemotivationsql.ID,
	welcomemotivationsql.Table,

	// select from table
	welcomemotivationsql.Table,
	// where id > current user welcome motivation id
	welcomemotivationsql.ID,

	welcomemotivationsql.ID,
	welcomemotivationsql.Motivation,

	welcomemotivationsql.Table,
	welcomemotivationsql.ID,
)

func (s *Storage) NextUserWelcomeMotivation(ctx context.Context, userId int64) (model.Motivation, error) {
	var motivation model.Motivation
	err := s.db.QueryRowContext(ctx, nextUserWelcomeMotivationQuery, userId).Scan(&motivation.ID, &motivation.Motivation)
	if err != nil {
		return motivation, Error(err, exception.NewCause("get next user motivation", "NextUserWelcomeMotivation", _PROVIDER))
	}
	return motivation, nil
}

var updateUserWelcomeMotivation = fmt.Sprintf(
	`UPDATE %s SET %s = $2 WHERE %s = $1`,
	usersql.Table,
	usersql.WelcomeMotivationId,
	usersql.ID,
)

func (s *Storage) UpdateUserWelcomeMotivation(ctx context.Context, userId int64, motivationId int) error {
	r, err := s.db.ExecContext(ctx, updateUserWelcomeMotivation, userId, motivationId)
	if err != nil {
		return Error(err, exception.NewCause("update user welcome motivation", "UpdateUserWelcomeMotivation", _PROVIDER))
	}
	if rows, _ := r.RowsAffected(); rows == 0 {
		return exception.Wrap(usererror.ExceptionUserNotFound(), exception.NewCause("check rows", "UpdateUserMotivation", _PROVIDER))
	}
	return nil
}
