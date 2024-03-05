package oops

func Explainf(err error, format string, args ...any) Error {
	if err == nil {
		return nil
	}

	v, ok := err.(Error)
	if !ok {
		return ErrUncaught.Wrapf(err, format, args...)
	}

	if v == nil {
		return nil
	}

	v.Explainf(format, args...)

	return v
}

func Nest(source ErrorDefined, nested ...Error) Error {
	if source == nil || len(nested) == 0 {
		return nil
	}

	finish, addf := source.Collect()
	for _, err := range nested {
		addf(err, "")
	}

	return finish()
}

func DeepIs(err Error, target ErrorDefined) bool {
	if err.Is(target) {
		return true
	}

	for _, nested := range err.Nested() {
		if DeepIs(nested, target) {
			return true
		}
	}

	return false
}

func As(err error) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	return v, ok
}
