package message

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/url"
)

type MessageParamsBuilder struct {
	v url.Values
}

func NewMessageParamsBuilder() *MessageParamsBuilder {
	return &MessageParamsBuilder{v: make(url.Values)}
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

func (mp *MessageParamsBuilder) SetKeyboard(keyboard Keyboard) *MessageParamsBuilder {
	js, err := json.Marshal(keyboard)
	if err != nil {
		return mp
	}
	mp.v.Set("keyboard", string(js))
	return mp
}
