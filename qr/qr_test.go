package qr

import (
	"bytes"
	"qrencode"
	"testing"
)

func TestVersion(t *testing.T) {
	badvers := 0
Version:
	for v := Version(1); v <= 40; v++ {
		c, err := qrencode.Encode(qrencode.Version(v), qrencode.L, qrencode.EightBit, "hi")
		if err != nil {
			t.Errorf("qrencode.Encode(%v, L, 8bit, hi): %v", v, err)
			continue
		}
		mask := (^c.Pixel[8][2]&1)<<2 | (c.Pixel[8][3]&1)<<1 | (^c.Pixel[8][4] & 1)
		p, err := NewPlan(v, L, Mask(mask))
		if err != nil {
			t.Errorf("NewPlan(%v, L, %d): %v", v, err, mask)
			continue
		}
		if len(p.Pixel) != len(c.Pixel) {
			t.Errorf("%v: NewPlan uses %dx%d, libqrencode uses %dx%d", v, len(p.Pixel), len(p.Pixel), len(c.Pixel), len(c.Pixel))
			continue
		}
		badpix := 0
	Pixel:
		for y, prow := range p.Pixel {
			for x, pix := range prow {
				cpix := c.Pixel[y][x]
				want := Pixel(0)
				switch {
				case cpix&qrencode.Finder != 0:
					want = Position.Pixel()
				case cpix&qrencode.Alignment != 0:
					want = Alignment.Pixel()
				case cpix&qrencode.Timing != 0:
					want = Timing.Pixel()
				case cpix&qrencode.Format != 0:
					want = Format.Pixel()
					want |= OffsetPixel(pix.Offset()) // sic
					want |= pix&Invert
				case cpix&qrencode.PVersion != 0:
					want = PVersion.Pixel()
				case cpix&qrencode.DataECC != 0:
					if pix.Role() == Check || pix.Role() == Extra {
						want = pix.Role().Pixel()
					} else {
						want = Data.Pixel()
					}
					want |= OffsetPixel(pix.Offset())
					want |= pix&Invert
					// KLUDGE
					if pix.Role() != Extra {
						want |= pix&Black
					}
				default:
					want = Unused.Pixel()
				}
				switch want.Role() {
				case Check, Data:
					//
				default:
					if cpix&qrencode.Black != 0 {
						want |= Black
					}
				}
				if pix != want {
					t.Errorf("%v/%v: Pixel[%d][%d] = %v, want %v %#x", v, mask, y, x, pix, want, cpix)
					if badpix++; badpix >= 10 {
						t.Errorf("stopping after %d bad pixels", badpix)
						break Pixel
					}
				}
			}
		}
		if badpix > 0 {
			if badvers++; badvers >= 5 {
				t.Errorf("stopping after %d bad versions", badvers)
				break Version
			}
		}
	}
}

func TestEncode(t *testing.T) {
	data := []byte{0x10, 0x20, 0x0c, 0x56, 0x61, 0x80, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11}
	check := []byte{0xa5, 0x24, 0xd4, 0xc1, 0xed, 0x36, 0xc7, 0x87, 0x2c, 0x55}
	const poly = 0x11d
	out := field.ECBytes(data, len(check))
	if !bytes.Equal(out, check) {
		t.Errorf("have %x want %x", out, check)
	}		
}
