package subscriptionstorage

import (
	"context"

	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
)

var localSubscriptionStorage = map[string]subscription.Subscription{
	"kurilka_month_subscription_2770": {
		ID:     "kurilka_month_subscription_2770",
		Title:  "Подписка на Месяц",
		Period: subscription.MONTH,
		Price:  23,
	},
}

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Subscription(ctx context.Context, subscriptionId string) (subscription.Subscription, error) {
	subscription, ok := localSubscriptionStorage[subscriptionId]
	if !ok {
		return subscription, subscriptionerror.SubscriptionNotFound
	}
	return subscription, nil
}
