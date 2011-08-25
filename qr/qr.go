package qr

import (
	"fmt"
	"gf256"
	"os"
	"strconv"
)

// field is the field for QR error correction.
var field = gf256.NewField(0x11d)

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
	Black Pixel = 1 << iota
	Invert
)

func (p Pixel) Offset() int {
	return int(p >> 6)
}

func OffsetPixel(o int) Pixel {
	return Pixel(o << 6)
}

func (r PixelRole) Pixel() Pixel {
	return Pixel(r << 2)
}

func (p Pixel) Role() PixelRole {
	return PixelRole(p>>2) & 15
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
	_         PixelRole = iota
	Position            // position squares (large)
	Alignment           // alignment squares (small)
	Timing              // timing strip between position squares
	Format              // format metadata
	PVersion   // version pattern
	Unused   // unused pixel
	Data                // data bit
	Check               // error correction check bit
	Extra
)

var roles = []string{
	"",
	"position",
	"alignment",
	"timing",
	"format",
	"pversion",
	"unused",
	"data",
	"check",
	"extra",
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
		return "LMQH"[l : l+1]
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
	Level   Level
	Mask    Mask

	DataBytes  int // number of data bytes
	CheckBytes int // number of error correcting (checksum) bytes
	Blocks     int // number of data blocks

	Pixel [][]Pixel // pixel map
}

// NewPlan returns a Plan for a QR code with the given
// version, level, and mask.
func NewPlan(version Version, level Level, mask Mask) (*Plan, os.Error) {
	p, err := vplan(version)
	if err != nil {
		return nil, err
	}
	if err := fplan(level, mask, p); err != nil {
		return nil, err
	}
	if err := lplan(version, level, p); err != nil {
		return nil, err
	}
	if err := mplan(mask, p); err != nil {
		return nil, err
	}
	return p, nil
}

// A version describes metadata associated with a version.
type version struct {
	apos    int
	astride int
	bytes int
	pattern int
	level [4]level
}

type level struct {
	nblock int
	check int
}

