package privacysettingstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/errorutils/userprivacysettingerror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/privacysettingsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userprivacysettingsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/lib/pq"
)

const _PROVIDER = "user/database/postgres/privacysettingstorage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Column {
		case userprivacysettingsql.SettingId.String():
			return exception.Wrap(userprivacysettingerror.ExceptionUserPrivacySettingNotFound(), cause)
		}
		switch pqErr.Constraint {
		case userprivacysettingsql.PrimaryKey:
			return exception.Wrap(userprivacysettingerror.ExceptionUserPrivacySettingExists(), cause)
		case userprivacysettingsql.ForeignKeyUsers:
			return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
		case userprivacysettingsql.ForeignKeyPrivacySettings:
			return exception.Wrap(userprivacysettingerror.ExceptionUserPrivacySettingNotFound(), cause)
		}
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return exception.Wrap(userprivacysettingerror.ExceptionUserPrivacySettingNotFound(), cause)
	default:
		return exception.Wrap(err, cause)
	}
}

var userPrivacySettingsQuery = fmt.Sprintf(
	`
	SELECT array_agg(%s) FROM %s
	INNER JOIN %s ON %s = %s
	WHERE %s = $1
	`,
	sqlutils.Full(privacysettingsql.Type),
	privacysettingsql.Table,

	userprivacysettingsql.Table,
	sqlutils.Full(privacysettingsql.ID),
	sqlutils.Full(userprivacysettingsql.SettingId),

	sqlutils.Full(userprivacysettingsql.UserId),
)

func (s *Storage) UserPrivacySettings(ctx context.Context, userId int64) ([]usermodel.PrivacySetting, error) {
	privacySettings := make([]usermodel.PrivacySetting, 0)
	err := s.db.QueryRowContext(ctx, userPrivacySettingsQuery, userId).Scan(pq.Array(&privacySettings))
	if err != nil {
		return privacySettings, Error(err, exception.NewCause("get user privacy settigns", "UserPrivacySettings", _PROVIDER))
	}
	// need for json marshalling!!! if we remove this, json.Marshal(list) was equal "null"
	if len(privacySettings) == 0 {
		privacySettings = make([]usermodel.PrivacySetting, 0)
	}
	return privacySettings, nil
}

var addUserPrivacySettingQuery = fmt.Sprintf(
	`
	INSERT INTO %s (%s,%s) VALUES (
		$1,
		(SELECT %s FROM %s WHERE %s = $2)
	)
	`,
	userprivacysettingsql.Table,
	userprivacysettingsql.UserId,
	userprivacysettingsql.SettingId,

	privacysettingsql.ID,
	privacysettingsql.Table,
	privacysettingsql.Type,
)

func (s *Storage) AddUserPrivacySetting(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	_, err := s.db.ExecContext(ctx, addUserPrivacySettingQuery, userId, setting)
	if err != nil {
		return Error(err, exception.NewCause("add user privacy setting", "AddUserPrivacySetting", _PROVIDER))
	}
	return nil
}

var removeUserPrivacySettingQuery = fmt.Sprintf(
	`
	DELETE FROM %s WHERE %s = $1 AND %s = (SELECT %s FROM %s WHERE %s = $2)
	`,
	userprivacysettingsql.Table,
	userprivacysettingsql.UserId,
	userprivacysettingsql.SettingId,

	privacysettingsql.ID,
	privacysettingsql.Table,
	privacysettingsql.Type,
)

func (s *Storage) RemoveUserPrivacySetting(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	r, err := s.db.ExecContext(ctx, removeUserPrivacySettingQuery, userId, setting)
	if err != nil {
		return Error(err, exception.NewCause("remove user privacy setting", "RemoveUserPrivacySetting", _PROVIDER))
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return Error(err, exception.NewCause("get rows affected amount", "RemoveUserPrivacySetting", _PROVIDER))
	}
	if rows == 0 {
		return exception.Wrap(userprivacysettingerror.ExceptionUserPrivacySettingNotFound(), exception.NewCause("check zero rows", "RemoveUserPrivacySetting", _PROVIDER))
	}
	return nil
}
