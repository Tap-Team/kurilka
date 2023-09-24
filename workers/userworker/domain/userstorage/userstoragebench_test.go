package userstorage_test

import (
	"sync"
	"testing"

	"github.com/Tap-Team/kurilka/workers"
	"github.com/Tap-Team/kurilka/workers/userworker/domain/userstorage"
)

func Benchmark_UserStorage_AddSingleByOneAsync(b *testing.B) {
	storage := userstorage.New()
	var wg sync.WaitGroup
	count := 10
	wg.Add(count)
	b.StartTimer()
	for i := 0; i < count; i++ {
		go func() {
			for i := 0; i < b.N/count; i++ {
				var user workers.User
				storage.AddUser(user)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	b.StopTimer()
}

func Benchmark_UserStorage_AddSingleByList(b *testing.B) {
	storage := userstorage.New()
	b.StartTimer()
	users := make([]*workers.User, 0)
	for i := 0; i < b.N; i++ {
		var user workers.User
		users = append(users, &user)
	}
	for _, u := range users {
		storage.AddUser(*u)
	}
	b.StopTimer()
}

func Benchmark_UserStorageAddMany(b *testing.B) {
	storage := userstorage.New()
	b.StartTimer()
	users := make([]*workers.User, 0)
	for i := 0; i < b.N; i++ {
		var user workers.User
		users = append(users, &user)
	}
	storage.AddUsers(users)
	b.StopTimer()
}
