package reachachievementexecutor_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/workers/userworker/achievementmessagesender"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/workers/userworker/executor/reachachievementexecutor"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_Executor_ExecuteUser(t *testing.T) {
	os.Setenv("TZ", time.UTC.String())
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	reacher := reachachievementexecutor.NewMockAchievementUserReacher(ctrl)
	user := userdatamanager.NewMockUserManager(ctrl)
	achievement := achievementdatamanager.NewMockAchievementManager(ctrl)
	sender := achievementmessagesender.NewMockAchievementMessageSenderAtTime(ctrl)

	executor := reachachievementexecutor.New(user, achievement, sender, reacher)

	cases := []struct {
		userData         *model.UserData
		userAchievements []*achievementmodel.Achievement

		userDataCall bool
		userDataErr  error

		userAchievementsCall bool
		userAchievementsErr  error

		reacherCall         bool
		reacherAchievements []int64

		reachAchievementsCall bool
		reachAchievementsErr  error

		err error
	}{
		{
			userDataCall: true,
			userDataErr:  usererror.ExceptionUserNotFound(),
			err:          usererror.ExceptionUserNotFound(),
		},
		{
			userDataCall: true,
			userData:     model.NewUserData(usermodel.PackPrice(1.00), usermodel.CigaretteDayAmount(20), usermodel.CigarettePackAmount(20), time.Now()),

			userAchievementsCall: true,
			userAchievementsErr:  usererror.ExceptionUserNotFound(),
			err:                  usererror.ExceptionUserNotFound(),
		},
		{
			userDataCall: true,
			userData:     model.NewUserData(usermodel.PackPrice(1.00), usermodel.CigaretteDayAmount(20), usermodel.CigarettePackAmount(20), time.Now()),

			userAchievementsCall: true,
			userAchievements:     AchievementList(1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17),

			reacherCall:         true,
			reacherAchievements: []int64{1, 9, 13, 17},

			reachAchievementsCall: true,
			reachAchievementsErr:  usererror.ExceptionUserNotFound(),

			err: usererror.ExceptionUserNotFound(),
		},
		{
			userDataCall: true,
			userData:     model.NewUserData(usermodel.PackPrice(1.00), usermodel.CigaretteDayAmount(20), usermodel.CigarettePackAmount(20), time.Now()),

			userAchievementsCall: true,
			userAchievements:     AchievementList(1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15, 16, 20),

			reacherCall:         true,
			reacherAchievements: []int64{1, 9, 13, 20},

			reachAchievementsCall: true,
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		if cs.userDataCall {
			user.EXPECT().UserData(gomock.Any(), userId).Return(cs.userData, cs.userDataErr).Times(1)
		}
		if cs.userAchievementsCall {
			achievement.EXPECT().UserAchievements(gomock.Any(), userId).Return(cs.userAchievements, cs.userAchievementsErr).Times(1)
		}
		if cs.reacherCall {
			reacher.EXPECT().ReachAchievements(gomock.Any(), userId, cs.userData, cs.userAchievements).Return(cs.reacherAchievements).Times(1)
		}
		if cs.reachAchievementsCall {
			achievement.EXPECT().ReachAchievements(gomock.Any(), userId, TimeSecondsMatcher{time.Now().Unix()}, cs.reacherAchievements).Return(cs.reachAchievementsErr).Times(1)
		}
		sendTime := reachachievementexecutor.NextSendTime(time.Now())
		if cs.reachAchievementsErr == nil {
			for i := 0; i < len(cs.reacherAchievements); i++ {
				sender.EXPECT().SendMessageAtTime(gomock.Any(), userId, gomock.Any(), TimeSecondsMatcher{sendTime.Unix()})
			}
		}
		err := executor.ExecuteUser(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "wrong error")
	}
}

func Test_NextSendTime(t *testing.T) {
	cases := []struct {
		now time.Time

		expected time.Time

		unix int64
	}{
		{
			now:      time.Date(2023, time.April, 1, 11, 1, 1, 1, time.UTC),
			expected: time.Date(2023, time.April, 2, 14, 0, 0, 0, reachachievementexecutor.MoscowLocation),

			unix: time.Date(2023, time.April, 2, 11, 0, 0, 0, time.UTC).Unix(),
		},
		{
			now:      time.Date(2023, time.April, 1, 10, 59, 59, 0, time.UTC),
			expected: time.Date(2023, time.April, 1, 14, 0, 0, 0, reachachievementexecutor.MoscowLocation),

			unix: time.Date(2023, time.April, 1, 11, 0, 0, 0, time.UTC).Unix(),
		},
	}

	for num, cs := range cases {
		actual := reachachievementexecutor.NextSendTime(cs.now)
		log.Println(cs.unix - actual.Unix())
		assert.Equal(t, cs.expected, actual, "wrong time, %d", num)
		assert.Equal(t, cs.unix, actual.Unix(), "unix not equal, %d", num)
	}
}
