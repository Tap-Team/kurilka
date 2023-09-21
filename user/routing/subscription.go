package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/user/transport/http/subscription"
	"github.com/Tap-Team/kurilka/user/transport/http/subscriptionadmin"
)

func SubscriptionRouting(s *setUpper) {

	const (
		GET_USER = "/user"
		UPDATE   = "/update"
	)
	ctx := context.Background()
	useCase := s.SubscriptionUseCase()

	handler := subscription.New(useCase)

	r := s.config.Mux.PathPrefix("/subscription").Subrouter()

	r.Handle(GET_USER, handler.UserSubscription(ctx)).Methods(http.MethodGet)

	adminHandler := subscriptionadmin.New(s.SubscriptionManager())

	r.Handle(UPDATE, adminHandler.UpdateUserSubscriptionHandler(ctx)).Methods(http.MethodPut)
}
