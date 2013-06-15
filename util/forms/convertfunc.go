package forms

import (
	"reflect"
	"strconv"
)

type ConvertFunc func(string) reflect.Value

var basicTypeConvert = map[reflect.Kind]ConvertFunc{
	reflect.String: func(s string) reflect.Value {
		return reflect.ValueOf(s)
	},
	reflect.Int: func(s string) reflect.Value {
		i, err := strconv.Atoi(s)
		if err != nil {
			return reflect.ValueOf(0)
		}
		return reflect.ValueOf(i)
	},
}
