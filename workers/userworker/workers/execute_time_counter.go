package workers

import (
	"time"
)

//go:generate mockgen -source execute_time_counter.go -destination execute_time_counter_mocks.go -package workers

type UserExecuteTimeCounter interface {
	CountUserExecuteTime(now, t time.Time) time.Time
}

type userExecuteTimeCounter struct {
	executePause time.Duration
}

func NewUserTimeExecutor(executePause time.Duration) UserExecuteTimeCounter {
	return &userExecuteTimeCounter{executePause: executePause}
}

/*
example

first case
now = 01.09.2023 15:00
userAbstinenceTime = 18.08.2023 16:00

function returns 01.09.2023 16:00 // return userAbstinenceTime without modify

second case
executePause = 1 day
now = 02.09.2023 17:00
userAbstinenceTime = 18.08.2023 16:00

function returns 03.09.2023 16:00 // add 1 day
*/
func (c *userExecuteTimeCounter) CountUserExecuteTime(now, t time.Time) time.Time {
	if t.After(now) {
		return t
	}
	sub := int64(now.Sub(t) / c.executePause)
	periodCount := sub + 1
	return t.Add(time.Duration(periodCount) * c.executePause)
}
