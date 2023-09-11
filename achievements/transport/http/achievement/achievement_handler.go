package achievement

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "achievements/transport"

type AchievementHandler struct {
	useCase achievementusecase.AchievementUseCase
}

func NewAchievementHandler(useCase achievementusecase.AchievementUseCase) *AchievementHandler {
	return &AchievementHandler{useCase: useCase}
}

var (
	ErrParseAchievementId error = errors.New("failed parse 'achievementId' query")
)

type query struct {
	url.Values
}

func (q *query) AchievementId() (int64, error) {
	achievementId, err := strconv.ParseInt(q.Get("achievementId"), 10, 64)
	if err != nil {
		return 0, ErrParseAchievementId
	}
	return achievementId, err
}

func NewQuery(uval url.Values) *query {
	return &query{uval}
}

//	@Summary		OpenSingle
//	@Description	open single achievement by user and achievement ids
//	@Tags			achievements
//	@Produce		json
//	@Param			achievementId	query	int64	true	"achievement id"
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/achievements/open-single [post]
func (h *AchievementHandler) OpenSingleHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "OpenSingleHandler", _PROVIDER)))
			return
		}
		query := NewQuery(r.URL.Query())
		achievementId, err := query.AchievementId()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse achievemetnId", "OpenSingleHandler", _PROVIDER)))
			return
		}
		openAchievementResponse, err := h.useCase.OpenSingle(ctx, userId, achievementId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("open single achievement", "OpenSingleHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, openAchievementResponse, http.StatusOK)
	}

	return http.HandlerFunc(handler)
}

// func (h *AchievementHandler) OpenTypeHandler(ctx context.Context) http.Handler {
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		userId, err := httphelpers.VKID(r)
// 		if err != nil {
// 			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "OpenTypeHandler", _PROVIDER)))
// 			return
// 		}
// 		achievementType := achievementmodel.AchievementType(r.URL.Query().Get("achievementType"))
// 		openAchievementResponse, err := h.useCase.OpenType(ctx, userId, achievementType)
// 		if err != nil {
// 			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("open achievement by type", "OpenTypeHandler", _PROVIDER)))
// 			return
// 		}
// 		httphelpers.WriteJSON(w, openAchievementResponse, http.StatusOK)
// 	}
// 	return http.HandlerFunc(handler)
// }

// func (h *AchievementHandler) OpenAllHandler(ctx context.Context) http.Handler {
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		userId, err := httphelpers.VKID(r)
// 		if err != nil {
// 			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "OpenAllHandler", _PROVIDER)))
// 			return
// 		}
// 		openAchievementResponse, err := h.useCase.OpenAll(ctx, userId)
// 		if err != nil {
// 			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("open all achievements", "OpenAllHandler", _PROVIDER)))
// 			return
// 		}
// 		httphelpers.WriteJSON(w, openAchievementResponse, http.StatusOK)
// 	}
// 	return http.HandlerFunc(handler)
// }

//	@Summary		UserAchievements
//	@Description	return all achievements that exists, if user not reach achievement reachDate is zero, is user not open achievement openDate is zero
//	@Tags			achievements
//	@Produce		json
//	@Success		200	{array}		achievementmodel.Achievement
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/achievements [get]
func (h *AchievementHandler) UserAchievementsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "UserAchievementsHandler", _PROVIDER)))
			return
		}
		achievements, err := h.useCase.UserAchievements(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get user achievements", "UserAchievementsHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, achievements, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

//	@Summary		MarkShown
//	@Description	set on all reach achievements show = true
//	@Tags			achievements
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/achievements/mark-shown [post]
func (h *AchievementHandler) MarkShownHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "MarkShownHandler", _PROVIDER)))
			return
		}
		err = h.useCase.MarkShown(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("mark shown", "MarkShownHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}
