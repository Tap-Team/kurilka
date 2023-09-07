package achievement_test

import "github.com/Tap-Team/kurilka/internal/model/achievementmodel"

func compareAchievements(a1, a2 *achievementmodel.Achievement) bool {
	if a1.ID != a2.ID {
		return false
	}
	if a1.Type != a2.Type {
		return false
	}
	if a1.Exp != a2.Exp {
		return false
	}
	if a1.Level != a2.Level {
		return false
	}
	if a1.Shown != a2.Shown {
		return false
	}
	return true
}
