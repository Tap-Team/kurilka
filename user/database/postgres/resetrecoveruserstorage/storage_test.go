package resetrecoveruserstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/user/database/postgres/resetrecoveruserstorage"
	"github.com/stretchr/testify/require"
)

var (
	db      *sql.DB
	storage *resetrecoveruserstorage.Storage
)

func TestMain(m *testing.M) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	d, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container")
	}
	defer term(ctx)
	db = d
	storage = resetrecoveruserstorage.New(db)
	os.Exit(m.Run())
}

func TestReset(t *testing.T) {
	ctx := context.Background()
	{
		subscription := usermodel.NewSubscription(usermodel.NONE, time.Now())

		userId := rand.Int63()
		user := random.StructTyped[usermodel.CreateUser]()
		err := insertUser(db, userId, &user, subscription, make([]usermodel.Trigger, 0))
		require.NoError(t, err, "failed insert user")

		err = storage.ResetUser(ctx, userId)
		require.NoError(t, err, "failed reset user")

		deleted, err := userDeleted(db, userId)
		require.NoError(t, err, "failed get user deleted")
		require.True(t, deleted, "user not mark as deleted")

		ach, err := achievementsCount(db, userId)
		require.True(t, err == nil && ach == 0)
		pr, err := privacySettingsCount(db, userId)
		require.True(t, err == nil && pr == 0)
	}

	{
		err := storage.ResetUser(ctx, rand.Int63())
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound())
	}

}

func TestRecover(t *testing.T) {
	ctx := context.Background()

	{
		subscription := usermodel.NewSubscription(usermodel.NONE, time.Now())
		userId := rand.Int63()
		user := random.StructTyped[usermodel.CreateUser]()
		err := insertUser(db, userId, &user, subscription, make([]usermodel.Trigger, 0))
		require.NoError(t, err, "failed insert user")

		_, err = storage.RecoverUser(ctx, userId, &user)
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound())
	}

	{
		subscription := usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour*40))

		triggers := []usermodel.Trigger{
			usermodel.SUPPORT_CIGGARETTE,
			usermodel.SUPPORT_HEALTH,
			usermodel.THANK_YOU,
			usermodel.SUPPORT_TRIAL,
		}
		userId := rand.Int63()
		createUser := random.StructTyped[usermodel.CreateUser]()
		err := insertUser(db, userId, &createUser, subscription, triggers)
		require.NoError(t, err, "failed insert user")

		err = storage.ResetUser(ctx, userId)
		require.NoError(t, err, "failed reset user")

		user, err := storage.RecoverUser(ctx, userId, &createUser)
		require.NoError(t, err, "failed recover user")

		sort.Slice(triggers, func(i, j int) bool { return triggers[i] > triggers[j] })
		sort.Slice(user.Triggers, func(i, j int) bool { return user.Triggers[i] > user.Triggers[j] })

		ok := slices.Equal(triggers, user.Triggers)
		require.True(t, ok, "slices not equal")

		require.Equal(t, user.Subscription.Type, subscription.Type)
		require.Equal(t, user.Subscription.Expired.Unix(), subscription.Expired.Unix())

		verifyRecoveredUser(t, user)
	}
}
