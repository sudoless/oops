package oops

// Trace returns the formatted stack frames captured at creation time.
func (err *Error) Trace() []string {
	if err == nil {
		return nil
	}

	return err.trace
}
