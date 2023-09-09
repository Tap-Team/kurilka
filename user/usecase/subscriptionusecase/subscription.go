package subscriptionusecase

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/subscriptiondatamanager"
)

//go:generate mockgen -source subscription.go -destination subscription_mocks.go -package subscriptionusecase

type VK_Subscription_Manager interface {
	UserSubscriptionById(ctx context.Context, userId int64) (time.Time, error)
}

type subscriptionUseCase struct {
	vk           VK_Subscription_Manager
	subscription subscriptiondatamanager.SubscriptionManager
}

type SubscriptionUseCase interface {
	UserSubscription(ctx context.Context, userId int64) usermodel.SubscriptionType
}

func New(vk VK_Subscription_Manager, subscription subscriptiondatamanager.SubscriptionManager) SubscriptionUseCase {
	return &subscriptionUseCase{
		vk:           vk,
		subscription: subscription,
	}
}

func (s *subscriptionUseCase) UserSubscription(ctx context.Context, userId int64) usermodel.SubscriptionType {
	subscription, err := s.subscription.UserSubscription(ctx, userId)
	if err != nil {
		return usermodel.NONE
	}
	if !subscription.IsNoneOrExpired() {
		return subscription.Type
	}
	expired, err := s.vk.UserSubscriptionById(ctx, userId)
	if err != nil {
		return usermodel.NONE
	}
	subscriptionType := usermodel.BASIC
	err = s.subscription.UpdateUserSubscription(ctx, userId, usermodel.NewSubscription(subscriptionType, expired))
	if err != nil {
		return usermodel.NONE
	}
	return subscriptionType
}
