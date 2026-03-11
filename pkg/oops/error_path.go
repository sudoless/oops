package oops

import (
	"fmt"
)

// WithPathf appends a formatted path segment.
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
