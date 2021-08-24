package oops

import "fmt"

type errorDefined struct {
	help      string
	blame     Blame
	namespace Namespace
	reason    Reason
}

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

func (e *errorDefined) Yeet() *Error {
	return &Error{
		defined:   e,
		blame:     e.blame,
		namespace: e.namespace,
		reason:    e.reason,
		help:      e.help,
	}
}

func (e *errorDefined) YeetExplain(explanation string) *Error {
	err := e.Yeet()
	err.explain(explanation)
	return err
}

func (e *errorDefined) YeetExplainFmt(explanation string, args ...interface{}) *Error {
	err := e.Yeet()

	err.explain(fmt.Sprintf(explanation, args...))
	return err
}

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

func (e *errorDefined) WrapExplain(target error, explanation string) *Error {
	err := e.Wrap(target)
	err.explain(explanation)
	return err
}

func (e *errorDefined) WrapExplainFmt(target error, explanation string, args ...interface{}) *Error {
	err := e.Wrap(target)
	err.explain(fmt.Sprintf(explanation, args...))
	return err
}

func (e *errorDefined) Error() string {
	panic("do not use errorDefined as error, use errorDefined.Yeet() and errorDefined.Wrap()")
}
