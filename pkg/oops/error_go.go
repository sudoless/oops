package oops

import (
	"errors"
)

func (err *errorImpl) Error() string {
	if err == nil {
		return "oops.Error(nil)"
	}

	if err.source.formatter == nil {
		panic("oops: un-formatted error")
	}

	return err.source.formatter(err)
}

func (err *errorImpl) Unwrap() error {
	return err.parent
}

func (err *errorImpl) Is(other error) bool {
	if err == nil {
		return other == nil
	}

	vOther, ok := other.(ErrorDefined)
	if ok {
		return err.source == vOther
	}

	vErr, ok := other.(Error)
	if ok {
		return err.source == vErr.Source()
	}

	if err.parent == nil {
		return false
	}

	return errors.Is(err.parent, other)
}

func (err *errorImpl) As(other any) bool {
	otherErrPtr, ok := other.(*Error)
	if !ok {
		if _, ok = other.(*ErrorDefined); ok {
			panic("oops: cannot use ErrorDefined as target")
		}

		return false
	}

	*otherErrPtr = err
	return true
}
