package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/user/transport/http/statistics"
)

func StatisticsRouting(s *setUpper) {
	const (
		TIME      = "/time"
		MONEY     = "/money"
		CIGARETTE = "/cigarette"
	)
	ctx := context.Background()

	r := s.config.Mux.PathPrefix("/statistics").Subrouter()

	handler := statistics.New(s.StatisticsUseCase())

	r.Handle(TIME, handler.TimeStatisticsHandler(ctx)).Methods(http.MethodGet)
	r.Handle(MONEY, handler.MoneyStatisticsHandler(ctx)).Methods(http.MethodGet)
	r.Handle(CIGARETTE, handler.CigaretteStatisticsHandler(ctx)).Methods(http.MethodGet)
}
