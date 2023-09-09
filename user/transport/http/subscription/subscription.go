package subscription

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/usecase/subscriptionusecase"
)

const _PROVIDER = "user/transport/http/subscription"

type SubscriptionHandler struct {
	useCase subscriptionusecase.SubscriptionUseCase
}

func New(useCase subscriptionusecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{useCase: useCase}
}

// UserSubscription godoc
//
//	@Summary		UserSubscription
//	@Description	get user subscription type
//	@Tags			subscription
//	@Produce		json
//	@Param			vk_user_id	query		int64	true	"vk user id"
//	@Success		200			{object}	usermodel.SubscriptionType
//	@Failure		400			{object}	errormodel.ErrorResponse
//	@Router			/subscription/user [get]
func (s *SubscriptionHandler) UserSubscription(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get vk user id", "UserSubscription", _PROVIDER)))
			return
		}
		subscriptionType := s.useCase.UserSubscription(ctx, userId)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(subscriptionType))
	}
	return http.HandlerFunc(handler)
}
