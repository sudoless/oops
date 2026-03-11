package oops

import (
	"slices"
	"strings"
)

// Error is a live error instance created from an ErrorDefinition. It holds
// accumulated context: explanations, cause/action tags, wrapped errors,
// path segments, arbitrary fields, an optional stack trace.
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
	if err == nil {
		return ""
	}

	return err.def.code
}

// Message returns the definition's public-facing message.
func (err *Error) Message() string {
	if err == nil {
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
func (err *Error) HasCause(cause string) bool {
	if err == nil {
		return false
	}

	return slices.Contains(err.causes, cause)
}

// HasAction reports whether the error has the given action tag.
func (err *Error) HasAction(action string) bool {
	if err == nil {
		return false
	}

	return slices.Contains(err.actions, action)
}
