package userusecase

import (
	"context"
	"sort"

	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
)

func NewUserFriendsProvider(friends UserFriendsProvider, user userdatamanager.UserManager) UserFriendsProvider {
	return &userFriendsProvider{friends: friends, user: user}
}

type userFriendsProvider struct {
	friends UserFriendsProvider
	user    userdatamanager.UserManager
}

func (f *userFriendsProvider) Friends(ctx context.Context, userId int64) []int64 {
	friends := f.user.FilterExists(ctx, f.friends.Friends(ctx, userId))
	sort.Slice(friends, func(i, j int) bool {
		return friends[i] < friends[j]
	})
	return friends
}
