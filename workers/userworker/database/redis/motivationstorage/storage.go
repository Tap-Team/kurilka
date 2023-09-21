package motivationstorage

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

const _PROVIDER = "workers/userworker/database/redis/motivationstorage.Storage"

type Storage struct {
	redis *redis.Client
	exp   time.Duration
}

func New(rc *redis.Client, expiration time.Duration) *Storage {
	return &Storage{redis: rc, exp: expiration}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

func (s *Storage) SaveUserMotivation(ctx context.Context, userId int64, motivation string) error {
	var user usermodel.UserData
	err := s.redis.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
	if err != nil {
		return Error(err, exception.NewCause("get user data from redis", "SaveUserWelcomeMotivation", _PROVIDER))
	}
	user.Motivation = motivation
	err = s.redis.Set(ctx, redishelper.UsersKey(userId), user, s.exp).Err()
	if err != nil {
		return Error(err, exception.NewCause("set saved user", "SaveUserWelcomeMotivation", _PROVIDER))
	}
	return nil
}

func (s *Storage) RemoveUserMotivation(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.UsersKey(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("delete user by id", "RemoveUserWelcomeMotivation", _PROVIDER))
	}
	return nil
}
