package oops

import "fmt"

type errorDefined struct {
	help      string
	blame     Blame
	namespace Namespace
	reason    Reason
}

// Define creates a top level error definition (*errorDefined) which should then be used to generate *Error using
// methods such as errorDefined.Yeet and errorDefined.Wrap. Passing a defined error as a builtin error will result
// in a panic.
func Define(blame Blame, namespace Namespace, reason Reason, help ...string) *errorDefined {
	var helpMsg string
	if len(help) > 0 {
		helpMsg = help[0]
	}

	return &errorDefined{
		blame:     blame,
		namespace: namespace,
		reason:    reason,
		help:      helpMsg,
	}
}

// Yeet generates a new *Error that inherits the Blame, Namespace, Reason and help message from the parent errorDefined.
func (e *errorDefined) Yeet() *Error {
	return &Error{
		defined:   e,
		blame:     e.blame,
		namespace: e.namespace,
		reason:    e.reason,
		help:      e.help,
	}
}

// YeetExplain similar to Yeet but provides the option to add an explanation which can then be read with
// Error.Explanation().
func (e *errorDefined) YeetExplain(explanation string) *Error {
	err := e.Yeet()
	err.explain(explanation)
	return err
}

// YeetExplainFmt similar to YeetExplain but provides the option to pass the explanation as a format string and then
// the args. It is not recommended to include private information or user input as the args.
func (e *errorDefined) YeetExplainFmt(explanation string, args ...interface{}) *Error {
	err := e.Yeet()

	err.explain(fmt.Sprintf(explanation, args...))
	return err
}

// Wrap generates a new *Error that inherits the Blame, Namespace, Reason and help message from the parent errorDefined
// and also sets the Error.parent to the target error. This can later be unwrapped using standard Go patterns.
func (e *errorDefined) Wrap(target error) *Error {
	return &Error{
		defined:   e,
		blame:     e.blame,
		namespace: e.namespace,
		reason:    e.reason,
		help:      e.help,
		parent:    target,
	}
}

// WrapExplain similar to Wrap but provides the option to add an explanation which can then be read with
// Error.Explanation().
func (e *errorDefined) WrapExplain(target error, explanation string) *Error {
	err := e.Wrap(target)
	err.explain(explanation)
	return err
}

// WrapExplainFmt similar to WrapExplain but provides the option to pass the explanation as a format string and then
// the args. It is not recommended to include private information or user input as the args.
func (e *errorDefined) WrapExplainFmt(target error, explanation string, args ...interface{}) *Error {
	err := e.Wrap(target)
	err.explain(fmt.Sprintf(explanation, args...))
	return err
}

// Error this will PANIC! Do not use! It is only defined to implement the builtin error interface so that errorDefined
// can beb used in errors.Is, etc.
func (e *errorDefined) Error() string {
	panic("do not use errorDefined as error, use errorDefined.Yeet() and errorDefined.Wrap()")
}
