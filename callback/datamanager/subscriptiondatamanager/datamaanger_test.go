package subscriptiondatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/callback/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	ZeroSubscription usermodel.Subscription
)

func Test_Manager_UserSubsctiption(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := subscriptiondatamanager.NewMockSubscriptionCache(ctrl)
	storage := subscriptiondatamanager.NewMockSubscriptionStorage(ctrl)

	manager := subscriptiondatamanager.New(storage, cache)

	{
		userId := rand.Int63()
		userSubscription := random.StructTyped[usermodel.Subscription]()

		cache.EXPECT().UserSubscription(gomock.Any(), userId).Return(userSubscription, nil).Times(1)

		subscription, err := manager.UserSubscription(ctx, userId)

		assert.Equal(t, subscription, userSubscription, "subscription not equal")
		assert.NilError(t, err, "non nil error")
	}

	{
		userId := rand.Int63()
		userSubscription := random.StructTyped[usermodel.Subscription]()

		cache.EXPECT().UserSubscription(gomock.Any(), userId).Return(ZeroSubscription, errors.New("any")).Times(1)
		storage.EXPECT().UserSubscription(gomock.Any(), userId).Return(userSubscription, nil).Times(1)

		subscription, err := manager.UserSubscription(ctx, userId)

		assert.Equal(t, subscription, userSubscription, "subscription not equal")
		assert.NilError(t, err, "non nil error")
	}

	{
		userId := rand.Int63()
		expectedErr := usererror.ExceptionUserNotFound()

		cache.EXPECT().UserSubscription(gomock.Any(), userId).Return(ZeroSubscription, errors.New("any")).Times(1)
		storage.EXPECT().UserSubscription(gomock.Any(), userId).Return(ZeroSubscription, expectedErr).Times(1)

		subscription, err := manager.UserSubscription(ctx, userId)

		assert.Equal(t, subscription, ZeroSubscription, "subscription not equal")
		assert.ErrorIs(t, err, expectedErr, "wrong err")
	}
}

func Test_Manager_UpdateUserSubscription(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := subscriptiondatamanager.NewMockSubscriptionCache(ctrl)
	storage := subscriptiondatamanager.NewMockSubscriptionStorage(ctrl)

	manager := subscriptiondatamanager.New(storage, cache)

	{
		userId := rand.Int63()
		userSubscription := random.StructTyped[usermodel.Subscription]()
		expectedErr := usererror.ExceptionUserNotFound()

		storage.EXPECT().UpdateUserSubscription(gomock.Any(), userId, userSubscription).Return(expectedErr).Times(1)

		err := manager.SetUserSubscription(ctx, userId, userSubscription)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	{
		userId := rand.Int63()
		userSubscription := random.StructTyped[usermodel.Subscription]()

		storage.EXPECT().UpdateUserSubscription(gomock.Any(), userId, userSubscription).Return(nil).Times(1)
		cache.EXPECT().UpdateUserSubscription(gomock.Any(), userId, userSubscription).Return(nil).Times(1)

		err := manager.SetUserSubscription(ctx, userId, userSubscription)

		assert.NilError(t, err, "non nil error")
	}

	{
		userId := rand.Int63()
		userSubscription := random.StructTyped[usermodel.Subscription]()

		storage.EXPECT().UpdateUserSubscription(gomock.Any(), userId, userSubscription).Return(nil).Times(1)
		cache.EXPECT().UpdateUserSubscription(gomock.Any(), userId, userSubscription).Return(errors.New("any error")).Times(1)
		cache.EXPECT().RemoveUserSubscription(gomock.Any(), userId).Return(nil)

		err := manager.SetUserSubscription(ctx, userId, userSubscription)

		assert.NilError(t, err, "non nil error")
	}
}
