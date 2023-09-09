package redishelper

import "fmt"

type UserKeyFunc func(int64) string

var (
	AchievementsKey    UserKeyFunc = func(userId int64) string { return fmt.Sprintf("achievements_%d", userId) }
	UsersKey           UserKeyFunc = func(userId int64) string { return fmt.Sprintf("users_%d", userId) }
	PrivacySettingsKey UserKeyFunc = func(userId int64) string { return fmt.Sprintf("privacy_settings_%d", userId) }
	SubscriptionKey    UserKeyFunc = func(userId int64) string { return fmt.Sprintf("subscription_%d", userId) }
)
