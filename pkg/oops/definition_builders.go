package oops

import "slices"

// Causes returns a copy of the definition with causes appended to its cause tags.
func (d *ErrorDefinition) Causes(causes ...Cause) *ErrorDefinition {
	c := *d
	c.causes = slices.Concat(d.causes, causes)
	return &c
}

// Actions returns a copy of the definition with actions appended to its action tags.
func (d *ErrorDefinition) Actions(actions ...Action) *ErrorDefinition {
	c := *d
	c.actions = slices.Concat(d.actions, actions)
	return &c
}

// Message returns a copy of the definition with its public-facing message set to msg.
func (d *ErrorDefinition) Message(msg string) *ErrorDefinition {
	c := *d
	c.message = msg
	return &c
}

// Traced returns a copy of the definition that captures a stack trace for each Error it creates.
func (d *ErrorDefinition) Traced() *ErrorDefinition {
	c := *d
	c.traced = true
	return &c
}

// Inherits returns a copy of the definition with defs appended to its parents.
// errors.Is and As match an Error against its definition's parents, transitively.
func (d *ErrorDefinition) Inherits(defs ...*ErrorDefinition) *ErrorDefinition {
	c := *d
	c.inherits = slices.Concat(d.inherits, defs)
	return &c
}

// Formatter returns a copy of the definition whose errors render with f instead
// of the default code[: message][; explanation].
func (d *ErrorDefinition) Formatter(f Formatter) *ErrorDefinition {
	c := *d
	c.formatter = f
	return &c
}
