package userstorage_test

import (
	"slices"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/workers"
	"github.com/Tap-Team/kurilka/workers/userworker/domain/userstorage"
	"gotest.tools/v3/assert"
)

func Test_UserStorage_AddUsers(t *testing.T) {
	storage := userstorage.New()

	cases := []struct {
		t          time.Time
		usersAdded bool
		users      []*workers.User
	}{
		{
			usersAdded: true,
			t:          time.Unix(1, 0),
		},
		{
			t: time.Unix(2, 0),
		},
	}

	for _, cs := range cases {
		var users []*workers.User
		if cs.usersAdded {
			users = generateUsers(100, func(u *workers.User) { u.AbstinenceTime = cs.t })
			storage.AddUsers(users)
		}

		u := storage.UsersByTime(cs.t.Unix())

		equal := slices.EqualFunc(users, u, func(u *workers.User, i int64) bool { return u.ID == i })
		assert.Equal(t, true, equal, "user ids not equal")
	}
}

func Test_UserStorage_AddUser(t *testing.T) {
	storage := userstorage.New()

	cases := []struct {
		users []workers.User
		t     time.Time

		expectedUsers []int64
	}{
		{
			t: time.Unix(1, 0),
			users: []workers.User{
				workers.NewUser(1, time.Unix(2, 0)),
				workers.NewUser(1, time.Unix(1, 0)),
			},
			expectedUsers: []int64{1},
		},

		{
			t: time.Unix(100, 0),
			users: []workers.User{
				workers.NewUser(100, time.Unix(100, 0)),
				workers.NewUser(200, time.Unix(100, 0)),
			},
			expectedUsers: []int64{100, 200},
		},
	}

	for _, cs := range cases {
		users := cs.users
		for _, u := range users {
			storage.AddUser(u)
		}

		u := storage.UsersByTime(cs.t.Unix())
		equal := slices.Equal(u, cs.expectedUsers)
		assert.Equal(t, true, equal, "users not equal")
	}
}

func Test_UserStorage_RemoveUser(t *testing.T) {
	storage := userstorage.New()

	cases := []struct {
		users        []workers.User
		removeUserId int64

		expectedUsers []int64

		t time.Time
	}{
		{
			t: time.Unix(1, 0),
			users: []workers.User{
				workers.NewUser(1, time.Unix(1, 0)),
				workers.NewUser(2, time.Unix(1, 0)),
				workers.NewUser(3, time.Unix(1, 0)),
				workers.NewUser(4, time.Unix(1, 0)),
				workers.NewUser(5, time.Unix(1, 0)),
				workers.NewUser(6, time.Unix(1, 0)),
			},
			removeUserId:  5,
			expectedUsers: []int64{1, 2, 3, 4, 6},
		},

		{
			t: time.Unix(2, 0),

			users: []workers.User{
				workers.NewUser(10, time.Unix(2, 0)),
				workers.NewUser(20, time.Unix(2, 0)),
				workers.NewUser(30, time.Unix(2, 0)),
				workers.NewUser(40, time.Unix(2, 0)),
				workers.NewUser(50, time.Unix(2, 0)),
				workers.NewUser(60, time.Unix(2, 0)),
			},
			removeUserId:  10,
			expectedUsers: []int64{60, 20, 30, 40, 50},
		},
	}

	for _, cs := range cases {
		for _, user := range cs.users {
			storage.AddUser(user)
		}

		storage.RemoveUser(cs.removeUserId)

		users := storage.UsersByTime(cs.t.Unix())

		equal := slices.Equal(users, cs.expectedUsers)
		assert.Equal(t, true, equal, "users not equal")
	}
}
