package oops

var (
	// ErrTODO is meant to be used as a placeholder error while developing software. It is recommended to add a lint
	// rule to catch any such errors in production or before committing.
	ErrTODO = Define().Code("todo").Type("test").StatusCode(482)

	// ErrUnexpected is the "default" error/behaviour when wrapping and/or explaining a non oops.Error, error. These
	// errors should be caught and investigated as they highlight bits of code where error handling is not exhaustive.
	ErrUnexpected = Define().Code("unexpected").Type("unexpected").StatusCode(500)

	// ErrMultiple reports that the current *Error is wrapping a multipleErrors, which contains its own explanation and
	// a slice of *Error.
	ErrMultiple = Define().Code("unexpected").Type("multiple").StatusCode(500)
)

// Explain is a helper method to wrap around Error or builtin error. Providing a builtin error will automatically
// generate an *Error using ErrUnexpected as the base, and calling Wrap in order to keep the target builtin error
// inheritance. If the given error is of type Error, then the explanation gets added to it.
func Explain(target error, explanation string, args ...any) *Error {
	if target == nil {
		return nil
	}

	oopsErr, ok := target.(*Error)
	if !ok {
		if multiErr, ok := target.(*multipleErrors); ok {
			return ErrMultiple.Wrap(multiErr, explanation, args...)
		}

		return ErrUnexpected.Wrap(target, explanation, args...)
	}

	if oopsErr == nil {
		return nil
	}

	return oopsErr.Explain(explanation, args...)
}

// Defer makes use of the Go error "handling" pattern that uses defer and a function that takes a named return error
// pointer and checks if it's nil or not, then performs a certain action, in this case you can define a standard
// format message and args, which will be added to the error explanation.
func Defer(err *error, format string, args ...any) {
	if err == nil {
		return
	}

	vErr := *err
	if vErr == nil {
		return
	}

	v, ok := vErr.(*Error)
	if !ok {
		*err = ErrUnexpected.Wrap(vErr, "deferred error").Explain(format, args...)
	} else {
		v.explain(format, args...)
	}
}

// As will take any type of error, if the error is not nil and not *Error, then a new ErrUnexpected is generated. In
// all cases the As function will return if the error isError (*Error) and/or if the error should be nil or not.
func As(target error) (err *Error, isError bool, isNil bool) {
	if target == nil {
		return nil, false, true
	}

	err, isError = target.(*Error)
	if !isError {
		if multiErr, ok := target.(*multipleErrors); ok {
			return ErrMultiple.Wrap(multiErr, "%d errors", len(multiErr.errs)), true, false
		}

		return ErrUnexpected.Wrap(target, ""), false, false
	}

	if err == nil {
		return nil, true, true
	}

	return err, true, false
}
