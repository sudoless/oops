package oops

import "slices"

// CollectorFinish finalizes a collection, returning nil if no errors were added.
type CollectorFinish = func() error

// CollectorAdd appends an error to the collection with an optional path label.
type CollectorAdd = func(err error, path string, args ...any)

// Collect returns a finish function and an add function for accumulating errors.
// The finish function returns nil if no errors were added; otherwise a new Error
// from d whose children are the errors added so far. add skips nil errors
// (including a typed-nil *Error) and wraps a non-oops error with ErrForeign.
// A non-empty path overwrites the added error's path in place.
// Neither function is safe for concurrent use.
func (d *ErrorDefinition) Collect() (CollectorFinish, CollectorAdd) {
	errs := make([]error, 0, 4)

	finish := func() error {
		if len(errs) == 0 {
			return nil
		}

		e := d.newError(1)
		e.wrapped = slices.Clone(errs)
		return e
	}

	addf := func(err error, path string, args ...any) {
		if isNil(err) {
			return
		}

		oErr, ok := err.(*Error) //nolint:errorlint // direct type check: sets path on the concrete *Error
		if !ok {
			oErr = ErrForeign.newError(1)
			oErr.wrapped = append(oErr.wrapped, err)
		}

		if path != "" {
			_ = oErr.Pathf(path, args...)
		}

		errs = append(errs, oErr)
	}

	return finish, addf
}
