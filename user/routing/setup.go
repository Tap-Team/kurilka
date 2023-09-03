package routing

import (
	"database/sql"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/Tap-Team/kurilka/user/database/postgres/achievementstorage"
	"github.com/Tap-Team/kurilka/user/database/postgres/privacysettingstorage"
	achievementcache "github.com/Tap-Team/kurilka/user/database/redis/achievementstorage"
	privacysettingcache "github.com/Tap-Team/kurilka/user/database/redis/privacysettingstorage"
	"github.com/Tap-Team/kurilka/user/database/vk/friendsstorage"
	"github.com/Tap-Team/kurilka/user/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/datamanager/triggerdatamanager"
	"github.com/Tap-Team/kurilka/user/usecase/userusecase"

	"github.com/Tap-Team/kurilka/user/database/postgres/resetrecoveruserstorage"
	"github.com/Tap-Team/kurilka/user/database/postgres/triggerstorage"
	"github.com/Tap-Team/kurilka/user/database/postgres/userstorage"
	triggercache "github.com/Tap-Team/kurilka/user/database/redis/triggerstorage"
	usercache "github.com/Tap-Team/kurilka/user/database/redis/userstorage"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Mux        *mux.Router
	Redis      *redis.Client
	DB         *sql.DB
	VK         *api.VK
	UserConfig struct {
		TrialPeriod     time.Duration
		CacheExpiration time.Duration
	}
}

type setUpper struct {
	config   Config
	managers struct {
		trigger         triggerdatamanager.TriggerManager
		user            userdatamanager.UserManager
		privacySettings privacysettingdatamanager.PrivacySettingManager
		achievement     achievementdatamanager.AchievementManager
	}
	usecases struct {
		user userusecase.UserUseCase
	}
}

func NewSetUpper(config Config) *setUpper {
	return &setUpper{config: config}
}

func (s *setUpper) Config() Config {
	return s.config
}

func (s *setUpper) TriggerManager() triggerdatamanager.TriggerManager {
	if s.managers.trigger != nil {
		return s.managers.trigger
	}
	cache := triggercache.New(s.config.Redis, s.config.UserConfig.CacheExpiration)
	storage := triggerstorage.New(s.config.DB)
	s.managers.trigger = triggerdatamanager.NewTriggerManager(storage, cache)
	return s.managers.trigger
}

func (s *setUpper) UserManager() userdatamanager.UserManager {
	if s.managers.user != nil {
		return s.managers.user
	}
	cache := usercache.New(s.config.Redis, s.config.UserConfig.CacheExpiration)
	storage := userstorage.New(s.config.DB, s.config.UserConfig.TrialPeriod)
	recoverResetter := resetrecoveruserstorage.New(s.config.DB)
	saver := userdatamanager.NewUserSaver(storage, recoverResetter)
	s.managers.user = userdatamanager.NewUserManager(
		recoverResetter,
		storage,
		cache,
		saver,
	)
	return s.managers.user
}

func (s *setUpper) AchievementManager() achievementdatamanager.AchievementManager {
	if s.managers.achievement != nil {
		return s.managers.achievement
	}
	storage := achievementstorage.New(s.config.DB)
	cache := achievementcache.New(s.config.Redis)
	s.managers.achievement = achievementdatamanager.NewAchievementManager(storage, cache)
	return s.managers.achievement
}

func (s *setUpper) PrivacySettingManager() privacysettingdatamanager.PrivacySettingManager {
	if s.managers.privacySettings != nil {
		return s.managers.privacySettings
	}
	storage := privacysettingstorage.New(s.config.DB)
	cache := privacysettingcache.New(s.config.Redis, s.config.UserConfig.CacheExpiration)
	s.managers.privacySettings = privacysettingdatamanager.NewPrivacyManager(storage, cache)
	return s.managers.privacySettings
}

func (s *setUpper) UserUseCase() userusecase.UserUseCase {
	if s.usecases.user != nil {
		return s.usecases.user
	}
	friendStorage := friendsstorage.New(s.config.VK)
	s.usecases.user = userusecase.NewUser(
		userusecase.NewUserFriendsProvider(friendStorage, s.UserManager()),
		s.UserManager(),
		s.PrivacySettingManager(),
		s.AchievementManager(),
		userusecase.NewFriendProvider(s.AchievementManager(), s.UserManager(), s.PrivacySettingManager()),
	)
	return s.usecases.user
}

func SetUpRouting(config Config) {
	setUpper := NewSetUpper(config)

	TriggerRouting(setUpper)
	PrivacySettingRouting(setUpper)
	UserRouting(setUpper)
}
