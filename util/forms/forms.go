package forms

import (
	"container/list"
	"errors"
	"reflect"
)

var cachedStruct = make(map[string]*StructInfo)

type StructInfo struct {
	Data map[string]FieldInfo
}

type FieldInfo struct {
	isBasic bool
	isSlice bool
	kind    reflect.Kind
	sInfo   *StructInfo
}

type processItem struct {
	t     reflect.StructField
	sName string
}

func Prepare(i interface{}) (*StructInfo, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return parse(t.Elem()), nil
	}
	return nil, errors.New("forms: function must recieve a pointer")
}

func parse(t reflect.Type) *StructInfo {
	sName := fullname(t)
	if sInfo, ok := cachedStruct[sName]; ok {
		return sInfo
	}

	sInfo := &StructInfo{}
	sInfo.Data = make(map[string]FieldInfo)
	cachedStruct[sName] = sInfo

	lst := list.New()
	pushToList(lst, t, sName)

	for lst.Len() != 0 {
		lstItem := lst.Front()
		lst.Remove(lstItem)
		item, ok := lstItem.Value.(processItem)
		if !ok {
			return cachedStruct[sName]
		}
		tags := taglist(item.t.Tag)
		if tags[0] == "-" {
			continue
		}

		fInfo := FieldInfo{}
		kind := item.t.Type.Kind()
		if isBasic(kind) {
			fInfo.isBasic = true
			fInfo.kind = kind
		} else if kind == reflect.Struct {
			sName := fullname(item.t.Type)
			if _, ok := cachedStruct[sName]; !ok {
				cachedStruct[sName] = &StructInfo{}
				cachedStruct[sName].Data = make(map[string]FieldInfo)
				pushToList(lst, item.t.Type, sName)
			}
			fInfo.sInfo = cachedStruct[sName]
		} else if kind == reflect.Slice || kind == reflect.Array {
			sName := fullname(item.t.Type.Elem())
			if _, ok := cachedStruct[sName]; !ok {
				cachedStruct[sName] = &StructInfo{}
				cachedStruct[sName].Data = make(map[string]FieldInfo)
				pushToList(lst, item.t.Type.Elem(), sName)
			}
			fInfo.sInfo = cachedStruct[sName]
		}

		cachedStruct[item.sName].Data[item.t.Name] = fInfo
	}
	return cachedStruct[sName]
}

func pushToList(lst *list.List, t reflect.Type, sName string) {
	for i := 0; i < t.NumField(); i++ {
		lst.PushFront(processItem{t.Field(i), sName})
	}
}

func Decode(dst interface{}, src map[string][]string) error {
	// sInfo, err := Prepare(dst)
	// if err != nil {
	// 	return err
	// }

	// for path, vals := range src {

	// }
	return nil
}
