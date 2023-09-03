package triggerdatamanager

import (
	"context"
	"slices"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "user/triggerdatamanager"

type TriggerCache interface {
	UserTriggers(ctx context.Context, userId int64) ([]usermodel.Trigger, error)
	SaveUserTriggers(ctx context.Context, userId int64, triggers []usermodel.Trigger) error
	RemoveUserTriggers(ctx context.Context, userId int64) error
}

type TriggerStorage interface {
	Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) error
}

type TriggerManager interface {
	Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) error
}

type triggerManager struct {
	cache   CacheWrapper
	storage TriggerStorage
}

func NewTriggerManager(storage TriggerStorage, cache TriggerCache) TriggerManager {
	return &triggerManager{storage: storage, cache: CacheWrapper{cache}}
}

type CacheWrapper struct{ TriggerCache }

func (t *CacheWrapper) Remove(ctx context.Context, userId int64, trigger usermodel.Trigger) {
	triggers, err := t.UserTriggers(ctx, userId)
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
	err = t.SaveUserTriggers(ctx, userId, triggers)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("save user triggers", "CacheWrapper.Remove", _PROVIDER)).Error(), "userId", userId)
		t.RemoveUserTriggers(ctx, userId)
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
