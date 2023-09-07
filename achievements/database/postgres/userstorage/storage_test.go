package userstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/database/postgres/userstorage"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"gotest.tools/v3/assert"
)

var (
	db      *sql.DB
	storage *userstorage.Storage
)

func TestMain(m *testing.M) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	database, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	db = database
	storage = userstorage.New(db)
	os.Exit(m.Run())
}

func TestUser(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		user *usermodel.CreateUser
		data *model.UserData
		err  error
	}{
		{
			user: usermodel.NewCreateUser(
				"dima",
				10,
				10,
				10,
			),
			data: model.NewUserData(10, 10, 10, time.Now()),
		},
		{
			err: usererror.ExceptionUserNotFound(),
		},
	}

	for _, cs := range cases {
		now := time.Now()
		userId := rand.Int63()
		err := insertUser(db, userId, cs.user, now)
		assert.NilError(t, err, "failed insert user")
		userData, err := storage.User(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "error not equal")

		if cs.data == nil {
			assert.Equal(t, true, userData == nil, "user data not equal")
			return
		}
		if userData == nil {
			t.Fatal("nil user data from storage")
		}
		assert.Equal(t, userData.AbstinenceTime.Unix(), now.Unix(), "time not equal")
		assert.Equal(t, userData.CigaretteDayAmount, cs.data.CigaretteDayAmount, "abstinence time not equal")
		assert.Equal(t, userData.CigarettePackAmount, cs.data.CigarettePackAmount, "pack price not equal")
		assert.Equal(t, userData.PackPrice, cs.data.PackPrice, "pack price not equal")
	}
}
