package achievementdatamanager_test

import (
	"time"

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
