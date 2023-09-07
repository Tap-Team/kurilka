package achievementdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/golang/mock/gomock"
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

	manager := achievementdatamanager.NewAchievementManager(storage, cache)

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

func Test_Manager_OpenSingle(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := achievementdatamanager.NewMockAchievementStorage(ctrl)
	cache := achievementdatamanager.NewMockAchievementCache(ctrl)

	manager := achievementdatamanager.NewAchievementManager(storage, cache)

	{
		userId := rand.Int63()
		achievementId := rand.Int63n(50) + 1
		openTime := time.Now()
		expectedErr := errors.New("any error")

		storage.EXPECT().OpenSingle(gomock.Any(), userId, NewOpenAchievementMatcher(achievementId, openTime)).Return(expectedErr).Times(1)

		openResponse, err := manager.OpenSingle(ctx, userId, achievementId)

		assert.Equal(t, NilOpenResponse, openResponse, "non nil open response")
		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	{
		userId := rand.Int63()
		achievementId := rand.Int63n(50) + 1
		openTime := time.Now()

		storage.EXPECT().OpenSingle(gomock.Any(), userId, NewOpenAchievementMatcher(achievementId, openTime)).Return(nil).Times(1)
		cache.EXPECT().OpenAchievements(gomock.Any(), userId, []int64{achievementId}, TimeMatcher(openTime)).Times(1)

		openResponse, err := manager.OpenSingle(ctx, userId, achievementId)

		assert.NilError(t, err, "non nil error")
		assert.Equal(t, openResponse.OpenTime.Unix(), openTime.Unix(), "wrong time")
	}
}

func Test_Manager_MarkShown(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := achievementdatamanager.NewMockAchievementStorage(ctrl)
	cache := achievementdatamanager.NewMockAchievementCache(ctrl)

	manager := achievementdatamanager.NewAchievementManager(storage, cache)

	{
		userId := rand.Int63()
		expectedErr := usererror.ExceptionUserNotFound()

		storage.EXPECT().MarkShown(gomock.Any(), userId).Return(expectedErr).Times(1)

		err := manager.MarkShown(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	{
		userId := rand.Int63()

		storage.EXPECT().MarkShown(gomock.Any(), userId).Return(nil).Times(1)
		cache.EXPECT().MarkShown(gomock.Any(), userId).Times(1)

		err := manager.MarkShown(ctx, userId)

		assert.NilError(t, err, "non nil error")
	}
}

func Test_Manager_ReachAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := achievementdatamanager.NewMockAchievementStorage(ctrl)
	cache := achievementdatamanager.NewMockAchievementCache(ctrl)

	manager := achievementdatamanager.NewAchievementManager(storage, cache)

	{
		userId := rand.Int63()
		reachDate := amidtime.Timestamp{Time: time.Now()}
		achievementIds := []int64{1, 2, 3, 4, 5, 6, 7}
		expectedErr := usererror.ExceptionUserNotFound()

		storage.EXPECT().InsertUserAchievements(gomock.Any(), userId, reachDate, achievementIds).Return(expectedErr).Times(1)

		err := manager.ReachAchievements(ctx, userId, reachDate, achievementIds)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	{
		userId := rand.Int63()
		reachDate := amidtime.Timestamp{Time: time.Now()}
		achievementIds := []int64{1, 2, 3, 4, 5, 6, 7}

		storage.EXPECT().InsertUserAchievements(gomock.Any(), userId, reachDate, achievementIds).Return(nil).Times(1)
		cache.EXPECT().ReachAchievements(gomock.Any(), userId, reachDate, achievementIds).Times(1)

		err := manager.ReachAchievements(ctx, userId, reachDate, achievementIds)

		assert.NilError(t, err, "non nil error")
	}
}
