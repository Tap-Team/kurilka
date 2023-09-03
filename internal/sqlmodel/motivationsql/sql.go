package motivationsql

/*
create table if not exists motivations (
    id smallserial primary key,
    motivation text not null,

    constraint motivations_unique unique (motivation)
);
*/

const Table = "motivations"

type motivation_column string

func (c motivation_column) Table() string {
	return Table
}

func (c motivation_column) String() string {
	return string(c)
}

const (
	ID         motivation_column = "id"
	Motivation motivation_column = "motivation"

	ConstraintMotivationUnique = "motivations_unique"
)
