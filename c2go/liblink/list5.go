package main

import "fmt"

// Inferno utils/5c/list.c
// http://code.google.com/p/inferno-os/source/browse/utils/5c/list.c
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

func Aconv_list5(a int) string {
	var s string
	s = "???"
	if a >= int(AXXX_8) && a < int(ALAST_8) {
		s = anames5[a]
	}
	return s
}

func Dconv_list5(p *Prog, a *Addr) string {
	var v int
	switch a.typ {
	default:
		return fmt.Sprintf("GOK-type(%d)", a.typ)
		break
	case D_NONE_8:
		if a.name != int(D_NONE_8) || a.reg != int(NREG_5) || a.sym != nil {
			return fmt.Sprintf("%s(R%d)(NONE)", Mconv_list5(a), a.reg)
		}
	case D_CONST_8:
		if a.reg != int(NREG_5) {
			return fmt.Sprintf("$%s(R%d)", Mconv_list5(a), a.reg)
		} else {
			return fmt.Sprintf("$%s", Mconv_list5(a))
		}
	case D_CONST2_8:
		return fmt.Sprintf("$%lld-%d", a.offset, a.offset2)
	case D_SHIFT_5:
		v = int(a.offset)
		op := "<<>>->@>"[((v>>5)&3)<<1:]
		var str string
		if v&(1<<4) != 0 /*untyped*/ {
			str = fmt.Sprintf("R%d%c%cR%d", v&15, op[0], op[1], (v>>8)&15)
		} else {
			str = fmt.Sprintf("R%d%c%c%d", v&15, op[0], op[1], (v>>7)&31)
		}
		if a.reg != int(NREG_5) {
			str += fmt.Sprintf("(R%d)", a.reg)
		}
	case D_OREG_5:
		if a.reg != int(NREG_5) {
			return fmt.Sprintf("%s(R%d)", Mconv_list5(a), a.reg)
		} else {
			return fmt.Sprintf("%s", Mconv_list5(a))
		}
	case D_REG_5:
		if a.name != int(D_NONE_8) || a.sym != nil {
			return fmt.Sprintf("%s(R%d)(REG)", Mconv_list5(a), a.reg)
		}
		return fmt.Sprintf("R%d", a.reg)
	case D_FREG_5:
		if a.name != int(D_NONE_8) || a.sym != nil {
			return fmt.Sprintf("%s(R%d)(REG)", Mconv_list5(a), a.reg)
		}
		return fmt.Sprintf("F%d", a.reg)
	case D_PSR_5:
		if a.name != int(D_NONE_8) || a.sym != nil {
			return fmt.Sprintf("%s(PSR)(REG)", Mconv_list5(a))
		}
		return fmt.Sprintf("PSR")
	case D_BRANCH_8:
		if a.sym != nil {
			return fmt.Sprintf("%s(SB)", a.sym.name)
		} else if p != nil && p.pcond != nil {
			return fmt.Sprintf("%lld", p.pcond.pc)
		} else if a.u.branch != nil {
			return fmt.Sprintf("%lld", a.u.branch.pc)
		} else {
			return fmt.Sprintf("%lld(PC)", a.offset) /*-pc*/
		}
	case D_FCONST_8:
		return fmt.Sprintf("$%.17g", a.u.dval)
	case D_SCONST_8:
		return fmt.Sprintf("$%q", a.u.sval)
	}
	return ""
}

func Mconv_list5(a *Addr) string {
	var s *LSym
	s = a.sym
	if s == nil {
		return fmt.Sprintf("%d", int(a.offset))
	}
	switch a.name {
	default:
		return fmt.Sprintf("GOK-name(%d)", a.name)
		break
	case D_NONE_8:
		return fmt.Sprintf("%lld", a.offset)
		break
	case D_EXTERN_8:
		return fmt.Sprintf("%s+%d(SB)", s.name, int(a.offset))
		break
	case D_STATIC_8:
		return fmt.Sprintf("%s<>+%d(SB)", s.name, int(a.offset))
		break
	case D_AUTO_8:
		return fmt.Sprintf("%s-%d(SP)", s.name, int(-a.offset))
		break
	case D_PARAM_8:
		return fmt.Sprintf("%s+%d(FP)", s.name, int(a.offset))
		break
	}
	return ""
}

