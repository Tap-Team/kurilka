package motivationerror

import "github.com/Tap-Team/kurilka/pkg/exception"

const EType = "motivation"

var (
	ExceptionMotivationNotExist = func() exception.Exception {
		return exception.New(400, EType, "not_exist")
	}
)
