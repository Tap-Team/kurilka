package userstorage_test

import (
	"context"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/redis/go-redis/v9"
)

func saveUser(ctx context.Context, rc *redis.Client, userId int64, user *usermodel.UserData) error {
	if user == nil {
		return nil
	}
	return rc.Set(ctx, redishelper.UsersKey(userId), user, 0).Err()
}
