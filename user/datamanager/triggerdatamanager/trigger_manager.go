package triggerdatamanager

import (
	"context"
	"slices"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source trigger_manager.go -destination trigger_manager_mocks.go -package triggerdatamanager

const _PROVIDER = "user/datamanager/triggerdatamanager.triggerManager"

type TriggerCache interface {
	UserTriggers(ctx context.Context, userId int64) ([]usermodel.Trigger, error)
	SaveUserTriggers(ctx context.Context, userId int64, triggers []usermodel.Trigger) error
	RemoveUserTriggers(ctx context.Context, userId int64) error
}

type TriggerStorage interface {
	Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) error
	Add(ctx context.Context, userId int64, trigger usermodel.Trigger) error
}

type TriggerManager interface {
	Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) error
	Add(ctx context.Context, userId int64, trigger usermodel.Trigger) error
}

type triggerManager struct {
	cache   CacheWrapper
	storage TriggerStorage
}

func NewTriggerManager(storage TriggerStorage, cache TriggerCache) TriggerManager {
	return &triggerManager{storage: storage, cache: CacheWrapper{cache}}
}

type CacheWrapper struct{ TriggerCache }

func (cw *CacheWrapper) Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) {
	triggers, err := cw.UserTriggers(ctx, userId)
	if err != nil {
		return
	}
	index := -1
	for i, tr := range triggers {
		if tr == trigger {
			index = i
			break
		}
	}
	if index == -1 {
		slog.Info("user remove non exists trigger", "userId", userId, "trigger", trigger)
		return
	}
	triggers = slices.Delete(triggers, index, index+1)
	err = cw.SaveUserTriggers(ctx, userId, triggers)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user triggers", "CacheWrapper.Remove", _PROVIDER)).Error(), "userId", userId)
		cw.RemoveUserTriggers(ctx, userId)
	}
}

func (cw *CacheWrapper) Add(ctx context.Context, userId int64, trigger usermodel.Trigger) {
	triggers, err := cw.UserTriggers(ctx, userId)
	if err != nil {
		return
	}
	for _, t := range triggers {
		if t == trigger {
			return
		}
	}
	triggers = append(triggers, trigger)
	err = cw.SaveUserTriggers(ctx, userId, triggers)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user triggers", "CacheWrapper.Remove", _PROVIDER)).Error(), "userId", userId)
		cw.RemoveUserTriggers(ctx, userId)
	}
}

func (t *triggerManager) Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) error {
	err := t.storage.Remove(ctx, userId, trigger)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("remove trigger from storage", "Remove", _PROVIDER))
	}
	t.cache.Remove(ctx, userId, trigger)
	return nil
}

func (t *triggerManager) Add(ctx context.Context, userId int64, trigger usermodel.Trigger) error {
	err := t.storage.Add(ctx, userId, trigger)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("add user trigger in storage", "Add", _PROVIDER))
	}
	t.cache.Add(ctx, userId, trigger)
	return nil
}
