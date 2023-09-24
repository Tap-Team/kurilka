package achievementmessagesender

import (
	"context"
	"sync"
	"time"
)

//go:generate mockgen -source scheduler.go -destination scheduler_mocks.go -package achievementmessagesender

type MessageData struct {
	achievementMessageData AchievementMessageData
	userId                 int64
}

func (m *MessageData) UserId() int64 {
	return m.userId
}

func (m *MessageData) AchievementMessageData() AchievementMessageData {
	return m.achievementMessageData
}

func NewMessageData(userId int64, achievementMessageData AchievementMessageData) *MessageData {
	return &MessageData{
		achievementMessageData: achievementMessageData,
		userId:                 userId,
	}
}

type MessageSchedulerStorage struct {
	mu sync.Mutex
	// map send time to array of message send data
	storage map[int64][]*MessageData
}

func NewMessageSchedulerStorage() *MessageSchedulerStorage {
	return &MessageSchedulerStorage{
		storage: make(map[int64][]*MessageData),
	}
}

func (m *MessageSchedulerStorage) AddMessageData(sendTimeSeconds int64, userId int64, messageData AchievementMessageData) {
	m.mu.Lock()
	m.storage[sendTimeSeconds] = append(m.storage[sendTimeSeconds], NewMessageData(userId, messageData))
	m.mu.Unlock()
}

func (m *MessageSchedulerStorage) PopMessageData(sendTimeSeconds int64) (MessageData, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	storage := m.storage[sendTimeSeconds]
	switch l := len(storage); l {
	case 0:
		return MessageData{}, false
	case 1:
		delete(m.storage, sendTimeSeconds)
		return *storage[0], true
	default:
		sendMessageData := m.storage[sendTimeSeconds][l-1]
		m.storage[sendTimeSeconds] = m.storage[sendTimeSeconds][:l-1]
		return *sendMessageData, true
	}
}

type AchievementMessageSenderAtTime interface {
	SendMessageAtTime(ctx context.Context, userId int64, messageData AchievementMessageData, t time.Time)
}

type achievementMessageScheduler struct {
	AchievementMessageSender
	isRunned bool
	storage  *MessageSchedulerStorage
	tick     time.Duration
}

func NewAchievementMessageSenderAtTime(
	ctx context.Context,
	messageSender AchievementMessageSender,
) AchievementMessageSenderAtTime {
	scheduler := &achievementMessageScheduler{
		AchievementMessageSender: messageSender,
		storage:                  NewMessageSchedulerStorage(),
		tick:                     time.Second,
	}
	go scheduler.run(ctx)
	return scheduler
}

func NewAchievementMessageSchedulerWithTickTime(
	ctx context.Context,
	messageSender AchievementMessageSender,
	tick time.Duration,
) AchievementMessageSenderAtTime {
	scheduler := &achievementMessageScheduler{
		AchievementMessageSender: messageSender,
		storage:                  NewMessageSchedulerStorage(),
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
		messageData, ok := ms.storage.PopMessageData(seconds)
		if !ok {
			break
		}
		wg.Add(1)
		go func(userId int64, messageData AchievementMessageData) {
			defer wg.Done()
			ms.SendMessage(ctx, userId, messageData)
		}(messageData.userId, messageData.AchievementMessageData())
	}
	wg.Wait()
}

func (ms *achievementMessageScheduler) SendMessageAtTime(ctx context.Context, userId int64, messageData AchievementMessageData, t time.Time) {
	ms.storage.AddMessageData(t.Unix(), userId, messageData)
}