var vtab = []version{
	{},
	{100, 100, 26, 0x0, [4]level{{1, 7}, {1, 10}, {1, 13}, {1, 17}}},  // 1
	{16, 100, 44, 0x0, [4]level{{1, 10}, {1, 16}, {1, 22}, {1, 28}}},  // 2
	{20, 100, 70, 0x0, [4]level{{1, 15}, {1, 26}, {2, 18}, {2, 22}}},  // 3
	{24, 100, 100, 0x0, [4]level{{1, 20}, {2, 18}, {2, 26}, {4, 16}}},  // 4
	{28, 100, 134, 0x0, [4]level{{1, 26}, {2, 24}, {4, 18}, {4, 22}}},  // 5
	{32, 100, 172, 0x0, [4]level{{2, 18}, {4, 16}, {4, 24}, {4, 28}}},  // 6
	{20, 16, 196, 0x7c94, [4]level{{2, 20}, {4, 18}, {6, 18}, {5, 26}}},  // 7
	{22, 18, 242, 0x85bc, [4]level{{2, 24}, {4, 22}, {6, 22}, {6, 26}}},  // 8
	{24, 20, 292, 0x9a99, [4]level{{2, 30}, {5, 22}, {8, 20}, {8, 24}}},  // 9
	{26, 22, 346, 0xa4d3, [4]level{{4, 18}, {5, 26}, {8, 24}, {8, 28}}},  // 10
	{28, 24, 404, 0xbbf6, [4]level{{4, 20}, {5, 30}, {8, 28}, {11, 24}}},  // 11
	{30, 26, 466, 0xc762, [4]level{{4, 24}, {8, 22}, {10, 26}, {11, 28}}},  // 12
	{32, 28, 532, 0xd847, [4]level{{4, 26}, {9, 22}, {12, 24}, {16, 22}}},  // 13
	{24, 20, 581, 0xe60d, [4]level{{4, 30}, {9, 24}, {16, 20}, {16, 24}}},  // 14
	{24, 22, 655, 0xf928, [4]level{{6, 22}, {10, 24}, {12, 30}, {18, 24}}},  // 15
	{24, 24, 733, 0x10b78, [4]level{{6, 24}, {10, 28}, {17, 24}, {16, 30}}},  // 16
	{28, 24, 815, 0x1145d, [4]level{{6, 28}, {11, 28}, {16, 28}, {19, 28}}},  // 17
	{28, 26, 901, 0x12a17, [4]level{{6, 30}, {13, 26}, {18, 28}, {21, 28}}},  // 18
	{28, 28, 991, 0x13532, [4]level{{7, 28}, {14, 26}, {21, 26}, {25, 26}}},  // 19
	{32, 28, 1085, 0x149a6, [4]level{{8, 28}, {16, 26}, {20, 30}, {25, 28}}},  // 20
	{26, 22, 1156, 0x15683, [4]level{{8, 28}, {17, 26}, {23, 28}, {25, 30}}},  // 21
	{24, 24, 1258, 0x168c9, [4]level{{9, 28}, {17, 28}, {23, 30}, {34, 24}}},  // 22
	{28, 24, 1364, 0x177ec, [4]level{{9, 30}, {18, 28}, {25, 30}, {30, 30}}},  // 23
	{26, 26, 1474, 0x18ec4, [4]level{{10, 30}, {20, 28}, {27, 30}, {32, 30}}},  // 24
	{30, 26, 1588, 0x191e1, [4]level{{12, 26}, {21, 28}, {29, 30}, {35, 30}}},  // 25
	{28, 28, 1706, 0x1afab, [4]level{{12, 28}, {23, 28}, {34, 28}, {37, 30}}},  // 26
	{32, 28, 1828, 0x1b08e, [4]level{{12, 30}, {25, 28}, {34, 30}, {40, 30}}},  // 27
	{24, 24, 1921, 0x1cc1a, [4]level{{13, 30}, {26, 28}, {35, 30}, {42, 30}}},  // 28
	{28, 24, 2051, 0x1d33f, [4]level{{14, 30}, {28, 28}, {38, 30}, {45, 30}}},  // 29
	{24, 26, 2185, 0x1ed75, [4]level{{15, 30}, {29, 28}, {40, 30}, {48, 30}}},  // 30
	{28, 26, 2323, 0x1f250, [4]level{{16, 30}, {31, 28}, {43, 30}, {51, 30}}},  // 31
	{32, 26, 2465, 0x209d5, [4]level{{17, 30}, {33, 28}, {45, 30}, {54, 30}}},  // 32
	{28, 28, 2611, 0x216f0, [4]level{{18, 30}, {35, 28}, {48, 30}, {57, 30}}},  // 33
	{32, 28, 2761, 0x228ba, [4]level{{19, 30}, {37, 28}, {51, 30}, {60, 30}}},  // 34
	{28, 24, 2876, 0x2379f, [4]level{{19, 30}, {38, 28}, {53, 30}, {63, 30}}},  // 35
	{22, 26, 3034, 0x24b0b, [4]level{{20, 30}, {40, 28}, {56, 30}, {66, 30}}},  // 36
	{26, 26, 3196, 0x2542e, [4]level{{21, 30}, {43, 28}, {59, 30}, {70, 30}}},  // 37
	{30, 26, 3362, 0x26a64, [4]level{{22, 30}, {45, 28}, {62, 30}, {74, 30}}},  // 38
	{24, 28, 3532, 0x27541, [4]level{{24, 30}, {47, 28}, {65, 30}, {77, 30}}},  // 39
	{28, 28, 3706, 0x28c69, [4]level{{25, 30}, {49, 28}, {68, 30}, {81, 30}}},  // 40
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
	const ti = 6 // timing is in row/column 6 (counting from 0)
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
	info := &vtab[v]
	for x := 4; x+5 < siz; {
		for y := 4; y+5 < siz; {
			// don't overwrite timing markers
			if (x < 7 && y < 7) || (x < 7 && y+5 >= siz-7) || (x+5 >= siz-7 && y < 7) {
			} else {
				alignBox(m, x, y)
			}
			if y == 4 {
				y = info.apos
			} else {
				y += info.astride
			}
		}
		if x == 4 {
			x = info.apos
		} else {
			x += info.astride
		}
	}

	// Version pattern.
	pat := vtab[v].pattern
	if pat != 0 {
		v := pat
		for x := 0; x < 6; x++ {
			for y := 0; y < 3; y++ {
				p := PVersion.Pixel()
				if v&1 != 0 {
					p |= Black
				}
				m[siz-11+y][x] = p
				m[x][siz-11+y] = p
				v >>= 1
			}
		}
	}

	// One lonely black pixel
	m[siz-8][8] = Unused.Pixel() | Black

	return p, nil
}

// fplan adds the format pixels
func fplan(l Level, m Mask, p *Plan) os.Error {
	// Format pixels.
	fb := uint32(l^1) << 13 // level: L=01, M=00, Q=11, H=10
	fb |= uint32(m) << 10   // mask
	const formatPoly = 0x537
	rem := fb
	for i := 14; i >= 10; i-- {
		if rem&(1<<uint(i)) != 0 {
			rem ^= formatPoly << uint(i-10)
		}
	}
	fb |= rem
	invert := uint32(0x5412)
	siz := len(p.Pixel)
	for i := 0; i < 15; i++ {
		pix := Format.Pixel() + OffsetPixel(i)
		if (fb>>uint(i))&1 == 1 {
			pix |= Black
		}
		if (invert>>uint(i))&1 == 1 {
			pix ^= Invert | Black
		}
		// top left
		switch {
		case i < 6:
			p.Pixel[i][8] = pix
		case i < 8:
			p.Pixel[i+1][8] = pix
		case i < 9:
			p.Pixel[8][7] = pix
		default:
			p.Pixel[8][14-i] = pix
		}
		// bottom right
		switch {
		case i < 8:
			p.Pixel[8][siz-1-i] = pix
		default:
			p.Pixel[siz-1-(14-i)][8] = pix
		}
	}
	return nil
}

