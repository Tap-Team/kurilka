package donutprolonged

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

//go:generate mockgen -source handler.go -destination mocks.go -package donutprolonged

const _PROVIDER = "callback/handler/donutprolonged.handler"

type Prolongationer interface {
	ProlongSubscription(ctx context.Context, userId int64, amount float64) error
}

type handler struct {
	prolong Prolongationer
}

func New(prolong Prolongationer) *handler {
	return &handler{
		prolong: prolong,
	}
}

type donutProlong struct {
	Amount           float64 `json:"amount"`
	AmountWithoutFee float64 `json:"amount_without_fee"`
	UserId           float64 `json:"user_id"`
}

func (h *handler) HandleEvent(ctx context.Context, object json.RawMessage) error {
	data, err := object.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed marshal object, %s", err)
	}
	var prolong donutProlong
	err = json.Unmarshal(data, &prolong)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("unmarshal object to donut prolong", "HandleEvent", _PROVIDER))
	}
	err = h.prolong.ProlongSubscription(ctx, int64(prolong.UserId), prolong.Amount)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("prolong subscription", "HandleEvent", _PROVIDER))
	}
	return nil
}
