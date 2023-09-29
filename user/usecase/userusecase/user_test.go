package userusecase_test

import (
	"context"
	"errors"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievementmessagesender"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/user/usecase/userusecase"
	"github.com/Tap-Team/kurilka/workers"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	NilUser   *usermodel.User
	NilLevel  *usermodel.LevelInfo
	NilFriend *usermodel.Friend
)

func TestUserMapper(t *testing.T) {
	cases := []struct {
		minutes             int
		cigaretteDayAmount  uint8
		cigarettePackAmount uint8
		packPrice           float32

		money     float64
		life      int
		cigarette int
		time      int
	}{
		{
			minutes:             3245,
			cigaretteDayAmount:  45,
			cigarettePackAmount: 20,
			packPrice:           178.50,

			money:     901.425,
			life:      45,
			cigarette: 101,
			time:      505,
		},

		{
			minutes:             8380,
			cigaretteDayAmount:  99,
			cigarettePackAmount: 99,
			packPrice:           4999.99,

			money:     29_090.8509,
			life:      116,
			cigarette: 576,
			time:      2880,
		},
	}

	for _, cs := range cases {
		userData := NewUserDataMinutes(cs.minutes, cs.cigaretteDayAmount, cs.cigarettePackAmount, cs.packPrice)
		mapper := userusecase.NewUserMapper(userData, time.Now())

		userId := rand.Int63()
		user := mapper.User(userId, usermodel.Subscription{})
		friend := mapper.Friend(userId, make([]*usermodel.Achievement, 0), make([]usermodel.PrivacySetting, 0), usermodel.NONE)

		assert.Equal(t, true, moneyEqual(user.Money, cs.money), "user money not equal")
		assert.Equal(t, true, moneyEqual(usermodel.Money(mapper.Money()), cs.money), "mapper money not equal")
		assert.Equal(t, true, moneyEqual(friend.Money, cs.money), "friend money not equal")

		assert.Equal(t, cs.life, user.Life, "user life not equal")
		assert.Equal(t, cs.life, mapper.Life(), "mapper life not equal")
		assert.Equal(t, cs.life, friend.Life, "friend life not equal")

		assert.Equal(t, cs.time, user.Time, "user time not equal")
		assert.Equal(t, cs.time, mapper.Time(), "mapper time not equal")
		assert.Equal(t, cs.time, friend.Time, "friend time not equal")

		assert.Equal(t, cs.cigarette, user.Cigarette, "user cigarette not equal")
		assert.Equal(t, cs.cigarette, mapper.Cigarette(), "mapper cigarette not equal")
		assert.Equal(t, cs.cigarette, friend.Cigarette, "friend cigarette not equal")
	}
}

func Test_UserMapper_Triggers(t *testing.T) {
	cases := []struct {
		triggers     []usermodel.Trigger
		subscription usermodel.Subscription

		expectedTriggers []usermodel.Trigger
	}{
		{
			triggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
			subscription: usermodel.NewSubscription(usermodel.NONE, time.Time{}),
			expectedTriggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
		},
		{
			triggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
			subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour)),
			expectedTriggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
		},
		{
			triggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
			subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour)),
			expectedTriggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_HEALTH,
			},
		},
		{
			triggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
			subscription: usermodel.NewSubscription(usermodel.TRIAL, time.Time{}),
			expectedTriggers: []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.SUPPORT_TRIAL,
				usermodel.SUPPORT_HEALTH,
			},
		},
	}

	for _, cs := range cases {
		mp := userusecase.NewUserMapper(usermodel.NewUserData("", 1, 1, 1, "", "", time.Now(), usermodel.LevelInfo{}, cs.triggers), time.Now())

		us := mp.User(0, cs.subscription)

		equal := slices.Equal(cs.expectedTriggers, us.Triggers)
		assert.Equal(t, true, equal, "triggers not equal")
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	userManager := userdatamanager.NewMockUserManager(ctrl)
	userFriendsProvider := userusecase.NewMockUserFriendsProvider(ctrl)
	achievementsProvider := achievementdatamanager.NewMockAchievementManager(ctrl)
	friendsProvider := userusecase.NewMockFriendProvider(ctrl)
	subscription := userusecase.NewMockSubscriptionStorage(ctrl)
	subscription.EXPECT().UserSubscription(gomock.Any(), gomock.Any()).Return(usermodel.Subscription{}, nil).AnyTimes()
	userWorker := workers.NewMockUserWorker(ctrl)
	achievementMessageSender := achievementmessagesender.NewMockAchievementMessageSenderAtTime(ctrl)

	useCase := userusecase.NewUser(userFriendsProvider, userManager, privacySettingsManager, achievementsProvider, friendsProvider, subscription, userWorker, achievementMessageSender)

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()

		expectedErr := errors.New("any error")

		userManager.EXPECT().Create(gomock.Any(), userId, &createUser).Return(nil, expectedErr).Times(1)

		user, err := useCase.Create(ctx, userId, &createUser)

		assert.ErrorIs(t, err, expectedErr, "wrong err")
		assert.Equal(t, user, NilUser, "wrong user")
	}

	{
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()

		expectedUser := random.StructTyped[usermodel.UserData]()

		userManager.EXPECT().Create(gomock.Any(), userId, &createUser).Return(&expectedUser, nil).Times(1)
		userWorker.EXPECT().AddUser(gomock.Any(), workers.NewUser(userId, expectedUser.AbstinenceTime.Time)).Times(1)

		user, err := useCase.Create(ctx, userId, &createUser)

		assert.NilError(t, err, "non nil err")
		assert.Equal(t, userId, user.ID, "id not equal")
	}
}

