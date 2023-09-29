package local

import (
	"context"

	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
)

type messageDataStorage map[int64][]*model.MessageData

func (m messageDataStorage) AppendToSendTime(ctx context.Context, sendTime int64, messageData *model.MessageData) (index int) {
	index = len(m[sendTime])
	m[sendTime] = append(m[sendTime], messageData)
	return
}

func (m messageDataStorage) PopFromSendTime(ctx context.Context, sendTime int64) (messageData model.MessageData, ok bool) {
	storage, ok := m[sendTime]
	if !ok {
		return
	}
	switch l := len(storage); l {
	case 0:
		return
	case 1:
		ok = true
		messageData = *storage[0]
		delete(m, sendTime)
	default:
		ok = true
		messageData = *m[sendTime][l-1]
		m[sendTime] = m[sendTime][:l-1]
	}
	return
}

func (m messageDataStorage) SendTimeMessageCount(ctx context.Context, sendTime int64) int {
	return len(m[sendTime])
}

func (m messageDataStorage) MarkMessagesAsDeleted(ctx context.Context, sendTime int64, indexes []int) {
	for _, index := range indexes {
		m[sendTime][index].MarkAsDeleted()
	}
}

func NewMessageDataStorage() messageDataStorage {
	return make(messageDataStorage)
}
