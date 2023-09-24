package messagesender

import (
	"context"
	"sync"
	"time"
)

type MessageData struct {
	message string
	userId  int64
}

func (m *MessageData) UserId() int64 {
	return m.userId
}
func (m *MessageData) Message() string {
	return m.message
}

func NewMessageData(userId int64, message string) *MessageData {
	return &MessageData{
		message: message,
		userId:  userId,
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

func (m *MessageSchedulerStorage) AddMessageData(sendTimeSeconds int64, userId int64, message string) {
	m.mu.Lock()
	m.storage[sendTimeSeconds] = append(m.storage[sendTimeSeconds], NewMessageData(userId, message))
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

type MessageScheduler struct {
	MessageSender
	isRunned bool
	storage  *MessageSchedulerStorage
	tick     time.Duration
}

func NewMessageScheduler(ctx context.Context, messageSender MessageSender) *MessageScheduler {
	scheduler := &MessageScheduler{
		storage:       NewMessageSchedulerStorage(),
		MessageSender: messageSender,
		tick:          time.Second,
	}
	go scheduler.run(ctx)
	return scheduler
}

func NewMessageSchedulerWithTickTime(ctx context.Context, messageSender MessageSender, tick time.Duration) *MessageScheduler {
	scheduler := &MessageScheduler{
		storage:       NewMessageSchedulerStorage(),
		MessageSender: messageSender,
		tick:          tick,
	}
	go scheduler.run(ctx)
	return scheduler
}

func (ms *MessageScheduler) run(ctx context.Context) {
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

func (ms *MessageScheduler) sendMessages(ctx context.Context, seconds int64) {
	var wg sync.WaitGroup
	for {
		messageData, ok := ms.storage.PopMessageData(seconds)
		if !ok {
			break
		}
		wg.Add(1)
		go func(userId int64, message string) {
			defer wg.Done()
			ms.SendMessage(ctx, message, userId)
		}(messageData.userId, messageData.message)
	}
	wg.Wait()
}

func (ms *MessageScheduler) SendMessageAtTime(ctx context.Context, message string, userId int64, t time.Time) {
	ms.storage.AddMessageData(t.Unix(), userId, message)
}
