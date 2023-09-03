package usersql

/*
create table if not exists users (
  id bigint not null,
    name varchar(15) not null,

    cigarette_day_amount int not null,
    cigarette_pack_amount int not null,
    pack_price real not null,
    abstinence_time timestamp(0) not null default now(),
    deleted boolean not null default FALSE,
    motivation_id smallint,
    welcome_motivation_id smallint,

    constraint users__motivations foreign key (motivation_id) references motivations(id),

    constraint users_welcome_motivations foreign key (welcome_motivation_id) references welcome_motivations(id),

    constraint users_key primary key (id)
);
*/

const Table = "users"

type users_column string

func (c users_column) String() string {
	return string(c)
}

func (c users_column) Table() string {
	return Table
}

const (
	ID                  users_column = "id"
	Name                users_column = "name"
	CigaretteDayAmount  users_column = "cigarette_day_amount"
	CigarettePackAmount users_column = "cigarette_pack_amount"
	PackPrice           users_column = "pack_price"
	AbstinenceTime      users_column = "abstinence_time"
	Deleted             users_column = "deleted"
	MotivationId        users_column = "motivation_id"
	WelcomeMotivationId users_column = "welcome_motivation_id"

	MotivationsForeignKey        = "users__motivations"
	WelcomeMotivationsForeignKey = "users_welcome_motivations"
	PrimaryKey                   = "users_key"
)
