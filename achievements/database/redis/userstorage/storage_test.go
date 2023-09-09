package userstorage_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/database/redis/userstorage"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/redis/go-redis/v9"
	"gotest.tools/v3/assert"
)

var (
	rc      *redis.Client
	storage *userstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	redisClient, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed start redis container, %s", err)
	}
	defer term(ctx)
	rc = redisClient
	storage = userstorage.New(rc)
	os.Exit(m.Run())
}

func TestUser(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		user *usermodel.UserData
		data *model.UserData
		err  error
	}{
		{
			user: usermodel.NewUserData(
				"dima",
				10,
				10,
				10,
				"",
				"",
				time.Now(),
				usermodel.LevelInfo{},
				[]usermodel.Trigger{},
			),
			data: model.NewUserData(10, 10, 10, time.Now()),
		},
		{
			err: usererror.ExceptionUserNotFound(),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		saveUser(ctx, rc, userId, cs.user)

		userData, err := storage.User(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "error not equal")

		if cs.data == nil {
			assert.Equal(t, true, userData == nil, "user data not equal")
			return
		}
		if userData == nil {
			t.Fatal("nil user data from storage")
		}
		assert.Equal(t, userData.AbstinenceTime.Unix(), cs.data.AbstinenceTime.Unix(), "time not equal")
		assert.Equal(t, userData.CigaretteDayAmount, cs.data.CigaretteDayAmount, "abstinence time not equal")
		assert.Equal(t, userData.CigarettePackAmount, cs.data.CigarettePackAmount, "pack price not equal")
		assert.Equal(t, userData.PackPrice, cs.data.PackPrice, "pack price not equal")
	}
}
