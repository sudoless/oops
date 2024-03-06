package oops

import (
	"errors"
)

// Explainf is a helper function to check the given error if it's an Error and then call Error.Explainf with the given
// format and arguments, if and only if it's also not nil. If the given error is not an Error, it will be wrapped with
// ErrUncaught and the format and arguments will be passed to it.
func Explainf(err error, format string, args ...any) Error {
	if err == nil {
		return nil
	}

	v, ok := err.(Error)
	if !ok {
		return ErrUncaught.Wrapf(err, format, args...)
	}

	if v == nil {
		return nil
	}

	v.Explainf(format, args...)

	return v
}

// Nest is a shortcut to ErrorDefined.Yeet followed by a call to Error.Append, if and only if the source is not nil and
// the given errors are not empty.
func Nest(source ErrorDefined, nested ...Error) Error {
	if source == nil || len(nested) == 0 {
		return nil
	}

	return source.Yeet().Append(nested...)
}

// DeepIs will check if the given err is an Error and if the Error.Source matches the target ErrorDefined. If the given
// error is not an Error, it will attempt to traverse the unwrap chain until an Error is found or nil is reached. Once
// an Error is found, the check is repeated strictly on Error.Nested errors and never up to the parent of any errors.
// If any of the nested errors' source matches the target, true is returned. The check is repeated recursively until
// either the check is successful, or the nested errors exhaust.
// This function respects nil as valid targets (compared to DeepAs which does not).
func DeepIs(err error, target ErrorDefined) bool {
	if err == nil {
		return target == nil
	}

	v, ok := err.(Error)
	if !ok {
		return DeepIs(errors.Unwrap(err), target)
	}

	if v == nil {
		return target == nil
	}

	if v.Source() == target {
		return true
	}

	for _, nested := range v.Nested() {
		if DeepIs(nested, target) {
			return true
		}
	}

	return false
}

// As will check if the given err is an Error and if the Error.Source matches the target ErrorDefined, at which point
// err gets returned as an Error. If the given err is not an Error, or if the Error.Source does not match, the check
// is repeated with the parent of err (if any) until either the check is successful, or the parent is nil.
func As(err error, target ErrorDefined) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	if !ok {
		return As(errors.Unwrap(err), target)
	}

	if v == nil {
		return nil, false
	}

	if v.Source() == target {
		return v, true
	}

	return As(v.Unwrap(), target)
}

// AsAny will check if the given err is an Error and if so, return it as an Error.
func AsAny(err error) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	return v, ok
}

// DeepAs will check if the given err is an Error and if the Error.Source matches the target ErrorDefined, at which
// point err gets returned as an Error. If the given err is not an Error, it will attempt to traverse the unwrap chain
// until an Error is found or nil is reached. Once an Error is found, the check is repeated strictly on Error.Nested
// errors and never up to the parent of any errors. If any of the nested errors' source matches the target, the
// nested error is returned. The check is repeated recursively until either the check is successful, or the nested
// errors exhaust.
func DeepAs(err error, target ErrorDefined) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	if !ok {
		return DeepAs(errors.Unwrap(err), target)
	}

	if v == nil {
		return nil, false
	}

	if v.Source() == target {
		return v, true
	}

	for _, nested := range v.Nested() {
		if result, ok := DeepAs(nested, target); ok {
			return result, true
		}
	}

	return nil, false
}
