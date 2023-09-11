package callback

import (
	"context"
	"database/sql"
	"time"

	"github.com/Tap-Team/kurilka/callback/database/postgres/subscriptionstorage"
	subscriptioncache "github.com/Tap-Team/kurilka/callback/database/redis/subscriptionstorage"
	"github.com/Tap-Team/kurilka/callback/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/callback/handler"
	"github.com/Tap-Team/kurilka/callback/usecase/subscriptionusecase"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Mux                *mux.Router
	DB                 *sql.DB
	Redis              *redis.Client
	SubscriptionConfig struct {
		CacheExpiration time.Duration
		CostPerMonth    int
	}
	VKConfig struct {
		GroupId        int64
		ConfirmKey     string
		Secret         string
		GroupAccessKey string
		ApiVersion     string
	}
}

type setUpper struct {
	cnf      *Config
	managers struct {
		subscription subscriptiondatamanager.SubscriptionManager
	}
	useCases struct {
		subscription subscriptionusecase.UseCase
	}
}

func NewSetUpper(cnf *Config) *setUpper {
	return &setUpper{cnf: cnf}
}

func (s *setUpper) SubscriptionManager() subscriptiondatamanager.SubscriptionManager {
	if s.managers.subscription != nil {
		return s.managers.subscription
	}
	storage := subscriptionstorage.New(s.cnf.DB)
	cache := subscriptioncache.New(s.cnf.Redis, s.cnf.SubscriptionConfig.CacheExpiration)
	s.managers.subscription = subscriptiondatamanager.New(storage, cache)
	return s.managers.subscription
}

func (s *setUpper) SubscriptionUseCase() subscriptionusecase.UseCase {
	if s.useCases.subscription != nil {
		return s.useCases.subscription
	}
	s.useCases.subscription = subscriptionusecase.New(s.SubscriptionManager(), s.cnf.SubscriptionConfig.CostPerMonth)
	return s.useCases.subscription
}

func SetUp(cnf *Config) {
	ctx := context.Background()

	setUpper := NewSetUpper(cnf)

	vk := cnf.VKConfig

	handler := handler.New(
		vk.ConfirmKey,
		vk.GroupId,
		vk.Secret,
		setUpper.SubscriptionUseCase(),
	)
	cnf.Mux.Handle("/vk/callback", handler.CallBackHandler(ctx))
}
