package reachachievementexecutor

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/internal/domain/achievementreacher"
	"github.com/Tap-Team/kurilka/internal/domain/userstatisticscounter"
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
	counter := userstatisticscounter.NewCounter(reachDate, user.AbstinenceTime, int(user.CigaretteDayAmount), int(user.CigarettePackAmount), float64(user.PackPrice), userstatisticscounter.Day)
	fabric := achievementreacher.NewPercentableFabric(counter.Cigarette(), int(counter.Money()), user.AbstinenceTime)
	reacher := achievementreacher.NewReacher(fabric)
	reachAchievements := reacher.ReachAchievements(reachDate, achievements)
	return reachAchievements
}
