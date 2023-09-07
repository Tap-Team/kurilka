package userdatamanager

import (
	"context"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "achievements/datamanager/userdatamanager"

type UserStorage interface {
	User(context.Context, int64) (*model.UserData, error)
}

type UserCache interface {
	User(context.Context, int64) (*model.UserData, error)
}

type UserManager interface {
	UserData(ctx context.Context, userId int64) (*model.UserData, error)
}

type userDataManager struct {
	storage UserStorage
	cache   UserCache
}

func New(
	storage UserStorage,
	cache UserCache,
) UserManager {
	return &userDataManager{
		storage: storage,
		cache:   cache,
	}
}

func (m *userDataManager) UserData(ctx context.Context, userId int64) (*model.UserData, error) {
	data, err := m.cache.User(ctx, userId)
	if err == nil {
		return data, nil
	}
	data, err = m.storage.User(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user from storage", "UserData", _PROVIDER))
	}
	return data, nil
}
