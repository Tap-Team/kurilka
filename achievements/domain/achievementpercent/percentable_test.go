package achievementpercent_test

import (
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/domain/achievementpercent"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestCigarettePercentable(t *testing.T) {
	cases := []struct {
		cigaretteAmount int
		level           int
		percent         int
	}{
		{
			level:           10,
			cigaretteAmount: 60,
			percent:         2,
		},
		{
			level:           9,
			cigaretteAmount: 55,
			percent:         3,
		},
		{
			level:           8,
			cigaretteAmount: 186,
			percent:         13,
		},
		{
			level:           1,
			cigaretteAmount: 10000,
			percent:         100,
		},
	}
	for _, cs := range cases {
		percent := achievementpercent.NewFabric(cs.cigaretteAmount, 0, time.Time{}).Percentable(achievementmodel.CIGARETTE, cs.level).Percent()
		assert.Equal(t, percent, cs.percent, "percent not equal")
	}
}
func timeSubMinutes(minutes int) time.Time {
	seconds := time.Now().Unix() - time.Unix(int64(minutes)*60, 0).Unix()
	return time.Unix(seconds, 0)
}
func TestHealthPercent(t *testing.T) {
	cases := []struct {
		minutes int
		level   int
		percent int
	}{
		{
			level:   10,
			minutes: 10000,
			percent: 4,
		},
		{
			level:   9,
			minutes: 9752,
			percent: 8,
		},
		{
			level:   8,
			minutes: 6800,
			percent: 16,
		},
		{
			level:   1,
			minutes: 6800,
			percent: 100,
		},
	}

	for _, cs := range cases {
		percent := achievementpercent.NewFabric(0, 0, timeSubMinutes(cs.minutes)).Percentable(achievementmodel.HEALTH, cs.level).Percent()
		assert.Equal(t, percent, cs.percent, "percent not equal")
	}
}

func TestSavingPercentable(t *testing.T) {
	cases := []struct {
		money   int
		level   int
		percent int
	}{
		{
			level:   10,
			money:   1875,
			percent: 8,
		},
		{
			level:   9,
			money:   1875,
			percent: 13,
		},
		{
			level:   8,
			money:   777,
			percent: 8,
		},
		{
			level:   1,
			money:   10000,
			percent: 100,
		},
	}
	for _, cs := range cases {
		percent := achievementpercent.NewFabric(0, cs.money, time.Time{}).Percentable(achievementmodel.SAVING, cs.level).Percent()
		assert.Equal(t, percent, cs.percent, "percent not equal")
	}
}

func TestDurationPercentable(t *testing.T) {
	cases := []struct {
		minutes int
		level   int
		percent int
	}{
		{
			level:   10,
			minutes: 19586,
			percent: 3,
		},
		{
			level:   9,
			minutes: 290000,
			percent: 56,
		},
		{
			level:   8,
			minutes: 176352,
			percent: 46,
		},
		{
			level:   1,
			minutes: 100000000,
			percent: 100,
		},
	}
	for _, cs := range cases {
		percent := achievementpercent.NewFabric(0, 0, timeSubMinutes(cs.minutes)).Percentable(achievementmodel.DURATION, cs.level).Percent()
		assert.Equal(t, percent, cs.percent, "percent not equal")
	}
}

func TestWellBeingPercentable(t *testing.T) {
	cases := []struct {
		minutes int
		level   int
		percent int
	}{
		{
			level:   10,
			minutes: 48000,
			percent: 67,
		},
		{
			level:   9,
			minutes: 17234,
			percent: 30,
		},
		{
			level:   8,
			minutes: 6896,
			percent: 14,
		},
		{
			level:   1,
			minutes: 1000000,
			percent: 100,
		},
	}

	for _, cs := range cases {
		percent := achievementpercent.NewFabric(0, 0, timeSubMinutes(cs.minutes)).Percentable(achievementmodel.WELL_BEING, cs.level).Percent()
		require.Equal(t, cs.percent, percent, "percent not equal")
	}
}
