package oops

// isNil reports whether err is nil or a typed-nil *Error.
func isNil(err error) bool {
	if err == nil {
		return true
	}

	v, ok := err.(*Error) //nolint:errorlint // direct check: only the top-level value is inspected
	return ok && v == nil
}

// Foreign returns err as an *Error: an *Error is returned as-is, any other
// error is wrapped with ErrForeign. It returns nil for a nil err, including a
// typed-nil *Error. Only the top-level value is inspected; wrapped chains are
// not traversed.
//
// The result is an *Error: check it for nil before returning it as an error,
// or a nil result becomes a non-nil error interface. The package helpers
// Explainf, AddCauses and Pathf return error and do this check for you.
func Foreign(err error) *Error {
	return foreign(err, 1)
}

// foreign implements Foreign. skip is passed on to newError, with 0
// identifying the caller of foreign.
func foreign(err error, skip int) *Error {
	if isNil(err) {
		return nil
	}

	if v, ok := err.(*Error); ok { //nolint:errorlint // direct check: Foreign does not traverse wrapped chains
		return v
	}

	e := ErrForeign.newError(skip + 1)
	e.wrapped = append(e.wrapped, err)

	return e
}

// Native reports whether err is an *Error and returns it. It returns (nil, false)
// for a nil err, a typed-nil *Error, or any other error type. Only the top-level
// value is inspected; wrapped chains are not traversed.
func Native(err error) (*Error, bool) {
	if isNil(err) {
		return nil, false
	}

	v, ok := err.(*Error) //nolint:errorlint // direct check: Native does not traverse wrapped chains
	return v, ok
}

// Explainf appends a formatted explanation to err, wrapping a non-oops error
// with ErrForeign first. It returns nil when err is nil.
func Explainf(err error, format string, args ...any) error {
	e := foreign(err, 1)
	if e == nil {
		return nil
	}

	return e.Explainf(format, args...)
}

// AddCauses appends cause tags to err, wrapping a non-oops error with ErrForeign
// first. It returns nil when err is nil.
func AddCauses(err error, causes ...Cause) error {
	e := foreign(err, 1)
	if e == nil {
		return nil
	}

	return e.AddCauses(causes...)
}

// Pathf sets the formatted path of err, wrapping a non-oops error with
// ErrForeign first. It returns nil when err is nil.
func Pathf(err error, format string, args ...any) error {
	e := foreign(err, 1)
	if e == nil {
		return nil
	}

	return e.Pathf(format, args...)
}

// As traverses the unwrap tree depth-first to find the first *Error whose
// definition is, or inherits, target. Nil and typed-nil nodes are skipped.
func As(err error, target *ErrorDefinition) (*Error, bool) {
	if isNil(err) || target == nil {
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
// It returns nil if def is nil or every error is nil (including typed-nil *Error).
func Nest(def *ErrorDefinition, errs ...error) error {
	if def == nil {
		return nil
	}

	var filtered []error
	for _, err := range errs {
		if !isNil(err) {
			filtered = append(filtered, err)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	e := def.newError(1)
	e.wrapped = filtered

	return e
}
