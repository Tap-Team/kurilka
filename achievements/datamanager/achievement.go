package datamanager

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

type AchievementCache interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	RemoveUserAchievements(ctx context.Context, userId int64) error
	SaveUserAchievements(ctx context.Context, userId int64, achievements []*achievementmodel.Achievement) error
}

type AchievementStorage interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	MarkShown(ctx context.Context, userId int64) error
	OpenSingle(ctx context.Context, userId int64, ach model.OpenAchievement) error
	InsertUserAchievements(ctx context.Context, userId int64, reachTime amidtime.Timestamp, achievementsIds []int64) error
}

type AchievementDataManager interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error)
	MarkShown(ctx context.Context, userId int64) error
	ReachAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) error
}

type achievementDataManager struct {
	storage AchievementStorage
	cache   AchievementCache
}

func NewAchievementManager(storage AchievementStorage, cache AchievementCache) AchievementDataManager {
	return &achievementDataManager{
		storage: storage,
		cache:   cache,
	}
}

func (dm *achievementDataManager) markShownCache(ctx context.Context, userId int64) {
	achievements, err := dm.cache.UserAchievements(ctx, userId)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "markShownCache", _PROVIDER)).Error())
		dm.cache.RemoveUserAchievements(ctx, userId)
	}
	for i := range achievements {
		achievements[i].Shown = true
	}
	err = dm.cache.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user achievements in storage", "markShownCache", _PROVIDER)).Error())
		dm.cache.RemoveUserAchievements(ctx, userId)
	}
}

func (dm *achievementDataManager) setAchievementsOpen(ctx context.Context, userId int64, achievementIds []int64, openTime time.Time) {
	achievements, err := dm.cache.UserAchievements(ctx, userId)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("get user achievements from cache", "markShownCache", _PROVIDER)).Error())
		dm.cache.RemoveUserAchievements(ctx, userId)
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
	err = dm.cache.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user achievements in storage", "markShownCache", _PROVIDER)).Error())
		dm.cache.RemoveUserAchievements(ctx, userId)
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

func (dm *achievementDataManager) reachAchievementsCache(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) {

	achievements, err := dm.cache.UserAchievements(ctx, userId)
	if err != nil {
		dm.cache.RemoveUserAchievements(ctx, userId)
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
	err = dm.cache.SaveUserAchievements(ctx, userId, achievements)
	if err != nil {
		dm.cache.RemoveUserAchievements(ctx, userId)
	}
}

func (dm *achievementDataManager) ReachAchievements(ctx context.Context, userId int64, reachDate amidtime.Timestamp, achievementsIds []int64) error {
	err := dm.storage.InsertUserAchievements(ctx, userId, reachDate, achievementsIds)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("add achievements to user", "AddAchievements", _PROVIDER))
	}
	dm.reachAchievementsCache(ctx, userId, reachDate, achievementsIds)
	return nil
}

func (dm *achievementDataManager) OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error) {
	openTime := time.Now()
	err := dm.storage.OpenSingle(ctx, userId, model.NewOpenAchievement(achievementId, openTime))
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("open single in storage", "OpenSingle", _PROVIDER))
	}
	dm.setAchievementsOpen(ctx, userId, []int64{achievementId}, openTime)
	return model.NewOpenAchievementResponse(openTime), nil
}

func (dm *achievementDataManager) MarkShown(ctx context.Context, userId int64) error {
	err := dm.storage.MarkShown(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("mark shown in storage", "MarkShown", _PROVIDER))
	}
	dm.markShownCache(ctx, userId)
	return nil
}