// lplan edits a version-only Plan to add information
// about the error correction levels.
func lplan(v Version, l Level, p *Plan) os.Error {
	p.Level = l

	nblock := vtab[v].level[l].nblock
	ne := vtab[v].level[l].check
	nde := (vtab[v].bytes - ne*nblock)/nblock
	extra := (vtab[v].bytes - ne*nblock)%nblock
	dataBits := (nde*nblock+extra)*8
	checkBits := ne*nblock*8

	// Make data + checksum pixels.
	data := make([]Pixel, dataBits)
	for i := range data {
		data[i] = Data.Pixel() | OffsetPixel(i)
	}
	check := make([]Pixel, checkBits)
	for i := range check {
		check[i] = Check.Pixel() | OffsetPixel(i)
	}
	
	// Split into blocks.
	dataList := make([][]Pixel, nblock)
	checkList := make([][]Pixel, nblock)
	for i := 0; i < nblock; i++ {
		// The last few blocks have an extra data byte (8 pixels).
		nd := nde
		if i >= nblock-extra {
			nd++
		}
		dataList[i], data = data[0:nd*8], data[nd*8:]
		checkList[i], check = check[0:ne*8], check[ne*8:]
	}
	
	// Build up bit sequence, taking first byte of each block,
	// then second byte, and so on.  Then checksums.
	bits := make([]Pixel, dataBits+checkBits)
	dst := bits
	for i := 0; i < nde+1; i++ {
		for _, b := range dataList {
			if i*8 < len(b) {
				copy(dst, b[i*8:(i+1)*8])
				dst = dst[8:]
			}
		}
	}
	for i := 0; i < ne; i++ {
		for _, b := range checkList {
			if i*8 < len(b) {
				copy(dst, b[i*8:(i+1)*8])
				dst = dst[8:]
			}
		}
	}
	if len(dst) != 0 {
		panic("dst math")
	}

	// Sweep up pair of columns,
	// then down, assigning to right then left pixel.
	// Repeat.
	// See Figure 2 of http://www.pclviewer.com/rs2/qrtopology.htm
	siz := len(p.Pixel)
	rem := make([]Pixel, 7)
	for i := range rem {
		rem[i] = Extra.Pixel()
	}
	src := append(bits, rem...)
	for x := siz; x > 0; {
		for y := siz - 1; y >= 0; y-- {
			if p.Pixel[y][x-1].Role() == 0 {
				p.Pixel[y][x-1], src = src[0], src[1:]
			}
			if p.Pixel[y][x-2].Role() == 0 {
				p.Pixel[y][x-2], src = src[0], src[1:]
			}
		}
		x -= 2
		if x == 7 { // vertical timing strip
			x--
		}
		for y := 0; y < siz; y++ {
			if p.Pixel[y][x-1].Role() == 0 {
				p.Pixel[y][x-1], src = src[0], src[1:]
			}
			if p.Pixel[y][x-2].Role() == 0 {
				p.Pixel[y][x-2], src = src[0], src[1:]
			}
		}
		x -= 2
	}
	return nil
}

// http://www.swetake.com/qr/qr5_en.html
var mfunc = []func(int, int) bool{
	func(i, j int) bool { return (i+j)%2 == 0 },
	func(i, j int) bool { return i%2 == 0 },
	func(i, j int) bool { return j%3 == 0 },
	func(i, j int) bool { return (i+j)%3 == 0 },
	func(i, j int) bool { return (i/2+j/3)%2 == 0 },
	func(i, j int) bool { return i*j%2+i*j%3 == 0 },
	func(i, j int) bool { return (i*j%2+i*j%3)%2 == 0 },
	func(i, j int) bool { return (i*j%3+(i+j)%2)%2 == 0 },
}

// mplan edits a version+level-only Plan to add the mask.
func mplan(m Mask, p *Plan) os.Error {
	f := mfunc[m]
	p.Mask = m
	for y, row := range p.Pixel {
		for x, pix := range row {
			if r := pix.Role(); (r == Data || r == Check || r == Extra) && f(y, x) {
				row[x] ^= Black | Invert
			}
		}
	}
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
