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
	exp int,
) LevelInfo {
	return LevelInfo{
		Level: level,
		Rank:  rank,
		Exp:   exp,
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

type User struct {
	ID            int64              `json:"id"`
	CigaretteTime amidtime.Timestamp `json:"cigaretteTime"`
	Life          int                `json:"life"`
	Cigarette     int                `json:"cigarette"`
	Money         int                `json:"money"`
	Friends       []int64            `json:"friends"`
	Level         LevelInfo          `json:"level"`
	Subscription  Subscription       `json:"subscription"`
}

func NewUser(
	id int64,
	cigaretteTime time.Time,
	life, cigarette, money int,
	level LevelInfo,
	subscription Subscription,
	friends []int64,
) User {
	return User{
		ID:            id,
		CigaretteTime: amidtime.Timestamp{Time: cigaretteTime},
		Life:          life,
		Cigarette:     cigarette,
		Money:         money,
		Level:         level,
		Friends:       friends,
	}
}
