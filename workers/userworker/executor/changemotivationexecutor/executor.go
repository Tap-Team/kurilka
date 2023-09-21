package changemotivationexecutor

import (
	"context"
	"time"

	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/motivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor"
)

const _PROVIDER = "workers/userworker/executor/changemotivationexecutor.Executor"

type Executor struct {
	motivation    motivationdatamanager.MotivationManager
	messageSender messagesender.MessageSenderAtTime
}

func New(motivation motivationdatamanager.MotivationManager, messageSender messagesender.MessageSenderAtTime) executor.UserExecutor {
	return &Executor{
		motivation:    motivation,
		messageSender: messageSender,
	}
}

func NextSendTime(now time.Time) time.Time {
	if now.Hour() < 11 {
		return time.Date(now.Year(), now.Month(), now.Day(), 11, 00, 00, 00, time.UTC)
	} else {
		return time.Date(now.Year(), now.Month(), now.Day()+1, 11, 00, 00, 00, time.UTC)
	}
}

func (e *Executor) ExecuteUser(ctx context.Context, userId int64) error {
	motivation, err := e.motivation.NextUserMotivation(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("get next user motivation", "ExecuteUser", _PROVIDER))
	}
	err = e.motivation.UpdateUserMotivation(ctx, userId, motivation)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("update user motivation", "ExecuteUser", _PROVIDER))
	}
	e.messageSender.SendMessageAtTime(ctx, motivation.Motivation, userId, NextSendTime(time.Now()))
	return nil
}
