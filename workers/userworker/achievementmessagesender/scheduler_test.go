package achievementmessagesender_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/workers/userworker/achievementmessagesender"
	"gotest.tools/v3/assert"
)

func RandomAchievementMessage() achievementmessagesender.AchievementMessageData {
	types := []achievementmodel.AchievementType{achievementmodel.CIGARETTE, achievementmodel.HEALTH, achievementmodel.SAVING, achievementmodel.DURATION, achievementmodel.WELL_BEING}
	i := rand.Intn(len(types))
	tp := types[i]
	return achievementmessagesender.NewAchievementMessageData(tp)
}

func Test_MessageSchedulerStorage(t *testing.T) {
	storage := achievementmessagesender.NewMessageSchedulerStorage()

	messages := map[int64][]*achievementmessagesender.MessageData{
		1: {
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		},
		2: {
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		},
	}

	for sec, messages := range messages {
		for _, msg := range messages {
			storage.AddMessageData(sec, msg.UserId(), msg.AchievementMessageData())
		}
	}

	for sec, messages := range messages {
		l := len(messages)
		for index := range messages {
			msg := messages[l-1-index]
			popMsg, ok := storage.PopMessageData(sec)
			assert.Equal(t, true, ok)
			assert.Equal(t, msg.AchievementMessageData(), popMsg.AchievementMessageData(), "message not equal")
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
	registeredCalls map[int64]achievementmessagesender.AchievementMessageData
}

func NewFakeMessageSender(t *testing.T) *fakeMessageSender {
	var mu sync.Mutex
	calls := make(map[int64]achievementmessagesender.AchievementMessageData)
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		if len(calls) != 0 {
			t.Fatal("no calls")
		}
	})
	return &fakeMessageSender{mu: &mu, t: t, registeredCalls: calls}
}

func (f *fakeMessageSender) RegisterCall(userId int64, messageData achievementmessagesender.AchievementMessageData) {
	f.mu.Lock()
	f.registeredCalls[userId] = messageData
	f.mu.Unlock()
}

func (f *fakeMessageSender) SendMessage(ctx context.Context, userId int64, messageData achievementmessagesender.AchievementMessageData) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	msg, ok := f.registeredCalls[userId]
	if !ok || msg != messageData {
		f.t.Fatalf("unexpected call to send message, msg:%s,userId:%d", messageData, userId)
	}
	delete(f.registeredCalls, userId)
	return nil
}

func Test_AchievementMessageScheduler(t *testing.T) {
	ctx := context.Background()
	sender := NewFakeMessageSender(t)

	messageScheduler := achievementmessagesender.NewAchievementMessageSchedulerWithTickTime(ctx, sender, time.Millisecond*100)

	now := time.Now()

	messages := []*achievementmessagesender.MessageData{
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		achievementmessagesender.NewMessageData(rand.Int63(), RandomAchievementMessage()),
	}
	for _, msg := range messages {
		msg := *msg
		messageScheduler.SendMessageAtTime(ctx, msg.UserId(), msg.AchievementMessageData(), now.Add(time.Second*2))
		sender.RegisterCall(msg.UserId(), msg.AchievementMessageData())
	}

	time.Sleep(time.Millisecond * 300)
}
