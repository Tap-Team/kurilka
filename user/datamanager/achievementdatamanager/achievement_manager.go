package achievementdatamanager

import (
	"context"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
)

type AchievementStorage interface {
	AchievementPreview(ctx context.Context, userId int64) []*usermodel.Achievement
}

type AchievementCache interface {
	AchievementPreview(ctx context.Context, userId int64) ([]*usermodel.Achievement, error)
	Delete(ctx context.Context, userId int64) error
}

type AchievementManager interface {
	AchievementPreview(ctx context.Context, userId int64) []*usermodel.Achievement
	Clear(ctx context.Context, userId int64)
}

type achievementDataManager struct {
	storage AchievementStorage
	cache   AchievementCache
}

func NewAchievementManager(
	storage AchievementStorage,
	cache AchievementCache,
) AchievementManager {
	return &achievementDataManager{
		storage: storage,
		cache:   cache,
	}
}

func (a *achievementDataManager) AchievementPreview(ctx context.Context, userId int64) []*usermodel.Achievement {
	achievements, err := a.cache.AchievementPreview(ctx, userId)
	if err == nil {
		return achievements
	}
	return a.storage.AchievementPreview(ctx, userId)
}

func (a *achievementDataManager) Clear(ctx context.Context, userId int64) {
	a.cache.Delete(ctx, userId)
}
