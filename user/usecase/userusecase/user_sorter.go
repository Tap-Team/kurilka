package userusecase

import (
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
)

type SortFriendsByIdAsc []*usermodel.Friend

func (s SortFriendsByIdAsc) Len() int {
	return len(s)
}

func (s SortFriendsByIdAsc) Swap(x, y int) {
	s[x], s[y] = s[y], s[x]
}

func (s SortFriendsByIdAsc) Less(x, y int) bool {
	return s[x].ID < s[y].ID
}
