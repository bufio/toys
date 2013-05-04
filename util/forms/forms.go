package forms

import (
	"reflect"
)

type StructData struct {
	Feilds    map[string]StructData
	converter func(string) reflect.Value
	validator func(string) error
}
