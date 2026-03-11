package oops

// Causes appends semantic cause tags to this definition.
// Call only once per definition, at package initialisation time; builder
// methods mutate d in place and are not safe for concurrent use.
func (d *ErrorDefinition) Causes(causes ...string) *ErrorDefinition {
	d.causes = append(d.causes, causes...)
	return d
}

// Actions appends semantic action tags to this definition.
func (d *ErrorDefinition) Actions(actions ...string) *ErrorDefinition {
	d.actions = append(d.actions, actions...)
	return d
}

// Message sets the public-facing message for this definition.
func (d *ErrorDefinition) Message(msg string) *ErrorDefinition {
	d.message = msg
	return d
}

// Traced enables stack trace capture for errors from this definition.
func (d *ErrorDefinition) Traced() *ErrorDefinition {
	d.traced = true
	return d
}

// Inherits adds parent definitions to this definition's inheritance chain.
// errors.Is checks traverse the inherits chain.
func (d *ErrorDefinition) Inherits(defs ...*ErrorDefinition) *ErrorDefinition {
	d.inherits = append(d.inherits, defs...)
	return d
}

// SetFormatter sets a custom formatter for errors from this definition.
func (d *ErrorDefinition) SetFormatter(f Formatter) *ErrorDefinition {
	d.formatter = f
	return d
}
