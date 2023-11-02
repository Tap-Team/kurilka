package privacysettingusecase_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/privacysettingerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/usecase/privacysettingusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type UserSubscriptionCall struct {
	Subscription usermodel.Subscription
	Err          error
}

func (c UserSubscriptionCall) Call() bool {
	return c != UserSubscriptionCall{}
}

type AddPrivacySettingCall struct {
	PrivacySetting usermodel.PrivacySetting
	Err            error
}

func (c AddPrivacySettingCall) Call() bool {
	return c != AddPrivacySettingCall{}
}

type RemovePrivacySettingCall struct {
	PrivacySetting usermodel.PrivacySetting
	Err            error
}

func (c RemovePrivacySettingCall) Call() bool {
	return c != RemovePrivacySettingCall{}
}

type PrivacySettingCall struct {
	PrivacySettings    []usermodel.PrivacySetting
	PrivacySettingsErr error
	Call               bool
	AddCall            AddPrivacySettingCall
	RemoveCall         RemovePrivacySettingCall
}

func Test_UseCase_Switch(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	subscriptionManager := privacysettingusecase.NewMockUserSubscriptionProvider(ctrl)

	usecase := privacysettingusecase.New(privacySettingsManager, subscriptionManager)

	cases := []struct {
		privacySetting       usermodel.PrivacySetting
		userSubscriptionCall UserSubscriptionCall
		privacySettingCall   PrivacySettingCall
		err                  error
	}{
		{
			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour)),
				Err:          usererror.ExceptionUserNotFound(),
			},
			err: privacysettingerror.ExceptionUserWithoutSubscription(),
		},
		{
			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.NONE, time.Now().Add(time.Hour)),
			},
			err: privacysettingerror.ExceptionUserWithoutSubscription(),
		},
		{
			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.BASIC, time.Time{}),
			},
			err: privacysettingerror.ExceptionUserWithoutSubscription(),
		},
		{
			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
			},
			privacySettingCall: PrivacySettingCall{
				PrivacySettingsErr: usererror.ExceptionUserNotFound(),
				Call:               true,
			},
			err: usererror.ExceptionUserNotFound(),
		},
		{
			privacySetting: usermodel.ACHIEVEMENTS_CIGARETTE,
			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
			},
			privacySettingCall: PrivacySettingCall{
				Call: true,
				AddCall: AddPrivacySettingCall{
					PrivacySetting: usermodel.ACHIEVEMENTS_CIGARETTE,
					Err:            privacysettingerror.ExceptionPrivacySettingNotExist(),
				},
			},

			err: privacysettingerror.ExceptionPrivacySettingNotExist(),
		},

		{
			privacySetting: usermodel.ACHIEVEMENTS_CIGARETTE,

			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
			},

			privacySettingCall: PrivacySettingCall{
				Call: true,
				AddCall: AddPrivacySettingCall{
					PrivacySetting: usermodel.ACHIEVEMENTS_CIGARETTE,
				},
			},
		},

		{
			privacySetting: usermodel.STATISTICS_LIFE,

			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
			},

			privacySettingCall: PrivacySettingCall{
				Call:            true,
				PrivacySettings: []usermodel.PrivacySetting{usermodel.STATISTICS_LIFE},
				RemoveCall: RemovePrivacySettingCall{
					PrivacySetting: usermodel.STATISTICS_LIFE,
					Err:            privacysettingerror.ExceptionPrivacySettingNotExist(),
				},
			},

			err: privacysettingerror.ExceptionPrivacySettingNotExist(),
		},

		{
			privacySetting: usermodel.STATISTICS_LIFE,

			userSubscriptionCall: UserSubscriptionCall{
				Subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
			},

			privacySettingCall: PrivacySettingCall{
				Call:            true,
				PrivacySettings: []usermodel.PrivacySetting{usermodel.STATISTICS_LIFE},
				RemoveCall: RemovePrivacySettingCall{
					PrivacySetting: usermodel.STATISTICS_LIFE,
				},
			},
		},
	}

	for caseNumber, cs := range cases {
		t.Log(caseNumber)
		userId := rand.Int63()
		if cs.userSubscriptionCall.Call() {
			subscriptionManager.EXPECT().
				UserSubscription(gomock.Any(), userId).
				Return(cs.userSubscriptionCall.Subscription, cs.userSubscriptionCall.Err).
				Times(1)
		}
		if cs.privacySettingCall.Call {
			privacySettingsManager.EXPECT().
				PrivacySettings(gomock.Any(), userId).
				Return(cs.privacySettingCall.PrivacySettings, cs.privacySettingCall.PrivacySettingsErr).
				Times(1)
		}
		if cs.privacySettingCall.AddCall.Call() {
			privacySettingsManager.EXPECT().
				Add(gomock.Any(), userId, cs.privacySettingCall.AddCall.PrivacySetting).
				Return(cs.privacySettingCall.AddCall.Err).
				Times(1)
		}

		if cs.privacySettingCall.RemoveCall.Call() {
			privacySettingsManager.EXPECT().
				Remove(gomock.Any(), userId, cs.privacySettingCall.RemoveCall.PrivacySetting).
				Return(cs.privacySettingCall.RemoveCall.Err).
				Times(1)
		}
		err := usecase.Switch(ctx, userId, cs.privacySetting)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}
}
