package usertriggersql

/*
create table if not exists user_triggers (
    trigger_id smallint not null,
    user_id bigint not null,

    constraint user_triggers_primary primary key (trigger_id, user_id)
);
*/

const Table = "user_triggers"

type user_triggers_column string

func (c user_triggers_column) String() string {
	return string(c)

}

func (c user_triggers_column) Table() string {
	return Table
}

const (
	UserId    user_triggers_column = "user_id"
	TriggerId user_triggers_column = "trigger_id"

	PrimaryKey = "user_triggers_primary"
)
