package oops

// Group will take the existing *ErrorDefined and use that as a base for generating new *ErrorDefined.
func (e *ErrorDefined) Group() *errorGroup {
	return &errorGroup{
		base: e,
	}
}

type errorGroup struct {
	base   *ErrorDefined
	prefix string
}

// PrefixCode will add the given prefix to all calls to errorGroup.Code.
func (e *errorGroup) PrefixCode(prefix string) *errorGroup {
	e.prefix = prefix
	return e
}

// Code will take the existing *ErrorDefined and use that as a base for generating new *ErrorDefined with the given
// code. The returned *ErrorDefined can still be changed using the chained builders.
func (e *errorGroup) Code(code string) *ErrorDefined {
	defined := &ErrorDefined{}
	*defined = *e.base

	if e.prefix != "" {
		defined.code = e.prefix + code
	} else {
		defined.code = code
	}

	return defined
}
