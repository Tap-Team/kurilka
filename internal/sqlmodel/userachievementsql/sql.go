package userachievementsql

/*
create table if not exists user_achievements (
    achievement_id smallint not null primary key,
    user_id bigint not null,
    status_id smallint not null,
    open_date timestamp(0) default now(),

    constraint user_achievements__user_achievements_status foreign key (status_id) references user_achievements_status(id),
    constraint user_achievements__users foreign key (user_id) references users(id) on delete cascade,
    constraint user_achievements__achievements foreign key (achievement_id) references achievements (id) on delete cascade
);
*/

const Table = "user_achievements"

type user_achievement_column string

func (c user_achievement_column) String() string {
	return string(c)
}

func (c user_achievement_column) Table() string {
	return Table
}

const (
	UserId        user_achievement_column = "user_id"
	AchievementId user_achievement_column = "achievement_id"
	StatusId      user_achievement_column = "status_id"
	OpenDate      user_achievement_column = "open_date"

	ForeignKeyUserAchievementStatus = "user_achievements__user_achievements_status"
	ForeignKeyUserAchievementUsers  = "user_achievements__users"
	ForeignKeyAchievements          = "user_achievements__achievements"
)
