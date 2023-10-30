package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/vote/model/notificationtype"
	"github.com/Tap-Team/kurilka/vote/model/subscription"
)

//go:generate mockgen -source handler.go -destination handler_mocks.go -package handler

type GetSubscriptionUseCase interface {
	Subscription(ctx context.Context, subscriptionId string) (subscription.Subscription, error)
}

type ChangeSubscriptionStatusUseCase interface {
	ChangeSubscriptionStatus(ctx context.Context, changeSubscriptionStatus subscription.ChangeSubscriptionStatus) (resp subscription.ChangeSubscriptionStatusResponse, err error)
}

type SubscriptionStorage interface {
	UserSubscriptionId(ctx context.Context, userId int64) (subscriptionId int64, err error)
}

type Handler struct {
	getSubscriptionUseCase          GetSubscriptionUseCase
	changeSubscriptionStatusUseCase ChangeSubscriptionStatusUseCase
	subscriptionStorage             SubscriptionStorage
	VKAppSecret                     string
}

func New(
	vkAppSecret string,
	getSubscriptionUseCase GetSubscriptionUseCase,
	changeSubscriptionStatusUseCase ChangeSubscriptionStatusUseCase,
	subscriptionStorage SubscriptionStorage,
) *Handler {
	return &Handler{
		getSubscriptionUseCase:          getSubscriptionUseCase,
		changeSubscriptionStatusUseCase: changeSubscriptionStatusUseCase,
		subscriptionStorage:             subscriptionStorage,
		VKAppSecret:                     vkAppSecret,
	}
}

func (h *Handler) GetSubscriptionHandler(params Params) http.Handler {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		subscriptionId := params.GetString("item")
		subscription, err := h.getSubscriptionUseCase.Subscription(ctx, subscriptionId)
		if err != nil {
			Error(w, err)
			return
		}
		WriteJSON(w, subscription, http.StatusOK)
	}
	return handler
}

func (h *Handler) SubscriptionStatusChangeHandler(params Params) http.Handler {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Printf("%v", params)
		changeSubscriptionStatus := subscription.NewChangeSubscriptionStatus(
			params.GetInt("subscription_id"),
			params.GetInt("user_id"),
			params.GetString("item_id"),
			subscription.SubscriptionStatus(params.GetString("status")),
			subscription.CancelReason(params.GetString("cancel_reason")),
		)
		resp, err := h.changeSubscriptionStatusUseCase.ChangeSubscriptionStatus(ctx, changeSubscriptionStatus)
		if err != nil {
			Error(w, err)
			return
		}
		WriteJSON(w, resp, http.StatusOK)
	}
	return handler
}

func (h *Handler) HandleNotificationHandler() http.Handler {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := make(Params)
		params.ReadFrom(r.Body)
		r.Body.Close()
		ntype := params.NotificationType()
		switch {
		case notificationtype.GET_SUBSCRIPTION.Is(ntype):
			h.GetSubscriptionHandler(params).ServeHTTP(w, r)
		case notificationtype.SUBSCRIPTION_STATUS_CHANGE.Is(ntype):
			h.SubscriptionStatusChangeHandler(params).ServeHTTP(w, r)
		}
	})
	handler = VerifyBodySigMiddleware(handler, h.VKAppSecret)
	return handler
}

func (h *Handler) UserSubscriptionIdHandler() http.Handler {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, err)
			return
		}
		subscriptionId, err := h.subscriptionStorage.UserSubscriptionId(ctx, userId)
		if err != nil {
			httphelpers.Error(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d", subscriptionId)
	}
	return handler
}
