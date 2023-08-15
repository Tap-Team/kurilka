package userachievementstatussql

/*
create table if not exists user_achievements_status (
    id smallserial primary key,
    status varchar(30) not null,

    constraint user_achievements_status_unique unique (status)
);
*/

const Table = "user_achievements_status"

type userachievementsstatus_column string

func (c userachievementsstatus_column) String() string {
	return string(c)
}

func (c userachievementsstatus_column) Table() string {
	return Table
}

const (
	ID     userachievementsstatus_column = "id"
	Status userachievementsstatus_column = "status"

	ContaintStatusUnique = "user_achievements_status_unique"
)
