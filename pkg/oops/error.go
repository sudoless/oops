package oops

import (
	"errors"
	"strings"
)

type Error struct {
	blame     Blame
	namespace Namespace
	reason    Reason

	defined *errorDefined
	trace   []string

	explanation strings.Builder
	multi       []string
	help        string

	parent error
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return e == nil
	}

	errDefined, ok := err.(*errorDefined)
	if ok {
		return errDefined == e.defined
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
