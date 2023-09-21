package achievementmessagesender

import (
	"context"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
)

type AchievementMessageSender interface {
	SendMessage(ctx context.Context, userId int64, message string, achievementType achievementmodel.AchievementType) error
}

type achievementMessageSender struct {
}

func New() AchievementMessageSender {

}
