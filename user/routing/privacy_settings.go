package routing

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/user/transport/http/privacysetting"
)

func PrivacySettingRouting(setUpper *setUpper) {
	ctx := context.Background()
	const (
		REMOVE = "/remove"
		ADD    = "/add"
		GET    = "/usersettings"
	)

	config := setUpper.Config()

	manager := setUpper.PrivacySettingManager()

	transport := privacysetting.NewPrivacySettingTransport(manager)

	r := config.Mux.PathPrefix("/privacysettings").Subrouter()

	r.Handle(GET, transport.GetPrivacySettingsHandler(ctx)).Methods(http.MethodGet)
	r.Handle(ADD, transport.AddRemovePrivacySettingHandler(ctx)).Methods(http.MethodPost)
	r.Handle(REMOVE, transport.RemovePrivacySettingHandler(ctx)).Methods(http.MethodDelete)

}
