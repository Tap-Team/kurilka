package userstorage

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

const _PROVIDER = "user/database/redis/userstorage"

type Storage struct {
	redis      *redis.Client
	expiration time.Duration
}

func New(redis *redis.Client, exp time.Duration) *Storage {
	return &Storage{redis: redis, expiration: exp}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

func (s *Storage) SaveUser(ctx context.Context, userId int64, user *usermodel.UserData) error {
	err := s.redis.Set(ctx, redishelper.UsersKey(userId), user, s.expiration).Err()
	if err != nil {
		return Error(err, exception.NewCause("set user data in storage", "SaveUser", _PROVIDER))
	}
	return nil
}

func (s *Storage) RemoveUser(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.UsersKey(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("delete user from hash table", "RemoveUser", _PROVIDER))
	}
	return nil
}

func (s *Storage) User(ctx context.Context, userId int64) (*usermodel.UserData, error) {
	var user usermodel.UserData
	err := s.redis.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
	if err != nil {
		return &user, Error(err, exception.NewCause("get user from cache", "User", _PROVIDER))
	}
	return &user, nil
}
