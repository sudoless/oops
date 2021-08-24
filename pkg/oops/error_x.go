package oops

// Multi appends to the internal Error.multi list which can then be read back using Error.Multiples().
func (e *Error) Multi(multi ...string) *Error {
	if len(multi) == 0 {
		return e
	}

	e.multi = append(e.multi, multi...)
	return e
}

func (e *Error) explain(explanation string) {
	if e.explanation.Len() != 0 {
		e.explanation.WriteString(", ")
	}

	e.explanation.WriteString(explanation)
}
