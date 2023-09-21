package statisticsusecase_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/user/model"
	"github.com/Tap-Team/kurilka/user/usecase/statisticsusecase"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_UseCase_MoneyStatistics(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := statisticsusecase.New(user)

	cases := []struct {
		user       *usermodel.UserData
		managerErr error

		err        error
		statistics model.FloatUserStatistics
	}{
		{
			managerErr: usererror.ExceptionUserNotFound(),
			err:        usererror.ExceptionUserNotFound(),
		},
		{
			user:       moneyUserData(168, 20, 15),
			statistics: model.NewFloatUserStatisctics(126),
		},
		{
			user:       moneyUserData(244, 20, 8),
			statistics: model.NewFloatUserStatisctics(97.6),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		user.EXPECT().User(gomock.Any(), userId).Return(cs.user, cs.managerErr).Times(1)

		statistics, err := useCase.MoneyStatistics(ctx, userId)
		assert.Equal(t, true, floatStatisticsEqual(statistics, cs.statistics), "statistics not equal")
		assert.ErrorIs(t, err, cs.err, "error not equal")
	}
}

func Test_UseCase_TimeStatistics(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := statisticsusecase.New(user)

	cases := []struct {
		user       *usermodel.UserData
		managerErr error

		err        error
		statistics model.IntUserStatistics
	}{
		{
			managerErr: usererror.ExceptionUserNotFound(),
			err:        usererror.ExceptionUserNotFound(),
		},
		{
			user:       timeUserData(100),
			statistics: model.NewIntUserStatistics(100 * 5),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		user.EXPECT().User(gomock.Any(), userId).Return(cs.user, cs.managerErr).Times(1)

		statistics, err := useCase.TimeStatistics(ctx, userId)
		assert.Equal(t, statistics, cs.statistics, "statistics not equal")
		assert.ErrorIs(t, err, cs.err, "error not equal")
	}

}

func Test_UseCase_CigaretteStatistics(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	user := userdatamanager.NewMockUserManager(ctrl)

	useCase := statisticsusecase.New(user)

	cases := []struct {
		user       *usermodel.UserData
		managerErr error

		err        error
		statistics model.IntUserStatistics
	}{
		{
			managerErr: usererror.ExceptionUserNotFound(),
			err:        usererror.ExceptionUserNotFound(),
		},
		{
			user:       timeUserData(100),
			statistics: model.NewIntUserStatistics(100),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()
		user.EXPECT().User(gomock.Any(), userId).Return(cs.user, cs.managerErr).Times(1)

		statistics, err := useCase.CigaretteStatistics(ctx, userId)
		assert.Equal(t, statistics, cs.statistics, "statistics not equal")
		assert.ErrorIs(t, err, cs.err, "error not equal")
	}
}
