package errs

import (
	"fmt"
)

// New returns an error that formats as the given text.
func New(text string) error {
	return newErrorInfo(nil, text)
}

func Newf(format string, a ...interface{}) error {
	return newErrorInfo(nil, fmt.Sprintf(format, a))
}

func Err(err error, text string) error {
	return newErrorInfo(err, text)
}

func Errf(err error, format string, a ...interface{}) error {
	return newErrorInfo(err, fmt.Sprintf(format, a))
}
