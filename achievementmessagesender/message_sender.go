package achievementmessagesender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Tap-Team/kurilka/achievementmessagesender/model"
	"github.com/Tap-Team/kurilka/internal/errorutils/vkerror"
	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/vk/message"
)

//go:generate mockgen -source message_sender.go -destination message_sender_mocks.go -package achievementmessagesender

const _PROVIDER = "workers/userworker/achievementmessagesender.achievementMessageSender"

type AchievementMessageWrapper struct{ model.AchievementMessageData }

func (a AchievementMessageWrapper) VKHashData() string {
	var tp string
	switch a.AchievementType() {
	case achievementmodel.DURATION:
		tp = "duration"
	case achievementmodel.CIGARETTE:
		tp = "cigarretes"
	case achievementmodel.SAVING:
		tp = "money"
	case achievementmodel.HEALTH:
		tp = "health"
	case achievementmodel.WELL_BEING:
		tp = "wellbeing"
	}
	return "/achievements/" + tp
}

func (a AchievementMessageWrapper) Message() string {
	return fmt.Sprintf("Поздравляем вы достигли новый уровень в %s.\nОткройте его и получите дополнительную мотивацию!", a.AchievementType())
}

type AchievementMessageSender interface {
	SendMessage(ctx context.Context, userId int64, messageData model.AchievementMessageData) error
}

type achievementMessageSender struct {
	client     *http.Client
	apiVersion string
	token      string
	ownerId    int
	appId      int
}

func NewMessageParamsBuilder(userId int64, accessToken string, text string, version string) *message.MessageParamsBuilder {
	return message.NewMessageParamsBuilder().
		SetApiVersion(version).
		SetMessage(text).
		SetAccessToken(accessToken).
		SetUser(userId)
}

func NewMessageSender(client *http.Client, apiVersion string, token string, ownerId int, appId int) AchievementMessageSender {
	return &achievementMessageSender{
		client:     client,
		apiVersion: apiVersion,
		token:      token,
		ownerId:    ownerId,
		appId:      appId,
	}
}

func NewAppKeyboard(
	appId int,
	groupId int,
	label string,
	hash string,
) message.Keyboard {
	return message.NewKeyboardBuilder().
		SetInline(true).
		AddButtons(
			message.NewButton(
				message.NewOpenAppAction(appId, -groupId, "", label, hash),
			),
		).
		Build()
}

func (a *achievementMessageSender) SendMessage(ctx context.Context, userId int64, messageData model.AchievementMessageData) error {
	achMessage := AchievementMessageWrapper{messageData}
	keyboard := NewAppKeyboard(a.appId, a.ownerId, "Бросить Курить", achMessage.VKHashData())
	params := NewMessageParamsBuilder(userId, a.token, achMessage.Message(), a.apiVersion).
		SetRandomIDByMessage(achMessage.Message()).
		SetKeyboard(keyboard).
		Build()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.vk.com/method/messages.send", strings.NewReader(params.Encode()))
	if err != nil {
		return exception.Wrap(err, exception.NewCause("create request with ctx", "SendMessage", _PROVIDER))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := a.client.Do(req)
	if err != nil {
		return exception.Wrap(err, exception.NewCause("failed do request", "SendMessage", _PROVIDER))
	}
	var vk_err vkerror.VKError
	json.NewDecoder(resp.Body).Decode(&vk_err)
	resp.Body.Close()
	if vk_err.Err.Code != 0 {
		return exception.Wrap(&vk_err, exception.NewCause("send message err", "SendMessage", _PROVIDER))
	}
	return nil
}
