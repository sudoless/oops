package oops

import (
	"maps"
	"slices"
	"strings"
)

// Error is a live error instance created from an ErrorDefinition. It holds
// accumulated context: explanations, cause/action tags, wrapped errors,
// a path label, arbitrary fields, an optional stack trace.
//
// Use Error only through *Error. Do not copy the struct value: the explanation
// is a strings.Builder, so explaining a copy of an Error that already has an
// explanation panics. Use Clone to get an independent copy.
type Error struct {
	def *ErrorDefinition

	causes  []Cause
	actions []Action
	wrapped []error

	path     string
	pathArgs []any
	fields   map[string]any

	trace       []string
	explanation strings.Builder
}

// Definition returns the ErrorDefinition that created this error.
func (err *Error) Definition() *ErrorDefinition {
	if err == nil {
		return nil
	}

	return err.def
}

// Code returns the definition's identity code.
func (err *Error) Code() string {
	if err == nil || err.def == nil {
		return ""
	}

	return err.def.code
}

// Message returns the definition's public-facing message.
func (err *Error) Message() string {
	if err == nil || err.def == nil {
		return ""
	}

	return err.def.message
}

// Explanation returns all explanations joined by ", ".
func (err *Error) Explanation() string {
	if err == nil {
		return ""
	}

	return err.explanation.String()
}

// Causes returns the cause tags.
func (err *Error) Causes() []Cause {
	if err == nil {
		return nil
	}

	return err.causes
}

// Actions returns the action tags.
func (err *Error) Actions() []Action {
	if err == nil {
		return nil
	}

	return err.actions
}

// Fields returns the raw field map.
func (err *Error) Fields() map[string]any {
	if err == nil {
		return nil
	}

	return err.fields
}

// Get returns the raw value for the given field key.
func (err *Error) Get(key string) (any, bool) {
	if err == nil || err.fields == nil {
		return nil, false
	}

	v, ok := err.fields[key]
	return v, ok
}

// HasCause reports whether the error has the given cause tag.
func (err *Error) HasCause(cause Cause) bool {
	if err == nil {
		return false
	}

	return slices.Contains(err.causes, cause)
}

// HasAction reports whether the error has the given action tag.
func (err *Error) HasAction(action Action) bool {
	if err == nil {
		return false
	}

	return slices.Contains(err.actions, action)
}

// Clone returns a shallow copy of err that can be annotated without affecting
// err: the copy owns its explanation, causes, actions, fields map, path, path
// args and list of wrapped errors. The wrapped errors themselves, the trace and
// field values are shared.
func (err *Error) Clone() *Error {
	if err == nil {
		return nil
	}

	c := &Error{
		def:      err.def,
		causes:   slices.Clone(err.causes),
		actions:  slices.Clone(err.actions),
		wrapped:  slices.Clone(err.wrapped),
		path:     err.path,
		pathArgs: slices.Clone(err.pathArgs),
		fields:   maps.Clone(err.fields),
		trace:    err.trace,
	}
	c.explanation.WriteString(err.explanation.String())

	return c
}
