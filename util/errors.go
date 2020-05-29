package util

import (
	"github.com/pkg/errors"
)

// ConcatErrors handles wrapping a list of errors into a single error.
func ConcatErrors(errs ...error) error {

	var err error

	for _, e := range errs {

		if e != nil && err != nil {
			err = e
		} else if e != nil {
			err = errors.Wrap(err, e.Error())
		}

	}

	return err

}
