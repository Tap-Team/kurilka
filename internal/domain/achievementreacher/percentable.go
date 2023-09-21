package achievementreacher

import (
	"math"
	"time"
)

type CigarettePercentable struct {
	cigaretteAmount int
	level           int
}

type CigaretteGoal int

func (c CigaretteGoal) Goal() int {
	goal := 0
	switch c {
	case 1:
		goal = 20
	case 2:
		goal = 100
	case 3:
		goal = 150
	case 4:
		goal = 250
	case 5:
		goal = 500
	case 6:
		goal = 750
	case 7:
		goal = 1000
	case 8:
		goal = 1500
	case 9:
		goal = 2000
	case 10:
		goal = 3000
	}
	return goal
}

func (c CigarettePercentable) Goal() int {
	return CigaretteGoal(c.level).Goal()
}

func (c CigarettePercentable) Percent() int {
	goal := c.Goal()
	percent := float64(c.cigaretteAmount*100) / float64(goal)
	percent = math.Ceil(percent)
	return min(100, int(percent))
}

type HealthPercentable struct {
	level          int
	abstinenceTime time.Time
}

const hour = 60
const day = 24 * hour
const week = 7 * day
const month = 30 * day
const year = 12*month + 5*day

type HealthGoal int

func (h HealthGoal) Goal() int {
	goal := 0
	switch h {
	case 1:
		goal = 20
	case 2:
		goal = 8 * hour
	case 3:
		goal = 9 * hour
	case 4:
		goal = day
	case 5:
		goal = 33 * hour
	case 6:
		goal = 2 * day
	case 7:
		goal = 3 * week
	case 8:
		goal = month
	case 9:
		goal = 3 * month
	case 10:
		goal = 6 * month
	}
	return goal
}

func (h HealthPercentable) Goal() int {
	return HealthGoal(h.level).Goal()
}

func (h HealthPercentable) Percent() int {
	goal := float64(h.Goal())
	minutes := time.Now().Sub(h.abstinenceTime).Minutes()
	percent := minutes * 100 / goal
	percent = math.Ceil(percent)
	return min(100, int(percent))
}

type SavingPercentable struct {
	level int
	money int
}

type SavingGoal int

func (s SavingGoal) Goal() int {
	goal := 0
	switch s {
	case 1:
		goal = 300
	case 2:
		goal = 600
	case 3:
		goal = 1200
	case 4:
		goal = 2000
	case 5:
		goal = 3000
	case 6:
		goal = 4000
	case 7:
		goal = 6000
	case 8:
		goal = 10000
	case 9:
		goal = 15000
	case 10:
		goal = 25000
	}
	return goal
}

func (s SavingPercentable) Goal() int {
	return SavingGoal(s.level).Goal()
}

func (s SavingPercentable) Percent() int {
	goal := float64(s.Goal())
	percent := float64(s.money) * 100 / goal
	percent = math.Ceil(percent)
	return min(100, int(percent))
}

type DurationPercentable struct {
	abstinenceTime time.Time
	level          int
}

type DurationGoal int

func (d DurationGoal) Goal() int {
	goal := 0
	switch d {
	case 1:
		goal = day
	case 2:
		goal = 3 * day
	case 3:
		goal = week
	case 4:
		goal = month
	case 5:
		goal = 2 * month
	case 6:
		goal = 3 * month
	case 7:
		goal = 6 * month
	case 8:
		goal = 9 * month
	case 9:
		goal = year
	case 10:
		goal = year + 6*month
	}
	return goal
}

func (d DurationPercentable) Goal() int {
	return DurationGoal(d.level).Goal()
}

func (d DurationPercentable) Percent() int {
	goal := float64(d.Goal())
	minutes := time.Now().Sub(d.abstinenceTime).Minutes()
	percent := minutes * 100 / goal
	percent = math.Ceil(percent)
	return min(100, int(percent))
}

type WellBeingPercentable struct {
	abstinenceTime time.Time
	level          int
}

type WellBeingGoal int

func (w WellBeingGoal) Goal() int {
	goal := 0
	switch w {
	case 1:
		goal = 5 * day
	case 2:
		goal = 7 * day
	case 3:
		goal = 10 * day
	case 4:
		goal = 14 * day
	case 5:
		goal = 17 * day
	case 6:
		goal = 25 * day
	case 7:
		goal = 31 * day
	case 8:
		goal = 35 * day
	case 9:
		goal = 40 * day
	case 10:
		goal = 50 * day
	}
	return goal
}

func (w WellBeingPercentable) Goal() int {
	return WellBeingGoal(w.level).Goal()
}

func (w WellBeingPercentable) Percent() int {
	goal := float64(w.Goal())
	minutes := time.Now().Sub(w.abstinenceTime).Minutes()
	percent := minutes * 100 / goal
	percent = math.Ceil(percent)
	return min(100, int(percent))
}
