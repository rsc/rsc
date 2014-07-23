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
const (
	STRINGSZ_list5 = 1000
)

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

func Pconv_list5(p *Prog) string {
	var str string
	var sc string
	var fp string
	var a int
	var s int

	a = p.as
	s = p.scond
	sc = extra_list5[s&C_SCOND_5]
	if s&C_SBIT_5 != 0 {
		sc += ".S"
	}
	if s&C_PBIT_5 != 0 {
		sc += ".P"
	}
	if s&C_WBIT_5 != 0 {
		sc += ".W"
	}
	if s&C_UBIT_5 != 0 { /* ambiguous with FBIT */
		sc += ".U"
	}
	if a == AMOVM_5 {
		if p.from.typ == D_CONST_5 {
			str = fmt.Sprintf("%.5d (%v)\t%v%s\t%v,%v", p.pc, p.Line(), Aconv_list5(a), sc, RAconv_list5(&p.from), Dconv_list5(p, 0, &p.to))
		} else if p.to.typ == D_CONST_5 {
			str = fmt.Sprintf("%.5d (%v)\t%v%s\t%v,%v", p.pc, p.Line(), Aconv_list5(a), sc, Dconv_list5(p, 0, &p.from), RAconv_list5(&p.to))
		} else {
			str = fmt.Sprintf("%.5d (%v)\t%v%s\t%v,%v", p.pc, p.Line(), Aconv_list5(a), sc, Dconv_list5(p, 0, &p.from), Dconv_list5(p, 0, &p.to))
		}
	} else if a == ADATA_5 {
		str = fmt.Sprintf("%.5d (%v)\t%v\t%v/%d,%v", p.pc, p.Line(), Aconv_list5(a), Dconv_list5(p, 0, &p.from), p.reg, Dconv_list5(p, 0, &p.to))
	} else if p.as == ATEXT_5 {
		str = fmt.Sprintf("%.5d (%v)\t%v\t%v,%d,%v", p.pc, p.Line(), Aconv_list5(a), Dconv_list5(p, 0, &p.from), p.reg, Dconv_list5(p, 0, &p.to))
	} else if p.reg == NREG_5 {
		str = fmt.Sprintf("%.5d (%v)\t%v%s\t%v,%v", p.pc, p.Line(), Aconv_list5(a), sc, Dconv_list5(p, 0, &p.from), Dconv_list5(p, 0, &p.to))
	} else if p.from.typ != D_FREG_5 {
		str = fmt.Sprintf("%.5d (%v)\t%v%s\t%v,R%d,%v", p.pc, p.Line(), Aconv_list5(a), sc, Dconv_list5(p, 0, &p.from), p.reg, Dconv_list5(p, 0, &p.to))
	} else {
		str = fmt.Sprintf("%.5d (%v)\t%v%s\t%v,F%d,%v", p.pc, p.Line(), Aconv_list5(a), sc, Dconv_list5(p, 0, &p.from), p.reg, Dconv_list5(p, 0, &p.to))
	}

	fp += str
	return fp
}

func Aconv_list5(a int) string {
	var s string
	var fp string

	s = "???"
	if a >= AXXX_5 && a < ALAST_5 {
		s = anames5[a]
	}
	fp += s
	return fp
}

