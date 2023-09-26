package workers_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/workers"
	"github.com/Tap-Team/kurilka/workers/userworker/domain/userstorage"
	"github.com/Tap-Team/kurilka/workers/userworker/executor"
	userworkers "github.com/Tap-Team/kurilka/workers/userworker/workers"
	"github.com/golang/mock/gomock"
)

func Test_RunnableWorker_Add(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	timeCounter := userworkers.NewMockUserExecuteTimeCounter(ctrl)
	executor := executor.NewMockUserExecutor(ctrl)
	storage := userstorage.NewMockUserStorage(ctrl)

	worker := userworkers.NewRunnableWorkerStruct(executor, time.Hour, storage, timeCounter, time.Second)

	cases := []struct {
		user workers.User
	}{
		{
			user: workers.NewUser(rand.Int63(), time.Now()),
		},
	}

	for _, cs := range cases {
		timeCounterTime := time.Time{}
		timeCounter.EXPECT().CountUserExecuteTime(gomock.Any(), cs.user.AbstinenceTime).Return(timeCounterTime).Times(1)
		storage.EXPECT().AddUser(workers.NewUser(cs.user.ID, timeCounterTime)).Times(1)

		worker.AddUser(ctx, cs.user)
	}
}

func Test_RunnableWorker_AddAll(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	timeCounter := userworkers.NewMockUserExecuteTimeCounter(ctrl)
	executor := executor.NewMockUserExecutor(ctrl)
	storage := userstorage.NewMockUserStorage(ctrl)

	worker := userworkers.NewRunnableWorkerStruct(executor, time.Hour, storage, timeCounter, time.Second)

	cases := []struct {
		users []*workers.User
	}{
		{
			users: []*workers.User{
				{ID: rand.Int63(), AbstinenceTime: time.Now()},
				{ID: rand.Int63(), AbstinenceTime: time.Now()},
				{ID: rand.Int63()},
			},
		},
	}

	for _, cs := range cases {
		for i := range cs.users {
			timeCounterTime := time.Time{}
			timeCounter.EXPECT().CountUserExecuteTime(gomock.Any(), cs.users[i].AbstinenceTime).Return(timeCounterTime).Times(1)
		}
		storage.EXPECT().AddUsers(cs.users)
		worker.AddAllUsers(ctx, cs.users)
	}
}

func Test_RunnableWorker_Remove(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	timeCounter := userworkers.NewMockUserExecuteTimeCounter(ctrl)
	executor := executor.NewMockUserExecutor(ctrl)
	storage := userstorage.NewMockUserStorage(ctrl)

	worker := userworkers.NewRunnableWorkerStruct(executor, time.Hour, storage, timeCounter, time.Second)

	cases := []struct {
		userId int64
	}{
		{
			userId: rand.Int63(),
		},
		{
			userId: rand.Int63(),
		},
		{
			userId: rand.Int63(),
		},
		{
			userId: rand.Int63(),
		},
		{
			userId: rand.Int63(),
		},
	}

	for _, cs := range cases {
		storage.EXPECT().RemoveUser(cs.userId)
		worker.RemoveUser(ctx, cs.userId)
	}
}

type fakeExecutor struct {
	t     *testing.T
	mu    *sync.Mutex
	calls map[int64][]error
}

func NewFakeExecutor(t *testing.T) *fakeExecutor {
	var mu sync.Mutex
	calls := make(map[int64][]error)

	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		nonCalls := 0
		for _, userCalls := range calls {
			nonCalls += len(userCalls)
		}
		if len(calls) != 0 {
			t.Fatalf("no calls, %d", nonCalls)
		}
	})
	return &fakeExecutor{
		t:     t,
		mu:    &mu,
		calls: calls,
	}
}

func (e *fakeExecutor) RegisterCall(userId int64, err error) {
	e.mu.Lock()
	e.calls[userId] = append(e.calls[userId], err)
	e.mu.Unlock()
}
func (e *fakeExecutor) ExecuteUser(ctx context.Context, userId int64) error {
	e.mu.Lock()
	errs, ok := e.calls[userId]
	if !ok {
		e.t.Fatalf("unexpected user id, %d", userId)
	}
	err := errs[0]
	e.calls[userId] = e.calls[userId][1:]
	if len(e.calls[userId]) == 0 {
		delete(e.calls, userId)
	}
	e.mu.Unlock()
	return err
}

func Test_RunnableWorker_Run(t *testing.T) {
	ctx := context.Background()
	executor := NewFakeExecutor(t)

	executePause := time.Second
	tickTime := time.Millisecond * 10
	worker := userworkers.NewRunnableWorkerStruct(executor, executePause, userstorage.New(), userworkers.NewUserTimeExecutor(executePause), tickTime)

	now := time.Now()
	users := []workers.User{
		workers.NewUser(1, now),
		workers.NewUser(2, now),
	}

	count := 5

	for _, u := range users {
		for i := 0; i < count; i++ {
			executor.RegisterCall(u.ID, nil)
		}
		worker.AddUser(ctx, u)
	}

	ctx, cancel := context.WithTimeout(ctx, tickTime*6)
	defer cancel()
	worker.Run(ctx)
}
