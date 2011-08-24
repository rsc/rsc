package qr

import (
	"qrencode"
	"testing"
)

func TestVersion(t *testing.T) {
	badvers := 0
Version:
	for v := Version(1); v <= 40; v++ {
		p, err := NewPlan(v, L, 0)
		if err != nil {
			t.Errorf("NewPlan(%v, L, 0): %v", v, err)
			continue
		}
		c, err := qrencode.Encode(qrencode.Version(v), qrencode.L, qrencode.EightBit, "hi")
		if err != nil {
			t.Errorf("qrencode.Encode(%v, L, 8bit, hi): %v", v, err)
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
				if cpix&qrencode.Finder != 0 {
					want = Position.Pixel()
				} else if cpix&qrencode.Alignment != 0 {
					want = Alignment.Pixel()
				} else if cpix&qrencode.Timing != 0 {
					want = Timing.Pixel()
				} else if cpix&qrencode.Format != 0 {
					want = Format.Pixel()
					want |= OffsetPixel(pix.Offset()) // sic
				}
				if want != 0 && want.Role() != Format && cpix&qrencode.Black != 0 {
					want |= Black
				}
				if want != 0 && pix != want {
					t.Errorf("%v: Pixel[%d][%d] = %v, want %v", v, y, x, pix, want)
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
