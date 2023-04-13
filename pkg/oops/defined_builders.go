package oops

// Code assigns what _should_ be a unique identifying code to the error.
func (e *ErrorDefined) Code(code string) *ErrorDefined {
	e.code = code
	return e
}

// Type assigns a general identifying type group to the error (eg: validation, authentication, etc).
func (e *ErrorDefined) Type(t string) *ErrorDefined {
	e.t = t
	return e
}

// StatusCode assigns a http status code to the error to be used in the response.
func (e *ErrorDefined) StatusCode(statusCode int) *ErrorDefined {
	e.statusCode = statusCode
	return e
}

// Trace will enable the generation of the stack trace for the eventually returned *Error.
func (e *ErrorDefined) Trace() *ErrorDefined {
	e.trace = true
	return e
}
