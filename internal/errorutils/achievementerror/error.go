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
)
