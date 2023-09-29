package schedulerstorage_test

import (
	"context"
	"log"
	"math/rand"
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage"
	"github.com/Tap-Team/kurilka/achievementmessagesender/schedulerstorage/local"
	"gotest.tools/v3/assert"
)

func Test_MessageSchedulerStorage_Local(t *testing.T) {
	ctx := context.Background()

	storage := schedulerstorage.NewMessageSchedulerStorage(local.NewMessageDataStorage(), local.NewMessageSendTimeStorage())

	messages := map[int64][]*model.MessageData{
		1: {
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		},
		2: {
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
			model.NewMessageData(rand.Int63(), RandomAchievementMessage()),
		},
	}

	for sec, messages := range messages {
		for _, msg := range messages {
			storage.AddMessageData(ctx, sec, msg.UserId(), msg.AchievementMessageData())
		}
	}

	for sec, messages := range messages {
		l := len(messages)
		for index := range messages {
			msg := messages[l-1-index]
			popMsg, ok := storage.PopMessageData(ctx, sec)
			assert.Equal(t, true, ok)
			assert.Equal(t, msg.AchievementMessageData(), popMsg.AchievementMessageData(), "message not equal")
			assert.Equal(t, msg.UserId(), popMsg.UserId(), "user id not equal")
		}
		_, ok := storage.PopMessageData(ctx, sec)
		assert.Equal(t, false, ok)
	}
}

func Test_MessageSchedulerStorage_ClearUserMessages_Local(t *testing.T) {
	ctx := context.Background()
	// len of sendTimeList must equal len of messages list
	cases := []struct {
		sendTimeList []int64
		messagesList [][]*model.MessageData

		expectedMessageList [][]*model.MessageData

		deletedUsers []int64
	}{
		{
			sendTimeList: []int64{
				rand.Int63(),
				rand.Int63(),
			},
			messagesList: [][]*model.MessageData{
				{
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(3, RandomAchievementMessage()),
					model.NewMessageData(4, RandomAchievementMessage()),
					model.NewMessageData(10, RandomAchievementMessage()),
					model.NewMessageData(10, RandomAchievementMessage()),
					model.NewMessageData(10, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(3, RandomAchievementMessage()),
				},
				{
					model.NewMessageData(11, RandomAchievementMessage()),
					model.NewMessageData(12, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(3, RandomAchievementMessage()),
					model.NewMessageData(4, RandomAchievementMessage()),
					model.NewMessageData(10, RandomAchievementMessage()),
					model.NewMessageData(10, RandomAchievementMessage()),
					model.NewMessageData(10, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(3, RandomAchievementMessage()),
				},
			},
			expectedMessageList: [][]*model.MessageData{
				{
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(4, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
				},
				{
					model.NewMessageData(11, RandomAchievementMessage()),
					model.NewMessageData(12, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
					model.NewMessageData(4, RandomAchievementMessage()),
					model.NewMessageData(1, RandomAchievementMessage()),
				},
			},
			deletedUsers: []int64{3, 10},
		},
	}

	for _, cs := range cases {
		storage := schedulerstorage.NewMessageSchedulerStorage(local.NewMessageDataStorage(), local.NewMessageSendTimeStorage())
		for index, sendTime := range cs.sendTimeList {
			messageList := cs.messagesList[index]
			for _, message := range messageList {
				storage.AddMessageData(ctx, sendTime, message.UserId(), message.AchievementMessageData())
			}
		}

		for _, user := range cs.deletedUsers {
			storage.ClearUserMessages(ctx, user)
		}

		for index, sendTime := range cs.sendTimeList {
			expectedMessageList := cs.expectedMessageList[index]
			l := len(expectedMessageList)
			actualMessageList := make([]*model.MessageData, l)
			for {
				message, ok := storage.PopMessageData(ctx, sendTime)
				if !ok {
					break
				}
				actualMessageList[l-1] = &message
				l--
			}
			log.Println("expected")
			for _, m := range expectedMessageList {
				log.Printf("%v", m)
			}
			log.Println("actual")
			for _, m := range actualMessageList {
				log.Printf("%v", m)
			}
			equal := slices.EqualFunc(expectedMessageList, actualMessageList, achievementMessageEqual)
			assert.Equal(t, true, equal, "slices not equal, %d", index)
		}
	}
}
