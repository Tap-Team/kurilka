package achievementusecase_test

import (
	context "context"
	"errors"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"gotest.tools/v3/assert"

	"github.com/golang/mock/gomock"
)

var (
	NilOpenResponse *model.OpenAchievementResponse
)

func Test_UseCase_MarkShown(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievement := achievementdatamanager.NewMockAchievementManager(ctrl)
	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := achievementusecase.New(achievement, user, nil, nil, nil)

	{
		userId := rand.Int63()
		expectedErr := errors.New("random error")

		achievement.EXPECT().MarkShown(gomock.Any(), userId).Return(expectedErr).Times(1)

		err := useCase.MarkShown(ctx, userId)
		assert.ErrorIs(t, err, expectedErr, "error not equal")
	}

	{
		userId := rand.Int63()

		achievement.EXPECT().MarkShown(gomock.Any(), userId).Return(nil).Times(1)

		err := useCase.MarkShown(ctx, userId)
		assert.NilError(t, err, "non nil error")
	}

}

func Test_UseCase_UserAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievement := achievementdatamanager.NewMockAchievementManager(ctrl)
	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := achievementusecase.New(achievement, user, nil, nil, nil)

	{
		userId := rand.Int63()
		expectedErr := errors.New("failed get achievement data")

		achievement.EXPECT().UserAchievements(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		achievements, err := useCase.UserAchievements(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
		assert.Equal(t, 0, len(achievements), "wrong achievements")
	}

	{
		userId := rand.Int63()
		expectedErr := errors.New("failed get user data")
		userAchievements := generateRandomAchievementList(50)

		achievement.EXPECT().UserAchievements(gomock.Any(), userId).Return(userAchievements, nil).Times(1)
		user.EXPECT().UserData(gomock.Any(), userId).Return(nil, expectedErr).Times(1)

		achievements, err := useCase.UserAchievements(ctx, userId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
		assert.Equal(t, 0, len(achievements), "wrong achievements")
	}

}

func Test_UseCase_UserReachedAchievements(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := achievementdatamanager.NewMockAchievementManager(ctrl)

	useCase := achievementusecase.New(manager, nil, nil, nil, nil)

	cases := []struct {
		achievements    []*achievementmodel.Achievement
		achievementsErr error

		err error

		reachedAchievements model.ReachedAchievements
	}{
		{
			achievementsErr: usererror.ExceptionUserNotFound(),
			err:             usererror.ExceptionUserNotFound(),
		},
		{
			achievements: []*achievementmodel.Achievement{
				NewAchievement(false, false, achievementmodel.WELL_BEING),
				NewAchievement(true, true, achievementmodel.WELL_BEING),
				NewAchievement(true, false, achievementmodel.HEALTH),
				NewAchievement(true, false, achievementmodel.WELL_BEING),

				NewAchievement(true, false, achievementmodel.SAVING),
				NewAchievement(true, false, achievementmodel.HEALTH),

				NewAchievement(false, true, achievementmodel.WELL_BEING),
				NewAchievement(false, true, achievementmodel.SAVING),
				NewAchievement(false, true, achievementmodel.DURATION),
				NewAchievement(false, true, achievementmodel.HEALTH),
				NewAchievement(false, true, achievementmodel.HEALTH),
				NewAchievement(false, true, achievementmodel.CIGARETTE),
				NewAchievement(false, true, achievementmodel.CIGARETTE),
				NewAchievement(false, true, achievementmodel.CIGARETTE),
			},
			reachedAchievements: model.ReachedAchievements{
				Saving:    1,
				WellBeing: 1,
				Health:    2,
				Cigarette: 3,
				Duration:  1,
				Type:      achievementmodel.CIGARETTE,
			},
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		manager.EXPECT().UserAchievements(gomock.Any(), userId).Return(cs.achievements, cs.achievementsErr).Times(1)

		rach, err := useCase.UserReachedAchievements(ctx, userId)
		assert.ErrorIs(t, err, cs.err, "wrong err")
		assert.Equal(t, rach, cs.reachedAchievements, "reach achievements not equal")
	}
}
