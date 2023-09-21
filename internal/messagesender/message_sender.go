package messagesender

import (
	"context"
	"time"
)

//go:generate mockgen -source message_sender.go -destination mocks.go -package messagesender

type MessageSender interface {
	SendMessage(ctx context.Context, message string, userId int64) error
}

type MessageSenderAtTime interface {
	SendMessageAtTime(ctx context.Context, message string, userId int64, t time.Time)
}
