package repeater

import "time"

type Repeater interface {
	Repeat(f func() error)
	Errors() []error
}

func New(n int, pause time.Duration) Repeater {
	return &repeater{times: n, time: pause}
}

type repeater struct {
	times  int
	time   time.Duration
	errors []error
}

func (r *repeater) Repeat(rf func() error) {
	for i := 0; i < r.times; i++ {
		err := rf()
		if err == nil {
			break
		}
		r.errors = append(r.errors, err)
	}
}

func (r repeater) Errors() []error {
	return r.errors
}
