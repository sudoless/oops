package oops

import (
	"fmt"
	"slices"
	"strings"
)

// Pathf sets the error's path label, replacing any previous path and path args,
// and returns the receiver for chaining. The rendered string is stored in Path();
// a copy of the raw args is stored in PathArgs() only when len(args) > 0, otherwise
// PathArgs() is nil. An empty format clears the path.
func (err *Error) Pathf(format string, args ...any) *Error {
	if err == nil {
		return nil
	}

	err.path = format
	if len(args) > 0 || strings.Contains(format, "%") {
		err.path = fmt.Sprintf(format, args...)
	}

	err.pathArgs = nil
	if len(args) > 0 {
		err.pathArgs = slices.Clone(args)
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

// PathArgs returns the path args. The returned slice is the error's own
// storage: do not modify it.
func (err *Error) PathArgs() []any {
	if err == nil {
		return nil
	}

	return err.pathArgs
}
