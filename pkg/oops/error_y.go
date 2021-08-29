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

func ExplainFmt(target error, format string, args ...interface{}) error {
	return Explain(target, fmt.Sprintf(format, args...))
}

// String is a helper method that will take any error type and return the normal .Error() for non Error errors. For
// Error type errors, it will instead return the Error.Code() and Error.Explanation().
func String(target error) string {
	err, isErr, isNil := As(target)
	if isNil {
		return ""
	}

	if isErr {
		return err.Code() + " " + err.Explanation()
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
