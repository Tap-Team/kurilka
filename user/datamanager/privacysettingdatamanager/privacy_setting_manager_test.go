package privacysettingdatamanager_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/golang/mock/gomock"
	"golang.org/x/exp/slices"
	"gotest.tools/v3/assert"
)

type PrivacySettingsMatcher struct {
	privacySettings []usermodel.PrivacySetting
}

func NewPrivacySettingsMatcher(settings []usermodel.PrivacySetting) gomock.Matcher {
	return &PrivacySettingsMatcher{privacySettings: settings}
}

func (c *PrivacySettingsMatcher) Matches(x interface{}) bool {

	arr, ok := x.([]usermodel.PrivacySetting)
	if !ok {
		return false
	}
	ok = slices.Equal(
		arr,
		c.privacySettings,
	)
	return ok
}

func (c *PrivacySettingsMatcher) String() string {
	return fmt.Sprintf("is equal to privacy settings list %v", c.privacySettings)
}

func Test_PrivacySettings_CacheWrapper_AddSingle(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := privacysettingdatamanager.NewMockPrivacySettingCache(ctrl)

	cacheWrapper := privacysettingdatamanager.CacheWrapper{cache}

	// case when we get privacy settings and successfully save it in cache
	{
		userId := rand.Int63()
		cachePrivacySettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.ACHIEVEMENTS_DURATION,
			usermodel.ACHIEVEMENTS_SAVING,
		}
		setting := usermodel.ACHIEVEMENTS_WELL_BEING

		saveList := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.ACHIEVEMENTS_DURATION,
			usermodel.ACHIEVEMENTS_SAVING,
			usermodel.ACHIEVEMENTS_WELL_BEING,
		}

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(cachePrivacySettings, nil).Times(1)
		cache.EXPECT().SaveUserPrivacySettings(ctx, userId, NewPrivacySettingsMatcher(saveList)).Return(nil).Times(1)

		cacheWrapper.AddSingle(ctx, userId, setting)
	}

	// case when user privacy settings return unknown error
	{
		userId := rand.Int63()
		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("unknown error")).Times(1)

		cacheWrapper.AddSingle(ctx, userId, usermodel.ACHIEVEMENTS_HEALTH)
	}

	// case when save in cache cause error
	{
		userId := rand.Int63()
		setting := usermodel.STATISTICS_MONEY

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return([]usermodel.PrivacySetting{}, nil).Times(1)
		cache.EXPECT().SaveUserPrivacySettings(ctx, userId, NewPrivacySettingsMatcher([]usermodel.PrivacySetting{usermodel.STATISTICS_MONEY})).Return(errors.New("unknown error")).Times(1)
		cache.EXPECT().RemoveUserPrivacySettings(ctx, userId).Return(nil).Times(1)
		cacheWrapper.AddSingle(ctx, userId, setting)
	}
}

func Test_PrivacySettings_CacheWrapper_RemoveSingle(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := privacysettingdatamanager.NewMockPrivacySettingCache(ctrl)

	cacheWrapper := privacysettingdatamanager.CacheWrapper{cache}

	{
		userId := rand.Int63()
		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return([]usermodel.PrivacySetting{}, nil).Times(1)

		cacheWrapper.RemoveSingle(ctx, userId, usermodel.ACHIEVEMENTS_HEALTH)
	}

	{
		userId := rand.Int63()
		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)

		cacheWrapper.RemoveSingle(ctx, userId, usermodel.ACHIEVEMENTS_CIGARETTE)
	}

	// case when remove non exists achievement
	{
		privacySettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.STATISTICS_TIME,
		}

		userId := rand.Int63()

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(privacySettings, nil).Times(1)

		cacheWrapper.RemoveSingle(ctx, userId, usermodel.STATISTICS_MONEY)
	}

	// case when all ok and remove single privacySetting from list
	{
		privacySettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.STATISTICS_MONEY,
		}

		removeSetting := privacySettings[0]

		expectedSaveSettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.STATISTICS_MONEY,
		}

		userId := rand.Int63()

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(privacySettings, nil).Times(1)
		cache.EXPECT().SaveUserPrivacySettings(gomock.Any(), userId, NewPrivacySettingsMatcher(expectedSaveSettings)).Return(nil).Times(1)

		cacheWrapper.RemoveSingle(ctx, userId, removeSetting)
	}

	// case when we save return err and we should remove cache
	{
		privacySettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.STATISTICS_MONEY,
		}

		removeSetting := privacySettings[0]

		expectedSaveSettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.STATISTICS_MONEY,
		}

		userId := rand.Int63()

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(privacySettings, nil).Times(1)
		cache.EXPECT().SaveUserPrivacySettings(gomock.Any(), userId, NewPrivacySettingsMatcher(expectedSaveSettings)).Return(errors.New("any err")).Times(1)
		cache.EXPECT().RemoveUserPrivacySettings(gomock.Any(), userId).Return(nil).Times(1)

		cacheWrapper.RemoveSingle(ctx, userId, removeSetting)
	}
}

func Test_PrivacySettings_Manager_Add(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := privacysettingdatamanager.NewMockPrivacySettingCache(ctrl)
	storage := privacysettingdatamanager.NewMockPrivacySettingStorage(ctrl)

	manager := privacysettingdatamanager.NewPrivacyManager(storage, cache)

	{
		userId := rand.Int63()

		setting := usermodel.ACHIEVEMENTS_CIGARETTE

		storage.EXPECT().AddUserPrivacySetting(gomock.Any(), userId, setting).Return(errors.New("any error")).Times(1)

		manager.Add(ctx, userId, setting)
	}

	{
		userId := rand.Int63()

		setting := usermodel.ACHIEVEMENTS_WELL_BEING

		storage.EXPECT().AddUserPrivacySetting(gomock.Any(), userId, setting).Return(nil).Times(1)
		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("cache umer")).Times(1)

		manager.Add(ctx, userId, setting)
	}
}

