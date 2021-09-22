package oops

import "fmt"

var (
	// ErrTODO is meant to be used as a placeholder error while developing software. It is recommended to add a lint
	// rule to catch any such errors in production or before committing.
	ErrTODO = Define(
		BlameDeveloper, NamespaceUnknown, ReasonUnexpected,
		"error 482, somebody just shot the server with a 12-gauge, please contact your administrator")

	// ErrUnexpected is the "default" error/behaviour when wrapping and/or explaining a non oops.Error, error. These
	// errors should be caught and investigated as they highlight bits of code where error handling is not exhaustive.
	ErrUnexpected = Define(
		BlameUnknown, NamespaceUnknown, ReasonUnexpected)
)

// Explain is a helper method to wrap around Error or builtin error. Providing a builtin error will automatically
// generate an *Error using ErrUnexpected as the base, and calling Wrap in order to keep the target builtin error
// inheritance. If the given error is of type Error, then the explanation gets added to it.
func Explain(target error, explanation string) error {
	err, isErr, isNil := As(target)
	if isNil {
		return nil
	}

	if isErr {
		err.explain(explanation)
		return err
	}

	return ErrUnexpected.WrapExplain(target, explanation)
}

// ExplainFmt similar to Explain, but takes a format string and optional args.
func ExplainFmt(target error, format string, args ...interface{}) error {
	return Explain(target, fmt.Sprintf(format, args...))
}

// String is a helper method that will take any error type and return the normal .Error() for non Error errors. For
// Error type errors, it will instead return the Error.Code() and Error.Explain().
func String(target error) string {
	err, isErr, isNil := As(target)
	if isNil {
		return ""
	}

	if isErr {
		return err.Code() + " " + err.Explain()
	}

	return target.Error()
}

// As will take any type of error, if the error is not nil and not *Error, then a new ErrUnexpected is generated. In
// all cases the As function will return if the error isError (*Error) and/or if the error should be nil or not.
func As(target error) (err *Error, isError bool, isNil bool) {
	if target == nil {
		return nil, false, true
	}

	err, ok := target.(*Error)
	if !ok {
		return ErrUnexpected.Wrap(target), false, false
	}

	if err == nil {
		return nil, true, true
	}

	return err, true, false
}

// Defer makes use of the Go error "handling" pattern that uses defer and a function that takes a named return error
// pointer and checks if it's nil or not, then performs a certain action, in this case you can define a standard
// format message and args, which will be added to the error explanation.
func Defer(err *error, format string, args ...interface{}) {
	if err == nil {
		return
	}

	vErr := *err
	oopsErr, isErr, isNil := As(vErr)
	if isNil {
		return
	}

	explanation := fmt.Sprintf(format, args...)

	if isErr {
		oopsErr.explain(explanation)
		return
	}

	*err = ErrUnexpected.WrapExplain(vErr, explanation)
}
