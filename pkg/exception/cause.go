package exception

import "fmt"

type Cause interface {
	error
	Action() string
	Method() string
	Pkg() string
}

type amidCause struct {
	action string
	method string
	pkg    string
}

func NewCause(action string, method string, pkg string) Cause {
	return &amidCause{action: action, method: method, pkg: pkg}
}

func (a *amidCause) Error() string {
	return fmt.Sprintf("Action %s of Method %s in %s package", a.action, a.method, a.pkg)
}

func (a *amidCause) Action() string {
	return a.action
}

func (e *amidCause) Method() string {
	return e.method
}

func (e *amidCause) Pkg() string {
	return e.pkg
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
	if e.pkg != err.Pkg() {
		return false
	}
	return true
}