func Pconv_list5(p *Prog) string {
	var a int
	var s int
	a = p.as
	s = p.scond
	sc := extra_list5[s&int(C_SCOND_5)]
	if s&int(C_SBIT_5) != 0 {
		sc += ".S"
	}
	if s&int(C_PBIT_5) != 0 {
		sc += ".P"
	}
	if s&int(C_WBIT_5) != 0 {
		sc += ".W"
	}
	if s&int(C_UBIT_5) != 0 { /* ambiguous with FBIT */
		sc += ".U"
	}
	if a == int(AMOVM_5) {
		if p.from.typ == int(D_CONST_8) {
			return fmt.Sprintf("%.5lld (%L)	%s%s	%s,%s", p.pc, p.lineno, Aconv_list5(a), sc, RAconv_list5(&p.from), Dconv_list5(p, &p.to))
		} else if p.to.typ == int(D_CONST_8) {
			return fmt.Sprintf("%.5lld (%L)	%s%s	%s,%s", p.pc, p.lineno, Aconv_list5(a), sc, Dconv_list5(p, &p.from), RAconv_list5(&p.to))
		} else {
			return fmt.Sprintf("%.5lld (%L)	%s%s	%s,%s", p.pc, p.lineno, Aconv_list5(a), sc, Dconv_list5(p, &p.from), Dconv_list5(p, &p.to))
		}
	} else if a == int(ADATA_8) {
		return fmt.Sprintf("%.5lld (%L)	%s	%s/%d,%s", p.pc, p.lineno, Aconv_list5(a), Dconv_list5(p, &p.from), p.reg, Dconv_list5(p, &p.to))
	} else if p.as == int(ATEXT_8) {
		return fmt.Sprintf("%.5lld (%L)	%s	%s,%d,%s", p.pc, p.lineno, Aconv_list5(a), Dconv_list5(p, &p.from), p.reg, Dconv_list5(p, &p.to))
	} else if p.reg == int(NREG_5) {
		return fmt.Sprintf("%.5lld (%L)	%s%s	%s,%s", p.pc, p.lineno, Aconv_list5(a), sc, Dconv_list5(p, &p.from), Dconv_list5(p, &p.to))
	} else if p.from.typ != int(D_FREG_5) {
		return fmt.Sprintf("%.5lld (%L)	%s%s	%s,R%d,%s", p.pc, p.lineno, Aconv_list5(a), sc, Dconv_list5(p, &p.from), p.reg, Dconv_list5(p, &p.to))
	} else {
		return fmt.Sprintf("%.5lld (%L)	%s%s	%s,F%d,%s", p.pc, p.lineno, Aconv_list5(a), sc, Dconv_list5(p, &p.from), p.reg, Dconv_list5(p, &p.to))
	}
}

func Rconv_list5(r int) string {
	return fmt.Sprintf("R%d", r)
}

func RAconv_list5(a *Addr) string {
	var i int
	var v int
	str := "GOK-reglist"
	switch a.typ {
	case D_CONST_8:
	case D_CONST2_8:
		if a.reg != int(NREG_5) {
			break
		}
		if a.sym != nil {
			break
		}
		v = int(a.offset)
		str = ""
		for i = 0; i < int(NREG_5); i++ {
			if v&(1<<uint(i)) != 0 /*untyped*/ {
				if str == "" {
					str += "[R"
				} else {
					str += ",R"
				}
				str += fmt.Sprintf("%d", i)
			}
		}
		str += "]"
	}
	return str
}

var extra_list5 = []string{
	".EQ",
	".NE",
	".CS",
	".CC",
	".MI",
	".PL",
	".VS",
	".VC",
	".HI",
	".LS",
	".GE",
	".LT",
	".GT",
	".LE",
	"",
	".NV",
}
