package subscriptionusecase_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/user/usecase/subscriptionusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	ZeroSubscription usermodel.Subscription
)

func Test_UseCase_UserSubscription(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	vk_manager := subscriptionusecase.NewMockVK_Subscription_Manager(ctrl)
	manager := subscriptiondatamanager.NewMockSubscriptionManager(ctrl)

	useCase := subscriptionusecase.New(vk_manager, manager)

	{
		userId := rand.Int63()
		expectedErr := errors.New("random error")

		manager.EXPECT().UserSubscription(gomock.Any(), userId).Return(ZeroSubscription, expectedErr).Times(1)

		subscriptionType := useCase.UserSubscription(ctx, userId)
		assert.Equal(t, subscriptionType, usermodel.NONE)
	}

	{
		cases := []struct {
			subscription usermodel.Subscription
			expectedType usermodel.SubscriptionType
		}{
			{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour)),
				expectedType: usermodel.BASIC,
			},
			{
				subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
				expectedType: usermodel.TRIAL,
			},
		}

		for _, cs := range cases {
			userId := rand.Int63()
			userSubscription := cs.subscription
			expectedType := cs.expectedType

			manager.EXPECT().UserSubscription(gomock.Any(), userId).Return(userSubscription, nil).Times(1)

			subscriptionType := useCase.UserSubscription(ctx, userId)

			assert.Equal(t, subscriptionType, expectedType)
		}

	}

	{
		cases := []struct {
			susbscription usermodel.Subscription
		}{
			{
				susbscription: usermodel.NewSubscription(usermodel.BASIC, time.Time{}),
			},
			{
				susbscription: usermodel.NewSubscription(usermodel.TRIAL, time.Time{}),
			},
			{
				susbscription: usermodel.NewSubscription(usermodel.NONE, time.Time{}),
			},
			{
				susbscription: usermodel.NewSubscription(usermodel.NONE, time.Now().Add(time.Hour)),
			},
		}
		for _, cs := range cases {
			userId := rand.Int63()
			userSubscription := cs.susbscription

			manager.EXPECT().UserSubscription(gomock.Any(), userId).Return(userSubscription, nil).Times(1)

			vk_manager.EXPECT().UserSubscriptionById(gomock.Any(), userId).Return(time.Time{}, errors.New("any error")).Times(1)
			if cs.susbscription.Type != usermodel.NONE {
				manager.EXPECT().UpdateUserSubscription(gomock.Any(), userId, usermodel.NewSubscription(usermodel.NONE, time.Time{})).Return(nil).Times(1)
			}

			subscriptionType := useCase.UserSubscription(ctx, userId)
			assert.Equal(t, usermodel.NONE, subscriptionType)
		}

	}

	{
		userId := rand.Int63()
		userSubscription := usermodel.NewSubscription(usermodel.BASIC, time.Time{})
		subscriptionExpired := time.Now().Add(time.Hour * time.Duration(rand.Intn(100)))

		manager.EXPECT().UserSubscription(gomock.Any(), userId).Return(userSubscription, nil).Times(1)

		vk_manager.EXPECT().UserSubscriptionById(gomock.Any(), userId).Return(subscriptionExpired, nil).Times(1)

		manager.EXPECT().UpdateUserSubscription(gomock.Any(), userId, usermodel.NewSubscription(usermodel.BASIC, subscriptionExpired)).Return(nil)

		subscripitionType := useCase.UserSubscription(ctx, userId)

		assert.Equal(t, subscripitionType, usermodel.BASIC)
	}

	{
		userId := rand.Int63()
		userSubscription := usermodel.NewSubscription(usermodel.BASIC, time.Time{})
		subscriptionExpired := time.Now().Add(time.Hour * time.Duration(rand.Intn(100)))

		manager.EXPECT().UserSubscription(gomock.Any(), userId).Return(userSubscription, nil).Times(1)

		vk_manager.EXPECT().UserSubscriptionById(gomock.Any(), userId).Return(subscriptionExpired, nil).Times(1)

		manager.EXPECT().UpdateUserSubscription(gomock.Any(), userId, usermodel.NewSubscription(usermodel.BASIC, subscriptionExpired)).Return(errors.New("any"))

		subscripitionType := useCase.UserSubscription(ctx, userId)

		assert.Equal(t, subscripitionType, usermodel.NONE)
	}

}
