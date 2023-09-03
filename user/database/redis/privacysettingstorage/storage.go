package privacysettingstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/redishelper"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/redis/go-redis/v9"
)

const _PROVIDER = "user/database/redis/privacysettingstorage"

type Storage struct {
	redis      *redis.Client
	expiration time.Duration
}

func New(redis *redis.Client, exp time.Duration) *Storage {
	return &Storage{redis: redis}
}

func Error(err error, cause exception.Cause) error {
	return exception.Wrap(err, cause)
}

type privacySettingList []usermodel.PrivacySetting

func (p privacySettingList) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	for _, stn := range p {
		buf.Write([]byte(`"` + stn + `"` + ","))
	}
	return buf.Bytes(), nil
}

func (p *privacySettingList) UnmarshalBinary(data []byte) error {
	s := strings.TrimSuffix(string(data), ",")
	list := strings.Split(s, ",")
	settingList := make([]usermodel.PrivacySetting, 0, len(list))
	for _, e := range list {
		stn := usermodel.PrivacySetting(strings.Trim(e, `"`))
		settingList = append(settingList, stn)
	}
	*p = settingList
	return nil
}

func (s *Storage) SaveUserPrivacySettings(ctx context.Context, userId int64, settings []usermodel.PrivacySetting) error {
	err := s.redis.Set(ctx, redishelper.PrivacySettingsKey(userId), privacySettingList(settings), s.expiration).Err()
	if err != nil {
		return Error(err, exception.NewCause("set user privacy settings", "SaveUserPrivacySettings", _PROVIDER))
	}
	return nil
}

func (s *Storage) RemoveUserPrivacySettings(ctx context.Context, userId int64) error {
	err := s.redis.Del(ctx, redishelper.PrivacySettingsKey(userId), fmt.Sprint(userId)).Err()
	if err != nil {
		return Error(err, exception.NewCause("delete user privacy settings", "RemoveUserPrivacySettings", _PROVIDER))
	}
	return nil
}

func (s *Storage) UserPrivacySettings(ctx context.Context, userId int64) ([]usermodel.PrivacySetting, error) {
	settingsList := make(privacySettingList, 0)
	err := s.redis.Get(ctx, redishelper.PrivacySettingsKey(userId)).Scan(&settingsList)
	if errors.Is(err, redis.Nil) {
		return settingsList, nil
	}
	if err != nil {
		return settingsList, Error(err, exception.NewCause("get user privacy settings", "UserPrivacySettings", _PROVIDER))
	}
	return settingsList, nil
}
