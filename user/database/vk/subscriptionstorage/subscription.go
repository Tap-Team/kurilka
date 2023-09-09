package subscriptionstorage

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const _PROVIDER = "user/database/vk/subscriptionstorage"

type SubscriptionStorage struct {
	client     *http.Client
	apiVersion string
	groupId    int64
	groupToken string
}

func New(client *http.Client, apiVersion string, groupId int64, groupToken string) *SubscriptionStorage {
	return &SubscriptionStorage{
		client:     client,
		apiVersion: apiVersion,
		groupId:    groupId,
		groupToken: groupToken,
	}
}

func Error(err error, cause exception.Cause) error {
	return exception.Wrap(err, cause)
}

type userSubscriptionResponse struct {
	Response struct {
		Count int     `json:"count"`
		Items []int64 `json:"items"`
	} `json:"response"`
}

func (u userSubscriptionResponse) Find(userId int64) bool {
	_, found := sort.Find(len(u.Response.Items), func(i int) int {
		return cmp.Compare(userId, u.Response.Items[i])
	})
	return found
}

func (s *SubscriptionStorage) UserSubscriptionById(ctx context.Context, userId int64) (time.Time, error) {
	urlValues := make(url.Values, 4)
	urlValues.Set("v", s.apiVersion)
	urlValues.Set("access_token", s.groupToken)
	urlValues.Set("group_id", fmt.Sprint(s.groupId))
	urlValues.Set("sort", "id_asc")
	urlValues.Set("filter", "donut")
	var response userSubscriptionResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.vk.com/method/groups.getMembers", strings.NewReader(urlValues.Encode()))
	if err != nil {
		return time.Time{}, Error(err, exception.NewCause("create request error", "UserSubscriptionById", _PROVIDER))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return time.Time{}, Error(err, exception.NewCause("request error", "UserSubscriptionById", _PROVIDER))
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return time.Time{}, Error(err, exception.NewCause("decode body", "UserSubscriptionById", _PROVIDER))
	}
	resp.Body.Close()
	if response.Find(userId) {
		return time.Now().Add(time.Hour * 24), nil
	} else {
		return time.Time{}, Error(errors.New("user not found"), exception.NewCause("search user", "userSubscriptionById", _PROVIDER))
	}
}
