package oops

var ErrUnexpected = Define(BlameUnknown, NamespaceUnknown, ReasonUnexpected)

// Explain is a helper method to wrap around Error or builtin error. Providing a builtin error will automatically
// generate an *Error using ErrUnexpected as the base, and calling Wrap in order to keep the target builtin error
// inheritance. If the given error is of type Error, then the explanation gets added to it.
func Explain(target error, explanation string) *Error {
	if target == nil {
		return nil
	}

	err, ok := target.(*Error)
	if !ok {
		return ErrUnexpected.Wrap(target)
	}

	err.explain(explanation)

	return err
}

// String is a helper method that will take any error type and return the normal .Error() for non Error errors. For
// Error type errors, it will instead return the Error.Code() and Error.Explanation().
func String(target error) string {
	if target == nil {
		return ""
	}

	err, ok := target.(*Error)
	if !ok {
		return err.Error()
	}

	return err.Code() + " " + err.Explanation()
}
