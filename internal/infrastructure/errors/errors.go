package errors

import (
	"errors"
	"fmt"
	"runtime"
)

func New(msg string) error {
	return errors.New(msg)
}

func Wrap(err error, msg string) error {
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)

	return fmt.Errorf("%s(): %s: %w", function.Name(), msg, err)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
