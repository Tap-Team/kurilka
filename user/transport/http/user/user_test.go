package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/validate"
	"github.com/Tap-Team/kurilka/user/transport/http/user"
	"github.com/Tap-Team/kurilka/user/usecase/userusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_UserGetHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := userusecase.NewMockUserUseCase(ctrl)

	transport := user.NewUserTransport(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		user usermodel.User
		err  error

		statusCode int
	}{

		{
			queryValues: map[string]string{
				"figna": "figna",
			},
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},

		{
			userId: 101,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(101),
			},
			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},

		{
			userId: 100000,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(100000),
			},
			user:       random.StructTyped[usermodel.User](),
			statusCode: http.StatusOK,
		},
	}

	for _, cs := range cases {

		if cs.userId != 0 {
			useCase.EXPECT().User(gomock.Any(), cs.userId).Return(&cs.user, cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/user?"+urlValues.Encode(), nil)
		transport.GetUserHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var user usermodel.User
			err := json.NewDecoder(rec.Result().Body).Decode(&user)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode body")

			fields, ok := compareUsers(user, cs.user)
			assert.Equal(t, true, ok, "fields %s not equal", fields)
		}
	}
}

func Test_UserExistsHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := userusecase.NewMockUserUseCase(ctrl)

	transport := user.NewUserTransport(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		exists bool
		err    error

		statusCode int
	}{
		{
			exists:     false,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 1,
			exists: false,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(1),
			},
			err: usererror.ExceptionUserNotFound(),

			statusCode: http.StatusOK,
		},
		{
			userId: 1,
			exists: true,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(1),
			},
			statusCode: http.StatusOK,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			useCase.EXPECT().User(gomock.Any(), cs.userId).Return(nil, cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/exists?"+urlValues.Encode(), nil)
		transport.UserExistsHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		var result bool
		err := json.NewDecoder(rec.Result().Body).Decode(&result)
		rec.Result().Body.Close()
		assert.NilError(t, err, "failed decode body")

		assert.Equal(t, cs.exists, result, "result not equal exists")
	}

}

func Test_UserLevelHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := userusecase.NewMockUserUseCase(ctrl)

	transport := user.NewUserTransport(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		level usermodel.LevelInfo
		err   error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 100,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(100),
			},
			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			userId: 100,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(100),
			},
			err:        errors.New("random error"),
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 101,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(101),
			},
			level:      random.StructTyped[usermodel.LevelInfo](),
			statusCode: http.StatusOK,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			useCase.EXPECT().Level(gomock.Any(), cs.userId).Return(&cs.level, cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/level?"+urlValues.Encode(), nil)
		transport.GetUserLevelHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			var level usermodel.LevelInfo
			err := json.NewDecoder(rec.Result().Body).Decode(&level)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode request body")

			assert.Equal(t, cs.level, level, "level not equal")
		}
	}
}

func Test_ResetUserHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := userusecase.NewMockUserUseCase(ctrl)

	transport := user.NewUserTransport(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		err error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 1,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(1),
			},
			err:        usererror.ExceptionUserNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			userId: 12,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(12),
			},
			err:        errors.New("error"),
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(123),
			},
			statusCode: http.StatusNoContent,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			useCase.EXPECT().Reset(gomock.Any(), cs.userId).Return(cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/reset?"+urlValues.Encode(), nil)
		transport.ResetUserHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}

}

func Test_CreateUserHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := userusecase.NewMockUserUseCase(ctrl)

	transport := user.NewUserTransport(useCase)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		createUser usermodel.CreateUser

		user usermodel.User
		err  error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			userId: 1,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(1),
			},
			createUser: random.StructTyped[usermodel.CreateUser](),
			err:        usererror.ExceptionUserExist(),
			statusCode: http.StatusBadRequest,
		},
		{
			userId: 12,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(12),
			},
			createUser: random.StructTyped[usermodel.CreateUser](),
			user:       random.StructTyped[usermodel.User](),
			statusCode: http.StatusCreated,
		},

		{
			userId: 0,
			queryValues: map[string]string{
				"vk_user_id": fmt.Sprint(1),
			},
			createUser: usermodel.CreateUser{},
			err:        &validate.WrongLenError{},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			useCase.EXPECT().Create(gomock.Any(), cs.userId, &cs.createUser).Return(&cs.user, cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()

		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(cs.createUser)
		assert.NilError(t, err, "failed encode create user")

		req := httptest.NewRequest(http.MethodGet, "/users/reset?"+urlValues.Encode(), &body)
		req.Header.Set("Content-Type", "application/json")

		transport.CreateUserHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, err, rec)
		} else {
			var user usermodel.User
			err := json.NewDecoder(rec.Result().Body).Decode(&user)
			assert.NilError(t, err, "failed decode user")

			fields, ok := compareUsers(user, cs.user)
			assert.Equal(t, true, ok, "fields %s no equal", fields)
		}
	}
}

func Test_FriendsHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	useCase := userusecase.NewMockUserUseCase(ctrl)

	transport := user.NewUserTransport(useCase)

	cases := []struct {
		friendsIds []int64

		friends []*usermodel.Friend
		err     error

		statusCode int
	}{
		{
			friendsIds: randomIntList(101),
			friends:    randomFriendsList(10),

			statusCode: http.StatusOK,
		},
	}

	for _, cs := range cases {
		useCase.EXPECT().Friends(gomock.Any(), cs.friendsIds).Return(cs.friends).Times(1)

		rec := httptest.NewRecorder()

		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(cs.friendsIds)
		assert.NilError(t, err, "failed encode create user")

		req := httptest.NewRequest(http.MethodGet, "/users/friends", &body)
		req.Header.Set("Content-Type", "application/json")

		transport.FriendsHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, err, rec)
		} else {
			friends := make([]*usermodel.Friend, 0)
			err := json.NewDecoder(rec.Result().Body).Decode(&friends)
			rec.Result().Body.Close()
			assert.NilError(t, err, "failed decode body")

			equal := slices.EqualFunc(friends, cs.friends, func(f1, f2 *usermodel.Friend) bool {
				fields, ok := compareFriends(f1, f2)
				if ok {
					return true
				}
				log.Printf("fields not equal, %s", fields)
				return false
			})
			assert.Equal(t, true, equal, "friends not equal")
		}
	}
}
