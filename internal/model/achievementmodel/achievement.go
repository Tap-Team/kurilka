package achievementmodel

import (
	"time"

	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type Achievement struct {
	ID   int64           `json:"id"`
	Type AchievementType `json:"type"`
	// количество экспы за открытие
	Exp int `json:"exp"`
	// уровень (от 1 до 10)
	Level int `json:"level"`
	// дата достижение пользователем ачивки по timestamp(0) в секундах, если достижение не достигнуто, равняется 0
	ReachDate amidtime.Timestamp `json:"reachDate"`
	// дата открытия ачивки по timestamp(0) в секундах, если достижение не открыто, равняется 0
	OpenDate amidtime.Timestamp `json:"openDate"`
	// была ли ачивка показана пользователю
	Shown bool `json:"shown"`
	// проценты до достижения (от 0 до 100), на открытых или достигнутых ачивках равняется 100
	Percent int `json:"percentage"`

	Description string `json:"description"`
}

func NewAchievement(
	id int64,
	achtype AchievementType,
	level, exp int,
	openDate amidtime.Timestamp,
	reachDate amidtime.Timestamp,
	shown bool,
	percent int,
	description string,
) *Achievement {
	return &Achievement{
		ID:          id,
		Type:        achtype,
		Exp:         exp,
		Level:       level,
		OpenDate:    openDate,
		Shown:       shown,
		ReachDate:   reachDate,
		Percent:     percent,
		Description: description,
	}
}

func (a *Achievement) SetOpenDate(t time.Time) {
	a.OpenDate = amidtime.Timestamp{Time: t}
}

func (a *Achievement) SetReachDate(t time.Time) {
	a.ReachDate = amidtime.Timestamp{Time: t}
}

func (a Achievement) Opened() bool {
	return !a.OpenDate.IsZero()
}

func (a Achievement) Reached() bool {
	return !a.ReachDate.IsZero()
}
