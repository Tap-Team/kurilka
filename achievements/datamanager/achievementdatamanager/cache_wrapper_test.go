package achievementdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_Wrapper_UserAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockCache(ctrl)

	wrapper := achievementdatamanager.NewCacheWrapper(cache)

	{
		userId := rand.Int63()
		achievements := generateRandomAchievementList(50)
		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)

		userAchievements, err := wrapper.UserAchievements(ctx, userId)
		assert.NilError(t, err, "non nil error")

		equal := slices.EqualFunc(achievements, userAchievements, compareAchievements)
		assert.Equal(t, true, equal, "achievements not equal")
	}

	{
		userId := rand.Int63()
		expectedErr := errors.New("any error")
		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		userAchievements, err := wrapper.UserAchievements(ctx, userId)

		assert.Equal(t, 0, len(userAchievements), "achievements non zero")
		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}
}

func Test_Wrapper_SaveUserAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockCache(ctrl)

	wrapper := achievementdatamanager.NewCacheWrapper(cache)

	{
		userId := rand.Int63()
		achievements := generateRandomAchievementList(50)

		expectedErr := errors.New("any error")

		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, achievements).Return(expectedErr).Times(1)

		err := wrapper.SaveUserAchievements(ctx, userId, achievements)
		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}
}

func Test_Wrapper_ReachAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockCache(ctrl)

	wrapper := achievementdatamanager.NewCacheWrapper(cache)

	{
		userId := rand.Int63()
		achievementIds := []int64{}

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, errors.New("any error"))
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, amidtime.Timestamp{Time: time.Now()}, achievementIds)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 12}
		reachDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewReachAchievementsMatcher(achievementIds, reachDate)).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, amidtime.Timestamp{Time: reachDate}, achievementIds)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 12}
		reachDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewReachAchievementsMatcher(achievementIds, reachDate)).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, amidtime.Timestamp{Time: reachDate}, achievementIds)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 13, 42, 50}
		reachDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewReachAchievementsMatcher(achievementIds, reachDate)).Return(errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, amidtime.Timestamp{Time: reachDate}, achievementIds)
	}
}

func Test_Wrapper_OpenAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockCache(ctrl)

	wrapper := achievementdatamanager.NewCacheWrapper(cache)

	{
		userId := rand.Int63()
		achievementsIds := []int64{}

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.OpenAchievements(ctx, userId, achievementsIds, time.Now())
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 12}
		openDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(ctx, userId, NewOpenAchievementsMatcher(achievementIds, openDate)).Return(nil).Times(1)

		wrapper.OpenAchievements(ctx, userId, achievementIds, openDate)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 12}
		openDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewOpenAchievementsMatcher(achievementIds, openDate)).Return(errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.OpenAchievements(ctx, userId, achievementIds, openDate)
	}
}

type achList []*achievementmodel.Achievement

func (a achList) Ids() []int64 {
	ids := make([]int64, 0, len(a))
	for _, ach := range a {
		ids = append(ids, ach.ID)
	}
	return ids
}

func Test_Wrapper_MarkShown(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockCache(ctrl)

	wrapper := achievementdatamanager.NewCacheWrapper(cache)

	reachOpt := func(ach *achievementmodel.Achievement) {
		if rand.Intn(3)%2 == 0 {
			ach.ReachDate = amidtime.Timestamp{Time: time.Now()}
		}
	}

	{
		userId := rand.Int63()

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.MarkShown(ctx, userId)
	}

	{
		userId := rand.Int63()

		userAchievements := generateRandomAchievementList(50, reachOpt)
		achievementIds := achList(userAchievements).Ids()

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(userAchievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewShownAchievementsMatcher(achievementIds)).Return(nil).Times(1)

		wrapper.MarkShown(ctx, userId)
	}

	{
		userId := rand.Int63()

		userAchievements := generateRandomAchievementList(50, reachOpt)
		achievementIds := achList(userAchievements).Ids()

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(userAchievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewShownAchievementsMatcher(achievementIds)).Return(errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.MarkShown(ctx, userId)
	}
}
