package usersubscriptionsql

/*
create table if not exists user_subscriptions (
    user_id bigint not null,
    type_id bigint not null,

    expired timestamp(0) not null,

    constraint user_subscriptions__users foreign key (user_id) references users(id),
    constraint user_subscriptions__subscription_types foreign key (type_id),

    constraint user_subscriptions_key primary key (user_id, subscription_type_id)
);
*/

const Table = "user_subscriptions"

type user_subscriptions_column string

func (c user_subscriptions_column) String() string {
	return string(c)
}

func (c user_subscriptions_column) Table() string {
	return Table
}

const (
	UserId  user_subscriptions_column = "user_id"
	TypeId  user_subscriptions_column = "type_id"
	Expired user_subscriptions_column = "expired"

	ForeignKeyUsers             = "user_subscriptions__users"
	ForeignKeySubscriptionTypes = "user_subscriptions__subscription_types"
	PrimaryKey                  = "user_subscriptions_key"
)
