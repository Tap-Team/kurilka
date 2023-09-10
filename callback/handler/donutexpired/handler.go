package donutexpired

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source handler.go -destination mocks.go -package donutexpired

const _PROVIDER = "callback/handler/donutexpired"

type Cleaner interface {
	CleanSubscription(ctx context.Context, userId int64) error
}

type handler struct {
	cleaner Cleaner
}

func New(cleaner Cleaner) *handler {
	return &handler{
		cleaner: cleaner,
	}
}

type donutExpired struct {
	UserId int64 `json:"user_id"`
}

func (h *handler) HandleEvent(ctx context.Context, object json.RawMessage) error {
	data, err := object.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed marshal object, %s", err)
	}
	var expired donutExpired
	err = json.Unmarshal(data, &expired)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("unmarshal object to donut expired", "HandleEvent", _PROVIDER))
	}
	err = h.cleaner.CleanSubscription(ctx, expired.UserId)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("clean subscription", "HandleEvent", _PROVIDER))
	}
	return nil
}
