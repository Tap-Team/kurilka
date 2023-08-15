package usererror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "user"

var (
	ExceptionUserNotFound = func() exception.Exception {
		return exception.New(http.StatusNotFound, EType, "not_found")
	}
)
