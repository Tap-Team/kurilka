package donutcreate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source handler.go -destination mocks.go -package donutcreate

const _PROVIDER = "callback/handler/donutcreate.handler"

type Creator interface {
	CreateSubscription(ctx context.Context, userId int64, amount int) error
}

type handler struct {
	creator Creator
}

func New(creator Creator) *handler {
	return &handler{
		creator: creator,
	}
}

type donutCreate struct {
	Amount           int     `json:"amount"`
	AmountWithoutFee float32 `json:"amount_without_fee"`
	UserId           int64   `json:"user_id"`
}

func (h *handler) HandleEvent(ctx context.Context, object json.RawMessage) error {
	data, err := object.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed marshal object, %s", err)
	}
	var create donutCreate
	err = json.Unmarshal(data, &create)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("unmarshal object to donut create object", "HandleEvent", _PROVIDER))
	}
	err = h.creator.CreateSubscription(ctx, create.UserId, create.Amount)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("create subscription", "HandleEvent", _PROVIDER))
	}
	return nil
}
