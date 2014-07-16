package main

import "fmt"

func Aconv_list6(i int) string {
	return anames6[i]
}

func Dconv_list6(p *Prog, isSplitConst bool, a *Addr) string {
	var i int
	i = a.typ
	str := ""
	if isSplitConst {
		if i == int(D_CONST_6) {
			str = fmt.Sprintf("$%d-%d", a.offset&0xffffffff, a.offset>>32)
		} else {
			// ATEXT dst is not constant
			str = fmt.Sprintf("!!%s", Dconv_list6(p, false, a))
		}
		goto brk
	}
	if i >= int(D_INDIR_6) {
		if a.offset != 0 {
			str = fmt.Sprintf("%d(%s)", a.offset, Rconv_list6(i-int(D_INDIR_6)))
		} else {
			str = fmt.Sprintf("(%s)", Rconv_list6(i-int(D_INDIR_6)))
		}
		goto brk
	}
	switch i {
	default:
		if a.offset != 0 {
			str = fmt.Sprintf("$%d,%s", a.offset, Rconv_list6(i))
		} else {
			str = fmt.Sprintf("%s", Rconv_list6(i))
		}
		break
	case D_NONE_6:
		str = ""
		break
	case D_BRANCH_6:
		if a.sym != nil {
			str = fmt.Sprintf("%s(SB)", a.sym.name)
		} else {
			if p != nil && p.pcond != nil {
				str = fmt.Sprintf("%d", p.pcond.pc)
			} else {
				if a.u.branch != nil {
					str = fmt.Sprintf("%d", a.u.branch.pc)
				} else {
					str = fmt.Sprintf("%d(PC)", a.offset)
				}
			}
		}
		break
	case D_EXTERN_6:
		str = fmt.Sprintf("%s+%d(SB)", a.sym.name, a.offset)
		break
	case D_STATIC_6:
		str = fmt.Sprintf("%s<>+%d(SB)", a.sym.name, a.offset)
		break
	case D_AUTO_6:
		if a.sym != nil {
			str = fmt.Sprintf("%s+%d(SP)", a.sym.name, a.offset)
		} else {
			str = fmt.Sprintf("%d(SP)", a.offset)
		}
		break
	case D_PARAM_6:
		if a.sym != nil {
			str = fmt.Sprintf("%s+%d(FP)", a.sym.name, a.offset)
		} else {
			str = fmt.Sprintf("%d(FP)", a.offset)
		}
		break
	case D_CONST_6:
		str = fmt.Sprintf("$%d", a.offset)
		break
	case D_FCONST_6:
		str = fmt.Sprintf("$(%.17g)", a.u.dval)
		break
	case D_SCONST_6:
		str = fmt.Sprintf("$%q", a.u.sval)
		break
	case D_ADDR_6:
		a.typ = a.index
		a.index = int(D_NONE_6)
		str = fmt.Sprintf("$%s", Dconv_list6(p, false, a))
		a.index = a.typ
		a.typ = int(D_ADDR_6)
		goto conv
	}
brk:
	if a.index != int(D_NONE_6) {
		s := fmt.Sprintf("(%s*%d)", Rconv_list6(int(a.index)), int(a.scale))
		str += s
	}
conv:
	return str
}

func Pconv_list6(ctxt *Link, p *Prog) string {
	switch p.as {
	case ADATA_6:
		return fmt.Sprintf("%.5d (%s)	%s	%s/%d,%s", p.pc, linklinefmt(ctxt, p.lineno, false, false), Aconv_list6(p.as), Dconv_list6(p, false, &p.from), p.from.scale, Dconv_list6(p, false, &p.to))
	case ATEXT_6:
		if p.from.scale != 0 {
			return fmt.Sprintf("%.5d (%s)	%s	%s,%d,%s", p.pc, linklinefmt(ctxt, p.lineno, false, false), Aconv_list6(p.as), Dconv_list6(p, false, &p.from), p.from.scale, Dconv_list6(p, true, &p.to))
		}
		return fmt.Sprintf("%.5d (%s)	%s	%s,%s", p.pc, linklinefmt(ctxt, p.lineno, false, false), Aconv_list6(p.as), Dconv_list6(p, false, &p.from), Dconv_list6(p, true, &p.to))
	default:
		return fmt.Sprintf("%.5d (%s)	%s	%s,%s", p.pc, linklinefmt(ctxt, p.lineno, false, false), Aconv_list6(p.as), Dconv_list6(p, false, &p.from), Dconv_list6(p, false, &p.to))
	}
}

func Rconv_list6(r int) string {
	if r >= int(D_AL_6) && r <= int(D_NONE_6) {
		return regstr[r-int(D_AL_6)]
	} else {
		return fmt.Sprintf("gok(%d)", r)
	}
}

// Inferno utils/6c/list.c
// http://code.google.com/p/inferno-os/source/browse/utils/6c/list.c
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
//
// Format conversions
//	%A int		Opcodes (instruction mnemonics)
//
//	%s Addr*	Addresses (instruction operands)
//		Flags: "%lD": seperate the high and low words of a constant by "-"
//
//	%P Prog*	Instructions
//
//	%R int		Registers
//
//	%$ char*	String constant addresses (for internal use only)
