package statisticsusecase

import (
	"context"

	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/userdatamanager"
	"github.com/Tap-Team/kurilka/user/model"
)

//go:generate mockgen -source usecase.go -destination mocks.go -package statisticsusecase

const _PROVIDER = "user/usecase/statisticsusecase"

type StatisticsUseCase interface {
	TimeStatistics(ctx context.Context, userId int64) (model.IntUserStatistics, error)
	CigaretteStatistics(ctx context.Context, userId int64) (model.IntUserStatistics, error)
	MoneyStatistics(ctx context.Context, userId int64) (model.FloatUserStatistics, error)
}

type statisticsUseCase struct {
	user userdatamanager.UserManager
}

func New(user userdatamanager.UserManager) StatisticsUseCase {
	return &statisticsUseCase{user: user}
}

func (s *statisticsUseCase) TimeStatistics(ctx context.Context, userId int64) (model.IntUserStatistics, error) {
	var statistics model.IntUserStatistics
	user, err := s.user.User(ctx, userId)
	if err != nil {
		return statistics, exception.Wrap(err, exception.NewCause("get user data", "TimeStatistics", _PROVIDER))
	}
	statistics = model.NewIntUserStatistics(int(user.CigaretteDayAmount) * 5)
	return statistics, nil
}

func (s *statisticsUseCase) CigaretteStatistics(ctx context.Context, userId int64) (model.IntUserStatistics, error) {
	var statistics model.IntUserStatistics
	user, err := s.user.User(ctx, userId)
	if err != nil {
		return statistics, exception.Wrap(err, exception.NewCause("get user data", "TimeStatistics", _PROVIDER))
	}
	statistics = model.NewIntUserStatistics(int(user.CigaretteDayAmount))
	return statistics, nil
}

func (s *statisticsUseCase) MoneyStatistics(ctx context.Context, userId int64) (model.FloatUserStatistics, error) {
	var statistics model.FloatUserStatistics
	user, err := s.user.User(ctx, userId)
	if err != nil {
		return statistics, exception.Wrap(err, exception.NewCause("get user data", "TimeStatistics", _PROVIDER))
	}
	cigarettePrice := float64(user.PackPrice) / float64(user.CigarettePackAmount)
	statistics = model.NewFloatUserStatisctics(cigarettePrice * float64(user.CigaretteDayAmount))
	return statistics, nil
}
