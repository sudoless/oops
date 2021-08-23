package oops

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

func (e *Error) Stack() *Error {
	if len(e.trace) != 0 {
		return e
	}

	e.trace = stack(1)
	return e
}
