package main

import "fmt"

var regstr = []string{
	"AL", /* [D_AL] */
	"CL",
	"DL",
	"BL",
	"SPB",
	"BPB",
	"SIB",
	"DIB",
	"R8B",
	"R9B",
	"R10B",
	"R11B",
	"R12B",
	"R13B",
	"R14B",
	"R15B",
	"AX", /* [D_AX] */
	"CX",
	"DX",
	"BX",
	"SP",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
	"AH",
	"CH",
	"DH",
	"BH",
	"F0", /* [D_F0] */
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"M0",
	"M1",
	"M2",
	"M3",
	"M4",
	"M5",
	"M6",
	"M7",
	"X0",
	"X1",
	"X2",
	"X3",
	"X4",
	"X5",
	"X6",
	"X7",
	"X8",
	"X9",
	"X10",
	"X11",
	"X12",
	"X13",
	"X14",
	"X15",
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
	"CR8",
	"CR9",
	"CR10",
	"CR11",
	"CR12",
	"CR13",
	"CR14",
	"CR15",
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
	"TLS",  /* [D_TLS] */
	"NONE", /* [D_NONE] */
}

func Aconv_list8(i int) string {
	return anames8[i]
}

func Dconv_list8(p *Prog, isSplitConst bool, a *Addr) string {
	var i int
	i = a.typ
	str := ""
	if isSplitConst {
		if i == int(D_CONST2_8) {
			str = fmt.Sprintf("$%lld-%d", a.offset, a.offset2)
		} else {
			// ATEXT dst is not constant
			str = fmt.Sprintf("!!%s", Dconv_list8(p, false, a))
		}
		goto brk
	}
	if i >= int(D_INDIR_8) {
		if a.offset != 0 {
			str = fmt.Sprintf("%lld(%R)", a.offset, i-int(D_INDIR_8))
		} else {
			str = fmt.Sprintf("(%R)", i-int(D_INDIR_8))
		}
		goto brk
	}
	switch i {
	default:
		if a.offset != 0 {
			str = fmt.Sprintf("$%lld,%R", a.offset, i)
		} else {
			str = fmt.Sprintf("%R", i)
		}
		break
	case D_NONE_8:
		str = ""
		break
	case D_BRANCH_8:
		if a.sym != nil {
			str = fmt.Sprintf("%s(SB)", a.sym.name)
		} else {
			if bigP_list8 != nil && bigP_list8.pcond != nil {
				str = fmt.Sprintf("%lld", bigP_list8.pcond.pc)
			} else {
				if a.u.branch != nil {
					str = fmt.Sprintf("%lld", a.u.branch.pc)
				} else {
					str = fmt.Sprintf("%lld(PC)", a.offset)
				}
			}
		}
		break
	case D_EXTERN_8:
		str = fmt.Sprintf("%s+%lld(SB)", a.sym.name, a.offset)
		break
	case D_STATIC_8:
		str = fmt.Sprintf("%s<>+%lld(SB)", a.sym.name, a.offset)
		break
	case D_AUTO_8:
		if a.sym != nil {
			str = fmt.Sprintf("%s+%lld(SP)", a.sym.name, a.offset)
		} else {
			str = fmt.Sprintf("%lld(SP)", a.offset)
		}
		break
	case D_PARAM_8:
		if a.sym != nil {
			str = fmt.Sprintf("%s+%lld(FP)", a.sym.name, a.offset)
		} else {
			str = fmt.Sprintf("%lld(FP)", a.offset)
		}
		break
	case D_CONST_8:
		str = fmt.Sprintf("$%lld", a.offset)
		break
	case D_CONST2_8:
		if !isSplitConst {
			// D_CONST2 outside of ATEXT should not happen
			str = fmt.Sprintf("!!$%lld-%d", a.offset, a.offset2)
		}
		break
	case D_FCONST_8:
		str = fmt.Sprintf("$(%.17g)", a.u.dval)
		break
	case D_SCONST_8:
		str = fmt.Sprintf("$\"%$\"", a.u.sval)
		break
	case D_ADDR_8:
		a.typ = a.index
		a.index = int(D_NONE_8)
		str = fmt.Sprintf("$%s", Dconv_list8(p, false, a))
		a.index = a.typ
		a.typ = int(D_ADDR_8)
		goto conv
	}
brk:
	if a.index != int(D_NONE_8) {
		s := fmt.Sprintf("(%R*%d)", int(a.index), int(a.scale))
		str += s
	}
conv:
	return str
} /* [D_AL] */ /* [D_AX] */ /* [D_F0] */ /* [D_CS] */ /* [D_GDTR] */ /* [D_IDTR] */ /* [D_LDTR] */ /* [D_MSW] */ /* [D_TASK] */ /* [D_CR] */ /* [D_DR] */ /* [D_TR] */ /* [D_X0] */ /* [D_TLS] */ /* [D_NONE] */

func Pconv_list8(p *Prog) string {
	switch p.as {
	case ADATA_8:
		return fmt.Sprintf("%.5lld (%L)	%A	%s/%d,%s", p.pc, p.lineno, Aconv_list8(p.as), Dconv_list8(p, false, &p.from), p.from.scale, Dconv_list8(p, false, &p.to))
	case ATEXT_8:
		if p.from.scale != 0 {
			return fmt.Sprintf("%.5lld (%L)	%A	%s,%d,%s", p.pc, p.lineno, Aconv_list8(p.as), Dconv_list8(p, false, &p.from), p.from.scale, Dconv_list8(p, true, &p.to))
		}
		return fmt.Sprintf("%.5lld (%L)	%A	%s,%s", p.pc, p.lineno, Aconv_list8(p.as), Dconv_list8(p, false, &p.from), Dconv_list8(p, true, &p.to))
	default:
		return fmt.Sprintf("%.5lld (%L)	%A	%s,%s", p.pc, p.lineno, Aconv_list8(p.as), Dconv_list8(p, false, &p.from), Dconv_list8(p, false, &p.to))
	}
}

func Rconv_list8(r int) string {
	if r >= int(D_AL_8) && r <= int(D_NONE_8) {
		return regstr[r-int(D_AL_8)]
	} else {
		return fmt.Sprintf("gok(%d)", r)
	}
}

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

var bigP_list8 *Prog
