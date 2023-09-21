package userdatamanager_test

import (
	context "context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	NilUserData *model.UserData
)

func Test_Manager_UserData(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cache := userdatamanager.NewMockUserCache(ctrl)
	storage := userdatamanager.NewMockUserCache(ctrl)

	manager := userdatamanager.New(storage, cache)

	{
		userId := rand.Int63()
		userData := random.StructTyped[model.UserData]()
		cache.EXPECT().User(gomock.Any(), userId).Return(&userData, nil).Times(1)

		data, err := manager.UserData(ctx, userId)

		assert.NilError(t, err, "non nil error")
		assert.Equal(t, data, &userData, "wrong user data")
	}

	{
		userId := rand.Int63()
		userData := random.StructTyped[model.UserData]()

		cache.EXPECT().User(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		storage.EXPECT().User(gomock.Any(), userId).Return(&userData, nil).Times(1)

		data, err := manager.UserData(ctx, userId)

		assert.NilError(t, err, "non nil error")
		assert.Equal(t, data, &userData, "wrong user data")
	}

	{
		userId := rand.Int63()
		expectedErr := usererror.ExceptionUserNotFound()

		cache.EXPECT().User(gomock.Any(), userId).Return(nil, errors.New("any error")).Times(1)
		storage.EXPECT().User(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		data, err := manager.UserData(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
		assert.Equal(t, data, NilUserData, "user data not equal")
	}

}
