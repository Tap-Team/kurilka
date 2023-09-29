package achievementmessagesender

import (
	"context"
	"sync"
	"time"

	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage"
)

//go:generate mockgen -source scheduler.go -destination scheduler_mocks.go -package achievementmessagesender

type AchievementMessageSenderAtTime interface {
	SendMessageAtTime(ctx context.Context, userId int64, messageData model.AchievementMessageData, t time.Time)
	CancelSendMessagesForUser(ctx context.Context, userId int64)
}

type achievementMessageScheduler struct {
	AchievementMessageSender
	isRunned bool
	storage  schedulerstorage.MessageSchedulerStorage
	tick     time.Duration
}

func NewAchievementMessageSenderAtTime(
	ctx context.Context,
	messageSender AchievementMessageSender,
	messageSchedulerStorage schedulerstorage.MessageSchedulerStorage,
) AchievementMessageSenderAtTime {
	scheduler := &achievementMessageScheduler{
		AchievementMessageSender: messageSender,
		storage:                  messageSchedulerStorage,
		tick:                     time.Second,
	}
	go scheduler.run(ctx)
	return scheduler
}

func NewAchievementMessageSchedulerWithTickTime(
	ctx context.Context,
	messageSender AchievementMessageSender,
	messageSchedulerStorage schedulerstorage.MessageSchedulerStorage,

	tick time.Duration,
) AchievementMessageSenderAtTime {
	scheduler := &achievementMessageScheduler{
		AchievementMessageSender: messageSender,
		storage:                  messageSchedulerStorage,
		tick:                     tick,
	}
	go scheduler.run(ctx)
	return scheduler
}

func (ms *achievementMessageScheduler) run(ctx context.Context) {
	if ms.isRunned {
		return
	}
	ms.isRunned = true
	seconds := time.Now().Unix()
	defer func() { ms.isRunned = false }()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			go ms.sendMessages(ctx, seconds)
			seconds++
			time.Sleep(ms.tick)
		}
	}
}

func (ms *achievementMessageScheduler) sendMessages(ctx context.Context, seconds int64) {
	var wg sync.WaitGroup
	for {
		messageData, ok := ms.storage.PopMessageData(ctx, seconds)
		if !ok {
			break
		}
		wg.Add(1)
		go func(userId int64, messageData model.AchievementMessageData) {
			defer wg.Done()
			ms.SendMessage(ctx, userId, messageData)
		}(messageData.UserId(), messageData.AchievementMessageData())
	}
	wg.Wait()
}

func (ms *achievementMessageScheduler) CancelSendMessagesForUser(ctx context.Context, userId int64) {
	ms.storage.ClearUserMessages(ctx, userId)
}

func (ms *achievementMessageScheduler) SendMessageAtTime(ctx context.Context, userId int64, messageData model.AchievementMessageData, t time.Time) {
	ms.storage.AddMessageData(ctx, t.Unix(), userId, messageData)
}
