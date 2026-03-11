package oops

import "errors"

// Error returns the string representation using the definition's formatter
// or the default formatter.
func (err *Error) Error() string {
	if err == nil {
		return "oops.Error(nil)"
	}

	if err.def.formatter != nil {
		return err.def.formatter(err)
	}

	return defaultFormatter(err)
}

// Unwrap implements the multi-error unwrap interface (Unwrap() []error).
func (err *Error) Unwrap() []error {
	return err.wrapped
}

// Is checks definition identity including the inherits chain, and for non-oops
// targets it checks whether any wrapped error matches via errors.Is.
func (err *Error) Is(other error) bool {
	if err == nil {
		return other == nil
	}

	switch v := other.(type) {
	case *ErrorDefinition:
		return err.def.is(v)
	case *Error:
		return err.def == v.def
	}

	wrapped := err.wrapped
	for idx := len(wrapped) - 1; idx >= 0; idx-- {
		w := wrapped[idx]
		if errors.Is(w, other) {
			return true
		}
	}

	return false
}
