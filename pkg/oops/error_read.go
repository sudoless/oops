package oops

// Code returns the identifying error code.
func (e *Error) Code() string {
	return e.source.code
}

// Type returns the error type.
func (e *Error) Type() string {
	return e.source.t
}

// Explanation returns the accumulated explanations, each original explanation string separated by a comma.
func (e *Error) Explanation() string {
	return e.explanation.String()
}

// Trace returns the stack trace array generated at the time of creating by calling the Error.Stack() method.
func (e *Error) Trace() []string {
	return e.trace
}

// StatusCode will return the mapped http status code for the given Error.
func (e *Error) StatusCode() int {
	return e.source.statusCode
}

// Help returns the defined error help message.
func (e *Error) Help() string {
	return e.source.help
}

// Err returns the Error as an error. Can be used to properly return a nil error when doing things like:
// return oops.Explain(possiblyNilError, "some explanation").Err()
func (e *Error) Err() error {
	if e == nil {
		return nil
	}

	return e
}

// Errs returns the Error as an error slice. If the Error is ErrMultiple then the slice will contain all the errors
// from the parent multi error. Otherwise, it will return a slice with a single error (itself).
func (e *Error) Errs() []error {
	if e == nil {
		return nil
	}

	if e.source == ErrMultiple && e.parent != nil {
		return e.parent.(*multipleErrors).Unwrap()
	}

	return []error{e}
}

func (e *Error) ErrsAs() []*Error {
	if e == nil {
		return nil
	}

	if e.source == ErrMultiple && e.parent != nil {
		return e.parent.(*multipleErrors).errs
	}

	return []*Error{e}
}
