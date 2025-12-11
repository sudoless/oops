package oops

import (
	"strings"

	"go.sdls.io/oops/internal/unsafe"
)

var _ ErrorDefined = &errorDefined{}

type errorDefined struct {
	traced    bool
	props     map[string]any
	formatter Formatter
}

func (defined *errorDefined) newError(parent error) *errorImpl {
	e := &errorImpl{
		source:      defined,
		parent:      parent,
		explanation: strings.Builder{},
	}

	if len(defined.props) != 0 {
		e.props = make(map[string]any, len(defined.props))
		for k, v := range defined.props {
			e.props[k] = v
		}
	}

	if defined.traced {
		e.trace = unsafe.Stack(3)
	}

	return e
}

func (defined *errorDefined) Error() string {
	panic("oops: do not use ErrorDefined as error, use ErrorDefined.Yeet and ErrorDefined.Wrap")
}

func (defined *errorDefined) Yeet() Error {
	return defined.newError(nil)
}

func (defined *errorDefined) Yeetf(format string, args ...any) Error {
	err := defined.newError(nil)
	err.Explainf(format, args...)

	return err
}

func (defined *errorDefined) Wrap(err error) Error {
	return defined.newError(err)
}

func (defined *errorDefined) Wrapf(other error, format string, args ...any) Error {
	err := defined.newError(other)
	err.Explainf(format, args...)

	return err
}

func (defined *errorDefined) Collect() (ErrorCollectorFinish, ErrorCollectorAdd) {
	errs := make([]Error, 0, 4)

	finish := func() Error {
		if len(errs) == 0 {
			return nil
		}

		err := defined.newError(nil)
		err.nested = errs

		return err
	}

	addf := func(err error, path string, args ...any) {
		if err == nil {
			return
		}

		v, ok := err.(Error)
		if !ok {
			vd, ok := err.(ErrorDefined)
			if !ok {
				panic("oops: uncaught unwrapped error")
			}

			v = vd.Yeet()
		}

		if v == nil {
			return
		}

		errs = append(errs, v.PathSetf(path, args...))
	}

	return finish, addf
}

func (defined *errorDefined) Is(other error) bool {
	if other == nil {
		return false
	}

	vOther, ok := other.(Error)
	if ok {
		return vOther.Source() == defined
	}

	return false
}
