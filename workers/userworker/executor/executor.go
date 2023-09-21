package executor

import "context"

//go:generate mockgen -source executor.go -destination mocks.go -package executor

type UserExecutor interface {
	ExecuteUser(ctx context.Context, userId int64) error
}
