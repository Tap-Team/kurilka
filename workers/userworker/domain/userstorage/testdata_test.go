package userstorage_test

import (
	"math/rand"
	"time"

	"github.com/Tap-Team/kurilka/workers"
)

type userOpt func(*workers.User)

func generateUser(opts ...userOpt) workers.User {
	user := workers.NewUser(rand.Int63(), time.Unix(int64(rand.Int31()), 0))
	for _, opt := range opts {
		opt(&user)
	}
	return user
}

func generateUsers(size int, opts ...userOpt) []*workers.User {
	users := make([]*workers.User, 0, size)
	for i := 0; i < size; i++ {
		user := generateUser(opts...)
		users = append(users, &user)
	}
	return users
}

func compareUsers(u1, u2 *workers.User) bool {
	if u1.AbstinenceTime != u2.AbstinenceTime {
		return false
	}
	if u1.ID != u2.ID {
		return false
	}
	return true
}
