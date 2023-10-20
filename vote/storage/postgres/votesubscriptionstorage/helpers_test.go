package votesubscriptionstorage_test

import (
	"database/sql"
	"log"

	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/votesubscriptionsql"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
)

type Checker struct {
	db *sql.DB
}

func NewChecker(db *sql.DB) *Checker {
	return &Checker{db: db}
}

var select_vote_subscription_by_id_query = new(sqlutils.QueryBuilder).
	Select(votesubscriptionsql.SubscriptionId).
	From(votesubscriptionsql.Table).
	WhereColumnEqual(votesubscriptionsql.SubscriptionId, "$1").
	Build()

func (c *Checker) VoteSubscriptionExists(subscriptionId int64) bool {
	var id int64
	err := c.db.QueryRow(select_vote_subscription_by_id_query, subscriptionId).Scan(&id)
	if err != nil {
		log.Printf("get vote subscription by id %d, err: %s", subscriptionId, err)
	}
	return err == nil
}

var select_vote_subscription_by_user_id_query = new(sqlutils.QueryBuilder).
	Select(votesubscriptionsql.SubscriptionId).
	From(votesubscriptionsql.Table).
	WhereColumnEqual(votesubscriptionsql.UserId, "$1").
	Build()

func (c *Checker) UserVoteSubscriptionExists(userId int64) bool {
	var id int64
	err := c.db.QueryRow(select_vote_subscription_by_user_id_query, userId).Scan(&id)
	if err != nil {
		log.Printf("get vote subscription by userId %d, err: %s", userId, err)
	}
	return err == nil
}

type Inserter struct {
	db *sql.DB
}

func NewInserter(db *sql.DB) *Inserter {
	return &Inserter{db: db}
}

// ID                  users_column = "id"
// Name                users_column = "name"
// CigaretteDayAmount  users_column = "cigarette_day_amount"
// CigarettePackAmount users_column = "cigarette_pack_amount"
// PackPrice           users_column = "pack_price"
// AbstinenceTime      users_column = "abstinence_time"
// Deleted             users_column = "deleted"
// MotivationId        users_column = "motivation_id"
// WelcomeMotivationId users_column = "welcome_motivation_id"

var insert_empty_user_query = new(sqlutils.QueryBuilder).
	InsertInto(
		usersql.Table,
		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
	).
	WriteQuery("VALUES ($1, '', 0, 0, 0)").
	Build()

func (i *Inserter) InsertEmptyUser(userId int64) (err error) {
	_, err = i.db.Exec(insert_empty_user_query, userId)
	return
}

var insert_subscription_query = new(sqlutils.QueryBuilder).
	InsertInto(votesubscriptionsql.Table, votesubscriptionsql.SubscriptionId, votesubscriptionsql.UserId).
	WriteQuery("VALUES ($1, $2)").
	Build()

func (i *Inserter) InsertSubscription(subscriptionId, userId int64) (err error) {
	_, err = i.db.Exec(insert_subscription_query, subscriptionId, userId)
	return
}
