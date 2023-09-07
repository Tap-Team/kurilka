package subscriptionusecase

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
)

type UserSubscriptionStorage interface {
	Subscription(ctx context.Context, userId int64) (*usermodel.Subscription, error)
	UpdateSubscription(ctx context.Context, userId int64, subscriptionType usermodel.SubscriptionType, expired time.Time) error
}

type VK_Subscription_Manager interface {
	UserSubscription(ctx context.Context, accessToken string) (time.Time, error)
}

type subscriptionUseCase struct {
	vk      VK_Subscription_Manager
	storage UserSubscriptionStorage
}

type SubscriptionUseCase interface {
	UserSubscription(ctx context.Context, userId int64, vkUserToken string) usermodel.SubscriptionType
}

func New(vk VK_Subscription_Manager, storage UserSubscriptionStorage) SubscriptionUseCase {
	return &subscriptionUseCase{
		vk:      vk,
		storage: storage,
	}
}

func (s *subscriptionUseCase) UserSubscription(ctx context.Context, userId int64, accessToken string) usermodel.SubscriptionType {
	subscription, err := s.storage.Subscription(ctx, userId)
	if err != nil {
		return usermodel.NONE
	}
	if !subscription.IsNoneOrExpired() {
		return subscription.Type
	}
	expired, err := s.vk.UserSubscription(ctx, accessToken)
	if err != nil {
		return usermodel.NONE
	}
	subscriptionType := usermodel.BASIC
	err = s.storage.UpdateSubscription(ctx, userId, subscriptionType, expired)
	if err != nil {
		return usermodel.NONE
	}
	return subscriptionType
}
