package subscriptiontypeerror

import (
	"net/http"

	"github.com/Tap-Team/kurilka/pkg/exception"
)

const EType = "subscription_type"

var (
	ExceptionSubscriptionTypeNotExists = func() exception.Exception {
		return exception.New(http.StatusBadRequest, EType, "not_exists")
	}
)
