package userstorage

import (
	"context"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/redis/go-redis/v9"
)

const _PROVIDER = "achievements/database/redis/userstorage"

type Storage struct {
	redis *redis.Client
}

func New(rc *redis.Client) *Storage {
	return &Storage{redis: rc}
}

func Error(err error, cause exception.Cause) error {
	return exception.Wrap(err, cause)
}

func (s *Storage) User(ctx context.Context, userId int64) (*model.UserData, error) {}

// func (s *Storage) UpdateUserLevel(ctx context.Context, userId int64, level usermodel.LevelInfo) error {
// 	user := usermodel.User{}
// 	err := s.redis.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
// 	if err != nil {
// 		return Error(err, exception.NewCause("get user from storage", "UpdateUserLevel", _PROVIDER))
// 	}
// 	user.Level = level
// 	err = s.redis.Set(ctx, redishelper.UsersKey(userId), user, 0).Err()
// 	if err != nil {
// 		s.redis.Del(ctx, redishelper.UsersKey(userId))
// 		return exception.Wrap(err, exception.NewCause("set user level", "UpdateUserLevel", _PROVIDER))
// 	}
// 	return nil
// }
