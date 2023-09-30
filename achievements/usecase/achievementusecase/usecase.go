package achievementusecase

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/domain/achievementreacher"
	"github.com/Tap-Team/kurilka/internal/domain/userstatisticscounter"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source usecase.go -destination mocks.go -package achievementusecase

const _PROVIDER = "achievements/usecase/achievementusecase.useCase"

type AchievementStorage interface {
	AchievementMotivation(ctx context.Context, achId int64) (string, error)
}

type useCase struct {
	achievementStorage AchievementStorage
	messageSender      messagesender.MessageSender
	achievement        achievementdatamanager.AchievementManager
	user               userdatamanager.UserManager
}

type AchievementUseCase interface {
	UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error)
	OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error)
	// OpenType(ctx context.Context, userId int64, achtype achievementmodel.AchievementType) (*model.OpenAchievementResponse, error)
	// OpenAll(ctx context.Context, userId int64) (*model.OpenAchievementResponse, error)
	MarkShown(ctx context.Context, userId int64) error
	UserReachedAchievements(ctx context.Context, userId int64) (model.ReachedAchievements, error)
}

func New(
	achievement achievementdatamanager.AchievementManager,
	user userdatamanager.UserManager,
	achievementStorage AchievementStorage,
	sender messagesender.MessageSender,
) AchievementUseCase {
	return &useCase{achievement: achievement, user: user, messageSender: sender, achievementStorage: achievementStorage}
}

func (u *useCase) OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error) {
	response, err := u.achievement.OpenSingle(ctx, userId, achievementId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("open single achievement", "OpenSingle", _PROVIDER))
	}
	motivation, err := u.achievementStorage.AchievementMotivation(ctx, achievementId)
	if err != nil {
		return response, nil
	}
	u.messageSender.SendMessage(ctx, motivation, userId)
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
	reachDate := time.Now()
	counter := userstatisticscounter.NewCounter(reachDate, user.AbstinenceTime, int(user.CigaretteDayAmount), int(user.CigarettePackAmount), float64(user.PackPrice), userstatisticscounter.Second)
	fabric := achievementreacher.NewPercentableFabric(counter.Cigarette(), int(counter.Money()), user.AbstinenceTime)
	reacher := achievementreacher.NewReacher(fabric)
	reachAchievements := reacher.ReachAchievements(reachDate, achievements)
	u.achievement.ReachAchievements(ctx, userId, amidtime.Timestamp{Time: reachDate}, reachAchievements)
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

func (u *useCase) UserReachedAchievements(ctx context.Context, userId int64) (model.ReachedAchievements, error) {
	var reachedAchievements model.ReachedAchievements
	achievements, err := u.achievement.UserAchievements(ctx, userId)
	if err != nil {
		return reachedAchievements, exception.Wrap(err, exception.NewCause("get user achievements", "UserReachedAchievements", _PROVIDER))
	}
	for _, ach := range achievements {
		if ach.Opened() || !ach.Reached() {
			continue
		}
		reachedAchievements.Type = ach.Type
		switch ach.Type {
		case achievementmodel.CIGARETTE:
			reachedAchievements.Cigarette++
		case achievementmodel.DURATION:
			reachedAchievements.Duration++
		case achievementmodel.HEALTH:
			reachedAchievements.Health++
		case achievementmodel.SAVING:
			reachedAchievements.Saving++
		case achievementmodel.WELL_BEING:
			reachedAchievements.WellBeing++
		}
	}
	return reachedAchievements, nil

}
