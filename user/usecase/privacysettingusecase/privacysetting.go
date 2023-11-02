package privacysettingusecase

import (
	"context"
	"log/slog"
	"slices"

	"github.com/Tap-Team/kurilka/internal/errorutils/privacysettingerror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
)

//go:generate mockgen -source privacysetting.go -destination mocks.go -package privacysettingusecase

const _PROVIDER = "user/usecase/privacysettingusecase.privacySettingUseCase"

type PrivacySettingUseCase interface {
	Switch(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
}

type UserSubscriptionProvider interface {
	UserSubscription(ctx context.Context, userid int64) (usermodel.Subscription, error)
}

type privacySettingUseCase struct {
	privacySettingManager   privacysettingdatamanager.PrivacySettingManager
	userSubscriptionManager UserSubscriptionProvider
}

func New(manager privacysettingdatamanager.PrivacySettingManager, userSubscriptionProvider UserSubscriptionProvider) PrivacySettingUseCase {
	return &privacySettingUseCase{
		privacySettingManager:   manager,
		userSubscriptionManager: userSubscriptionProvider,
	}
}

func nilOrError(err error, cause exception.Cause) error {
	if err == nil {
		return nil
	}
	return exception.Wrap(err, cause)
}

func (p *privacySettingUseCase) Switch(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	if !p.isUserHaveSubscription(ctx, userId) {
		return privacysettingerror.ExceptionUserWithoutSubscription()
	}
	return p.switchUserPrivacySetting(ctx, userId, setting)
}

func (p *privacySettingUseCase) isUserHaveSubscription(ctx context.Context, userId int64) bool {
	subscription, err := p.userSubscriptionManager.UserSubscription(ctx, userId)
	if err != nil {
		err := exception.Wrap(err, exception.NewCause("failed get user subscription", "isUserHaveSubscription", _PROVIDER))
		slog.ErrorContext(ctx, err.Error(), "user_id", userId)
		return false
	}
	return !subscription.IsNoneOrExpired()
}

func (p *privacySettingUseCase) switchUserPrivacySetting(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	privacySettings, err := p.privacySettingManager.PrivacySettings(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get user privacy settings", "Switch", _PROVIDER))
	}
	if slices.Contains(privacySettings, setting) {
		err := p.privacySettingManager.Remove(ctx, userId, setting)
		return nilOrError(err, exception.NewCause("remove privacy setting", "Switch", _PROVIDER))
	} else {
		err := p.privacySettingManager.Add(ctx, userId, setting)
		return nilOrError(err, exception.NewCause("add privacy setting", "Switch", _PROVIDER))
	}
}
