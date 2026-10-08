package oops

// Formatter controls how an Error is rendered as a string: it replaces the
// default code[: message][; explanation] returned by Error(). Under %+v it
// renders only the first line of the error's entry in the tree; the attribute
// and child lines are always rendered by oops (see (*Error).Format).
type Formatter = func(*Error) string

func defaultFormatter(err *Error) string {
	if err == nil {
		return "oops.Error(nil)"
	}

	explanation := err.Explanation()
	if explanation != "" {
		if err.def.message != "" {
			return err.def.code + ": " + err.def.message + "; " + explanation
		}

		return err.def.code + ": " + explanation
	}

	if err.def.message != "" {
		return err.def.code + ": " + err.def.message
	}

	return err.def.code
}
