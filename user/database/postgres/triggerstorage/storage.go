package triggerstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/errorutils/triggererror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usertriggererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/triggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usertriggersql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/lib/pq"
)

const _PROVIDER = "user/database/postgres/triggerstorage"

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
		case usertriggersql.UsersForeignKey:
			return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
		case usertriggersql.TriggersForeignKey:
			return exception.Wrap(triggererror.ExceptionTriggerNotExist(), cause)
		case usertriggersql.PrimaryKey:
			return exception.Wrap(usertriggererror.UserTriggerExists(), cause)
		}
	}
	return exception.Wrap(err, cause)
}

var removeTriggerQuery = fmt.Sprintf(
	`
	WITH type_select as (
		SELECT %s as type FROM %s WHERE %s = $2
	)
	DELETE FROM %s WHERE %s = $1 AND %s = (SELECT type FROM type_select)
	`,
	triggersql.ID,
	triggersql.Table,
	triggersql.Name,
	usertriggersql.Table,
	usertriggersql.UserId,
	usertriggersql.TriggerId,
)

func (s *Storage) Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) error {
	res, err := s.db.ExecContext(ctx, removeTriggerQuery, userId, trigger)
	if err != nil {
		return Error(err, exception.NewCause("remove user trigger", "Remove", _PROVIDER))
	}
	rows, err := res.RowsAffected()
	if err != nil {
		err = Error(err, exception.NewCause("failed get rows affected", "Remove", _PROVIDER))
		slog.Error(err.Error())
		return nil
	}
	if rows != 1 {
		slog.Info("remove trigger affected wrong amount of rows", "rows", rows, "provider", _PROVIDER)
		return exception.Wrap(usertriggererror.UserTriggerNotFound(), exception.NewCause("wrong rows amount", "Remove", _PROVIDER))
	}
	return nil
}

var addTriggerQuery = fmt.Sprintf(`
WITH type_select as (
	SELECT %s as type FROM %s WHERE %s = $2
)
INSERT INTO %s (%s, %s) VALUES ($1, coalesce((SELECT type FROM type_select),0))
`,
	triggersql.ID,
	triggersql.Table,
	triggersql.Name,

	usertriggersql.Table,
	usertriggersql.UserId,
	usertriggersql.TriggerId,
)

func (s *Storage) Add(ctx context.Context, userId int64, trigger usermodel.Trigger) error {
	_, err := s.db.ExecContext(ctx, addTriggerQuery, userId, trigger)
	if err != nil {
		return Error(err, exception.NewCause("add trigger to user", "Add", _PROVIDER))
	}
	return nil
}
