package achievementtypesql

/*
create table if not exists achievements_type (
    id smallserial primary key,
    type varchar(30) not null,

    constraint achievements_type_unique unique (type)
);
*/

const Table = "achievements_type"

type achievements_type_column string

func (c achievements_type_column) String() string {
	return string(c)
}

func (c achievements_type_column) Table() string {
	return Table
}

const (
	ID   achievements_type_column = "id"
	Type achievements_type_column = "type"

	ConstraintTypeUnique = "achievements_type_unique"
)
