package local

import (
	"context"
	"sort"

	"golang.org/x/exp/slices"
)

type userMessagesSendTimeStorage map[int64]map[int64][]int

func (u userMessagesSendTimeStorage) AddMessageToUser(ctx context.Context, userId int64, sendTime int64, index int) {
	_, ok := u[userId]
	if !ok {
		u[userId] = make(map[int64][]int, 1)
	}
	u[userId][sendTime] = append(u[userId][sendTime], index)
}

func (u userMessagesSendTimeStorage) RemoveUserMessage(ctx context.Context, userId int64, sendTime int64, index int) {
	_, ok := u[userId]
	if !ok {
		return
	}
	indexPosition := sort.Search(len(u[userId][sendTime]), func(i int) bool { return u[userId][sendTime][i] == index })
	u[userId][sendTime] = slices.Delete(u[userId][sendTime], indexPosition, indexPosition+1)
	if len(u[userId][sendTime]) == 0 {
		delete(u[userId], sendTime)
	}
}

func (u userMessagesSendTimeStorage) UserIndexesBySendTime(ctx context.Context, userId int64, sendTime int64) []int {
	return u[userId][sendTime]
}

func (u userMessagesSendTimeStorage) UserMessagesSendTime(ctx context.Context, userId int64) []int64 {
	messagesSendTime := make([]int64, 0, len(u[userId]))
	for sendTime := range u[userId] {
		messagesSendTime = append(messagesSendTime, sendTime)
	}
	return messagesSendTime
}

func NewMessageSendTimeStorage() userMessagesSendTimeStorage {
	return make(userMessagesSendTimeStorage)
}
