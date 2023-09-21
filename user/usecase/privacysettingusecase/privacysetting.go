package privacysettingusecase

import (
	"context"
	"slices"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
)

//go:generate mockgen -source privacysetting.go -destination mocks.go -package privacysettingusecase

const _PROVIDER = "user/usecase/privacysettingusecase.privacySettingUseCase"

type PrivacySettingUseCase interface {
	Switch(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
}

type privacySettingUseCase struct {
	manager privacysettingdatamanager.PrivacySettingManager
}

func New(manager privacysettingdatamanager.PrivacySettingManager) PrivacySettingUseCase {
	return &privacySettingUseCase{manager: manager}
}

func (p *privacySettingUseCase) Switch(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	privacySettings, err := p.manager.PrivacySettings(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get user privacy settings", "Switch", _PROVIDER))
	}
	if contains := slices.Contains(privacySettings, setting); contains {
		err := p.manager.Remove(ctx, userId, setting)
		if err != nil {
			return exception.Wrap(err, exception.NewCause("remove privacy setting", "Switch", _PROVIDER))
		}
	} else {
		err := p.manager.Add(ctx, userId, setting)
		if err != nil {
			return exception.Wrap(err, exception.NewCause("add privacy setting", "Switch", _PROVIDER))
		}
	}
	return nil
}
