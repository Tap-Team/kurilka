package httphelpers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Tap-Team/kurilka/internal/model/errormodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"gotest.tools/v3/assert"
)

func AssertError(t *testing.T, err error, rec *httptest.ResponseRecorder) {
	var errResponse errormodel.ErrorResponse
	decodeErr := json.NewDecoder(rec.Result().Body).Decode(&errResponse)
	rec.Result().Body.Close()
	assert.NilError(t, decodeErr, "failed decode body")

	if err, ok := err.(exception.CodeTypedError); ok {
		assert.Equal(t, errResponse.Code, exception.MakeCode(err), "error code not equal")
	}
}
