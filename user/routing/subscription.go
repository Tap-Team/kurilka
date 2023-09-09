package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/user/transport/http/subscription"
)

func SubscriptionRouting(s *setUpper) {

	const (
		GET_USER = "/user"
	)
	ctx := context.Background()
	useCase := s.SubscriptionUseCase()

	handler := subscription.New(useCase)

	r := s.config.Mux.PathPrefix("/subscription").Subrouter()

	r.Handle(GET_USER, handler.UserSubscription(ctx)).Methods(http.MethodGet)
}
