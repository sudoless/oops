package oops

import "fmt"

func (e *Error) explain(explanation string) {
	if e.explanation.Len() != 0 {
		e.explanation.WriteString(", ")
	}

	e.explanation.WriteString(explanation)
}

// Fields takes a list of key followed by value pairs and keeps them in a fields list. There is no field deduplication.
func (e *Error) Fields(kv ...string) *Error {
	kvLen := len(kv)
	if len(kv) == 0 || kvLen%2 != 0 {
		return e
	}

	if e.fields == nil {
		e.fields = make([]string, 0, kvLen)
	}

	e.fields = append(e.fields, kv...)

	return e
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
