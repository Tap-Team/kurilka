package welcomemotivationstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/welcomemotivationerror"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/welcomemotivationstorage"
	"gotest.tools/v3/assert"
)

var (
	db      *sql.DB
	storage *welcomemotivationstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	database, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	db = database
	storage = welcomemotivationstorage.New(db)
	os.Exit(m.Run())
}

func Test_Storage_NextUserWelcomeMotivation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		welcomeMotivationId int

		expectedMotivationId int
	}{
		{
			welcomeMotivationId:  1,
			expectedMotivationId: 2,
		},
		{
			welcomeMotivationId:  2,
			expectedMotivationId: 3,
		},
		{
			welcomeMotivationId:  3,
			expectedMotivationId: 4,
		},
		{
			welcomeMotivationId:  5,
			expectedMotivationId: 6,
		},
		{
			welcomeMotivationId:  20,
			expectedMotivationId: 1,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		err := insertUser(db, userId, cs.welcomeMotivationId)
		assert.NilError(t, err, "failed insert user")

		motivation, err := storage.NextUserWelcomeMotivation(ctx, userId)
		assert.NilError(t, err, "failed get user motivation")

		assert.Equal(t, motivation.ID, cs.expectedMotivationId, "wrong motivation id")
	}
}

func Test_Storage_UpdateUserWelcomeMotivation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		welcomeMotivationId int

		updatedMotivationId         int
		expectedWelcomeMotivationId int
		updateMotivationErr         error
	}{
		{
			welcomeMotivationId:         1,
			updatedMotivationId:         2,
			expectedWelcomeMotivationId: 2,
		},
		{
			welcomeMotivationId:         2,
			updatedMotivationId:         200,
			updateMotivationErr:         welcomemotivationerror.ExceptionMotivationNotExist(),
			expectedWelcomeMotivationId: 2,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		err := insertUser(db, userId, cs.welcomeMotivationId)
		assert.NilError(t, err, "failed insert user")

		err = storage.UpdateUserWelcomeMotivation(ctx, userId, cs.updatedMotivationId)
		assert.ErrorIs(t, err, cs.updateMotivationErr)

		motivation, err := userWelcomeMotivation(db, userId)
		assert.NilError(t, err, "failed get user welcome motivation")
		assert.Equal(t, cs.expectedWelcomeMotivationId, motivation.ID, "wrong motivation id")
	}
}
