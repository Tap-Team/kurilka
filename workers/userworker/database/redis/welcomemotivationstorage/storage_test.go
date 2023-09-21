package welcomemotivationstorage_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/Tap-Team/kurilka/workers/userworker/database/redis/welcomemotivationstorage"
	"github.com/redis/go-redis/v9"
	"gotest.tools/v3/assert"
)

var (
	rc      *redis.Client
	storage *welcomemotivationstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	redis, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed create redis container, %s", err)
	}
	defer term(ctx)
	rc = redis
	storage = welcomemotivationstorage.New(rc, 0)
	os.Exit(m.Run())
}

func Test_Storage_SaveUserWelcomeMotivation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		welcomeMotivation string

		updateWelcomeMotivation string

		expectedWelcomeMotivation string

		expectedErr error
	}{
		{
			welcomeMotivation:         "Hello!",
			updateWelcomeMotivation:   "Bye!",
			expectedWelcomeMotivation: "Bye!",
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()

		err := saveUser(rc, userId, cs.welcomeMotivation)
		assert.NilError(t, err, "failed save user welcome motivation")

		err = storage.SaveUserWelcomeMotivation(ctx, userId, cs.updateWelcomeMotivation)
		assert.ErrorIs(t, err, cs.expectedErr)

		m, err := userMotivation(rc, userId)
		assert.NilError(t, err, "failed get user motivation")
		assert.Equal(t, m, cs.expectedWelcomeMotivation)
	}
}

func Test_Storage_RemoveUserWelcomeMotivation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		welcomeMotivation string
		expectedErr       error
	}{
		{
			welcomeMotivation: random.String(10),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()

		err := saveUser(rc, userId, cs.welcomeMotivation)
		assert.NilError(t, err, "failed save user welcome motivation")

		err = storage.RemoveUserWelcomeMotivation(ctx, userId)
		assert.ErrorIs(t, err, cs.expectedErr, "wrong err welcome motivaiton")

		_, err = userMotivation(rc, userId)
		assert.ErrorIs(t, err, redis.Nil, "wrong err")
	}
}
