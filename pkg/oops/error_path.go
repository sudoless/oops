package oops

import (
	"fmt"
)

// Pathf sets the error's path, replacing any previous one, and returns the
// receiver for chaining. The rendered string is stored in Path(); the raw args are
// stored in PathArgs() only when len(args) > 0 — callers that need to reconstruct
// the original format string should store it separately.
func (err *Error) Pathf(format string, args ...any) *Error {
	if err == nil {
		return nil
	}

	if format == "" {
		return err
	}

	if len(args) > 0 {
		err.path = fmt.Sprintf(format, args...)
		err.pathArgs = args
	} else {
		err.path = format
	}

	return err
}

// Path returns the formatted path.
func (err *Error) Path() string {
	if err == nil {
		return ""
	}

	return err.path
}

// PathArgs returns the path args.
func (err *Error) PathArgs() []any {
	if err == nil {
		return nil
	}

	return err.pathArgs
}
