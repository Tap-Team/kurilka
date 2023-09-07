package achievementdatamanager_test

import (
	"fmt"
	"time"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/golang/mock/gomock"
)

type achievementsMatcher struct {
	achievementIds []int64
	matchFunc      func(exists bool, ach *achievementmodel.Achievement) bool
}

func (m *achievementsMatcher) Matches(x interface{}) bool {
	achlist, ok := x.([]*achievementmodel.Achievement)
	if !ok {
		return false
	}
	ids := make(map[int64]struct{}, len(m.achievementIds))
	for _, i := range m.achievementIds {
		ids[i] = struct{}{}
	}

	for _, ach := range achlist {
		_, ok := ids[ach.ID]
		ok = m.matchFunc(ok, ach)
		if !ok {
			return false
		}
		delete(ids, ach.ID)
	}
	return len(ids) == 0
}

func (m *achievementsMatcher) String() string {
	return "every achievement is match func"
}

func NewAchievementsMatcher(achievementsIds []int64, matchFunc func(exists bool, ach *achievementmodel.Achievement) bool) gomock.Matcher {
	return &achievementsMatcher{
		achievementIds: achievementsIds,
		matchFunc:      matchFunc,
	}
}

func NewOpenAchievementsMatcher(openAchievementIds []int64, openDate time.Time) gomock.Matcher {
	return NewAchievementsMatcher(
		openAchievementIds,
		func(exists bool, ach *achievementmodel.Achievement) bool {
			if !exists {
				return !ach.Opened()
			}
			if !ach.OpenDate.Equal(openDate) {
				return false
			}
			return true
		},
	)
}

func NewReachAchievementsMatcher(reachAchievementIds []int64, reachDate time.Time) gomock.Matcher {
	return NewAchievementsMatcher(
		reachAchievementIds,
		func(exists bool, ach *achievementmodel.Achievement) bool {
			if !exists {
				return !ach.Reached()
			}
			if !ach.ReachDate.Equal(reachDate) {
				return false
			}
			return true
		},
	)
}

func NewShownAchievementsMatcher(shownAchievementIds []int64) gomock.Matcher {
	return NewAchievementsMatcher(
		shownAchievementIds,
		func(exists bool, ach *achievementmodel.Achievement) bool {
			return ach.Shown == ach.Reached()
		},
	)
}

type TimeMatcher time.Time

func (t TimeMatcher) Matches(x interface{}) bool {
	tm, ok := x.(time.Time)
	if !ok {
		return false
	}
	return tm.Unix() == time.Time(t).Unix()
}

func (t TimeMatcher) String() string {
	return fmt.Sprintf("is unix %d second", time.Time(t).Unix())
}

type OpenAchievementMatcher struct {
	model.OpenAchievement
}

func (m *OpenAchievementMatcher) Matches(x interface{}) bool {
	openAch, ok := x.(model.OpenAchievement)
	if !ok {
		return false
	}
	if !TimeMatcher(openAch.OpenTime.Time).Matches(m.OpenTime.Time) {
		return false
	}
	if openAch.AchievementId != m.AchievementId {
		return false
	}
	return true
}

func (m *OpenAchievementMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.OpenAchievement)
}

func NewOpenAchievementMatcher(achievementId int64, openDate time.Time) gomock.Matcher {
	return &OpenAchievementMatcher{
		OpenAchievement: model.NewOpenAchievement(achievementId, openDate),
	}
}
