package motivationdatamanager_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/motivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_Manager_NextUserMotivation(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := motivationdatamanager.NewMockMotivationStorage(ctrl)
	cache := motivationdatamanager.NewMockMotivationCache(ctrl)

	manager := motivationdatamanager.New(storage, cache)

	cases := []struct {
		motivation model.Motivation
		err        error
	}{
		{
			err: errors.New("any error"),
		},
		{
			motivation: random.StructTyped[model.Motivation](),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		storage.EXPECT().NextUserMotivation(gomock.Any(), userId).Return(cs.motivation, cs.err).Times(1)

		motivation, err := manager.NextUserMotivation(ctx, userId)

		assert.Equal(t, motivation, cs.motivation, "motivation not equal")
		assert.ErrorIs(t, err, cs.err, "error not equal")
	}
}

func Test_Manager_UpdateUserMotivation(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	storage := motivationdatamanager.NewMockMotivationStorage(ctrl)
	cache := motivationdatamanager.NewMockMotivationCache(ctrl)

	manager := motivationdatamanager.New(storage, cache)

	cases := []struct {
		updateStorageCall bool
		updateStorageErr  error

		saveCacheCall bool
		saveCacheErr  error

		removeCacheCall bool
		removeCacheErr  error

		err error
	}{
		{
			updateStorageCall: true,
			updateStorageErr:  usererror.ExceptionUserNotFound(),

			err: usererror.ExceptionUserNotFound(),
		},
		{
			updateStorageCall: true,
			saveCacheCall:     true,
		},
		{
			updateStorageCall: true,
			saveCacheCall:     true,
			saveCacheErr:      errors.New("any error"),
			removeCacheCall:   true,
			removeCacheErr:    errors.New("any error"),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		motivation := random.StructTyped[model.Motivation]()

		if cs.updateStorageCall {
			storage.EXPECT().UpdateUserMotivation(gomock.Any(), userId, motivation.ID).Return(cs.updateStorageErr).Times(1)
		}
		if cs.saveCacheCall {
			cache.EXPECT().SaveUserMotivation(gomock.Any(), userId, motivation.Motivation).Return(cs.saveCacheErr).Times(1)
		}
		if cs.removeCacheCall {
			cache.EXPECT().RemoveUserMotivation(gomock.Any(), userId).Return(cs.removeCacheErr).Times(1)
		}

		err := manager.UpdateUserMotivation(ctx, userId, motivation)
		assert.ErrorIs(t, err, cs.err, "wrong error")
	}
}
