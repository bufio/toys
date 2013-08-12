// +build debug

package errs

import (
	"fmt"
	"runtime"
)

// errorInfo is a trivial implementation of error.
type errorInfo struct {
	text string
	file string
	line int
}

func (e *errorInfo) Error() string {
	return fmt.Sprintf("%s:%d %s", e.file, e.line, e.text)
}

func newErrorInfo(olderr error, text string) *errorInfo {
	err := &errorInfo{}

	if olderr != nil {
		err.text = olderr.Error() + "\n"
	}

	err.text += text
	_, file, line, ok := runtime.Caller(2)
	if ok {
		err.file = file
		err.line = line
	}
	return err
}
