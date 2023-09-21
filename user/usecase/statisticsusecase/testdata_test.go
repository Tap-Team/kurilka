package statisticsusecase_test

import (
	"math"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/model"
)

func floatStatisticsUnitEqual(u1, u2 model.FloatStatisticsUnit) bool {
	return math.Abs(float64(u1-u2)) < 0.01
}

func floatStatisticsEqual(s1, s2 model.FloatUserStatistics) bool {
	return floatStatisticsUnitEqual(s1.Day, s2.Day) &&
		floatStatisticsUnitEqual(s1.Week, s2.Week) &&
		floatStatisticsUnitEqual(s1.Month, s2.Month) &&
		floatStatisticsUnitEqual(s1.Year, s1.Year)
}

func moneyUserData(packPrice float64, cigarettePackAmount int, cigaretteDayAmount int) *usermodel.UserData {
	return usermodel.NewUserData("", uint8(cigaretteDayAmount), uint8(cigarettePackAmount), float32(packPrice), "", "", time.Time{}, usermodel.LevelInfo{}, make([]usermodel.Trigger, 0))
}

func timeUserData(cigaretteDayAmount int) *usermodel.UserData {
	return usermodel.NewUserData("", uint8(cigaretteDayAmount), 0, 0, "", "", time.Time{}, usermodel.LevelInfo{}, make([]usermodel.Trigger, 0))
}
