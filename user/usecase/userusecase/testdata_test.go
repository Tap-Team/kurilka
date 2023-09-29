package userusecase_test

import (
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
)

func moneyEqual(money usermodel.Money, target float64) bool {
	difference := money - usermodel.Money(target)
	if difference < 0 {
		difference = -difference
	}
	return difference < 1
}

func daysTime(days int) time.Time {
	seconds := time.Now().Sub(time.Unix(86400*int64(days), 0)).Seconds()
	return time.Unix(int64(seconds), 0)
}

func minutesTime(minutes int) time.Time {
	seconds := time.Now().Sub(time.Unix(60*int64(minutes), 0)).Seconds()
	return time.Unix(int64(seconds), 0)
}

func NewUserDataDays(
	days int,
	cigaretteDayAmount uint8,
	cigarettePackAmount uint8,
	packPrice float32,
) *usermodel.UserData {
	return usermodel.NewUserData(
		"",
		cigaretteDayAmount,
		cigarettePackAmount,
		packPrice,
		"",
		"",
		daysTime(days),
		usermodel.LevelInfo{},
		[]usermodel.Trigger{},
	)
}

func NewUserDataMinutes(
	minutes int,
	cigaretteDayAmount uint8,
	cigarettePackAmount uint8,
	packPrice float32,
) *usermodel.UserData {
	return usermodel.NewUserData(
		"",
		cigaretteDayAmount,
		cigarettePackAmount,
		packPrice,
		"",
		"",
		minutesTime(minutes),
		usermodel.LevelInfo{},
		[]usermodel.Trigger{},
	)
}

func randomAchievementList(size int) []*usermodel.Achievement {
	achlist := make([]*usermodel.Achievement, 0, size)
	for i := 0; i < size; i++ {
		ach := random.StructTyped[usermodel.Achievement]()
		achlist = append(achlist, &ach)
	}
	return achlist
}
