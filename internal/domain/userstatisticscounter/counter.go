package userstatisticscounter

import "time"

type AccuracyDegree int

const (
	Day AccuracyDegree = iota
	Minute
	Second
)

type Counter struct {
	noSmokingTime time.Duration

	accuracyDegree AccuracyDegree

	cigaretteDayAmount  int
	cigarettePackAmount int
	packPrice           float64
}

func NewCounter(now, userAbstinenceTime time.Time, cigaretteDayAmount, cigarettePackAmount int, packPrice float64, accuracyDegree AccuracyDegree) Counter {
	return Counter{
		noSmokingTime:       now.Sub(userAbstinenceTime),
		cigaretteDayAmount:  cigaretteDayAmount,
		cigarettePackAmount: cigaretteDayAmount,
		packPrice:           packPrice,
		accuracyDegree:      accuracyDegree,
	}
}

func (u Counter) index() float64 {
	switch u.accuracyDegree {
	case Day:
		return float64(int(u.noSmokingTime.Hours() / 24))
	case Minute:
		return u.noSmokingTime.Minutes() / 1440
	case Second:
		return u.noSmokingTime.Seconds() / 86400
	default:
		return float64(int(u.noSmokingTime.Hours() / 24))
	}
}

func (u Counter) Cigarette() int {
	return int(u.index() * float64(u.cigaretteDayAmount))
}

func (u Counter) Life() int {
	return int(u.index() * 20)
}

func (u Counter) Time() int {
	timePerDay := u.cigaretteDayAmount * 5
	return int(float64(timePerDay) * u.index())
}

func (u Counter) Money() float64 {
	cigaretteCost := float64(u.packPrice) / float64(u.cigarettePackAmount)
	moneyPerDay := cigaretteCost * float64(u.cigaretteDayAmount)
	return moneyPerDay * u.index()
}
