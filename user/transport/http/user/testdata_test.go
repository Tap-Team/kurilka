package user_test

import (
	"math"
	"math/rand"
	"strings"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/random"
)

func compareUsers(u1, u2 usermodel.User) (string, bool) {
	var fields strings.Builder
	if u1.ID != u2.ID {
		fields.WriteString("ID ")
	}
	if u1.AbstinenceTime.Unix() != u2.AbstinenceTime.Unix() {
		fields.WriteString("AbstinenceTime ")
	}
	if u1.Life != u2.Life {
		fields.WriteString("Life ")
	}
	if u1.Cigarette != u2.Cigarette {
		fields.WriteString("Cigarette ")
	}
	if u1.Time != u2.Time {
		fields.WriteString("Time ")
	}
	if math.Abs(float64(u1.Money)-float64(u2.Money)) > 1 {
		fields.WriteString("Money ")
	}

	return fields.String(), fields.Len() == 0
}

func compareFriends(f1, f2 *usermodel.Friend) (string, bool) {
	var fields strings.Builder
	if f1.ID != f2.ID {
		fields.WriteString("ID ")
	}
	if f1.AbstinenceTime.Unix() != f2.AbstinenceTime.Unix() {
		fields.WriteString("AbstinenceTime ")
	}
	if f1.Life != f2.Life {
		fields.WriteString("Life ")
	}
	if f1.Cigarette != f2.Cigarette {
		fields.WriteString("Cigarette ")
	}
	if f1.Time != f2.Time {
		fields.WriteString("Time ")
	}
	if math.Abs(float64(f1.Money)-float64(f2.Money)) > 1 {
		fields.WriteString("Money ")
	}
	return fields.String(), fields.Len() == 0

}

func randomIntList(size int) []int64 {
	list := make([]int64, 0, size)
	for i := 0; i < size; i++ {
		list = append(list, rand.Int63())
	}
	return list
}

func randomFriendsList(size int) []*usermodel.Friend {
	friends := make([]*usermodel.Friend, 0, size)

	for i := 0; i < size; i++ {
		friend := random.StructTyped[usermodel.Friend]()
		friends = append(friends, &friend)
	}
	return friends
}
