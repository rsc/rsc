package main

import (
	"os"
	"path"
	"strings"
)

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
func addlib(ctxt *Link, src string, obj string, pathname string) {
	var i int

	name := path.Clean(pathname)

	// runtime.a -> runtime
	short := strings.TrimSuffix(name, ".a")

	// already loaded?
	for i := range ctxt.library {
		if ctxt.library[i].pkg == short {
			return
		}
	}

	var pname string
	// runtime -> runtime.a for search
	if (!(ctxt.windows != 0) && name[0] == '/') || (ctxt.windows != 0 && name[1] == ':') {
		pname = name
	} else {
		// try dot, -L "libdir", and then goroot.
		for i = 0; i < ctxt.nlibdir; i++ {
			pname = ctxt.libdir[i] + "/" + name
			if _, err := os.Stat(pname); !os.IsNotExist(err) {
				break
			}
		}
	}
	pname = path.Clean(pname)

	// runtime.a -> runtime
	pname = strings.TrimSuffix(pname, ".a")

	if ctxt.debugvlog > 1 && ctxt.bso != nil {
		Bprint(ctxt.bso, "%5.2f addlib: %s %s pulls in %s\n", cputime(), obj, src, pname)
	}
	addlibpath(ctxt, src, obj, pname, name)
}

/*
 * add library to library list.
 *	srcref: src file referring to package
 *	objref: object file referring to package
 *	file: object file, e.g., /home/rsc/go/pkg/container/vector.a
 *	pkg: package import path, e.g. container/vector
 */
func addlibpath(ctxt *Link, srcref string, objref string, file string, pkg string) {
	for i := range ctxt.library {
		if file == ctxt.library[i].file {
			return
		}
	}
	if ctxt.debugvlog > 1 && ctxt.bso != nil {
		Bprint(ctxt.bso, "%5.2f addlibpath: srcref: %s objref: %s file: %s pkg: %s\n", cputime(), srcref, objref, file, pkg)
	}
	ctxt.library = append(ctxt.library, Library{
		objref: objref,
		srcref: srcref,
		file:   file,
		pkg:    pkg,
	})
}

func mkfwd(sym *LSym) {
	var p *Prog
	var i int
	var dwn [LOG_ld]int32
	var cnt [LOG_ld]int32
	var lst [LOG_ld]*Prog
	for i = 0; i < int(LOG_ld); i++ {
		if i == 0 {
			cnt[i] = 1
		} else {
			cnt[i] = int32(LOG_ld) * cnt[i-1]
		}
		dwn[i] = 1
		lst[i] = (*Prog)(nil)
	}
	i = 0
	for p = sym.text; p != nil && p.link != nil; p = p.link {
		i--
		if i < 0 {
			i = int(LOG_ld - 1)
		}
		p.forwd = (*Prog)(nil)
		dwn[i]--
		if dwn[i] <= 0 {
			dwn[i] = cnt[i]
			if lst[i] != nil {
				lst[i].forwd = p
			}
			lst[i] = p
		}
	}
}

func copyp(ctxt *Link, q *Prog) *Prog {
	var p *Prog
	p = ctxt.arch.prg()
	*p = *q
	return p
}

func appendp(ctxt *Link, q *Prog) *Prog {
	var p *Prog
	p = ctxt.arch.prg()
	p.link = q.link
	q.link = p
	p.lineno = q.lineno
	p.mode = q.mode
	return p
}

func atolwhex(s string) int64 {
	var n int64
	var f int
	n = 0
	f = 0
	for s[0] == ' ' || s[0] == '\t' {
		s = s[1:]
	}
	if s[0] == '-' || s[0] == '+' {
		var tmp string = s
		s = s[1:]
		if tmp[0] == '-' {
			f = 1
		}
		for s[0] == ' ' || s[0] == '\t' {
			s = s[1:]
		}
	}
	if s[0] == '0' && s[1] != 0 {
		if s[1] == 'x' || s[1] == 'X' {
			s = s[2:]
			for {
				if s[0] >= '0' && s[0] <= '9' {
					n = n*16 + int64(s[0]) - '0'
					s = s[1:]
				} else if s[0] >= 'a' && s[0] <= 'f' {
					n = n*16 + int64(s[0]) - 'a' + 10
					s = s[1:]
				} else if s[0] >= 'A' && s[0] <= 'F' {
					n = n*16 + int64(s[0]) - 'A' + 10
					s = s[1:]
				} else {
					break
				}
			}
		} else {
			for s[0] >= '0' && s[0] <= '7' {
				n = n*8 + int64(s[0]) - '0'
				s = s[1:]
			}
		}
	} else {
		for s[0] >= '0' && s[0] <= '9' {
			n = n*10 + int64(s[0]) - '0'
			s = s[1:]
		}
	}
	if f != 0 {
		n = -n
	}
	return n
}

const (
	LOG_ld = 5
)
