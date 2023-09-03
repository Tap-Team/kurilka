package achievementdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_Achievement_Manager_AchievementPreview(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockAchievementCache(ctrl)
	storage := achievementdatamanager.NewMockAchievementStorage(ctrl)

	manager := achievementdatamanager.NewAchievementManager(storage, cache)

	// case then manager should return cache response
	{

		const cacheAchievementSize = 100
		cacheAchievements := randomAchievementList(cacheAchievementSize)
		userId := rand.Int63()
		cache.EXPECT().AchievementPreview(gomock.Any(), userId).Return(cacheAchievements, nil).Times(1)
		managerAchievements := manager.AchievementPreview(ctx, userId)
		require.Equal(t, cacheAchievements, managerAchievements, "manager achievements not equal, cache case")
	}

	// case then manager return achievement preview
	{
		const storageAchievementSize = 100
		storageAchievements := randomAchievementList(storageAchievementSize)
		userId := rand.Int63()
		cache.EXPECT().AchievementPreview(gomock.Any(), userId).Return(nil, errors.New("some error from cache")).Times(1)
		storage.EXPECT().AchievementPreview(gomock.Any(), userId).Return(storageAchievements).Times(1)
		managerAchievements := manager.AchievementPreview(ctx, userId)
		require.Equal(t, storageAchievements, managerAchievements, "manager achievements not equal, storage case")
	}
}

func TestAchievementClear(t *testing.T) {
	ctrl := gomock.NewController(t)

	cache := achievementdatamanager.NewMockAchievementCache(ctrl)

	datamanager := achievementdatamanager.NewAchievementManager(nil, cache)

	cache.EXPECT().Delete(gomock.Any(), gomock.Any())
	datamanager.Clear(context.Background(), rand.Int63())
}
