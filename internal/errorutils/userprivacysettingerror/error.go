package userprivacysettingerror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "_user_privacy_setting"

var (
	ExceptionUserPrivacySettingNotFound = func() exception.Exception {
		return exception.New(http.StatusNotFound, EType, "setting_not_found")
	}
	ExceptionUserPrivacySettingExists = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "setting_exists")
	}
)
