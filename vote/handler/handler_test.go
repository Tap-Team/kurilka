package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/Tap-Team/kurilka/vote/handler"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type UserSubscriptionCall struct {
	SubscriptionId int64
	Err            error
}

func (c UserSubscriptionCall) Call() bool {
	return c != UserSubscriptionCall{}
}

func Test_UserSubscriptionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	vkSecret := ""

	subscriptionStorage := handler.NewMockSubscriptionStorage(ctrl)
	getSubscriptionUseCase := handler.NewMockGetSubscriptionUseCase(ctrl)
	changeSubscriptionStatusUseCase := handler.NewMockChangeSubscriptionStatusUseCase(ctrl)

	h := handler.New(vkSecret, getSubscriptionUseCase, changeSubscriptionStatusUseCase, subscriptionStorage)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		userSubscriptionCall UserSubscriptionCall

		statusCode int

		subscriptionId int64
		err            error
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 257824,
			queryValues: map[string]string{
				"vk_user_id": "257824",
			},

			userSubscriptionCall: UserSubscriptionCall{
				Err: subscriptionerror.SubscriptionNotFound,
			},
			statusCode: http.StatusNotFound,
			err:        subscriptionerror.SubscriptionNotFound,
		},
		{
			userId: 348571418,
			queryValues: map[string]string{
				"vk_user_id": "348571418",
			},

			userSubscriptionCall: UserSubscriptionCall{
				SubscriptionId: 25485452485,
			},
			statusCode:     http.StatusOK,
			subscriptionId: 25485452485,
		},
	}

	for _, cs := range cases {
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}
		if cs.userSubscriptionCall.Call() {
			subscriptionStorage.EXPECT().
				UserSubscriptionId(gomock.Any(), cs.userId).
				Return(cs.userSubscriptionCall.SubscriptionId, cs.userSubscriptionCall.Err).
				Times(1)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/vote-subscription-id?"+urlValues.Encode(), nil)

		h.UserSubscriptionIdHandler().ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode

		assert.Equal(t, cs.statusCode, statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var id int64
			err := json.NewDecoder(rec.Result().Body).Decode(&id)
			assert.NilError(t, err, "failed decode body")
			assert.Equal(t, cs.subscriptionId, id, "subscription id not equal")
		}
	}
}

type GetSubscriptionHandlerCall struct {
	SubscriptionId string

	Subscription subscription.Subscription
	Err          error
}

func (c GetSubscriptionHandlerCall) Call() bool {
	return c != GetSubscriptionHandlerCall{}
}

func Test_GetSubscriptionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	vkSecret := ""

	subscriptionStorage := handler.NewMockSubscriptionStorage(ctrl)
	getSubscriptionUseCase := handler.NewMockGetSubscriptionUseCase(ctrl)
	changeSubscriptionStatusUseCase := handler.NewMockChangeSubscriptionStatusUseCase(ctrl)

	h := handler.New(vkSecret, getSubscriptionUseCase, changeSubscriptionStatusUseCase, subscriptionStorage)

	cases := []struct {
		params handler.Params

		getSubscriptionhandlerCall GetSubscriptionHandlerCall

		statusCode int

		resp subscription.Subscription
		err  error
	}{
		{
			params: handler.Params{
				"item": "efjasddfllasdfjlakdj",
			},

			getSubscriptionhandlerCall: GetSubscriptionHandlerCall{
				SubscriptionId: "efjasddfllasdfjlakdj",
				Err:            subscriptionerror.SubscriptionNotFound,
			},
			err:        subscriptionerror.SubscriptionNotFound,
			statusCode: http.StatusNotFound,
		},
		{
			params: handler.Params{
				"item": "ryuwrnnabkaqwuqhqrqjklfsgiuqq",
			},

			getSubscriptionhandlerCall: GetSubscriptionHandlerCall{
				SubscriptionId: "ryuwrnnabkaqwuqhqrqjklfsgiuqq",
				Subscription: subscription.Subscription{
					ID:     "ryuwrnnabkaqwuqhqrqjklfsgiuqq",
					Title:  "abiba",
					Price:  23,
					Period: subscription.MONTH,
				},
			},
			statusCode: http.StatusOK,
			resp: subscription.Subscription{
				ID:     "ryuwrnnabkaqwuqhqrqjklfsgiuqq",
				Title:  "abiba",
				Price:  23,
				Period: subscription.MONTH,
			},
		},
	}

	for _, cs := range cases {
		if cs.getSubscriptionhandlerCall.Call() {
			getSubscriptionUseCase.EXPECT().
				Subscription(gomock.Any(), cs.getSubscriptionhandlerCall.SubscriptionId).
				Return(cs.getSubscriptionhandlerCall.Subscription, cs.getSubscriptionhandlerCall.Err).
				Times(1)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/any-url", nil)

		h.GetSubscriptionHandler(cs.params).ServeHTTP(rec, req)

		statusCode := cs.statusCode
		assert.Equal(t, cs.statusCode, statusCode, "status code not equal")

		if cs.err != nil {
			handler.AssertError(t, cs.err, rec.Result().Body)
		} else {
			handler.AssertResponse(t, cs.resp, rec.Result().Body)
		}
	}
}

