package triggerstorage

import (
	"context"
	"database/sql"
	"fmt"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/errorutils/usertriggererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/triggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usertriggersql"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "user/database/postgres/triggerstorage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
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
		slog.Error(Error(err, exception.NewCause("failed get rows affected", "Remove", _PROVIDER)).Error())
		return nil
	}
	if rows != 1 {
		slog.Info("remove trigger affected wrong amount of rows", "rows", rows, "provider", _PROVIDER)
		return exception.Wrap(usertriggererror.UserTriggerNotFound(), exception.NewCause("wrong rows amount", "Remove", _PROVIDER))
	}
	return nil
}
