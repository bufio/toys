package forms

import (
	"container/list"
	"errors"
	"reflect"
)

var cachedStruct = make(map[string]*StructInfo)

type StructInfo struct {
	Data map[string]FieldInfo
	val  *reflect.Value
}

type FieldInfo struct {
	isBasic bool
	val     *reflect.Value
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
		fInfo.kind = kind
		if isBasic(kind) {
			fInfo.isBasic = true
		} else if kind == reflect.Struct {
			sName := fullname(item.t.Type)
			if _, ok := cachedStruct[sName]; !ok {
				cachedStruct[sName] = &StructInfo{}
				cachedStruct[sName].Data = make(map[string]FieldInfo)
				sType := reflect.Zero(item.t.Type)
				cachedStruct[sName].val = &sType
				pushToList(lst, item.t.Type, sName)
			}
			fInfo.sInfo = cachedStruct[sName]
		} else if elem := item.t.Type.Elem(); kind == reflect.Slice || kind == reflect.Array {
			var sName string
			var sType reflect.Value
			if elem.Kind() == reflect.Struct {
				sName = fullname(elem)
				sType = reflect.Zero(elem)
			} else if elem.Kind() == reflect.Ptr && elem.Elem().Kind() == reflect.Struct {
				sName = fullname(elem.Elem())
				sType = reflect.Zero(elem.Elem())
				elem = elem.Elem()
			}
			if elem.Kind() == reflect.Struct {
				if _, ok := cachedStruct[sName]; !ok {
					cachedStruct[sName] = &StructInfo{}
					cachedStruct[sName].Data = make(map[string]FieldInfo)
					cachedStruct[sName].val = &sType
					pushToList(lst, elem, sName)
				}
				fInfo.sInfo = cachedStruct[sName]
			}
		} else if kind == reflect.Ptr && elem.Kind() == reflect.Struct {
			sName := fullname(elem)
			if _, ok := cachedStruct[sName]; !ok {
				cachedStruct[sName] = &StructInfo{}
				cachedStruct[sName].Data = make(map[string]FieldInfo)
				sType := reflect.New(elem)
				cachedStruct[sName].val = &sType
				pushToList(lst, elem, sName)
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
