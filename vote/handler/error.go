package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

type VKPaymentError interface {
	error
	Code() int
	Critical() bool
}

type vkPaymentError struct {
	Err struct {
		Code     int    `json:"error_code"`
		Msg      string `json:"error_msg"`
		Critical bool   `json:"critical"`
	} `json:"error"`
}

func (v *vkPaymentError) SetCode(code int) {
	v.Err.Code = code
}

func (v *vkPaymentError) SetCritical(critical bool) {
	v.Err.Critical = critical
}

func (v vkPaymentError) Code() int {
	return v.Err.Code
}

func (v vkPaymentError) Critical() bool {
	return v.Err.Critical
}

func (v vkPaymentError) Error() string {
	return v.Err.Msg
}

func NewPaymentError(code int, msg string, critical bool) vkPaymentError {
	return vkPaymentError{
		Err: struct {
			Code     int    "json:\"error_code\""
			Msg      string "json:\"error_msg\""
			Critical bool   "json:\"critical\""
		}{
			Code:     code,
			Msg:      msg,
			Critical: critical,
		},
	}
}

func Error(w http.ResponseWriter, err error) {
	var httpCode int = http.StatusInternalServerError
	vkPaymentError := NewPaymentError(1, err.Error(), false)
	if err, ok := err.(interface{ Code() int }); ok {
		vkPaymentError.SetCode(err.Code())
	}
	if err, ok := err.(interface{ Critical() bool }); ok {
		vkPaymentError.SetCritical(err.Critical())
	}

	if err, ok := err.(interface{ HttpCode() int }); ok {
		httpCode = err.HttpCode()
	}
	w.WriteHeader(httpCode)
	slog.Error(err.Error())
	json.NewEncoder(w).Encode(vkPaymentError)
}

func AssertError(t *testing.T, err error, body io.ReadCloser) {
	var bodyErr vkPaymentError
	if decodeErr := json.NewDecoder(body).Decode(&bodyErr); decodeErr != nil {
		t.Fatalf("failed assert error: decode body, err: %s", decodeErr)
	}
	body.Close()
	vkPaymentError := NewPaymentError(1, err.Error(), false)
	if err, ok := err.(interface{ Code() int }); ok {
		vkPaymentError.SetCode(err.Code())
	}
	if err, ok := err.(interface{ Critical() bool }); ok {
		vkPaymentError.SetCritical(err.Critical())
	}
	if bodyErr.Code() != vkPaymentError.Code() {
		t.Fatal("error code not equal")
	}
	if bodyErr.Critical() != vkPaymentError.Critical() {
		t.Fatal("error critical not equal")
	}
	if bodyErr.Error() != vkPaymentError.Error() {
		t.Fatal("error not equal")
	}
}
