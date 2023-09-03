package usermodel

import (
	"slices"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type Achievement struct {
	Type  achievementmodel.AchievementType `json:"type"`
	Level int                              `json:"level"`
}

func NewAсhievement(tp achievementmodel.AchievementType, level int) Achievement {
	return Achievement{Type: tp, Level: level}
}

type Friend struct {
	ID               int64              `json:"id"`
	AbstinenceTime   amidtime.Timestamp `json:"cigaretteTime"`
	Life             int                `json:"life"`
	Cigarette        int                `json:"cigarette"`
	Money            Money              `json:"money"`
	Time             int                `json:"time"`
	SubscriptionType SubscriptionType   `json:"subscriptionType"`
	Level            LevelInfo          `json:"level"`
	Achievements     []*Achievement     `json:"achivements"`
	PrivacySettings  []PrivacySetting   `json:"privacySettings"`
}

func NewFriend(
	id int64,
	abstinenceTime time.Time,
	life, cigarette, time int,
	money float64,
	subscriptionType SubscriptionType,
	level LevelInfo,
	achievements []*Achievement,
	privacySettings []PrivacySetting,
) *Friend {
	return &Friend{
		ID:               id,
		AbstinenceTime:   amidtime.Timestamp{Time: abstinenceTime},
		Life:             life,
		Time:             time,
		Cigarette:        cigarette,
		Money:            Money(money),
		SubscriptionType: subscriptionType,
		Level:            level,
		Achievements:     achievements,
		PrivacySettings:  privacySettings,
	}
}

type PrivacySettingFilter func(f *Friend)

var (
	STATISTICS_MONEY_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Money = 0
	}
	STATISTICS_CIGARETTE_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Cigarette = 0
	}
	STATISTICS_LIFE_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Life = 0
	}
	STATISTICS_TIME_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Time = 0
	}
	ACHIEVEMENTS_DURATION_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Achievements = slices.DeleteFunc(f.Achievements, func(a *Achievement) bool {
			return a.Type == achievementmodel.DURATION
		})
	}
	ACHIEVEMENTS_HEALTH_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Achievements = slices.DeleteFunc(f.Achievements, func(a *Achievement) bool {
			return a.Type == achievementmodel.HEALTH
		})
	}
	ACHIEVEMENTS_WELL_BEING_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Achievements = slices.DeleteFunc(f.Achievements, func(a *Achievement) bool {
			return a.Type == achievementmodel.WELL_BEING
		})
	}
	ACHIEVEMENTS_SAVING_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Achievements = slices.DeleteFunc(f.Achievements, func(a *Achievement) bool {
			return a.Type == achievementmodel.SAVING
		})
	}
	ACHIEVEMENTS_CIGARETTE_FILTER PrivacySettingFilter = func(f *Friend) {
		f.Achievements = slices.DeleteFunc(f.Achievements, func(a *Achievement) bool {
			return a.Type == achievementmodel.CIGARETTE
		})
	}
)

var (
	filters = map[PrivacySetting]PrivacySettingFilter{
		STATISTICS_MONEY:        STATISTICS_MONEY_FILTER,
		STATISTICS_CIGARETTE:    STATISTICS_CIGARETTE_FILTER,
		STATISTICS_LIFE:         STATISTICS_LIFE_FILTER,
		STATISTICS_TIME:         STATISTICS_TIME_FILTER,
		ACHIEVEMENTS_DURATION:   ACHIEVEMENTS_DURATION_FILTER,
		ACHIEVEMENTS_HEALTH:     ACHIEVEMENTS_HEALTH_FILTER,
		ACHIEVEMENTS_WELL_BEING: ACHIEVEMENTS_WELL_BEING_FILTER,
		ACHIEVEMENTS_SAVING:     ACHIEVEMENTS_SAVING_FILTER,
		ACHIEVEMENTS_CIGARETTE:  ACHIEVEMENTS_CIGARETTE_FILTER,
	}
)

func (f *Friend) UseFilters() {
	for _, st := range f.PrivacySettings {
		filter, ok := filters[st]
		if !ok {
			continue
		}
		filter(f)
	}
	if f.SubscriptionType == NONE {
		STATISTICS_TIME_FILTER(f)
		STATISTICS_CIGARETTE_FILTER(f)

		ACHIEVEMENTS_CIGARETTE_FILTER(f)
		ACHIEVEMENTS_HEALTH_FILTER(f)
	}
}
