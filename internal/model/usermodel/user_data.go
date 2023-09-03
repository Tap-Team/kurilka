package usermodel

import (
	"encoding/json"
	"time"

	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type UserData struct {
	Name                Name                `json:"name"`
	CigaretteDayAmount  CigaretteDayAmount  `json:"cigaretteDayAmount"`
	CigarettePackAmount CigarettePackAmount `json:"cigarettePackAmount"`
	PackPrice           PackPrice           `json:"packPrice"`
	AbstinenceTime      amidtime.Timestamp  `json:"abstinenceTime"`

	Motivation        string `json:"motivation"`
	WelcomeMotivation string `json:"welcomeMotivation"`

	Level        LevelInfo    `json:"level"`
	Subscription Subscription `json:"subscription"`
	Triggers     []Trigger    `json:"triggers"`
}

func NewUserData(
	name string,
	cigaretteDayAmount uint8,
	cigarettePackAmount uint8,
	packPrice float32,
	motivation, welcomeMotivation string,
	abstinenceTime time.Time,
	level LevelInfo,
	subscription Subscription,
	triggers []Trigger,
) *UserData {
	return &UserData{
		Name:                Name(name),
		CigaretteDayAmount:  CigaretteDayAmount(cigaretteDayAmount),
		CigarettePackAmount: CigarettePackAmount(cigarettePackAmount),
		PackPrice:           PackPrice(packPrice),
		AbstinenceTime:      amidtime.Timestamp{Time: abstinenceTime},
		Motivation:          motivation,
		WelcomeMotivation:   welcomeMotivation,
		Level:               level,
		Subscription:        subscription,
		Triggers:            triggers,
	}
}

func (u UserData) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
