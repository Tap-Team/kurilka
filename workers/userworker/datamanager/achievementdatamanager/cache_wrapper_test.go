package achievementdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"slices"

	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/achievementdatamanager"
	gomock "github.com/golang/mock/gomock"
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

		wrapper.ReachAchievements(ctx, userId, time.Now(), achievementIds)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 12}
		reachDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewReachAchievementsMatcher(achievementIds, reachDate)).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, reachDate, achievementIds)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 12}
		reachDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewReachAchievementsMatcher(achievementIds, reachDate)).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, reachDate, achievementIds)
	}

	{
		userId := rand.Int63()
		achievementIds := []int64{1, 2, 7, 9, 13, 42, 50}
		reachDate := time.Now()
		achievements := generateRandomAchievementList(50)

		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(achievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, NewReachAchievementsMatcher(achievementIds, reachDate)).Return(errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserAchievements(gomock.Any(), userId).Return(nil).Times(1)

		wrapper.ReachAchievements(ctx, userId, reachDate, achievementIds)
	}
}
