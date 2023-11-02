package routing

import (
	"database/sql"
	"time"

	"github.com/Tap-Team/kurilka/achievements/database/postgres/achievementstorage"
	"github.com/Tap-Team/kurilka/achievements/database/postgres/subscriptionstorage"
	"github.com/Tap-Team/kurilka/achievements/database/postgres/userstorage"
	achievementcache "github.com/Tap-Team/kurilka/achievements/database/redis/achievementstorage"
	subscriptioncache "github.com/Tap-Team/kurilka/achievements/database/redis/subscriptionstorage"
	usercache "github.com/Tap-Team/kurilka/achievements/database/redis/userstorage"
	"github.com/Tap-Team/kurilka/internal/messagesender"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/subscriptiondatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Mux               *mux.Router
	Redis             *redis.Client
	DB                *sql.DB
	MessageSender     messagesender.MessageSender
	AchievementConfig struct {
		CacheExpiration time.Duration
	}
	SubscriptionConfig struct {
		CacheExpiration time.Duration
	}
}

type setUpper struct {
	cnf                *Config
	achievementStorage *achievementstorage.Storage
	managers           struct {
		user         userdatamanager.UserManager
		achievement  achievementdatamanager.AchievementManager
		subscription subscriptiondatamanager.SubscriptionManager
	}
	useCases struct {
		achievement achievementusecase.AchievementUseCase
	}
}

func (s *setUpper) AchievementStorage() *achievementstorage.Storage {
	if s.achievementStorage != nil {
		return s.achievementStorage
	}
	s.achievementStorage = achievementstorage.New(s.cnf.DB)
	return s.achievementStorage
}

func (s *setUpper) UserManager() userdatamanager.UserManager {
	if s.managers.user != nil {
		return s.managers.user
	}
	cache := usercache.New(s.cnf.Redis)
	storage := userstorage.New(s.cnf.DB)
	s.managers.user = userdatamanager.New(storage, cache)
	return s.managers.user
}

func (s *setUpper) AchievementManager() achievementdatamanager.AchievementManager {
	if s.managers.achievement != nil {
		return s.managers.achievement
	}
	cache := achievementcache.New(s.cnf.Redis, s.cnf.AchievementConfig.CacheExpiration)
	storage := s.AchievementStorage()
	s.managers.achievement = achievementdatamanager.NewAchievementManager(storage, achievementdatamanager.NewCacheWrapper(cache))
	return s.managers.achievement
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

func (s *setUpper) AchievementUseCase() achievementusecase.AchievementUseCase {
	if s.useCases.achievement != nil {
		return s.useCases.achievement
	}
	user := s.UserManager()
	achievement := s.AchievementManager()
	s.useCases.achievement = achievementusecase.New(achievement, user, s.AchievementStorage(), s.cnf.MessageSender, s.SubscriptionManager())
	return s.useCases.achievement
}

func NewSetUpper(cnf *Config) *setUpper {
	return &setUpper{
		cnf: cnf,
	}
}

func SetUpAchievement(cnf *Config) {
	setUpper := NewSetUpper(cnf)
	AchievementRouting(setUpper)
}
