package userstorage_test

import (
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/subscriptiontypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersubscriptionsql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
)

func addUserExp(db *sql.DB, userId int64, exp int) error {
	achievementType := random.String(30)
	achievementDescription := random.String(30)
	query := fmt.Sprintf(`
	WITH type_insert as (
		INSERT INTO %s (%s) VALUES ($1) RETURNING %s as id
	),
	achievement_insert as (
		INSERT INTO %s (%s,%s,%s,%s,%s) VALUES (1,(SELECT id from type_insert), $2, $4, $4) RETURNING %s as id
	)
	INSERT INTO %s (%s, %s, %s)
		VALUES (
			$3,
			(SELECT id FROM achievement_insert),
			now()
		)
	`,

		// type insert
		achievementtypesql.Table,

		achievementtypesql.Type,
		achievementtypesql.ID,

		// achievement insert
		achievementsql.Table,

		achievementsql.Level,
		achievementsql.TypeId,
		achievementsql.Exp,
		achievementsql.Description,
		achievementsql.Motivation,

		achievementsql.ID,

		// user achievement insert
		userachievementsql.Table,

		userachievementsql.UserId,
		userachievementsql.AchievementId,
		userachievementsql.OpenDate,
	)
	_, err := db.Exec(query, achievementType, exp, userId, achievementDescription)
	return err
}

func fakeAddUserExp(db *sql.DB, userId int64, exp int) error {
	achievementType := random.String(30)

	query := fmt.Sprintf(`
	WITH type_insert as (
		INSERT INTO %s (%s) VALUES ($1) RETURNING %s as id
	),
	achievement_insert as (
		INSERT INTO %s (%s,%s,%s) VALUES (1,(SELECT id from type_insert), $2) RETURNING %s as id
	)
	INSERT INTO %s (%s, %s, %s)
		VALUES (
			$3,
			(SELECT id FROM achievement_insert)
		)
	`,

		// type insert
		achievementtypesql.Table,

		achievementtypesql.Type,
		achievementtypesql.ID,

		// achievement insert
		achievementsql.Table,

		achievementsql.Level,
		achievementsql.TypeId,
		achievementsql.Exp,

		achievementsql.ID,

		// user achievement insert
		userachievementsql.Table,

		userachievementsql.UserId,
		userachievementsql.AchievementId,
		userachievementsql.OpenDate,
	)
	_, err := db.Exec(query, achievementType, exp, userId)
	return err
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

func removeUser(db *sql.DB, userId int64) error {
	query := fmt.Sprintf(
		`UPDATE %s SET %s = TRUE WHERE %s = $1`,
		usersql.Table,
		usersql.Deleted,
		usersql.ID,
	)
	_, err := db.Exec(query, userId)
	return err
}

// func (u user) SubscriptionQuery() string {
// 	return fmt.Sprintf(
// 		`SELECT %s,%s FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1 GROUP BY %s,%s`,
// 		sqlutils.Full(usersubscriptionsql.Expired),
// 		sqlutils.Full(subscriptiontypesql.Type),

// 		usersubscriptionsql.Table,

// 		subscriptiontypesql.Table,
// 		sqlutils.Full(usersubscriptionsql.TypeId),
// 		sqlutils.Full(subscriptiontypesql.ID),

// 		sqlutils.Full(usersubscriptionsql.UserId),

// 		sqlutils.Full(usersubscriptionsql.Expired),
// 		sqlutils.Full(subscriptiontypesql.Type),
// 	)
// }

// func (u *user) Subscription(ctx context.Context) usermodel.Subscription {
// 	if u.err != nil {
// 		return usermodel.Subscription{}
// 	}
// var subscription usermodel.Subscription
// err := u.tx.QueryRowContext(ctx, u.SubscriptionQuery(), u.userId).Scan(&subscription.Expired, &subscription.Type)
// if err != nil {
// 	u.err = Error(err, exception.NewCause("subscription", "Subscription", _PROVIDER))
// }
// return subscription
// }

func userSubscription(db *sql.DB, userId int64) (usermodel.Subscription, error) {
	query := fmt.Sprintf(
		`SELECT %s,%s FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1 GROUP BY %s,%s`,
		sqlutils.Full(usersubscriptionsql.Expired),
		sqlutils.Full(subscriptiontypesql.Type),

		usersubscriptionsql.Table,

		subscriptiontypesql.Table,
		sqlutils.Full(usersubscriptionsql.TypeId),
		sqlutils.Full(subscriptiontypesql.ID),

		sqlutils.Full(usersubscriptionsql.UserId),

		sqlutils.Full(usersubscriptionsql.Expired),
		sqlutils.Full(subscriptiontypesql.Type),
	)
	var subscription usermodel.Subscription
	err := db.QueryRow(query, userId).Scan(&subscription.Expired, &subscription.Type)
	if err != nil {
		return usermodel.Subscription{}, err
	}
	return subscription, nil
}
