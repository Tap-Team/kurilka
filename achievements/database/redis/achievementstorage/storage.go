package achievementstorage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Tap-Team/kurilka/internal/errorutils/userachievementerror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/redis/go-redis/v9"
)

const _PROVIDER = "achievements/database/redis/achievementstorage"

type Storage struct {
	redis *redis.Client
}

func New(rc *redis.Client) *Storage {
	return &Storage{redis: rc}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(userachievementerror.ExceptionAchievementNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

type userAchievementList []*achievementmodel.Achievement

func (a userAchievementList) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *userAchievementList) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

func (s *Storage) SaveUserAchievements(ctx context.Context, userId int64, achievements []*achievementmodel.Achievement) error {
	err := s.redis.Set(ctx, redishelper.AchievementsKey(userId), userAchievementList(achievements), 0).Err()
	if err != nil {
		return Error(err, exception.NewCause("set achievements key", "SaveUserAchievements", _PROVIDER))
	}
	return nil
}

func (s *Storage) RemoveUserAchievements(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.AchievementsKey(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("remove user achievements", "RemoveUserAchievements", _PROVIDER))
	}
	return nil
}

func (s *Storage) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	achievementList := make(userAchievementList, 0)
	err := s.redis.Get(ctx, redishelper.AchievementsKey(userId)).Scan(&achievementList)
	if err != nil {
		return nil, Error(err, exception.NewCause("get user achievements", "UserAchievements", _PROVIDER))
	}
	return achievementList, nil
}
