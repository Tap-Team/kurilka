package transport

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "achievements/transport"

type AchievementUseCase interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error)
	OpenType(ctx context.Context, userId int64, achtype achievementmodel.AchievementType) (*model.OpenAchievementResponse, error)
	OpenAll(ctx context.Context, userId int64) (*model.OpenAchievementResponse, error)
	MarkShown(ctx context.Context, userId int64) error
}

type AchievementHandler struct {
	useCase AchievementUseCase
}

func NewAchievementHandler(useCase AchievementUseCase) *AchievementHandler {
	return &AchievementHandler{useCase: useCase}
}

func (h *AchievementHandler) OpenSingleHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "OpenSingleHandler", _PROVIDER)))
			return
		}
		achievementId, err := strconv.ParseInt(r.URL.Query().Get("achievementId"), 10, 64)
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

func (h *AchievementHandler) OpenTypeHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "OpenTypeHandler", _PROVIDER)))
			return
		}
		achievementType := achievementmodel.AchievementType(r.URL.Query().Get("achievementType"))
		openAchievementResponse, err := h.useCase.OpenType(ctx, userId, achievementType)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("open achievement by type", "OpenTypeHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, openAchievementResponse, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

func (h *AchievementHandler) OpenAllHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "OpenAllHandler", _PROVIDER)))
			return
		}
		openAchievementResponse, err := h.useCase.OpenAll(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("open all achievements", "OpenAllHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, openAchievementResponse, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

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
