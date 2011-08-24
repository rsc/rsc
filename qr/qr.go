package qr

import (
	"fmt"
	"os"
	"strconv"
)

// A Version represents a QR version.
// The version specifies the size of the QR code:
// a QR code with version v has 4v+17 pixels on a side.
// Versions number from 1 to 40: the larger the version,
// the more information the code can store.
// A non-positive version means to select the
// version automatically.
type Version int

func (v Version) String() string {
	if v < 1 {
		return "auto"
	}
	return strconv.Itoa(int(v))
}

// A Mode represents a QR mode.
// The mode specifies the character set and, indirectly,
// the encoding.  The more precise the mode, the shorter
// the encoded data.
type Mode int

const (
	Numeric      Mode = 1
	Alphanumeric Mode = 2
	EightBit     Mode = 4
)

func (m Mode) String() string {
	switch m {
	case Numeric:
		return "numeric"
	case Alphanumeric:
		return "alpha"
	case EightBit:
		return "8bit"
	}
	return strconv.Itoa(int(m))
}

// A Pixel describes a single pixel in a QR code.
type Pixel uint32

const (
	Black Pixel = 1<<iota
	Invert
)

func (p Pixel) Offset() int {
	return int(p>>5)
}

func OffsetPixel(o int) Pixel {
	return Pixel(o<<5)
}

func (r PixelRole) Pixel() Pixel {
	return Pixel(r<<2)
}

func (p Pixel) Role() PixelRole {
	return PixelRole(p>>2)&7
}

func (p Pixel) String() string {
	s := p.Role().String()
	if p&Black != 0 {
		s += "+black"
	}
	if p&Invert != 0 {
		s += "+invert"
	}
	s += "+" + strconv.Itoa(p.Offset())
	return s
}

// A PixelRole describes the role of a QR pixel.
type PixelRole uint32
const (
	_ PixelRole = iota
	Position // position squares (large)
	Alignment  // alignment squares (small)
	Timing  // timing strip between position squares
	Format  // format metadata
	Data  // data bit
	Check  // error correction check bit
)

var roles = []string{
	"",
	"position",
	"alignment",
	"timing",
	"format",
	"data",
	"check",
}

func (r PixelRole) String() string {
	if Position <= r && r <= Check {
		return roles[r]
	}
	return strconv.Itoa(int(r))
}

// A Level represents a QR error correction level.
// From least to most tolerant of errors, they are L, M, Q, H.
type Level int

const (
	L Level = 0
	M
	Q
	H
)

func (l Level) String() string {
	if L <= l && l <= H {
		return "LMQH"[l:l+1]
	}
	return strconv.Itoa(int(l))
}

// A Mask describes a mask that is applied to the QR
// code to avoid QR artifacts being interpreted as
// alignment and timing patterns (such as the squares
// in the corners).
type Mask int

// TODO: fill in masks

// A Plan describes how to construct a QR code
// with a specific version, level, and mask.
type Plan struct {
	Version Version
	Level Level
	Mask Mask
	
	DataBytes int // number of data bytes
	CheckBytes int // number of error correcting (checksum) bytes
	Blocks  int // number of data blocks
	
	Pixel [][]Pixel  // pixel map
}

// NewPlan returns a Plan for a QR code with the given
// version, level, and mask.
func NewPlan(version Version, level Level, mask Mask) (*Plan, os.Error) {
	p, err := vplan(version)
	if err != nil {
		return nil, err
	}
	if err := lplan(level, p); err != nil {
		return nil, err
	}
	if err := mplan(mask, p); err != nil {
		return nil, err
	}
	return p, nil
}

// vplan creates a Plan for the given version.
func vplan(v Version) (*Plan, os.Error) {
	p := &Plan{Version: v}
	if v < 1 || v > 40 {
		return nil, fmt.Errorf("invalid QR version %d", int(v))
	}
	siz := 17 + int(v)*4
	m := make([][]Pixel, siz)
	pix := make([]Pixel, siz*siz)
	for i := range m {
		m[i], pix = pix[:siz], pix[siz:]
	}
	p.Pixel = m
	
	// Timing markers (overwritten by boxes).
	// TODO: are there more in higher versions?
	const ti = 6  // timing is in row/column 6 (counting from 0)
	for i := range m {
		p := Timing.Pixel()
		if i&1 == 0 {
			p |= Black
		}
		m[i][ti] = p
		m[ti][i] = p
	}
	
	// Position boxes.
	posBox(m, 0, 0)
	posBox(m, siz-7, 0)
	posBox(m, 0, siz-7)

	// Alignment boxes.
	// TODO.
	
	return p, nil
}

// lplan edits a version-only Plan to add information
// about the error correction levels.
func lplan(l Level, p *Plan) os.Error {
	p.Level = l
	// TODO: fill in info
	return nil
}

// mplan edits a version+level-only Plan to add the mask.
func mplan(m Mask, p *Plan) os.Error {
	p.Mask = m
	// TODO: invert pixels
	return nil
}

// posBox draws a position (large) box at upper left x, y.
func posBox(m [][]Pixel, x, y int) {
	pos := Position.Pixel()
	// box
	for dy := 0; dy < 7; dy++ {
		for dx := 0; dx < 7; dx++ {
			p := pos
			if dx == 0 || dx == 6 || dy == 0 || dy == 6 || 2 <= dx && dx <= 4 && 2 <= dy && dy <= 4 {
				p |= Black
			}
			m[y+dy][x+dx] = p
		}
	}
	// white border
	for dy := -1; dy < 8; dy++ {
		if 0 <= y+dy && y+dy < len(m) {
			if x > 0 {
				m[y+dy][x-1] = pos
			}
			if x+7 < len(m) {
				m[y+dy][x+7] = pos
			}
		}
	}
	for dx := -1; dx < 8; dx++ {
		if 0 <= x+dx && x+dx < len(m) {
			if y > 0 {
				m[y-1][x+dx] = pos
			}
			if y+7 < len(m) {
				m[y+7][x+dx] = pos
			}
		}
	}
}

// alignBox draw an alignment (small) box at upper left x, y.
func alignBox(m [][]Pixel, x, y int) {
	// box
	align := Alignment.Pixel()
	for dy := 0; dy < 5; dy++ {
		for dx := 0; dx < 5; dx++ {
			p := align
			if dx == 0 || dx == 4 || dy == 0 || dy == 4 || dx == 2 && dy == 2 {
				p |= Black
			}
			m[y+dy][x+dx] = p
		}
	}
}
