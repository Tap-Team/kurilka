package achievementreacher

import (
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
)

type userReacher struct {
	fabric PercentableFabric
}

type Reacher interface {
	ReachAchievements(reachDate time.Time, achievements []*achievementmodel.Achievement) []int64
}

func NewReacher(fabric PercentableFabric) Reacher {
	return &userReacher{fabric: fabric}
}

func (ur *userReacher) ReachAchievements(reachDate time.Time, achievements []*achievementmodel.Achievement) []int64 {
	reachAchievements := make([]int64, 0)
	for i, ach := range achievements {
		achtype := ach.Type
		level := ach.Level
		percent := ur.fabric.Percentable(achtype, level).Percent()
		achievements[i].Percent = percent
		if percent == 100 && !ach.Reached() {
			achievements[i].SetReachDate(reachDate)
			reachAchievements = append(reachAchievements, ach.ID)
		}
	}
	return reachAchievements
}
