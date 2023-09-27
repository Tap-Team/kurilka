package workers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/workers"
	"github.com/Tap-Team/kurilka/workers/userworker/domain/userstorage"
	"github.com/Tap-Team/kurilka/workers/userworker/executor"
)

type RunnableUserWorker interface {
	workers.UserWorker
	Run(ctx context.Context)
}

type runnableWorker struct {
	isRunned     bool
	executePause time.Duration
	tickTime     time.Duration

	executor               executor.UserExecutor
	storage                userstorage.UserStorage
	userExecuteTimeCounter UserExecuteTimeCounter
}

func NewRunnableWorkerStruct(
	executor executor.UserExecutor,
	executePause time.Duration,
	userStorage userstorage.UserStorage,
	userTimeExecutor UserExecuteTimeCounter,
	tickTime time.Duration,
) *runnableWorker {
	return &runnableWorker{
		executor:               executor,
		storage:                userStorage,
		executePause:           executePause,
		userExecuteTimeCounter: userTimeExecutor,
		tickTime:               tickTime,
	}
}

func NewRunnableWorker(executor executor.UserExecutor, executePause time.Duration) RunnableUserWorker {
	return NewRunnableWorkerStruct(
		executor,
		executePause,
		userstorage.New(),
		NewUserTimeExecutor(executePause),
		time.Second,
	)
}

func (r *runnableWorker) Run(ctx context.Context) {
	if r.isRunned {
		return
	}
	r.isRunned = true
	defer func() { r.isRunned = false }()
	seconds := time.Now().Unix()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			users := r.storage.UsersByTime(seconds)
			go r.ExecuteUsers(ctx, seconds, users)
			seconds++
			time.Sleep(r.tickTime)
		}
	}
}

func (r *runnableWorker) ExecuteUsers(ctx context.Context, seconds int64, users []int64) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	nextExecuteTime := time.Unix(seconds, 0).Add(r.executePause)

	var wg sync.WaitGroup
	wg.Add(len(users))

	for _, userId := range users {
		go func(userId int64) {
			defer wg.Done()
			err := r.executor.ExecuteUser(ctx, userId)
			if errors.Is(err, usererror.ExceptionUserNotFound()) {
				r.storage.RemoveUser(userId)
			}
			r.storage.UpdateUserTime(userId, nextExecuteTime)
		}(userId)
	}
	wg.Wait()
	cancel()
}

func (r *runnableWorker) AddUser(ctx context.Context, user workers.User) {
	now := time.Now()
	user.AbstinenceTime = r.userExecuteTimeCounter.CountUserExecuteTime(now, user.AbstinenceTime)
	r.storage.AddUser(user)
}

func (r *runnableWorker) RemoveUser(ctx context.Context, userId int64) {
	r.storage.RemoveUser(userId)
}

func (r *runnableWorker) AddAllUsers(ctx context.Context, users []*workers.User) {
	now := time.Now()
	for i := range users {
		users[i].AbstinenceTime = r.userExecuteTimeCounter.CountUserExecuteTime(now, users[i].AbstinenceTime)
	}
	r.storage.AddUsers(users)
}
