package friendsstorage

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/api"
)

type Storage struct {
	vk *api.VK
}

func New(vk *api.VK) *Storage {
	return &Storage{vk: vk}
}

type UserFriendsParams api.Params

func (u UserFriendsParams) SetUser(userId int64) {
	u["user_id"] = userId
}

func (s *Storage) Friends(ctx context.Context, userId int64) []int64 {
	userParams := make(UserFriendsParams)
	userParams.SetUser(userId)
	r, err := s.vk.FriendsGet(api.Params(userParams))
	if err != nil {
		return make([]int64, 0)
	}
	friends := make([]int64, 0, r.Count)
	for _, id := range r.Items {
		friends = append(friends, int64(id))
	}
	return friends
}
