package achievementsql

/*
create table if not exists achievements (
    id smallserial primary key,
    level smallint not null,
    type_id smallint not null,
    exp integer not null,

    constraint achievements__achievements_type references achievements_type(id),
    constraint achievements_unique unique (level, type_id)
);
*/

const Table = "achievements"

type achievements_column string

func (c achievements_column) String() string {
	return string(c)
}

func (c achievements_column) Table() string {
	return Table
}

const (
	ID     achievements_column = "id"
	Level  achievements_column = "level"
	TypeId achievements_column = "type_id"
	Exp    achievements_column = "exp"

	ConstraintAchievementsTypeForeignKey = "achievements__achievements_type"
	ConstraintAchievementsUnique         = "achievements_unique"
)
