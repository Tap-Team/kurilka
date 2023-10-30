package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/internal/middleware"
	"github.com/Tap-Team/kurilka/user/transport/http/user"
)

func UserRouting(setUpper *setUpper) {
	ctx := context.Background()
	const (
		GET     = "/user"
		CREATE  = "/create"
		RESET   = "/reset"
		LEVEL   = "/level"
		FRIENDS = "/friends"
		EXISTS  = "/exists"
	)

	config := setUpper.Config()

	useCase := setUpper.UserUseCase()

	transport := user.NewUserTransport(useCase)

	r := config.Mux.PathPrefix("/users").Subrouter()

	r.Handle(GET, transport.GetUserHandler(ctx)).
		Methods(http.MethodGet)
	r.Handle(RESET, transport.ResetUserHandler(ctx)).
		Methods(http.MethodDelete)
	r.Handle(LEVEL, transport.GetUserLevelHandler(ctx)).
		Methods(http.MethodGet)
	r.Handle(EXISTS, transport.UserExistsHandler(ctx)).
		Methods(http.MethodGet)

	jsonRoute := r.NewRoute().Subrouter()
	jsonRoute.Use(middleware.JSON)

	jsonRoute.Handle(FRIENDS, transport.FriendsHandler(ctx)).Methods(http.MethodGet)
	jsonRoute.Handle(CREATE, transport.CreateUserHandler(ctx)).Methods(http.MethodPost)

}
