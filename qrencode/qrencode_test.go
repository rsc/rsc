// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qrencode

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	c, err := Encode(2, L, Alphanumeric, "HELLO WORLD")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	for _, pix := range c.Pixel {
		fmt.Print("    ")
		for _, p := range pix {
			if p&Black != 0 {
				fmt.Print("")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	
	f, err := os.Create("x.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, c); err != nil {
		t.Fatal(err)
	}
	f.Close()
}
