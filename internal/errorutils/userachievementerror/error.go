package userachievementerror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "user_achievement"

var (
	ExceptionAchievementNotFound = func() exception.Exception { return exception.New(http.StatusNotFound, EType, "not_found") }
)
