package motivationdatamanager

import (
	"context"

	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
)

//go:generate mockgen -source manager.go -destination mocks.go -package motivationdatamanager

const _PROVIDER = "workers/userworker/datamanager/achievementdatamanager.manager"

type MotivationStorage interface {
	NextUserMotivation(ctx context.Context, userId int64) (model.Motivation, error)
	UpdateUserMotivation(ctx context.Context, userId int64, motivationId int) error
}

type MotivationCache interface {
	SaveUserMotivation(ctx context.Context, userId int64, motivation string) error
	RemoveUserMotivation(ctx context.Context, userId int64) error
}

type MotivationManager interface {
	NextUserMotivation(ctx context.Context, userId int64) (model.Motivation, error)
	UpdateUserMotivation(ctx context.Context, userId int64, motivation model.Motivation) error
}

type manager struct {
	storage MotivationStorage
	cache   MotivationCache
}

func New(storage MotivationStorage, cache MotivationCache) MotivationManager {
	return &manager{storage: storage, cache: cache}
}

func (m *manager) NextUserMotivation(ctx context.Context, userId int64) (model.Motivation, error) {
	motivation, err := m.storage.NextUserMotivation(ctx, userId)
	if err != nil {
		return model.Motivation{}, exception.Wrap(err, exception.NewCause("get next user motivation", "NextUserMotivation", _PROVIDER))
	}
	return motivation, nil
}

func (m *manager) UpdateUserMotivation(ctx context.Context, userId int64, motivation model.Motivation) error {
	err := m.storage.UpdateUserMotivation(ctx, userId, motivation.ID)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("update user motivation", "UpdateUserMotivation", _PROVIDER))
	}
	err = m.cache.SaveUserMotivation(ctx, userId, motivation.Motivation)
	if err != nil {
		m.cache.RemoveUserMotivation(ctx, userId)
	}
	return nil
}
