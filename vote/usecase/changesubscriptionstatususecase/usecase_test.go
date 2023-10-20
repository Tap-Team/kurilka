package changesubscriptionstatususecase_test

import (
	"context"
	"testing"

	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
	"github.com/Tap-Team/kurilka/vote/usecase/changesubscriptionstatususecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type StrategyCall struct {
	Response subscription.ChangeSubscriptionStatusResponse
	Err      error
}

func (c StrategyCall) Call() bool {
	return c != StrategyCall{}
}

func Test_UseCase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	addStrategy := changesubscriptionstatususecase.NewMockChangeSubscriptionStatusStrategy(ctrl)
	cancelStrategy := changesubscriptionstatususecase.NewMockChangeSubscriptionStatusStrategy(ctrl)

	useCase := changesubscriptionstatususecase.New(nil, nil, nil)
	useCase.SetAddStrategy(addStrategy)
	useCase.SetCancelStrategy(cancelStrategy)

	cases := []struct {
		changeSubscriptionStatus subscription.ChangeSubscriptionStatus
		addCall, cancelCall      StrategyCall

		resp subscription.ChangeSubscriptionStatusResponse
		err  error
	}{
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				19581,
				10713,
				"subaaba",
				subscription.CHARGEABLE,
				subscription.UNKNOWN,
			),

			addCall: StrategyCall{
				Response: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 19581},
			},
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 19581},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				19581,
				10713,
				"subaaba",
				subscription.CHARGEABLE,
				subscription.UNKNOWN,
			),

			addCall: StrategyCall{
				Err: subscriptionerror.SubscriptionNotFound,
			},
			err: subscriptionerror.SubscriptionNotFound,
		},

		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				69119,
				107616142,
				"asdhs",
				subscription.ACTIVE,
				subscription.UNKNOWN,
			),
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				69119,
				107616142,
				"asdhs",
				subscription.ACTIVE,
				subscription.PAYMENT_FAIL,
			),
			cancelCall: StrategyCall{
				Response: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
			},
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				69119,
				107616142,
				"asdhs",
				subscription.ACTIVE,
				subscription.APP_DECISION,
			),
			cancelCall: StrategyCall{
				Response: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
			},
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				69119,
				107616142,
				"asdhs",
				subscription.ACTIVE,
				subscription.USER_DECISION,
			),
			cancelCall: StrategyCall{
				Response: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
			},
			resp: subscription.ChangeSubscriptionStatusResponse{SubscriptionId: 69119},
		},
		{
			changeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
				69119,
				107616142,
				"asdhs",
				subscription.ACTIVE,
				subscription.PAYMENT_FAIL,
			),
			cancelCall: StrategyCall{
				Err: subscriptionerror.SubscriptionNotFound,
			},
			err: subscriptionerror.SubscriptionNotFound,
		},
	}

	for _, cs := range cases {
		if cs.addCall.Call() {
			addStrategy.EXPECT().Change(gomock.Any(), cs.changeSubscriptionStatus).Return(cs.addCall.Response, cs.addCall.Err).Times(1)
		}
		if cs.cancelCall.Call() {
			cancelStrategy.EXPECT().Change(gomock.Any(), cs.changeSubscriptionStatus).Return(cs.cancelCall.Response, cs.cancelCall.Err).Times(1)
		}

		resp, err := useCase.ChangeSubscriptionStatus(ctx, cs.changeSubscriptionStatus)
		assert.ErrorIs(t, err, cs.err, "err not equal")
		assert.Equal(t, resp, cs.resp, "response not equal")
	}
}
