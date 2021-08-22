package oops

var ErrUnexpected = Define(BlameUnknown, NamespaceUnknown, ReasonUnexpected)

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
