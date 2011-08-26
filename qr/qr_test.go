package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"qrencode"
	"testing"
)

func test(t *testing.T, v Version, l Level, text ...Encoding) bool {
	s := ""
	ty := qrencode.EightBit
	switch x := text[0].(type) {
	case String:
		s = string(x)
	case Alpha:
		s = string(x)
		ty = qrencode.Alphanumeric
	case Num:
		s = string(x)
		ty = qrencode.Numeric
	}
	key, err := qrencode.Encode(qrencode.Version(v), qrencode.Level(l), ty, s)
	if err != nil {
		t.Errorf("qrencode.Encode(%v, %v, %d, %#q): %v", v, l, ty, s, err)
		return false
	}
	mask := (^key.Pixel[8][2]&1)<<2 | (key.Pixel[8][3]&1)<<1 | (^key.Pixel[8][4] & 1)
	p, err := NewPlan(v, l, Mask(mask))
	if err != nil {
		t.Errorf("NewPlan(%v, L, %d): %v", v, err, mask)
		return false
	}
	if len(p.Pixel) != len(key.Pixel) {
		t.Errorf("%v: NewPlan uses %dx%d, libqrencode uses %dx%d", v, len(p.Pixel), len(p.Pixel), len(key.Pixel), len(key.Pixel))
		return false
	}
	c, err := p.Encode(text...)
	if err != nil {
		t.Errorf("Encode: %v", err)
		return false
	}
	badpix := 0
Pixel:
	for y, prow := range c.Pixel {
		for x, pix := range prow {
			keypix := key.Pixel[y][x]
			want := Pixel(0)
			switch {
			case keypix&qrencode.Finder != 0:
				want = Position.Pixel()
			case keypix&qrencode.Alignment != 0:
				want = Alignment.Pixel()
			case keypix&qrencode.Timing != 0:
				want = Timing.Pixel()
			case keypix&qrencode.Format != 0:
				want = Format.Pixel()
				want |= OffsetPixel(pix.Offset()) // sic
				want |= pix & Invert
			case keypix&qrencode.PVersion != 0:
				want = PVersion.Pixel()
			case keypix&qrencode.DataECC != 0:
				if pix.Role() == Check || pix.Role() == Extra {
					want = pix.Role().Pixel()
				} else {
					want = Data.Pixel()
				}
				want |= OffsetPixel(pix.Offset())
				want |= pix & Invert
			default:
				want = Unused.Pixel()
			}
			if keypix&qrencode.Black != 0 {
				want |= Black
			}
			if pix != want {
				t.Errorf("%v/%v: Pixel[%d][%d] = %v, want %v %#x", v, mask, y, x, pix, want, keypix)
				if badpix++; badpix >= 100 {
					t.Errorf("stopping after %d bad pixels", badpix)
					break Pixel
				}
			}
		}
	}
	if false {
		write(key, "v%d-key.png", int(v))
		write(c, "v%d-out.png", int(v))
		for y, row := range c.Pixel {
			for x, pix := range row {
				if pix&Invert != 0 {
					row[x] ^= Black
					key.Pixel[y][x] ^= qrencode.Black
				}
			}
		}
		write(key, "v%du-key.png", int(v))
		write(c, "v%du-out.png", int(v))
	}
	return badpix == 0
}

var input = []Encoding{
	String("hello"),
	Num("1"),
	Num("12"),
	Num("123"),
	Alpha("AB"),
	Alpha("ABC"),
}

func TestVersion(t *testing.T) {
	badvers := 0
Version:
	for v := Version(1); v <= 40; v++ {
		for l := L; l <= H; l++ {
			for _, in := range input {
				if !test(t, v, l, in) {
					if badvers++; badvers >= 10 {
						t.Errorf("stopping after %d bad versions", badvers)
						break Version
					}
				}
			}
		}
	}
}

func TestEncode(t *testing.T) {
	data := []byte{0x10, 0x20, 0x0c, 0x56, 0x61, 0x80, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11}
	check := []byte{0xa5, 0x24, 0xd4, 0xc1, 0xed, 0x36, 0xc7, 0x87, 0x2c, 0x55}
	out := field.ECBytes(data, len(check))
	if !bytes.Equal(out, check) {
		t.Errorf("have %x want %x", out, check)
	}
}

func write(m image.Image, format string, args ...interface{}) {
	f, err := os.Create(fmt.Sprintf(format, args...))
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, m); err != nil {
		panic(err)
	}
	f.Close()
}
