package subscriptionadmin

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source subscriptionadmin.go -destination mocks.go -package subscriptionadmin

const _PROVIDER = "user/transport/http/subscriptionadmin.Handler"

var (
	ErrFailedParseExpired = errors.New("failed parse expired as int")
)

type UserSubscriptionUpdater interface {
	UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error
}

type Handler struct {
	updater UserSubscriptionUpdater
}

func New(updater UserSubscriptionUpdater) *Handler {
	return &Handler{updater: updater}
}

type query struct {
	url.Values
}

func (q query) Subscription() (usermodel.Subscription, error) {
	tp := q.Get("subscriptionType")
	expiredString := q.Get("expired")
	expired, err := strconv.ParseInt(expiredString, 10, 64)
	if err != nil {
		return usermodel.Subscription{}, ErrFailedParseExpired
	}
	subscriptionType := usermodel.SubscriptionType(tp)
	if err := subscriptionType.Validate(); err != nil {
		return usermodel.Subscription{}, err
	}
	expiredTime := time.Unix(expired, 0)
	return usermodel.NewSubscription(subscriptionType, expiredTime), nil
}

// UpdateUserSubscriptionHandler godoc
//
//	@Summary		UpdateUserSubscription
//	@Description	manual update user subscription (only admin)
//	@Tags			subscription
//	@Produce		json
//	@Param			subscriptionType	query	usermodel.SubscriptionType	true	"subscription type"
//	@Param			expired				query	int64						true	"time when subscription expired"
//	@Success		204
//	@Falure			400 {object} errormodel.ErrorResponse
//	@Router			/subscription/update [put]
func (h *Handler) UpdateUserSubscriptionHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk id", "UpdateUserSubscriptionHandler", _PROVIDER)))
			return
		}
		subscription, err := query{r.URL.Query()}.Subscription()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse subscription", "UpdateUserSubscriptionHandler", _PROVIDER)))
			return
		}
		err = h.updater.UpdateUserSubscription(ctx, userId, subscription)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("update user subscription", "UpdateUserSubscriptionHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}
