package workers_test

import (
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/workers/userworker/workers"
	"gotest.tools/v3/assert"
)

func TestExecuteUserTime(t *testing.T) {
	os.Setenv("TZ", time.UTC.String())
	cases := []struct {
		now                time.Time
		userAbstinenceTime time.Time

		executePause time.Duration

		expected time.Time
	}{
		{
			now:                time.Date(2023, time.September, 1, 15, 0, 0, 0, time.UTC),
			userAbstinenceTime: time.Date(2023, time.August, 18, 16, 1, 2, 0, time.UTC),

			executePause: time.Hour * 24,

			expected: time.Date(2023, time.September, 1, 16, 1, 2, 0, time.UTC),
		},
		{
			now:                time.Date(2023, time.September, 2, 17, 0, 0, 0, time.UTC),
			userAbstinenceTime: time.Date(2023, time.August, 18, 16, 1, 2, 0, time.UTC),

			executePause: time.Hour * 24,

			expected: time.Date(2023, time.September, 3, 16, 1, 2, 0, time.UTC),
		},
		{
			now:                time.Date(2023, time.September, 26, 23, 33, 0, 0, time.UTC),
			userAbstinenceTime: time.Date(2023, time.September, 26, 23, 33, 0, 0, time.UTC),

			executePause: time.Hour * 48,
			expected:     time.Date(2023, time.September, 28, 23, 33, 0, 0, time.UTC),
		},
		{
			now:                time.Date(2023, time.September, 27, 0, 14, 0, 0, time.UTC),
			userAbstinenceTime: time.Date(2023, time.September, 27, 0, 11, 0, 0, time.UTC),

			executePause: time.Hour,

			expected: time.Date(2023, time.September, 27, 1, 11, 0, 0, time.UTC),
		},
		{
			now:                time.Date(2023, time.September, 27, 0, 14, 0, 0, time.UTC),
			userAbstinenceTime: time.Date(2023, time.September, 26, 0, 20, 0, 0, time.UTC),

			executePause: time.Hour,

			expected: time.Date(2023, time.September, 27, 0, 20, 0, 0, time.UTC),
		},
	}

	for _, cs := range cases {
		userTimeExecutor := workers.NewUserTimeExecutor(cs.executePause)
		actual := userTimeExecutor.CountUserExecuteTime(cs.now, cs.userAbstinenceTime)
		assert.Equal(t, actual, cs.expected, "wrong time")
	}
}
