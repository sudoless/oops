package oops

import "fmt"

// Explainf appends a formatted explanation.
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

	if len(args) == 0 {
		err.explanation.WriteString(format)
		return err
	}

	err.explanation.WriteString(fmt.Sprintf(format, args...))
	return err
}

// Set stores a field value.
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

// AddCause appends semantic cause tags.
func (err *Error) AddCause(causes ...string) *Error {
	if err == nil {
		return nil
	}

	err.causes = append(err.causes, causes...)
	return err
}

// SetActions replaces the action tags (not accumulated).
func (err *Error) SetActions(actions ...string) *Error {
	if err == nil {
		return nil
	}

	err.actions = actions
	return err
}

// Nest adds an error to the wrapped slice.
func (err *Error) Nest(other error) *Error {
	if err == nil {
		return nil
	}

	if other != nil {
		err.wrapped = append(err.wrapped, other)
	}
	return err
}

// Append adds typed errors to the wrapped slice.
func (err *Error) Append(errs ...*Error) *Error {
	if err == nil {
		return nil
	}

	for _, e := range errs {
		if e != nil {
			err.wrapped = append(err.wrapped, e)
		}
	}
	return err
}