func Dconv_list5(p *Prog, flag int, a *Addr) string {
	var str string
	var fp string
	var op string
	var v int

	switch a.typ {
	default:
		str = fmt.Sprintf("GOK-type(%d)", a.typ)
	case D_NONE_5:
		str = ""
		if a.name != D_NONE_5 || a.reg != NREG_5 || a.sym != nil {
			str = fmt.Sprintf("%v(R%d)(NONE)", Mconv_list5(a), a.reg)
		}
	case D_CONST_5:
		if a.reg != NREG_5 {
			str = fmt.Sprintf("$%v(R%d)", Mconv_list5(a), a.reg)
		} else {
			str = fmt.Sprintf("$%v", Mconv_list5(a))
		}
	case D_CONST2_5:
		str = fmt.Sprintf("$%d-%d", a.offset, a.offset2)
	case D_SHIFT_5:
		v = int(a.offset)
		op = string("<<>>->@>"[((v>>5)&3)<<1:])
		if v&(1<<4) != 0 {
			str = fmt.Sprintf("R%d%c%cR%d", v&15, op[0], op[1], (v>>8)&15)
		} else {
			str = fmt.Sprintf("R%d%c%c%d", v&15, op[0], op[1], (v>>7)&31)
		}
		if a.reg != NREG_5 {
			str += fmt.Sprintf("(R%d)", a.reg)
		}
	case D_OREG_5:
		if a.reg != NREG_5 {
			str = fmt.Sprintf("%v(R%d)", Mconv_list5(a), a.reg)
		} else {
			str = fmt.Sprintf("%v", Mconv_list5(a))
		}
	case D_REG_5:
		str = fmt.Sprintf("R%d", a.reg)
		if a.name != D_NONE_5 || a.sym != nil {
			str = fmt.Sprintf("%v(R%d)(REG)", Mconv_list5(a), a.reg)
		}
	case D_FREG_5:
		str = fmt.Sprintf("F%d", a.reg)
		if a.name != D_NONE_5 || a.sym != nil {
			str = fmt.Sprintf("%v(R%d)(REG)", Mconv_list5(a), a.reg)
		}
	case D_PSR_5:
		str = fmt.Sprintf("PSR")
		if a.name != D_NONE_5 || a.sym != nil {
			str = fmt.Sprintf("%v(PSR)(REG)", Mconv_list5(a))
		}
	case D_BRANCH_5:
		if a.sym != nil {
			str = fmt.Sprintf("%s(SB)", a.sym.name)
		} else if p != nil && p.pcond != nil {
			str = fmt.Sprintf("%d", p.pcond.pc)
		} else if a.u.branch != nil {
			str = fmt.Sprintf("%d", a.u.branch.pc)
		} else {
			str = fmt.Sprintf("%d(PC)", a.offset) /*-pc*/
		}
	case D_FCONST_5:
		str = fmt.Sprintf("$%.17g", a.u.dval)
	case D_SCONST_5:
		str = fmt.Sprintf("$\"%q\"", a.u.sval)
		break
	}
	fp += str
	return fp
}

func RAconv_list5(a *Addr) string {
	var str string
	var fp string
	var i int
	var v int

	str = fmt.Sprintf("GOK-reglist")
	switch a.typ {
	case D_CONST_5,
		D_CONST2_5:
		if a.reg != NREG_5 {
			break
		}
		if a.sym != nil {
			break
		}
		v = int(a.offset)
		str = ""
		for i = 0; i < NREG_5; i++ {
			if v&(1<<uint(i)) != 0 {
				if str[0] == 0 {
					str += "[R"
				} else {
					str += ",R"
				}
				str += fmt.Sprintf("%d", i)
			}
		}
		str += "]"
	}
	fp += str
	return fp
}

func Rconv_list5(r int) string {
	var fp string
	var str string

	str = fmt.Sprintf("R%d", r)
	fp += str
	return fp
}

func Mconv_list5(a *Addr) string {
	var str string
	var fp string
	var s *LSym

	s = a.sym
	if s == nil {
		str = fmt.Sprintf("%d", int32(a.offset))
		goto out
	}
	switch a.name {
	default:
		str = fmt.Sprintf("GOK-name(%d)", a.name)
	case D_NONE_5:
		str = fmt.Sprintf("%d", int32(a.offset))
	case D_EXTERN_5:
		str = fmt.Sprintf("%s+%d(SB)", s.name, int(a.offset))
	case D_STATIC_5:
		str = fmt.Sprintf("%s<>+%d(SB)", s.name, int(a.offset))
	case D_AUTO_5:
		str = fmt.Sprintf("%s-%d(SP)", s.name, int(-a.offset))
	case D_PARAM_5:
		str = fmt.Sprintf("%s+%d(FP)", s.name, int(a.offset))
		break
	}
out:
	fp += str
	return fp
}
