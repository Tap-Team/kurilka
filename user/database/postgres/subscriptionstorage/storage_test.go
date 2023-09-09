package subscriptionstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/user/database/postgres/subscriptionstorage"
	"gotest.tools/v3/assert"
)

var (
	db      *sql.DB
	storage *subscriptionstorage.Storage
)

func TestMain(m *testing.M) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	migrationFolder := amidsql.DEFAULT_MIGRATION_PATH
	database, term, err := amidsql.NewContainer(ctx, migrationFolder)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)

	db = database
	storage = subscriptionstorage.New(db)

	os.Exit(m.Run())
}

func TestUpdateSubscription(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		before struct {
			userId       int64
			createUser   *usermodel.CreateUser
			subscription usermodel.Subscription
		}
		expected struct {
			subscription usermodel.Subscription
			updateErr    error
			getErr       error
		}
	}{
		{
			before: struct {
				userId       int64
				createUser   *usermodel.CreateUser
				subscription usermodel.Subscription
			}{
				userId: rand.Int63(),
				createUser: usermodel.NewCreateUser(
					"pofig",
					10,
					20,
					140.89,
				),
				subscription: usermodel.NewSubscription(
					usermodel.NONE,
					time.Time{},
				),
			},
			expected: struct {
				subscription usermodel.Subscription
				updateErr    error
				getErr       error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour*24*30)),
			},
		},

		{
			before: struct {
				userId       int64
				createUser   *usermodel.CreateUser
				subscription usermodel.Subscription
			}{
				userId: rand.Int63(),
				createUser: usermodel.NewCreateUser(
					"pofig",
					10,
					20,
					140.89,
				),
				subscription: usermodel.NewSubscription(
					usermodel.TRIAL,
					time.Time{},
				),
			},
			expected: struct {
				subscription usermodel.Subscription
				updateErr    error
				getErr       error
			}{
				subscription: usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour*24*30)),
			},
		},

		{
			before: struct {
				userId       int64
				createUser   *usermodel.CreateUser
				subscription usermodel.Subscription
			}{},
			expected: struct {
				subscription usermodel.Subscription
				updateErr    error
				getErr       error
			}{
				updateErr: usererror.ExceptionUserNotFound(),
				getErr:    usererror.ExceptionUserNotFound(),
			},
		},
	}

	for _, cs := range cases {
		err := insertUserWithSubscription(db, cs.before.userId, cs.before.createUser, cs.before.subscription)
		assert.NilError(t, err, "failed insert user subscription")

		err = storage.UpdateUserSubscription(ctx, cs.before.userId, cs.expected.subscription)
		assert.ErrorIs(t, err, cs.expected.updateErr)

		subscription, err := storage.UserSubscription(ctx, cs.before.userId)
		assert.ErrorIs(t, err, cs.expected.updateErr, "update error not equal")

		assert.Equal(t, subscription.Expired.Time.Unix(), cs.expected.subscription.Expired.Time.Unix(), "subscription time not equal")
		assert.Equal(t, subscription.Type, cs.expected.subscription.Type, "subscription time not equal")

		assert.ErrorIs(t, err, cs.expected.getErr, "get error not equal")
	}
}

func TestUserSubscription(t *testing.T) {

}
