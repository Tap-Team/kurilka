package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	achievementrouting "github.com/Tap-Team/kurilka/achievements/routing"
	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/middleware"
	"github.com/Tap-Team/kurilka/internal/swagger"
	"github.com/Tap-Team/kurilka/internal/userworkerinit"
	"github.com/Tap-Team/kurilka/sharestatic"
	userrouting "github.com/Tap-Team/kurilka/user/routing"
	"github.com/Tap-Team/kurilka/vote"
	"github.com/Tap-Team/kurilka/workers/userworker"
	"github.com/rs/cors"
)

func filePath() string {
	filePath := os.Getenv("CONFIG_PATH")
	return filePath
}

const (
	subscriptionCacheExpiration   = time.Hour * 24
	userCacheExpiration           = time.Hour * 24
	privacySettingCacheExpiration = time.Hour * 24
	achievementCacheExpiration    = time.Hour * 24
)

func Run() {
	ctx := context.Background()
	os.Setenv("TZ", time.UTC.String())
	start := time.Now()
	SetLogger()
	cnf := config.ParseFromFile(filePath())

	log.Print("Config Parsed: ", cnf)

	db := Postgres(cnf.PostgresConfig())
	rc := Redis(cnf.RedisConfig())
	vkcnf := cnf.VKConfig()
	router := Router()
	messageSender := MessageSender(cnf.VKConfig())
	messageScheduler := messagesender.NewMessageScheduler(ctx, messageSender)

	achievementMessageScheduler := AchievementMessageSenderScheduler(ctx, cnf.VKConfig())

	userworker := userworker.Worker(&userworker.Config{
		DB:                             db,
		Redis:                          rc,
		MessageSender:                  messageScheduler,
		AchievementMessageSenderAtTime: achievementMessageScheduler,
		VKConfig: struct {
			ApiVersion string
			GroupToken string
			GroupID    int
			AppID      int
		}{
			ApiVersion: vkcnf.ApiVersion,
			GroupToken: vkcnf.GroupAccessKey,
			GroupID:    int(vkcnf.GroupID),
			AppID:      int(vkcnf.AppID),
		},
		UserConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: userCacheExpiration,
		},
		AchievementConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: achievementCacheExpiration,
		},
	})
	userworkerinit.InitUserWorkerWorker(db, userworker)

	apiRouter := router.NewRoute().Subrouter()
	apiRouter.Use(middleware.LaunchParams(vkcnf.AppSecretKey))

	userrouting.SetUpRouting(&userrouting.Config{
		Mux:                      apiRouter,
		Redis:                    rc,
		DB:                       db,
		UserWorker:               userworker,
		AchievementMessageSender: achievementMessageScheduler,
		UserConfig: struct {
			TrialPeriod     time.Duration
			CacheExpiration time.Duration
		}{
			TrialPeriod:     time.Hour * 24 * 3,
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
			ServiceKey string
			GroupToken string
		}{
			ApiVersion: vkcnf.ApiVersion,
			GroupID:    vkcnf.GroupID,
			GroupToken: vkcnf.GroupAccessKey,
			ServiceKey: vkcnf.AppAccessKey,
		},
	})
	achievementrouting.SetUpAchievement(&achievementrouting.Config{
		Mux:           apiRouter,
		Redis:         rc,
		DB:            db,
		MessageSender: messageSender,
		AchievementConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: achievementCacheExpiration,
		},
		SubscriptionConfig: struct{ CacheExpiration time.Duration }{
			CacheExpiration: subscriptionCacheExpiration,
		},
	})
	// callback.SetUp(&callback.Config{
	// 	Mux:   router,
	// 	Redis: rc,
	// 	DB:    db,
	// 	SubscriptionConfig: struct {
	// 		CacheExpiration time.Duration
	// 		CostPerMonth    int
	// 	}{
	// 		CacheExpiration: subscriptionCacheExpiration,
	// 		CostPerMonth:    vkcnf.SubscriptionPrice,
	// 	},
	// 	VKConfig: struct {
	// 		GroupId        int64
	// 		ConfirmKey     string
	// 		Secret         string
	// 		GroupAccessKey string
	// 		ApiVersion     string
	// 	}{
	// 		GroupId:        vkcnf.GroupID,
	// 		ConfirmKey:     vkcnf.CallBackConfirmKey,
	// 		Secret:         vkcnf.CallBackSecretKey,
	// 		GroupAccessKey: vkcnf.GroupAccessKey,
	// 		ApiVersion:     vkcnf.ApiVersion,
	// 	},
	// })
	vote.SetUp(&vote.Config{
		Redis:     rc,
		DB:        db,
		ApiRouter: apiRouter,
		Mux:       router,
		SubscriptionConfig: struct{ Expiration time.Duration }{
			Expiration: subscriptionCacheExpiration,
		},
		VKConfig: struct {
			VKAppSecret     string
			VKAppServiceKey string
			Version         string
		}{
			VKAppSecret:     vkcnf.AppSecretKey,
			VKAppServiceKey: vkcnf.AppAccessKey,
			Version:         vkcnf.ApiVersion,
		},
	})
	fsCnf := cnf.FileStaticConfig()
	sharestatic.EnableShareStatic(&sharestatic.Config{
		Mux:    router,
		ApiMux: apiRouter,
		FileSystemConfig: struct {
			StoragePath    string
			FileExpiration time.Duration
		}{
			StoragePath:    fsCnf.StaticDir,
			FileExpiration: time.Second * 60 * 2,
		},
		StaticRouteConfig: struct {
			StaticUrlPrefix string
			StaticRoute     string
		}{
			StaticRoute:     fsCnf.StaticRoute,
			StaticUrlPrefix: fsCnf.RouterPath,
		},
	})

	swagger.Swagger(router, cnf.ServerConfig())

	server := Server(cors.AllowAll().Handler(router), cnf.ServerConfig())

	slog.Info("server launched", "duration", time.Since(start).String(), "host", cnf.ServerConfig().Host, "port", cnf.ServerConfig().Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed start server, %s", err)
	}
}
