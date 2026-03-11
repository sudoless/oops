package oops

// Formatter controls how an Error is rendered as a string.
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
