package oops

import (
	"errors"
	"strings"
)

type Error struct {
	parent      error
	defined     *errorDefined
	explanation strings.Builder
	multi       []string
	trace       []string
	blame       Blame
	namespace   Namespace
	reason      Reason
	code        string
}

// Is acts as a shortcut to calling errors.Is(e, err). Is will check if the target err is a errorDefined or another
// Error type, in which case matching is done as such, otherwise errors.Is, is used as a last call.
func (e *Error) Is(err error) bool {
	if err == nil {
		return e == nil
	}

	errDefined, ok := err.(*errorDefined)
	if ok {
		if errDefined.defined != nil {
			return errDefined.defined == e.defined
		}

		panic("oops: defined error must have the defined *errorDefined field allocated")
	}

	errError, ok := err.(*Error)
	if ok {
		return e.defined == errError.defined
	}

	if e.parent == nil {
		return false
	}

	return errors.Is(e.parent, err)
}

func (e *Error) Error() string {
	return e.Code()
}

func (e *Error) String() string {
	return e.Code()
}

func (e *Error) Unwrap() error {
	return e.parent
}
