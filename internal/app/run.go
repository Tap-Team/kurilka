package app

import (
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/Tap-Team/kurilka/internal/swagger"
	userrouting "github.com/Tap-Team/kurilka/user/routing"
)

const (
	productionConfigPath = "config/production/config.yaml"
	debugConfigPath      = "config/debug/config.yaml"
)

func filePath() string {
	mode := os.Getenv("MODE")
	mode = strings.Trim(mode, "")
	mode = strings.ToUpper(mode)
	switch mode {
	case "PRODUCTION":
		return productionConfigPath
	case "DEBUG":
		return debugConfigPath
	default:
		return ""
	}
}

func Run() {
	start := time.Now()
	SetLogger()
	cnf := config.ParseFromFile(filePath())
	db := Postgres(cnf.PostgresConfig())
	rc := Redis(cnf.RedisConfig())
	vk := VK()
	router := Router()
	userrouting.SetUpRouting(userrouting.Config{
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
	})
	swagger.Swagger(router, cnf.ServerConfig())

	server := Server(router, cnf.ServerConfig())

	slog.Info("server launched", "duration", time.Since(start).String())
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed start server, %s", err)
	}
}
