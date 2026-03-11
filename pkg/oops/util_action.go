package oops

// Action is a semantic tag describing what the caller should do about it.
type Action = string

const (
	ActionRetry Action = "retry"
	ActionAbort Action = "abort"
	ActionFatal Action = "fatal"
	ActionFix   Action = "fix"
	ActionWait  Action = "wait"
	ActionAuth  Action = "auth"
	ActionSkip  Action = "skip"
)
