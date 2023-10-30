package handler

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

type Params map[string]any

func (p Params) NotificationType() string {
	tp, ok := p["notification_type"].(string)
	if !ok {
		return ""
	}
	return tp
}

func (p Params) ReadFrom(r io.Reader) (n int64, err error) {
	buf := new(bytes.Buffer)
	n, err = buf.ReadFrom(r)
	s := buf.String()
	params := strings.Split(s, "&")
	for _, keyVal := range params {
		var pair KeyValuePair
		ok := pair.Parse(keyVal)
		if !ok {
			continue
		}
		p[pair.Key] = pair.Value
	}
	return
}

func (p Params) GetString(key string) string {
	value, ok := p[key].(string)
	if !ok {
		return ""
	}
	return value
}

func (p Params) GetInt(key string) int64 {
	value, ok := p[key].(string)
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
