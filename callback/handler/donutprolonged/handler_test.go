package donutprolonged_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Tap-Team/kurilka/callback/handler/donutprolonged"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestDonutProlonged(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	prolong := donutprolonged.NewMockProlongationer(ctrl)

	eventHandler := donutprolonged.New(prolong)

	cases := []struct {
		userId int64
		amount float64

		object json.RawMessage

		prolongCall bool
		err         error
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

			prolongCall: true,
			err:         usererror.ExceptionUserNotFound(),
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

			prolongCall: true,
		},
	}

	for _, cs := range cases {
		if cs.prolongCall {
			prolong.EXPECT().ProlongSubscription(gomock.Any(), cs.userId, cs.amount).Return(cs.err).Times(1)
		}
		err := eventHandler.HandleEvent(ctx, cs.object)

		assert.ErrorIs(t, err, cs.err)
	}
}
