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

// FieldsList returns the raw fields list.
func (e *Error) FieldsList() []string {
	return e.fields
}

// FieldsMap returns the fields as a map.
func (e *Error) FieldsMap() map[string]string {
	if len(e.fields) == 0 {
		return nil
	}

	fields := make(map[string]string, len(e.fields)/2)
	for idx := 0; idx < len(e.fields); idx += 2 {
		fields[e.fields[idx]] = e.fields[idx+1]
	}

	return fields
}
