package reachachievementexecutor

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/internal/domain/achievementreacher"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
)

//go:generate mockgen -source achievement_reacher.go -destination achievement_reacher_mocks.go -package reachachievementexecutor

type AchievementUserReacher interface {
	ReachAchievements(ctx context.Context, userId int64, user *model.UserData, achievements []*achievementmodel.Achievement) []int64
}

type reacher struct{}

func NewAchievementReacher() AchievementUserReacher {
	return &reacher{}
}

func (r *reacher) ReachAchievements(ctx context.Context, userId int64, user *model.UserData, achievements []*achievementmodel.Achievement) []int64 {
	reachDate := time.Now()
	days := int(time.Now().Sub(user.AbstinenceTime).Hours() / 24)
	cigarette := days * int(user.CigaretteDayAmount)
	singleCigaretteCost := float64(user.PackPrice) / float64(user.CigarettePackAmount)
	money := int(float64(cigarette) * singleCigaretteCost)
	fabric := achievementreacher.NewPercentableFabric(cigarette, money, user.AbstinenceTime)

	reacher := achievementreacher.NewReacher(fabric)
	reachAchievements := reacher.ReachAchievements(reachDate, achievements)
	return reachAchievements
}
