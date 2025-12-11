package oops

func defaultFormatter(err Error) string {
	if err == nil {
		return "oops.Error(nil)"
	}

	explanation := err.Explanation()
	if explanation == "" {
		return "oops.Error"
	}

	return explanation
}

func Define(props ...any) *errorDefined {
	defined := &errorDefined{
		formatter: defaultFormatter,
	}

	if len(props) > 0 {
		if len(props)%2 != 0 {
			panic("oops: Define requires an even number of arguments")
		}

		defined.props = make(map[string]any, len(props)/2)
	}

	for idx := 0; idx < len(props); idx += 2 {
		defined.props[props[idx].(string)] = props[idx+1]
	}

	return defined
}
