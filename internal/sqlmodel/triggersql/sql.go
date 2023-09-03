package triggersql

/*
create table if not exists triggers (
    id smallserial primary key,
    name varchar(50) not null,

    constraint triggers_name_unique unique(name)
);
*/

const Table = "triggers"

type triggers_column string

func (c triggers_column) Table() string {
	return Table
}

func (c triggers_column) String() string {
	return string(c)
}

const (
	ID   triggers_column = "id"
	Name triggers_column = "name"

	ConstraintNameUnique = "triggers_name_unique"
)
