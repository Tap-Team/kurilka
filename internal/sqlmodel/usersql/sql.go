package usersql

/*
create table if not exists users (
    id bigint not null,
    name varchar(15) not null,

    cigarette_day_amount int not null,
    cigarette_pack_amount int not null,
    pack_price real not null,
    cigarrete_time timestamp(0) not null default now(),

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
	CigaretteTime       users_column = "cigarrete_time"

	PrimaryKey users_column = "users_key"
)
