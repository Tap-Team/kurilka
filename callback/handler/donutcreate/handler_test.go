package donutcreate_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Tap-Team/kurilka/callback/handler/donutcreate"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	creator := donutcreate.NewMockCreator(ctrl)

	eventHandler := donutcreate.New(creator)

	cases := []struct {
		userId int64
		amount int

		object json.RawMessage

		createCall bool
		err        error
	}{
		{
			userId: 1,
			amount: 100,
			object: json.RawMessage(`
				{
					"user_id":1,
					"amount": 100,
					"amount_without_fee": 120.98
				}
			`),

			createCall: true,
			err:        usererror.ExceptionUserNotFound(),
		},
		{
			userId: 1001231,
			amount: 123,
			object: json.RawMessage(`
				{
					"user_id":1001231,
					"amount": 123,
					"amount_without_fee": 123.60
				}
			`),

			createCall: true,
		},
	}

	for _, cs := range cases {
		if cs.createCall {
			creator.EXPECT().CreateSubscription(gomock.Any(), cs.userId, cs.amount).Return(cs.err)
		}

		err := eventHandler.HandleEvent(ctx, cs.object)
		assert.ErrorIs(t, err, cs.err, "wrong error")
	}
}
