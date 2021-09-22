package oops

import (
	"fmt"
	"strings"
)

type errorDefined struct {
	help      string
	code      string
	blame     Blame
	namespace Namespace
	reason    Reason
	defined   *errorDefined
	noStack   bool
}

// Define creates a top level error definition (*errorDefined) which should then be used to generate *Error using
// methods such as errorDefined.Yeet and errorDefined.Wrap. Passing a defined error as a builtin error will result
// in a panic.
func Define(blame Blame, namespace Namespace, reason Reason, help ...string) *errorDefined {
	var helpMsg string
	if len(help) > 0 {
		helpMsg = help[0]
	}

	err := &errorDefined{
		blame:     blame,
		namespace: namespace,
		reason:    reason,
		help:      helpMsg,
	}
	err.defined = err

	var builder strings.Builder

	builder.WriteString(blame.String())
	builder.WriteRune('.')
	builder.WriteString(namespace.String())
	builder.WriteRune('.')
	builder.WriteString(reason.String())

	err.code = builder.String()

	return err
}

func (e *errorDefined) error() *Error {
	err := &Error{
		defined:   e,
		blame:     e.blame,
		namespace: e.namespace,
		reason:    e.reason,
		help:      e.help,
	}
	if !e.noStack {
		err.trace = stack(3)
	}

	return err
}

// Code returns the three part code as defined by the Blame, Namespace and Reason.
func (e *errorDefined) Code() string {
	return e.code
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

// NoStack will disable the generation of the stack trace for the eventually returned *Error.
func (e *errorDefined) NoStack() *errorDefined {
	eNew := &errorDefined{
		help:      e.help,
		blame:     e.blame,
		namespace: e.namespace,
		reason:    e.reason,
		defined:   e,
		noStack:   true,
	}

	return eNew
}
