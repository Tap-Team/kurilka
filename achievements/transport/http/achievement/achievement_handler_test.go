package achievement_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/transport/http/achievement"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestOpenSingleHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievementUseCase := achievementusecase.NewMockAchievementUseCase(ctrl)

	handler := achievement.NewAchievementHandler(achievementUseCase)

	cases := []struct {
		userId        int64
		achievementId int64
		queryValues   map[string]string

		useCaseResponse *model.OpenAchievementResponse
		useCaseError    error
		useCaseCall     bool

		err error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			queryValues: map[string]string{
				"vk_user_id": "10",
			},
			err:        achievement.ErrParseAchievementId,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId:        1,
			achievementId: 2,
			queryValues: map[string]string{
				"vk_user_id":    "1",
				"achievementId": "2",
			},

			useCaseResponse: model.NewOpenAchievementResponse(time.Now()),
			useCaseError:    nil,
			useCaseCall:     true,

			statusCode: 200,
		},
	}

	for _, cs := range cases {
		if cs.useCaseCall {
			achievementUseCase.EXPECT().OpenSingle(gomock.Any(), cs.userId, cs.achievementId).Return(cs.useCaseResponse, cs.useCaseError).Times(1)
		}
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/achievements/open-single?"+urlValues.Encode(), nil)
		handler.OpenSingleHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, cs.statusCode, rec.Result().StatusCode, "status code not equal")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var response model.OpenAchievementResponse
			err := json.NewDecoder(rec.Result().Body).Decode(&response)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode result body")

			assert.Equal(t, response.OpenTime.Unix(), cs.useCaseResponse.OpenTime.Unix(), "wrong response")
		}
	}
}

func TestUserAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievementUseCase := achievementusecase.NewMockAchievementUseCase(ctrl)

	handler := achievement.NewAchievementHandler(achievementUseCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		useCaseResponse []*achievementmodel.Achievement
		useCaseErr      error
		useCaseCall     bool

		err error

		statusCode int
	}{}

	for _, cs := range cases {
		if cs.useCaseCall {
			achievementUseCase.EXPECT().UserAchievements(gomock.Any(), cs.userId).Return(cs.useCaseResponse, cs.useCaseErr).Times(1)
		}
		uval := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			uval.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/achievements/userachievements?"+uval.Encode(), nil)
		handler.UserAchievementsHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, cs.statusCode, rec.Result().StatusCode, "status code not equal")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var response []*achievementmodel.Achievement
			err := json.NewDecoder(rec.Result().Body).Decode(&response)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode body")

			equal := slices.EqualFunc(response, cs.useCaseResponse, compareAchievements)
			assert.Equal(t, true, equal, "response not equal")
		}
	}
}

func TestMarkShown(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievementUseCase := achievementusecase.NewMockAchievementUseCase(ctrl)

	handler := achievement.NewAchievementHandler(achievementUseCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		useCaseCall bool
		useCaseErr  error

		err error

		statusCode int
	}{}

	for _, cs := range cases {
		if cs.useCaseCall {
			achievementUseCase.EXPECT().MarkShown(gomock.Any(), cs.userId).Return(cs.useCaseErr).Times(1)
		}
		uval := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			uval.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/achievements/mark-shown?"+uval.Encode(), nil)
		handler.MarkShownHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}
}
