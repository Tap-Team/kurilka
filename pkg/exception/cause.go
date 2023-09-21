package exception

import "fmt"

type Cause interface {
	error
	Action() string
	Method() string
	Provider() string
}

type amidCause struct {
	action   string
	method   string
	provider string
}

func NewCause(action string, method string, provider string) Cause {
	return &amidCause{action: action, method: method, provider: provider}
}

func (a *amidCause) Error() string {
	return fmt.Sprintf("Action %s of Method %s in %s provider", a.action, a.method, a.provider)
}

func (a *amidCause) Action() string {
	return a.action
}

func (e *amidCause) Method() string {
	return e.method
}

func (e *amidCause) Provider() string {
	return e.provider
}

func (e *amidCause) Is(target error) bool {
	err, ok := target.(Cause)
	if !ok {
		return false
	}
	if e.action != err.Action() {
		return false
	}
	if e.method != err.Method() {
		return false
	}
	if e.provider != err.Provider() {
		return false
	}
	return true
}
