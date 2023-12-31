package usermodel

import (
	"time"

	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type LevelInfo struct {
	Level  Level `json:"level"`
	Rank   Rank  `json:"rank"`
	MinExp int   `json:"minExp"`
	MaxExp int   `json:"maxExp"`
	Exp    int   `json:"exp"`
}

func NewLevelInfo(
	level Level,
	rank Rank,
	min, max, exp int,

) LevelInfo {
	return LevelInfo{
		Level:  level,
		Rank:   rank,
		Exp:    exp,
		MinExp: min,
		MaxExp: max,
	}
}

type Subscription struct {
	Type    SubscriptionType   `json:"type"`
	Expired amidtime.Timestamp `json:"expired"`
}

func NewSubscription(status SubscriptionType, expired time.Time) Subscription {
	return Subscription{
		Type:    status,
		Expired: amidtime.Timestamp{Time: expired},
	}
}

func (s *Subscription) SetExpired(expired time.Time) {
	s.Expired = amidtime.Timestamp{Time: expired}
}

func (s Subscription) IsExpired() bool {
	return time.Now().After(s.Expired.Time)
}

func (s Subscription) IsNoneOrExpired() bool {
	return s.Type == NONE || s.IsExpired()
}

type User struct {
	ID int64 `json:"id"`

	// Имя пользователя
	Name Name `json:"name"`

	// Момент когда пользователь перестал курить, просто момент времени, ты должен отнимать от текущего времени пользователя по UTC это время и получать время которое пользователь воздерживается
	AbstinenceTime amidtime.Timestamp `json:"abstinenceTime"`
	// Параметр жизни пользователя, измеряется в минутах
	Life int `json:"life"`
	// Количество не выкуренных пользователем сигарет
	Cigarette int `json:"cigarette"`
	// Сэкономленные пользователем средства
	Money Money `json:"money"`
	// Время которое пользователь секономил на сигаретах, измеряется в минутах
	Time int `json:"time"`
	// Текст Баннера мотивации
	Motivation string `json:"motivation"`
	// Текст приветственной мотивашки
	WelcomeMotivation string `json:"welcomeMotivation"`

	Level LevelInfo `json:"level"`
	// Subscription Subscription `json:"subscription"`

	// Триггеры от которых зависит должен баннер показываться или нет
	Triggers []Trigger `json:"triggers"`
}

func NewUser(
	id int64,
	name Name,
	abstinenceTime time.Time,
	life, cigarette, time int,
	money float64,
	motivation, welcomeMotivation string,
	level LevelInfo,
	triggers []Trigger,
) *User {
	return &User{
		ID:                id,
		Name:              name,
		AbstinenceTime:    amidtime.Timestamp{Time: abstinenceTime},
		Life:              life,
		Cigarette:         cigarette,
		Time:              time,
		Money:             Money(money),
		Motivation:        motivation,
		WelcomeMotivation: welcomeMotivation,
		Level:             level,
		Triggers:          triggers,
	}
}
