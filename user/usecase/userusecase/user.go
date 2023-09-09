package userusecase

import (
	"context"
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
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
	userFriends    UserFriendsProvider
	user           userdatamanager.UserManager
	privacySetting privacysettingdatamanager.PrivacySettingManager
	achievement    achievementdatamanager.AchievementManager
	friend         FriendProvider
}

type UserUseCase interface {
	Create(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.User, error)
	Reset(ctx context.Context, userId int64) error
	User(ctx context.Context, userId int64) (*usermodel.User, error)
	Level(ctx context.Context, userId int64) (*usermodel.LevelInfo, error)
	Friends(ctx context.Context, friendsIds []int64) []*usermodel.Friend
}

func NewUser(
	userFriends UserFriendsProvider,
	user userdatamanager.UserManager,
	privacySetting privacysettingdatamanager.PrivacySettingManager,
	achievement achievementdatamanager.AchievementManager,
	friendProvider FriendProvider,
) UserUseCase {
	return &userUseCase{
		userFriends:    userFriends,
		user:           user,
		privacySetting: privacySetting,
		achievement:    achievement,
		friend:         friendProvider,
	}
}

func NewUserMapper(data *usermodel.UserData) UserMapper {
	return UserMapper{data}
}

type UserMapper struct {
	data *usermodel.UserData
}

func (u UserMapper) days() int {
	return int(time.Now().Sub(u.data.AbstinenceTime.Time).Hours() / 24)
}

func (u UserMapper) Cigarette() int {
	return u.days() * int(u.data.CigaretteDayAmount)
}

func (u UserMapper) Life() int {
	return u.days() * 20
}

func (u UserMapper) Time() int {
	return u.Cigarette() * 5
}

func (u UserMapper) Money() float64 {
	cigaretteCost := float64(u.data.PackPrice) / float64(u.data.CigarettePackAmount)
	money := float64(u.Cigarette()) * cigaretteCost
	return money
}

func (u UserMapper) User(userId int64, friends []int64) *usermodel.User {
	return usermodel.NewUser(
		userId,
		u.data.AbstinenceTime.Time,
		u.Life(),
		u.Cigarette(),
		u.Time(),
		u.Money(),
		u.data.Motivation,
		u.data.WelcomeMotivation,
		u.data.Level,
		friends,
		u.data.Triggers,
	)
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
	friendsIds := u.userFriends.Friends(ctx, userId)
	IdsSorter(friendsIds).Sort()
	return UserMapper{userData}.User(userId, friendsIds), nil
}

func (u *userUseCase) Reset(ctx context.Context, userId int64) error {
	err := u.user.Reset(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("reset user", "Reset", _PROVIDER))
	}
	u.achievement.Clear(ctx, userId)
	u.privacySetting.Clear(ctx, userId)
	return nil
}

func (u *userUseCase) User(ctx context.Context, userId int64) (*usermodel.User, error) {
	userData, err := u.user.User(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user data", "User", _PROVIDER))
	}
	friendsIds := u.userFriends.Friends(ctx, userId)
	IdsSorter(friendsIds).Sort()
	return UserMapper{userData}.User(userId, friendsIds), nil
}

func (u *userUseCase) Level(ctx context.Context, userId int64) (*usermodel.LevelInfo, error) {
	level, err := u.user.Level(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user level", "Level", _PROVIDER))
	}
	return level, nil
}

func (u *userUseCase) Friends(ctx context.Context, friendsIds []int64) []*usermodel.Friend {
	var wg sync.WaitGroup
	var mu sync.Mutex
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
	FriendsSorter(friends).Sort()
	return friends
}
