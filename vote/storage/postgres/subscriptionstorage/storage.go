package subscriptionstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/subscriptiontypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersubscriptionsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
)

const _PROVIDER = "user/database/postgres/subscriptionstorage"

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

var (
	userSubscriptionQuery = fmt.Sprintf(
		`SELECT %s,%s FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1 GROUP BY %s,%s`,
		sqlutils.Full(usersubscriptionsql.Expired),
		sqlutils.Full(subscriptiontypesql.Type),

		usersubscriptionsql.Table,

		subscriptiontypesql.Table,
		sqlutils.Full(usersubscriptionsql.TypeId),
		sqlutils.Full(subscriptiontypesql.ID),

		sqlutils.Full(usersubscriptionsql.UserId),

		sqlutils.Full(usersubscriptionsql.Expired),
		sqlutils.Full(subscriptiontypesql.Type),
	)
)

func (s *Storage) UserSubscription(ctx context.Context, userid int64) (usermodel.Subscription, error) {
	var subscription usermodel.Subscription
	err := s.db.QueryRowContext(ctx, userSubscriptionQuery, userid).Scan(
		&subscription.Expired,
		&subscription.Type,
	)
	if err != nil {
		return subscription, Error(err, exception.NewCause("get user subscription", "UserSubscription", _PROVIDER))
	}
	return subscription, nil
}

var (
	updateSubscriptionQuery = fmt.Sprintf(
		`
		WITH type_select as (
			SELECT %s as id FROM %s WHERE %s = $3 
		)
		UPDATE %s SET %s = $2, %s = (SELECT id FROM type_select) WHERE %s = $1
		`,

		subscriptiontypesql.ID,
		subscriptiontypesql.Table,
		subscriptiontypesql.Type,

		usersubscriptionsql.Table,
		usersubscriptionsql.Expired,
		usersubscriptionsql.TypeId,

		usersubscriptionsql.UserId,
	)
)

func (s *Storage) UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	tx, err := s.db.Begin()
	if err != nil {
		return Error(err, exception.NewCause("begin tx", "UpdateUserSubscription", _PROVIDER))
	}
	defer tx.Rollback()
	r, err := tx.ExecContext(ctx, updateSubscriptionQuery, userId, subscription.Expired, subscription.Type)
	if err != nil {
		return Error(err, exception.NewCause("update user subscription query", "UpdateUserSubscription", _PROVIDER))
	}
	rows, _ := r.RowsAffected()
	if rows == 0 || rows > 1 {
		return exception.Wrap(usererror.ExceptionUserNotFound(), exception.NewCause(fmt.Sprintf("%d rows affected", rows), "UpdateUserSubscription", _PROVIDER))
	}
	err = tx.Commit()
	if err != nil {
		return Error(err, exception.NewCause("commit tx", "UpdateUserSubscription", _PROVIDER))
	}
	return nil
}
