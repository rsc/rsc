package gf256

import "testing"

var f = NewField(0x11d)	// x^8 + x^4 + x^3 + x^2 + 1

func TestBasic(t *testing.T) {
	if f.Exp(0) != 1 || f.Exp(1) != 2 || f.Exp(255) != 1 {
		panic("bad Exp")
	}
}
