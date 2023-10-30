package votesubscriptionstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tap-Team/kurilka/internal/sqlmodel/votesubscriptionsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/lib/pq"
)

const _PROVIDER = "vote/storage/postgres/votesubscriptionstorage"

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
		case votesubscriptionsql.PrimaryKey:
			return subscriptionerror.SubscriptionIdExists
		case votesubscriptionsql.SubscriptionUnique:
			return subscriptionerror.SubscriptionIdExists
		}
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return subscriptionerror.SubscriptionNotFound
	}
	return err
}

var insert_subscription_query = new(sqlutils.QueryBuilder).
	InsertInto(votesubscriptionsql.Table, votesubscriptionsql.SubscriptionId, votesubscriptionsql.UserId).
	WriteQuery("VALUES ($1, $2)").
	Build()

func (s *Storage) CreateSubscription(ctx context.Context, subscriptionId, userId int64) error {
	_, err := s.db.ExecContext(ctx, insert_subscription_query, subscriptionId, userId)
	if err != nil {
		return Error(err, exception.NewCause("create subscription", "CreateSubscription", _PROVIDER))
	}
	return nil
}

var delete_subscription_query = new(sqlutils.QueryBuilder).
	DeleteFrom(votesubscriptionsql.Table).
	WhereColumnEqual(votesubscriptionsql.SubscriptionId, "$1").
	Build()

func (s *Storage) DeleteSubscription(ctx context.Context, subscriptionId int64) error {
	_, err := s.db.ExecContext(ctx, delete_subscription_query, subscriptionId)
	if err != nil {
		return Error(err, exception.NewCause("delete subscription", "DeleteSubscription", _PROVIDER))
	}
	return nil
}

var update_subscription_id_query = new(sqlutils.QueryBuilder).
	Update(votesubscriptionsql.Table).
	SetColumn(votesubscriptionsql.SubscriptionId, "$2").
	WhereColumnEqual(votesubscriptionsql.UserId, "$1").
	Build()

func (s *Storage) UpdateUserSubscriptionId(ctx context.Context, userId, subscriptionId int64) error {
	_, err := s.db.ExecContext(ctx, update_subscription_id_query, userId, subscriptionId)
	if err != nil {
		return Error(err, exception.NewCause("update user subscription", "UpdateUserSubscriptionId", _PROVIDER))
	}
	return nil
}

var user_subscription_id_query = new(sqlutils.QueryBuilder).
	Select(votesubscriptionsql.SubscriptionId).
	From(votesubscriptionsql.Table).
	WhereColumnEqual(votesubscriptionsql.UserId, "$1").
	Build()

func (s *Storage) UserSubscriptionId(ctx context.Context, userId int64) (subscriptionId int64, err error) {
	err = s.db.QueryRowContext(ctx, user_subscription_id_query, userId).Scan(&subscriptionId)
	if err != nil {
		err = Error(err, exception.NewCause("get user subscription", "UserSubscriptionId", _PROVIDER))
		return
	}
	return
}
