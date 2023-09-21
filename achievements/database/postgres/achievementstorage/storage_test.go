package achievementstorage_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/database/postgres/achievementstorage"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/random"
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

func TestMarkShown(t *testing.T) {
	ctx := context.Background()
	for i := 1; i < 51; i++ {
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievements := make(map[int]map[achievementmodel.AchievementType]struct{})
		achievementList := generateRandomAchievementList(i)
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
			if _, ok := achievements[ach.level]; !ok {
				achievements[ach.level] = make(map[achievementmodel.AchievementType]struct{})
			}
			achievements[ach.level][ach.achType] = struct{}{}
		}
		userAchievements, err := storage.UserAchievements(ctx, userId)
		for _, ach := range userAchievements {
			_, ok := achievements[ach.Level][ach.Type]
			assert.Equal(t, !ok, ach.Shown)
		}

		err = storage.MarkShown(ctx, userId)
		require.NoError(t, err, "failed mark shown")

		userAchievements, err = storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")
		for _, ach := range userAchievements {
			assert.Equal(t, ach.Shown, true, "achievement not shown")
		}
	}
}

func achievementOpenTime(db *sql.DB, userId int64, achievementId int64) (amidtime.Timestamp, error) {
	var openTime amidtime.Timestamp
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = $1 AND %s = $2`,
		userachievementsql.OpenDate,
		userachievementsql.Table,
		userachievementsql.UserId,
		userachievementsql.AchievementId,
	)
	err := db.QueryRow(query, userId, achievementId).Scan(&openTime)
	return openTime, err
}

func TestOpenSingle_BasicCase(t *testing.T) {
	ctx := context.Background()

	for i := 1; i < 50; i++ {
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievementList := generateRandomAchievementList(i)
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
		}

		userAchievements, err := userReachedAchievements(db, userId)
		require.NoError(t, err, "failed get user reached achievements")

		openTime := amidtime.Timestamp{Time: time.Now()}
		achIndex := rand.Intn(len(userAchievements))
		achID := userAchievements[achIndex].ID
		err = storage.OpenSingle(ctx, userId, model.NewOpenAchievement(achID, openTime.Time))
		require.NoError(t, err, "failed open single achievement")

		userAchievements, err = storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		for _, ach := range userAchievements {
			achOpenTime, _ := achievementOpenTime(db, userId, ach.ID)
			if ach.ID == achID {
				assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time not equal")
			} else {
				assert.Equal(t, achOpenTime.IsZero(), true, "achievement is openn")
			}
		}

	}
}

func TestOpenSingle_UserNotExistsCase(t *testing.T) {
	ctx := context.Background()

	for i := 1; i < 50; i++ {
		userId := rand.Int63()
		openTime := amidtime.Timestamp{Time: time.Now()}
		achID := i
		err := storage.OpenSingle(ctx, userId, model.NewOpenAchievement(int64(achID), openTime.Time))
		require.ErrorIs(t, err, usererror.ExceptionUserNotFound(), "wrong error")
	}
}

// func TestOpenType_BasicCase(t *testing.T) {
// 	ctx := context.Background()

// 	achtypes := []achievementmodel.AchievementType{
// 		achievementmodel.DURATION,
// 		achievementmodel.CIGARETTE,
// 		achievementmodel.HEALTH,
// 		achievementmodel.WELL_BEING,
// 		achievementmodel.SAVING,
// 	}

// 	for i := 1; i < 50; i++ {
// 		achtype := achtypes[rand.Intn(len(achtypes))]
// 		userId := rand.Int63()
// 		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

// 		achievementList := generateRandomAchievementList(i)
// 		for _, ach := range achievementList {
// 			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
// 		}

// 		userAchievements, err := storage.UserAchievements(ctx, userId)
// 		require.NoError(t, err, "failed get user achievements")

// 		openTime := amidtime.Timestamp{Time: time.Now()}

// 		achIds, err := storage.OpenType(ctx, userId, model.NewAchievementType(achtype, openTime.Time))
// 		require.NoError(t, err, "failed open by type")

// 		for _, achId := range achIds {
// 			achOpenTime, err := achievementOpenTime(db, userId, achId)
// 			require.NoError(t, err, "failed get achievement open time")
// 			assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time from achIds not equal")
// 		}
// 		var i int
// 		for _, ach := range userAchievements {
// 			achOpenTime, err := achievementOpenTime(db, userId, ach.ID)
// 			require.NoError(t, err, "failed get achievement open time")
// 			if ach.Type == achtype {
// 				i++
// 				assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time not equal")
// 			} else {
// 				assert.Equal(t, achOpenTime.IsZero(), true, "achievement opened")
// 			}
// 		}

// 		require.Equal(t, len(achIds), i)
// 	}
// }

// func TestOpenAll(t *testing.T) {
// 	ctx := context.Background()
// 	for i := 1; i < 50; i++ {
// 		userId := rand.Int63()
// 		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

// 		achievementList := generateRandomAchievementList(i)
// 		for _, ach := range achievementList {
// 			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
// 		}

// 		userAchievements, err := storage.UserAchievements(ctx, userId)
// 		require.NoError(t, err, "failed get user achievements")

// 		openTime := amidtime.Timestamp{Time: time.Now()}

// 		achIds, err := storage.OpenAll(ctx, userId, openTime)
// 		require.NoError(t, err, "failed open all")

// 		for _, achId := range achIds {
// 			achOpenTime, err := achievementOpenTime(db, userId, achId)
// 			require.NoError(t, err, "failed get achievement open time")
// 			assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time from achIds not equal")
// 		}
// 		for _, ach := range userAchievements {
// 			achOpenTime, err := achievementOpenTime(db, userId, ach.ID)
// 			require.NoError(t, err, "failed get achievement open time")
// 			assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time not equal")
// 		}

// 		require.Equal(t, len(achIds), len(userAchievements))
// 	}
// }
