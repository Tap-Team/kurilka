package userstorage

import (
	"context"
	"errors"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
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
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

func (s *Storage) User(ctx context.Context, userId int64) (*model.UserData, error) {
	var user usermodel.UserData
	var userData model.UserData
	err := s.redis.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
	if err != nil {
		return nil, Error(err, exception.NewCause("get user by id", "User", _PROVIDER))
	}
	userData.AbstinenceTime = user.AbstinenceTime.Time
	userData.CigaretteDayAmount = user.CigaretteDayAmount
	userData.CigarettePackAmount = user.CigarettePackAmount
	userData.PackPrice = user.PackPrice
	return &userData, nil
}

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
