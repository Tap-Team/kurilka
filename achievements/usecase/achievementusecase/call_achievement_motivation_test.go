package achievementusecase_test

import (
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/golang/mock/gomock"
)

type AchievementMotivationCall struct {
	WillBeCalled  bool
	AchievementId int64

	Motivation string
	Err        error
}

func (c *AchievementMotivationCall) RegisterCall(achievementStorage *achievementusecase.MockAchievementStorage) {
	if c.WillBeCalled {
		achievementStorage.EXPECT().
			AchievementMotivation(gomock.Any(), c.AchievementId).
			Return(c.Motivation, c.Err).
			Times(1)
	}
}

type AchievementMotivationCallBuilder struct {
	achievementId int64

	motivation string
	err        error
}

func (b *AchievementMotivationCallBuilder) SetInput(achievementId int64) *AchievementMotivationCallBuilder {
	b.achievementId = achievementId
	return b
}

func (b *AchievementMotivationCallBuilder) SetOutput(motivation string, err error) *AchievementMotivationCallBuilder {
	b.motivation = motivation
	b.err = err
	return b
}

func (b *AchievementMotivationCallBuilder) Build() AchievementMotivationCall {
	if b == nil {
		return AchievementMotivationCall{}
	}
	return AchievementMotivationCall{
		AchievementId: b.achievementId,
		Motivation:    b.motivation,
		Err:           b.err,

		WillBeCalled: true,
	}
}
