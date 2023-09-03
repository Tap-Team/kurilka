package resetrecoveruserstorage_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/privacysettingsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/subscriptiontypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/triggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userprivacysettingsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersubscriptionsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usertriggersql"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

type userTransaction struct {
	*sql.Tx
	userId int64
	user   *usermodel.CreateUser
}

func (u userTransaction) Insert() error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s,%s,%s,%s,%s) VALUES ($1,$2,$3,$4,$5)`,
		usersql.Table,

		usersql.ID,
		usersql.Name,
		usersql.PackPrice,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
	)
	_, err := u.Exec(
		query,
		u.userId,
		u.user.Name,
		u.user.PackPrice,
		u.user.CigaretteDayAmount,
		u.user.CigarettePackAmount,
	)
	return err
}

func (u userTransaction) AddPrivacySetting(setting usermodel.PrivacySetting) error {
	query := fmt.Sprintf(
		`
		INSERT INTO %s (%s,%s) VALUES (
			(SELECT %s FROM %s WHERE %s = $1),
			$2
		)
		`,
		userprivacysettingsql.Table,
		userprivacysettingsql.SettingId,
		userprivacysettingsql.UserId,

		privacysettingsql.ID,
		privacysettingsql.Table,
		privacysettingsql.Type,
	)
	_, err := u.Exec(query, setting, u.userId)
	return err
}

func (u userTransaction) AddAchievement(achtype achievementmodel.AchievementType, level int) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s,%s) VALUES (
			(SELECT %s FROM %s INNER JOIN %s ON %s = %s WHERE %s = $2 AND %s = $3),
			$1
		)`,
		userachievementsql.Table,

		userachievementsql.AchievementId,
		userachievementsql.UserId,

		sqlutils.Full(achievementsql.ID),

		achievementsql.Table,

		achievementtypesql.Table,
		sqlutils.Full(achievementsql.TypeId),
		sqlutils.Full(achievementtypesql.ID),

		sqlutils.Full(achievementtypesql.Type),
		sqlutils.Full(achievementsql.Level),
	)
	_, err := u.Exec(query, u.userId, achtype, level)
	return err
}

func (u userTransaction) AddSubscription(subscription usermodel.Subscription) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s,%s,%s) VALUES (
			(SELECT %s FROM %s WHERE %s = $2),
			$1,
			$3
		)`,
		usersubscriptionsql.Table,

		usersubscriptionsql.TypeId,
		usersubscriptionsql.UserId,
		usersubscriptionsql.Expired,

		subscriptiontypesql.ID,
		subscriptiontypesql.Table,
		subscriptiontypesql.Type,
	)
	_, err := u.Exec(query, u.userId, subscription.Type, subscription.Expired)
	return err
}

func (u userTransaction) AddTrigger(tr usermodel.Trigger) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES (
			(SELECT %s FROM %s WHERE %s = $2),
			$1
		)`,
		usertriggersql.Table,
		usertriggersql.TriggerId,
		usertriggersql.UserId,

		triggersql.ID,
		triggersql.Table,
		triggersql.Name,
	)

	_, err := u.Exec(query, u.userId, tr)
	return err

}

func insertUser(db *sql.DB, userId int64, user *usermodel.CreateUser, subscription usermodel.Subscription, triggers []usermodel.Trigger) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	userTx := userTransaction{Tx: tx, userId: userId, user: user}
	err = userTx.Insert()
	if err != nil {
		slog.Error("user insert")
		return err
	}
	err = userTx.AddSubscription(subscription)
	if err != nil {
		slog.Error("subscription")
		return err
	}
	err = userTx.AddAchievement(achievementmodel.CIGARETTE, 1)
	if err != nil {
		slog.Error("add achievement")
		return err
	}
	err = userTx.AddAchievement(achievementmodel.CIGARETTE, 2)
	if err != nil {
		slog.Error("add achievement")
		return err
	}
	err = userTx.AddPrivacySetting(usermodel.ACHIEVEMENTS_CIGARETTE)
	if err != nil {
		slog.Error("privacy setting")
		return err
	}
	for _, tr := range triggers {
		err = userTx.AddTrigger(tr)
		if err != nil {
			slog.Error("trigger")
			return err
		}
	}
	return tx.Commit()
}

func userDeleted(db *sql.DB, userId int64) (bool, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = $1`,
		usersql.Deleted,
		usersql.Table,
		usersql.ID,
	)
	deleted := false
	err := db.QueryRow(query, userId).Scan(&deleted)
	return deleted, err
}

func verifyRecoveredUser(t *testing.T, user *usermodel.UserData) {
	require.Equal(t, usermodel.One, user.Level.Level, "wrong level")
	require.Equal(t, 0, user.Level.Exp, "wrong exp")
	require.Equal(t, usermodel.Noob, user.Level.Rank, "wrong rank")
}

func achievementsCount(db *sql.DB, userId int64) (int, error) {
	query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE %s = $1`, userachievementsql.Table, userachievementsql.UserId)
	count := 0
	err := db.QueryRow(query, userId).Scan(&count)
	return count, err
}

func privacySettingsCount(db *sql.DB, userId int64) (int, error) {
	query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE %s = $1`, userprivacysettingsql.Table, userprivacysettingsql.UserId)
	count := 0
	err := db.QueryRow(query, userId).Scan(&count)
	return count, err
}
