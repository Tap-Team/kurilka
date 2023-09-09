package userstorage_test

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
	amidsql "github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/user/database/postgres/userstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	migrationFolder = amidsql.DEFAULT_MIGRATION_PATH
	trialPeriod     = time.Hour * 24 * 30
)

var (
	db      *sql.DB
	storage *userstorage.Storage
)

func TestMain(m *testing.M) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	conn, term, err := amidsql.NewContainer(ctx, migrationFolder)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	db = conn
	storage = userstorage.New(db, trialPeriod)
	os.Exit(m.Run())
}

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	triggers := map[usermodel.Trigger]struct{}{
		usermodel.SUPPORT_CIGGARETTE: {},
		usermodel.SUPPORT_HEALTH:     {},
		usermodel.SUPPORT_TRIAL:      {},
		usermodel.THANK_YOU:          {},
	}
	for i := 0; i < 50; i++ {
		createUser := random.StructTyped[usermodel.CreateUser]()

		userId := rand.Int63()

		user, err := storage.InsertUser(ctx, userId, &createUser)
		require.NoError(t, err, "failed insert user")

		for _, tr := range user.Triggers {
			_, ok := triggers[tr]
			require.True(t, ok, "trigger %s not found", tr)
		}
		userSubscription, err := userSubscription(db, userId)
		require.NoError(t, err, "failed get user subscription")

		require.Equal(t, userSubscription.Type, usermodel.TRIAL, "wrong subscription type")
		require.Equal(t, user.Level.Level, usermodel.One, "wrong init level")
		require.Equal(t, user.Level.MinExp, 0, "wrong first level min exp")
		require.False(t, len(user.Motivation) == 0, "wrong user motivation")
		require.False(t, len(user.WelcomeMotivation) == 0, "wrong user welcome motivation")
		require.NotEqual(t, user.Motivation, user.WelcomeMotivation, "motivations are equal")

		userFromDatabase, err := storage.User(ctx, userId)
		require.NoError(t, err, "failed get user from storage")

		for _, tr := range userFromDatabase.Triggers {
			_, ok := triggers[tr]
			require.True(t, ok, "trigger %s not found", tr)
		}

		user.AbstinenceTime = userFromDatabase.AbstinenceTime
		user.Triggers = userFromDatabase.Triggers
		require.Equal(t, user, userFromDatabase, "users not equal")
		require.False(t, user.AbstinenceTime.IsZero(), "wrong abstinence time")

	}

}

func TestUserExp(t *testing.T) {
	const listSize = 10
	var err error
	ctx := context.Background()

	userIds := make([]int64, 0, listSize)
	for i := 0; i < listSize; i++ {
		userId := rand.Int63()
		user := random.StructTyped[usermodel.CreateUser]()

		err = insertUser(db, userId, user)
		require.NoError(t, err, "failed insert user")

		userIds = append(userIds, userId)
		exp, err := storage.UserExp(ctx, userId)
		require.NoError(t, err, "failed get exp")
		require.Equal(t, 0, exp, "wrong exp")

		expectedExp := exp

		for i := 0; i < 10; i++ {
			exp := random.Range{Max: 1000}.Int()
			addUserExp(db, userId, exp)
			fakeAddUserExp(db, userId, exp)
			expectedExp += exp
		}

		exp, err = storage.UserExp(ctx, userId)
		require.NoError(t, err, "failed get user exp after change")
		require.Equal(t, expectedExp, exp, "wrong user exp")
	}

}

func TestUserLevel(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		minExp, maxExp int
		level          usermodel.Level
	}{
		{
			minExp: 0,
			maxExp: 99,
			level:  usermodel.One,
		},
		{
			minExp: 100,
			maxExp: 199,
			level:  usermodel.Two,
		},
		{
			minExp: 200,
			maxExp: 299,
			level:  usermodel.Three,
		},
		{
			minExp: 300,
			maxExp: 399,
			level:  usermodel.Four,
		},
		{
			minExp: 400,
			maxExp: 499,
			level:  usermodel.Five,
		},
		{
			minExp: 500,
			maxExp: 599,
			level:  usermodel.Six,
		},
		{
			minExp: 600,
			maxExp: 699,
			level:  usermodel.Seven,
		},
		{
			minExp: 700,
			maxExp: 799,
			level:  usermodel.Eight,
		},
		{
			minExp: 800,
			maxExp: 899,
			level:  usermodel.Nine,
		},
		{
			minExp: 900,
			maxExp: 1000,
			level:  usermodel.Ten,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		createuser := random.StructTyped[usermodel.CreateUser]()

		_, err := storage.InsertUser(ctx, userId, &createuser)
		require.NoError(t, err, "failed insert user")

		rng := random.Range{Min: int64(cs.minExp), Max: int64(cs.maxExp)}
		exp := rng.Int()
		err = addUserExp(db, userId, exp)
		require.NoError(t, err, "failed add exp to user")

		level, err := storage.UserLevel(ctx, userId)
		require.NoError(t, err, "failed get user from storage")

		require.Equal(t, cs.level, level.Level, "user level not equal")
		require.Equal(t, exp, level.Exp, "wrong user exp")
		require.Equal(t, cs.minExp, level.MinExp, "min exp not equal")
		require.Equal(t, cs.maxExp, level.MaxExp, "max exp not equal")
	}
}

func TestExists(t *testing.T) {
	ctx := context.Background()

	friendsIds := make([]int64, 0)
	existsUsers := make([]int64, 0)
	for i := 0; i < 1000; i++ {
		id := rand.Int63()
		if rand.Intn(3)%2 == 0 {
			err := insertUser(db, id, random.StructTyped[usermodel.CreateUser]())
			require.NoError(t, err, "failed insert user")
			if rand.Intn(3)%2 == 0 {
				existsUsers = append(existsUsers, id)
			} else {
				removeUser(db, id)
			}
		}
		friendsIds = append(friendsIds, id)
	}
	sort.Slice(existsUsers, func(i, j int) bool {
		return existsUsers[i] < existsUsers[j]
	})
	users := storage.Exists(ctx, friendsIds)
	ok := slices.Equal(existsUsers, users)
	require.True(t, ok, "slices not equal")
}

func TestUserDeleted(t *testing.T) {
	ctx := context.Background()
	{
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()), "failed insert user")

		deleted, err := storage.UserDeleted(ctx, userId)
		require.NoError(t, err, "get user deleted")
		assert.False(t, deleted, "wrong deleted")
	}

	{
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()), "failed insert user")

		removeUser(db, userId)
		deleted, err := storage.UserDeleted(ctx, userId)
		require.NoError(t, err, "get user deleted")
		assert.True(t, deleted, "wrong deleted")
	}

	{
		userId := rand.Int63()
		_, err := storage.UserDeleted(ctx, userId)
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound(), "wrong error")
	}

}
