package achievementstorage_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/userachievementerror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/Tap-Team/kurilka/workers/userworker/database/redis/achievementstorage"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	rc      *redis.Client
	storage *achievementstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	client, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed create redis container, %s", err)
	}
	defer term(ctx)
	rc = client
	storage = achievementstorage.New(rc, 0)
	os.Exit(m.Run())
}

func generateRandomAchievementList(size int) []*achievementmodel.Achievement {
	types := []achievementmodel.AchievementType{
		achievementmodel.DURATION,
		achievementmodel.CIGARETTE,
		achievementmodel.HEALTH,
		achievementmodel.WELL_BEING,
		achievementmodel.SAVING,
	}

	// Создаем слайс для хранения случайных достижений
	achievements := make([]*achievementmodel.Achievement, 0, size)

	// Генерируем случайные достижения
	for i := 0; i < size; i++ {
		index := rand.Intn(len(types))
		level := rand.Intn(11)
		exp := rand.Intn(1001)
		openDate := amidtime.Timestamp{Time: time.Now()}
		shown := rand.Intn(3)%2 == 0

		achievement := achievementmodel.Achievement{
			ID:       int64(i + 1),
			Type:     types[index],
			Exp:      exp,
			Level:    level,
			OpenDate: openDate,
			Shown:    shown,
		}

		achievements = append(achievements, &achievement)
	}

	return achievements
}

func compareAchievements(a1, a2 *achievementmodel.Achievement) (string, bool) {
	fields := strings.Builder{}
	if a1.ID != a2.ID {
		fields.WriteString("ID ")
	}
	if a1.Type != a2.Type {
		fields.WriteString("Type ")
	}
	if a1.Exp != a2.Exp {
		fields.WriteString("Exp ")
	}
	if a1.Level != a2.Level {
		fields.WriteString("Level ")
	}
	if a1.OpenDate.Unix() != a2.OpenDate.Unix() {
		fields.WriteString("OpenDate ")
	}
	if a1.Shown != a2.Shown {
		fields.WriteString("Shown ")
	}
	return fields.String(), fields.Len() == 0
}

func compareAchievementList(a1, a2 []*achievementmodel.Achievement) bool {
	return slices.EqualFunc(a1, a2, func(a1, a2 *achievementmodel.Achievement) bool {
		fields, ok := compareAchievements(a1, a2)
		if !ok {
			log.Printf("fields not equal, %s", fields)
		}
		return ok
	})
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()

	for i := 1; i < 50; i++ {
		userId := rand.Int63()

		genAchievements := generateRandomAchievementList(i)
		err := storage.SaveUserAchievements(ctx, userId, genAchievements)
		require.NoError(t, err, "failed save user achievements")

		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		ok := compareAchievementList(genAchievements, userAchievements)
		require.True(t, ok, "list not equal")

		err = storage.RemoveUserAchievements(ctx, userId)
		require.NoError(t, err, "failed remove user achievements")

		_, err = storage.UserAchievements(ctx, userId)
		require.ErrorIs(t, err, userachievementerror.ExceptionAchievementNotFound(), "wrong not found error")
	}
}

func TestSave(t *testing.T) {
	ctx := context.Background()
	for i := 1; i < 50; i++ {
		userId := rand.Int63()

		genAchievements := generateRandomAchievementList(i)
		err := storage.SaveUserAchievements(ctx, userId, genAchievements)
		require.NoError(t, err, "failed save user achievements")

		userAchievements, err := storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements")

		ok := compareAchievementList(genAchievements, userAchievements)
		require.True(t, ok, "list not equal")

		genAchievements = generateRandomAchievementList(i)
		err = storage.SaveUserAchievements(ctx, userId, genAchievements)
		require.NoError(t, err, "failed update user achievements")

		userAchievements, err = storage.UserAchievements(ctx, userId)
		require.NoError(t, err, "failed get user achievements after update")

		ok = compareAchievementList(genAchievements, userAchievements)
		require.True(t, ok, "list not equal after update")
	}
}
