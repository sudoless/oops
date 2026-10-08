package oops

// Error returns the string representation using the definition's formatter
// or the default formatter: code, code: message, code: explanation or
// code: message; explanation, depending on which parts are set. Text of wrapped
// errors never appears in it.
func (err *Error) Error() string {
	if err == nil {
		return "oops.Error(nil)"
	}

	if err.def == nil {
		return "oops.Error(undefined)"
	}

	if err.def.formatter != nil {
		return err.def.formatter(err)
	}

	return defaultFormatter(err)
}

// Unwrap implements the multi-error unwrap interface (Unwrap() []error). The
// returned slice is the error's own storage: do not modify it.
func (err *Error) Unwrap() []error {
	if err == nil {
		return nil
	}

	return err.wrapped
}

// Is reports whether err's definition is, or inherits, the target's definition.
// The target is an *ErrorDefinition or an *Error (compared by its definition).
// It compares definitions only and never looks at wrapped errors: use
// errors.Is to traverse the tree. A nil or typed-nil target matches only a
// nil err.
func (err *Error) Is(other error) bool {
	if isNil(other) {
		return err == nil
	}

	if err == nil {
		return false
	}

	switch v := other.(type) {
	case *ErrorDefinition:
		return err.def.is(v)
	case *Error:
		return err.def.is(v.def)
	}

	return false
}
