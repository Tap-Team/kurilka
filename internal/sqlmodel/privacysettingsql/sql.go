package privacysettingsql

/*
create table if not exists privacy_settings (
    id smallserial primary key,
    type varchar(30) not null,

    constraint privacy_settings_unique unique (type)
);
*/

const Table = "privacy_settings"

type privacy_settings_column string

func (c privacy_settings_column) String() string {
	return string(c)
}

func (c privacy_settings_column) Table() string {
	return Table
}

const (
	ID   privacy_settings_column = "id"
	Type privacy_settings_column = "type"

	ConstraintSettingsUnique = "privacy_settings_unique"
)
