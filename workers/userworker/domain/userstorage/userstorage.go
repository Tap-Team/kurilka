package userstorage

import (
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/workers"
)

//go:generate mockgen -source userstorage.go -destination mocks.go -package userstorage

type UserStorage interface {
	AddUser(user workers.User)
	AddUsers(users []*workers.User)
	RemoveUser(userId int64)
	UpdateUserTime(userId int64, t time.Time)
	UsersByTime(t int64) []int64
}

func New() UserStorage {
	return &userStorage{
		usersCheckTime: make(map[int64]int64),
		users:          make(map[int64][]int64),
	}
}

type userStorage struct {
	mu sync.Mutex
	// map user id to own check time
	usersCheckTime map[int64]int64
	// map next check time to user ids
	users map[int64][]int64
}

func (s *userStorage) add(user workers.User) {
	checkTime := user.AbstinenceTime.Unix()
	if _, ok := s.usersCheckTime[user.ID]; ok {
		s.remove(user.ID, checkTime)
	}
	s.usersCheckTime[user.ID] = checkTime
	s.users[checkTime] = append(s.users[checkTime], user.ID)
}

func (s *userStorage) AddUser(user workers.User) {
	s.mu.Lock()
	s.add(user)
	s.mu.Unlock()
}

func (s *userStorage) AddUsers(users []*workers.User) {
	s.mu.Lock()
	for _, user := range users {
		s.add(*user)
	}
	s.mu.Unlock()
}

func (s *userStorage) remove(userId int64, checkTime int64) {
	delete(s.usersCheckTime, userId)
	for index, id := range s.users[checkTime] {
		if userId == id {
			l := len(s.users[checkTime])
			s.users[checkTime][index] = s.users[checkTime][l-1]
			s.users[checkTime] = s.users[checkTime][:l-1]
		}
	}
}

func (s *userStorage) RemoveUser(userId int64) {
	s.mu.Lock()
	checkTime := s.usersCheckTime[userId]
	s.remove(userId, checkTime)
	s.mu.Unlock()
}

func (s *userStorage) updateUserTime(userId int64, t int64) {
	checkTime := s.usersCheckTime[userId]
	for index, id := range s.users[checkTime] {
		if userId == id {
			l := len(s.users[checkTime])
			s.users[checkTime][index] = s.users[checkTime][l-1]
			s.users[checkTime] = s.users[checkTime][:l-1]
		}
	}
	checkTime = t
	s.users[checkTime] = append(s.users[checkTime], userId)
	s.usersCheckTime[userId] = checkTime
}

func (s *userStorage) UpdateUserTime(userId int64, t time.Time) {
	s.mu.Lock()
	s.updateUserTime(userId, t.Unix())
	s.mu.Unlock()
}

func (s *userStorage) UsersByTime(t int64) []int64 {
	s.mu.Lock()
	users := make([]int64, len(s.users[t]))
	copy(users, s.users[t])
	s.mu.Unlock()
	return users
}
