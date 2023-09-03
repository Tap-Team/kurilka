package errormodel

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewResponse(message, code string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Code:    code,
	}
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) Is(target error) bool {
	err, ok := target.(ErrorResponse)
	if !ok {
		return false
	}
	if err.Code != e.Code {
		return false
	}
	return true
}
