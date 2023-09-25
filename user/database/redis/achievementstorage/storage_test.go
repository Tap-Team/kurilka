package achievementstorage_test

import (
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/database/redis/achievementstorage"
	"gotest.tools/v3/assert"
)

func Test_FilterMaxLevelFromAchievementList(t *testing.T) {
	cases := []struct {
		achievements []*achievementmodel.Achievement

		expected []*usermodel.Achievement
	}{
		{
			achievements: []*achievementmodel.Achievement{
				Achievement(achievementmodel.CIGARETTE, 10, false),
				Achievement(achievementmodel.CIGARETTE, 4, true),
				Achievement(achievementmodel.CIGARETTE, 6, true),

				Achievement(achievementmodel.HEALTH, 9, false),
				Achievement(achievementmodel.HEALTH, 3, true),
				Achievement(achievementmodel.HEALTH, 4, true),

				Achievement(achievementmodel.WELL_BEING, 7, false),
				Achievement(achievementmodel.WELL_BEING, 6, true),
				Achievement(achievementmodel.WELL_BEING, 3, true),
			},
			expected: []*usermodel.Achievement{
				UserAchievement(achievementmodel.HEALTH, 4),
				UserAchievement(achievementmodel.WELL_BEING, 6),
				UserAchievement(achievementmodel.CIGARETTE, 6),
			},
		},
	}

	for _, cs := range cases {
		actual := achievementstorage.FilterMaxLevelFromAchievementList(cs.achievements)
		assert.Equal(t, true, slices.EqualFunc(actual, cs.expected, compareUserAchievements), "not equal")
	}
}
