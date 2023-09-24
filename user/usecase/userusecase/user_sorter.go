package userusecase

import (
	"sort"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
)

type FriendsSorter []*usermodel.Friend

func (s FriendsSorter) Len() int {
	return len(s)
}

func (s FriendsSorter) Swap(x, y int) {
	s[x], s[y] = s[y], s[x]
}

func (s FriendsSorter) Less(x, y int) bool {
	return s[x].ID < s[y].ID
}

func (s FriendsSorter) Sort() {
	sort.Sort(s)
}