func Test_PrivacySetitngs_UseCase_Remove(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := privacysettingdatamanager.NewMockPrivacySettingCache(ctrl)
	storage := privacysettingdatamanager.NewMockPrivacySettingStorage(ctrl)

	manager := privacysettingdatamanager.NewPrivacyManager(storage, cache)

	{
		userId := rand.Int63()

		setting := usermodel.STATISTICS_CIGARETTE

		expectedErr := errors.New("server umer")

		storage.EXPECT().RemoveUserPrivacySetting(gomock.Any(), userId, setting).Return(expectedErr).Times(1)

		err := manager.Remove(ctx, userId, setting)

		assert.ErrorIs(t, err, expectedErr, "wrong error from remove")
	}

	{
		userId := rand.Int63()

		setting := usermodel.STATISTICS_TIME

		storage.EXPECT().RemoveUserPrivacySetting(gomock.Any(), userId, setting).Return(nil).Times(1)
		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("cache umer")).Times(1)

		err := manager.Remove(ctx, userId, setting)
		assert.ErrorIs(t, err, nil, "wrong err from remove")
	}
}

func Test_PrivacySettings_UseCase_PrivacySettings(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := privacysettingdatamanager.NewMockPrivacySettingCache(ctrl)
	storage := privacysettingdatamanager.NewMockPrivacySettingStorage(ctrl)

	manager := privacysettingdatamanager.NewPrivacyManager(storage, cache)

	{
		userId := rand.Int63()

		expectedPrivacySettings := []usermodel.PrivacySetting{
			usermodel.STATISTICS_CIGARETTE,
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.STATISTICS_MONEY,
		}

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(expectedPrivacySettings, nil).Times(1)

		settings, err := manager.PrivacySettings(ctx, userId)

		equal := slices.Equal(expectedPrivacySettings, settings)
		assert.Equal(t, true, equal, "settings not equal")
		assert.ErrorIs(t, err, nil, "wrong err from privacy settings")
	}

	{
		userId := rand.Int63()

		expectedPrivacySettings := []usermodel.PrivacySetting{
			usermodel.STATISTICS_CIGARETTE,
			usermodel.STATISTICS_LIFE,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.ACHIEVEMENTS_CIGARETTE,
		}

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(expectedPrivacySettings, nil).Times(1)

		settings, err := manager.PrivacySettings(ctx, userId)

		equal := slices.Equal(expectedPrivacySettings, settings)
		assert.Equal(t, true, equal, "settings not equal")
		assert.ErrorIs(t, err, nil, "wrong err from privacy settings")
	}

	{
		userId := rand.Int63()

		expectedErr := errors.New("any error")

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("failed get data from cache")).Times(1)
		storage.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		settings, err := manager.PrivacySettings(ctx, userId)

		assert.Equal(t, 0, len(settings), "wrong settings")
		assert.ErrorIs(t, err, expectedErr, "wrong err")
	}

	{
		userId := rand.Int63()

		expectedPrivacySettings := []usermodel.PrivacySetting{
			usermodel.STATISTICS_CIGARETTE,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_DURATION,
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.STATISTICS_LIFE,
		}

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("failed get data from cache")).Times(1)
		storage.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(expectedPrivacySettings, nil).Times(1)
		cache.EXPECT().SaveUserPrivacySettings(gomock.Any(), userId, expectedPrivacySettings).Return(nil).Times(1)

		settings, err := manager.PrivacySettings(ctx, userId)

		equal := slices.Equal(expectedPrivacySettings, settings)

		assert.Equal(t, true, equal, "settings not equal")
		assert.ErrorIs(t, err, nil, "wrong error")
	}

	{
		userId := rand.Int63()

		expectedPrivacySettings := []usermodel.PrivacySetting{
			usermodel.STATISTICS_CIGARETTE,
			usermodel.ACHIEVEMENTS_WELL_BEING,
			usermodel.ACHIEVEMENTS_HEALTH,
			usermodel.ACHIEVEMENTS_DURATION,
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.STATISTICS_LIFE,
		}

		cache.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(nil, errors.New("failed get data from cache")).Times(1)
		storage.EXPECT().UserPrivacySettings(gomock.Any(), userId).Return(expectedPrivacySettings, nil).Times(1)
		cache.EXPECT().SaveUserPrivacySettings(gomock.Any(), userId, expectedPrivacySettings).Return(errors.New("failed save")).Times(1)
		cache.EXPECT().RemoveUserPrivacySettings(gomock.Any(), userId).Times(1)

		settings, err := manager.PrivacySettings(ctx, userId)

		equal := slices.Equal(expectedPrivacySettings, settings)

		assert.Equal(t, true, equal, "settings not equal")
		assert.ErrorIs(t, err, nil, "wrong error")
	}
}

func Test_PrivacySettings_UseCase_Clear(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := privacysettingdatamanager.NewMockPrivacySettingCache(ctrl)
	storage := privacysettingdatamanager.NewMockPrivacySettingStorage(ctrl)

	manager := privacysettingdatamanager.NewPrivacyManager(storage, cache)

	{
		userId := rand.Int63()

		cache.EXPECT().RemoveUserPrivacySettings(gomock.Any(), userId).Return(errors.New("random err")).Times(1)

		manager.Clear(ctx, userId)
	}

}
