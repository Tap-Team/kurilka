package privacysettingusecase_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/privacysettingerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/usecase/privacysettingusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_UseCase_Switch(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)

	usecase := privacysettingusecase.New(manager)

	cases := []struct {
		privacySetting usermodel.PrivacySetting

		privacySettings    []usermodel.PrivacySetting
		privacySettingsErr error

		removeCall bool
		removeErr  error

		addCall bool
		addErr  error

		err error
	}{
		{
			privacySettingsErr: usererror.ExceptionUserNotFound(),
			err:                usererror.ExceptionUserNotFound(),
		},
		{
			privacySetting: usermodel.ACHIEVEMENTS_CIGARETTE,

			addCall: true,
			addErr:  privacysettingerror.ExceptionPrivacySettingNotExist(),

			err: privacysettingerror.ExceptionPrivacySettingNotExist(),
		},

		{
			privacySetting:  usermodel.ACHIEVEMENTS_CIGARETTE,
			privacySettings: []usermodel.PrivacySetting{},

			addCall: true,
		},

		{
			privacySetting: usermodel.STATISTICS_LIFE,

			privacySettings: []usermodel.PrivacySetting{usermodel.STATISTICS_LIFE},

			removeCall: true,
			removeErr:  privacysettingerror.ExceptionPrivacySettingNotExist(),

			err: privacysettingerror.ExceptionPrivacySettingNotExist(),
		},

		{
			privacySetting: usermodel.STATISTICS_LIFE,

			privacySettings: []usermodel.PrivacySetting{usermodel.STATISTICS_LIFE},

			removeCall: true,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		manager.EXPECT().PrivacySettings(gomock.Any(), userId).Return(cs.privacySettings, cs.privacySettingsErr).Times(1)
		if cs.addCall {
			manager.EXPECT().Add(gomock.Any(), userId, cs.privacySetting).Return(cs.addErr).Times(1)
		}
		if cs.removeCall {
			manager.EXPECT().Remove(gomock.Any(), userId, cs.privacySetting).Return(cs.removeErr).Times(1)
		}
		err := usecase.Switch(ctx, userId, cs.privacySetting)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}
}
