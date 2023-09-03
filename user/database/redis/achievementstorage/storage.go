package achievementstorage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Tap-Team/kurilka/internal/errorutils/userachievementerror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/redis/go-redis/v9"
)

const _PROVIDER = "user/database/redis/achievementstorage"

type Storage struct {
	redis *redis.Client
}

func New(redisClient *redis.Client) *Storage {
	return &Storage{redis: redisClient}
}

func Error(err error, cause exception.Cause) error {
	switch {
	case errors.Is(err, redis.Nil):
		return exception.Wrap(userachievementerror.ExceptionAchievementNotFound(), cause)
	}
	return exception.Wrap(err, cause)
}

type achievementList []*achievementmodel.Achievement

func (a achievementList) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *achievementList) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

func filterMaxLevelFromAchievementList(achievements []*achievementmodel.Achievement) []*usermodel.Achievement {
	maxLevelAchievements := make(map[achievementmodel.AchievementType]*achievementmodel.Achievement)

	for _, ach := range achievements {
		currentMaxLevel := maxLevelAchievements[ach.Type].Level
		if ach.Level > currentMaxLevel {
			maxLevelAchievements[ach.Type] = ach
		}
	}
	achievementPreview := make([]*usermodel.Achievement, 0)
	for _, ach := range maxLevelAchievements {
		achievement := usermodel.NewAсhievement(ach.Type, ach.Level)
		achievementPreview = append(achievementPreview, &achievement)
	}
	return achievementPreview
}

func (s *Storage) AchievementPreview(ctx context.Context, userId int64) ([]*usermodel.Achievement, error) {
	achievementList := make(achievementList, 0)
	err := s.redis.Get(ctx, redishelper.AchievementsKey(userId)).Scan(&achievementList)
	if err != nil {
		return nil, Error(err, exception.NewCause("get achievements", "AchievementPreview", _PROVIDER))
	}
	return filterMaxLevelFromAchievementList(achievementList), nil
}

func (s *Storage) Delete(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.AchievementsKey(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("delete achievements", "Delete", _PROVIDER))
	}
	return nil
}
