package changesubscriptionstatususecase_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/vote/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
	"github.com/Tap-Team/kurilka/vote/usecase/changesubscriptionstatususecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type GetSubscriptionItemCall struct {
	ItemId string

	SubscriptionItem subscription.Subscription
	Err              error
}

func (c GetSubscriptionItemCall) Call() bool {
	return c != GetSubscriptionItemCall{}
}

type CreateVoteSubscriptionCall struct {
	SubscriptionId, UserId int64

	Err error
}

func (c CreateVoteSubscriptionCall) Call() bool {
	return c != CreateVoteSubscriptionCall{}
}

type UpdateVoteSubscriptionCall struct {
	SubscriptionId, UserId int64
	Err                    error
}

func (c UpdateVoteSubscriptionCall) Call() bool {
	return c != UpdateVoteSubscriptionCall{}
}

type DeleteVoteSubscriptionCall struct {
	SubscriptionId int64

	Err error
}

func (c DeleteVoteSubscriptionCall) Call() bool {
	return c != DeleteVoteSubscriptionCall{}
}

type SubscriptionMatcher usermodel.Subscription

func (s SubscriptionMatcher) Matches(x any) bool {
	subs, ok := x.(usermodel.Subscription)
	if !ok {
		return false
	}
	if subs.Type != s.Type {
		return false
	}
	if subs.Expired.Unix() != s.Expired.Unix() {
		return false
	}
	return true
}

func (s SubscriptionMatcher) String() string {
	return fmt.Sprintf("is equal %v", usermodel.Subscription(s))
}

type SetUserSubscriptionCall struct {
	UserId       int64
	Subscription SubscriptionMatcher

	Err error
}

func (c SetUserSubscriptionCall) Call() bool {
	return c != SetUserSubscriptionCall{}
}

var subscriptionItem = subscription.Subscription{
	ID:     "aboba2000",
	Title:  "pofig",
	Price:  23,
	Period: subscription.MONTH,
}

func Test_Strategy_Add(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	voteStorage := changesubscriptionstatususecase.NewMockVoteSubscriptionStorage(ctrl)
	subscriptionManager := subscriptiondatamanager.NewMockSubscriptionManager(ctrl)
	subscriptionItemStorage := changesubscriptionstatususecase.NewMockSubscriptionItemStroage(ctrl)

	strategy := changesubscriptionstatususecase.NewAddStrategy(voteStorage, subscriptionManager, subscriptionItemStorage)

	cases := []struct {
		changeSubscriptionStatus subscription.ChangeSubscriptionStatus

		getSubscriptionItemCall    GetSubscriptionItemCall
		createVoteSubscriptionCall CreateVoteSubscriptionCall
		updateVoteSubscriptionCall UpdateVoteSubscriptionCall
		setUserSubscriptionCall    SetUserSubscriptionCall

		resp subscription.ChangeSubscriptionStatusResponse
		err  error
	}{
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(1584239618, 18345123, "aboba2000", subscription.CHARGEABLE, subscription.UNKNOWN),
			getSubscriptionItemCall: GetSubscriptionItemCall{
				ItemId: "aboba2000",
				Err:    subscriptionerror.SubscriptionNotFound,
			},
			err: subscriptionerror.SubscriptionNotFound,
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(1584239618, 18345123, "aboba2000", subscription.CHARGEABLE, subscription.UNKNOWN),
			getSubscriptionItemCall: GetSubscriptionItemCall{
				ItemId:           "aboba2000",
				SubscriptionItem: subscriptionItem,
			},
			createVoteSubscriptionCall: CreateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,

				Err: sql.ErrConnDone,
			},
			err: sql.ErrConnDone,
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(1584239618, 18345123, "aboba2000", subscription.CHARGEABLE, subscription.UNKNOWN),
			getSubscriptionItemCall: GetSubscriptionItemCall{
				ItemId:           "aboba2000",
				SubscriptionItem: subscriptionItem,
			},
			createVoteSubscriptionCall: CreateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,
			},
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 18345123,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.BASIC,
					Expired: amidtime.Timestamp{Time: time.Now().Add(time.Hour * 24 * 30)},
				},
				Err: sql.ErrTxDone,
			},
			err: sql.ErrTxDone,
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(1584239618, 18345123, "aboba2000", subscription.CHARGEABLE, subscription.UNKNOWN),
			getSubscriptionItemCall: GetSubscriptionItemCall{
				ItemId:           "aboba2000",
				SubscriptionItem: subscriptionItem,
			},
			createVoteSubscriptionCall: CreateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,
				Err:            subscriptionerror.SubscriptionIdExists,
			},
			updateVoteSubscriptionCall: UpdateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,
				Err:            subscriptionerror.SubscriptionIdExists,
			},
			err: subscriptionerror.SubscriptionIdExists,
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(1584239618, 18345123, "aboba2000", subscription.CHARGEABLE, subscription.UNKNOWN),
			getSubscriptionItemCall: GetSubscriptionItemCall{
				ItemId:           "aboba2000",
				SubscriptionItem: subscriptionItem,
			},
			createVoteSubscriptionCall: CreateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,
				Err:            subscriptionerror.SubscriptionIdExists,
			},
			updateVoteSubscriptionCall: UpdateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,
			},
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 18345123,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.BASIC,
					Expired: amidtime.Timestamp{Time: time.Now().Add(time.Hour * 24 * 30)},
				},
				Err: sql.ErrTxDone,
			},
			err: sql.ErrTxDone,
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(1584239618, 18345123, "aboba2000", subscription.CHARGEABLE, subscription.UNKNOWN),
			getSubscriptionItemCall: GetSubscriptionItemCall{
				ItemId:           "aboba2000",
				SubscriptionItem: subscriptionItem,
			},
			createVoteSubscriptionCall: CreateVoteSubscriptionCall{
				SubscriptionId: 1584239618,
				UserId:         18345123,
			},
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 18345123,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.BASIC,
					Expired: amidtime.Timestamp{Time: time.Now().Add(time.Hour * 24 * 30)},
				},
			},
			resp: subscription.ChangeSubscriptionStatusResponse{
				SubscriptionId: 1584239618,
			},
		},
	}

	for _, cs := range cases {
		if cs.getSubscriptionItemCall.Call() {
			subscriptionItemStorage.EXPECT().
				Subscription(gomock.Any(), cs.getSubscriptionItemCall.ItemId).
				Return(cs.getSubscriptionItemCall.SubscriptionItem, cs.getSubscriptionItemCall.Err).
				Times(1)
		}
		if cs.createVoteSubscriptionCall.Call() {
			voteStorage.EXPECT().
				CreateSubscription(gomock.Any(), cs.createVoteSubscriptionCall.SubscriptionId, cs.createVoteSubscriptionCall.UserId).
				Return(cs.createVoteSubscriptionCall.Err).
				Times(1)
		}
		if cs.updateVoteSubscriptionCall.Call() {
			voteStorage.EXPECT().
				UpdateUserSubscriptionId(gomock.Any(), cs.updateVoteSubscriptionCall.UserId, cs.updateVoteSubscriptionCall.SubscriptionId).
				Return(cs.updateVoteSubscriptionCall.Err).
				Times(1)
		}
		if cs.setUserSubscriptionCall.Call() {
			subscriptionManager.EXPECT().
				SetUserSubscription(gomock.Any(), cs.setUserSubscriptionCall.UserId, cs.setUserSubscriptionCall.Subscription).
				Return(cs.setUserSubscriptionCall.Err).
				Times(1)
		}

		resp, err := strategy.Change(ctx, cs.changeSubscriptionStatus)
		assert.ErrorIs(t, err, cs.err, "wrong err")
		assert.Equal(t, resp, cs.resp, "response not equal")
	}
}

