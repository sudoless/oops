package oops

import "fmt"

func (e *Error) explain(explanation string) {
	if explanation == "" {
		return
	}

	if e.explanation.Len() != 0 {
		e.explanation.WriteString(", ")
	}

	e.explanation.WriteString(explanation)
}

// Explain appends to the explanation string.
func (e *Error) Explain(explanation string, args ...any) *Error {
	if len(args) > 0 {
		e.explain(fmt.Sprintf(explanation, args...))
	} else {
		e.explain(explanation)
	}

	return e
}
