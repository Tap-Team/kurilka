package changemotivationexecutor_test

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/motivationerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/motivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor/changemotivationexecutor"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_Executor(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	messageSender := messagesender.NewMockMessageSenderAtTime(ctrl)
	motivation := motivationdatamanager.NewMockMotivationManager(ctrl)

	executor := changemotivationexecutor.New(motivation, messageSender)

	cases := []struct {
		motivation    model.Motivation
		motivationErr error

		updateMotivationCall bool
		updateMotivationErr  error

		messageSenderCall bool

		err error
	}{
		{
			motivation:    random.StructTyped[model.Motivation](),
			motivationErr: usererror.ExceptionUserNotFound(),
			err:           usererror.ExceptionUserNotFound(),
		},
		{
			motivation:           random.StructTyped[model.Motivation](),
			updateMotivationCall: true,
			updateMotivationErr:  motivationerror.ExceptionMotivationNotExist(),
			err:                  motivationerror.ExceptionMotivationNotExist(),
		},
		{
			motivation:           random.StructTyped[model.Motivation](),
			updateMotivationCall: true,
			messageSenderCall:    true,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		motivation.EXPECT().NextUserMotivation(gomock.Any(), userId).Return(cs.motivation, cs.motivationErr).Times(1)
		if cs.updateMotivationCall {
			motivation.EXPECT().UpdateUserMotivation(gomock.Any(), userId, cs.motivation).Return(cs.updateMotivationErr).Times(1)
		}
		if cs.messageSenderCall {
			messageSender.EXPECT().SendMessageAtTime(gomock.Any(), gomock.Any(), userId, changemotivationexecutor.NextSendTime(time.Now())).Times(1)
		}

		err := executor.ExecuteUser(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}
}

func Test_NextSendTime(t *testing.T) {
	cases := []struct {
		now time.Time

		expected time.Time

		unix int64
	}{
		{
			now:      time.Date(2023, time.April, 1, 11, 1, 1, 1, time.UTC),
			expected: time.Date(2023, time.April, 2, 14, 0, 0, 0, changemotivationexecutor.MoscowLocation),

			unix: time.Date(2023, time.April, 2, 11, 0, 0, 0, time.UTC).Unix(),
		},
		{
			now:      time.Date(2023, time.April, 1, 10, 59, 59, 0, time.UTC),
			expected: time.Date(2023, time.April, 1, 14, 0, 0, 0, changemotivationexecutor.MoscowLocation),

			unix: time.Date(2023, time.April, 1, 11, 0, 0, 0, time.UTC).Unix(),
		},
	}

	for num, cs := range cases {
		actual := changemotivationexecutor.NextSendTime(cs.now)
		log.Println(cs.unix - actual.Unix())
		assert.Equal(t, cs.expected, actual, "wrong time, %d", num)
		assert.Equal(t, cs.unix, actual.Unix(), "unix not equal, %d", num)
	}
}
