package forms

import (
	"testing"
)

type baz struct {
	U int
	V int
}

type bar struct {
	X int
	Y int
	Z baz
}

type foz struct {
	I int
	O string
}

type foo struct {
	A int
	B bar
	C baz
	D []foz
}

func TestForms(t *testing.T) {
	Prepare(&foo{})
	printCache()
}
