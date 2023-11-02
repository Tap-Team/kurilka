package privacysettingerror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "privacy_setting"

var (
	ExceptionPrivacySettingNotExist = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "not_exist")
	}
	ExceptionUserWithoutSubscription = func() exception.Exception {
		return exception.New(http.StatusForbidden, EType, "user_without_subscription")
	}
)
