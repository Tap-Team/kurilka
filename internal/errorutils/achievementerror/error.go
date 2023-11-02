package achievementerror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "achievement"

var (
	ExceptionAchievementNotExists = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "not_exists")
	}
	ExceptionCantOpenAchievementForUser = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "cant_open_for_user")
	}
)
