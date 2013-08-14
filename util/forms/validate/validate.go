package validate

import (
	"github.com/kidstuff/toys/util/errs"
	"reflect"
)

type Invalid map[string]error

func (i Invalid) Error() string {
	return ""
}

func Valid(i interface{}) (*Invalid, error) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errs.New("validate: function must recieve a struct or a pointer to struct")
	}
}

func check(v reflect.Value, m Invalid) error {
	if v.Kind() == reflect.Ptr {

	}
}
