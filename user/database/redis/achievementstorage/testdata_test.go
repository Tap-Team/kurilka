package achievementstorage_test

import (
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

func Achievement(tp achievementmodel.AchievementType, level int, open bool) *achievementmodel.Achievement {
	var openDate time.Time
	if open {
		openDate = time.Now()
	}
	return achievementmodel.NewAchievement(0, tp, level, 0, amidtime.Timestamp{Time: openDate}, amidtime.Timestamp{}, false, 0, "")
}
func UserAchievement(tp achievementmodel.AchievementType, level int) *usermodel.Achievement {
	a := usermodel.NewAсhievement(tp, level)
	return &a
}

func compareUserAchievements(uach1, uach2 *usermodel.Achievement) bool {
	if uach1.Level != uach2.Level {
		return false
	}
	if uach1.Type != uach2.Type {
		return false
	}
	return true
}
