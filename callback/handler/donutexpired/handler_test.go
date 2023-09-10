package donutexpired_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Tap-Team/kurilka/callback/handler/donutexpired"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestExpired(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cleaner := donutexpired.NewMockCleaner(ctrl)

	eventHandler := donutexpired.New(cleaner)

	cases := []struct {
		userId int64

		object json.RawMessage

		err         error
		cleanerCall bool
	}{
		{
			userId: 32123,
			object: json.RawMessage(`
				{
					"user_id":32123
				}
			`),

			cleanerCall: true,
		},
		{
			userId: 43234,
			object: json.RawMessage(`
				{
					"user_id":43234
				}
			`),

			err:         errors.New("error while clean subscription"),
			cleanerCall: true,
		},
	}

	for _, cs := range cases {
		if cs.cleanerCall {
			cleaner.EXPECT().CleanSubscription(gomock.Any(), cs.userId).Return(cs.err).Times(1)
		}
		err := eventHandler.HandleEvent(ctx, cs.object)
		assert.ErrorIs(t, err, cs.err, "error not equal")
	}
}
