package userdatamanager_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	NilUser  *usermodel.UserData
	NilLevel *usermodel.LevelInfo
)

func Test_User_Saver_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	storage := userdatamanager.NewMockUserStorage(ctrl)
	recoverer := userdatamanager.NewMockUserRecoverer(ctrl)

	saver := userdatamanager.NewUserSaver(storage, recoverer)

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()
		deleted := rand.Intn(3)%2 == 0

		expectedErr := errors.New("database umir")

		storage.EXPECT().UserDeleted(gomock.Any(), userId).Return(deleted, expectedErr).Times(1)

		user, err := saver.Save(ctx, userId, &createUser)

		assert.Equal(t, user, NilUser, "user not nil")
		assert.ErrorIs(t, err, expectedErr, "wrong err from save")
	}

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()
		deleted := rand.Intn(3)%2 == 0

		expectedUser := random.StructTyped[usermodel.UserData]()
		expectedError := errors.New("random err")

		storage.EXPECT().UserDeleted(gomock.Any(), userId).Return(deleted, usererror.ExceptionUserNotFound()).Times(1)
		storage.EXPECT().InsertUser(gomock.Any(), userId, &createUser).Return(&expectedUser, expectedError).Times(1)

		user, err := saver.Save(ctx, userId, &createUser)

		assert.ErrorIs(t, err, expectedError, "wrong error")
		assert.Equal(t, user, &expectedUser, "wrong user")
	}

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()

		deleted := false

		expectedErr := usererror.ExceptionUserExist()

		storage.EXPECT().UserDeleted(gomock.Any(), userId).Return(deleted, nil).Times(1)

		user, err := saver.Save(ctx, userId, &createUser)

		assert.Equal(t, user, NilUser, "non nil user")
		assert.ErrorIs(t, err, expectedErr, "wrong err from save")
	}

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()

		deleted := true

		expectedUser := random.StructTyped[usermodel.UserData]()
		expectedErr := errors.New("gofman polozhil prod")

		storage.EXPECT().UserDeleted(gomock.Any(), userId).Return(deleted, nil).Times(1)
		recoverer.EXPECT().RecoverUser(gomock.Any(), userId, &createUser).Return(&expectedUser, expectedErr)

		user, err := saver.Save(ctx, userId, &createUser)

		assert.Equal(t, user, &expectedUser, "wrong user from save")
		assert.ErrorIs(t, err, expectedErr, "wrong err from saves")
	}
}

func Test_User_Manager_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	storage := userdatamanager.NewMockUserStorage(ctrl)
	cache := userdatamanager.NewMockUserCache(ctrl)
	recoverResetter := userdatamanager.NewMockUserRecoverReseter(ctrl)
	saver := userdatamanager.NewMockUserSaver(ctrl)

	manager := userdatamanager.NewUserManager(recoverResetter, storage, cache, saver)

	{
		userId := rand.Int63()

		expectedUser := random.StructTyped[usermodel.UserData]()

		cache.EXPECT().User(gomock.Any(), userId).Return(&expectedUser, nil).Times(1)

		user, err := manager.User(ctx, userId)

		assert.NilError(t, err, "non nil err")
		assert.Equal(t, user, &expectedUser, "wrong user")
	}

	{
		userId := rand.Int63()

		expectedErr := errors.New("random error")

		cache.EXPECT().User(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		storage.EXPECT().User(ctx, userId).Return(nil, expectedErr).Times(1)

		user, err := manager.User(ctx, userId)

		assert.Equal(t, user, NilUser, "non nil user")
		assert.ErrorIs(t, err, expectedErr, "wrong err")
	}

	{
		userId := rand.Int63()

		userData := random.StructTyped[usermodel.UserData]()

		cache.EXPECT().User(gomock.Any(), userId).Return(nil, errors.New("gof break all")).Times(1)
		storage.EXPECT().User(ctx, userId).Return(&userData, nil).Times(1)
		cache.EXPECT().SaveUser(gomock.Any(), userId, &userData).Return(nil).Times(1)

		user, err := manager.User(ctx, userId)

		assert.Equal(t, user, &userData, "wrong user")
		assert.NilError(t, err, "non nil error")
	}

}

type userLevelMatcher struct {
	level usermodel.LevelInfo
}

func (l *userLevelMatcher) Matches(x interface{}) bool {
	user, ok := x.(*usermodel.UserData)
	if !ok {
		return false
	}
	return reflect.DeepEqual(user.Level, l.level)
}

func (l *userLevelMatcher) String() string {
	return fmt.Sprintf("is equal %v", l.level)
}
func NewUserLevelMatcher(level usermodel.LevelInfo) gomock.Matcher {
	return &userLevelMatcher{level: level}
}

