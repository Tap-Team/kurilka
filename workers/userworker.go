package workers

import (
	"context"
	"time"
)

//go:generate mockgen -source userworker.go -destination mocks.go -package workers

type User struct {
	AbstinenceTime time.Time
	ID             int64
}

func NewUser(id int64, t time.Time) User {
	return User{
		AbstinenceTime: t,
		ID:             id,
	}
}

type UserWorker interface {
	AddAllUsers(ctx context.Context, users []*User)
	AddUser(ctx context.Context, user User)
	RemoveUser(ctx context.Context, userId int64)
}
