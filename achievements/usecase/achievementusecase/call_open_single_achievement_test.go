package achievementusecase_test

import (
	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/golang/mock/gomock"
)

type OpenSingleAchievementCall struct {
	WillBeCalled  bool
	UserId        int64
	AchievementId int64

	OpenAchievementResponse *model.OpenAchievementResponse
	Err                     error
}

func (c *OpenSingleAchievementCall) RegisterCall(achievementManager *achievementdatamanager.MockAchievementDataManager) {
	if c.WillBeCalled {
		achievementManager.EXPECT().
			OpenSingle(gomock.Any(), c.UserId, c.AchievementId).
			Return(c.OpenAchievementResponse, c.Err).
			Times(1)
	}
}

type OpenSingleAchievementCallBuilder struct {
	userId        int64
	achievementId int64

	openAchievementResponse *model.OpenAchievementResponse
	err                     error
}

func (b *OpenSingleAchievementCallBuilder) SetInput(userId, achievementId int64) *OpenSingleAchievementCallBuilder {
	b.userId = userId
	b.achievementId = achievementId
	return b
}

func (b *OpenSingleAchievementCallBuilder) SetOutput(openAchievementResponse *model.OpenAchievementResponse, err error) *OpenSingleAchievementCallBuilder {
	b.err = err
	b.openAchievementResponse = openAchievementResponse
	return b
}

func (b *OpenSingleAchievementCallBuilder) Build() OpenSingleAchievementCall {
	if b == nil {
		return OpenSingleAchievementCall{}
	}
	return OpenSingleAchievementCall{
		UserId:        b.userId,
		AchievementId: b.achievementId,
		Err:           b.err,

		WillBeCalled: true,
	}
}
