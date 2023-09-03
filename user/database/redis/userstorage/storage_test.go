package userstorage_test

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/Tap-Team/kurilka/user/database/redis/userstorage"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	storage     *userstorage.Storage
	redisClient *redis.Client
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	rc, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed start redis container, %s", err)
	}
	defer term(ctx)
	redisClient = rc
	storage = userstorage.New(rc, 0)

	m.Run()
}

func TestCRUD(t *testing.T) {
	var err error
	ctx := context.Background()

	require.NoError(t, storage.RemoveUser(ctx, rand.Int63()), "error from unexists user")

	_, err = storage.User(ctx, rand.Int63())
	require.ErrorIs(t, err, usererror.ExceptionUserNotFound(), "wrong error")

	const listSize = 100

	userIds := make([]int64, 0, listSize)

	for i := 0; i < listSize; i++ {
		userId := rand.Int63()
		user := random.StructTyped[usermodel.UserData]()
		err := storage.SaveUser(ctx, userId, &user)
		require.NoError(t, err, "failed save user in storage")

		u, err := storage.User(ctx, userId)
		require.NoError(t, err, "failed get user from storage")

		require.Equal(t, user.AbstinenceTime.Unix(), u.AbstinenceTime.Unix(), "abstinence time not equal")
		require.Equal(t, user.Subscription.Expired.Unix(), u.Subscription.Expired.Time.Unix(), "subscription expired not equal")
		user.AbstinenceTime = u.AbstinenceTime
		user.Subscription.Expired = u.Subscription.Expired
		require.Equal(t, user, *u, "users not equal")
	}

	for _, userId := range userIds {
		err := storage.RemoveUser(ctx, userId)
		require.NoError(t, err, "failed remove user")

		_, err = storage.User(ctx, userId)
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound(), "wrong not found user error")
	}
}

func getUser(ctx context.Context, rc *redis.Client, userId int64) (user usermodel.UserData, err error) {
	err = rc.Get(ctx, redishelper.UsersKey(userId)).Scan(&user)
	return
}

func TestSaveUser_Case_Update_After_Save(t *testing.T) {
	var err error
	ctx := context.Background()

	const count = 100
	for i := 0; i < count; i++ {
		userId := rand.Int63()
		user := random.StructTyped[usermodel.UserData]()

		err = storage.SaveUser(ctx, userId, &user)
		require.NoError(t, err, "failed save user")

		u, err := getUser(ctx, redisClient, userId)
		require.NoError(t, err, "failed get user from storage")

		require.Equal(t, user.AbstinenceTime.Unix(), u.AbstinenceTime.Unix(), "abstinence time not equal")
		require.Equal(t, user.Subscription.Expired.Unix(), u.Subscription.Expired.Time.Unix(), "subscription expired not equal")
		user.AbstinenceTime = u.AbstinenceTime
		user.Subscription.Expired = u.Subscription.Expired
		require.Equal(t, user, u, "users not equal")

		user = random.StructTyped[usermodel.UserData]()
		err = storage.SaveUser(ctx, userId, &user)
		require.NoError(t, err, "failed save user")

		u, err = getUser(ctx, redisClient, userId)
		require.NoError(t, err, "failed get user from storage")

		require.Equal(t, user.AbstinenceTime.Unix(), u.AbstinenceTime.Unix(), "abstinence time not equal")
		require.Equal(t, user.Subscription.Expired.Unix(), u.Subscription.Expired.Time.Unix(), "subscription expired not equal")
		user.AbstinenceTime = u.AbstinenceTime
		user.Subscription.Expired = u.Subscription.Expired
		require.Equal(t, user, u, "users not equal")
	}

}
