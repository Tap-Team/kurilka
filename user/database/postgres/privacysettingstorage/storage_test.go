package privacysettingstorage_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"sort"
	"testing"

	amidsql "github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/stretchr/testify/require"

	"github.com/Tap-Team/kurilka/internal/errorutils/userprivacysettingerror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/user/database/postgres/privacysettingstorage"
)

var (
	db      *sql.DB
	storage *privacysettingstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container, %s", err)
	}
	defer term(ctx)
	db = conn
	storage = privacysettingstorage.New(db)
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
	var err error
	ctx := context.Background()

	user := random.StructTyped[usermodel.CreateUser]()
	userId := rand.Int63()

	err = insertUser(db, userId, user)
	require.NoError(t, err, "failed insert user")

	settings, err := storage.UserPrivacySettings(ctx, userId)
	require.NoError(t, err, "failed get user privacy settings")

	require.Equal(t, 0, len(settings), "wrong privacy settings length")

	const settingsSize = 4
	psets := getRandomPrivacySettingsList(settingsSize)

	for _, stn := range psets {
		err = storage.AddUserPrivacySetting(ctx, userId, stn)
		require.NoError(t, err, "failed add user privacy setting")
	}

	settings, err = storage.UserPrivacySettings(ctx, userId)
	require.NoError(t, err, "failed get user privacy settings")

	sort.Slice(settings, func(i, j int) bool {
		return settings[i] > settings[j]
	})
	sort.Slice(psets, func(i, j int) bool {
		return psets[i] > psets[j]
	})
	res := slices.Compare(settings, psets)
	require.Equal(t, 0, res, "settings not equal")

	psets = make([]usermodel.PrivacySetting, 0)
	for _, stn := range settings {

		if rand.Intn(3)%2 == 0 {
			psets = append(psets, stn)
		} else {
			err := storage.RemoveUserPrivacySetting(ctx, userId, stn)
			require.NoError(t, err, "failed remove user privacy setting")
		}
	}

	settings, err = storage.UserPrivacySettings(ctx, userId)
	require.NoError(t, err, "failed get user privacy settings")

	sort.Slice(settings, func(i, j int) bool {
		return settings[i] > settings[j]
	})
	sort.Slice(psets, func(i, j int) bool {
		return psets[i] > psets[j]
	})
	res = slices.Compare(settings, psets)
	require.Equal(t, 0, res, "settings not equal after delete")
}

func TestAddUserPrivacySetting(t *testing.T) {
	ctx := context.Background()

	var id1, id2 int64 = rand.Int63(), rand.Int63()
	require.NoError(t, insertUser(db, id1, random.StructTyped[usermodel.CreateUser]()), "failed insert user 1")
	require.NoError(t, insertUser(db, id2, random.StructTyped[usermodel.CreateUser]()), "failed insert user 2")

	cases := []struct {
		userId   int64
		settings []usermodel.PrivacySetting
		err      error
	}{
		{
			userId:   id1,
			settings: getRandomPrivacySettingsList(9),
		},

		{
			userId:   id1,
			settings: getRandomPrivacySettingsList(9),
			err:      userprivacysettingerror.ExceptionUserPrivacySettingExists(),
		},
		{
			userId:   id2,
			settings: []usermodel.PrivacySetting{"asdfasdf"},
			err:      userprivacysettingerror.ExceptionUserPrivacySettingNotFound(),
		},
	}

	for _, cs := range cases {
		uset := make([]usermodel.PrivacySetting, 0, len(cs.settings))
		sets, err := storage.UserPrivacySettings(ctx, cs.userId)
		require.NoError(t, err, "failed get user privacy settings")
		uset = append(uset, sets...)
		for _, stn := range cs.settings {
			err := storage.AddUserPrivacySetting(ctx, cs.userId, stn)
			require.ErrorIs(t, err, cs.err)
			if err == nil {
				uset = append(uset, stn)
			}
			sets, err := storage.UserPrivacySettings(ctx, cs.userId)
			require.NoError(t, err, "failed get user privacy settings")

			sort.Slice(sets, func(i, j int) bool { return sets[i] > sets[j] })
			sort.Slice(uset, func(i, j int) bool { return uset[i] > uset[j] })
			res := slices.Compare(sets, uset)

			require.Equal(t, 0, res, "settings not equal")
		}
	}
}

func TestRemoveUserPrivacySetting(t *testing.T) {
	ctx := context.Background()

	var id1, id2 int64 = rand.Int63(), rand.Int63()
	require.NoError(t, insertUser(db, id1, random.StructTyped[usermodel.CreateUser]()), "failed insert user 1")
	require.NoError(t, insertUser(db, id2, random.StructTyped[usermodel.CreateUser]()), "failed insert user 2")

	err := storage.RemoveUserPrivacySetting(ctx, id1, usermodel.ACHIEVEMENTS_CIGARETTE)
	require.ErrorIs(t, err, userprivacysettingerror.ExceptionUserPrivacySettingNotFound(), "wrong user privacy err")

	cases := []struct {
		userId   int64
		settings []usermodel.PrivacySetting
		err      error
	}{
		{
			userId:   id1,
			settings: getRandomPrivacySettingsList(9),
		},
		{
			userId:   id2,
			settings: getRandomPrivacySettingsList(9),
		},
	}

	for _, cs := range cases {
		uset := make([]usermodel.PrivacySetting, 0, len(cs.settings))
		for _, stn := range cs.settings {
			require.NoError(t, storage.AddUserPrivacySetting(ctx, cs.userId, stn), "failed add user privacy")
			uset = append(uset, stn)
		}
		for i, stn := range cs.settings {
			err := storage.RemoveUserPrivacySetting(ctx, cs.userId, stn)
			require.ErrorIs(t, err, cs.err, "wrong remove err, %d", i)

			if err == nil {
				uset = slices.DeleteFunc(uset, func(ps usermodel.PrivacySetting) bool { return ps == stn })
			}

			sets, err := storage.UserPrivacySettings(ctx, cs.userId)
			require.NoError(t, err, "failed get user privacy settings")

			sort.Slice(sets, func(i, j int) bool { return sets[i] > sets[j] })
			sort.Slice(uset, func(i, j int) bool { return uset[i] > uset[j] })
			res := slices.Compare(sets, uset)

			require.Equal(t, 0, res, "settings not equal")
		}
	}

}
