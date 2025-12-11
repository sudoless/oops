package oops

type Error interface {
	// Error returns the string representation of the error. The format of the string is implementation-specific.
	Error() string

	// Unwrap should return the parent error, if any. This method may have implementation-specific behaviour.
	Unwrap() error

	// Append should add the given errors as "nested" errors returnable by calling Error.Nested.
	Append(errs ...Error) Error

	// Nested should return the list of errors that were added using Error.Append or by other means (such as collection
	// of a ErrorDefined.Collect).
	Nested() []Error

	// Is returns true if both errors are nil, or if the other error is a ErrorDefined and the Error.Source matches it,
	// or if the other error is an Error, then both their Error.Source must match. Otherwise, the behaviour can be
	// implementation-specific, but it's recommended to at least check any parent error if any.
	Is(other error) bool

	// As should check the target type to determine if it's a *Error, and if so, set the target to itself. This method
	// must be implemented such that errors.As can be used to "cast" any error to an Error.
	As(any) bool

	// Explainf should update the error's explanation with the given format and arguments. The format and arguments
	// should be used to create a human-readable explanation of the error. This method may have implementation-specific
	// behaviour.
	Explainf(format string, args ...any)

	// Set should add or update the value for the given key in the context of the Error. If set, values should be
	// retrievable using Error.Get or Error.GetAll. This method may have implementation-specific behaviour.
	Set(key string, value any) Error

	// Get should return the value for the given key in the context of the Error. If the key is not found, the second
	// return value should be false. This method may have implementation-specific behaviour.
	Get(key string) (value any, ok bool)

	// GetAll should return all the values in the context of the Error. This method may have implementation-specific
	// behaviour.
	GetAll() map[string]any

	// PathSetf should set the path and args for the error. The path meaning is implementation-specific.
	PathSetf(path string, args ...any) Error

	// Path should return the formatted path of the error. The path meaning is implementation-specific.
	Path() string

	// PathArgs should return the arguments set by Error.PathSetf. The path meaning is implementation-specific.
	PathArgs() []any

	// Explanation should return the complete current explanation of the error.
	Explanation() string

	// Trace optionally returns a list of strings, where each string represents a step in the error's trace. The trace
	// format is implementation-specific.
	Trace() []string

	// Source must return the ErrorDefined that created/spawned/returned this error.
	Source() ErrorDefined
}

type ErrorDefined interface {
	// Error will panic, as ErrorDefined is not intended to be used as a replacement for `error`, but it can
	// be passed as an `error` to specific functions.
	Error() string

	Yeet() Error
	Yeetf(format string, args ...any) Error

	Wrap(err error) Error
	Wrapf(err error, format string, args ...any) Error

	// Collect returns a ErrorCollectorAdd function that appends errors to Error.Nested and a ErrorCollectorFinish
	// that will return an Error with ErrorDefined as the source, if any non-nil Error were added with the collector.
	// Otherwise, nil is returned. It is safe to use both functions with nils and without checks.
	Collect() (finish ErrorCollectorFinish, addf ErrorCollectorAdd)

	Is(other error) bool
}

type ErrorCollectorFinish func() Error

type ErrorCollectorAdd func(err error, path string, args ...any)

type Formatter func(err Error) string
