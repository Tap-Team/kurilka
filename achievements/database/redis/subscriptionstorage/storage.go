package subscriptionstorage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/redis/go-redis/v9"
)

const _PROVIDER = "user/database/redis/subscriptionstorage"

type Storage struct {
	redis      *redis.Client
	expiration time.Duration
}

func New(rc *redis.Client, expiration time.Duration) *Storage {
	return &Storage{redis: rc, expiration: expiration}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

type subscriptionBinaryMarshaller usermodel.Subscription

func (s subscriptionBinaryMarshaller) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *subscriptionBinaryMarshaller) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Storage) UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error) {
	var subscription subscriptionBinaryMarshaller
	err := s.redis.Get(ctx, redishelper.SubscriptionKey(userId)).Scan(&subscription)
	if err != nil {
		return usermodel.Subscription(subscription), Error(err, exception.NewCause("get user subscription", "UserSubscription", _PROVIDER))
	}
	return usermodel.Subscription(subscription), nil
}

func (s *Storage) UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	err := s.redis.Set(ctx, redishelper.SubscriptionKey(userId), subscriptionBinaryMarshaller(subscription), s.expiration).Err()
	if err != nil {
		return Error(err, exception.NewCause("set user subscription value to key", "UpdateUserSubscription", _PROVIDER))
	}
	return nil
}

func (s *Storage) RemoveUserSubscription(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.SubscriptionKey(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("delete user subscription from cache", "RemoveUserSubscription", _PROVIDER))
	}
	return nil
}
