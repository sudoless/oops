package oops

var (
	// ErrForeign wraps non-oops errors at the package boundary (Foreign, the
	// package helpers, and Collect).
	ErrForeign = Define("foreign").Causes(CauseInternal).Actions(ActionAbort).Traced()

	// ErrTODO is a placeholder for unimplemented paths.
	ErrTODO = Define("todo").Causes(CauseInternal).Actions(ActionAbort).Traced().Message("not implemented")
)
