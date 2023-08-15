package validate

import (
	"fmt"
	"net/http"
	"strings"
)

type WrongLenError struct {
	RangeValidatable
	actual int
}

func (v *WrongLenError) Error() string {
	return fmt.Sprintf("wrong len of parameter %s min %d max %d actual %d", v.Name(), v.Min(), v.Max(), v.actual)
}
func (v *WrongLenError) HttpCode() int {
	return http.StatusBadRequest
}
func (v *WrongLenError) Code() string {
	return "wrong_len"
}
func (v *WrongLenError) Type() string {
	return "common"
}

// ${param} должна находится в диапазоне от ${min} до ${max}
func (v *WrongLenError) Replace(target string) string {
	replacer := strings.NewReplacer("${param}", v.Name(), "${min}", fmt.Sprint(v.Min()), "${max}", fmt.Sprint(v.Max()))
	return replacer.Replace(target)
}

type RangeValidatable interface {
	Name() string
	Min() int
	Max() int
}