func TestReset(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	userManager := userdatamanager.NewMockUserManager(ctrl)
	userFriendsProvider := userusecase.NewMockUserFriendsProvider(ctrl)
	achievementsProvider := achievementdatamanager.NewMockAchievementManager(ctrl)
	friendProvider := userusecase.NewMockFriendProvider(ctrl)
	subscription := userusecase.NewMockSubscriptionStorage(ctrl)
	userWorker := workers.NewMockUserWorker(ctrl)
	achievementMessageSender := achievementmessagesender.NewMockAchievementMessageSenderAtTime(ctrl)

	useCase := userusecase.NewUser(userFriendsProvider, userManager, privacySettingsManager, achievementsProvider, friendProvider, subscription, userWorker, achievementMessageSender)

	{
		userId := rand.Int63()
		expectedErr := errors.New("any")

		userManager.EXPECT().Reset(gomock.Any(), userId).Return(expectedErr).Times(1)

		err := useCase.Reset(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "error not equal")
	}

	{
		userId := rand.Int63()

		userManager.EXPECT().Reset(gomock.Any(), userId).Return(nil).Times(1)

		privacySettingsManager.EXPECT().Clear(gomock.Any(), userId).Times(1)
		achievementsProvider.EXPECT().Clear(gomock.Any(), userId).Times(1)
		userWorker.EXPECT().RemoveUser(gomock.Any(), userId).Times(1)
		subscription.EXPECT().Clear(gomock.Any(), userId).Times(1)
		achievementMessageSender.EXPECT().CancelSendMessagesForUser(gomock.Any(), userId).Times(1)

		err := useCase.Reset(ctx, userId)

		assert.NilError(t, err, "non nil error")
	}
}

func TestLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	userManager := userdatamanager.NewMockUserManager(ctrl)
	userFriendsProvider := userusecase.NewMockUserFriendsProvider(ctrl)
	achievementsProvider := achievementdatamanager.NewMockAchievementManager(ctrl)
	friendsProvider := userusecase.NewMockFriendProvider(ctrl)
	subscription := userusecase.NewMockSubscriptionStorage(ctrl)

	useCase := userusecase.NewUser(userFriendsProvider, userManager, privacySettingsManager, achievementsProvider, friendsProvider, subscription, nil, nil)

	{
		userId := rand.Int63()
		expectedErr := errors.New("sdf;aksldj;flkaslkdfjasjdfljspoiqrwenzxncv")

		userManager.EXPECT().Level(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		level, err := useCase.Level(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "error not equal")
		assert.Equal(t, level, NilLevel, "level not nil")
	}

	{
		userId := rand.Int63()
		expectedLevel := random.StructTyped[usermodel.LevelInfo]()

		userManager.EXPECT().Level(gomock.Any(), userId).Return(&expectedLevel, nil).Times(1)

		level, err := useCase.Level(ctx, userId)

		assert.NilError(t, err, "non nil err")
		assert.Equal(t, level, &expectedLevel)
	}
}

func TestUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	userManager := userdatamanager.NewMockUserManager(ctrl)
	userFriendsProvider := userusecase.NewMockUserFriendsProvider(ctrl)
	achievementsProvider := achievementdatamanager.NewMockAchievementManager(ctrl)
	friendsProvider := userusecase.NewMockFriendProvider(ctrl)
	subscription := userusecase.NewMockSubscriptionStorage(ctrl)
	subscription.EXPECT().UserSubscription(gomock.Any(), gomock.Any()).Return(usermodel.Subscription{}, nil).AnyTimes()

	useCase := userusecase.NewUser(userFriendsProvider, userManager, privacySettingsManager, achievementsProvider, friendsProvider, subscription, nil, nil)

	{
		userId := rand.Int63()

		expectedErr := errors.New("failed")

		userManager.EXPECT().User(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		user, err := useCase.User(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
		assert.Equal(t, user, NilUser, "wrong user")

	}

	{

		userId := rand.Int63()

		expectedUser := random.StructTyped[usermodel.UserData]()

		userManager.EXPECT().User(gomock.Any(), userId).Return(&expectedUser, nil).Times(1)

		user, err := useCase.User(ctx, userId)

		assert.NilError(t, err, "non nil err")
		assert.Equal(t, userId, user.ID, "id not equal")
	}
}

func Test_Wrapper_UserFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	userManager := userdatamanager.NewMockUserManager(ctrl)
	userFriendsProvider := userusecase.NewMockUserFriendsProvider(ctrl)

	friendsProvider := userusecase.NewUserFriendsProvider(userFriendsProvider, userManager)

	{
		userId := rand.Int63()

		storageFriends := []int64{123, 55, 7, 81, 1, 35, 5, 550}
		existsFriends := []int64{123, 7, 1, 5}

		userFriendsProvider.EXPECT().Friends(gomock.Any(), userId).Return(storageFriends).Times(1)
		userManager.EXPECT().FilterExists(gomock.Any(), storageFriends).Return(existsFriends).Times(1)

		friends := friendsProvider.Friends(ctx, userId)

		equal := slices.Equal(existsFriends, friends)
		assert.Equal(t, true, equal, "equals")
	}
}

