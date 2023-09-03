package triggerstorage_test

import (
	"context"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/redis/go-redis/v9"
)

func setUser(rc *redis.Client, userId int64, user *usermodel.UserData) error {
	return rc.Set(context.Background(), redishelper.UsersKey(userId), user, 0).Err()
}
