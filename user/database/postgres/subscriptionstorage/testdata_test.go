package subscriptionstorage_test

import (
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/subscriptiontypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersubscriptionsql"
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

func insertUserWithSubscription(db *sql.DB, userId int64, createUser *usermodel.CreateUser, subscription usermodel.Subscription) error {
	if createUser == nil {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx, %s", err)
	}
	defer tx.Rollback()
	insert := userTransaction{Tx: tx, userId: userId, user: createUser}
	err = insert.Insert()
	if err != nil {
		return fmt.Errorf("insert user, %s", err)
	}
	err = insert.AddSubscription(subscription)
	if err != nil {
		return fmt.Errorf("insert subscription, %s", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx, %s", err)
	}
	return nil
}
