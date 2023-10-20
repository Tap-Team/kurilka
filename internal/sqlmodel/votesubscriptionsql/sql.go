package votesubscriptionsql

// create table if not exists vote_subscriptions (
//     subscription_id bigint not null,

//     user_id bigint not null,
//     last_answer jsonb,

//     constraint fk_vote_subscriptions__users foreign key(user_id) references users(id),

//     constraint vote_subscriptions_key primary key(subscription_id)
// );

const Table = "vote_subscriptions"

type vote_subscriptions_column string

func (c vote_subscriptions_column) String() string {
	return string(c)
}

func (c vote_subscriptions_column) Table() string {
	return Table
}

const (
	SubscriptionId vote_subscriptions_column = "subscription_id"
	UserId         vote_subscriptions_column = "user_id"
	LastAnswer     vote_subscriptions_column = "last_answer"

	FKUsers            = "fk_vote_subscriptions__users"
	SubscriptionUnique = "vote_subscriptions__subscription_id_unique"
	PrimaryKey         = "vote_subscriptions_key"
)
