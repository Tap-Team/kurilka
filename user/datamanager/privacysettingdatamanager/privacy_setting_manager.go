package privacysettingdatamanager

import (
	"context"
	"slices"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "user/privacysettingdatamanager"

type PrivacySettingStorage interface {
	UserPrivacySettings(ctx context.Context, userId int64) ([]usermodel.PrivacySetting, error)
	AddUserPrivacySetting(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
	RemoveUserPrivacySetting(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
}

type PrivacySettingCache interface {
	UserPrivacySettings(ctx context.Context, userId int64) ([]usermodel.PrivacySetting, error)
	SaveUserPrivacySettings(ctx context.Context, userId int64, settings []usermodel.PrivacySetting) error
	RemoveUserPrivacySettings(ctx context.Context, userId int64) error
}

type PrivacySettingManager interface {
	PrivacySettings(ctx context.Context, userId int64) ([]usermodel.PrivacySetting, error)
	Add(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
	Remove(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
	Clear(ctx context.Context, userId int64)
}

type userPrivacySettingManager struct {
	cache   *CacheWrapper
	storage PrivacySettingStorage
}

func NewPrivacyManager(storage PrivacySettingStorage, cache PrivacySettingCache) PrivacySettingManager {
	return &userPrivacySettingManager{storage: storage, cache: &CacheWrapper{PrivacySettingCache: cache}}
}

type CacheWrapper struct {
	PrivacySettingCache
}

func (p *CacheWrapper) AddSingle(ctx context.Context, userId int64, setting usermodel.PrivacySetting) {
	settings, err := p.UserPrivacySettings(ctx, userId)
	if err != nil {
		return
	}
	settings = append(settings, setting)
	err = p.SaveUserPrivacySettings(ctx, userId, settings)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("add user privacy settings", "AddSingle", _PROVIDER)).Error())
		p.RemoveUserPrivacySettings(ctx, userId)
	}
}

func (p *CacheWrapper) RemoveSingle(ctx context.Context, userId int64, setting usermodel.PrivacySetting) {
	settings, err := p.UserPrivacySettings(ctx, userId)
	if err != nil {
		return
	}
	index := -1
	for i := range settings {
		if settings[i] == setting {
			index = i
		}
	}
	if index == -1 {
		slog.Info("User Try remove user achievement", "userId", userId, "privacySetting", setting)
		return
	}
	settings = slices.Delete(settings, index, index+1)
	err = p.SaveUserPrivacySettings(ctx, userId, settings)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("remove user privacy settings", "RemoveSingle", _PROVIDER)).Error())
		p.RemoveUserPrivacySettings(ctx, userId)
	}
}

func (u *userPrivacySettingManager) Add(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	err := u.storage.AddUserPrivacySetting(ctx, userId, setting)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("add user privacy setting", "Add", _PROVIDER))
	}
	u.cache.AddSingle(ctx, userId, setting)
	return nil
}

func (u *userPrivacySettingManager) Remove(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	err := u.storage.RemoveUserPrivacySetting(ctx, userId, setting)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("add user privacy setting", "Add", _PROVIDER))
	}
	u.cache.RemoveSingle(ctx, userId, setting)
	return nil
}

func (u *userPrivacySettingManager) Clear(ctx context.Context, userId int64) {
	u.cache.RemoveUserPrivacySettings(ctx, userId)
}

func (u *userPrivacySettingManager) PrivacySettings(ctx context.Context, userId int64) ([]usermodel.PrivacySetting, error) {
	settings, err := u.cache.UserPrivacySettings(ctx, userId)
	if err == nil {
		return settings, nil
	}
	settings, err = u.storage.UserPrivacySettings(ctx, userId)
	if err != nil {
		return settings, exception.Wrap(err, exception.NewCause("get user privacy settings", "PrivacySettings", _PROVIDER))
	}
	err = u.cache.SaveUserPrivacySettings(ctx, userId, settings)
	if err != nil {
		slog.ErrorContext(ctx, "failed remove user privacy settings", "err", err)
		u.cache.RemoveUserPrivacySettings(ctx, userId)
	}
	return settings, nil
}
