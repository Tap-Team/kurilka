package achievementdatamanager

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"golang.org/x/exp/slog"
)

//go:generate mockgen -source cache_wrapper.go -destination cache_wrapper_mocks.go -package achievementdatamanager

const _CACHE_PROVIDER = "workers/userworker/datamanager/achievementdatamanager.cacheWrapper"

type Cache interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	RemoveUserAchievements(ctx context.Context, userId int64) error
	SaveUserAchievements(ctx context.Context, userId int64, achievements []*achievementmodel.Achievement) error
}

type cacheWrapper struct {
	Cache
}

func NewCacheWrapper(cache Cache) AchievementCache {
	return &cacheWrapper{cache}
}

type AchievementCache interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	SaveUserAchievements(ctx context.Context, userId int64, achievements []*achievementmodel.Achievement) error
	ReachAchievements(ctx context.Context, userId int64, reachDate time.Time, achievementsIds []int64)
}

func (cw *cacheWrapper) ReachAchievements(ctx context.Context, userId int64, reachDate time.Time, achievementsIds []int64) {
	achievements, err := cw.UserAchievements(ctx, userId)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "ReachAchievements", _CACHE_PROVIDER)).Error())
		cw.RemoveUserAchievements(ctx, userId)
		return
	}
	achievementsIdsSet := make(map[int64]struct{}, len(achievementsIds))
	for _, id := range achievementsIds {
		achievementsIdsSet[id] = struct{}{}
	}
	for i, ach := range achievements {
		_, ok := achievementsIdsSet[ach.ID]
		if ok {
			achievements[i].SetReachDate(reachDate)
		}
	}
	err = cw.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		cw.RemoveUserAchievements(ctx, userId)
	}
}
