BEGIN;


create table if not exists subscription_types ( 
    id smallserial primary key,
    type varchar(16) not null,

    constraint subscription_types_unique unique (type)
);

create table if not exists users (
    id bigint not null,
    name varchar(15) not null,

    cigarette_day_amount int not null,
    cigarette_pack_amount int not null,
    pack_price real not null,
    cigarrete_time timestamp(0) not null default now(),

    constraint users_key primary key (id)

);

create table if not exists user_subscriptions (
    user_id bigint not null,
    type_id bigint not null,

    expired timestamp(0) not null,

    constraint user_subscriptions__users foreign key (user_id) references users(id) on delete cascade,
    constraint user_subscriptions__subscription_types foreign key (type_id) references subscription_types(id),

    constraint user_subscriptions_key primary key (user_id, type_id)
);


create table if not exists achievements_type (
    id smallserial primary key,
    type varchar(30) not null,

    constraint achievements_type_unique unique (type)
);

create table if not exists achievements (
    id smallserial primary key,
    level smallint not null,
    type_id smallint not null,
    exp integer not null,

    constraint achievements__achievements_type foreign key (type_id) references achievements_type(id),
    constraint achievements_unique unique (level, type_id)
);

create table if not exists user_achievements_status (
    id smallserial primary key,
    status varchar(30) not null,

    constraint user_achievements_status_unique unique (status)
);

create table if not exists user_achievements (
    achievement_id smallint not null primary key,
    user_id bigint not null, 
    status_id smallint not null,
    open_date timestamp(0) default now(),

    constraint user_achievements__user_achievements_status foreign key (status_id) references user_achievements_status(id),
    constraint user_achievements__users foreign key (user_id) references users(id) on delete cascade,
    constraint user_achievements__achievements foreign key (achievement_id) references achievements (id) on delete cascade
);

create table if not exists levels (
    level smallint not null,
    rank varchar(30) not null,
    min_exp integer not null,
    max_exp integer not null,

    constraint level_key primary key (level)
);

COMMIT;


BEGIN;

INSERT INTO user_achievements_status (status) VALUES ('NONE'),('ACHIEVED'),('OPENED');

INSERT INTO subscription_types (type) VALUES ('NONE'), ('TRIAL'), ('BASIC');

INSERT INTO levels (level, rank, min_exp, max_exp) VALUES 
(1, 'Новичок', 0, 99),
(2, '', 100, 199),
(3, '', 200, 299),
(4, '', 300, 399),
(5, '', 400, 499),
(6, '', 500, 599),
(7, '', 600, 699),
(8, '', 700, 799),
(9, '', 800, 899),
(10, '', 900, 1000);

COMMIT;


