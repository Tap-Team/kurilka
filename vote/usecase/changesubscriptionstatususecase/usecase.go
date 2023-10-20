package changesubscriptionstatususecase

import (
	"context"

	"github.com/Tap-Team/kurilka/vote/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
)

//go:generate mockgen -source usecase.go -destination usecase_mocks.go -package changesubscriptionstatususecase

type useCase struct {
	addStrategy, cancelStrategy ChangeSubscriptionStatusStrategy
}

func (u *useCase) SetAddStrategy(strategy ChangeSubscriptionStatusStrategy) {
	u.addStrategy = strategy
}

func (u *useCase) SetCancelStrategy(strategy ChangeSubscriptionStatusStrategy) {
	u.cancelStrategy = strategy
}

func New(
	voteSubscriptionStorage VoteSubscriptionStorage,
	subscriptionDataManager subscriptiondatamanager.SubscriptionManager,
	subscriptionItemStorage SubscriptionItemStroage,
) *useCase {
	return &useCase{
		addStrategy: &AddSubscriptionStrategy{
			voteSubscriptionStorage: voteSubscriptionStorage,
			subscriptionDataManager: subscriptionDataManager,
			subscriptionItemStorage: subscriptionItemStorage,
		},
		cancelStrategy: &CancelSubscriptionStrategy{
			voteSubscriptionStorage: voteSubscriptionStorage,
			subscriptionDataManager: subscriptionDataManager,
			subscriptionItemStorage: subscriptionItemStorage,
		},
	}
}

func (u *useCase) ChangeSubscriptionStatus(ctx context.Context, changeSubscriptionStatus subscription.ChangeSubscriptionStatus) (resp subscription.ChangeSubscriptionStatusResponse, err error) {
	switch changeSubscriptionStatus.Status {
	case subscription.CHARGEABLE:
		return u.addStrategy.Change(ctx, changeSubscriptionStatus)
	case subscription.ACTIVE:
		if changeSubscriptionStatus.CancelReason == subscription.UNKNOWN {
			return subscription.ChangeSubscriptionStatusResponse{SubscriptionId: changeSubscriptionStatus.SubscriptionId}, nil
		}
		return u.cancelStrategy.Change(ctx, changeSubscriptionStatus)
	case subscription.CANCELLED:
		return u.cancelStrategy.Change(ctx, changeSubscriptionStatus)
	}
	return
}
