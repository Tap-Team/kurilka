package triggerstorage

import (
	"context"
	"errors"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/redis/go-redis/v9"
)

const _PROVIDER = "user/database/redis/triggerstorage"

type Storage struct {
	redis      *redis.Client
	expiration time.Duration
}

func New(rc *redis.Client, exp time.Duration) *Storage {
	return &Storage{redis: rc, expiration: exp}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

func (s *Storage) UserTriggers(ctx context.Context, userId int64) ([]usermodel.Trigger, error) {
	var user usermodel.UserData
	err := s.redis.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
	if err != nil {
		return nil, Error(err, exception.NewCause("get user from storage", "UserTriggers", _PROVIDER))
	}
	return user.Triggers, nil
}

func (s *Storage) SaveUserTriggers(ctx context.Context, userId int64, triggers []usermodel.Trigger) error {
	var user usermodel.UserData
	err := s.redis.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
	if err != nil {
		return Error(err, exception.NewCause("get user from storage", "SaveUserTriggers", _PROVIDER))
	}
	user.Triggers = triggers
	err = s.redis.Set(ctx, redishelper.UsersKey(userId), user, s.expiration).Err()
	if err != nil {
		return Error(err, exception.NewCause("set user", "SaveUserTriggers", _PROVIDER))
	}
	return nil
}

func (s *Storage) RemoveUserTriggers(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.UsersKey(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("delete user triggers", "RemoveUserTriggers", _PROVIDER))
	}
	return nil
}
