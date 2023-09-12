package workers

import "context"

type UserWorker interface {
	AddUser(ctx context.Context, userId int64) error
	RemoveUser(ctx context.Context, userId int64) error
}
