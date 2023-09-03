package triggererror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "trigger"

var (
	ExceptionTriggerNotExist = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "not_exist")
	}
)
