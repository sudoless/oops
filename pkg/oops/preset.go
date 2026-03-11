package oops

var (
	// ErrUncaught wraps non-oops errors automatically.
	ErrUncaught = Define("uncaught").Causes(CauseInternal).Actions(ActionAbort).Traced()

	// ErrTODO is a placeholder for unimplemented paths.
	ErrTODO = Define("todo").Causes(CauseInternal).Actions(ActionAbort).Traced().Message("not implemented")
)
