package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tap-Team/kurilka/callback/handler"
	"github.com/Tap-Team/kurilka/callback/usecase/subscriptionusecase"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_Handler_CallBackHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := subscriptionusecase.NewMockUseCase(ctrl)

	const (
		confirmCode = "confirm_code"
		groupId     = 1
		secret      = "1"
	)

	callBackHandler := handler.New(confirmCode, groupId, secret, useCase)

	cases := []struct {
		before func()
		body   []byte
		object json.RawMessage

		err        error
		statusCode int
		response   string
	}{
		{
			body: []byte(`
			{ "type": "confirmation", "group_id": 221790286 }	
			`),
			statusCode: http.StatusOK,
			response:   confirmCode,
		},
		{
			body: []byte(`
				{
					"type": "group_join",
					"object": {
						"user_id": 1,
						"join_type": "approved"
					},
					"group_id": 1
				}
				`,
			),
			statusCode: http.StatusInternalServerError,
			err:        handler.ErrWrongSecret,
		},

		{
			body: []byte(`
				{
					"type":"donut_subscription_create",
					"object": {
						"user_id":123,
						"amount":169,
						"amount_without_fee": 180.00
					},
					"secret":"1"
				}
			`),
			before: func() {
				useCase.EXPECT().CreateSubscription(gomock.Any(), int64(123), float64(169)).Return(nil).Times(1)
			},
			statusCode: http.StatusOK,
			response:   "ok",
		},
		{
			body: []byte(`
				{
					"type":"donut_subscription_create",
					"object": {
						"user_id":123,
						"amount":169,
						"amount_without_fee": 180.00
					},
					"secret":"1"
				}
			`),
			before: func() {
				useCase.EXPECT().CreateSubscription(gomock.Any(), int64(123), float64(169)).Return(usererror.ExceptionUserNotFound()).Times(1)
			},
			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			body: []byte(`
			{
				"type":"donut_subscription_expired",
				"object": {
					"user_id":1
				},
				"secret":"1"
			}
		`),
			before: func() {
				useCase.EXPECT().CleanSubscription(gomock.Any(), int64(1)).Return(nil).Times(1)
			},
			statusCode: http.StatusOK,
			response:   "ok",
		},
		{
			body: []byte(`
			{
				"type":"donut_subscription_expired",
				"object": {
					"user_id":1
				},
				"secret":"1"
			}
		`),
			before: func() {
				useCase.EXPECT().CleanSubscription(gomock.Any(), int64(1)).Return(usererror.ExceptionUserNotFound()).Times(1)
			},
			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},

		{
			body: []byte(`
			{
				"type":"donut_subscription_prolonged",
				"object": {
					"user_id":321,
					"amount":169,
					"amount_without_fee": 180.00
				},
				"secret":"1"
			}
		`),
			before: func() {
				useCase.EXPECT().ProlongSubscription(gomock.Any(), int64(321), float64(169)).Return(nil).Times(1)
			},
			statusCode: http.StatusOK,
			response:   "ok",
		},
		{
			body: []byte(`
			{
				"type":"donut_subscription_prolonged",
				"object": {
					"user_id":321,
					"amount":169.000000,
					"amount_without_fee": 180.00
				},
				"secret":"1"
			}
		`),
			before: func() {
				useCase.EXPECT().ProlongSubscription(gomock.Any(), int64(321), float64(169)).Return(usererror.ExceptionUserNotFound()).Times(1)
			},
			statusCode: http.StatusNotFound,
			err:        usererror.ExceptionUserNotFound(),
		},
	}

	for _, cs := range cases {
		if cs.before != nil {
			cs.before()
		}
		body := bytes.NewReader(cs.body)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/vk/callback", body)

		callBackHandler.CallBackHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, cs.statusCode, rec.Result().StatusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			b, err := io.ReadAll(rec.Result().Body)
			assert.NilError(t, err, "failed read result")
			response := string(b)

			assert.Equal(t, response, cs.response, "wrong response")
		}
	}
}
