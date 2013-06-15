package forms

type FormError struct {
}

func (e *FormError) Error() string {
	return ""
}
