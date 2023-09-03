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
