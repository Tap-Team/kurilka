package achievementusecase_test

import (
	context "context"
	"errors"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager"
	"github.com/Tap-Team/kurilka/achievements/model"
	"github.com/Tap-Team/kurilka/achievements/usecase/achievementusecase"
	"github.com/Tap-Team/kurilka/internal/errorutils/achievementerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	gomock "github.com/golang/mock/gomock"
)

var (
	openDateNow   = amidtime.Timestamp{Time: time.Now()}
	emptyOpenDate = amidtime.Timestamp{}

	reachDateNow   = amidtime.Timestamp{Time: time.Now()}
	emptyReachDate = amidtime.Timestamp{}
)

var achievements = []*achievementmodel.Achievement{
	achievementmodel.NewAchievement(1, achievementmodel.HEALTH, 1, 20, openDateNow, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(2, achievementmodel.HEALTH, 2, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(3, achievementmodel.HEALTH, 3, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(4, achievementmodel.HEALTH, 4, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(5, achievementmodel.HEALTH, 5, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(6, achievementmodel.HEALTH, 6, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(7, achievementmodel.HEALTH, 7, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(8, achievementmodel.HEALTH, 8, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(9, achievementmodel.HEALTH, 9, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(10, achievementmodel.HEALTH, 10, 20, emptyOpenDate, emptyOpenDate, true, 100, ""),

	achievementmodel.NewAchievement(11, achievementmodel.CIGARETTE, 1, 20, openDateNow, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(12, achievementmodel.CIGARETTE, 2, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(13, achievementmodel.CIGARETTE, 3, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(14, achievementmodel.CIGARETTE, 4, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(15, achievementmodel.CIGARETTE, 5, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(16, achievementmodel.CIGARETTE, 6, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(17, achievementmodel.CIGARETTE, 7, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(18, achievementmodel.CIGARETTE, 8, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(19, achievementmodel.CIGARETTE, 9, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(20, achievementmodel.CIGARETTE, 10, 20, emptyOpenDate, emptyOpenDate, true, 100, ""),

	achievementmodel.NewAchievement(21, achievementmodel.WELL_BEING, 1, 20, openDateNow, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(22, achievementmodel.WELL_BEING, 2, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(23, achievementmodel.WELL_BEING, 3, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(24, achievementmodel.WELL_BEING, 4, 20, emptyOpenDate, reachDateNow, true, 100, ""),
	achievementmodel.NewAchievement(25, achievementmodel.WELL_BEING, 5, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(26, achievementmodel.WELL_BEING, 6, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(27, achievementmodel.WELL_BEING, 7, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(28, achievementmodel.WELL_BEING, 8, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(29, achievementmodel.WELL_BEING, 9, 20, emptyOpenDate, emptyReachDate, true, 100, ""),
	achievementmodel.NewAchievement(30, achievementmodel.WELL_BEING, 10, 20, emptyOpenDate, emptyOpenDate, true, 100, ""),
}

func Test_UseCase_OpenSingle(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	achievement := achievementdatamanager.NewMockAchievementManager(ctrl)
	achievementStorage := achievementusecase.NewMockAchievementStorage(ctrl)
	messageSender := messagesender.NewMockMessageSender(ctrl)
	subscriptionProvider := achievementusecase.NewMockSubscriptionProvider(ctrl)

	cases := NewOpenSingleAchievementCases(achievement, subscriptionProvider, achievementStorage, messageSender)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("при получение ачивок пользователя у нас возвращается ошибка, метод должен вернуть ошибку не удалось открыть ачивку").
			SetInput(23181, 21).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(23181).
					SetOutput(nil, usererror.ExceptionUserNotFound()),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("пробуем открыть ачивку которую по идее не можем открыть, потому что есть ачивка которая ещё не открыта с меньшим уровнем").
			SetInput(1485128, 23).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(1485128).
					SetOutput(achievements, nil),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("пробуем открыть ачивку которую по идее не можем открыть, потому что есть ачивка которая ещё не открыта с меньшим уровнем").
			SetInput(1485128, 24).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(1485128).
					SetOutput(achievements, nil),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("если ачивка относится к типу Сигареты или Здоровье то мы должны проверить пользовательскую подписку, в данном случае проверяем ачивку типа Сигареты, метод должен вернуть ошибку").
			SetInput(132332, 12).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(132332).
					SetOutput(achievements, nil),
			).
			SetUserSubscriptionCallBuilder(
				new(UserSubscriptionCallBuilder).
					SetInput(132332).
					SetOutput(usermodel.Subscription{}, usererror.ExceptionUserNotFound()),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("если ачивка относится к типу Сигареты или Здоровье то мы должны проверить пользовательскую подписку, в данном случае проверяем ачивку типа Здоровье, метод должен вернуть ошибку").
			SetInput(132332, 2).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(132332).
					SetOutput(achievements, nil),
			).
			SetUserSubscriptionCallBuilder(
				new(UserSubscriptionCallBuilder).
					SetInput(132332).
					SetOutput(usermodel.Subscription{}, usererror.ExceptionUserNotFound()),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("если ачивка относится к типу Сигареты или Здоровье то мы должны проверить пользовательскую подписку, в данном случае проверяем ачивку типа Здоровье, метод должен вернуть ошибку поскольку подписка недействительна").
			SetInput(132332, 2).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(132332).
					SetOutput(achievements, nil),
			).
			SetUserSubscriptionCallBuilder(
				new(UserSubscriptionCallBuilder).
					SetInput(132332).
					SetOutput(usermodel.NewSubscription(usermodel.NONE, time.Now().Add(time.Hour)), nil),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("если ачивка относится к типу Сигареты или Здоровье то мы должны проверить пользовательскую подписку, в данном случае проверяем ачивку типа Здоровье, метод должен вернуть ошибку поскольку подписка недействительна").
			SetInput(132332, 2).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(132332).
					SetOutput(achievements, nil),
			).
			SetUserSubscriptionCallBuilder(
				new(UserSubscriptionCallBuilder).
					SetInput(132332).
					SetOutput(usermodel.NewSubscription(usermodel.BASIC, time.Time{}), nil),
			).
			SetOutput(nil, achievementerror.ExceptionCantOpenAchievementForUser()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("пробуем открыть ачивку которая удовлетворяет всем условиям проверки, метод должен пойти дальше и попробовать открыть ачивку, которая по итогу возвращает ошибку").
			SetInput(148528, 22).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(148528).
					SetOutput(achievements, nil),
			).
			SetOpenSingleAchievementCallBuilder(
				new(OpenSingleAchievementCallBuilder).
					SetInput(148528, 22).
					SetOutput(model.NewOpenAchievementResponse(time.Now()), usererror.ExceptionUserNotFound()),
			).
			SetOutput(nil, usererror.ExceptionUserNotFound()).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("метод доходит до открытия ачивки, ачивка отрылась успешно, далее мы должны перейти в отправке сообщения, вне зависимости от результата операции метод завершается успешно").
			SetInput(14828, 22).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(14828).
					SetOutput(achievements, nil),
			).
			SetOpenSingleAchievementCallBuilder(
				new(OpenSingleAchievementCallBuilder).
					SetInput(14828, 22).
					SetOutput(model.NewOpenAchievementResponse(time.Now()), nil),
			).
			SetAchievementMotivationCallBuilder(
				new(AchievementMotivationCallBuilder).
					SetInput(22).
					SetOutput("Ухх ебать, мотивация", errors.New("choto ne poperlo")),
			).
			SetOutput(model.NewOpenAchievementResponse(time.Now()), nil).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("метод доходит до открытия ачивки, ачивка отрылась успешно, далее мы должны перейти в отправке сообщения, вне зависимости от результата операции метод завершается успешно").
			SetInput(14828, 22).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(14828).
					SetOutput(achievements, nil),
			).
			SetOpenSingleAchievementCallBuilder(
				new(OpenSingleAchievementCallBuilder).
					SetInput(14828, 22).
					SetOutput(model.NewOpenAchievementResponse(time.Now()), nil),
			).
			SetAchievementMotivationCallBuilder(
				new(AchievementMotivationCallBuilder).
					SetInput(22).
					SetOutput("Ухх ебать, мотивация", nil),
			).
			SetSendMessageCallBuilder(
				new(SendMessageCallBuilder).
					SetInput("Ухх ебать, мотивация", 14828).
					SetOutput(errors.New("choto ne poperlo")),
			).
			SetOutput(model.NewOpenAchievementResponse(time.Now()), nil).
			Build(),
	)

	cases.AddCase(
		new(OpenSingleAchievementCaseBuilder).
			SetDescription("метод доходит до открытия ачивки, ачивка отрылась успешно, далее мы должны перейти в отправке сообщения, вне зависимости от результата операции метод завершается успешно").
			SetInput(14828, 22).
			SetUserAchievementsCallBuilder(
				new(UserAchievementsCallBuilder).
					SetInput(14828).
					SetOutput(achievements, nil),
			).
			SetOpenSingleAchievementCallBuilder(
				new(OpenSingleAchievementCallBuilder).
					SetInput(14828, 22).
					SetOutput(model.NewOpenAchievementResponse(time.Now()), nil),
			).
			SetAchievementMotivationCallBuilder(
				new(AchievementMotivationCallBuilder).
					SetInput(22).
					SetOutput("Ухх ебать, мотивация", nil),
			).
			SetSendMessageCallBuilder(
				new(SendMessageCallBuilder).
					SetInput("Ухх ебать, мотивация", 14828).
					SetOutput(nil),
			).
			SetOutput(model.NewOpenAchievementResponse(time.Now()), nil).
			Build(),
	)

	cases.Test(t, ctx)
}
