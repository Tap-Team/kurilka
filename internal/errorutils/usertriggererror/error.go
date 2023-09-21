package usertriggererror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "usertrigger"

var (
	UserTriggerNotFound = func() exception.Exception {
		return exception.New(http.StatusNotFound, EType, "not_found")
	}
	UserTriggerExists = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "exist")
	}
)
