package worker

import (
	"context"
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/workers"
)

type User struct {
	ID                int64
	Motivation        int
	WelcomeMotivation int
	AbstinenceTime    time.Time
}

type UserStorage struct {
	mu sync.Mutex

	users map[int64]*User
	// map check time to user
	storage map[int64][]*User
}

func (s *UserStorage) Add(checkTime int64, user User) {
	s.mu.Lock()

	s.mu.Unlock()
}

func (s *UserStorage) SetUser()

func (s *UserStorage) Remove(userId int64)

func (s *UserStorage) Users(checkTime int64) ([]User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, ok := s.storage[checkTime]
	if !ok {
		return make([]User, 0), false
	}
	userList := make([]User, 0, len(users))

	for i := range users {
		userList = append(userList, *users[i])
	}
	return userList, true
}

type Worker struct {
	storage UserStorage
}

func New() workers.UserWorker {
	return &Worker{}
}

func (w *Worker) AddAll()

func (w *Worker) Add()

func (w *Worker) Remove()

func (w *Worker) Run(ctx context.Context) {
	seconds := time.Now().Unix()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			defer func() { seconds++ }()

		}
	}
}
