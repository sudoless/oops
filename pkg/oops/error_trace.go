package oops

// Trace returns the formatted stack frames captured at creation time. The
// returned slice is the error's own storage: do not modify it.
func (err *Error) Trace() []string {
	if err == nil {
		return nil
	}

	return err.trace
}
