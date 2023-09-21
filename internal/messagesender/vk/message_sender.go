package vk

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"

	"github.com/Tap-Team/kurilka/internal/errorutils/vkerror"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/pkg/exception"
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

type MessageParamsBuilder struct {
	v url.Values
}

func (m *MessageParamsBuilder) Build() url.Values {
	return m.v
}

func (mp *MessageParamsBuilder) SetApiVersion(version string) *MessageParamsBuilder {
	mp.v.Set("v", version)
	return mp
}

func (mp *MessageParamsBuilder) SetUser(userId int64) *MessageParamsBuilder {
	mp.v.Set("user_id", fmt.Sprint(userId))
	return mp
}
func (mp *MessageParamsBuilder) SetAccessToken(accessToken string) *MessageParamsBuilder {
	mp.v.Set("access_token", accessToken)
	return mp
}
func (mp *MessageParamsBuilder) SetMessage(message string) *MessageParamsBuilder {
	mp.v.Set("message", message)
	return mp
}

func (mp *MessageParamsBuilder) SetRandomID(randomId int64) *MessageParamsBuilder {
	mp.v.Set("random_id", fmt.Sprint(randomId))
	return mp
}

func (mp *MessageParamsBuilder) SetRandomIDByMessage(message string) *MessageParamsBuilder {
	f := fnv.New32a()
	f.Write([]byte(message))
	randomId := f.Sum32()
	mp.SetRandomID(int64(randomId))
	return mp
}

func NewMessageParamsBuilder(userId int64, accessToken string, message string, version string) *MessageParamsBuilder {
	params := &MessageParamsBuilder{make(url.Values)}
	return params.
		SetApiVersion(version).
		SetMessage(message).
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
