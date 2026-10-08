package oops

// Cause is a semantic tag describing why an error occurred.
type Cause = string

const (
	CauseInternal    Cause = "internal"
	CauseNotFound    Cause = "not_found"
	CauseAuth        Cause = "auth"
	CauseForbidden   Cause = "forbidden"
	CauseConflict    Cause = "conflict"
	CauseRateLimit   Cause = "rate_limit"
	CauseTimeout     Cause = "timeout"
	CauseBadRequest  Cause = "bad_request"
	CauseUnavailable Cause = "unavailable"
	CauseBadGateway  Cause = "bad_gateway"
	CauseExpired     Cause = "expired"
	CauseIO          Cause = "io"
	CauseValidation  Cause = "validation"
)
