package oops

import (
	"fmt"
	"strconv"
)

type multipleErrors struct {
	errs []*Error
}

func (m *multipleErrors) String() string {
	return "[multiple errors] " + strconv.Itoa(len(m.errs))
}

func (m *multipleErrors) Error() string {
	return fmt.Sprintf("multiple errors %d", len(m.errs))
}

func (m *multipleErrors) Unwrap() []error {
	errs := make([]error, len(m.errs))
	for idx := range m.errs {
		errs[idx] = m.errs[idx]
	}
	return errs
}

// Join combines all the given errors into a multipleErrors error. If any of the given errors is ErrMultiple, then it
// will be flattened and all the errors will be added to the multipleErrors error also applying the multiErr explanation.
func Join(errs ...error) *Error {
	nonNilErrs := make([]*Error, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}

		v, ok := err.(*Error)
		if !ok {
			v = ErrUnexpected.Wrap(err, "")
		}

		if v.source == ErrMultiple && v.parent != nil {
			nonNilErrs = append(nonNilErrs, v.parent.(*multipleErrors).errs...)
		} else {
			nonNilErrs = append(nonNilErrs, v)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	}

	if len(nonNilErrs) == 1 {
		return nonNilErrs[0]
	}

	return ErrMultiple.Wrap(&multipleErrors{errs: nonNilErrs}, "")
}
