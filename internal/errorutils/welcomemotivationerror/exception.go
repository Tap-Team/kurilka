package welcomemotivationerror

import "github.com/Tap-Team/kurilka/pkg/exception"

const EType = "welcome_motivation"

var (
	ExceptionMotivationNotExist = func() exception.Exception {
		return exception.New(400, EType, "not_exist")
	}
)
