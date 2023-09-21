package achievementdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"slices"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/achievementdatamanager"
	gomock "github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	NilOpenResponse *model.OpenAchievementResponse
)

func Test_Manager_UserAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := achievementdatamanager.NewMockAchievementStorage(ctrl)
	cache := achievementdatamanager.NewMockAchievementCache(ctrl)

	manager := achievementdatamanager.New(storage, cache)

	{
		userId := rand.Int63()
		userAchievements := generateRandomAchievementList(50)
		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(userAchievements, nil).Times(1)

		achievements, err := manager.UserAchievements(ctx, userId)

		assert.NilError(t, err, "non nil error")

		equal := slices.EqualFunc(achievements, userAchievements, compareAchievements)
		assert.Equal(t, true, equal)
	}

	{
		userId := rand.Int63()
		userAchievements := generateRandomAchievementList(50)
		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		storage.EXPECT().UserAchievements(gomock.Any(), userId).Return(userAchievements, nil).Times(1)
		cache.EXPECT().SaveUserAchievements(gomock.Any(), userId, userAchievements).Return(nil).Times(1)

		achievements, err := manager.UserAchievements(ctx, userId)

		assert.NilError(t, err, "non nil error")

		equal := slices.EqualFunc(achievements, userAchievements, compareAchievements)
		assert.Equal(t, true, equal)
	}

	{
		userId := rand.Int63()
		expectedErr := errors.New("database fall")
		cache.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		storage.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		achievements, err := manager.UserAchievements(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")

		assert.Equal(t, 0, len(achievements))
	}

}

func Test_Manager_ReachAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := achievementdatamanager.NewMockAchievementStorage(ctrl)
	cache := achievementdatamanager.NewMockAchievementCache(ctrl)

	manager := achievementdatamanager.New(storage, cache)

	{
		userId := rand.Int63()
		reachDate := time.Now()
		achievementIds := []int64{1, 2, 3, 4, 5, 6, 7}
		expectedErr := usererror.ExceptionUserNotFound()

		storage.EXPECT().InsertUserAchievements(gomock.Any(), userId, reachDate, achievementIds).Return(expectedErr).Times(1)

		err := manager.ReachAchievements(ctx, userId, reachDate, achievementIds)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	{
		userId := rand.Int63()
		reachDate := time.Now()
		achievementIds := []int64{1, 2, 3, 4, 5, 6, 7}
		storage.EXPECT().InsertUserAchievements(gomock.Any(), userId, reachDate, achievementIds).Return(nil).Times(1)
		cache.EXPECT().ReachAchievements(gomock.Any(), userId, reachDate, achievementIds).Times(1)

		err := manager.ReachAchievements(ctx, userId, reachDate, achievementIds)

		assert.NilError(t, err, "non nil error")
	}
}
