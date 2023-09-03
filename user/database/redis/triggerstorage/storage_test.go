package triggerstorage_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/Tap-Team/kurilka/user/database/redis/triggerstorage"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	rc      *redis.Client
	storage *triggerstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	redisClient, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed to start user container, %s", err)
	}
	defer term(ctx)
	rc = redisClient
	storage = triggerstorage.New(rc, 0)
	os.Exit(m.Run())
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()

	triggers := []usermodel.Trigger{
		usermodel.SUPPORT_CIGGARETTE,
		usermodel.SUPPORT_HEALTH,
		usermodel.THANK_YOU,
	}

	{
		userId := rand.Int63()
		_, err := storage.UserTriggers(ctx, userId)
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound())

		err = storage.RemoveUserTriggers(ctx, userId)
		require.NoError(t, err, "failed remove user")

		err = storage.SaveUserTriggers(ctx, userId, triggers)
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound())
	}

	{
		userId := rand.Int63()
		user := random.StructTyped[usermodel.UserData]()
		user.Triggers = triggers
		setUser(rc, userId, &user)

		trs, err := storage.UserTriggers(ctx, userId)
		require.NoError(t, err)

		ok := slices.Equal(triggers, trs)
		require.True(t, ok, "triggers not equal")

		err = storage.SaveUserTriggers(ctx, userId, make([]usermodel.Trigger, 0))
		require.NoError(t, err, "failed save triggers")

		trs, err = storage.UserTriggers(ctx, userId)
		require.NoError(t, err, "failed get user triggers")
		require.Equal(t, 0, len(trs), "triggers not updated")

		err = storage.RemoveUserTriggers(ctx, userId)
		require.NoError(t, err, "failed remove user triggers")

		_, err = storage.UserTriggers(ctx, userId)
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound())
	}
}
