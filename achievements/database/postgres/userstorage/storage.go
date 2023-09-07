package userstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "achievements/database/postgres/userstorage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

var userQuery = fmt.Sprintf(
	`SELECT %s,%s,%s,%s FROM %s WHERE %s = $1`,
	usersql.CigaretteDayAmount,
	usersql.CigarettePackAmount,
	usersql.PackPrice,
	usersql.AbstinenceTime,

	usersql.Table,
	usersql.ID,
)

func (s *Storage) User(ctx context.Context, userId int64) (*model.UserData, error) {
	var data model.UserData
	absTime := amidtime.Timestamp{}
	err := s.db.QueryRowContext(ctx, userQuery, userId).Scan(
		&data.CigaretteDayAmount,
		&data.CigarettePackAmount,
		&data.PackPrice,
		&absTime,
	)
	if err != nil {
		return nil, Error(err, exception.NewCause("get user data", "User", _PROVIDER))
	}
	data.AbstinenceTime = absTime.Time
	return &data, nil
}
