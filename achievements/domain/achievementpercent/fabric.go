package achievementpercent

import (
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
)

type Percentable interface {
	Percent() int
}

type nilPercent struct{}

func (n nilPercent) Percent() int {
	return 0
}

type PercentableFabric interface {
	Percentable(achtype achievementmodel.AchievementType, level int) Percentable
}

type percentFabric struct {
	cigarette      int
	abstinenceTime time.Time
	money          int
}

func NewFabric(cigarette, money int, abstinenceTime time.Time) PercentableFabric {
	return &percentFabric{cigarette: cigarette, money: money, abstinenceTime: abstinenceTime}
}

func (p *percentFabric) Percentable(achtype achievementmodel.AchievementType, level int) Percentable {

	switch achtype {
	case achievementmodel.DURATION:
		return DurationPercentable{abstinenceTime: p.abstinenceTime, level: level}
	case achievementmodel.CIGARETTE:
		return CigarettePercentable{cigaretteAmount: p.cigarette, level: level}
	case achievementmodel.HEALTH:
		return HealthPercentable{abstinenceTime: p.abstinenceTime, level: level}
	case achievementmodel.WELL_BEING:
		return WellBeingPercentable{abstinenceTime: p.abstinenceTime, level: level}
	case achievementmodel.SAVING:
		return SavingPercentable{money: p.money, level: level}
	default:
		return nilPercent{}
	}
}
