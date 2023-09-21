package statistics_test

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/user/model"
	"github.com/Tap-Team/kurilka/user/transport/http/statistics"
	"github.com/Tap-Team/kurilka/user/usecase/statisticsusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_TimeStatisticsHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := statisticsusecase.NewMockStatisticsUseCase(ctrl)

	handler := statistics.New(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		useCaseCall     bool
		useCaseResponse model.IntUserStatistics
		useCaseErr      error

		response   model.IntUserStatistics
		statusCode int
		err        error
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id": "123",
			},
			useCaseCall: true,
			useCaseErr:  usererror.ExceptionUserNotFound(),

			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			userId: 321,
			queryValues: map[string]string{
				"vk_user_id": "321",
			},
			useCaseCall:     true,
			useCaseResponse: model.NewIntUserStatistics(100),

			statusCode: http.StatusOK,
			response:   model.NewIntUserStatistics(100),
		},
	}

	for _, cs := range cases {
		if cs.useCaseCall {
			useCase.EXPECT().TimeStatistics(gomock.Any(), cs.userId).Return(cs.useCaseResponse, cs.useCaseErr).Times(1)
		}
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/statistics/time?"+urlValues.Encode(), nil)

		handler.TimeStatisticsHandler(ctx).ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode
		assert.Equal(t, statusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var response model.IntUserStatistics
			err := json.NewDecoder(rec.Result().Body).Decode(&response)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode body")

			assert.Equal(t, cs.response, response, "response not equal")
		}
	}
}

func floatStatisticsUnitEqual(u1, u2 model.FloatStatisticsUnit) bool {
	return math.Abs(float64(u1-u2)) < 0.01
}

func floatStatisticsEqual(s1, s2 model.FloatUserStatistics) bool {
	return floatStatisticsUnitEqual(s1.Day, s2.Day) &&
		floatStatisticsUnitEqual(s1.Week, s2.Week) &&
		floatStatisticsUnitEqual(s1.Month, s2.Month) &&
		floatStatisticsUnitEqual(s1.Year, s1.Year)
}

func Test_MoneyStatisticsHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := statisticsusecase.NewMockStatisticsUseCase(ctrl)

	handler := statistics.New(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		useCaseCall     bool
		useCaseResponse model.FloatUserStatistics
		useCaseErr      error

		response   model.FloatUserStatistics
		statusCode int
		err        error
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id": "123",
			},
			useCaseCall: true,
			useCaseErr:  usererror.ExceptionUserNotFound(),

			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			userId: 321,
			queryValues: map[string]string{
				"vk_user_id": "321",
			},
			useCaseCall:     true,
			useCaseResponse: model.NewFloatUserStatisctics(172.98),

			statusCode: http.StatusOK,
			response:   model.NewFloatUserStatisctics(172.98),
		},
	}

	for _, cs := range cases {
		if cs.useCaseCall {
			useCase.EXPECT().MoneyStatistics(gomock.Any(), cs.userId).Return(cs.useCaseResponse, cs.useCaseErr).Times(1)
		}
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/statistics/time?"+urlValues.Encode(), nil)

		handler.MoneyStatisticsHandler(ctx).ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode
		assert.Equal(t, statusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var response model.FloatUserStatistics
			err := json.NewDecoder(rec.Result().Body).Decode(&response)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode body")

			assert.Equal(t, true, floatStatisticsEqual(response, cs.response), "response not equal")
		}
	}
}

func Test_CigaretteStatisticsHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := statisticsusecase.NewMockStatisticsUseCase(ctrl)

	handler := statistics.New(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		useCaseCall     bool
		useCaseResponse model.IntUserStatistics
		useCaseErr      error

		response   model.IntUserStatistics
		statusCode int
		err        error
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id": "123",
			},
			useCaseCall: true,
			useCaseErr:  usererror.ExceptionUserNotFound(),

			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			userId: 321,
			queryValues: map[string]string{
				"vk_user_id": "321",
			},
			useCaseCall:     true,
			useCaseResponse: model.NewIntUserStatistics(100),

			statusCode: http.StatusOK,
			response:   model.NewIntUserStatistics(100),
		},
	}

	for _, cs := range cases {
		if cs.useCaseCall {
			useCase.EXPECT().CigaretteStatistics(gomock.Any(), cs.userId).Return(cs.useCaseResponse, cs.useCaseErr).Times(1)
		}
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/statistics/time?"+urlValues.Encode(), nil)

		handler.CigaretteStatisticsHandler(ctx).ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode
		assert.Equal(t, statusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var response model.IntUserStatistics
			err := json.NewDecoder(rec.Result().Body).Decode(&response)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode body")

			assert.Equal(t, cs.response, response, "response not equal")
		}
	}
}