func Test_Strategy_Cancel(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	voteStorage := changesubscriptionstatususecase.NewMockVoteSubscriptionStorage(ctrl)
	subscriptionManager := subscriptiondatamanager.NewMockSubscriptionManager(ctrl)
	subscriptionItemStorage := changesubscriptionstatususecase.NewMockSubscriptionItemStroage(ctrl)

	strategy := changesubscriptionstatususecase.NewCancelStrategy(voteStorage, subscriptionManager, subscriptionItemStorage)

	cases := []struct {
		changeSubscriptionStatus subscription.ChangeSubscriptionStatus

		setUserSubscriptionCall    SetUserSubscriptionCall
		deleteVoteSubscriptionCall DeleteVoteSubscriptionCall

		resp subscription.ChangeSubscriptionStatusResponse
		err  error
	}{
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(12958, 96861, "random_item_id", subscription.ACTIVE, subscription.USER_DECISION),
			resp:                     subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 12958},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(951432516, 951342, "sssss", subscription.ACTIVE, subscription.PAYMENT_FAIL),
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 951342,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.NONE,
					Expired: amidtime.Timestamp{},
				},
				Err: usererror.ExceptionUserNotFound(),
			},
			err: usererror.ExceptionUserNotFound(),
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(951432516, 951342, "sssss", subscription.ACTIVE, subscription.PAYMENT_FAIL),
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 951342,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.NONE,
					Expired: amidtime.Timestamp{},
				},
			},
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 951432516},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(951432516, 951342, "sssss", subscription.CANCELLED, subscription.PAYMENT_FAIL),
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 951342,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.NONE,
					Expired: amidtime.Timestamp{},
				},
			},
			deleteVoteSubscriptionCall: DeleteVoteSubscriptionCall{
				SubscriptionId: 951432516,
			},
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 951432516},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(951432516, 951342, "sssss", subscription.CANCELLED, subscription.PAYMENT_FAIL),
			setUserSubscriptionCall: SetUserSubscriptionCall{
				UserId: 951342,
				Subscription: SubscriptionMatcher{
					Type:    usermodel.NONE,
					Expired: amidtime.Timestamp{},
				},
			},
			deleteVoteSubscriptionCall: DeleteVoteSubscriptionCall{
				SubscriptionId: 951432516,
				Err:            subscriptionerror.SubscriptionNotFound,
			},
			err: subscriptionerror.SubscriptionNotFound,
		},
	}

	for _, cs := range cases {
		if cs.setUserSubscriptionCall.Call() {
			subscriptionManager.EXPECT().
				SetUserSubscription(gomock.Any(), cs.setUserSubscriptionCall.UserId, cs.setUserSubscriptionCall.Subscription).
				Return(cs.setUserSubscriptionCall.Err).
				Times(1)
		}
		if cs.deleteVoteSubscriptionCall.Call() {
			voteStorage.EXPECT().
				DeleteSubscription(gomock.Any(), cs.deleteVoteSubscriptionCall.SubscriptionId).
				Return(cs.deleteVoteSubscriptionCall.Err).
				Times(1)
		}

		resp, err := strategy.Change(ctx, cs.changeSubscriptionStatus)
		assert.ErrorIs(t, err, cs.err, "wrong err")
		assert.Equal(t, resp, cs.resp, "response not equal")
	}

}