func Test_User_Manager_Level(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	storage := userdatamanager.NewMockUserStorage(ctrl)
	cache := userdatamanager.NewMockUserCache(ctrl)
	recoverResetter := userdatamanager.NewMockUserRecoverReseter(ctrl)
	saver := userdatamanager.NewMockUserSaver(ctrl)

	manager := userdatamanager.NewUserManager(recoverResetter, storage, cache, saver)

	{
		userId := rand.Int63()

		expectedErr := errors.New("random error")

		storage.EXPECT().UserLevel(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		level, err := manager.Level(ctx, userId)

		assert.Equal(t, level, NilLevel, "non nil error")
		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	{
		userId := rand.Int63()
		expectedLevel := random.StructTyped[usermodel.LevelInfo]()

		storage.EXPECT().UserLevel(gomock.Any(), userId).Return(&expectedLevel, nil).Times(1)
		cache.EXPECT().User(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)

		level, err := manager.Level(ctx, userId)

		assert.NilError(t, err, "wrong error from level")
		assert.Equal(t, level, &expectedLevel)
	}

	{
		userId := rand.Int63()
		expectedLevel := random.StructTyped[usermodel.LevelInfo]()

		user := random.StructTyped[usermodel.UserData]()

		storage.EXPECT().UserLevel(gomock.Any(), userId).Return(&expectedLevel, nil).Times(1)
		cache.EXPECT().User(gomock.Any(), userId).Return(&user, nil).Times(1)
		cache.EXPECT().SaveUser(gomock.Any(), userId, NewUserLevelMatcher(expectedLevel)).Return(nil).Times(1)

		level, err := manager.Level(ctx, userId)

		assert.NilError(t, err, "non nil error from level")
		assert.Equal(t, level, &expectedLevel)
	}

	{
		userId := rand.Int63()
		expectedLevel := random.StructTyped[usermodel.LevelInfo]()

		user := random.StructTyped[usermodel.UserData]()

		storage.EXPECT().UserLevel(gomock.Any(), userId).Return(&expectedLevel, nil).Times(1)
		cache.EXPECT().User(gomock.Any(), userId).Return(&user, nil).Times(1)
		cache.EXPECT().SaveUser(gomock.Any(), userId, NewUserLevelMatcher(expectedLevel)).Return(errors.New("i dont know")).Times(1)
		cache.EXPECT().RemoveUser(gomock.Any(), userId).Return(nil).Times(1)

		level, err := manager.Level(ctx, userId)

		assert.NilError(t, err, "non nil error from level")
		assert.Equal(t, level, &expectedLevel)
	}
}

func Test_User_Manager_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	storage := userdatamanager.NewMockUserStorage(ctrl)
	cache := userdatamanager.NewMockUserCache(ctrl)
	recoverResetter := userdatamanager.NewMockUserRecoverReseter(ctrl)
	saver := userdatamanager.NewMockUserSaver(ctrl)

	manager := userdatamanager.NewUserManager(recoverResetter, storage, cache, saver)

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()

		expectedErr := errors.New("any error")

		saver.EXPECT().Save(gomock.Any(), userId, &createUser).Return(nil, expectedErr).Times(1)

		user, err := manager.Create(ctx, userId, &createUser)

		assert.ErrorIs(t, err, expectedErr, "wrong err from user")
		assert.Equal(t, user, NilUser, "wrong user")
	}

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()

		expectedUserData := random.StructTyped[usermodel.UserData]()

		saver.EXPECT().Save(gomock.Any(), userId, &createUser).Return(&expectedUserData, nil).Times(1)
		cache.EXPECT().SaveUser(gomock.Any(), userId, &expectedUserData).Return(nil).Times(1)

		user, err := manager.Create(ctx, userId, &createUser)

		assert.NilError(t, err, "non nil error")
		assert.Equal(t, user, &expectedUserData, "user not equal")
	}
}

func Test_User_Manager_Reset(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	storage := userdatamanager.NewMockUserStorage(ctrl)
	cache := userdatamanager.NewMockUserCache(ctrl)
	recoverResetter := userdatamanager.NewMockUserRecoverReseter(ctrl)
	saver := userdatamanager.NewMockUserSaver(ctrl)

	manager := userdatamanager.NewUserManager(recoverResetter, storage, cache, saver)

	{
		userId := rand.Int63()

		expectedErr := errors.New("failed reset user")

		recoverResetter.EXPECT().ResetUser(gomock.Any(), userId).Return(expectedErr).Times(1)

		err := manager.Reset(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error from reset")
	}

	{
		userId := rand.Int63()

		recoverResetter.EXPECT().ResetUser(gomock.Any(), userId).Return(nil).Times(1)
		cache.EXPECT().RemoveUser(gomock.Any(), userId).Return(nil).Times(1)

		err := manager.Reset(ctx, userId)

		assert.NilError(t, err, "non nil error")
	}
}

func Test_User_Manager_FilterExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	storage := userdatamanager.NewMockUserStorage(ctrl)
	cache := userdatamanager.NewMockUserCache(ctrl)
	recoverResetter := userdatamanager.NewMockUserRecoverReseter(ctrl)
	saver := userdatamanager.NewMockUserSaver(ctrl)

	manager := userdatamanager.NewUserManager(recoverResetter, storage, cache, saver)

	{
		userIds := []int64{1, 2, 232, 134, 123, 5, 7, 21341, 6, 73, 456}
		filterIds := []int64{1, 2, 5, 7, 6, 73}

		storage.EXPECT().Exists(gomock.Any(), userIds).Return(filterIds).Times(1)

		filterUserIds := manager.FilterExists(ctx, userIds)

		equal := slices.Equal(filterIds, filterUserIds)
		assert.Equal(t, true, equal, "ids not equal")
	}
}
