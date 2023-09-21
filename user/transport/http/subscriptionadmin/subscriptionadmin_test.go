package subscriptionadmin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	usermodel "github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/transport/http/subscriptionadmin"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestUpdateSubscription(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	updater := subscriptionadmin.NewMockUserSubscriptionUpdater(ctrl)

	handler := subscriptionadmin.New(updater)

	cases := []struct {
		userId       int64
		subscription usermodel.Subscription

		queryValues map[string]string

		updaterCall bool
		updaterErr  error

		err error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 1,
			queryValues: map[string]string{
				"vk_user_id": "1",
			},
			err:        subscriptionadmin.ErrFailedParseExpired,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 200,
			queryValues: map[string]string{
				"vk_user_id": "200",
				"expired":    "1",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id":       "123",
				"expired":          "1",
				"subscriptionType": string(usermodel.BASIC),
			},
			subscription: usermodel.NewSubscription(usermodel.BASIC, time.Unix(1, 0)),
			updaterCall:  true,
			statusCode:   http.StatusNoContent,
		},
		{
			userId: 321,
			queryValues: map[string]string{
				"vk_user_id":       "321",
				"expired":          "12345",
				"subscriptionType": string(usermodel.BASIC),
			},
			subscription: usermodel.NewSubscription(usermodel.BASIC, time.Unix(12345, 0)),
			updaterCall:  true,
			updaterErr:   usererror.ExceptionUserNotFound(),
			err:          usererror.ExceptionUserNotFound(),
			statusCode:   http.StatusNotFound,
		},
	}

	for _, cs := range cases {
		if cs.updaterCall {
			updater.EXPECT().UpdateUserSubscription(gomock.Any(), cs.userId, cs.subscription).Return(cs.updaterErr).Times(1)
		}
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscription/update?"+urlValues.Encode(), nil)

		handler.UpdateUserSubscriptionHandler(ctx).ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode
		assert.Equal(t, statusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}
}
