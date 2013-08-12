// +build !debug

package errs

// errorInfo is a trivial implementation of error.
type errorInfo struct {
	text string
}

func (e *errorInfo) Error() string {
	return e.text
}

func newErrorInfo(olderr error, text string) *errorInfo {
	err := &errorInfo{}

	if olderr != nil {
		err.text = olderr.Error() + "\n"
	}

	err.text += text

	return err
}
