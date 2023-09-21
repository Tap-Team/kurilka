package vkerror

import "fmt"

type VKError struct {
	Err struct {
		Code    int    `json:"error_code"`
		Message string `json:"error_msg"`
	} `json:"error"`
}

func (e *VKError) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Err.Code, e.Err.Message)
}
