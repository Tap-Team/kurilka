package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/user/transport/http/trigger"
)

func TriggerRouting(setUpper *setUpper) {
	ctx := context.Background()
	const (
		REMOVE = "/remove"
		ADD    = "/add"
	)

	config := setUpper.Config()

	manager := setUpper.TriggerManager()

	transport := trigger.NewTriggerTransport(manager)

	r := config.Mux.PathPrefix("/triggers").Subrouter()

	r.Handle(REMOVE, transport.RemoveTriggerHandler(ctx)).Methods(http.MethodDelete)
	r.Handle(ADD, transport.AddTriggerHandler(ctx)).Methods(http.MethodPost)
}
