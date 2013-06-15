package forms

import (
	"errors"
	"reflect"
	"strings"
)

var cacheStruct = map[string]*StructData{}

type StructData struct {
	Fields map[string]*StructData
	// converter func(string) interface{}
	// validator func(string) error
	IsBasic bool
	Kind    reflect.Kind
}

func Prepare(i interface{}) (*StructData, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return parse(t.Elem()), nil
	}
	return nil, errors.New("forms: function must recieve a pointer")
}

func parse(t reflect.Type) *StructData {
	if v, ok := cacheStruct[t.PkgPath()+"."+t.Name()]; ok {
		return v
	}

	sdata := &StructData{}
	sdata.Kind = t.Kind()
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
			if v, ok := cacheStruct[t.PkgPath()+"."+t.Name()]; ok {
				return v
			}
		}
		if numField > 0 {
			sdata.Fields = make(map[string]*StructData)
			cacheStruct[t.PkgPath()+"."+t.Name()] = sdata
			for i := 0; i < numField; i++ {
				f := t.Field(i)
				if len(f.PkgPath) == 0 { // exported field
					sdata.Fields[f.Name] = parse(f.Type)
				}
			}
		}
	}
	return sdata
}

func Fill(src map[string][]string, dst interface{}) error {
	//TODO: alot
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.Elem().Kind() != reflect.Struct {
		return errors.New("forms: function must recieve a pointer")
	}
	dstVal = dstVal.Elem()
	dstData := parse(dstVal.Type())
	for path, input := range src {
		fields := strings.Split(path, ".")
		sdata := dstData
		fieldVal := dstVal
		var i int
		for i = 0; i < len(fields); i++ {
			var ok bool
			sdata, ok = sdata.Fields[fields[i]]
			if !ok {
				break
			}
			fieldVal = fieldVal.FieldByName(fields[i])
		}
		if sdata == nil || i != len(fields) {
			println("filed", sdata, i, len(fields))
			continue
		}
		if sdata.IsBasic {
			fieldVal.Set(basicTypeConvert[sdata.Kind](input[0]))
		}
	}
	return nil
}

func GenCode(v interface{}, ferr *FormError) (string, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		//sdata := parse(t.Elem())
	}
	return "", errors.New("forms: function must recieve a pointer")
}

func SetConverter(path string, f ConvertFunc) {

}

func Cache() map[string]*StructData {
	return cacheStruct
}
