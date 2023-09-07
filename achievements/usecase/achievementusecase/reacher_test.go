package achievementusecase_test

import (
	"slices"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"gotest.tools/v3/assert"
)

func Test_Reacher_ReachAchievements(t *testing.T) {

	cases := []struct {
		userData     *model.UserData
		achievements []*achievementmodel.Achievement
		expectedIds  []int64
	}{
		{
			userData:     nil,
			achievements: nil,
			expectedIds:  nil,
		},
	}

	for _, cs := range cases {
		reachDate := amidtime.Timestamp{Time: time.Now()}
		reacher := achievementusecase.NewReacher(cs.userData)
		reachAchievementIds := reacher.ReachAchievements(reachDate, cs.achievements)
		equal := slices.Equal(reachAchievementIds, cs.expectedIds)
		assert.Equal(t, true, equal, "achievements not equal")
	}
}
