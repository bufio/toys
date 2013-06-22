package forms

import (
	"fmt"
	"reflect"
	"strings"
)

func fullname(t reflect.Type) string {
	//make it look like url to document
	return t.PkgPath() + "/#" + t.Name()
}

func taglist(tag reflect.StructTag) []string {
	return strings.Split(tag.Get("forms"), ",")
}

func isBasic(k reflect.Kind) bool {
	// String and Numeric are basic
	return (1 <= k && k <= 16) || k == 24
}

func printLog(sInfo *StructInfo) {
	for name, _ := range sInfo.Data {
		fmt.Println(name)
	}
}

func printCache() {
	for name, sInfo := range cachedStruct {
		fmt.Println(name)
		printLog(sInfo)
	}
}
