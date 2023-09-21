package changewelcomemotivationexecutor_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/welcomemotivationdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor/changewelcomemotivationexecutor"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

var (
	anyErr = errors.New("hello world failed")
)

func Test_Executor(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := welcomemotivationdatamanager.NewMockWelcomeMotivationManager(ctrl)

	executor := changewelcomemotivationexecutor.New(manager)

	cases := []struct {
		nextUserMotivationCall     bool
		nextUserMotivationResponse struct {
			motivation model.Motivation
			err        error
		}

		updateUserMotivationCall bool
		updateUserMotivationErr  error

		err error
	}{
		{
			nextUserMotivationCall: true,
			nextUserMotivationResponse: struct {
				motivation model.Motivation
				err        error
			}{
				err: usererror.ExceptionUserNotFound(),
			},
			err: usererror.ExceptionUserNotFound(),
		},
		{
			nextUserMotivationCall: true,
			nextUserMotivationResponse: struct {
				motivation model.Motivation
				err        error
			}{
				motivation: random.StructTyped[model.Motivation](),
			},
			updateUserMotivationCall: true,
			updateUserMotivationErr:  anyErr,

			err: anyErr,
		},
		{
			nextUserMotivationCall: true,
			nextUserMotivationResponse: struct {
				motivation model.Motivation
				err        error
			}{
				motivation: random.StructTyped[model.Motivation](),
			},
			updateUserMotivationCall: true,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()

		if cs.nextUserMotivationCall {
			r := cs.nextUserMotivationResponse
			manager.EXPECT().NextUserWelcomeMotivation(gomock.Any(), userId).Return(r.motivation, r.err).Times(1)
		}
		if cs.updateUserMotivationCall {
			r := cs.nextUserMotivationResponse
			manager.EXPECT().UpdateUserWelcomeMotivation(gomock.Any(), userId, r.motivation).Return(cs.updateUserMotivationErr).Times(1)
		}

		err := executor.ExecuteUser(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}

}
