package subscriptiondatamanager

import (
	"context"
	"log/slog"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source subscription.go -destination subscription_mocks.go -package subscriptiondatamanager

const _PROVIDER = "user/datamanager/subscriptiondatamanager"

type SubscriptionStorage interface {
	UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error)
	UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error
}

type SubscriptionCache interface {
	UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error)
	UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error
	RemoveUserSubscription(ctx context.Context, userId int64) error
}

type SubscriptionManager interface {
	UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error)
	UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error
}

type subscriptionManager struct {
	storage SubscriptionStorage
	cache   SubscriptionCache
}

func New(storage SubscriptionStorage, cache SubscriptionCache) SubscriptionManager {
	return &subscriptionManager{
		storage: storage,
		cache:   cache,
	}
}

func (s *subscriptionManager) UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error) {
	subscription, err := s.cache.UserSubscription(ctx, userId)
	if err == nil {
		return subscription, nil
	}
	subscription, err = s.storage.UserSubscription(ctx, userId)
	if err != nil {
		return usermodel.Subscription{}, exception.Wrap(err, exception.NewCause("get user subscription", "UserSubscription", _PROVIDER))
	}
	err = s.cache.UpdateUserSubscription(ctx, userId, subscription)
	if err != nil {
		slog.ErrorContext(ctx, "failed update user subscription", "err", err)
		s.cache.RemoveUserSubscription(ctx, userId)
	}
	return subscription, nil
}

func (s *subscriptionManager) UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	err := s.storage.UpdateUserSubscription(ctx, userId, subscription)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("update subscription in storage", "UpdateUserSubscription", _PROVIDER))
	}
	err = s.cache.UpdateUserSubscription(ctx, userId, subscription)
	if err != nil {
		slog.ErrorContext(ctx, "failed update user subscription", "err", err)
		s.cache.RemoveUserSubscription(ctx, userId)
	}
	return nil
}
