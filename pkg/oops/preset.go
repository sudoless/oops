package oops

var (
	ErrTODO = &errorDefined{
		traced: true,
		formatter: func(err Error) string {
			explain := err.Explanation()
			if explain != "" {
				return "TODO: " + explain
			}

			return "TODO"
		},
	}

	ErrUncaught = &errorDefined{
		traced: true,
		formatter: func(err Error) string {
			explain := err.Explanation()
			if explain != "" {
				return "uncaught unwrapped: " + explain
			}

			return "uncaught unwrapped"
		},
	}

	NilErr = Error((*errorImpl)(nil))
)
