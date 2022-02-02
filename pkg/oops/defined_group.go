package oops

// Group will take the existing *errorDefined and use that as a base for generating new *errorDefined.
func (e *errorDefined) Group() *errorGroup {
	return &errorGroup{
		base: e,
	}
}

type errorGroup struct {
	base   *errorDefined
	prefix string
}

// PrefixCode will add the given prefix to all calls to errorGroup.Code.
func (e *errorGroup) PrefixCode(prefix string) *errorGroup {
	e.prefix = prefix
	return e
}

// Code will take the existing *errorDefined and use that as a base for generating new *errorDefined with the given
// code. The returned *errorDefined can still be changed using the chained builders.
func (e *errorGroup) Code(code string) *errorDefined {
	defined := &errorDefined{}
	*defined = *e.base

	if e.prefix != "" {
		defined.code = e.prefix + code
	} else {
		defined.code = code
	}

	return defined
}
