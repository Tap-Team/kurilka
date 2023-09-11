package subscriptionusecase_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/callback/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/callback/usecase/subscriptionusecase"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type monthSubscriptionMatcher struct {
	month time.Month
	t     usermodel.SubscriptionType
}

func NewMonthSubscriptionMatcher(month time.Month, stype usermodel.SubscriptionType) *monthSubscriptionMatcher {
	return &monthSubscriptionMatcher{month: month, t: stype}
}

func (m *monthSubscriptionMatcher) Matches(x interface{}) bool {
	subscription, ok := x.(usermodel.Subscription)
	if !ok {
		return false
	}
	return subscription.Expired.Month() == m.month && subscription.Type == m.t
}

func (m *monthSubscriptionMatcher) String() string {
	return fmt.Sprintf("month is equal %d and type is equal %s", m.month, m.t)
}

func month(m int) time.Month {
	if m <= 12 {
		return time.Month(m)
	}
	m = m % 12
	if m == 0 {
		m = 12
	}
	return time.Month(m)
}

func Test_UseCase_CreateSubscription(t *testing.T) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := subscriptiondatamanager.NewMockSubscriptionManager(ctrl)
	const subscriptionCost = 178

	useCase := subscriptionusecase.New(manager, subscriptionCost)

	currentMonth := int(time.Now().Month())

	cases := []struct {
		amount int

		month int

		err error
	}{
		{
			month: currentMonth,
			err:   errors.New("random error while set user subscription"),
		},
		{
			amount: 88,
			month:  currentMonth,
		},
		{
			amount: 89,
			month:  currentMonth + 1,
		},
		{
			amount: 266,
			month:  currentMonth + 1,
		},
		{
			amount: 267,
			month:  currentMonth + 2,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		manager.EXPECT().SetUserSubscription(gomock.Any(), userId, NewMonthSubscriptionMatcher(month(cs.month), usermodel.BASIC)).Return(cs.err).Times(1)

		err := useCase.CreateSubscription(ctx, userId, cs.amount)
		assert.ErrorIs(t, err, cs.err, "wrong error")
	}

}

func Test_UseCase_CleanSubscription(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := subscriptiondatamanager.NewMockSubscriptionManager(ctrl)

	useCase := subscriptionusecase.New(manager, 0)

	cases := []struct {
		err error
	}{
		{},
		{
			err: errors.New("any error"),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		manager.EXPECT().SetUserSubscription(gomock.Any(), userId, usermodel.NewSubscription(usermodel.NONE, time.Time{})).Return(cs.err).Times(1)

		err := useCase.CleanSubscription(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}

}

func Test_UseCase_ProlongSubscription(t *testing.T) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := subscriptiondatamanager.NewMockSubscriptionManager(ctrl)
	const subscriptionCost = 178

	useCase := subscriptionusecase.New(manager, subscriptionCost)

	currentMonth := int(time.Now().Month())

	cases := []struct {
		amount int
		month  int

		userSubscriptionCall     bool
		userSubscriptionResponse struct {
			subscription usermodel.Subscription
			err          error
		}
		setSubscriptionCall bool
		setSubscriptionErr  error

		err error
	}{
		{
			amount: 88,
			month:  currentMonth,

			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now()),
			},
			setSubscriptionCall: true,
		},
		{
			amount:               89,
			month:                currentMonth + 1,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now()),
			},
			setSubscriptionCall: true,
		},
		{
			amount:               266,
			month:                currentMonth + 1,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now()),
			},
			setSubscriptionCall: true,
		},
		{
			amount:               267,
			month:                currentMonth + 2,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now()),
			},
			setSubscriptionCall: true,
		},

		{
			amount:               178,
			month:                currentMonth + 1,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.NONE, time.Now()),
			},
			setSubscriptionCall: true,
		},
		{
			amount:               178,
			month:                currentMonth + 1,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Time{}),
			},
			setSubscriptionCall: true,
		},
		{
			amount:               178,
			month:                currentMonth + 1,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.NONE, time.Time{}),
			},
			setSubscriptionCall: true,
		},
		{
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				err: usererror.ExceptionUserNotFound(),
			},
			err: usererror.ExceptionUserNotFound(),
		},

		{
			amount:               178,
			month:                currentMonth + 1,
			userSubscriptionCall: true,
			userSubscriptionResponse: struct {
				subscription usermodel.Subscription
				err          error
			}{
				subscription: usermodel.NewSubscription(usermodel.NONE, time.Time{}),
			},
			setSubscriptionCall: true,
			setSubscriptionErr:  usererror.ExceptionUserNotFound(),
			err:                 usererror.ExceptionUserNotFound(),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		if cs.userSubscriptionCall {
			manager.EXPECT().UserSubscription(gomock.Any(), userId).Return(cs.userSubscriptionResponse.subscription, cs.userSubscriptionResponse.err).Times(1)
		}
		if cs.setSubscriptionCall {
			manager.EXPECT().SetUserSubscription(gomock.Any(), userId, NewMonthSubscriptionMatcher(month(cs.month), usermodel.BASIC)).Return(cs.setSubscriptionErr).Times(1)
		}
		err := useCase.ProlongSubscription(ctx, userId, cs.amount)
		assert.ErrorIs(t, err, cs.err, "wrong error")
	}
}
