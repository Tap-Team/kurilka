package httphelpers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"gotest.tools/v3/assert"
)

func TestVkUserId(t *testing.T) {
	cases := []struct {
		id       int64
		query    string
		errIsNil bool
	}{
		{
			id:       1,
			query:    fmt.Sprintf("/example/pofig?vk_user_id=%d", 1),
			errIsNil: true,
		},
		{
			id:    0,
			query: fmt.Sprintf("/example/pofig/boolveas?user_id=%d", 1),
		},
		{
			id:       0,
			query:    fmt.Sprintf("/example/pofig?vk_user_id=%d", 0),
			errIsNil: true,
		},
	}

	for _, cs := range cases {
		req := httptest.NewRequest(http.MethodGet, cs.query, nil)

		id, err := httphelpers.VKID(req)

		assert.Equal(t, err == nil, cs.errIsNil, "wrong error")
		assert.Equal(t, id, cs.id, "id not equal")

	}
}
