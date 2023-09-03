package triggerstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/user/database/postgres/triggerstorage"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var (
	db      *sql.DB
	storage *triggerstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	database, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container")
	}
	defer term(ctx)

	db = database
	storage = triggerstorage.New(db)
	os.Exit(m.Run())
}

func sortUserTriggers(triggers []usermodel.Trigger) {
	sort.Slice(triggers, func(i, j int) bool { return triggers[i] > triggers[j] })

}

func TestRemove(t *testing.T) {
	ctx := context.Background()
	const amount = 100
	var wg sync.WaitGroup
	wg.Add(amount)
	for i := 0; i < amount; i++ {
		go func() {
			defer wg.Done()
			triggers := []usermodel.Trigger{
				usermodel.SUPPORT_CIGGARETTE,
				usermodel.THANK_YOU,
				usermodel.SUPPORT_HEALTH,
				usermodel.SUPPORT_TRIAL,
			}
			sortUserTriggers(triggers)
			userId := rand.Int63()
			err := insertUserWithAllTriggers(db, userId)
			require.NoError(t, err, "failed insert user with triggers")

			usertriggers, err := userTriggers(db, userId)
			require.NoError(t, err, "failed get user triggers")
			sortUserTriggers(usertriggers)

			ok := slices.Equal(usertriggers, triggers)
			require.True(t, ok, "slices not equal")

			for i := len(triggers) - 1; i != 0; i-- {
				trigger := triggers[i]
				err := storage.Remove(ctx, userId, trigger)
				require.NoError(t, err, "failed remove trigger")
				triggers = slices.Delete(triggers, i, i+1)

				usertriggers, err := userTriggers(db, userId)
				require.NoError(t, err, "failed get user triggers")
				sortUserTriggers(triggers)
				sortUserTriggers(usertriggers)

				ok := slices.Equal(usertriggers, triggers)
				require.True(t, ok, "slices not equal")
			}
		}()
	}
	wg.Wait()

}
