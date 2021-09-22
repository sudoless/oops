package oops

import "strings"

// Code builds and returns the three part code (BLAME.NAMESPACE.REASON) as a string.
func (e *Error) Code() string {
	if e.code != "" {
		return e.code
	}

	var builder strings.Builder

	builder.WriteString(e.blame.String())
	builder.WriteRune('.')
	builder.WriteString(e.namespace.String())
	builder.WriteRune('.')
	builder.WriteString(e.reason.String())

	e.code = builder.String()
	return e.code
}

// Explain returns the accumulated explanations, each original explanation string separated by a comma.
func (e *Error) Explain() string {
	return e.explanation.String()
}

// Trace returns the stack trace array generated at the time of creating by calling the Error.Stack() method.
func (e *Error) Trace() []string {
	return e.trace
}

// Multiples returns the list of strings (considered to be a multi-reason/explanation-error) defined by calling
// the Error.Multi() method.
func (e *Error) Multiples() []string {
	return e.multi
}

// StatusCode will return the mapped http status code for the given Error Reason, by using the lookup table
// mapReasonToHttpStatus used by Reason.HttpStatusCode.
func (e *Error) StatusCode() int {
	return e.reason.HttpStatusCode()
}
