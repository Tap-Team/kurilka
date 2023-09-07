package achievementusecase

import (
	"time"

	"github.com/Tap-Team/kurilka/achievements/domain/achievementpercent"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type userReacher struct {
	data *model.UserData
}

type Reacher interface {
	ReachAchievements(reachDate amidtime.Timestamp, achievements []*achievementmodel.Achievement) []int64
}

func NewReacher(data *model.UserData) Reacher {
	return &userReacher{data: data}
}

func (u *userReacher) ReachAchievements(reachDate amidtime.Timestamp, achievements []*achievementmodel.Achievement) []int64 {
	days := int(time.Now().Sub(u.data.AbstinenceTime).Hours() / 24)
	cigarette := days * int(u.data.CigaretteDayAmount)
	singleCigaretteCost := float64(u.data.PackPrice) / float64(u.data.CigarettePackAmount)
	money := int(float64(cigarette) * singleCigaretteCost)
	fabric := achievementpercent.NewFabric(cigarette, money, u.data.AbstinenceTime)

	reachAchievements := make([]int64, 0)
	for i, ach := range achievements {
		achtype := ach.Type
		level := ach.Level
		percent := fabric.Percentable(achtype, level).Percent()
		achievements[i].Percent = percent
		if percent == 100 && !ach.Reached() {
			achievements[i].ReachDate = reachDate
			reachAchievements = append(reachAchievements, ach.ID)
		}
	}
	return reachAchievements
}
