package schedulerstorage

import (
	context "context"
	"sync"

	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
)

type MessageDataStorage interface {
	AppendToSendTime(ctx context.Context, sendTime int64, messageData *model.MessageData) (index int)
	SendTimeMessageCount(ctx context.Context, sendTime int64) int
	MarkMessagesAsDeleted(ctx context.Context, sendTime int64, indexes []int)
	PopFromSendTime(ctx context.Context, sendTime int64) (messageData model.MessageData, ok bool)
}

type UserMessagesSendTimeStorage interface {
	AddMessageToUser(ctx context.Context, userId int64, sendTime int64, index int)
	RemoveUserMessage(ctx context.Context, userId int64, sendTime int64, index int)
	UserIndexesBySendTime(ctx context.Context, userId int64, sendTime int64) []int
	UserMessagesSendTime(ctx context.Context, userId int64) []int64
}

type messageSchedulerStorage struct {
	mu sync.Mutex
	// map user id to time when we need send messages
	userMessagesSendTimeStorage UserMessagesSendTimeStorage
	// map send time to array of message send data
	storage MessageDataStorage
}

type MessageSchedulerStorage interface {
	AddMessageData(ctx context.Context, int64, userId int64, messageData model.AchievementMessageData)
	PopMessageData(ctx context.Context, sendTimeSeconds int64) (model.MessageData, bool)
	ClearUserMessages(ctx context.Context, userId int64)
}

func NewMessageSchedulerStorage(messageDataStorage MessageDataStorage, userMessageSendTimeStorage UserMessagesSendTimeStorage) MessageSchedulerStorage {
	return &messageSchedulerStorage{
		storage:                     messageDataStorage,
		userMessagesSendTimeStorage: userMessageSendTimeStorage,
	}
}

func (m *messageSchedulerStorage) AddMessageData(ctx context.Context, sendTimeSeconds int64, userId int64, messageData model.AchievementMessageData) {
	m.mu.Lock()
	index := m.storage.AppendToSendTime(ctx, sendTimeSeconds, model.NewMessageData(userId, messageData))
	m.userMessagesSendTimeStorage.AddMessageToUser(ctx, userId, sendTimeSeconds, index)
	m.mu.Unlock()
}

func (m *messageSchedulerStorage) PopMessageData(ctx context.Context, sendTimeSeconds int64) (model.MessageData, bool) {
	m.mu.Lock()
	messageData, ok := m.storage.PopFromSendTime(ctx, sendTimeSeconds)
	if ok {
		index := m.storage.SendTimeMessageCount(ctx, sendTimeSeconds)
		m.userMessagesSendTimeStorage.RemoveUserMessage(ctx, messageData.UserId(), sendTimeSeconds, index)
	}
	m.mu.Unlock()
	if !ok {
		return model.MessageData{}, false
	}
	if messageData.IsDeleted() {
		return m.PopMessageData(ctx, sendTimeSeconds)
	}
	return messageData, ok
}

func (m *messageSchedulerStorage) ClearUserMessages(ctx context.Context, userId int64) {
	m.mu.Lock()
	for _, sendTime := range m.userMessagesSendTimeStorage.UserMessagesSendTime(ctx, userId) {
		indexes := m.userMessagesSendTimeStorage.UserIndexesBySendTime(ctx, userId, sendTime)
		m.storage.MarkMessagesAsDeleted(ctx, sendTime, indexes)
	}
	m.mu.Unlock()
}
