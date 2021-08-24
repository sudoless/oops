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
