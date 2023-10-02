package achievementdatamanager

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/errorutils/userachievementerror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source manager.go -destination manager_mocks.go -package achievementdatamanager

const _PROVIDER = "workers/userworker/datamanager/achievementdatamanager.manager"

type AchievementStorage interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	InsertUserAchievements(ctx context.Context, userId int64, reachDate time.Time, achievementsIds []int64) error
}

type AchievementManager interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	ReachAchievements(ctx context.Context, userId int64, reachDate time.Time, achievementsIds []int64) error
}

type manager struct {
	cache   AchievementCache
	storage AchievementStorage
}

func New(storage AchievementStorage, cache AchievementCache) AchievementManager {
	return &manager{cache: cache, storage: storage}
}

func (m *manager) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	achievements, err := m.cache.UserAchievements(ctx, userId)
	if err == nil {
		return achievements, nil
	}
	if !errors.Is(err, userachievementerror.ExceptionAchievementNotFound()) {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "UserAchievements", _PROVIDER)).String())
	}
	achievements, err = m.storage.UserAchievements(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user achievements from storage", "UserAchievements", _PROVIDER))
	}
	m.cache.SaveUserAchievements(ctx, userId, achievements)
	return achievements, nil
}

func (m *manager) ReachAchievements(ctx context.Context, userId int64, reachDate time.Time, achievementsIds []int64) error {
	err := m.storage.InsertUserAchievements(ctx, userId, reachDate, achievementsIds)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("add achievements to user", "AddAchievements", _PROVIDER))
	}
	m.cache.ReachAchievements(ctx, userId, reachDate, achievementsIds)
	return nil
}
