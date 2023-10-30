package handler_test

import (
	"io"
	"strings"
	"testing"

	"github.com/Tap-Team/kurilka/vote/handler"
	"gotest.tools/v3/assert"
)

func TestVerifyBodySig(t *testing.T) {

	cases := []struct {
		body      io.Reader
		appSecret string
	}{
		{
			body:      strings.NewReader("app_id=51755509&item=kurilka_month_subscription_2770&lang=ru_RU&notification_type=get_subscription_test&order_id=2138953&receiver_id=211250278&user_id=211250278&sig=dce5072b45e4a37f3396ca43c36954b8"),
			appSecret: "8JQXTQDKu3fGvdsjCOhT",
		},
	}

	for _, cs := range cases {
		params := new(handler.SignParameters)
		_, err := params.ReadFrom(cs.body)
		assert.NilError(t, err, "failed read body")
		assert.Equal(t, true, params.Verify(cs.appSecret))
	}
}
