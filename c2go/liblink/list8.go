package main

import "fmt"

// Inferno utils/8c/list.c
// http://code.google.com/p/inferno-os/source/browse/utils/8c/list.c
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
const (
	STRINGSZ_list8 = 1000
)

func Pconv_list8(p *Prog) string {
	var str string
	var fp string

	switch p.as {
	case ADATA_8:
		str = fmt.Sprintf("%.5d (%v)\t%v\t%v/%d,%v", p.pc, p.Line(), Aconv_list8(p.as), Dconv_list8(p, 0, &p.from), p.from.scale, Dconv_list8(p, 0, &p.to))
	case ATEXT_8:
		if p.from.scale != 0 {
			str = fmt.Sprintf("%.5d (%v)\t%v\t%v,%d,%v", p.pc, p.Line(), Aconv_list8(p.as), Dconv_list8(p, 0, &p.from), p.from.scale, Dconv_list8(p, fmtLong, &p.to))
			break
		}
		str = fmt.Sprintf("%.5d (%v)\t%v\t%v,%v", p.pc, p.Line(), Aconv_list8(p.as), Dconv_list8(p, 0, &p.from), Dconv_list8(p, fmtLong, &p.to))
	default:
		str = fmt.Sprintf("%.5d (%v)\t%v\t%v,%v", p.pc, p.Line(), Aconv_list8(p.as), Dconv_list8(p, 0, &p.from), Dconv_list8(p, 0, &p.to))
		break
	}

	fp += str
	return fp
}

func Aconv_list8(i int) string {
	var fp string

	fp += anames8[i]
	return fp
}

func Dconv_list8(p *Prog, flag int, a *Addr) string {
	var str string
	var s string
	var fp string
	var i int

	i = a.typ
	if flag&fmtLong != 0 /*untyped*/ {
		if i == D_CONST2_8 {
			str = fmt.Sprintf("$%d-%d", a.offset, a.offset2)
		} else {
			// ATEXT dst is not constant
			str = fmt.Sprintf("!!%v", Dconv_list8(p, 0, a))
		}
		goto brk
	}
	if i >= D_INDIR_8 {
		if a.offset != 0 {
			str = fmt.Sprintf("%d(%v)", a.offset, Rconv_list8(i-D_INDIR_8))
		} else {
			str = fmt.Sprintf("(%v)", Rconv_list8(i-D_INDIR_8))
		}
		goto brk
	}
	switch i {
	default:
		if a.offset != 0 {
			str = fmt.Sprintf("$%d,%v", a.offset, Rconv_list8(i))
		} else {
			str = fmt.Sprintf("%v", Rconv_list8(i))
		}
	case D_NONE_8:
		str = ""
	case D_BRANCH_8:
		if a.sym != nil {
			str = fmt.Sprintf("%s(SB)", a.sym.name)
		} else if p != nil && p.pcond != nil {
			str = fmt.Sprintf("%d", p.pcond.pc)
		} else if a.u.branch != nil {
			str = fmt.Sprintf("%d", a.u.branch.pc)
		} else {
			str = fmt.Sprintf("%d(PC)", a.offset)
		}
	case D_EXTERN_8:
		str = fmt.Sprintf("%s+%d(SB)", a.sym.name, a.offset)
	case D_STATIC_8:
		str = fmt.Sprintf("%s<>+%d(SB)", a.sym.name, a.offset)
	case D_AUTO_8:
		if a.sym != nil {
			str = fmt.Sprintf("%s+%d(SP)", a.sym.name, a.offset)
		} else {
			str = fmt.Sprintf("%d(SP)", a.offset)
		}
	case D_PARAM_8:
		if a.sym != nil {
			str = fmt.Sprintf("%s+%d(FP)", a.sym.name, a.offset)
		} else {
			str = fmt.Sprintf("%d(FP)", a.offset)
		}
	case D_CONST_8:
		str = fmt.Sprintf("$%d", a.offset)
	case D_CONST2_8:
		if flag&fmtLong == 0 /*untyped*/ {
			// D_CONST2 outside of ATEXT should not happen
			str = fmt.Sprintf("!!$%d-%d", a.offset, a.offset2)
		}
	case D_FCONST_8:
		str = fmt.Sprintf("$(%.17g)", a.u.dval)
	case D_SCONST_8:
		str = fmt.Sprintf("$\"%q\"", a.u.sval)
	case D_ADDR_8:
		a.typ = a.index
		a.index = D_NONE_8
		str = fmt.Sprintf("$%v", Dconv_list8(p, 0, a))
		a.index = a.typ
		a.typ = D_ADDR_8
		goto conv
	}
brk:
	if a.index != D_NONE_8 {
		s = fmt.Sprintf("(%v*%d)", Rconv_list8(int(a.index)), int(a.scale))
		str += s
	}
conv:
	fp += str
	return fp
}

var regstr_list8 = []string{
	"AL", /* [D_AL] */
	"CL",
	"DL",
	"BL",
	"AH",
	"CH",
	"DH",
	"BH",
	"AX", /* [D_AX] */
	"CX",
	"DX",
	"BX",
	"SP",
	"BP",
	"SI",
	"DI",
	"F0", /* [D_F0] */
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"CS", /* [D_CS] */
	"SS",
	"DS",
	"ES",
	"FS",
	"GS",
	"GDTR", /* [D_GDTR] */
	"IDTR", /* [D_IDTR] */
	"LDTR", /* [D_LDTR] */
	"MSW",  /* [D_MSW] */
	"TASK", /* [D_TASK] */
	"CR0",  /* [D_CR] */
	"CR1",
	"CR2",
	"CR3",
	"CR4",
	"CR5",
	"CR6",
	"CR7",
	"DR0", /* [D_DR] */
	"DR1",
	"DR2",
	"DR3",
	"DR4",
	"DR5",
	"DR6",
	"DR7",
	"TR0", /* [D_TR] */
	"TR1",
	"TR2",
	"TR3",
	"TR4",
	"TR5",
	"TR6",
	"TR7",
	"X0", /* [D_X0] */
	"X1",
	"X2",
	"X3",
	"X4",
	"X5",
	"X6",
	"X7",
	"TLS",  /* [D_TLS] */
	"NONE", /* [D_NONE] */
}

func Rconv_list8(r int) string {
	var str string
	var fp string

	if r >= D_AL_8 && r <= D_NONE_8 {
		str = fmt.Sprintf("%s", regstr_list8[r-D_AL_8])
	} else {
		str = fmt.Sprintf("gok(%d)", r)
	}
	fp += str
	return fp
}
