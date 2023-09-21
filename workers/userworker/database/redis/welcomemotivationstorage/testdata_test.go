package welcomemotivationstorage_test

import (
	"context"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/redis/go-redis/v9"
)

func saveUser(rc *redis.Client, userId int64, welcomeMotivation string) error {
	var user usermodel.UserData
	user.WelcomeMotivation = welcomeMotivation
	err := rc.Set(context.Background(), redishelper.UsersKey(userId), user, 0).Err()
	return err
}

func userMotivation(rc *redis.Client, userId int64) (string, error) {
	var user usermodel.UserData
	err := rc.Get(context.Background(), redishelper.UsersKey(userId)).Scan(&user)
	return user.WelcomeMotivation, err
}
