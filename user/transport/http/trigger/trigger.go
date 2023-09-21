package trigger

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/triggerdatamanager"
)

const _PROVIDER = "user/transport/http/trigger"

type TriggerTransport struct {
	trigger triggerdatamanager.TriggerManager
}

func NewTriggerTransport(trigger triggerdatamanager.TriggerManager) *TriggerTransport {
	return &TriggerTransport{trigger: trigger}
}

type query url.Values

func (q query) Trigger() usermodel.Trigger {
	trigger := url.Values(q).Get("trigger")
	return usermodel.Trigger(trigger)
}

// RemoveTriggerHandler godoc
//
//	@Summary		RemoveTrigger
//	@Description	remove user trigger, if user not exists, or trigger has been removed return error
//	@Tags			triggers
//	@Produce		json
//	@Param			trigger	query	usermodel.Trigger	true	"trigger"
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/triggers/remove [delete]
func (t *TriggerTransport) RemoveTriggerHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse id", "RemoveTriggerHandler", _PROVIDER)))
			return
		}
		trigger := query(r.URL.Query()).Trigger()
		err = trigger.Validate()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("validate trigger", "RemoveTriggerHandler", _PROVIDER)))
			return
		}
		err = t.trigger.Remove(ctx, userId, trigger)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("remove trigger", "RemoveTriggerHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}

// RemoveTriggerHandler godoc
//
//	@Summary		RemoveTrigger
//	@Description	remove user trigger, if user not exists, or trigger has been removed return error
//	@Tags			triggers
//	@Produce		json
//	@Param			trigger	query	usermodel.Trigger	true	"trigger"
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/triggers/add [post]
func (t *TriggerTransport) AddTriggerHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse id", "AddTriggerHandler", _PROVIDER)))
			return
		}
		trigger := query(r.URL.Query()).Trigger()
		err = trigger.Validate()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("validate trigger", "AddTriggerHandler", _PROVIDER)))
			return
		}
		err = t.trigger.Add(ctx, userId, trigger)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("add trigger", "AddTriggerHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}
