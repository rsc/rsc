// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gf256 implements arithmetic over the Galois Field GF(256).
package gf256

import "strconv"

type Field struct {
	log [256]byte	// log[0] is unused
	exp [255]byte
}

func NewField(poly int) *Field {
	if poly < 0x100 || poly >= 0x200 {
		panic("gf256: invalid polynomial: " + strconv.Itoa(poly))
	}
	var f Field
	x := 1
	for i := range f.exp {
		if x == 1 && i != 0 {
			panic("gf256: reducible polynomial: " + strconv.Itoa(poly))
		}
		f.exp[i] = byte(x)
		f.log[x] = byte(i)
		x *= 2
		if x >= 0x100 {
			x ^= poly
		}
	}
	return &f
}

// Add returns the sum of x and y in the field.
func (f *Field) Add(x, y byte) byte {
	return x ^ y
}

// Exp returns the base 2 exponential of e in the field.
// If e < 0, Exp returns 0.
func (f *Field) Exp(e int) byte {
	if e < 0 {
		return 0
	}
	return f.exp[e%255]
}

// Log returns the base 2 logarithm of x in the field.
// If x == 0, Log returns -1.
func (f *Field) Log(x byte) int {
	if x == 0 {
		return -1
	}
	return int(f.log[x])
}

// Inv returns the multiplicative inverse of x in the field.
// If x == 0, Inv returns 0.
func (f *Field) Inv(x byte) byte {
	if x == 0 {
		return 0
	}
	return f.Exp(255 - f.Log(x))
}

// Mul returns the product of x and y in the field.
func (f *Field) Mul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	return f.Exp((f.Log(x) + f.Log(y)) % 255)
}

type Poly []byte

var Zero = Poly{}
var One = Poly{1}

func (z Poly) Norm() Poly {
	i := len(z)
	for i > 0 && z[i-1] == 0 {
		i--
	}
	return z[0:i]
}

func (x Poly) Add(y Poly) Poly {
	if len(x) < len(y) {
		x, y = y, x
	}
	z := make(Poly, len(x))
	for i := range y {
		z[i] = x[i] ^ y[i]
	}
	for i := len(y); i < len(x); i++ {
		z[i] = x[i]
	}
	return z.Norm()
}

func Mono(a byte, i int) Poly {
	p := make(Poly, i+1)
	p[i] = a
	return p
}

func (f *Field) MulPoly(x, y Poly) Poly {
	if len(x) == 0 || len(y) == 0 {
		return nil
	}
	z := make(Poly, len(x) + len(y) - 1)
	for i, xi := range x {
		if xi == 0 {
			continue
		}
		for j, yj := range y {
			z[i+j] = z[i+j] ^ f.Mul(xi, yj)
		}
	}
	return z
}

func (f *Field) DivPoly(x, y Poly) (q, r Poly) {
	y = y.Norm()
	if len(y) == 0 {
		panic("divide by zero")
	}

	r = x
	inv := f.Inv(y[len(y)-1])
	for len(r) >= len(y) {
		iq := Mono(f.Mul(r[len(r)-1], inv), len(r)-len(y))
		q = q.Add(iq)
		r = r.Add(f.MulPoly(iq, y))
	}
	return
}

func (p Poly) String() string {
	s := ""
	for i := len(p)-1; i >= 0; i-- {
		v := p[i]
		if v != 0 {
			if s != "" {
				s += " + "
			}
			if v != 1 {
				s += strconv.Itoa(int(v)) + " "
			}
			s += "x^" + strconv.Itoa(i)
		}
	}
	return s
}


func (f *Field) Gen(e int) Poly {
	p := Poly{1}
	for i := 0; i < e; i++ {
		p = f.MulPoly(p, Poly{f.Exp(i), 1})
	}
	return p
}

func (f *Field) ECBytes(data []byte, ecBytes int) []byte {
	if ecBytes == 0 {
		return nil
	}

	p := make(Poly, len(data))
	n := len(p)-1
	for i, v := range data {
		p[n-i] = v
	}
	p = p.Norm()
	p = f.MulPoly(p, Mono(1, ecBytes))
	
	_, r := f.DivPoly(p, f.Gen(ecBytes))
	ec := make([]byte, ecBytes)
	n = ecBytes-1
	for i, v := range r {
		ec[n-i] = v
	}
	return ec
}
