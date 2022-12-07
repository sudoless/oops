package oops

import "fmt"

type multipleErrors struct {
	errs []error
}

func (m multipleErrors) Error() string {
	return fmt.Sprintf("multiple errors %d", len(m.errs))
}

func (m multipleErrors) Unwrap() []error {
	return m.errs
}

func Multi(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return &multipleErrors{errs: errs}
}
