package achievementusecase_test

import (
	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	achievementmodel "github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	gomock "github.com/golang/mock/gomock"
)

type UserAchievementsCall struct {
	WillBeCalled bool
	UserId       int64

	UserAchievements []*achievementmodel.Achievement
	Err              error
}

func (c *UserAchievementsCall) RegisterCall(achievementManager *achievementdatamanager.MockAchievementDataManager) {
	if c.WillBeCalled {
		achievementManager.EXPECT().
			UserAchievements(gomock.Any(), c.UserId).
			Return(c.UserAchievements, c.Err).
			Times(1)
	}
}

type UserAchievementsCallBuilder struct {
	userId int64

	userAchievements []*achievementmodel.Achievement
	err              error
}

func (b *UserAchievementsCallBuilder) SetInput(userId int64) *UserAchievementsCallBuilder {
	b.userId = userId
	return b
}
func (b *UserAchievementsCallBuilder) SetOutput(userAchievements []*achievementmodel.Achievement, err error) *UserAchievementsCallBuilder {
	b.userAchievements = userAchievements
	b.err = err
	return b
}

func (b *UserAchievementsCallBuilder) Build() UserAchievementsCall {
	if b == nil {
		return UserAchievementsCall{}
	}
	return UserAchievementsCall{
		UserId:           b.userId,
		UserAchievements: b.userAchievements,
		Err:              b.err,

		WillBeCalled: true,
	}
}
