package userusecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
)

//go:generate mockgen -source friend_provider.go -destination friend_provider_mocks.go -package userusecase

type SubscriptionStorage interface {
	UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error)
	Clear(ctx context.Context, userId int64) error
}

type friendProvider struct {
	achievement    achievementdatamanager.AchievementManager
	user           userdatamanager.UserManager
	privacySetting privacysettingdatamanager.PrivacySettingManager
	subscription   SubscriptionStorage
}

func NewFriendProvider(
	achievement achievementdatamanager.AchievementManager,
	user userdatamanager.UserManager,
	privacySetting privacysettingdatamanager.PrivacySettingManager,
	subscription SubscriptionStorage,
) FriendProvider {
	return &friendProvider{
		achievement:    achievement,
		user:           user,
		privacySetting: privacySetting,
		subscription:   subscription,
	}
}

func (f *friendProvider) UserSubscriptionType(ctx context.Context, userId int64) usermodel.SubscriptionType {
	subscription, err := f.subscription.UserSubscription(ctx, userId)
	if err != nil {
		return usermodel.NONE
	}
	if subscription.IsNoneOrExpired() {
		return usermodel.NONE
	}
	return subscription.Type
}

func (f *friendProvider) Friend(ctx context.Context, userId int64) (*usermodel.Friend, error) {
	userData, err := f.user.User(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user", "friend", _PROVIDER))
	}
	achievements := f.achievement.AchievementPreview(ctx, userId)
	subscriptionType := f.UserSubscriptionType(ctx, userId)
	privacySettings, err := f.privacySetting.PrivacySettings(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, exception.Wrap(err, exception.NewCause("get privacy settings", "friend", _PROVIDER)).Error(), "userId", userId)
	}
	friend := NewUserMapper(userData, time.Now()).Friend(userId, achievements, privacySettings, subscriptionType)
	friend.UseFilters()
	return friend, nil
}