type ChangeSubscriptionStatusCall struct {
	ChangeSubscriptionStatus subscription.ChangeSubscriptionStatus
	Response                 subscription.ChangeSubscriptionStatusResponse
	Err                      error
}

func (c ChangeSubscriptionStatusCall) Call() bool {
	return c != ChangeSubscriptionStatusCall{}
}

func Test_ChangeSubscriptionStatusHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	vkSecret := ""

	subscriptionStorage := handler.NewMockSubscriptionStorage(ctrl)
	getSubscriptionUseCase := handler.NewMockGetSubscriptionUseCase(ctrl)
	changeSubscriptionStatusUseCase := handler.NewMockChangeSubscriptionStatusUseCase(ctrl)

	h := handler.New(vkSecret, getSubscriptionUseCase, changeSubscriptionStatusUseCase, subscriptionStorage)

	cases := []struct {
		params handler.Params

		changeSubscriptionStatusCall ChangeSubscriptionStatusCall

		statusCode int

		resp subscription.ChangeSubscriptionStatusResponse
		err  error
	}{
		{
			params: handler.Params{
				"subscription_id": int64(1238123712),
				"user_id":         int64(182318231),
				"cancel_reason":   string(subscription.USER_DECISION),
				"item_id":         "sdkfjhajksdhfjaeejrr",
				"status":          string(subscription.ACTIVE),
			},

			changeSubscriptionStatusCall: ChangeSubscriptionStatusCall{
				ChangeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
					1238123712,
					182318231,
					"sdkfjhajksdhfjaeejrr",
					subscription.ACTIVE,
					subscription.USER_DECISION,
				),
				Err: subscriptionerror.SubscriptionNotFound,
			},
			err:        subscriptionerror.SubscriptionNotFound,
			statusCode: http.StatusNotFound,
		},

		{
			params: handler.Params{
				"subscription_id": int64(234874259785),
				"user_id":         int64(990926164296967),
				"cancel_reason":   string(subscription.PAYMENT_FAIL),
				"item_id":         "sdkfjhajksdhfjaeejrr",
				"status":          string(subscription.CANCELLED),
			},

			changeSubscriptionStatusCall: ChangeSubscriptionStatusCall{
				ChangeSubscriptionStatus: subscription.NewChangeSubscriptionStatus(
					234874259785,
					990926164296967,
					"sdkfjhajksdhfjaeejrr",
					subscription.CANCELLED,
					subscription.PAYMENT_FAIL,
				),
				Err: subscriptionerror.SubscriptionNotFound,
			},
			err:        subscriptionerror.SubscriptionNotFound,
			statusCode: http.StatusNotFound,
		},
	}

	for _, cs := range cases {
		if cs.changeSubscriptionStatusCall.Call() {
			changeSubscriptionStatusUseCase.EXPECT().
				ChangeSubscriptionStatus(gomock.Any(), cs.changeSubscriptionStatusCall.ChangeSubscriptionStatus).
				Return(cs.changeSubscriptionStatusCall.Response, cs.changeSubscriptionStatusCall.Err).
				Times(1)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/any-url", nil)

		h.SubscriptionStatusChangeHandler(cs.params).ServeHTTP(rec, req)

		statusCode := cs.statusCode
		assert.Equal(t, cs.statusCode, statusCode, "status code not equal")

		if cs.err != nil {
			handler.AssertError(t, cs.err, rec.Result().Body)
		} else {
			handler.AssertResponse(t, cs.resp, rec.Result().Body)
		}
	}
}
