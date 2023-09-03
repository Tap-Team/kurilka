package achievementusecase

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager"
	"github.com/Tap-Team/kurilka/achievements/domain/achievementpercent"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = ""

type UseCase struct {
	achievement datamanager.AchievementDataManager
	user        datamanager.UserDataManager
}

func New(datamanager datamanager.AchievementDataManager) *UseCase {
	return &UseCase{achievement: datamanager}
}

func (u *UseCase) OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error) {
	response, err := u.achievement.OpenSingle(ctx, userId, achievementId)
	if err != nil {
		return response, exception.Wrap(err, exception.NewCause("open single achievement", "OpenSingle", _PROVIDER))
	}
	return response, nil
}

func (u *UseCase) MarkShown(ctx context.Context, userId int64) error {
	err := u.achievement.MarkShown(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("mark achievement shown", "MarkShown", _PROVIDER))
	}
	return nil
}

type UserReacher struct {
	userId int64
	data   *model.UserData
}

func (u UserReacher) ReachAchievements(ctx context.Context, reachDate amidtime.Timestamp, achievements []*achievementmodel.Achievement) []int64 {
	days := int(time.Now().Sub(u.data.AbstinenceTime.Time).Hours() / 24)
	cigarette := days * int(u.data.CigaretteDayAmount)
	singleCigaretteCost := float64(u.data.PackPrice) / float64(u.data.CigarettePackAmount)
	money := int(float64(cigarette) * singleCigaretteCost)
	fabric := achievementpercent.NewFabric(cigarette, money, u.data.AbstinenceTime.Time)

	reachAchievements := make([]int64, 0)
	for i, ach := range achievements {
		achtype := ach.Type
		level := ach.Level
		percent := fabric.Percentable(achtype, level).Percent()
		achievements[i].Percent = percent
		if percent == 100 && !ach.Reached() {
			achievements[i].ReachDate = reachDate
			reachAchievements = append(reachAchievements, ach.ID)
		}
	}
	return reachAchievements
}

func (u *UseCase) checkReached(ctx context.Context, userId int64, user *model.UserData, achievements []*achievementmodel.Achievement) {
	reachDate := amidtime.Timestamp{Time: time.Now()}
	reachAchievements := UserReacher{userId: userId, data: user}.ReachAchievements(ctx, reachDate, achievements)
	u.achievement.ReachAchievements(ctx, userId, reachDate, reachAchievements)
}

func (u *UseCase) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	achievements, err := u.achievement.UserAchievements(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user achievements", "UserAchievements", _PROVIDER))
	}
	user, err := u.user.UserData(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user data", "UserAchievements", _PROVIDER))
	}
	u.checkReached(ctx, userId, user, achievements)
	return achievements, nil
}
