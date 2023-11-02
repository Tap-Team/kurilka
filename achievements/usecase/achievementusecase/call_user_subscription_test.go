package achievementusecase_test

import (
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	usermodel "github.com/Tap-Team/kurilka/internal/model/usermodel"
	gomock "github.com/golang/mock/gomock"
)

type UserSubscriptionCall struct {
	WillBeCalled bool

	UserId int64

	Subscription usermodel.Subscription
	Err          error
}

func (c *UserSubscriptionCall) RegisterCall(subscriptionProvider *achievementusecase.MockSubscriptionProvider) {
	if c.WillBeCalled {
		subscriptionProvider.EXPECT().
			UserSubscription(gomock.Any(), c.UserId).
			Return(c.Subscription, c.Err).
			Times(1)
	}
}

type UserSubscriptionCallBuilder struct {
	UserId int64

	Subscription usermodel.Subscription
	Err          error
}

func (b *UserSubscriptionCallBuilder) SetInput(userId int64) *UserSubscriptionCallBuilder {
	b.UserId = userId
	return b
}

func (b *UserSubscriptionCallBuilder) SetOutput(subscription usermodel.Subscription, err error) *UserSubscriptionCallBuilder {
	b.Subscription = subscription
	b.Err = err
	return b
}

func (b *UserSubscriptionCallBuilder) Build() UserSubscriptionCall {
	if b == nil {
		return UserSubscriptionCall{}
	}
	return UserSubscriptionCall{
		UserId:       b.UserId,
		Subscription: b.Subscription,
		Err:          b.Err,

		WillBeCalled: true,
	}
}
