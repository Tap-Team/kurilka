package vk

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Tap-Team/kurilka/internal/errorutils/vkerror"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/vk/message"
)

const _PROVIDER = "workers/userworker/messagesender/vk.sender"

type sender struct {
	client     *http.Client
	apiVersion string
	token      string
}

func NewMessageSender(
	client *http.Client,
	apiVersion, token string,
) messagesender.MessageSender {
	return &sender{
		client:     client,
		apiVersion: apiVersion,
		token:      token,
	}
}

func NewMessageParamsBuilder(userId int64, accessToken string, text string, version string) *message.MessageParamsBuilder {
	params := message.NewMessageParamsBuilder()
	return params.
		SetApiVersion(version).
		SetMessage(text).
		SetAccessToken(accessToken).
		SetUser(userId)
}

func (s *sender) SendMessage(ctx context.Context, message string, userId int64) error {
	params := NewMessageParamsBuilder(userId, s.token, message, s.apiVersion).SetRandomIDByMessage(message)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.vk.com/method/messages.send", strings.NewReader(params.Build().Encode()))
	if err != nil {
		return exception.Wrap(err, exception.NewCause("create request with ctx", "SendMessage", _PROVIDER))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
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
