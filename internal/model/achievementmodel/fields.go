package achievementmodel

import "github.com/Tap-Team/kurilka/internal/errorutils/achievementerror"

type AchievementType string

const (
	DURATION   AchievementType = "Длительность"
	CIGARETTE  AchievementType = "Сигареты"
	HEALTH     AchievementType = "Здоровье"
	WELL_BEING AchievementType = "Самочувствие"
	SAVING     AchievementType = "Экономия"
)

func (a AchievementType) Validate() error {
	for _, tp := range []AchievementType{DURATION, CIGARETTE, HEALTH, WELL_BEING, SAVING} {
		if a == tp {
			return nil
		}
	}
	return achievementerror.ExceptionAchievementNotExists()
}
