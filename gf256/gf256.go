// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gf256 implements arithmetic over the Galois Field GF(256).
package gf256

import "strconv"

type Field struct {
	log [256]byte // log[0] is unused
	exp [510]byte
}

func NewField(poly int) *Field {
	if poly < 0x100 || poly >= 0x200 {
		panic("gf256: invalid polynomial: " + strconv.Itoa(poly))
	}
	var f Field
	x := 1
	for i := 0; i < 255; i++ {
		if x == 1 && i != 0 {
			panic("gf256: reducible polynomial: " + strconv.Itoa(poly))
		}
		f.exp[i] = byte(x)
		f.exp[i+255] = byte(x)
		f.log[x] = byte(i)
		x *= 2
		if x >= 0x100 {
			x ^= poly
		}
	}
	f.log[0] = 255
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
	return f.exp[255-f.log[x]]
}

// Mul returns the product of x and y in the field.
func (f *Field) Mul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	return f.exp[int(f.log[x])+int(f.log[y])]
}

// An RSEncoder implements Reed-Solomon encoding
// over a given field using a given number of error correction bytes.
type RSEncoder struct {
	f    *Field
	c    int
	lgen []byte
	p    []byte
}

func (f *Field) lgen(e int) []byte {
	// p = 1
	p := make([]byte, e+1)
	p[e] = 1

	for i := 0; i < e; i++ {
		// p *= (x + Exp(i))
		// p[j] = p[j]*Exp(i) + p[j+1].
		c := f.Exp(i)
		for j := 0; j < e; j++ {
			p[j] = f.Mul(p[j], c) ^ p[j+1]
		}
		p[e] = f.Mul(p[e], c)
	}

	// replace p with log p.
	for i, c := range p {
		if c == 0 {
			p[i] = 255
		} else {
			p[i] = byte(f.Log(c))
		}
	}
	return p
}

// NewRSEncoder returns a new Reed-Solomon encoder
// over the given field and number of error correction bytes.
func NewRSEncoder(f *Field, c int) *RSEncoder {
	return &RSEncoder{f: f, c: c, lgen: f.lgen(c)}
}

// ECC writes to check the error correcting code bytes
// for data using the given Reed-Solomon parameters.
func (rs *RSEncoder) ECC(data []byte, check []byte) {
	if len(check) < rs.c {
		panic("gf256.RSEncoder: invalid check byte length")
	}
	if rs.c == 0 {
		return
	}

	// The check bytes are the remainder after dividing
	// data padded with c zeros by the generator polynomial.  

	// p = data padded with c zeros.
	var p []byte
	n := len(data) + rs.c
	if len(rs.p) >= n {
		p = rs.p
	} else {
		p = make([]byte, n)
	}
	copy(p, data)
	for i := len(data); i < len(p); i++ {
		p[i] = 0
	}

	// Divide p by gen, leaving the remainder in p[len(data):].
	// p[0] is the most significant term in p, and
	// gen[0] is the most significant term in the generator,
	// which is always 1.
	// To avoid repeated work, we store various values as
	// lv, not v, where lv = log[v].
	f := rs.f
	lgen := rs.lgen[1:]
	for i := 0; i < len(data); i++ {
		c := p[i]
		if c == 0 {
			continue
		}
		q := p[i+1:]
		exp := f.exp[f.log[c]:]
		for j, lg := range lgen {
			if lg != 255 { // lgen uses 255 for log 0
				q[j] ^= exp[lg]
			}
		}
	}
	copy(check, p[len(data):])
	rs.p = p
}
