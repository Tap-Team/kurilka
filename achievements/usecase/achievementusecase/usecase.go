package achievementusecase

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"

	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/internal/domain/achievementreacher"
	"github.com/Tap-Team/kurilka/internal/domain/userstatisticscounter"
	"github.com/Tap-Team/kurilka/internal/errorutils/achievementerror"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source usecase.go -destination mocks.go -package achievementusecase

const _PROVIDER = "achievements/usecase/achievementusecase.useCase"

type AchievementStorage interface {
	AchievementMotivation(ctx context.Context, achId int64) (string, error)
}

type SubscriptionProvider interface {
	UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error)
}

type useCase struct {
	achievementStorage   AchievementStorage
	messageSender        messagesender.MessageSender
	achievement          achievementdatamanager.AchievementManager
	user                 userdatamanager.UserManager
	subscriptionProvider SubscriptionProvider
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
	subscriptionProvider SubscriptionProvider,
) AchievementUseCase {
	return &useCase{achievement: achievement, user: user, messageSender: sender, achievementStorage: achievementStorage, subscriptionProvider: subscriptionProvider}
}

func (u *useCase) OpenSingle(ctx context.Context, userId int64, achievementId int64) (*model.OpenAchievementResponse, error) {
	if !u.canUserOpenAchievement(ctx, userId, achievementId) {
		return nil, achievementerror.ExceptionCantOpenAchievementForUser()
	}
	response, err := u.achievement.OpenSingle(ctx, userId, achievementId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("open single achievement", "OpenSingle", _PROVIDER))
	}
	u.sendAchievementMotivation(ctx, userId, achievementId)
	return response, nil
}

func (u *useCase) sendAchievementMotivation(ctx context.Context, userId, achievementId int64) {
	motivation, err := u.achievementStorage.AchievementMotivation(ctx, achievementId)
	if err != nil {
		slog.ErrorContext(ctx, "failed get achievement motivation for user", "user_id", userId, "achievement_id", achievementId)
		return
	}
	u.messageSender.SendMessage(ctx, motivation, userId)
}

func (u *useCase) canUserOpenAchievement(ctx context.Context, userId int64, achievementId int64) bool {
	achievements, err := u.achievement.UserAchievements(ctx, userId)
	if err != nil {
		return false
	}
	achievement := u.findAchievementById(achievements, achievementId)
	if !u.isAchievementLevelTheLeastInOwnType(achievements, &achievement) {
		return false
	}
	achievementType := achievement.Type
	if achievementType == achievementmodel.CIGARETTE || achievementType == achievementmodel.HEALTH {
		return u.isUserHaveSubscription(ctx, userId)
	}
	return true
}

func (u *useCase) isUserHaveSubscription(ctx context.Context, userId int64) bool {
	subscription, err := u.subscriptionProvider.UserSubscription(ctx, userId)
	if err != nil {
		return false
	}
	return !subscription.IsNoneOrExpired()
}

func (u *useCase) isAchievementLevelTheLeastInOwnType(achievements []*achievementmodel.Achievement, achievement *achievementmodel.Achievement) bool {
	reachedAndNotOpenedAchievements := u.filterReachedAndNotOpenedAchievements(achievements)
	achievementsByType := u.filterAchievementsByType(reachedAndNotOpenedAchievements, achievement.Type)
	leastLevel := u.leastAchievementsLevel(achievementsByType)
	return achievement.Level == leastLevel
}

func (u *useCase) leastAchievementsLevel(achievements []*achievementmodel.Achievement) (level int) {
	level = math.MaxInt
	for i := range achievements {
		ach := achievements[i]
		if ach.Level < level {
			level = ach.Level
		}
	}
	return
}

func (u *useCase) findAchievementById(achievements []*achievementmodel.Achievement, achievementId int64) (achievement achievementmodel.Achievement) {
	for _, ach := range achievements {
		if ach.ID == achievementId {
			achievement = *ach
			break
		}
	}
	return
}

func (u *useCase) filterReachedAndNotOpenedAchievements(achievements []*achievementmodel.Achievement) []*achievementmodel.Achievement {
	filteredList := make([]*achievementmodel.Achievement, 0)
	for i := range achievements {
		ach := achievements[i]
		if ach.Opened() {
			continue
		}
		if ach.Reached() {
			filteredList = append(filteredList, ach)
		}
	}
	return filteredList
}

func (u *useCase) filterAchievementsByType(achievements []*achievementmodel.Achievement, achievementType achievementmodel.AchievementType) []*achievementmodel.Achievement {
	filteredList := make([]*achievementmodel.Achievement, 0)
	for i := range achievements {
		ach := achievements[i]
		if ach.Type == achievementType {
			filteredList = append(filteredList, ach)
		}
	}
	return filteredList
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
	counter := userstatisticscounter.NewCounter(
		reachDate,
		user.AbstinenceTime,
		int(user.CigaretteDayAmount),
		int(user.CigarettePackAmount),
		float64(user.PackPrice),
		userstatisticscounter.Second,
	)
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
