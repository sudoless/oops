package oops

// Matcher is a predicate that tests whether an Error matches some criteria.
type Matcher = func(*Error) bool

// Match tests whether an error matches the given Matcher.
func Match(err error, m Matcher) bool {
	if err == nil {
		return false
	}
	v, ok := err.(*Error) //nolint:errorlint // Match operates on direct *Error only by design
	if !ok {
		return false
	}
	return m(v)
}

// All returns a Matcher that requires all sub-matchers to match.
func All(matchers ...Matcher) Matcher {
	return func(err *Error) bool {
		for _, m := range matchers {
			if !m(err) {
				return false
			}
		}
		return true
	}
}

// Any returns a Matcher that requires at least one sub-matcher to match.
func Any(matchers ...Matcher) Matcher {
	return func(err *Error) bool {
		for _, m := range matchers {
			if m(err) {
				return true
			}
		}
		return false
	}
}

// Not returns a Matcher that inverts the given matcher.
func Not(m Matcher) Matcher {
	return func(err *Error) bool { return !m(err) }
}

// ByCause returns a Matcher that checks for a specific cause tag.
func ByCause(cause string) Matcher {
	return func(err *Error) bool { return err.HasCause(cause) }
}

// ByAction returns a Matcher that checks for a specific action tag.
func ByAction(action string) Matcher {
	return func(err *Error) bool { return err.HasAction(action) }
}

// ByCode returns a Matcher that checks for a specific error code.
func ByCode(code string) Matcher {
	return func(err *Error) bool { return err.Code() == code }
}

// ByDefinition returns a Matcher that checks definition identity (including inherits).
func ByDefinition(def *ErrorDefinition) Matcher {
	return func(err *Error) bool { return err.Is(def) }
}
