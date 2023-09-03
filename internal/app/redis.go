package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/redis/go-redis/v9"
)

func Redis(cnf config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		Password: cnf.Password,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("failed ping redis, %s", err)
	}
	return client
}
