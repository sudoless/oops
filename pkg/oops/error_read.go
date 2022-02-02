package oops

// Code returns the identifying error code.
func (e *Error) Code() string {
	return e.source.code
}

// Type returns the error type.
func (e *Error) Type() string {
	return e.source.t
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

// StatusCode will return the mapped http status code for the given Error.
func (e *Error) StatusCode() int {
	return e.source.statusCode
}

// Help returns the defined error help message.
func (e *Error) Help() string {
	return e.source.help
}
