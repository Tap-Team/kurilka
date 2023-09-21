package userstorage

import (
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/workers"
)

type UserStorage struct {
	mu sync.Mutex
	// map user id to own check time
	usersCheckTime map[int64]int64
	// map next check time to user ids
	users map[int64][]int64
}

func New() *UserStorage {
	return &UserStorage{
		usersCheckTime: make(map[int64]int64),
		users:          make(map[int64][]int64),
	}
}

func (s *UserStorage) add(user workers.User) {
	checkTime := user.AbstinenceTime.Unix()
	if _, ok := s.usersCheckTime[user.ID]; ok {
		s.remove(user.ID, checkTime)
	}
	s.usersCheckTime[user.ID] = checkTime
	s.users[checkTime] = append(s.users[checkTime], user.ID)
}

func (s *UserStorage) AddUser(user workers.User) {
	s.mu.Lock()
	s.add(user)
	s.mu.Unlock()
}

func (s *UserStorage) AddUsers(users []*workers.User) {
	s.mu.Lock()
	for _, user := range users {
		s.add(*user)
	}
	s.mu.Unlock()
}

func (s *UserStorage) remove(userId int64, checkTime int64) {
	delete(s.usersCheckTime, userId)
	for index, id := range s.users[checkTime] {
		if userId == id {
			l := len(s.users[checkTime])
			s.users[checkTime][index] = s.users[checkTime][l-1]
			s.users[checkTime] = s.users[checkTime][:l-1]
		}
	}
}

func (s *UserStorage) RemoveUser(userId int64) {
	s.mu.Lock()
	checkTime := s.usersCheckTime[userId]
	s.remove(userId, checkTime)
	s.mu.Unlock()
}

func (s *UserStorage) updateUserTime(userId int64, t int64) {
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

func (s *UserStorage) UpdateUserTime(userId int64, t time.Time) {
	s.mu.Lock()
	s.updateUserTime(userId, t.Unix())
	s.mu.Unlock()
}

func (s *UserStorage) UsersByTime(t int64) []int64 {
	s.mu.Lock()
	users := make([]int64, len(s.users[t]))
	copy(users, s.users[t])
	s.mu.Unlock()
	return users
}
