package motivationstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/motivationerror"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/motivationstorage"
	"gotest.tools/v3/assert"
)

var (
	db      *sql.DB
	storage *motivationstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	database, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	db = database
	storage = motivationstorage.New(db)
	os.Exit(m.Run())
}

func Test_Storage_NextUserMotivation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		motivationId int

		expectedMotivationId int
		err                  error
	}{
		{
			motivationId:         1,
			expectedMotivationId: 2,
		},
		{
			motivationId:         2,
			expectedMotivationId: 3,
		},
		{
			motivationId:         3,
			expectedMotivationId: 4,
		},
		{
			motivationId:         5,
			expectedMotivationId: 6,
		},
		{
			motivationId:         29,
			expectedMotivationId: 30,
		},
		{
			motivationId: 30,
			err:          motivationerror.ExceptionMotivationNotExist(),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		err := insertUser(db, userId, cs.motivationId)
		assert.NilError(t, err, "failed insert user")

		motivation, err := storage.NextUserMotivation(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "failed get user motivation")

		assert.Equal(t, motivation.ID, cs.expectedMotivationId, "wrong motivation id")
	}
}

func Test_Storage_UpdateUserMotivation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		motivationId int

		updatedMotivationId  int
		expectedMotivationId int
		updateMotivationErr  error
	}{
		{
			motivationId:         1,
			updatedMotivationId:  2,
			expectedMotivationId: 2,
		},
		{
			motivationId:         2,
			updatedMotivationId:  200,
			updateMotivationErr:  motivationerror.ExceptionMotivationNotExist(),
			expectedMotivationId: 2,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		err := insertUser(db, userId, cs.motivationId)
		assert.NilError(t, err, "failed insert user")

		err = storage.UpdateUserMotivation(ctx, userId, cs.updatedMotivationId)
		assert.ErrorIs(t, err, cs.updateMotivationErr)

		motivation, err := userMotivation(db, userId)
		assert.NilError(t, err, "failed get user welcome motivation")
		assert.Equal(t, cs.expectedMotivationId, motivation.ID, "wrong motivation id")
	}
}
