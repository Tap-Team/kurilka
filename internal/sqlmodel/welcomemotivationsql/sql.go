package welcomemotivationsql

/*
create table if not exists welcome_motivations (
    id smallserial primary key,
    motivation text not null,

    constraint welcome_motivations_unique unique (motivation)
);
*/

const Table = "welcome_motivations"

type welcome_motivation_column string

func (c welcome_motivation_column) String() string {
	return string(c)
}

func (c welcome_motivation_column) Table() string {
	return Table
}

const (
	ID         welcome_motivation_column = "id"
	Motivation welcome_motivation_column = "motivation"

	ConstraintMotivationUnique = "welcome_motivations_unique"
)
