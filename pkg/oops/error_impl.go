package oops

import (
	"fmt"
	"strings"
)

var _ Error = &errorImpl{}

//nolint:errname
type errorImpl struct {
	source *errorDefined

	parent error
	nested []Error

	path     string
	pathArgs []any
	props    map[string]any

	trace       []string
	explanation strings.Builder
}

func (err *errorImpl) Nested() []Error {
	return err.nested
}

func (err *errorImpl) Append(errs ...Error) Error { //nolint:ireturn
	err.nested = append(err.nested, errs...)
	return err
}

func (err *errorImpl) GetAll() map[string]any {
	return err.props
}

func (err *errorImpl) Set(key string, value any) Error { //nolint:ireturn
	if err.props == nil {
		err.props = make(map[string]any, 4)
	}

	err.props[key] = value

	return err
}

func (err *errorImpl) Get(key string) (value any, ok bool) {
	if err.props == nil {
		return nil, false
	}

	value, ok = err.props[key]
	return value, ok
}

func (err *errorImpl) Explanation() string {
	return err.explanation.String()
}

func (err *errorImpl) Trace() []string {
	return err.trace
}

func (err *errorImpl) Source() ErrorDefined { //nolint:ireturn
	return err.source
}

func (err *errorImpl) Path() string {
	return err.path
}

func (err *errorImpl) PathArgs() []any {
	return err.pathArgs
}

func (err *errorImpl) PathSetf(path string, args ...any) Error { //nolint:ireturn
	if len(args) == 0 {
		err.path = path
	} else {
		err.path = fmt.Sprintf(path, args...)
		err.pathArgs = args
	}

	return err
}

func (err *errorImpl) Explainf(format string, args ...any) {
	if format == "" {
		return
	}

	if err.explanation.Len() != 0 {
		err.explanation.WriteString(", ")
	}

	if len(args) == 0 {
		err.explanation.WriteString(format)
		return
	}

	err.explanation.WriteString(fmt.Sprintf(format, args...))
}
