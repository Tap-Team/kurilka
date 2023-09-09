package subscription_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/transport/http/subscription"
	"github.com/Tap-Team/kurilka/user/usecase/subscriptionusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestUserSubscription(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := subscriptionusecase.NewMockSubscriptionUseCase(ctrl)

	handler := subscription.New(useCase)

	cases := []struct {
		userId      int64
		accessToken string

		queryValues      map[string]string
		subscriptionType usermodel.SubscriptionType
		useCaseCall      bool
		err              error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId:      100,
			accessToken: "amidman_2004",
			queryValues: map[string]string{
				"vk_user_id":   "100",
				"access_token": "amidman_2004",
			},
			subscriptionType: usermodel.BASIC,
			useCaseCall:      true,
			statusCode:       http.StatusOK,
		},
	}

	for _, cs := range cases {
		if cs.useCaseCall {
			useCase.EXPECT().UserSubscription(gomock.Any(), cs.userId).Return(cs.subscriptionType).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		req := httptest.NewRequest(http.MethodGet, "/subscription/user?"+urlValues.Encode(), nil)
		rec := httptest.NewRecorder()
		handler.UserSubscription(ctx).ServeHTTP(rec, req)

		assert.Equal(t, cs.statusCode, rec.Result().StatusCode)

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			b, err := io.ReadAll(rec.Result().Body)
			assert.NilError(t, err, "failed decode body")
			subscriptionType := usermodel.SubscriptionType(b)
			assert.Equal(t, subscriptionType, cs.subscriptionType, "subscription type not equal")
		}
	}
}
