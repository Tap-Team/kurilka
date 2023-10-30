package userusecase

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/achievementmessagesender"
	"github.com/Tap-Team/kurilka/internal/domain/userstatisticscounter"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/workers"
)

//go:generate mockgen -source user.go -destination user_mocks.go -package userusecase

const _PROVIDER = "user/userusecase"

type UserFriendsProvider interface {
	Friends(ctx context.Context, userId int64) []int64
}

type FriendProvider interface {
	Friend(ctx context.Context, userId int64) (*usermodel.Friend, error)
}

type userUseCase struct {
	userFriends              UserFriendsProvider
	user                     userdatamanager.UserManager
	privacySetting           privacysettingdatamanager.PrivacySettingManager
	achievement              achievementdatamanager.AchievementManager
	friend                   FriendProvider
	subscription             SubscriptionStorage
	userWorker               workers.UserWorker
	achievementMessageSender achievementmessagesender.AchievementMessageSenderAtTime
}

type UserUseCase interface {
	Create(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.User, error)
	Reset(ctx context.Context, userId int64) error
	User(ctx context.Context, userId int64) (*usermodel.User, error)
	Level(ctx context.Context, userId int64) (*usermodel.LevelInfo, error)
	Friends(ctx context.Context, userId int64) []*usermodel.Friend
}

func NewUser(
	userFriends UserFriendsProvider,
	user userdatamanager.UserManager,
	privacySetting privacysettingdatamanager.PrivacySettingManager,
	achievement achievementdatamanager.AchievementManager,
	friendProvider FriendProvider,
	subscription SubscriptionStorage,
	userWorker workers.UserWorker,
	achievementMessageSender achievementmessagesender.AchievementMessageSenderAtTime,
) UserUseCase {
	return &userUseCase{
		userFriends:              userFriends,
		user:                     user,
		privacySetting:           privacySetting,
		achievement:              achievement,
		friend:                   friendProvider,
		subscription:             subscription,
		userWorker:               userWorker,
		achievementMessageSender: achievementMessageSender,
	}
}

func NewUserMapper(data *usermodel.UserData, now time.Time) UserMapper {
	counter := userstatisticscounter.NewCounter(
		now,
		data.AbstinenceTime.Time,
		int(data.CigaretteDayAmount),
		int(data.CigarettePackAmount),
		float64(data.PackPrice),
		userstatisticscounter.Second,
	)
	return UserMapper{data: data, now: now, Counter: counter}
}

type UserMapper struct {
	data *usermodel.UserData
	now  time.Time
	userstatisticscounter.Counter
}

func (u UserMapper) User(userId int64, subscription usermodel.Subscription) *usermodel.User {
	user := usermodel.NewUser(
		userId,
		u.data.Name,
		u.data.AbstinenceTime.Time,
		u.Life(),
		u.Cigarette(),
		u.Time(),
		u.Money(),
		u.data.Motivation,
		u.data.WelcomeMotivation,
		u.data.Level,
		u.data.Triggers,
	)
	if subscription.Type == usermodel.TRIAL && !subscription.IsExpired() {
		user.Triggers = slices.DeleteFunc(user.Triggers, func(t usermodel.Trigger) bool { return t == usermodel.SUPPORT_TRIAL })
	}
	return user
}

func (u UserMapper) Friend(
	friendId int64,
	achievements []*usermodel.Achievement,
	privacySettings []usermodel.PrivacySetting,
	subscriptionType usermodel.SubscriptionType,
) *usermodel.Friend {
	return usermodel.NewFriend(
		friendId,
		u.data.AbstinenceTime.Time,
		u.Life(),
		u.Cigarette(),
		u.Time(),
		u.Money(),
		subscriptionType,
		u.data.Level,
		achievements,
		privacySettings,
	)
}

func (u *userUseCase) Create(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.User, error) {
	userData, err := u.user.Create(ctx, userId, createUser)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("create user", "Create", _PROVIDER))
	}
	subscription, _ := u.subscription.UserSubscription(ctx, userId)
	u.userWorker.AddUser(ctx, workers.NewUser(userId, userData.AbstinenceTime.Time))
	return NewUserMapper(userData, time.Now()).User(userId, subscription), nil
}

func (u *userUseCase) Reset(ctx context.Context, userId int64) error {
	err := u.user.Reset(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("reset user", "Reset", _PROVIDER))
	}
	u.achievement.Clear(ctx, userId)
	u.privacySetting.Clear(ctx, userId)
	u.subscription.Clear(ctx, userId)
	u.userWorker.RemoveUser(ctx, userId)
	u.achievementMessageSender.CancelSendMessagesForUser(ctx, userId)
	return nil
}

func (u *userUseCase) User(ctx context.Context, userId int64) (*usermodel.User, error) {
	userData, err := u.user.User(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user data", "User", _PROVIDER))
	}
	subscription, _ := u.subscription.UserSubscription(ctx, userId)
	return NewUserMapper(userData, time.Now()).User(userId, subscription), nil
}

func (u *userUseCase) Level(ctx context.Context, userId int64) (*usermodel.LevelInfo, error) {
	level, err := u.user.Level(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user level", "Level", _PROVIDER))
	}
	return level, nil
}

func (u *userUseCase) Friends(ctx context.Context, userId int64) []*usermodel.Friend {
	var wg sync.WaitGroup
	var mu sync.Mutex
	friendsIds := u.userFriends.Friends(ctx, userId)
	friends := make([]*usermodel.Friend, 0)
	wg.Add(len(friendsIds))
	for _, id := range friendsIds {
		go func(id int64) {
			defer wg.Done()
			friend, err := u.friend.Friend(ctx, id)
			if err != nil {
				return
			}
			mu.Lock()
			friends = append(friends, friend)
			mu.Unlock()
		}(id)
	}
	wg.Wait()
	sort.Sort(SortFriendsByIdAsc(friends))
	return friends
}
