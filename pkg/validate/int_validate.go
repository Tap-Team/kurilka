package validate

type ValidatableInt interface {
	RangeValidatable
	Int() int64
}

func IntValidate(vi ValidatableInt) error {
	min, max := vi.Min(), vi.Max()
	i := int(vi.Int())
	if i < min || i > max {
		return &WrongLenError{
			RangeValidatable: vi,
			actual:           i,
		}
	}
	return nil
}
