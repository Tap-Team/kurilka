package messagesender_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/pkg/random"
	"gotest.tools/v3/assert"
)

func Test_MessageSchedulerStorage(t *testing.T) {
	storage := messagesender.NewMessageSchedulerStorage()

	messages := map[int64][]*messagesender.MessageData{
		1: {
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
		},
		2: {
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
			messagesender.NewMessageData(rand.Int63(), random.String(100)),
		},
	}

	for sec, messages := range messages {
		for _, msg := range messages {
			storage.AddMessageData(sec, msg.UserId(), msg.Message())
		}
	}

	for sec, messages := range messages {
		l := len(messages)
		for index := range messages {
			msg := messages[l-1-index]
			popMsg, ok := storage.PopMessageData(sec)
			assert.Equal(t, true, ok)
			assert.Equal(t, msg.Message(), popMsg.Message(), "message not equal")
			assert.Equal(t, msg.UserId(), popMsg.UserId(), "user id not equal")
		}
		_, ok := storage.PopMessageData(sec)
		assert.Equal(t, false, ok)
	}
}

type fakeMessageSender struct {
	t  *testing.T
	mu *sync.Mutex
	// map user id to message
	registeredCalls map[int64]string
}

func NewFakeMessageSender(t *testing.T) *fakeMessageSender {
	var mu sync.Mutex
	calls := make(map[int64]string)
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		if len(calls) != 0 {
			t.Fatal("no calls")
		}
	})
	return &fakeMessageSender{mu: &mu, t: t, registeredCalls: calls}
}

func (f *fakeMessageSender) RegisterCall(message string, userId int64) {
	f.mu.Lock()
	f.registeredCalls[userId] = message
	f.mu.Unlock()
}

func (f *fakeMessageSender) SendMessage(ctx context.Context, message string, userId int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	msg, ok := f.registeredCalls[userId]
	if !ok || msg != message {
		f.t.Fatalf("unexpected call to send message, msg:%s,userId:%d", message, userId)
	}
	delete(f.registeredCalls, userId)
	return nil
}

func TestMessageScheduler(t *testing.T) {
	ctx := context.Background()

	sender := NewFakeMessageSender(t)

	now := time.Now()
	messageScheduler := messagesender.NewMessageSchedulerWithTickTime(ctx, sender, time.Millisecond*100)

	messages := []*messagesender.MessageData{
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
		messagesender.NewMessageData(rand.Int63(), random.String(100)),
	}

	for _, msg := range messages {
		msg := *msg
		messageScheduler.SendMessageAtTime(ctx, msg.Message(), msg.UserId(), now.Add(time.Second*2))
		sender.RegisterCall(msg.Message(), msg.UserId())
	}

	time.Sleep(time.Millisecond * 300)
}
