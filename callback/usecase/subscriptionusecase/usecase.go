package subscriptionusecase

import (
	"context"
	"math"
	"time"

	"github.com/Tap-Team/kurilka/callback/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source usecase.go -destination mocks.go -package subscriptionusecase

const _PROVIDER = "callback/datamanager/subscriptiondatamanager.useCase"

type UseCase interface {
	CreateSubscription(ctx context.Context, userId int64, amount int) error
	CleanSubscription(ctx context.Context, userId int64) error
	ProlongSubscription(ctx context.Context, userId int64, amount int) error
	SetSubscriptionMonthCost(cost int)
}

type useCase struct {
	subscription          subscriptiondatamanager.SubscriptionManager
	subscriptionMonthCost int
}

func New(subscription subscriptiondatamanager.SubscriptionManager, subscriptionPricePerMonth int) UseCase {
	return &useCase{subscription: subscription, subscriptionMonthCost: subscriptionPricePerMonth}
}

func (u *useCase) SetSubscriptionMonthCost(cost int) {
	u.subscriptionMonthCost = cost
}

func (u *useCase) Months(amount int) int {
	months := math.Round(float64(amount) / float64(u.subscriptionMonthCost))
	return int(months)
}

func (u *useCase) CreateSubscription(ctx context.Context, userId int64, amount int) error {
	months := u.Months(amount)
	expired := time.Now().AddDate(0, months, 0)
	subscription := usermodel.NewSubscription(usermodel.BASIC, expired)
	err := u.subscription.SetUserSubscription(ctx, userId, subscription)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("set user subscription", "CreateSubscription", _PROVIDER))
	}
	return nil
}

func (u *useCase) CleanSubscription(ctx context.Context, userId int64) error {
	subscription := usermodel.NewSubscription(usermodel.NONE, time.Time{})
	err := u.subscription.SetUserSubscription(ctx, userId, subscription)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("set user subscription", "CleanSubscription", _PROVIDER))
	}
	return nil
}

func (u *useCase) ProlongSubscription(ctx context.Context, userId int64, amount int) error {
	subscription, err := u.subscription.UserSubscription(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get user subscription", "ProlongSubscription", _PROVIDER))
	}
	months := u.Months(amount)

	if subscription.IsNoneOrExpired() {
		subscription.SetExpired(time.Now())
	}
	subscription.SetExpired(subscription.Expired.AddDate(0, months, 0))
	subscription.Type = usermodel.BASIC
	err = u.subscription.SetUserSubscription(ctx, userId, subscription)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("set user subscription", "ProlongSubscription", _PROVIDER))
	}
	return nil
}