func TestFriend(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	userManager := userdatamanager.NewMockUserManager(ctrl)
	achievementsProvider := achievementdatamanager.NewMockAchievementManager(ctrl)
	subscriptionStorage := userusecase.NewMockSubscriptionStorage(ctrl)

	friendProvider := userusecase.NewFriendProvider(achievementsProvider, userManager, privacySettingsManager, subscriptionStorage)

	{
		friendId := rand.Int63()

		expectedErr := errors.New("failed get user data")

		userManager.EXPECT().User(gomock.Any(), friendId).Return(nil, expectedErr).Times(1)

		friend, err := friendProvider.Friend(ctx, friendId)

		assert.Equal(t, friend, NilFriend, "non nil friend")
		assert.ErrorIs(t, err, expectedErr, "wrong error")
	}

	subscriptionCases := []struct {
		subscription usermodel.Subscription
		err          error

		expectedSubscription usermodel.SubscriptionType
	}{
		{
			err:                  errors.New("any error"),
			expectedSubscription: usermodel.NONE,
		},
		{
			subscription: usermodel.Subscription{
				Type: usermodel.NONE,
			},
			expectedSubscription: usermodel.NONE,
		},
		{
			subscription: usermodel.Subscription{
				Type: usermodel.TRIAL,
			},
			expectedSubscription: usermodel.NONE,
		},
		{
			subscription: usermodel.Subscription{
				Type: usermodel.BASIC,
			},
			expectedSubscription: usermodel.NONE,
		},
		{
			subscription: usermodel.Subscription{
				Type:    usermodel.BASIC,
				Expired: amidtime.Timestamp{Time: time.Now().Add(time.Hour)},
			},
			expectedSubscription: usermodel.BASIC,
		},
		{
			subscription: usermodel.Subscription{
				Type:    usermodel.TRIAL,
				Expired: amidtime.Timestamp{Time: time.Now().Add(time.Hour)},
			},
			expectedSubscription: usermodel.TRIAL,
		},
	}

	for _, cs := range subscriptionCases {
		friendId := rand.Int63()

		userData := random.StructTyped[usermodel.UserData]()

		privacySettings := []usermodel.PrivacySetting{
			usermodel.ACHIEVEMENTS_CIGARETTE,
			usermodel.STATISTICS_CIGARETTE,
			usermodel.ACHIEVEMENTS_SAVING,
			usermodel.STATISTICS_MONEY,
		}

		achievementList := randomAchievementList(10)
		subscription := cs.subscription

		userManager.EXPECT().User(gomock.Any(), friendId).Return(&userData, nil).Times(1)

		achievementsProvider.EXPECT().AchievementPreview(gomock.Any(), friendId).Return(achievementList).Times(1)
		privacySettingsManager.EXPECT().PrivacySettings(gomock.Any(), friendId).Return(privacySettings, nil).Times(1)
		subscriptionStorage.EXPECT().UserSubscription(gomock.Any(), friendId).Return(subscription, cs.err).Times(1)

		friend, err := friendProvider.Friend(ctx, friendId)

		assert.NilError(t, err, "error not nil")

		equal := slices.Equal(privacySettings, friend.PrivacySettings)
		assert.Equal(t, true, equal, "privacy settings not equal")
		assert.Equal(t, cs.expectedSubscription, friend.SubscriptionType, "subscription type not equal")

		equal = slices.EqualFunc(achievementList, friend.Achievements, func(a1, a2 *usermodel.Achievement) bool {
			if a1.Level != a2.Level {
				return false
			}
			if a1.Type != a2.Type {
				return false
			}
			return true
		})

		assert.Equal(t, true, equal, "achievements not equal")
	}

}

type Number interface {
	int64
}

type IntMather[T Number] struct {
	matches func(T) bool
}

func (m *IntMather[T]) Matches(x interface{}) bool {
	i, ok := x.(T)
	if !ok {
		return false
	}
	return m.matches(i)
}

func (m *IntMather[T]) String() string {
	return "is matches func"
}

func NewIntMatcher[T Number](matches func(T) bool) gomock.Matcher {
	return &IntMather[T]{matches: matches}
}

func TestFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	privacySettingsManager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)
	userManager := userdatamanager.NewMockUserManager(ctrl)
	userFriendsProvider := userusecase.NewMockUserFriendsProvider(ctrl)
	achievementsProvider := achievementdatamanager.NewMockAchievementManager(ctrl)
	friendsProvider := userusecase.NewMockFriendProvider(ctrl)
	subscription := userusecase.NewMockSubscriptionStorage(ctrl)

	useCase := userusecase.NewUser(userFriendsProvider, userManager, privacySettingsManager, achievementsProvider, friendsProvider, subscription, nil, nil)

	{
		friends := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

		oddElementsCount := reduce(friends, func(v int, el int64) int {
			if el%2 == 1 {
				v++
			}
			return v
		})
		evenElementsCount := len(friends) - oddElementsCount

		oddMatcher := NewIntMatcher[int64](func(i int64) bool { return i%2 == 1 })
		evenMatcher := NewIntMatcher[int64](func(i int64) bool { return i%2 == 0 })
		friendsProvider.EXPECT().Friend(gomock.Any(), oddMatcher).Return(&usermodel.Friend{ID: rand.Int63()}, nil).Times(oddElementsCount)
		friendsProvider.EXPECT().Friend(gomock.Any(), evenMatcher).Return(nil, errors.New("failed get userdata")).Times(evenElementsCount)

		frs := useCase.Friends(ctx, friends)

		assert.Equal(t, len(frs), oddElementsCount, "wrong len of friends")

		// check slice is sorted
		minId := 0
		for _, fr := range frs {
			assert.Equal(t, true, fr.ID > int64(minId))
		}
	}
}

func reduce[E any, V any](collection []E, f func(V, E) V) V {
	var value V
	for _, e := range collection {
		value = f(value, e)
	}
	return value
}
