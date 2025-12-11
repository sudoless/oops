package oops

import "errors"

// Nest is a shortcut to ErrorDefined.Yeet followed by a call to Error.Append, if and only if the source is not nil and
// the given errors are not empty.
func Nest(source ErrorDefined, nested ...Error) Error {
	if source == nil || len(nested) == 0 {
		return nil
	}

	return source.Yeet().Append(nested...)
}

// NestedAs will check if the given err is an Error and if the Error.Source matches the target ErrorDefined, at which
// point err gets returned as an Error. If the given err is not an Error, it will attempt to traverse the unwrap chain
// until an Error is found or nil is reached. Once an Error is found, the check is repeated strictly on Error.Nested
// errors and never up to the parent of any errors. If any of the nested errors' source matches the target, the
// nested error is returned. The check is repeated recursively until either the check is successful, or the nested
// errors exhaust.
func NestedAs(err error, target ErrorDefined) (Error, bool) {
	if err == nil {
		return nil, false
	}

	v, ok := err.(Error)
	if !ok {
		return NestedAs(errors.Unwrap(err), target)
	}

	if v == nil {
		return nil, false
	}

	if v.Source() == target {
		return v, true
	}

	for _, nested := range v.Nested() {
		if result, ok := NestedAs(nested, target); ok {
			return result, true
		}
	}

	return nil, false
}

// NestedIs will check if the given err is an Error and if the Error.Source matches the target ErrorDefined. If the given
// error is not an Error, it will attempt to traverse the unwrap chain until an Error is found or nil is reached. Once
// an Error is found, the check is repeated strictly on Error.Nested errors and never up to the parent of any errors.
// If any of the nested errors' source matches the target, true is returned. The check is repeated recursively until
// either the check is successful, or the nested errors exhaust.
// This function respects nil as valid targets (compared to NestedAs which does not).
func NestedIs(err error, target ErrorDefined) bool {
	if err == nil {
		return target == nil
	}

	v, ok := err.(Error)
	if !ok {
		return NestedIs(errors.Unwrap(err), target)
	}

	if v == nil {
		return target == nil
	}

	if v.Source() == target {
		return true
	}

	for _, nested := range v.Nested() {
		if NestedIs(nested, target) {
			return true
		}
	}

	return false
}
