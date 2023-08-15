package usermodel

import (
	"time"

	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type Achivement struct {
	Type  string `json:"type"`
	Level int    `json:"level"`
}

func NewAсhievement(tp string, level int) Achivement {
	return Achivement{Type: tp, Level: level}
}

type Friend struct {
	ID                 int64              `json:"id"`
	CigaretteTime      amidtime.Timestamp `json:"cigaretteTime"`
	Life               int                `json:"life"`
	Cigarette          int                `json:"cigarette"`
	Money              int                `json:"money"`
	SubscriptionStatus SubscriptionType   `json:"subscriptionStatus"`
	Level              LevelInfo          `json:"level"`
	Achievements       []Achivement       `json:"achivements"`
}

func NewFriend(
	id int64,
	cigaretteTime time.Time,
	life, cigarette, money int,
	level LevelInfo,
	achievements []Achivement,
) Friend {
	return Friend{
		ID:            id,
		CigaretteTime: amidtime.Timestamp{Time: cigaretteTime},
		Life:          life,
		Cigarette:     cigarette,
		Money:         money,
		Level:         level,
		Achievements:  achievements,
	}
}
