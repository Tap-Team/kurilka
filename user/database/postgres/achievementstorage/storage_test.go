package achievementstorage_test

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	amidsql "github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/Tap-Team/kurilka/user/database/postgres/achievementstorage"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var (
	db      *sql.DB
	storage *achievementstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	db = conn
	storage = achievementstorage.New(db)
	m.Run()
}

func insertUser(db *sql.DB, userId int64, userData usermodel.CreateUser) error {
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
		userData.Name,
		userData.CigaretteDayAmount,
		userData.CigarettePackAmount,
		userData.PackPrice,
	)
	return err
}

func insertAchievement(db *sql.DB, userId int64, achievement usermodel.Achievement, isOpened bool) error {
	time := pq.NullTime{Time: time.Now(), Valid: isOpened}
	query := fmt.Sprintf(
		`
		WITH achievement as (
			SELECT %s as id FROM %s 
			INNER JOIN %s ON %s = %s AND %s = $2
			WHERE %s = $3
			GROUP BY %s
		)
		INSERT INTO %s (%s,%s,%s) VALUES (
			$1,
			(SELECT id FROM achievement),
			$4
		)`,
		// select id from achievement
		sqlutils.Full(achievementsql.ID),
		achievementsql.Table,

		achievementtypesql.Table,
		sqlutils.Full(achievementsql.TypeId),
		sqlutils.Full(achievementtypesql.ID),
		// and type = $3
		sqlutils.Full(achievementtypesql.Type),
		// where level = $2
		sqlutils.Full(achievementsql.Level),
		// group by
		sqlutils.Full(achievementsql.ID),

		userachievementsql.Table,

		userachievementsql.UserId,
		userachievementsql.AchievementId,
		userachievementsql.OpenDate,
	)
	_, err := db.Exec(query, userId, achievement.Type, achievement.Level, time)
	return err
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

func (a *achievementGenerator) Achievement() (usermodel.Achievement, bool) {
	if len(a.availableTypes) == 0 {
		return usermodel.Achievement{}, false
	}

	achIndex := rand.Intn(len(a.availableTypes))
	achtype := a.availableTypes[achIndex]

	levelIndex := rand.Intn(len(a.state[achtype]))
	level := a.state[achtype][levelIndex]

	a.state[achtype] = slices.Delete(a.state[achtype], levelIndex, levelIndex+1)

	if len(a.state[achtype]) == 0 {
		a.availableTypes = slices.Delete(a.availableTypes, achIndex, achIndex+1)
	}

	achievement := usermodel.NewAсhievement(achtype, level)
	return achievement, true
}

func generateRandomAchievementList(size int) []usermodel.Achievement {

	achievements := make([]usermodel.Achievement, 0, size)

	achGen := NewAchievementGenerator()

	// Генерируем случайные достижения
	for i := 0; i < size; i++ {
		achievement, ok := achGen.Achievement()
		if !ok {
			break
		}
		achievements = append(achievements, achievement)
	}

	return achievements
}

func filterMaxLevelFromAchievementList(achievements []*usermodel.Achievement) []*usermodel.Achievement {
	maxLevelAchievements := make(map[achievementmodel.AchievementType]usermodel.Achievement)

	for i := range achievements {
		ach := achievements[i]
		currentMaxLevel := maxLevelAchievements[ach.Type].Level
		if ach.Level > currentMaxLevel {
			maxLevelAchievements[ach.Type] = *ach
		}
	}
	achievements = make([]*usermodel.Achievement, 0, len(maxLevelAchievements))
	for _, ach := range maxLevelAchievements {
		ach := ach
		achievements = append(achievements, &ach)
	}
	return achievements
}

func TestGenerateData(t *testing.T) {
	for i := 1; i < 100; i++ {
		list := generateRandomAchievementList(i)
		require.Equal(t, len(list), min(i, 50))
		listchecker := make(map[achievementmodel.AchievementType]map[int]struct{}, len(list))
		for _, ach := range list {
			if _, ok := listchecker[ach.Type]; !ok {
				listchecker[ach.Type] = make(map[int]struct{})
			}
			_, ok := listchecker[ach.Type][ach.Level]
			if ok {
				log.Println(list)
			}
			require.False(t, ok, "duplicate found, %v", ach)

			listchecker[ach.Type][ach.Level] = struct{}{}
		}
	}
}

func TestAchievementPreview(t *testing.T) {
	ctx := context.Background()

	for i := 1; i < 50; i++ {
		user := random.StructTyped[usermodel.CreateUser]()
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, user))

		generator := NewAchievementGenerator()
		list := make([]*usermodel.Achievement, 0)
		for in := 0; in < i; in++ {
			ach, ok := generator.Achievement()
			if !ok {
				break
			}
			opened := rand.Intn(3)%2 == 0
			require.NoError(t, insertAchievement(db, userId, ach, opened), "failed insert user achievement")
			if opened {
				list = append(list, &ach)
			}
		}

		achievementPreview := storage.AchievementPreview(ctx, userId)

		achievementList := filterMaxLevelFromAchievementList(list)
		sort.SliceStable(achievementList, func(i, j int) bool {
			levelCompare := cmp.Compare(achievementList[i].Level, achievementList[j].Level)
			typeCompare := cmp.Compare(achievementList[i].Type, achievementList[j].Type)
			if levelCompare == 0 {
				return typeCompare == 1
			}
			return levelCompare == 1
		})
		sort.SliceStable(achievementPreview, func(i, j int) bool {
			levelCompare := cmp.Compare(achievementPreview[i].Level, achievementPreview[j].Level)
			typeCompare := cmp.Compare(achievementPreview[i].Type, achievementPreview[j].Type)
			if levelCompare == 0 {
				return typeCompare == 1
			}
			return levelCompare == 1
		})

		ok := slices.EqualFunc(achievementPreview, achievementList, func(a1, a2 *usermodel.Achievement) bool {
			return reflect.DeepEqual(a1, a2)
		})
		require.True(t, ok, "slices not equal")
	}
}

func TestAchievementPreviewOnlyOpen(t *testing.T) {
	ctx := context.Background()

	for i := 1; i < 50; i++ {
		user := random.StructTyped[usermodel.CreateUser]()
		userId := rand.Int63()
		require.NoError(t, insertUser(db, userId, user))

		list := generateRandomAchievementList(i)
		for _, ach := range list {
			require.NoError(t, insertAchievement(db, userId, ach, false), "failed insert user achievement")
		}

		achievementPreview := storage.AchievementPreview(ctx, userId)

		require.Equal(t, 0, len(achievementPreview), "achievement preview returned not opened achievement")
	}
}
