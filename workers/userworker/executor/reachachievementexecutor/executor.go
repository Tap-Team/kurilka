package reachachievementexecutor

import (
	"context"
	"fmt"
	"time"

	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor"
)

const _PROVIDER = "workers/userworker/executor/reachachievementexecutor.Executor"

type Executor struct {
	user          userdatamanager.UserManager
	achievement   achievementdatamanager.AchievementManager
	reacher       AchievementUserReacher
	messageSender messagesender.MessageSenderAtTime
}

func NewReachAchievementExecutor(
	user userdatamanager.UserManager,
	achievement achievementdatamanager.AchievementManager,
	messageSender messagesender.MessageSenderAtTime,
	reacher AchievementUserReacher,
) *Executor {
	return &Executor{
		user:          user,
		achievement:   achievement,
		messageSender: messageSender,
		reacher:       reacher,
	}
}

func New(user userdatamanager.UserManager, achievement achievementdatamanager.AchievementManager, messageSender messagesender.MessageSenderAtTime, reacher AchievementUserReacher) executor.UserExecutor {
	return NewReachAchievementExecutor(user, achievement, messageSender, reacher)
}

func ReachAchievementMessage(achType achievementmodel.AchievementType, level int) string {
	return fmt.Sprintf("%s\nПоздравляем, вы достигли %d уровня, откройте и получите опыт и мотивацию!", achType, level)
}

func NextSendTime(now time.Time) time.Time {
	if now.Hour() < 11 {
		return time.Date(now.Year(), now.Month(), now.Day(), 11, 00, 00, 00, time.UTC)
	} else {
		return time.Date(now.Year(), now.Month(), now.Day()+1, 11, 00, 00, 00, time.UTC)
	}
}

func (e *Executor) ExecuteUser(ctx context.Context, userId int64) error {
	userData, err := e.user.UserData(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get user data", "ExecuteUser", _PROVIDER))
	}
	achievements, err := e.achievement.UserAchievements(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get user achievements", "ExecuteUser", _PROVIDER))
	}
	reachAchievements := e.reacher.ReachAchievements(ctx, userId, userData, achievements)
	err = e.achievement.ReachAchievements(ctx, userId, time.Now(), reachAchievements)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("reach achievements", "ExecuteUser", _PROVIDER))
	}
	achs := make(map[int64]struct{}, len(reachAchievements))
	sendTime := NextSendTime(time.Now())
	for _, id := range reachAchievements {
		achs[id] = struct{}{}
	}
	for _, ach := range achievements {
		if _, ok := achs[ach.ID]; !ok {
			continue
		}
		achievement := ach
		message := ReachAchievementMessage(achievement.Type, achievement.Level)
		e.messageSender.SendMessageAtTime(ctx, message, userId, sendTime)
	}
	return nil
}
