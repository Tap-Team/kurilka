package routing

import (
	"database/sql"

	"github.com/Tap-Team/kurilka/achievements/database/postgres/achievementstorage"
	"github.com/Tap-Team/kurilka/achievements/database/postgres/userstorage"
	achievementcache "github.com/Tap-Team/kurilka/achievements/database/redis/achievementstorage"
	usercache "github.com/Tap-Team/kurilka/achievements/database/redis/userstorage"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Mux   *mux.Router
	Redis *redis.Client
	DB    *sql.DB
}

type setUpper struct {
	cnf      *Config
	managers struct {
		user        userdatamanager.UserManager
		achievement achievementdatamanager.AchievementManager
	}
	useCases struct {
		achievement achievementusecase.AchievementUseCase
	}
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
	cache := achievementcache.New(s.cnf.Redis)
	storage := achievementstorage.New(s.cnf.DB)
	s.managers.achievement = achievementdatamanager.NewAchievementManager(storage, achievementdatamanager.NewCacheWrapper(cache))
	return s.managers.achievement
}

func (s *setUpper) AchievementUseCase() achievementusecase.AchievementUseCase {
	if s.useCases.achievement != nil {
		return s.useCases.achievement
	}
	user := s.UserManager()
	achievement := s.AchievementManager()
	s.useCases.achievement = achievementusecase.New(achievement, user)
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
