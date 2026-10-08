package oops

import (
	"fmt"
	"strings"
)

// Explainf appends a formatted explanation. Mutates Error, returned for chaining.
func (err *Error) Explainf(format string, args ...any) *Error {
	if err == nil {
		return nil
	}

	if format == "" {
		return err
	}

	if err.explanation.Len() != 0 {
		err.explanation.WriteString(", ")
	}

	if len(args) == 0 && !strings.Contains(format, "%") {
		err.explanation.WriteString(format)
		return err
	}

	_, _ = fmt.Fprintf(&err.explanation, format, args...)
	return err
}

// Set stores a field value. Mutates Error, returned for chaining.
func (err *Error) Set(key string, value any) *Error {
	if err == nil {
		return nil
	}

	if err.fields == nil {
		err.fields = make(map[string]any, 4)
	}
	err.fields[key] = value
	return err
}

// AddCauses appends semantic cause tags. Mutates Error, returned for chaining.
func (err *Error) AddCauses(causes ...string) *Error {
	if err == nil {
		return nil
	}

	err.causes = append(err.causes, causes...)
	return err
}

// SetActions replaces the action tags (not accumulated). Mutates Error, returned for chaining.
func (err *Error) SetActions(actions ...string) *Error {
	if err == nil {
		return nil
	}

	err.actions = actions
	return err
}

// Nest appends errs to the wrapped children, skipping nil and typed-nil *Error
// values. Mutates Error, returned for chaining.
func (err *Error) Nest(errs ...error) *Error {
	if err == nil {
		return nil
	}

	for _, other := range errs {
		if !isNil(other) {
			err.wrapped = append(err.wrapped, other)
		}
	}
	return err
}
