package userstorage_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/pkg/random"
	amidsql "github.com/Tap-Team/kurilka/pkg/sql"
	"github.com/Tap-Team/kurilka/user/database/userstorage"
	"github.com/Tap-Team/kurilka/user/model/usermodel"
	"github.com/stretchr/testify/require"
)

const (
	migrationFolder = amidsql.DEFAULT_MIGRATION_PATH
	trialPeriod     = time.Hour * 24 * 30
)

var storage *userstorage.Storage

func TestMain(m *testing.M) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	db, term, err := amidsql.NewContainer(ctx, migrationFolder)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	storage = userstorage.New(db, trialPeriod)
	m.Run()
}

func TestUserCRUD(t *testing.T) {
	var err error
	ctx := context.Background()

	createUser := random.StructTyped[usermodel.CreateUser]()

	userId := rand.Int63()

	user, err := storage.InsertUser(ctx, userId, &createUser)
	require.NoError(t, err, "failed insert user")

	require.Equal(t, user.Subscription.Type, usermodel.TRIAL, "wrong subscription type")
	require.Equal(t, user.Level.Level, usermodel.One, "wrong init level")
	require.Equal(t, user.Level.MinExp, 0, "wrong first level min exp")

	userFromDatabase, err := storage.User(ctx, userId)
	require.NoError(t, err, "failed get user from storage")
	require.Equal(t, user.Subscription.Expired.Unix(), userFromDatabase.Subscription.Expired.Unix(), "time not equal")
	user.Subscription.Expired = userFromDatabase.Subscription.Expired
	require.Equal(t, user, userFromDatabase, "users not equal")

	err = storage.DeleteUser(ctx, userId)
	require.NoError(t, err, "failed delete user")

	_, err = storage.User(ctx, userId)
	require.ErrorIs(t, err, usererror.ExceptionUserNotFound(), "wrong select user error")
}
