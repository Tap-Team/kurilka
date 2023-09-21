package welcomemotivationdatamanager

import (
	"context"

	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
)

//go:generate mockgen -source manager.go -destination mocks.go -package welcomemotivationdatamanager

const _PROVIDER = "workers/userworker/datamanager/welcomemotivationdatamanager.manager"

type WelcomeMotivationStorage interface {
	NextUserWelcomeMotivation(ctx context.Context, userId int64) (model.Motivation, error)
	UpdateUserWelcomeMotivation(ctx context.Context, userId int64, motivationId int) error
}

type WelcomeMotivationCache interface {
	SaveUserWelcomeMotivation(ctx context.Context, userId int64, welcomeMotivation string) error
	RemoveUserWelcomeMotivation(ctx context.Context, userId int64) error
}

type WelcomeMotivationManager interface {
	NextUserWelcomeMotivation(ctx context.Context, userId int64) (model.Motivation, error)
	UpdateUserWelcomeMotivation(ctx context.Context, userId int64, motivation model.Motivation) error
}

type manager struct {
	storage WelcomeMotivationStorage
	cache   WelcomeMotivationCache
}

func New(storage WelcomeMotivationStorage, cache WelcomeMotivationCache) WelcomeMotivationManager {
	return &manager{storage: storage, cache: cache}
}

func (m *manager) NextUserWelcomeMotivation(ctx context.Context, userId int64) (model.Motivation, error) {
	motivation, err := m.storage.NextUserWelcomeMotivation(ctx, userId)
	if err != nil {
		return model.Motivation{}, exception.Wrap(err, exception.NewCause("get next user welcome motivation", "NextUserWelcomeMotivation", _PROVIDER))
	}
	return motivation, nil
}

func (m *manager) UpdateUserWelcomeMotivation(ctx context.Context, userId int64, motivation model.Motivation) error {
	err := m.storage.UpdateUserWelcomeMotivation(ctx, userId, motivation.ID)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("update user welcome motivation in storage", "UpdateUserWelcomeMotivation", _PROVIDER))
	}
	err = m.cache.SaveUserWelcomeMotivation(ctx, userId, motivation.Motivation)
	if err != nil {
		m.cache.RemoveUserWelcomeMotivation(ctx, userId)
	}
	return nil
}
