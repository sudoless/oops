package oops

import (
	"errors"
	"strings"
)

type Error struct {
	parent      error
	source      *ErrorDefined
	explanation strings.Builder
	trace       []string
	fields      []field
}

// Is acts as a shortcut to calling errors.Is(e, err). Is will check if the target err is a ErrorDefined or another
// Error type, in which case matching is done as such, otherwise errors.Is, is used as a last call.
func (e *Error) Is(err error) bool {
	if err == nil {
		return e == nil
	}

	errDefined, ok := err.(*ErrorDefined)
	if ok {
		return e.source == errDefined
	}

	errError, ok := err.(*Error)
	if ok {
		return e.source == errError.source
	}

	if e.parent == nil {
		return false
	}

	return errors.Is(e.parent, err)
}

// Error returns the error string message as returned by Error.String.
func (e *Error) Error() string {
	return e.String()
}

// String returns the error type in square brackets, followed by the code, followed by a :, followed by the explanation.
// The String method aims to be a very generic error representation and as such it's not recommended for production
// use, instead you should define your own representation, appropriate for your use case, using Error.Code, Error.Type,
// Error.Explain, Error.Multiples, and Error.Trace.
func (e *Error) String() string {
	return "[" + e.source.t + "] " + e.source.code + " : " + e.explanation.String()
}

// Unwrap returns the parent error which exists if the Error was created using a Wrap method.
func (e *Error) Unwrap() error {
	return e.parent
}
