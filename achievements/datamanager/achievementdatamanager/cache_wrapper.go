package achievementdatamanager

import (
	"context"
	"log/slog"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

type Cache interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	RemoveUserAchievements(ctx context.Context, userId int64) error
	SaveUserAchievements(ctx context.Context, userId int64, achievements []*achievementmodel.Achievement) error
}

type cacheWrapper struct {
	Cache
}

type AchievementCache interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	SaveUserAchievements(ctx context.Context, userId int64, achievements []*achievementmodel.Achievement) error
	ReachAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64)
	OpenAchievements(ctx context.Context, userId int64, achievementIds []int64, openTime time.Time)
	MarkShown(ctx context.Context, userId int64)
}

func NewCacheWrapper(cache Cache) AchievementCache {
	return &cacheWrapper{cache}
}

func (cw *cacheWrapper) ReachAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) {
	achievements, err := cw.UserAchievements(ctx, userId)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "ReachAchievements", _PROVIDER)).Error())
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
			achievements[i].ReachDate = reachDate
		}
	}
	err = cw.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		cw.RemoveUserAchievements(ctx, userId)
	}
}

func (cw *cacheWrapper) OpenAchievements(ctx context.Context, userId int64, achievementIds []int64, openTime time.Time) {
	achievements, err := cw.UserAchievements(ctx, userId)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "OpenAchievements", _PROVIDER)).Error())
		cw.RemoveUserAchievements(ctx, userId)
		return
	}
	ids := make(map[int64]struct{}, len(achievementIds))
	for _, id := range achievementIds {
		ids[id] = struct{}{}
	}
	for i := range achievements {
		if _, ok := ids[achievements[i].ID]; ok {
			achievements[i].OpenDate = amidtime.Timestamp{Time: openTime}
		}
	}
	err = cw.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user achievements in storage", "OpenAchievements", _PROVIDER)).Error())
		cw.RemoveUserAchievements(ctx, userId)
	}
}

func (cw *cacheWrapper) MarkShown(ctx context.Context, userId int64) {
	achievements, err := cw.UserAchievements(ctx, userId)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "MarkShown", _PROVIDER)).Error())
		cw.RemoveUserAchievements(ctx, userId)
		return
	}
	for i := range achievements {
		achievements[i].Shown = achievements[i].Reached()
	}
	err = cw.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user achievements in storage", "MarkShown", _PROVIDER)).Error())
		cw.RemoveUserAchievements(ctx, userId)
	}
}
