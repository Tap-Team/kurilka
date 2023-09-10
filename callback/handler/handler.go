package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Tap-Team/kurilka/callback/handler/donutcreate"
	"github.com/Tap-Team/kurilka/callback/handler/donutexpired"
	"github.com/Tap-Team/kurilka/callback/handler/donutprolonged"
	"github.com/Tap-Team/kurilka/callback/usecase/subscriptionusecase"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"golang.org/x/exp/slog"
)

const _HANDLER_PROVIDER = "callback/handler.Handler"

type EventType string

const (
	CONFIRMATION    EventType = "confirmation"
	DONUT_CREATE    EventType = "donut_subscription_create"
	DONUT_EXPIRED   EventType = "donut_subscription_expired"
	DONUT_PROLONGED EventType = "donut_subscription_prolonged"
)

type Event struct {
	Type    EventType       `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
	EventID string          `json:"event_id"`
	V       string          `json:"v"`
	Secret  string          `json:"secret"`
}

type EventHandler interface {
	HandleEvent(ctx context.Context, object json.RawMessage) error
}

type Handler struct {
	handlers    map[EventType]EventHandler
	confirmCode string
	groupId     int64
	secret      string
}

func New(confirmCode string, groupId int64, secret string, useCase subscriptionusecase.UseCase) *Handler {
	donutCreateHandler := donutcreate.New(useCase)
	donutExpiredHandler := donutexpired.New(useCase)
	donutProlongedHandler := donutprolonged.New(useCase)
	return &Handler{
		confirmCode: confirmCode,
		groupId:     groupId,
		secret:      secret,
		handlers: map[EventType]EventHandler{
			DONUT_CREATE:    donutCreateHandler,
			DONUT_EXPIRED:   donutExpiredHandler,
			DONUT_PROLONGED: donutProlongedHandler,
		},
	}
}

type nilEventHandler struct{}

func (nilEventHandler) HandleEvent(ctx context.Context, object json.RawMessage) error {
	return nil
}

func (h *Handler) EventHandler(etype EventType) EventHandler {
	handler, ok := h.handlers[etype]
	if ok {
		return handler
	}
	slog.Info("event handler not found", "type", etype)
	return nilEventHandler{}
}

func ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) confirm(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(h.confirmCode))
}

func ip(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

var (
	ErrWrongSecret = errors.New("wrong secret")
)

func (h *Handler) CallBackHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			err = exception.Wrap(err, exception.NewCause("decode json", "CallBackHandler", _HANDLER_PROVIDER))
			httphelpers.Error(w, err)
			return
		}
		if event.Type == CONFIRMATION {
			h.confirm(w)
			return
		}
		if event.Secret != h.secret {
			slog.InfoContext(ctx, "wrong secret send", "user-agent", r.UserAgent(), "ip", ip(r))
			httphelpers.Error(w, ErrWrongSecret)
			return
		}
		err = h.EventHandler(event.Type).HandleEvent(ctx, event.Object)
		if err != nil {
			err = exception.Wrap(err, exception.NewCause(fmt.Sprintf("handle event %s", event.Type), "CallBackHandler", _HANDLER_PROVIDER))
			httphelpers.Error(w, err)
			return
		}
		ok(w)
	}
	return http.HandlerFunc(handler)
}
