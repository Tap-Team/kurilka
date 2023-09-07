package achievementusecase_test

import (
	context "context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"gotest.tools/v3/assert"

	"github.com/golang/mock/gomock"
)

var (
	NilOpenResponse *model.OpenAchievementResponse
)

func Test_UseCase_OpenSingle(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievement := achievementdatamanager.NewMockAchievementManager(ctrl)
	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := achievementusecase.New(achievement, user)

	{
		userId := rand.Int63()
		achievementId := rand.Int63()
		expectedErr := errors.New("fatal error")
		achievement.EXPECT().OpenSingle(gomock.Any(), userId, achievementId).Return(nil, expectedErr).Times(1)

		response, err := useCase.OpenSingle(ctx, userId, achievementId)

		assert.ErrorIs(t, err, expectedErr, "wrong error")
		assert.Equal(t, response, NilOpenResponse, "wrong open response")
	}

	{
		userId := rand.Int63()
		achievementId := rand.Int63()
		openResponse := model.NewOpenAchievementResponse(time.Now())
		achievement.EXPECT().OpenSingle(gomock.Any(), userId, achievementId).Return(openResponse, nil).Times(1)

		response, err := useCase.OpenSingle(ctx, userId, achievementId)

		assert.NilError(t, err, "non nil error")
		assert.Equal(t, openResponse, response, "response not equal")
	}
}

func Test_UseCase_MarkShown(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievement := achievementdatamanager.NewMockAchievementManager(ctrl)
	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := achievementusecase.New(achievement, user)

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

	useCase := achievementusecase.New(achievement, user)

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

	{

	}
}
