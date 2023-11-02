package achievementusecase_test

import (
	"context"
	"testing"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"gotest.tools/v3/assert"
)

type OpenSingleAchievementCases struct {
	achievement          *achievementdatamanager.MockAchievementDataManager
	subscriptionProvider *achievementusecase.MockSubscriptionProvider
	achievementStorage   *achievementusecase.MockAchievementStorage
	messageSender        *messagesender.MockMessageSender

	cases []*OpenSingleAchievementCase
}

func NewOpenSingleAchievementCases(
	achievement *achievementdatamanager.MockAchievementDataManager,
	subscriptionProvider *achievementusecase.MockSubscriptionProvider,
	achievementStorage *achievementusecase.MockAchievementStorage,
	messageSender *messagesender.MockMessageSender,
) *OpenSingleAchievementCases {
	return &OpenSingleAchievementCases{
		achievement:          achievement,
		subscriptionProvider: subscriptionProvider,
		achievementStorage:   achievementStorage,
		messageSender:        messageSender,
	}
}

func (c *OpenSingleAchievementCases) AddCase(openSingleAchievementCase *OpenSingleAchievementCase) {
	c.cases = append(c.cases, openSingleAchievementCase)
}

func (c *OpenSingleAchievementCases) Test(t *testing.T, ctx context.Context) {
	achievementUseCase := achievementusecase.New(c.achievement, nil, c.achievementStorage, c.messageSender, c.subscriptionProvider)
	for _, cs := range c.cases {
		cs.Test(t, ctx, achievementUseCase, c.achievement, c.subscriptionProvider, c.achievementStorage, c.messageSender)
	}
}

type OpenSingleAchievementCase struct {
	Description string

	UserId            int64
	OpenAchievementId int64

	OpenAchievementResponse *model.OpenAchievementResponse
	Err                     error

	UserAchievementsCall      UserAchievementsCall
	UserSubscriptionCall      UserSubscriptionCall
	OpenSingleAchievementCall OpenSingleAchievementCall
	AchievementMotivationCall AchievementMotivationCall
	SendMessageCall           SendMessageCall
}

func (c *OpenSingleAchievementCase) registerCalls(
	achievement *achievementdatamanager.MockAchievementDataManager,
	subscriptionProvider *achievementusecase.MockSubscriptionProvider,
	achievementStorage *achievementusecase.MockAchievementStorage,
	messageSender *messagesender.MockMessageSender,
) {
	c.UserAchievementsCall.RegisterCall(achievement)
	c.UserSubscriptionCall.RegisterCall(subscriptionProvider)
	c.OpenSingleAchievementCall.RegisterCall(achievement)
	c.AchievementMotivationCall.RegisterCall(achievementStorage)
	c.SendMessageCall.RegisterCall(messageSender)
}

func (c *OpenSingleAchievementCase) Test(
	t *testing.T,
	ctx context.Context,
	achievementUseCase achievementusecase.AchievementUseCase,
	achievement *achievementdatamanager.MockAchievementDataManager,
	subscriptionProvider *achievementusecase.MockSubscriptionProvider,
	achievementStorage *achievementusecase.MockAchievementStorage,
	messageSender *messagesender.MockMessageSender,
) {
	c.registerCalls(achievement, subscriptionProvider, achievementStorage, messageSender)
	c.assertResult(t, ctx, achievementUseCase)
}

func (c *OpenSingleAchievementCase) assertResult(t *testing.T, ctx context.Context, achievementUseCase achievementusecase.AchievementUseCase) {
	openAchievementResponse, err := achievementUseCase.OpenSingle(ctx, c.UserId, c.OpenAchievementId)
	assert.ErrorIs(t, err, c.Err, "test %s, wrong error", c.Description)
	c.assertOpenAchievementResponse(t, openAchievementResponse)
}

func (c *OpenSingleAchievementCase) assertOpenAchievementResponse(t *testing.T, openAchievementResponse *model.OpenAchievementResponse) {
	if openAchievementResponse != nil {
		assert.Equal(t, openAchievementResponse.OpenTime.Unix(), c.OpenAchievementResponse.OpenTime.Unix(), "wrong response")
	}
}

type OpenSingleAchievementCaseBuilder struct {
	description string

	userId            int64
	openAchievementId int64

	openAchievementResponse *model.OpenAchievementResponse
	err                     error

	userAchievementsCallBuilder      *UserAchievementsCallBuilder
	userSubscriptionCallBuilder      *UserSubscriptionCallBuilder
	openSingleAchievementCallBuilder *OpenSingleAchievementCallBuilder
	achievementMotivationCallBuilder *AchievementMotivationCallBuilder
	sendMessageCallBuilder           *SendMessageCallBuilder
}

func (b *OpenSingleAchievementCaseBuilder) SetDescription(description string) *OpenSingleAchievementCaseBuilder {
	b.description = description
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetInput(userId, achievementId int64) *OpenSingleAchievementCaseBuilder {
	b.userId = userId
	b.openAchievementId = achievementId
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetOutput(openAchievementResponse *model.OpenAchievementResponse, err error) *OpenSingleAchievementCaseBuilder {
	b.openAchievementResponse = openAchievementResponse
	b.err = err
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetUserAchievementsCallBuilder(builder *UserAchievementsCallBuilder) *OpenSingleAchievementCaseBuilder {
	b.userAchievementsCallBuilder = builder
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetUserSubscriptionCallBuilder(builder *UserSubscriptionCallBuilder) *OpenSingleAchievementCaseBuilder {
	b.userSubscriptionCallBuilder = builder
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetOpenSingleAchievementCallBuilder(builder *OpenSingleAchievementCallBuilder) *OpenSingleAchievementCaseBuilder {
	b.openSingleAchievementCallBuilder = builder
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetAchievementMotivationCallBuilder(builder *AchievementMotivationCallBuilder) *OpenSingleAchievementCaseBuilder {
	b.achievementMotivationCallBuilder = builder
	return b
}

func (b *OpenSingleAchievementCaseBuilder) SetSendMessageCallBuilder(builder *SendMessageCallBuilder) *OpenSingleAchievementCaseBuilder {
	b.sendMessageCallBuilder = builder
	return b
}

func (b *OpenSingleAchievementCaseBuilder) Build() *OpenSingleAchievementCase {
	return &OpenSingleAchievementCase{
		Description:               b.description,
		UserId:                    b.userId,
		OpenAchievementId:         b.openAchievementId,
		OpenAchievementResponse:   b.openAchievementResponse,
		Err:                       b.err,
		UserAchievementsCall:      b.userAchievementsCallBuilder.Build(),
		UserSubscriptionCall:      b.userSubscriptionCallBuilder.Build(),
		OpenSingleAchievementCall: b.openSingleAchievementCallBuilder.Build(),
		AchievementMotivationCall: b.achievementMotivationCallBuilder.Build(),
		SendMessageCall:           b.sendMessageCallBuilder.Build(),
	}
}
