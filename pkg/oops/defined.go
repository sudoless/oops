package oops

import (
	"fmt"
	"strings"
)

type errorDefined struct {
	t          string
	code       string
	help       string
	statusCode int
	trace      bool
}

// Define creates a top level error definition (*errorDefined) which should then be used to generate *Error using
// methods such as errorDefined.Yeet and errorDefined.Wrap. Passing a defined error as a builtin error will result
// in a panic. *errorDefined has a series of chained builders to create the error template. If multiple errors use
// the same type, status code, etc, the last call in the chain can be errorDefined.Group() to return a template for
// defined errors.
func Define() *errorDefined {
	return &errorDefined{}
}

func (e *errorDefined) error() *Error {
	err := &Error{
		source:      e,
		explanation: strings.Builder{},
	}
	if e.trace {
		err.trace = stack(3)
	}
	return err
}

// Yeet generates a new *Error that inherits the defined values from the parent errorDefined. If an explanation is
// given it gets appended to the error internal explanation which can be read with Error.Explain(). Furthermore, args
// can be passed which will be formatted into the explanation.
func (e *errorDefined) Yeet(explanation string, args ...any) *Error {
	err := e.error()

	if len(explanation) > 0 {
		if len(args) > 0 {
			err.explain(fmt.Sprintf(explanation, args...))
		} else {
			err.explain(explanation)
		}
	}

	return err
}

// Wrap generates a new *Error that inherits the values from the parent errorDefined and also sets the Error.parent to
// the target error. This can later be unwrapped using standard Go patterns (errors.Unwrap). If an explanation is
// given it gets appended to the error internal explanation which can be read with Error.Explain(). Furthermore, args
// can be passed which will be formatted into the explanation.
func (e *errorDefined) Wrap(target error, explanation string, args ...any) *Error {
	err := e.error()
	err.parent = target

	if len(explanation) > 0 {
		if len(args) > 0 {
			err.explain(fmt.Sprintf(explanation, args...))
		} else {
			err.explain(explanation)
		}
	}

	return err
}

// Error will PANIC! Do not use! It is only defined to implement the builtin error interface so that errorDefined
// can beb used in errors.Is, etc.
func (e *errorDefined) Error() string {
	panic("oops: do not use errorDefined as error, use errorDefined.Yeet() and errorDefined.Wrap()")
}
