package oops

import stack "go.sdls.io/oops/v2/internal/unsafe"

// ErrorDefinition is a sentinel error definition created once at package level via Define.
// It holds identity (code), semantic tags (causes, actions), a public-facing message,
// and optional configuration (tracing, formatting, inheritance).
//
//nolint:errname // ErrorDefinition is a sentinel definition, not an error value
type ErrorDefinition struct {
	code      string
	causes    []Cause
	actions   []Action
	message   string
	traced    bool
	inherits  []*ErrorDefinition
	formatter Formatter
}

// Define creates a new ErrorDefinition with the given code.
// Call only once per definition, at package initialisation time; builder
// methods mutate d in place and are not safe for concurrent use.
func Define(code string) *ErrorDefinition {
	return &ErrorDefinition{code: code}
}

func (d *ErrorDefinition) newError() *Error {
	e := &Error{def: d}

	if len(d.causes) > 0 {
		e.causes = make([]string, len(d.causes))
		copy(e.causes, d.causes)
	}

	if len(d.actions) > 0 {
		e.actions = make([]string, len(d.actions))
		copy(e.actions, d.actions)
	}

	if d.traced {
		// skip=3: Stack(0) + newError(1) + public method(2) → user at frame 3
		e.trace = stack.Stack(3)
	}

	return e
}

// Code returns the definition's identity code.
func (d *ErrorDefinition) Code() string { return d.code }

// Error returns "code: message" or just "code" if message is empty.
func (d *ErrorDefinition) Error() string {
	if d.message == "" {
		return d.code
	}
	return d.code + ": " + d.message
}

// Is reports whether other is this definition or a definition it inherits.
// An *Error target never matches: match a live error against a definition
// with errors.Is(err, def).
func (d *ErrorDefinition) Is(other error) bool {
	if isNil(other) {
		return false
	}

	if v, ok := other.(*ErrorDefinition); ok {
		return d.is(v)
	}

	return false
}

// is checks identity including the inherits chain. A nil d or target never matches.
func (d *ErrorDefinition) is(target *ErrorDefinition) bool {
	if d == nil || target == nil {
		return false
	}

	if d == target {
		return true
	}

	for _, parent := range d.inherits {
		if parent.is(target) {
			return true
		}
	}

	return false
}

// Yeet creates a new Error from this definition.
func (d *ErrorDefinition) Yeet() *Error {
	return d.newError()
}

// Yeetf creates a new Error with a formatted explanation.
func (d *ErrorDefinition) Yeetf(format string, args ...any) *Error {
	e := d.newError()
	return e.Explainf(format, args...)
}

// Wrap creates a new Error that wraps the given error. A nil err, including a
// typed-nil *Error, is not wrapped.
func (d *ErrorDefinition) Wrap(err error) *Error {
	e := d.newError()
	if !isNil(err) {
		e.wrapped = append(e.wrapped, err)
	}
	return e
}

// Wrapf creates a new Error that wraps the given error with a formatted explanation.
func (d *ErrorDefinition) Wrapf(err error, format string, args ...any) *Error {
	e := d.newError()
	if !isNil(err) {
		e.wrapped = append(e.wrapped, err)
	}
	return e.Explainf(format, args...)
}
