package oops

import (
	"fmt"
	"strings"
)

type ErrorDefined struct {
	t          string
	code       string
	help       string
	statusCode int
	trace      bool
}

// Define creates a top level error definition (*ErrorDefined) which should then be used to generate *Error using
// methods such as ErrorDefined.Yeet and ErrorDefined.Wrap. Passing a defined error as a builtin error will result
// in a panic. *ErrorDefined has a series of chained builders to create the error template. If multiple errors use
// the same type, status code, etc, the last call in the chain can be ErrorDefined.Group() to return a template for
// defined errors.
func Define() *ErrorDefined {
	return &ErrorDefined{}
}

func (e *ErrorDefined) error() *Error {
	err := &Error{
		source:      e,
		explanation: strings.Builder{},
	}
	if e.trace {
		err.trace = stack(3)
	}
	return err
}

// Yeet generates a new *Error that inherits the defined values from the parent ErrorDefined. If an explanation is
// given it gets appended to the error internal explanation which can be read with Error.Explain(). Furthermore, args
// can be passed which will be formatted into the explanation.
func (e *ErrorDefined) Yeet(explanation string, args ...any) *Error {
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

// Wrap generates a new *Error that inherits the values from the parent ErrorDefined and also sets the Error.parent to
// the target error. This can later be unwrapped using standard Go patterns (errors.Unwrap). If an explanation is
// given it gets appended to the error internal explanation which can be read with Error.Explain(). Furthermore, args
// can be passed which will be formatted into the explanation.
func (e *ErrorDefined) Wrap(target error, explanation string, args ...any) *Error {
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

// Error will PANIC! Do not use! It is only defined to implement the builtin error interface so that ErrorDefined
// can beb used in errors.Is, etc.
func (e *ErrorDefined) Error() string {
	panic("oops: do not use ErrorDefined as error, use ErrorDefined.Yeet() and ErrorDefined.Wrap()")
}
