package oops

import (
	"fmt"
)

// WithPathf sets a formatted path segment on the error and returns the receiver
// for chaining. The rendered string is stored in Path(); the raw args are stored
// in PathArgs() only when len(args) > 0 — callers that need to reconstruct the
// original format string should store it separately.
func (err *Error) WithPathf(format string, args ...any) *Error {
	if format == "" {
		return err
	}

	err.path = fmt.Sprintf(format, args...)
	if len(args) > 0 {
		err.pathArgs = args
	}

	return err
}

// Path returns the formatted path.
func (err *Error) Path() string { return err.path }

// PathArgs returns the path args.
func (err *Error) PathArgs() []any {
	return err.pathArgs
}
