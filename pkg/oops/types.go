package oops

type Error interface {
	Error() string

	Unwrap() error
	Nested() []Error

	Is(other error) bool
	As(any) bool

	Explainf(format string, args ...any)
	Set(key string, value any) Error
	Get(key string) (value any, ok bool)
	GetAll() map[string]any

	Path() string
	PathArgs() []any
	PathSetf(path string, args ...any) Error

	Explanation() string
	Trace() []string
	Source() ErrorDefined
}

type ErrorDefined interface {
	Error() string

	Yeet() Error
	Yeetf(format string, args ...any) Error

	Wrap(err error) Error
	Wrapf(err error, format string, args ...any) Error

	Collect() (finish ErrorCollectorFinish, addf ErrorCollectorAdd)

	Is(other error) bool
}

type ErrorCollectorFinish func() Error

type ErrorCollectorAdd func(err error, path string, args ...any)

type Formatter func(err Error) string
