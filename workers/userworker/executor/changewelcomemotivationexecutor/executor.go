package changewelcomemotivationexecutor

import (
	"context"

	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/welcomemotivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor"
)

const _PROVIDER = "workers/userworker/executor/changewelcomemotivationexecutor.Executor"

type Executor struct {
	welcomeMotivation welcomemotivationdatamanager.WelcomeMotivationManager
}

func New(welcomeMotivation welcomemotivationdatamanager.WelcomeMotivationManager) executor.UserExecutor {
	return &Executor{welcomeMotivation: welcomeMotivation}
}

func (e *Executor) ExecuteUser(ctx context.Context, userId int64) error {
	motivation, err := e.welcomeMotivation.NextUserWelcomeMotivation(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get next user welcome motivation", "ExecuteUser", _PROVIDER))
	}
	err = e.welcomeMotivation.UpdateUserWelcomeMotivation(ctx, userId, motivation)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("update user welcome motivation", "ExecuteUser", _PROVIDER))
	}
	return nil
}
