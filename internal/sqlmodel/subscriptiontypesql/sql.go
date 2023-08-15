package subscriptiontypesql

/*
create table if not exists subscription_types (
    id smallserial primary key,
    type varchar(16) not null,

    constraint subscription_types_unique unique (type)
);
*/

const Table = "subscription_types"

type subscription_types_column string

func (c subscription_types_column) String() string {
	return string(c)
}

func (c subscription_types_column) Table() string {
	return Table
}

const (
	ID   subscription_types_column = "id"
	Type subscription_types_column = "type"

	ConstraintTypeUnique = "subscription_types_unique"
)
