package userusecase

import (
	"context"
	"log/slog"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
)

type friendProvider struct {
	achievement    achievementdatamanager.AchievementManager
	user           userdatamanager.UserManager
	privacySetting privacysettingdatamanager.PrivacySettingManager
}

func NewFriendProvider(
	achievement achievementdatamanager.AchievementManager,
	user userdatamanager.UserManager,
	privacySetting privacysettingdatamanager.PrivacySettingManager,
) FriendProvider {
	return &friendProvider{
		achievement:    achievement,
		user:           user,
		privacySetting: privacySetting,
	}
}

func (f *friendProvider) Friend(ctx context.Context, userId int64) (*usermodel.Friend, error) {
	userData, err := f.user.User(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user", "friend", _PROVIDER))
	}
	achievements := f.achievement.AchievementPreview(ctx, userId)
	privacySettings, err := f.privacySetting.PrivacySettings(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, exception.Wrap(err, exception.NewCause("get privacy settings", "friend", _PROVIDER)).Error(), "userId", userId)
	}
	friend := NewUserMapper(userData).Friend(userId, achievements, privacySettings)
	friend.UseFilters()
	return friend, nil
}
