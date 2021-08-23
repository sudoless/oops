package oops

import "strings"

func (e *Error) Code() string {
	var builder strings.Builder

	builder.WriteString(e.blame.String())
	builder.WriteRune('.')
	builder.WriteString(e.namespace.String())
	builder.WriteRune('.')
	builder.WriteString(e.reason.String())

	return builder.String()
}

func (e *Error) Explanation() string {
	return e.explanation.String()
}

func (e *Error) Trace() []string {
	return e.trace
}

func (e *Error) Multiples() []string {
	return e.multi
}
