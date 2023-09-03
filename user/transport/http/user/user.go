package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/usecase/userusecase"
)

const _PROVIDER = "user/transport/http/user"

type UserTransport struct {
	userUseCase userusecase.UserUseCase
}

func NewUserTransport(userUseCase userusecase.UserUseCase) *UserTransport {
	return &UserTransport{
		userUseCase: userUseCase,
	}
}

// GetUserHandler godoc
//
//	@Summary		GetUser
//	@Description	get user by vk_user_id
//	@Tags			users
//	@Produce		json
//	@Param			vk_user_id	query		int64	true	"vk user id"
//	@Success		200			{object}	usermodel.User
//	@Failure		400			{object}	errormodel.ErrorResponse
//	@Router			/users/user [get]
func (t *UserTransport) GetUserHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "GetUserHandler", _PROVIDER)))
			return
		}
		user, err := t.userUseCase.User(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get user", "GetUserHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, user, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

// UserExistsHandler godoc
//
//	@Summary		UserExists
//	@Description	check user exists
//	@Tags			users
//	@Produce		json
//	@Param			vk_user_id	query		int64	true	"vk user id"
//	@Success		200			{object}	bool
//	@Failure		400			{object}	bool
//	@Router			/users/exists [get]
func (t *UserTransport) UserExistsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			err := exception.Wrap(err, exception.NewCause("parse vk_user_id", "GetUserHandler", _PROVIDER))
			slog.ErrorContext(ctx, err.Error(), "url", r.URL.String())
			httphelpers.WriteJSON(w, false, http.StatusInternalServerError)
			return
		}
		_, err = t.userUseCase.User(ctx, userId)
		if errors.Is(err, usererror.ExceptionUserNotFound()) {
			httphelpers.WriteJSON(w, false, http.StatusOK)
			return
		}
		if err == nil {
			httphelpers.WriteJSON(w, true, http.StatusOK)
			return
		}
		slog.Error("failed get user from storage", "userId", userId)
		httphelpers.WriteJSON(w, false, http.StatusInternalServerError)
	}
	return http.HandlerFunc(handler)
}

// ResetUserHandler godoc
//
//	@Summary		ResetUser
//	@Description	"reset user data"
//	@Tags			users
//	@Produce		json
//	@Param			vk_user_id	query	int64	true	"vk user id"
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/users/reset [delete]
func (t *UserTransport) ResetUserHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "ResetUserHandler", _PROVIDER)))
			return
		}
		err = t.userUseCase.Reset(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("reset user", "ResetUserHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}

// CreateUserHandler godoc
//
//	@Summary		CreateUser
//	@Description	"create user"
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			vk_user_id	query		int64	true	"vk user id"
//	@Success		201			{object}	usermodel.User
//	@Failure		400			{object}	errormodel.ErrorResponse
//	@Router			/users/create [post]
func (t *UserTransport) CreateUserHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "CreateUserHandler", _PROVIDER)))
			return
		}
		var createUser usermodel.CreateUser
		err = json.NewDecoder(r.Body).Decode(&createUser)
		r.Body.Close()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("decode body", "CreateUserHandler", _PROVIDER)))
			return
		}
		err = createUser.Validate()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("validate body", "CreateUserHandler", _PROVIDER)))
			return
		}
		user, err := t.userUseCase.Create(ctx, userId, &createUser)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("create user", "CreateUserHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, user, http.StatusCreated)
	}
	return http.HandlerFunc(handler)
}

// GetUserLevelHandler godoc
//
//	@Summary		UserLevel
//	@Description	"get user level from postgres, update cache and return level"
//	@Tags			users
//	@Produce		json
//	@Param			vk_user_id	query		int64	true	"vk user id"
//	@Success		200			{object}	usermodel.LevelInfo
//	@Failure		400			{object}	errormodel.ErrorResponse
//	@Router			/users/level [get]
func (t *UserTransport) GetUserLevelHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk_user_id", "GetUserLevelHandler", _PROVIDER)))
			return
		}
		level, err := t.userUseCase.Level(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get user level", "GetUserLevelHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, level, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

// FriendsHandler godoc
//
//	@Summary		Friends
//	@Description	"map friends id to dto list"
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			friends	body		[]int64	true	"list of user friends"
//	@Success		200		{array}		usermodel.Friend
//	@Failure		400		{object}	errormodel.ErrorResponse
//	@Router			/users/friends [get]
func (t *UserTransport) FriendsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		friendsIds := make([]int64, 0)
		err := json.NewDecoder(r.Body).Decode(&friendsIds)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("decode body", "FriendsHandler", _PROVIDER)))
			return
		}
		friends := t.userUseCase.Friends(ctx, friendsIds)
		httphelpers.WriteJSON(w, friends, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}
