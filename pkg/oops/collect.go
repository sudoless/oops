package oops

// CollectorFinish finalizes a collection, returning nil if no errors were added.
type CollectorFinish = func() error

// CollectorAdd appends an error to the collection with an optional path segment.
type CollectorAdd = func(err error, path string, args ...any)

// Collect returns a finish function and an add function for accumulating errors.
// The finish function returns nil if no errors were added.
// Neither function is safe for concurrent use. CollectorAdd called with a path
// will mutate the errors path in place.
func (d *ErrorDefinition) Collect() (CollectorFinish, CollectorAdd) {
	errs := make([]error, 0, 4)

	finish := func() error {
		if len(errs) == 0 {
			return nil
		}

		e := d.newError()
		e.wrapped = errs
		return e
	}

	addf := func(err error, path string, args ...any) {
		if err == nil {
			return
		}

		if oErr, ok := err.(*Error); ok { //nolint:errorlint // direct type check: sets path on the concrete *Error
			_ = oErr.WithPathf(path, args...)
			errs = append(errs, oErr)
			return
		}

		wrapped := ErrUncaught.newError()
		wrapped.wrapped = append(wrapped.wrapped, err)
		_ = wrapped.WithPathf(path, args...)

		errs = append(errs, wrapped)
	}

	return finish, addf
}
