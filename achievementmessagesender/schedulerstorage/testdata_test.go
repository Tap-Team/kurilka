package schedulerstorage_test

import (
	"math/rand"

	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
)

func RandomAchievementMessage() model.AchievementMessageData {
	types := []achievementmodel.AchievementType{achievementmodel.CIGARETTE, achievementmodel.HEALTH, achievementmodel.SAVING, achievementmodel.DURATION, achievementmodel.WELL_BEING}
	i := rand.Intn(len(types))
	tp := types[i]
	return model.NewAchievementMessageData(tp)
}

func achievementMessageEqual(m1, m2 *model.MessageData) bool {
	if m1.UserId() != m2.UserId() {
		return false
	}
	if m1.IsDeleted() != m2.IsDeleted() {
		return false
	}
	return true
}
