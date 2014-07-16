package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func linklookup(ctxt *Link, name string, v int) *LSym {
	return _lookup_sym(ctxt, name, v, 1)
}

func linknew(arch *LinkArch) *Link {
	var ctxt *Link
	var p string
	var buf string
	ctxt = new(Link)
	ctxt.arch = arch
	ctxt.version = int(HistVersion)
	ctxt.goroot = getgoroot()
	ctxt.goroot_final = os.Getenv("GOROOT_FINAL")
	p = getgoarch()
	if p != arch.name {
		sysfatal("invalid goarch %s (want %s)", p, arch.name)
	}
	buf, err := os.Getwd()
	if err != nil {
		buf = "/???"
	}
	if yy_isalpha_sym(int(buf[0])) && buf[1] == ':' {
		// On Windows.
		ctxt.windows = 1
		// Canonicalize path by converting \ to / (Windows accepts both).
		buf = strings.Replace(buf, `\`, `/`, -1)
	}
	ctxt.pathname = buf
	ctxt.headtype = headtype(getgoos())
	if ctxt.headtype < 0 {
		sysfatal("unknown goos %s", getgoos())
	}
	// Record thread-local storage offset.
	// TODO(rsc): Move tlsoffset back into the linker.
	switch ctxt.headtype {
	default:
		sysfatal("unknown thread-local storage offset for %s", headstr(ctxt.headtype))
	case Hplan9:
	case Hwindows:
		break
	/*
	 * ELF uses TLS offset negative from FS.
	 * Translate 0(FS) and 8(FS) into -16(FS) and -8(FS).
	 * Known to low-level assembly in package runtime and runtime/cgo.
	 */
	case Hlinux:
	case Hfreebsd:
	case Hnetbsd:
	case Hopenbsd:
	case Hdragonfly:
	case Hsolaris:
		ctxt.tlsoffset = -2 * ctxt.arch.ptrsize
		break
	case Hnacl:
		switch ctxt.arch.thechar {
		default:
			sysfatal("unknown thread-local storage offset for nacl/%s", ctxt.arch.name)
		case '6':
			ctxt.tlsoffset = 0
			break
		case '8':
			ctxt.tlsoffset = -8
			break
		case '5':
			ctxt.tlsoffset = 0
			break
		}
		break
	/*
	 * OS X system constants - offset from 0(GS) to our TLS.
	 * Explained in ../../pkg/runtime/cgo/gcc_darwin_*.c.
	 */
	case Hdarwin:
		switch ctxt.arch.thechar {
		default:
			sysfatal("unknown thread-local storage offset for darwin/%s", ctxt.arch.name)
		case '6':
			ctxt.tlsoffset = 0x8a0
			break
		case '8':
			ctxt.tlsoffset = 0x468
			break
		}
		break
	}
	// On arm, record goarm.
	if ctxt.arch.thechar == '5' {
		p = getgoarm()
		if p != "" {
			x, _ := strconv.Atoi(p)
			ctxt.goarm = int32(x)
		} else {
			ctxt.goarm = 6
		}
	}
	return ctxt
}

func linknewsym(ctxt *Link, symb string, v int) *LSym {
	var s *LSym
	s = new(LSym)
	*s = LSym{}
	s.dynid = -1
	s.plt = -1
	s.got = -1
	s.name = symb
	s.typ = 0
	s.version = v
	s.value = 0
	s.sig = 0
	s.size = 0
	ctxt.nsymbol++
	s.allsym = ctxt.allsym
	ctxt.allsym = s
	return s
}

// read-only lookup
func linkrlookup(ctxt *Link, name string, v int) *LSym {
	return _lookup_sym(ctxt, name, v, 0)
}

func (s *LSym) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.name
}

func headtype(name string) int {
	var i int
	for i = 0; headers_sym[i].name != ""; i++ {
		if name == headers_sym[i].name {
			return headers_sym[i].val
		}
	}
	return -1
}

func headstr(v int) string {
	var i int
	for i = 0; headers_sym[i].name != ""; i++ {
		if v == headers_sym[i].val {
			return headers_sym[i].name
		}
	}
	return fmt.Sprintf("%d", v)
}

// Derived from Inferno utils/6l/obj.c and utils/6l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/obj.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/span.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
func yy_isalpha_sym(c int) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

var headers_sym = []struct {
	name string
	val  int
}{

	{"android", Hlinux},
	{"darwin", Hdarwin},
	{"dragonfly", Hdragonfly},
	{"elf", Helf},
	{"freebsd", Hfreebsd},
	{"linux", Hlinux},
	{"nacl", Hnacl},
	{"netbsd", Hnetbsd},
	{"openbsd", Hopenbsd},
	{"plan9", Hplan9},
	{"solaris", Hsolaris},
	{"windows", Hwindows},
	{"windowsgui", Hwindows},
}

func _lookup_sym(ctxt *Link, symb string, v int, creat int) *LSym {
	var s *LSym
	var p string
	var h int
	var c int
	h = v
	for p = symb; len(p) > 0; p = p[1:] {
		c = int(p[0])
		h = h + h + h + c
	}
	h &= 0xffffff
	h %= int(LINKHASH)
	for s = ctxt.hash[h]; s != nil; s = s.hash {
		if s.version == v && s.name == symb {
			return s
		}
	}
	if !(creat != 0) {
		return (*LSym)(nil)
	}
	s = linknewsym(ctxt, symb, v)
	s.extname = s.name
	s.hash = ctxt.hash[h]
	ctxt.hash[h] = s
	return s
}
