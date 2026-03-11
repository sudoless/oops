package oops

import "go.sdls.io/oops/internal/unsafe"

// ErrorDefinition is a sentinel error definition created once at package level via Define.
// It holds identity (code), semantic tags (causes, actions), a public-facing message,
// and optional configuration (tracing, formatting, inheritance).
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
		e.trace = unsafe.Stack(3)
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

// Is checks if other matches this definition or is an Error from this definition.
func (d *ErrorDefinition) Is(other error) bool {
	if other == nil {
		return false
	}

	switch v := other.(type) {
	case *ErrorDefinition:
		return d.is(v)
	case *Error:
		return v.def == d
	}

	return false
}

// is checks identity including the inherits chain.
func (d *ErrorDefinition) is(target *ErrorDefinition) bool {
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

// Wrap creates a new Error that wraps the given error.
func (d *ErrorDefinition) Wrap(err error) *Error {
	e := d.newError()
	if err != nil {
		e.wrapped = append(e.wrapped, err)
	}
	return e
}

// Wrapf creates a new Error that wraps the given error with a formatted explanation.
func (d *ErrorDefinition) Wrapf(err error, format string, args ...any) *Error {
	e := d.newError()
	if err != nil {
		e.wrapped = append(e.wrapped, err)
	}
	return e.Explainf(format, args...)
}
