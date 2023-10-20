package subscription

import "time"

type Period int

const (
	WEEK   Period = 7
	MONTH  Period = 30
	SEASON Period = 90
)

type TrialDuration int

const (
	THREE  TrialDuration = 3
	SEVEN  TrialDuration = 7
	THIRTY TrialDuration = 30
)

type Subscription struct {
	ID            string        `json:"item_id,omitempty"`
	Title         string        `json:"title"`
	Photo         string        `json:"photo_url,omitempty"`
	Price         int           `json:"price"`
	Period        Period        `json:"period"`
	TrialDuration TrialDuration `json:"trial_duration,omitempty"`
	Expiration    int           `json:"expiration,omitempty"`
}

func (s *Subscription) SetExpiration(exp time.Duration) {
	s.Expiration = int(exp.Seconds())
}
