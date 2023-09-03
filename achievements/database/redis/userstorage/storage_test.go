package userstorage_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Tap-Team/kurilka/achievements/database/redis/userstorage"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/redis/go-redis/v9"
)

var (
	rc      *redis.Client
	storage *userstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	redisClient, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed start redis container, %s", err)
	}
	defer term(ctx)
	rc = redisClient
	storage = userstorage.New(rc)
	os.Exit(m.Run())
}

func TestUpdateUserLevel(t *testing.T) {

}
