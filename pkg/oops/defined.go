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
	noTrace    bool
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
	if !e.noTrace {
		err.trace = stack(3)
	}
	return err
}

// Yeet generates a new *Error that inherits the Blame, Namespace, Reason and help message from the parent errorDefined.
func (e *errorDefined) Yeet() *Error {
	return e.error()
}

// YeetExplain similar to Yeet but provides the option to add an explanation which can then be read with
// Error.Explain().
func (e *errorDefined) YeetExplain(explanation string) *Error {
	err := e.error()
	err.explain(explanation)
	return err
}

// YeetExplainFmt similar to YeetExplain but provides the option to pass the explanation as a format string and then
// the args. It is not recommended to include private information or user input as the args.
func (e *errorDefined) YeetExplainFmt(explanation string, args ...interface{}) *Error {
	err := e.error()
	err.explain(fmt.Sprintf(explanation, args...))
	return err
}

// Wrap generates a new *Error that inherits the Blame, Namespace, Reason and help message from the parent errorDefined
// and also sets the Error.parent to the target error. This can later be unwrapped using standard Go patterns.
func (e *errorDefined) Wrap(target error) *Error {
	err := e.error()
	err.parent = target
	return err
}

// WrapExplain similar to Wrap but provides the option to add an explanation which can then be read with
// Error.Explain().
func (e *errorDefined) WrapExplain(target error, explanation string) *Error {
	err := e.error()
	err.parent = target
	err.explain(explanation)
	return err
}

// WrapExplainFmt similar to WrapExplain but provides the option to pass the explanation as a format string and then
// the args. It is not recommended to include private information or user input as the args.
func (e *errorDefined) WrapExplainFmt(target error, explanation string, args ...interface{}) *Error {
	err := e.error()
	err.parent = target
	err.explain(fmt.Sprintf(explanation, args...))
	return err
}

// Error this will PANIC! Do not use! It is only defined to implement the builtin error interface so that errorDefined
// can beb used in errors.Is, etc.
func (e *errorDefined) Error() string {
	panic("oops: do not use errorDefined as error, use errorDefined.Yeet() and errorDefined.Wrap()")
}
