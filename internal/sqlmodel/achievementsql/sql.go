package achievementsql

/*
create table if not exists achievements (
    id smallserial primary key,
    level smallint not null,
    description text not null,
    motivation text not null,
    type_id smallint not null,
    exp integer not null,

    constraint achievements__achievements_type foreign key (type_id) references achievements_type(id),
    constraint achievements_description_unique unique (description),
    constraint achievements_motivation_unique unique (motivation),
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
	ID          achievements_column = "id"
	Level       achievements_column = "level"
	TypeId      achievements_column = "type_id"
	Exp         achievements_column = "exp"
	Description achievements_column = "description"
	Motivation  achievements_column = "motivation"

	ForeignKeyAchievementsType   = "achievements__achievements_type"
	ConstraintAchievementsUnique = "achievements_unique"
	ConstraintMotivationUnique   = "achievements_motivation_unique"
	ConstraintDescriptonUnique   = "achievements_description_unique"
)
