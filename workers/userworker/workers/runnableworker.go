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
	isRunned bool
	executor executor.UserExecutor
	storage  *userstorage.UserStorage
}

func NewRunnableWorker(executor executor.UserExecutor, executePause time.Duration) RunnableUserWorker {
	return &runnableWorker{
		executor: executor,
		storage:  userstorage.New(),
	}
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
			go r.ExecuteUsers(ctx, users)
			seconds++
			time.Sleep(time.Second)
		}
	}
}

func (r *runnableWorker) ExecuteUsers(ctx context.Context, users []int64) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, userId := range users {
		go func(userId int64) {
			defer wg.Done()
			err := r.executor.ExecuteUser(ctx, userId)
			if errors.Is(err, usererror.ExceptionUserNotFound()) {
				r.storage.RemoveUser(userId)
			}
		}(userId)
	}
	wg.Wait()
	cancel()
}

func (r *runnableWorker) AddUser(ctx context.Context, user workers.User) {
	r.storage.AddUser(user)
}

func (r *runnableWorker) RemoveUser(ctx context.Context, userId int64) {
	r.storage.RemoveUser(userId)
}

func (r *runnableWorker) AddAllUsers(ctx context.Context, users []*workers.User) {
	r.storage.AddUsers(users)
}
