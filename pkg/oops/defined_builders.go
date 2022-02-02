package oops

// Code assigns what _should_ be a unique identifying code to the error.
func (e *errorDefined) Code(code string) *errorDefined {
	e.code = code
	return e
}

// Type assigns a general identifying type group to the error (eg: validation, authentication, etc).
func (e *errorDefined) Type(t string) *errorDefined {
	e.t = t
	return e
}

// Help assigns a helpful message to the error.
func (e *errorDefined) Help(help string) *errorDefined {
	e.help = help
	return e
}

// StatusCode assigns a http status code to the error to be used in the response.
func (e *errorDefined) StatusCode(statusCode int) *errorDefined {
	e.statusCode = statusCode
	return e
}

// NoTrace will disable the generation of the stack trace for the eventually returned *Error.
func (e *errorDefined) NoTrace() *errorDefined {
	e.noTrace = true
	return e
}
