package httphelpers_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/errormodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {

	cases := []struct {
		err      error
		code     int
		response errormodel.ErrorResponse
	}{
		{
			err:      errors.New("failed something"),
			code:     500,
			response: errormodel.NewResponse("failed something", "internal"),
		},
		{
			err:      exception.New(404, "type", "code"),
			code:     404,
			response: errormodel.NewResponse("something error", "type_code"),
		},
	}
	for _, cs := range cases {
		w := httptest.NewRecorder()
		httphelpers.Error(w, cs.err)

		recorderContentType := w.Header().Get("Content-Type")
		require.Equal(t, "application/json", recorderContentType, "wrong recorder content type header")

		code := w.Result().StatusCode
		require.Equal(t, cs.code, code, "http status code not equal")

		var response errormodel.ErrorResponse
		err := json.NewDecoder(w.Result().Body).Decode(&response)
		require.NoError(t, err, "failed decode recorder body")

		require.True(t, errors.Is(response, cs.response), "response not equal")
	}
}
