package achievementstorage_test

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/achievementstorage"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

var (
	db      *sql.DB
	storage *achievementstorage.Storage
)

func TestMain(m *testing.M) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	d, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create new postgres container")
	}
	defer term(ctx)
	db = d
	storage = achievementstorage.New(db)
	os.Exit(m.Run())
}

func TestUserAchievements(t *testing.T) {
	ctx := context.Background()
	for i := 1; i < 50; i++ {
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievementList := generateRandomAchievementList(i)
		achievements := make(map[int]map[achievementmodel.AchievementType]struct{})
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
			if _, ok := achievements[ach.level]; !ok {
				achievements[ach.level] = make(map[achievementmodel.AchievementType]struct{})
			}
			achievements[ach.level][ach.achType] = struct{}{}
		}
		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		for _, ach := range userAchievements {
			_, ok := achievements[ach.Level][ach.Type]

			assert.Equal(t, false, ach.Opened(), "ach is open")
			assert.Equal(t, !ok, ach.Shown, "ach is shown")
			assert.Equal(t, ok, ach.Reached(), "ach is reached")
			assert.Equal(t, false, len(ach.Description) == 0, "description is zero")
		}

	}
}

func TestUserAchievementsNotAchieveList(t *testing.T) {
	ctx := context.Background()
	userId := rand.Int63()
	err := insertUser(db, userId, random.StructTyped[usermodel.CreateUser]())
	assert.NilError(t, err, "failed insert user")

	achievements, err := storage.UserAchievements(ctx, userId)
	assert.NilError(t, err, "failed get user achievements")

	assert.Equal(t, 50, len(achievements), "wrong size of achievements")

	for _, ach := range achievements {
		assert.Equal(t, false, ach.Opened(), "ach opened")
		assert.Equal(t, false, ach.Reached(), "ach reached")
		assert.Equal(t, true, ach.Shown, "ach shown")
	}
}
