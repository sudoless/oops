package oops

import "go.sdls.io/oops/v2/internal/stack"

// ErrorDefinition is a sentinel error definition created once at package level via Define.
// It holds identity (code), semantic tags (causes, actions), a public-facing message,
// and optional configuration (tracing, formatting, inheritance).
//
// A definition's identity is its pointer. Every builder (Causes, Actions, Message,
// Traced, Inherits, Formatter) returns a NEW definition and leaves the receiver
// unchanged, so the sentinel is the pointer returned by the last builder in the
// chain. Errors created from a derived definition do not match the receiver
// under errors.Is or As; to derive a related sentinel that does, use
// Define(code).Inherits(parent).
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

// Define creates a new ErrorDefinition with the given code. Call it once per
// definition, at package level. Builder methods return a new definition, so
// the variable must hold the result of the whole chain.
func Define(code string) *ErrorDefinition {
	return &ErrorDefinition{code: code}
}

func (d *ErrorDefinition) newError() *Error {
	e := &Error{def: d}

	if len(d.causes) > 0 {
		e.causes = make([]Cause, len(d.causes))
		copy(e.causes, d.causes)
	}

	if len(d.actions) > 0 {
		e.actions = make([]Action, len(d.actions))
		copy(e.actions, d.actions)
	}

	if d.traced {
		// skip=3: Stack(0) + newError(1) + public method(2) → user at frame 3
		e.trace = stack.Stack(3)
	}

	return e
}

// Code returns the definition's identity code.
func (d *ErrorDefinition) Code() string {
	if d == nil {
		return ""
	}

	return d.code
}

// Error panics. An ErrorDefinition implements error only so it can be the
// target of errors.Is; it is never an error value itself. Create an Error with
// Yeet, Yeetf, Wrap or Wrapf and return that instead.
func (d *ErrorDefinition) Error() string {
	panic("oops: ErrorDefinition " + d.Code() + " used as an error value; create an Error with Yeet or Wrap")
}

// Is reports whether other is this definition or a definition it inherits.
// An *Error target never matches: match a live error against a definition
// with errors.Is(err, def).
func (d *ErrorDefinition) Is(other error) bool {
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
