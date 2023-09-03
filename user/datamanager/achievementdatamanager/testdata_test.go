package achievementdatamanager_test

import (
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
)

func randomAchievementList(size int) []*usermodel.Achievement {
	achievements := make([]*usermodel.Achievement, 0, size)

	for i := 0; i < size; i++ {
		ach := random.StructTyped[usermodel.Achievement]()
		achievements = append(achievements, &ach)
	}
	return achievements
}
