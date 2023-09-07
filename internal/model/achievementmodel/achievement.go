package achievementmodel

import (
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type Achievement struct {
	ID        int64              `json:"id"`
	Type      AchievementType    `json:"type"`
	Exp       int                `json:"exp"`
	Level     int                `json:"level"`
	ReachDate amidtime.Timestamp `json:"reachDate"`
	OpenDate  amidtime.Timestamp `json:"openDate"`
	Shown     bool               `json:"shown"`
	Percent   int                `json:"percentage"`
}

func NewAchievement(
	id int64,
	achtype AchievementType,
	level, exp int,
	openDate amidtime.Timestamp,
	reachDate amidtime.Timestamp,
	shown bool,
	percent int,
) *Achievement {
	return &Achievement{
		ID:        id,
		Type:      achtype,
		Exp:       exp,
		Level:     level,
		OpenDate:  openDate,
		Shown:     shown,
		ReachDate: reachDate,
		Percent:   percent,
	}
}

func (a Achievement) Opened() bool {
	return !a.OpenDate.IsZero()
}

func (a Achievement) Reached() bool {
	return !a.ReachDate.IsZero()
}
