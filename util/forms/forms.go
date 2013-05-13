package forms

import (
	"reflect"
)

var cacheStruct = map[string]*StructData{}

type StructData struct {
	Feilds map[string]*StructData
	// converter func(string) interface{}
	// validator func(string) error
	IsBasic bool
}

func Prepare(i interface{}) (*StructData, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Struct {
		parse(t)
	}
	return nil, nil
}

func parse(t reflect.Type) *StructData {
	if v, ok := cacheStruct[t.PkgPath()+"."+t.Name()]; ok {
		return v
	}

	sdata := &StructData{}

	if (1 <= t.Kind() && t.Kind() <= 16) || t.Kind() == 24 {
		sdata.IsBasic = true
	} else {
		sdata.IsBasic = false
		numField := 0
		if t.Kind() == reflect.Struct {
			numField = t.NumField()
		} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			t = t.Elem()
			numField = t.NumField()
			//TODO check the cache one more time
		}
		if numField > 0 {
			sdata.Feilds = make(map[string]*StructData)
			cacheStruct[t.PkgPath()+"."+t.Name()] = sdata
			for i := 0; i < numField; i++ {
				f := t.Field(i)
				if len(f.PkgPath) == 0 { // exported field
					sdata.Feilds[f.Name] = parse(f.Type)
				}
			}
		}
	}
	return sdata
}

func Cache() map[string]*StructData {
	return cacheStruct
}
