package achievementstorage_test

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"slices"

	"github.com/Tap-Team/kurilka/achievements/database/postgres/achievementstorage"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
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

func insertUser(db *sql.DB, userId int64, createUser usermodel.CreateUser) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1,$2,$3,$4,$5)`,
		// insert into users
		usersql.Table,
		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
	)
	_, err := db.Exec(query,
		userId,
		createUser.Name,
		createUser.CigaretteDayAmount,
		createUser.CigarettePackAmount,
		createUser.PackPrice,
	)
	return err
}

type insertAchievement struct {
	level   int
	achType achievementmodel.AchievementType
}

type achievementGenerator struct {
	availableTypes []achievementmodel.AchievementType
	state          map[achievementmodel.AchievementType][]int
}

func NewAchievementGenerator() *achievementGenerator {
	achtypeList := []achievementmodel.AchievementType{
		achievementmodel.DURATION,
		achievementmodel.CIGARETTE,
		achievementmodel.HEALTH,
		achievementmodel.WELL_BEING,
		achievementmodel.SAVING,
	}
	state := make(map[achievementmodel.AchievementType][]int, len(achtypeList))
	for _, tp := range achtypeList {
		state[tp] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	}
	return &achievementGenerator{
		availableTypes: achtypeList,
		state:          state,
	}
}

func (a *achievementGenerator) Achievement() (insertAchievement, bool) {
	if len(a.availableTypes) == 0 {
		return insertAchievement{}, false
	}

	achIndex := rand.Intn(len(a.availableTypes))
	achtype := a.availableTypes[achIndex]

	levelIndex := rand.Intn(len(a.state[achtype]))
	level := a.state[achtype][levelIndex]

	a.state[achtype] = slices.Delete(a.state[achtype], levelIndex, levelIndex+1)

	if len(a.state[achtype]) == 0 {
		a.availableTypes = slices.Delete(a.availableTypes, achIndex, achIndex+1)
	}

	achievement := insertAchievement{level: level, achType: achtype}
	return achievement, true
}

func generateRandomAchievementList(size int) []*insertAchievement {

	achievements := make([]*insertAchievement, 0, size)

	achGen := NewAchievementGenerator()

	// Генерируем случайные достижения
	for i := 0; i < size; i++ {
		achievement, ok := achGen.Achievement()
		if !ok {
			break
		}
		achievements = append(achievements, &achievement)
	}

	return achievements
}

func userAchieve(db *sql.DB, userId int64, ach *insertAchievement) error {
	query := fmt.Sprintf(
		`
		WITH achievement_select as (
			SELECT %s as id FROM %s 
			INNER JOIN %s ON %s = %s 
			WHERE %s = $2 AND %s = $3
			GROUP BY %s
		)
		INSERT INTO %s (%s,%s) VALUES ($1, (SELECT id FROM achievement_select)) 
		`,
		sqlutils.Full(achievementsql.ID),
		achievementsql.Table,

		// inner join achievement type
		achievementtypesql.Table,
		sqlutils.Full(achievementsql.TypeId),
		sqlutils.Full(achievementtypesql.ID),

		// where level and type eq $2 and $3
		sqlutils.Full(achievementtypesql.Type),
		sqlutils.Full(achievementsql.Level),

		// group by
		sqlutils.Full(achievementsql.ID),

		userachievementsql.Table,
		userachievementsql.UserId,
		userachievementsql.AchievementId,
	)
	_, err := db.Exec(query, userId, ach.achType, ach.level)
	return err
}

func compareInsertWithAchievement(iach *insertAchievement, ach *achievementmodel.Achievement) (string, bool) {
	var fields strings.Builder
	if iach.achType != ach.Type {
		fields.WriteString("Type ")
	}
	if iach.level != ach.Level {
		fields.WriteString("Level ")
	}
	return fields.String(), fields.Len() == 0
}

func sortAchievements(achList []*achievementmodel.Achievement) {
	sort.SliceStable(achList, func(i, j int) bool {
		levelCompare := cmp.Compare(achList[i].Level, achList[j].Level)
		typeCompare := cmp.Compare(achList[i].Type, achList[j].Type)
		if levelCompare == 0 {
			return typeCompare == 1
		}
		return levelCompare == 1
	})
}

func sortInsertAchievements(iachList []*insertAchievement) {
	sort.SliceStable(iachList, func(i, j int) bool {
		levelCompare := cmp.Compare(iachList[i].level, iachList[j].level)
		typeCompare := cmp.Compare(iachList[i].achType, iachList[j].achType)
		if levelCompare == 0 {
			return typeCompare == 1
		}
		return levelCompare == 1
	})
}

func TestUserAchievements(t *testing.T) {
	ctx := context.Background()
	for i := 1; i < 50; i++ {
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievementList := generateRandomAchievementList(i)
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
		}
		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")
		sortAchievements(userAchievements)
		sortInsertAchievements(achievementList)
		ok := slices.EqualFunc(achievementList, userAchievements, func(ia *insertAchievement, a *achievementmodel.Achievement) bool {
			assert.Assert(t, a.ID != 0, "zero achievement id")
			assert.Equal(t, a.Shown, false, "achievement shown wrong value")
			assert.Equal(t, a.OpenDate.IsZero(), true, "achievement open date not null")
			fields, ok := compareInsertWithAchievement(ia, a)
			if ok {
				return ok
			}
			t.Logf("fields not equal, %s", fields)
			return false
		})
		require.True(t, ok, "user achievements not equal")
	}
}

func TestMarkShown(t *testing.T) {
	ctx := context.Background()
	for i := 1; i < 50; i++ {
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievementList := generateRandomAchievementList(i)
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
		}
		err := storage.MarkShown(ctx, userId)
		require.NoError(t, err, "failed mark shown")

		userAchievements, err := storage.UserAchievements(ctx, userId)
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

		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		openTime := amidtime.Timestamp{Time: time.Now()}
		achIndex := rand.Intn(len(userAchievements))
		achID := userAchievements[achIndex].ID
		err = storage.OpenSingle(ctx, userId, model.NewOpenAchievement(achID, openTime.Time))
		require.NoError(t, err, "failed open single achievement")

		for _, ach := range userAchievements {
			achOpenTime, err := achievementOpenTime(db, userId, ach.ID)
			require.NoError(t, err, "failed get achievement open time")
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

func TestOpenType_BasicCase(t *testing.T) {
	ctx := context.Background()

	achtypes := []achievementmodel.AchievementType{
		achievementmodel.DURATION,
		achievementmodel.CIGARETTE,
		achievementmodel.HEALTH,
		achievementmodel.WELL_BEING,
		achievementmodel.SAVING,
	}

	for i := 1; i < 50; i++ {
		achtype := achtypes[rand.Intn(len(achtypes))]
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievementList := generateRandomAchievementList(i)
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
		}

		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		openTime := amidtime.Timestamp{Time: time.Now()}

		achIds, err := storage.OpenType(ctx, userId, model.NewAchievementType(achtype, openTime.Time))
		require.NoError(t, err, "failed open by type")

		for _, achId := range achIds {
			achOpenTime, err := achievementOpenTime(db, userId, achId)
			require.NoError(t, err, "failed get achievement open time")
			assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time from achIds not equal")
		}
		var i int
		for _, ach := range userAchievements {
			achOpenTime, err := achievementOpenTime(db, userId, ach.ID)
			require.NoError(t, err, "failed get achievement open time")
			if ach.Type == achtype {
				i++
				assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time not equal")
			} else {
				assert.Equal(t, achOpenTime.IsZero(), true, "achievement opened")
			}
		}

		require.Equal(t, len(achIds), i)
	}
}

func TestOpenAll(t *testing.T) {
	ctx := context.Background()
	for i := 1; i < 50; i++ {
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, random.StructTyped[usermodel.CreateUser]()))

		achievementList := generateRandomAchievementList(i)
		for _, ach := range achievementList {
			require.NoError(t, userAchieve(db, userId, ach), "failed inser achievement")
		}

		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		openTime := amidtime.Timestamp{Time: time.Now()}

		achIds, err := storage.OpenAll(ctx, userId, openTime)
		require.NoError(t, err, "failed open all")

		for _, achId := range achIds {
			achOpenTime, err := achievementOpenTime(db, userId, achId)
			require.NoError(t, err, "failed get achievement open time")
			assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time from achIds not equal")
		}
		for _, ach := range userAchievements {
			achOpenTime, err := achievementOpenTime(db, userId, ach.ID)
			require.NoError(t, err, "failed get achievement open time")
			assert.Equal(t, openTime.Unix(), achOpenTime.Unix(), "open time not equal")
		}

		require.Equal(t, len(achIds), len(userAchievements))
	}
}
