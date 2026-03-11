package oops

// Trace returns the formatted stack frames captured at creation time.
func (err *Error) Trace() []string {
	return err.trace
}
