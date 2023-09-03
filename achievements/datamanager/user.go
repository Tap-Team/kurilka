package datamanager

import (
	"context"

	"github.com/Tap-Team/kurilka/achievements/model"
)

type UserDataManager interface {
	UserData(ctx context.Context, userId int64) (*model.UserData, error)
}
