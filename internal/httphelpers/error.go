package httphelpers

import (
	"net/http"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/model/errormodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
)

func Error(w http.ResponseWriter, err error) {
	httpCode := http.StatusInternalServerError
	messageCode := "internal"
	if httpErr, ok := err.(exception.HttpError); ok {
		httpCode = httpErr.HttpCode()
	}
	if codeTypeErr, ok := err.(exception.CodeTypedError); ok {
		messageCode = exception.MakeCode(codeTypeErr)
	}

	slog.Error(err.Error())
	WriteJSON(w, errormodel.NewResponse(err.Error(), messageCode), httpCode)
}
