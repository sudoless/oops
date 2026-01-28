package oops

func (defined *errorDefined) Set(key string, value any) *errorDefined {
	if defined.props == nil {
		defined.props = make(map[string]any, 4)
	}

	defined.props[key] = value
	return defined
}

func (defined *errorDefined) Trace() *errorDefined {
	defined.traced = true
	return defined
}

func (defined *errorDefined) Formatter(formatter Formatter) *errorDefined {
	defined.formatter = formatter
	return defined
}
