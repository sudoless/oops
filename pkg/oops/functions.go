package oops

// Explainf is a helper function to check the given error if it's an Error and then call Error.Explainf with the given
// format and arguments, if and only if it's also not nil. If the given error is not an Error, it will be wrapped with
// ErrUncaught and the format and arguments will be passed to it.
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

// As will check if the given err is an Error and if the Error.Source matches the target ErrorDefined, at which point
// err gets returned as an Error. If the given err is not an Error, or if the Error.Source does not match, the check
// is repeated with the parent of err (if any) until either the check is successful, or the parent is nil.
// As does not check Error.Nested errors.
func As(err error, target ErrorDefined) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	if !ok {
		switch vv := err.(type) {
		case interface{ Unwarp() error }:
			return As(vv.Unwarp(), target)
		case interface{ Unwrap() []error }:
			for _, er := range vv.Unwrap() {
				if er == nil {
					continue
				}

				aer, ok := As(er, target)
				if ok {
					return aer, true
				}
			}

			return nil, false
		case interface{ Unwraps() []error }:
			for _, er := range vv.Unwraps() {
				if er == nil {
					continue
				}

				aer, ok := As(er, target)
				if ok {
					return aer, true
				}
			}

			return nil, false
		}
	}

	if v == nil {
		return nil, false
	}

	if v.Source() == target {
		return v, true
	}

	return As(v.Unwrap(), target)
}

// AsAny will check if the given err is an Error and if so, return it as an Error. AsAny does not check the unwrap chain.
func AsAny(err error) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	return v, ok
}

// AsMust will cast the given error as an Error, if the error does not implement Error, then it will become ErrUncaught.
// AsMust does not check the unwrap chain.
func AsMust(err error) Error {
	if err == nil {
		return nil
	}

	v, ok := err.(Error)
	if !ok {
		return ErrUncaught.Wrap(err)
	}

	return v
}
