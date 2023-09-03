package userdatamanager

import (
	"context"
	"errors"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "user/datamanager"

type UserReseter interface {
	ResetUser(ctx context.Context, userId int64) error
}
type UserRecoverer interface {
	RecoverUser(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.UserData, error)
}

type UserRecoverReseter interface {
	UserReseter
	UserRecoverer
}

type UserStorage interface {
	Exists(ctx context.Context, userIds []int64) []int64
	User(ctx context.Context, userId int64) (*usermodel.UserData, error)
	InsertUser(ctx context.Context, userId int64, user *usermodel.CreateUser) (*usermodel.UserData, error)
	UserDeleted(ctx context.Context, userId int64) (bool, error)
	UserLevel(ctx context.Context, userId int64) (*usermodel.LevelInfo, error)
}

type UserCache interface {
	User(ctx context.Context, userId int64) (*usermodel.UserData, error)
	SaveUser(ctx context.Context, userId int64, user *usermodel.UserData) error
	RemoveUser(ctx context.Context, userId int64) error
}

type UserManager interface {
	Create(ctx context.Context, userId int64, user *usermodel.CreateUser) (*usermodel.UserData, error)
	User(ctx context.Context, userId int64) (*usermodel.UserData, error)
	Level(ctx context.Context, userId int64) (*usermodel.LevelInfo, error)
	Reset(ctx context.Context, userId int64) error
	FilterExists(ctx context.Context, userIds []int64) []int64
}

type userManager struct {
	saver   UserSaver
	reseter UserReseter
	storage UserStorage
	cache   UserCache
}

func NewUserManager(recoverReseter UserRecoverReseter, storage UserStorage, cache UserCache, saver UserSaver) UserManager {
	return &userManager{
		reseter: recoverReseter,
		storage: storage,
		cache:   cache,
		saver:   saver,
	}
}

type UserSaver interface {
	Save(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.UserData, error)
}

type userSaver struct {
	storage   UserStorage
	recoverer UserRecoverer
}

func NewUserSaver(storage UserStorage, recoverer UserRecoverer) UserSaver {
	return &userSaver{
		storage:   storage,
		recoverer: recoverer,
	}
}

func (us *userSaver) Save(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.UserData, error) {
	deleted, err := us.storage.UserDeleted(ctx, userId)
	if errors.Is(err, usererror.ExceptionUserNotFound()) {
		return us.storage.InsertUser(ctx, userId, createUser)
	}
	if err != nil {
		return nil, err
	}
	if deleted {
		return us.recoverer.RecoverUser(ctx, userId, createUser)
	} else {
		return nil, usererror.ExceptionUserExist()
	}
}

func (u *userManager) Create(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.UserData, error) {
	userData, err := u.saver.Save(ctx, userId, createUser)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("save user", "CreateUser", _PROVIDER))
	}
	u.cache.SaveUser(ctx, userId, userData)
	return userData, nil
}

func (u *userManager) User(ctx context.Context, userId int64) (*usermodel.UserData, error) {
	user, err := u.cache.User(ctx, userId)
	if err == nil {
		return user, nil
	}
	user, err = u.storage.User(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("get user from storage", "userdata", _PROVIDER))
	}
	u.cache.SaveUser(ctx, userId, user)
	return user, nil
}

func (u *userManager) Level(ctx context.Context, userId int64) (*usermodel.LevelInfo, error) {
	level, err := u.storage.UserLevel(ctx, userId)
	if err != nil {
		return nil, exception.Wrap(err, exception.NewCause("user level", "Level", _PROVIDER))
	}
	user, err := u.cache.User(ctx, userId)
	if err != nil {
		return level, nil
	}
	user.Level = *level
	err = u.cache.SaveUser(ctx, userId, user)
	if err != nil {
		slog.Error(exception.Wrap(err, exception.NewCause("failed save user", "Level", _PROVIDER)).Error())
		u.cache.RemoveUser(ctx, userId)
	}
	return level, nil
}

func (u *userManager) Reset(ctx context.Context, userId int64) error {
	err := u.reseter.ResetUser(ctx, userId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("reset user", "Reset", _PROVIDER))
	}
	u.cache.RemoveUser(ctx, userId)
	return nil
}

func (u *userManager) FilterExists(ctx context.Context, userIds []int64) []int64 {
	return u.storage.Exists(ctx, userIds)
}
