package levelsql

/*
create table if not exists levels (
    level smallint not null,
    rank varchar(30) not null,
    min_exp integer not null,
    max_exp integer not null,

    constraint level_key primary key (level)
);
*/

const Table = "levels"

type levels_column string

func (c levels_column) String() string {
	return string(c)
}

func (c levels_column) Table() string {
	return Table
}

const (
	Level  levels_column = "level"
	Rank   levels_column = "rank"
	MinExp levels_column = "min_exp"
	MaxExp levels_column = "max_exp"

	PrimaryKey = "level_key"
)
