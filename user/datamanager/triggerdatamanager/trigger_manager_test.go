package triggerdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/triggerdatamanager"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_Triggers_CacheWrapper_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	cache := triggerdatamanager.NewMockTriggerCache(ctrl)

	cacheWrapper := triggerdatamanager.CacheWrapper{cache}

	{
		userId := rand.Int63()

		cache.EXPECT().UserTriggers(gomock.Any(), userId).Return(nil, errors.New("user not found")).Times(1)

		cacheWrapper.Remove(ctx, userId, usermodel.SUPPORT_CIGGARETTE)
	}

	{
		userId := rand.Int63()

		cache.EXPECT().UserTriggers(gomock.Any(), userId).Return([]usermodel.Trigger{}, nil).Times(1)

		cacheWrapper.Remove(ctx, userId, usermodel.SUPPORT_HEALTH)
	}

	{
		userId := rand.Int63()

		triggers := []usermodel.Trigger{
			usermodel.SUPPORT_CIGGARETTE,
			usermodel.SUPPORT_HEALTH,
			usermodel.SUPPORT_TRIAL,
		}

		removeTrigger := usermodel.SUPPORT_HEALTH

		expectedTriggers := []usermodel.Trigger{
			usermodel.SUPPORT_CIGGARETTE,
			usermodel.SUPPORT_TRIAL,
		}

		cache.EXPECT().UserTriggers(gomock.Any(), userId).Return(triggers, nil).Times(1)
		cache.EXPECT().SaveUserTriggers(gomock.Any(), userId, expectedTriggers).Return(nil).Times(1)

		cacheWrapper.Remove(ctx, userId, removeTrigger)
	}

	{
		userId := rand.Int63()

		triggers := []usermodel.Trigger{
			usermodel.SUPPORT_CIGGARETTE,
			usermodel.SUPPORT_HEALTH,
			usermodel.SUPPORT_TRIAL,
		}

		removeTrigger := usermodel.SUPPORT_HEALTH

		expectedTriggers := []usermodel.Trigger{
			usermodel.SUPPORT_CIGGARETTE,
			usermodel.SUPPORT_TRIAL,
		}

		cache.EXPECT().UserTriggers(gomock.Any(), userId).Return(triggers, nil).Times(1)
		cache.EXPECT().SaveUserTriggers(gomock.Any(), userId, expectedTriggers).Return(errors.New("any err")).Times(1)
		cache.EXPECT().RemoveUserTriggers(gomock.Any(), userId).Return(nil).Times(1)
		cacheWrapper.Remove(ctx, userId, removeTrigger)
	}
}

func Test_Triggers_CacheWrapper_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	cache := triggerdatamanager.NewMockTriggerCache(ctrl)

	cacheWrapper := triggerdatamanager.CacheWrapper{cache}

	cases := []struct {
		userTriggers    []usermodel.Trigger
		userTriggersErr error

		addTrigger usermodel.Trigger

		saveTriggers []usermodel.Trigger

		saveUserTriggersCall bool
		saveUserTriggersErr  error

		removeUserTriggersCall bool
	}{

		{
			addTrigger: usermodel.SUPPORT_CIGGARETTE,

			saveTriggers: []usermodel.Trigger{usermodel.SUPPORT_CIGGARETTE},

			saveUserTriggersCall: true,
		},
		{
			userTriggers: []usermodel.Trigger{
				usermodel.SUPPORT_HEALTH,
			},
			addTrigger: usermodel.SUPPORT_HEALTH,
		},

		{
			userTriggers:         []usermodel.Trigger{usermodel.SUPPORT_HEALTH},
			addTrigger:           usermodel.SUPPORT_TRIAL,
			saveUserTriggersCall: true,
			saveTriggers: []usermodel.Trigger{
				usermodel.SUPPORT_HEALTH,
				usermodel.SUPPORT_TRIAL,
			},
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		cache.EXPECT().UserTriggers(gomock.Any(), userId).Return(cs.userTriggers, cs.userTriggersErr).Times(1)
		if cs.saveUserTriggersCall {
			cache.EXPECT().SaveUserTriggers(gomock.Any(), userId, cs.saveTriggers).Return(cs.saveUserTriggersErr).Times(1)
		}
		if cs.removeUserTriggersCall {
			cache.EXPECT().RemoveUserTriggers(gomock.Any(), userId).Return(nil).Times(1)
		}
		cacheWrapper.Add(ctx, userId, cs.addTrigger)
	}
}

func Test_Triggers_Manager_Remove(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := triggerdatamanager.NewMockTriggerCache(ctrl)
	storage := triggerdatamanager.NewMockTriggerStorage(ctrl)

	manager := triggerdatamanager.NewTriggerManager(storage, cache)

	{
		userId := rand.Int63()

		removeTrigger := usermodel.SUPPORT_TRIAL

		expectedErr := errors.New("failed remove trigger manager")

		storage.EXPECT().Remove(gomock.Any(), userId, removeTrigger).Return(expectedErr).Times(1)

		err := manager.Remove(ctx, userId, removeTrigger)

		assert.ErrorIs(t, err, expectedErr, "wrong error from remove")
	}

	{
		userId := rand.Int63()

		removeTrigger := usermodel.THANK_YOU

		storage.EXPECT().Remove(gomock.Any(), userId, removeTrigger).Return(nil).Times(1)
		cache.EXPECT().UserTriggers(gomock.Any(), userId).Return(nil, errors.New("random err")).Times(1)

		err := manager.Remove(ctx, userId, removeTrigger)

		assert.ErrorIs(t, err, nil, "wrong err")
	}
}

func Test_Triggers_Manager_Add(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	cache := triggerdatamanager.NewMockTriggerCache(ctrl)
	storage := triggerdatamanager.NewMockTriggerStorage(ctrl)

	manager := triggerdatamanager.NewTriggerManager(storage, cache)

	cases := []struct {
		storageErr       error
		cacheWrapperCall bool
		err              error
	}{}

	for _, cs := range cases {
		trigger := usermodel.THANK_YOU
		userId := rand.Int63()
		storage.EXPECT().Add(gomock.Any(), userId, trigger).Return(cs.storageErr).Times(1)

		if cs.cacheWrapperCall {
			cache.EXPECT().UserTriggers(gomock.Any(), userId).Return(nil, errors.New("random err")).Times(1)
		}

		err := manager.Add(ctx, userId, trigger)
		assert.ErrorIs(t, err, cs.err)
	}
}
