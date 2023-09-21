package workers

import (
	"context"

	"github.com/Tap-Team/kurilka/workers"
)

type Worker struct {
	workers []RunnableUserWorker
}

func New(ctx context.Context, workers ...RunnableUserWorker) workers.UserWorker {
	for _, w := range workers {
		go w.Run(ctx)
	}
	return &Worker{
		workers: workers,
	}
}

func (w *Worker) AddAllUsers(ctx context.Context, users []*workers.User) {
	for _, w := range w.workers {
		w.AddAllUsers(ctx, users)
	}
}

func (w *Worker) AddUser(ctx context.Context, user workers.User) {
	for _, w := range w.workers {
		w.AddUser(ctx, user)
	}
}

func (w *Worker) RemoveUser(ctx context.Context, userId int64) {
	for _, w := range w.workers {
		w.RemoveUser(ctx, userId)
	}
}
