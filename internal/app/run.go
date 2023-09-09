package app

import (
	"log"
	"log/slog"
	"os"
	"time"

	achievementrouting "github.com/Tap-Team/kurilka/achievements/routing"
	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/Tap-Team/kurilka/internal/swagger"
	userrouting "github.com/Tap-Team/kurilka/user/routing"
)

func filePath() string {
	filePath := os.Getenv("CONFIG_PATH")
	return filePath
}

func Run() {
	os.Setenv("TZ", time.UTC.String())
	start := time.Now()
	SetLogger()
	cnf := config.ParseFromFile(filePath())
	db := Postgres(cnf.PostgresConfig())
	rc := Redis(cnf.RedisConfig())
	vk := VK(cnf.VKConfig())
	router := Router()
	userrouting.SetUpRouting(&userrouting.Config{
		Mux:   router,
		Redis: rc,
		DB:    db,
		VK:    vk,
		UserConfig: struct {
			TrialPeriod     time.Duration
			CacheExpiration time.Duration
		}{
			TrialPeriod:     time.Hour * 24 * 5,
			CacheExpiration: time.Hour * 12,
		},
		PrivacySettingsConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: time.Hour * 12,
		},
		SubscriptionConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: time.Hour * 24,
		},
		VKConfig: struct {
			ApiVersion string
			GroupID    int64
			GroupToken string
		}{
			ApiVersion: cnf.VKConfig().ApiVersion,
			GroupID:    cnf.VKConfig().GroupID,
			GroupToken: cnf.VKConfig().GroupAccessKey,
		},
	})
	achievementrouting.SetUpAchievement(&achievementrouting.Config{
		Mux:   router,
		Redis: rc,
		DB:    db,
	})
	swagger.Swagger(router, cnf.ServerConfig())

	server := Server(router, cnf.ServerConfig())

	slog.Info("server launched", "duration", time.Since(start).String(), "host", cnf.ServerConfig().Host, "port", cnf.ServerConfig().Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed start server, %s", err)
	}
}
