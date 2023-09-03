package redishelper

import "fmt"

var (
	AchievementsKey    = func(userId int64) string { return fmt.Sprintf("achievements_%d", userId) }
	UsersKey           = func(userId int64) string { return fmt.Sprintf("users_%d", userId) }
	PrivacySettingsKey = func(userId int64) string { return fmt.Sprintf("privacy_settings_%d", userId) }
)
