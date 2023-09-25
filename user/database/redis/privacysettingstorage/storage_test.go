package privacysettingstorage_test

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"slices"

	"github.com/Tap-Team/kurilka/internal/errorutils/userprivacysettingerror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/Tap-Team/kurilka/user/database/redis/privacysettingstorage"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	storage     *privacysettingstorage.Storage
	redisClient *redis.Client
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	rc, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed create redis container, %s", err)
	}
	defer term(ctx)
	redisClient = rc
	storage = privacysettingstorage.New(redisClient, 0)

	m.Run()
}

func getRandomPrivacySettingsList(size int) []usermodel.PrivacySetting {
	settings := []usermodel.PrivacySetting{
		usermodel.STATISTICS_MONEY,
		usermodel.STATISTICS_CIGARETTE,
		usermodel.STATISTICS_LIFE,
		usermodel.STATISTICS_TIME,
		usermodel.ACHIEVEMENTS_DURATION,
		usermodel.ACHIEVEMENTS_HEALTH,
		usermodel.ACHIEVEMENTS_WELL_BEING,
		usermodel.ACHIEVEMENTS_SAVING,
		usermodel.ACHIEVEMENTS_CIGARETTE,
	}

	// Проверяем, что размер не превышает количество доступных настроек
	if size > len(settings) {
		size = len(settings)
	}

	// Создаем слайс для хранения случайных настроек
	randomSettings := make([]usermodel.PrivacySetting, 0, size)

	// Выбираем случайные настройки
	for len(randomSettings) < size {
		index := rand.Intn(len(settings))

		// Проверяем, что настройка еще не была выбрана
		found := false
		for _, setting := range randomSettings {
			if setting == settings[index] {
				found = true
				break
			}
		}

		// Если настройка еще не была выбрана, добавляем ее в список
		if !found {
			randomSettings = append(randomSettings, settings[index])
		}
	}

	return randomSettings
}

func TestCRUD(t *testing.T) {
	const count = 100
	var err error
	ctx := context.Background()

	userIds := make([]int64, 0)

	for i := 0; i < count; i++ {
		userId := rand.Int63()

		sets, err := storage.UserPrivacySettings(ctx, userId)
		require.Equal(t, 0, len(sets), "wrong settings length")

		privacySettings := getRandomPrivacySettingsList(9)
		err = storage.SaveUserPrivacySettings(ctx, userId, privacySettings)
		require.NoError(t, err, "failed save privacy settings")

		sets, err = storage.UserPrivacySettings(ctx, userId)
		require.NoError(t, err, "failed get user privacy settings")

		equal := slices.Equal(privacySettings, sets)
		require.True(t, equal, "settings not equal")

		userIds = append(userIds, userId)
	}

	for _, userId := range userIds {
		err = storage.RemoveUserPrivacySettings(ctx, userId)
		require.NoError(t, err, "failed remove user privacy settings")

		sets, err := storage.UserPrivacySettings(ctx, userId)
		require.ErrorIs(t, err, userprivacysettingerror.ExceptionUserPrivacySettingNotFound(), "failed get user privacy settings")
		require.Equal(t, 0, len(sets), "wrong settings length")
	}

}

func TestSaveUserPrivacySettings(t *testing.T) {
	const count = 100
	var err error

	ctx := context.Background()

	for i := 0; i < count; i++ {
		userId := rand.Int63()

		// save privacy settings in first time
		{
			sets, err := storage.UserPrivacySettings(ctx, userId)
			require.Equal(t, 0, len(sets), "wrong settings length")

			privacySettings := getRandomPrivacySettingsList(4)
			err = storage.SaveUserPrivacySettings(ctx, userId, privacySettings)
			require.NoError(t, err, "failed save privacy settings")

			sets, err = storage.UserPrivacySettings(ctx, userId)
			require.NoError(t, err, "failed get user privacy settings")

			equal := slices.Equal(privacySettings, sets)
			require.True(t, equal, "settings not equal")
		}

		// update privacy settings
		{
			privacySettings := getRandomPrivacySettingsList(4)
			err = storage.SaveUserPrivacySettings(ctx, userId, privacySettings)
			require.NoError(t, err, "failed save privacy settings")

			sets, err := storage.UserPrivacySettings(ctx, userId)
			require.NoError(t, err, "failed get user privacy settings")

			equal := slices.Equal(privacySettings, sets)
			require.True(t, equal, "settings not equal")
		}
	}
}
