package app

import (
	"context"
	"net/http"

	"github.com/Tap-Team/kurilka/achievementmessagesender"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage/local"
	"github.com/Tap-Team/kurilka/internal/config"
)

func AchievementMessageSenderScheduler(ctx context.Context, cnf config.VKConfig) achievementmessagesender.AchievementMessageSenderAtTime {
	messageSender := achievementmessagesender.NewMessageSender(http.DefaultClient, cnf.ApiVersion, cnf.GroupAccessKey, int(cnf.GroupID), int(cnf.AppID))

	messageDataStorage := local.NewMessageDataStorage()
	userMessageSendTimeStorage := local.NewMessageSendTimeStorage()
	schedulerStorage := schedulerstorage.NewMessageSchedulerStorage(messageDataStorage, userMessageSendTimeStorage)
	return achievementmessagesender.NewAchievementMessageSenderAtTime(ctx, messageSender, schedulerStorage)
}
