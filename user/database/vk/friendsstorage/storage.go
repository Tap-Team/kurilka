package friendsstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "user/database/vk/friendsstorage.Storage"

type Storage struct {
	client     *http.Client
	serviceKey string
	apiVersion string
}

func New(client *http.Client, serviceKey, apiVersion string) *Storage {
	return &Storage{
		client:     client,
		serviceKey: serviceKey,
		apiVersion: apiVersion,
	}
}

type Response struct {
	Response struct {
		Items []int64 `json:"items"`
	} `json:"response"`
}

func (s *Storage) Friends(ctx context.Context, userId int64) []int64 {
	friends := make([]int64, 0)
	urlValues := make(url.Values, 4)
	urlValues.Set("v", s.apiVersion)
	urlValues.Set("access_token", s.serviceKey)
	urlValues.Set("user_id", fmt.Sprint(userId))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.vk.com/method/friends.get", strings.NewReader(urlValues.Encode()))
	if err != nil {
		return friends
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		err := exception.Wrap(err, exception.NewCause("request error", "UserSubscriptionById", _PROVIDER))
		slog.ErrorContext(ctx, err.Error(), "user_id", userId)
		return friends
	}
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		err := exception.Wrap(err, exception.NewCause("decode body", "UserSubscriptionById", _PROVIDER))
		slog.ErrorContext(ctx, err.Error(), "user_id", userId)
		return friends
	}
	resp.Body.Close()
	return response.Response.Items
}
