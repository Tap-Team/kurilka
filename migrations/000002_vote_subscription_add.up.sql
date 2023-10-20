create table if not exists vote_subscriptions (
    user_id bigint not null,
    subscription_id bigint not null,

    constraint fk_vote_subscriptions__users foreign key(user_id) references users(id),
    constraint vote_subscriptions__subscription_id_unique unique (subscription_id),

    constraint vote_subscriptions_key primary key(user_id)
);


