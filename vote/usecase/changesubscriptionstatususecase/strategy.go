package changesubscriptionstatususecase

import (
	"context"
	"errors"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/vote/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
)

//go:generate mockgen -source strategy.go -destination strategy_mocks.go -package changesubscriptionstatususecase

type ChangeSubscriptionStatusStrategy interface {
	Change(ctx context.Context, changeSubscriptionStatus subscription.ChangeSubscriptionStatus) (subscription.ChangeSubscriptionStatusResponse, error)
}

type VoteSubscriptionStorage interface {
	CreateSubscription(ctx context.Context, subscriptionId, userId int64) error
	DeleteSubscription(ctx context.Context, subscriptionId int64) error
	UpdateUserSubscriptionId(ctx context.Context, userId, subscriptionId int64) error
}

type SubscriptionItemStroage interface {
	Subscription(ctx context.Context, subscriptionId string) (subscription.Subscription, error)
}

type AddSubscriptionStrategy struct {
	voteSubscriptionStorage VoteSubscriptionStorage
	subscriptionDataManager subscriptiondatamanager.SubscriptionManager
	subscriptionItemStorage SubscriptionItemStroage
}

func (s *AddSubscriptionStrategy) Change(ctx context.Context, changeSubscriptionStatus subscription.ChangeSubscriptionStatus) (resp subscription.ChangeSubscriptionStatusResponse, err error) {
	subsItem, err := s.subscriptionItemStorage.Subscription(ctx, changeSubscriptionStatus.ItemId)
	if err != nil {
		return resp, err
	}
	err = s.voteSubscriptionStorage.CreateSubscription(ctx, changeSubscriptionStatus.SubscriptionId, changeSubscriptionStatus.UserId)
	switch {
	case errors.Is(err, subscriptionerror.SubscriptionIdExists):
		err = s.voteSubscriptionStorage.UpdateUserSubscriptionId(ctx, changeSubscriptionStatus.UserId, changeSubscriptionStatus.SubscriptionId)
		if err != nil {
			return
		}
	case err == nil:
	default:
		return
	}
	userSubscription := usermodel.NewSubscription(
		usermodel.BASIC,
		time.Now().Add(time.Hour*24*time.Duration(subsItem.Period)),
	)
	err = s.subscriptionDataManager.SetUserSubscription(ctx, changeSubscriptionStatus.UserId, userSubscription)
	if err != nil {
		return
	}
	resp = subscription.ChangeSubscriptionStatusResponse{
		SubscriptionId: changeSubscriptionStatus.SubscriptionId,
	}
	return
}

func NewAddStrategy(
	voteSubscriptionStorage VoteSubscriptionStorage,
	subscriptionDataManager subscriptiondatamanager.SubscriptionManager,
	subscriptionItemStorage SubscriptionItemStroage,
) ChangeSubscriptionStatusStrategy {
	return &AddSubscriptionStrategy{
		voteSubscriptionStorage: voteSubscriptionStorage,
		subscriptionDataManager: subscriptionDataManager,
		subscriptionItemStorage: subscriptionItemStorage,
	}
}

type CancelSubscriptionStrategy struct {
	voteSubscriptionStorage VoteSubscriptionStorage
	subscriptionDataManager subscriptiondatamanager.SubscriptionManager
	subscriptionItemStorage SubscriptionItemStroage
}

func (s *CancelSubscriptionStrategy) Change(ctx context.Context, changeSubscriptionStatus subscription.ChangeSubscriptionStatus) (subscription.ChangeSubscriptionStatusResponse, error) {
	if changeSubscriptionStatus.CancelReason == subscription.USER_DECISION {
		return subscription.ChangeSubscriptionStatusResponse{SubscriptionId: changeSubscriptionStatus.SubscriptionId}, nil
	}
	userSubscription := usermodel.NewSubscription(usermodel.NONE, time.Time{})
	err := s.subscriptionDataManager.SetUserSubscription(ctx, changeSubscriptionStatus.UserId, userSubscription)
	if err != nil {
		return subscription.ChangeSubscriptionStatusResponse{}, err
	}
	if changeSubscriptionStatus.Status == subscription.CANCELLED {
		err = s.voteSubscriptionStorage.DeleteSubscription(ctx, changeSubscriptionStatus.SubscriptionId)
		if err != nil {
			return subscription.ChangeSubscriptionStatusResponse{}, err
		}
	}
	return subscription.ChangeSubscriptionStatusResponse{SubscriptionId: changeSubscriptionStatus.SubscriptionId}, nil
}

func NewCancelStrategy(
	voteSubscriptionStorage VoteSubscriptionStorage,
	subscriptionDataManager subscriptiondatamanager.SubscriptionManager,
	subscriptionItemStorage SubscriptionItemStroage,
) ChangeSubscriptionStatusStrategy {
	return &CancelSubscriptionStrategy{
		voteSubscriptionStorage: voteSubscriptionStorage,
		subscriptionDataManager: subscriptionDataManager,
		subscriptionItemStorage: subscriptionItemStorage,
	}
}
