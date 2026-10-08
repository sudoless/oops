package oops

// Reassigning a preset replaces its identity: errors created before the
// reassignment no longer match it. Reassign only at package initialisation,
// before any error is created.
var (
	// ErrForeign wraps non-oops errors at the package boundary (Foreign, the
	// package helpers, and Collect).
	ErrForeign = Define("foreign").Causes(CauseInternal).Actions(ActionAbort).Traced()

	// ErrTODO is a placeholder for unimplemented paths.
	ErrTODO = Define("todo").Causes(CauseInternal).Actions(ActionAbort).Traced().Message("not implemented")
)
