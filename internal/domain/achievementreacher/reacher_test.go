package achievementreacher_test

import (
	"testing"
)

func Test_Reacher_ReachAchievements(t *testing.T) {

	// cases := []struct {
	// 	userData     *model.UserData
	// 	achievements []*achievementmodel.Achievement
	// 	expectedIds  []int64
	// }{
	// 	{
	// 		userData:     nil,
	// 		achievements: nil,
	// 		expectedIds:  nil,
	// 	},
	// }

	// for _, cs := range cases {
	// 	reachDate := time.Now()
	// 	days := int(time.Now().Sub(cs.userData.AbstinenceTime).Hours() / 24)
	// 	cigarette := days * int(cs.userData.CigaretteDayAmount)
	// 	singleCigaretteCost := float64(cs.userData.PackPrice) / float64(cs.userData.CigarettePackAmount)
	// 	money := int(float64(cigarette) * singleCigaretteCost)
	// 	fabric := achievementreacher.NewPercentableFabric(cigarette, money, cs.userData.AbstinenceTime)
	// 	reacher := achievementreacher.NewReacher(fabric)
	// 	reachAchievementIds := reacher.ReachAchievements(reachDate, cs.achievements)
	// 	equal := slices.Equal(reachAchievementIds, cs.expectedIds)
	// 	assert.Equal(t, true, equal, "achievements not equal")
	// }
}
