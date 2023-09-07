package achievementdatamanager

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/userachievementerror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "achievements/datamanager"

type AchievementStorage interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	MarkShown(ctx context.Context, userId int64) error
	OpenSingle(ctx context.Context, userId int64, ach model.OpenAchievement) error
	InsertUserAchievements(ctx context.Context, userId int64, reachTime amidtime.Timestamp, achievementsIds []int64) error
}

type AchievementManager interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error)
	MarkShown(ctx context.Context, userId int64) error
	ReachAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) error
}

type achievementDataManager struct {
	storage AchievementStorage
	cache   AchievementCache
}

func NewAchievementManager(storage AchievementStorage, cache AchievementCache) AchievementManager {
	return &achievementDataManager{
		storage: storage,
		cache:   cache,
	}
}

func (dm *achievementDataManager) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	achievements, err := dm.cache.UserAchievements(ctx, userId)
	if err == nil {
		return achievements, nil
	}
	if !errors.Is(err, userachievementerror.ExceptionAchievementNotFound()) {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "UserAchievements", _PROVIDER)).String())
	}
	achievements, err = dm.storage.UserAchievements(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user achievements from storage", "UserAchievements", _PROVIDER))
	}
	dm.cache.SaveUserAchievements(ctx, userId, achievements)
	return achievements, nil
}

func (dm *achievementDataManager) ReachAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) error {
	err := dm.storage.InsertUserAchievements(ctx, userId, reachDate, achievementsIds)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("add achievements to user", "AddAchievements", _PROVIDER))
	}
	dm.cache.ReachAchievements(ctx, userId, reachDate, achievementsIds)
	return nil
}

func (dm *achievementDataManager) OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error) {
	openTime := time.Now()
	err := dm.storage.OpenSingle(ctx, userId, model.NewOpenAchievement(achievementId, openTime))
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("open single in storage", "OpenSingle", _PROVIDER))
	}
	dm.cache.OpenAchievements(ctx, userId, []int64{achievementId}, openTime)
	return model.NewOpenAchievementResponse(openTime), nil
}

func (dm *achievementDataManager) MarkShown(ctx context.Context, userId int64) error {
	err := dm.storage.MarkShown(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("mark shown in storage", "MarkShown", _PROVIDER))
	}
	dm.cache.MarkShown(ctx, userId)
	return nil
}
