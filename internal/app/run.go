package app

import (
	"log"
	"log/slog"
	"os"
	"time"

	achievementrouting "github.com/Tap-Team/kurilka/achievements/routing"
	"github.com/Tap-Team/kurilka/callback"
	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/Tap-Team/kurilka/internal/swagger"
	userrouting "github.com/Tap-Team/kurilka/user/routing"
)

func filePath() string {
	filePath := os.Getenv("CONFIG_PATH")
	return filePath
}

const (
	subscriptionCacheExpiration   = time.Hour * 24
	userCacheExpiration           = time.Hour * 24
	privacySettingCacheExpiration = time.Hour * 24
)

func Run() {
	os.Setenv("TZ", time.UTC.String())
	start := time.Now()
	SetLogger()
	cnf := config.ParseFromFile(filePath())
	db := Postgres(cnf.PostgresConfig())
	rc := Redis(cnf.RedisConfig())
	vk := VK(cnf.VKConfig())
	vkcnf := cnf.VKConfig()
	router := Router(vkcnf.AppSecretKey)

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
			CacheExpiration: userCacheExpiration,
		},
		PrivacySettingsConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: privacySettingCacheExpiration,
		},
		SubscriptionConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: subscriptionCacheExpiration,
		},
		VKConfig: struct {
			ApiVersion string
			GroupID    int64
			GroupToken string
		}{
			ApiVersion: vkcnf.ApiVersion,
			GroupID:    vkcnf.GroupID,
			GroupToken: vkcnf.GroupAccessKey,
		},
	})
	achievementrouting.SetUpAchievement(&achievementrouting.Config{
		Mux:   router,
		Redis: rc,
		DB:    db,
	})
	callback.SetUp(&callback.Config{
		Mux:   router,
		Redis: rc,
		DB:    db,
		SubscriptionConfig: struct {
			CacheExpiration time.Duration
			CostPerMonth    int
		}{
			CacheExpiration: subscriptionCacheExpiration,
			CostPerMonth:    vkcnf.SubscriptionPrice,
		},
		VKConfig: struct {
			GroupId        int64
			ConfirmKey     string
			Secret         string
			GroupAccessKey string
			ApiVersion     string
		}{
			GroupId:        vkcnf.GroupID,
			ConfirmKey:     vkcnf.CallBackConfirmKey,
			Secret:         vkcnf.CallBackSecretKey,
			GroupAccessKey: vkcnf.GroupAccessKey,
			ApiVersion:     vkcnf.ApiVersion,
		},
	})
	swagger.Swagger(router, cnf.ServerConfig())

	server := Server(router, cnf.ServerConfig())

	slog.Info("server launched", "duration", time.Since(start).String(), "host", cnf.ServerConfig().Host, "port", cnf.ServerConfig().Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed start server, %s", err)
	}
}
