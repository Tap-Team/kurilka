package statistics

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/usecase/statisticsusecase"
)

const _PROVIDER = "user/transport/http/statistics.StatisticsHandler"

type StatisticsHandler struct {
	statistics statisticsusecase.StatisticsUseCase
}

func New(statistics statisticsusecase.StatisticsUseCase) *StatisticsHandler {
	return &StatisticsHandler{statistics: statistics}
}

// TimeStatisticsHandler godoc
//
//	@Summary		TimeStatistics
//	@Description	get time statistics
//	@Tags			statistics
//	@Produce		json
//	@Success		200	{object}	model.IntUserStatistics
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/statistics/time [get]
func (s *StatisticsHandler) TimeStatisticsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get vk user id", "TimeStatisticsHandler", _PROVIDER)))
			return
		}
		statistics, err := s.statistics.TimeStatistics(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get time user statistics", "TimeStatisticsHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, statistics, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

// CigaretteStatisticsHandler godoc
//
//	@Summary		CigaretteStatistics
//	@Description	get cigarette statistics
//	@Tags			statistics
//	@Produce		json
//	@Success		200	{object}	model.IntUserStatistics
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/statistics/cigarette [get]
func (s *StatisticsHandler) CigaretteStatisticsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get vk user id", "CigaretteStatisticsHandler", _PROVIDER)))
			return
		}
		statistics, err := s.statistics.CigaretteStatistics(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get cigarette user statistics", "CigaretteStatisticsHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, statistics, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

// MoneyStatisticsHandler godoc
//
//	@Summary		MoneyStatistics
//	@Description	get money statistics
//	@Tags			statistics
//	@Produce		json
//	@Success		200	{object}	model.FloatUserStatistics
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/statistics/money [get]
func (s *StatisticsHandler) MoneyStatisticsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get vk user id", "MoneyStatisticsHandler", _PROVIDER)))
			return
		}
		statistics, err := s.statistics.MoneyStatistics(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get money user statistics", "MoneyStatisticsHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, statistics, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}
