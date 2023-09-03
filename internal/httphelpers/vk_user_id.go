package httphelpers

import (
	"errors"
	"net/http"
	"strconv"
)

const vk_user_id_query = "vk_user_id"

var (
	ErrFailedParseVK_ID = errors.New("failed parse vk_user_id query")
)

func VKID(r *http.Request) (int64, error) {
	queryId := r.URL.Query().Get(vk_user_id_query)
	userId, err := strconv.ParseInt(queryId, 10, 64)
	if err != nil {
		return 0, ErrFailedParseVK_ID
	}
	return userId, nil
}
