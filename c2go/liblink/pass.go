package main

// Inferno utils/6l/pass.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/pass.c
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
// Code and data passes.
func brchain(ctxt *Link, p *Prog) *Prog {
	var i int
	for i = 0; i < 20; i++ {
		if p == nil || p.as != ctxt.arch.AJMP {
			return p
		}
		p = p.pcond
	}
	return (*Prog)(nil)
}

func brloop(ctxt *Link, p *Prog) *Prog {
	var c int
	var q *Prog
	c = 0
	for q = p; q != nil; q = q.pcond {
		if q.as != ctxt.arch.AJMP {
			break
		}
		c++
		if c >= 5000 {
			return (*Prog)(nil)
		}
	}
	return q
}

func linkpatch(ctxt *Link, sym *LSym) {
	var c int
	var p *Prog
	var q *Prog
	ctxt.cursym = sym
	for p = sym.text; p != nil; p = p.link {
		if ctxt.arch.progedit != nil {
			ctxt.arch.progedit(ctxt, p)
		}
		if p.to.typ != ctxt.arch.D_BRANCH {
			continue
		}
		if p.to.u.branch != nil {
			// TODO: Remove to.u.branch in favor of p->pcond.
			p.pcond = p.to.u.branch
			continue
		}
		if p.to.sym != nil {
			continue
		}
		c = int(p.to.offset)
		for q = sym.text; q != nil; {
			if c == q.pc {
				break
			}
			if q.forwd != nil && c >= q.forwd.pc {
				q = q.forwd
			} else {
				q = q.link
			}
		}
		if q == nil {
			var tmp string
			if p.to.sym != nil {
				tmp = p.to.sym.name
			} else {
				tmp = "<nil>"
			}
			ctxt.diag("branch out of range (%#ux)\n%P [%s]", c, p, tmp)
			p.to.typ = ctxt.arch.D_NONE
		}
		p.to.u.branch = q
		p.pcond = q
	}
	for p = sym.text; p != nil; p = p.link {
		p.mark = 0 /* initialization for follow */
		if p.pcond != nil {
			p.pcond = brloop(ctxt, p.pcond)
			if p.pcond != nil {
				if p.to.typ == ctxt.arch.D_BRANCH {
					p.to.offset = int64(p.pcond.pc)
				}
			}
		}
	}
}
