package achievementusecase_test

import (
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/golang/mock/gomock"
)

type SendMessageCall struct {
	WillBeCalled bool

	Message string
	UserId  int64

	Err error
}

func (c *SendMessageCall) RegisterCall(messageSender *messagesender.MockMessageSender) {
	if c.WillBeCalled {
		messageSender.EXPECT().
			SendMessage(gomock.Any(), c.Message, c.UserId).
			Return(c.Err).
			Times(1)
	}
}

type SendMessageCallBuilder struct {
	message string
	userId  int64

	err error
}

func (b *SendMessageCallBuilder) SetInput(message string, userId int64) *SendMessageCallBuilder {
	b.message = message
	b.userId = userId
	return b
}

func (b *SendMessageCallBuilder) SetOutput(err error) *SendMessageCallBuilder {
	b.err = err
	return b
}

func (b *SendMessageCallBuilder) Build() SendMessageCall {
	if b == nil {
		return SendMessageCall{}
	}
	return SendMessageCall{
		UserId:  b.userId,
		Message: b.message,
		Err:     b.err,

		WillBeCalled: true,
	}
}
