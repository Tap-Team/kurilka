package userworker

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/workers"
	"github.com/Tap-Team/kurilka/workers/userworker/achievementmessagesender"
	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/achievementstorage"
	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/motivationstorage"
	achievementcache "github.com/Tap-Team/kurilka/workers/userworker/database/redis/achievementstorage"
	motivationcache "github.com/Tap-Team/kurilka/workers/userworker/database/redis/motivationstorage"

	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/userstorage"
	usercache "github.com/Tap-Team/kurilka/workers/userworker/database/redis/userstorage"

	"github.com/Tap-Team/kurilka/workers/userworker/database/postgres/welcomemotivationstorage"
	welcomemotivationcache "github.com/Tap-Team/kurilka/workers/userworker/database/redis/welcomemotivationstorage"
	"github.com/Tap-Team/kurilka/workers/userworker/executor/changemotivationexecutor"
	"github.com/Tap-Team/kurilka/workers/userworker/executor/changewelcomemotivationexecutor"
	"github.com/Tap-Team/kurilka/workers/userworker/executor/reachachievementexecutor"

	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/motivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/welcomemotivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor"
	userworkers "github.com/Tap-Team/kurilka/workers/userworker/workers"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB            *sql.DB
	Redis         *redis.Client
	MessageSender messagesender.MessageSenderAtTime
	VKConfig      struct {
		ApiVersion string
		GroupToken string
		GroupID    int
		AppID      int
	}
	UserConfig struct {
		CacheExpiration time.Duration
	}
	AchievementConfig struct {
		CacheExpiration time.Duration
	}
}

type setUpper struct {
	cnf                            *Config
	achievementMessageSender       achievementmessagesender.AchievementMessageSender
	achievementMessageSenderAtTime achievementmessagesender.AchievementMessageSenderAtTime
	managers                       struct {
		achievement       achievementdatamanager.AchievementManager
		user              userdatamanager.UserManager
		motivation        motivationdatamanager.MotivationManager
		welcomeMotivation welcomemotivationdatamanager.WelcomeMotivationManager
	}
	executors struct {
		achievement       executor.UserExecutor
		motivation        executor.UserExecutor
		welcomeMotivation executor.UserExecutor
	}
	workers struct {
		achievement       userworkers.RunnableUserWorker
		motivation        userworkers.RunnableUserWorker
		welcomeMotivation userworkers.RunnableUserWorker
	}
}

func NewSetUpper(cnf *Config) *setUpper {
	return &setUpper{cnf: cnf}
}

func (s *setUpper) WelcomeMotivationManager() welcomemotivationdatamanager.WelcomeMotivationManager {
	if s.managers.welcomeMotivation != nil {
		return s.managers.welcomeMotivation
	}
	storage := welcomemotivationstorage.New(s.cnf.DB)
	cache := welcomemotivationcache.New(s.cnf.Redis, s.cnf.UserConfig.CacheExpiration)
	s.managers.welcomeMotivation = welcomemotivationdatamanager.New(storage, cache)
	return s.managers.welcomeMotivation
}

func (s *setUpper) AchievementManager() achievementdatamanager.AchievementManager {
	if s.managers.achievement != nil {
		return s.managers.achievement
	}
	storage := achievementstorage.New(s.cnf.DB)
	cache := achievementcache.New(s.cnf.Redis, s.cnf.AchievementConfig.CacheExpiration)
	s.managers.achievement = achievementdatamanager.New(storage, achievementdatamanager.NewCacheWrapper(cache))
	return s.managers.achievement
}

func (s *setUpper) UserManager() userdatamanager.UserManager {
	if s.managers.user != nil {
		return s.managers.user
	}
	storage := userstorage.New(s.cnf.DB)
	cache := usercache.New(s.cnf.Redis)
	s.managers.user = userdatamanager.New(storage, cache)
	return s.managers.user
}

func (s *setUpper) MotivationManager() motivationdatamanager.MotivationManager {
	if s.managers.motivation != nil {
		return s.managers.motivation
	}
	storage := motivationstorage.New(s.cnf.DB)
	cache := motivationcache.New(s.cnf.Redis, s.cnf.UserConfig.CacheExpiration)
	s.managers.motivation = motivationdatamanager.New(storage, cache)
	return s.managers.motivation
}

func (s *setUpper) AchievementMessageSender() achievementmessagesender.AchievementMessageSender {
	if s.achievementMessageSender != nil {
		return s.achievementMessageSender
	}
	s.achievementMessageSender = achievementmessagesender.NewMessageSender(
		http.DefaultClient,
		s.cnf.VKConfig.ApiVersion,
		s.cnf.VKConfig.GroupToken,
		s.cnf.VKConfig.GroupID,
		s.cnf.VKConfig.AppID,
	)
	return s.achievementMessageSender
}

func (s *setUpper) AchievementMessageSenderAtTime() achievementmessagesender.AchievementMessageSenderAtTime {
	if s.achievementMessageSenderAtTime != nil {
		return s.achievementMessageSenderAtTime
	}
	s.achievementMessageSenderAtTime = achievementmessagesender.NewAchievementMessageSenderAtTime(context.Background(), s.AchievementMessageSender())
	return s.achievementMessageSenderAtTime
}

func (s *setUpper) AchievementExecutor() executor.UserExecutor {
	if s.executors.achievement != nil {
		return s.executors.achievement
	}
	s.executors.achievement = reachachievementexecutor.New(
		s.UserManager(),
		s.AchievementManager(),
		s.AchievementMessageSenderAtTime(),
		reachachievementexecutor.NewAchievementReacher(),
	)
	return s.executors.achievement
}

func (s *setUpper) ChangeMotivationExecutor() executor.UserExecutor {
	if s.executors.motivation != nil {
		return s.executors.motivation
	}
	s.executors.motivation = changemotivationexecutor.New(s.MotivationManager(), s.cnf.MessageSender)
	return s.executors.motivation
}

func (s *setUpper) ChangeWelcomeMotivationExecutor() executor.UserExecutor {
	if s.executors.welcomeMotivation != nil {
		return s.executors.welcomeMotivation
	}
	s.executors.welcomeMotivation = changewelcomemotivationexecutor.New(s.WelcomeMotivationManager())
	return s.executors.welcomeMotivation
}

func Worker(cnf *Config) workers.UserWorker {
	ctx := context.Background()
	s := NewSetUpper(cnf)
	const day = time.Hour * 24
	motivationWorker := userworkers.NewRunnableWorker(s.ChangeMotivationExecutor(), day*2)
	welcomeMotivationWorker := userworkers.NewRunnableWorker(s.ChangeWelcomeMotivationExecutor(), day)
	achievementWorker := userworkers.NewRunnableWorker(s.AchievementExecutor(), day)
	return userworkers.New(ctx, motivationWorker, welcomeMotivationWorker, achievementWorker)
}
