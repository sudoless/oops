package oops

// Catch extracts an *Error from err, or wraps a non-oops error with ErrUncaught.
func Catch(err error) *Error {
	if err == nil {
		return nil
	}

	if v, ok := err.(*Error); ok { //nolint:errorlint // direct check: Catch does not traverse wrapped chains
		return v
	}

	return ErrUncaught.Wrap(err)
}

// Assert extracts an *Error from err and return true, or false if nil, or false and ErrCaught if non-oops error.
func Assert(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}

	if v, ok := err.(*Error); ok { //nolint:errorlint // direct check: Catch does not traverse wrapped chains
		return v, true
	}

	return ErrUncaught.Wrap(err), false
}

// Explainf adds a formatted explanation, wrapping non-oops errors with ErrUncaught.
func Explainf(err error, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return Catch(err).Explainf(format, args...)
}

// AddCause appends cause tags, wrapping non-oops errors with ErrUncaught.
func AddCause(err error, causes ...string) *Error {
	if err == nil {
		return nil
	}
	return Catch(err).AddCause(causes...)
}

// Pathf appends a formatted path segment, wrapping non-oops errors with ErrUncaught.
func Pathf(err error, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return Catch(err).WithPathf(format, args...)
}

// As traverses the unwrap chain to find an *Error whose definition matches target.
func As(err error, target *ErrorDefinition) (*Error, bool) {
	if err == nil || target == nil {
		return nil, false
	}

	if v, ok := err.(*Error); ok { //nolint:errorlint // As implements custom traversal; direct node check
		return asOopsError(v, target)
	}

	return asWrapped(err, target)
}

// asOopsError searches an *Error node and its wrapped children for target.
func asOopsError(v *Error, target *ErrorDefinition) (*Error, bool) {
	if v.def.is(target) {
		return v, true
	}

	for _, w := range v.wrapped {
		if found, ok := As(w, target); ok {
			return found, true
		}
	}

	return nil, false
}

// asWrapped handles non-*Error nodes by dispatching on standard unwrap interfaces.
func asWrapped(err error, target *ErrorDefinition) (*Error, bool) {
	switch vv := err.(type) { //nolint:errorlint // type switch is the traversal mechanism for non-oops errors
	case interface{ Unwrap() error }:
		return As(vv.Unwrap(), target)
	case interface{ Unwrap() []error }:
		for _, e := range vv.Unwrap() {
			if found, ok := As(e, target); ok {
				return found, true
			}
		}
	}

	return nil, false
}

// Nest creates a new Error from def with the given errors as wrapped children.
// Returns nil if def is nil or all errors are nil.
func Nest(def *ErrorDefinition, errs ...error) *Error {
	if def == nil {
		return nil
	}

	var filtered []error
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	e := def.newError()
	e.wrapped = filtered

	return e
}
