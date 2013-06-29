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
	B *bar
	C []foz
	D []string
}

/*
A

B.X
B.Y
B.Z.U
B.Z.V

C,0.I
C,0.O
C,1.I
C,1.O

D
*/

func TestForms(t *testing.T) {
	Prepare(&foo{})
	printCache()
}
