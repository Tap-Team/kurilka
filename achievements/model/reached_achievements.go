package model

import "github.com/Tap-Team/kurilka/internal/model/achievementmodel"

type ReachedAchievements struct {
	Cigarette int `json:"cigarette"`
	Duration  int `json:"duration"`
	Health    int `json:"health"`
	WellBeing int `json:"well-being"`
	Saving    int `json:"saving"`

	Type achievementmodel.AchievementType `json:"achievementType"`
}

func NewReachedAchievements(cigarette, duration, health, wellBeing, saving int) ReachedAchievements {
	return ReachedAchievements{
		Cigarette: cigarette,
		Duration:  duration,
		Health:    health,
		WellBeing: wellBeing,
		Saving:    saving,
	}
}
