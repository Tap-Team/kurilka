package transport_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/transport"
	"github.com/Tap-Team/kurilka/internal/model/errormodel"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestOpenSingle(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievementUseCase := transport.NewMockAchievementUseCase(ctrl)

	handler := transport.NewAchievementHandler(achievementUseCase)

	openAchievementResponse := model.NewOpenAchievementResponse(time.Now())
	cases := []struct {
		success bool
		before  func()
		request struct {
			req           *http.Request
			userId        int64
			achievementId int64
		}
		response struct {
			code            int
			successResponse *model.OpenAchievementResponse
			errorResponse   errormodel.ErrorResponse
		}
	}{
		{
			success: true,
			before: func() {
				achievementUseCase.EXPECT().OpenSingle(gomock.Any(), 1, 1).Return(openAchievementResponse, nil).Times(1)
			},
			request: struct {
				req           *http.Request
				userId        int64
				achievementId int64
			}{
				req:           httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/achievements/open-single?vk_user_id=%d&achievementId=%d", 1, 1), nil),
				userId:        1,
				achievementId: 1,
			},
			response: struct {
				code            int
				successResponse *model.OpenAchievementResponse
				errorResponse   errormodel.ErrorResponse
			}{
				code:            http.StatusOK,
				successResponse: openAchievementResponse,
			},
		},
		{},
	}

	for _, cs := range cases {

		cs.before()

		w := httptest.NewRecorder()
		handler.OpenSingleHandler(ctx).ServeHTTP(w, cs.request.req)
		assert.Equal(t, cs.response.code, w.Result().StatusCode, "wrong code")

		if cs.success {
			var errorResponse errormodel.ErrorResponse
			err := json.NewDecoder(w.Result().Body).Decode(&errorResponse)
			require.NoError(t, err, "failed decode body")

			assert.ErrorIs(t, errorResponse, cs.response.errorResponse, "wrong error response")
		} else {
			var openAchievementResponse model.OpenAchievementResponse
			err := json.NewDecoder(w.Result().Body).Decode(&openAchievementResponse)
			require.NoError(t, err, "failed decode body")

			openTime := cs.response.successResponse.OpenTime
			require.Equal(t, openTime.Unix(), openAchievementResponse.OpenTime.Unix(), "open time not equal")
			openAchievementResponse.OpenTime = openTime
			require.Equal(t, cs.response.successResponse, openAchievementResponse, "open achievement response not equal")
		}

	}

}

func TestOpenType(t *testing.T) {

}

func TestOpenAll(t *testing.T) {

}

func TestUserAchievements(t *testing.T) {

}

func TestMarkShown(t *testing.T) {

}
