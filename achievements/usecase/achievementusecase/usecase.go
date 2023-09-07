package achievementusecase

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = ""

type useCase struct {
	achievement achievementdatamanager.AchievementManager
	user        userdatamanager.UserManager
}

type AchievementUseCase interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error)
	// OpenType(ctx context.Context, userId int64, achtype achievementmodel.AchievementType) (*model.OpenAchievementResponse, error)
	// OpenAll(ctx context.Context, userId int64) (*model.OpenAchievementResponse, error)
	MarkShown(ctx context.Context, userId int64) error
}

func New(achievement achievementdatamanager.AchievementManager, user userdatamanager.UserManager) AchievementUseCase {
	return &useCase{achievement: achievement, user: user}
}

func (u *useCase) OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error) {
	response, err := u.achievement.OpenSingle(ctx, userId, achievementId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("open single achievement", "OpenSingle", _PROVIDER))
	}
	return response, nil
}

func (u *useCase) MarkShown(ctx context.Context, userId int64) error {
	err := u.achievement.MarkShown(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("mark achievement shown", "MarkShown", _PROVIDER))
	}
	return nil
}

func (u *useCase) ReachAchievements(ctx context.Context, userId int64, user *model.UserData, achievements []*achievementmodel.Achievement) {
	reachDate := amidtime.Timestamp{Time: time.Now()}
	reacher := NewReacher(user)
	reachAchievements := reacher.ReachAchievements(reachDate, achievements)
	u.achievement.ReachAchievements(ctx, userId, reachDate, reachAchievements)
}

func (u *useCase) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	achievements, err := u.achievement.UserAchievements(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user achievements", "UserAchievements", _PROVIDER))
	}
	user, err := u.user.UserData(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user data", "UserAchievements", _PROVIDER))
	}
	u.ReachAchievements(ctx, userId, user, achievements)
	return achievements, nil
}
