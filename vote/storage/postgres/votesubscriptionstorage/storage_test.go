package votesubscriptionstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/vote/error/subscriptionerror"
	"github.com/Tap-Team/kurilka/vote/storage/postgres/votesubscriptionstorage"
	"gotest.tools/v3/assert"
)

var (
	storage *votesubscriptionstorage.Storage
	db      *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	d, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create new container, %s", err)
	}
	defer term(ctx)
	db = d
	storage = votesubscriptionstorage.New(db)
	os.Exit(m.Run())
}

func Test_CreateSubscription(t *testing.T) {
	ctx := context.Background()
	checker := NewChecker(db)
	inserter := NewInserter(db)

	cases := []struct {
		subscriptionId, userId int64
		err                    error
		userRegistered         bool
		subExists, userExists  bool
	}{
		{
			subscriptionId: rand.Int63(),
			userId:         rand.Int63(),

			subExists:  true,
			userExists: true,
		},
		{
			subscriptionId: 1951341261,
			userId:         58519581,
			subExists:      true,
			userExists:     true,
		},
		{
			subscriptionId: 1951341261,
			userId:         85181583,
			err:            subscriptionerror.SubscriptionIdExists,
			subExists:      true,
			userExists:     false,
		},
		{
			subscriptionId: 1858132012341,
			userId:         58519581,
			userRegistered: true,
			err:            subscriptionerror.SubscriptionIdExists,
			subExists:      false,
			userExists:     true,
		},
	}

	for _, cs := range cases {
		if !cs.userRegistered {
			assert.NilError(t, inserter.InsertEmptyUser(cs.userId))
		}
		err := storage.CreateSubscription(ctx, cs.subscriptionId, cs.userId)
		assert.ErrorIs(t, err, cs.err, "error not equal")

		assert.Equal(t, cs.subExists, checker.VoteSubscriptionExists(cs.subscriptionId), "subExists")
		assert.Equal(t, cs.userExists, checker.UserVoteSubscriptionExists(cs.userId), "userExists")
	}
}

func Test_DeleteSubscription(t *testing.T) {
	ctx := context.Background()
	checker := NewChecker(db)
	inserter := NewInserter(db)

	cases := []struct {
		subscriptionId, userId int64
		err                    error
		userRegistered         bool
		subExists, userExists  bool
	}{
		{
			subscriptionId: 19958391328581,
			userId:         230618141,
		},
		{
			subscriptionId: 19958391328581,
			userId:         230618141,
			userRegistered: true,
		},
	}

	for _, cs := range cases {
		if !cs.userRegistered {
			assert.NilError(t, inserter.InsertEmptyUser(cs.userId))
			assert.NilError(t, inserter.InsertSubscription(cs.subscriptionId, cs.userId), "failed insert subscription")
		}
		err := storage.DeleteSubscription(ctx, cs.subscriptionId)
		assert.ErrorIs(t, err, cs.err, "wrong err")
		assert.Equal(t, cs.subExists, checker.VoteSubscriptionExists(cs.subscriptionId), "subExists")
		assert.Equal(t, cs.userExists, checker.UserVoteSubscriptionExists(cs.userId), "userExists")
	}
}

func Test_UpdateSubscriptionId(t *testing.T) {
	ctx := context.Background()
	checker := NewChecker(db)
	inserter := NewInserter(db)

	cases := []struct {
		subscriptionId, userId int64
		err                    error
		userRegistered         bool

		updateSubscriptionId int64

		subExists, userExists bool
	}{
		{
			userId:         498624,
			subscriptionId: 1581358335,

			updateSubscriptionId: 78907342,

			subExists:  true,
			userExists: true,
		},

		{
			userId: 67589390,

			updateSubscriptionId: 38296821,

			subExists:  true,
			userExists: true,
		},

		{
			userId: 498624,

			updateSubscriptionId: 38296821,
			err:                  subscriptionerror.SubscriptionIdExists,
			userRegistered:       true,

			subExists:  true,
			userExists: true,
		},
	}

	for _, cs := range cases {
		if !cs.userRegistered {
			assert.NilError(t, inserter.InsertEmptyUser(cs.userId))
			assert.NilError(t, inserter.InsertSubscription(cs.subscriptionId, cs.userId), "failed insert subscription")
		}

		err := storage.UpdateUserSubscriptionId(ctx, cs.userId, cs.updateSubscriptionId)
		assert.ErrorIs(t, err, cs.err, "wrong err")
		assert.Equal(t, cs.subExists, checker.VoteSubscriptionExists(cs.updateSubscriptionId), "subExists")
		assert.Equal(t, cs.userExists, checker.UserVoteSubscriptionExists(cs.userId), "userExists")
	}
}

func Test_UserSubscriptionId(t *testing.T) {
	ctx := context.Background()
	inserter := NewInserter(db)

	cases := []struct {
		userId, subscriptionId int64

		userRegistered bool
	}{
		{
			userId:         438852349,
			subscriptionId: 238412,
		},
		{
			userId:         438852349,
			userRegistered: true,
			subscriptionId: 238412,
		},
	}

	for _, cs := range cases {
		if !cs.userRegistered {
			_, err := storage.UserSubscriptionId(ctx, cs.userId)
			assert.ErrorIs(t, err, subscriptionerror.SubscriptionNotFound, "wrong err")
			assert.NilError(t, inserter.InsertEmptyUser(cs.userId), "failed insert user")
			assert.NilError(t, inserter.InsertSubscription(cs.subscriptionId, cs.userId), "failed insert subscription")
		}

		subscriptionId, err := storage.UserSubscriptionId(ctx, cs.userId)
		assert.NilError(t, err, "failed get user subscription")
		assert.Equal(t, cs.subscriptionId, subscriptionId, "wrong subscription id")
	}
}
