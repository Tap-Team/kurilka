package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/achievements/transport/http/achievement"
)

func AchievementRouting(s *setUpper) {
	ctx := context.Background()
	const (
		OPEN_SINGLE = "/open-single"
		MARK_SHOWN  = "/mark-shown"
		GET         = ""
	)
	cnf := s.cnf
	useCase := s.AchievementUseCase()

	handler := achievement.NewAchievementHandler(useCase)

	r := cnf.Mux.PathPrefix("/achievements").Subrouter()

	r.Handle(GET, handler.UserAchievementsHandler(ctx)).Methods(http.MethodGet)
	r.Handle(MARK_SHOWN, handler.MarkShownHandler(ctx)).Methods(http.MethodPost)
	r.Handle(OPEN_SINGLE, handler.OpenSingleHandler(ctx)).Methods(http.MethodPost)
}
