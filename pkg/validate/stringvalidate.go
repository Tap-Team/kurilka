package validate

import (
	"fmt"
)

type ValidatableString interface {
	fmt.Stringer
	RangeValidatable
}

func StringValidate(s ValidatableString) error {
	l := len([]rune(s.String()))
	min, max := s.Min(), s.Max()
	if l < min || l > max {
		return &WrongLenError{
			RangeValidatable: s,
			actual:           l,
		}
	}
	return nil
}
