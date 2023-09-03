package userprivacysettingsql

/*
create table if not exists user_privacy_settings (
    user_id bigint not null,
    setting_id smallint not null,

    constraint user_privacy_settings_unique primary key (user_id, settings_id),

    constraint user_privacy_settings__users foreign key (user_id) references users(id),

    constraint user_privacy_settings__privacy_settings foreign key (setting_id) references privacy_settings(id)
);
*/

const Table = "user_privacy_settings"

type user_privacy_settings_column string

func (c user_privacy_settings_column) String() string {
	return string(c)
}

func (c user_privacy_settings_column) Table() string {
	return Table
}

const (
	UserId    user_privacy_settings_column = "user_id"
	SettingId user_privacy_settings_column = "setting_id"

	PrimaryKey = "user_privacy_settings_unique"

	ForeignKeyUsers           = "user_privacy_settings__users"
	ForeignKeyPrivacySettings = "user_privacy_settings__privacy_settings"
)
