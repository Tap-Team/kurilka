package reachachievementexecutor_test

import (
	"fmt"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
)

func AchievementList(ids ...int64) []*achievementmodel.Achievement {
	achs := make([]*achievementmodel.Achievement, 0, len(ids))

	for _, id := range ids {
		achs = append(achs, &achievementmodel.Achievement{
			ID: id,
		})
	}
	return achs
}

type TimeSecondsMatcher struct {
	seconds int64
}

func (t TimeSecondsMatcher) Matches(x interface{}) bool {
	tm, ok := x.(time.Time)
	if !ok {
		return false
	}
	return tm.Unix() == t.seconds
}

func (t TimeSecondsMatcher) String() string {
	return fmt.Sprintf("time.Unix() is equal %d", t.seconds)
}
