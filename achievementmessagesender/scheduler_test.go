package achievementmessagesender_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievementmessagesender"
	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage/local"
)

type fakeMessageSender struct {
	t  *testing.T
	mu *sync.Mutex
	// map user id to message
	registeredCalls map[int64]model.AchievementMessageData
}

func NewFakeMessageSender(t *testing.T) *fakeMessageSender {
	var mu sync.Mutex
	calls := make(map[int64]model.AchievementMessageData)
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		if len(calls) != 0 {
			t.Fatal("no calls")
		}
	})
	return &fakeMessageSender{mu: &mu, t: t, registeredCalls: calls}
}

func (f *fakeMessageSender) RegisterCall(userId int64, messageData model.AchievementMessageData) {
	f.mu.Lock()
	f.registeredCalls[userId] = messageData
	f.mu.Unlock()
}

func (f *fakeMessageSender) SendMessage(ctx context.Context, userId int64, messageData model.AchievementMessageData) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	msg, ok := f.registeredCalls[userId]
	if !ok || msg != messageData {
		f.t.Fatalf("unexpected call to send message, msg:%s,userId:%d", messageData, userId)
	}
	delete(f.registeredCalls, userId)
	return nil
}

func Test_AchievementMessageScheduler_Local(t *testing.T) {
	ctx := context.Background()
	sender := NewFakeMessageSender(t)

	messageStorage := schedulerstorage.NewMessageSchedulerStorage(local.NewMessageDataStorage(), local.NewMessageSendTimeStorage())
	messageScheduler := achievementmessagesender.NewAchievementMessageSchedulerWithTickTime(ctx, sender, messageStorage, time.Millisecond*100)

	now := time.Now()

	messages := []*model.MessageData{
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
	}
	for _, msg := range messages {
		msg := *msg
		messageScheduler.SendMessageAtTime(ctx, msg.UserId(), msg.AchievementMessageData(), now.Add(time.Second*2))
		sender.RegisterCall(msg.UserId(), msg.AchievementMessageData())
	}

	time.Sleep(time.Millisecond * 300)
}
