package main

import (
	"fmt"
	"math"
)

var linkarm LinkArch

// Derived from Inferno utils/5c/swt.c
// http://code.google.com/p/inferno-os/source/browse/utils/5c/swt.c
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
var zprg_obj5 = Prog{
	as:    AGOK_8,
	scond: C_SCOND_NONE_5,
	reg:   NREG_5,
	from: Addr{
		name: D_NONE_8,
		typ:  D_NONE_8,
		reg:  NREG_5,
	},
	to: Addr{
		name: D_NONE_8,
		typ:  D_NONE_8,
		reg:  NREG_5,
	},
}

func symtype_obj5(a *Addr) int {
	return a.name
}

func isdata_obj5(p *Prog) bool {
	return p.as == int(ADATA_8) || p.as == int(AGLOBL_8)
}

func iscall_obj5(p *Prog) bool {
	return p.as == int(ABL_5)
}

func datasize_obj5(p *Prog) int {
	return p.reg
}

func textflag_obj5(p *Prog) int {
	return p.reg
}

func settextflag_obj5(p *Prog, f int) {
	p.reg = f
}

func progedit_obj5(ctxt *Link, p *Prog) {
	var literal string
	var s *LSym
	var tlsfallback *LSym
	p.from.class = 0
	p.to.class = 0
	// Rewrite B/BL to symbol as D_BRANCH.
	switch p.as {
	case AB_5:
	case ABL_5:
	case ADUFFZERO_8:
	case ADUFFCOPY_8:
		if p.to.typ == int(D_OREG_5) && (p.to.name == int(D_EXTERN_8) || p.to.name == int(D_STATIC_8)) && p.to.sym != nil {
			p.to.typ = int(D_BRANCH_8)
		}
		break
	}
	// Replace TLS register fetches on older ARM procesors.
	switch p.as {
	// If the instruction matches MRC 15, 0, <reg>, C13, C0, 3, replace it.
	case AMRC_5:
		if ctxt.goarm < 7 && (p.to.offset&0xffff0fff) == 0xee1d0f70 {
			tlsfallback = linklookup(ctxt, "runtime.read_tls_fallback", 0)
			// BL runtime.read_tls_fallback(SB)
			p.as = int(ABL_5)
			p.to.typ = int(D_BRANCH_8)
			p.to.sym = tlsfallback
			p.to.offset = 0
		} else {
			// Otherwise, MRC/MCR instructions need no further treatment.
			p.as = int(AWORD_8)
		}
		break
	}
	// Rewrite float constants to values stored in memory.
	switch p.as {
	case AMOVF_5:
		if p.from.typ == int(D_FCONST_8) && chipfloat5(ctxt, p.from.u.dval) < 0 && (chipzero5(ctxt, p.from.u.dval) < 0 || (p.scond&int(C_SCOND_5)) != int(C_SCOND_NONE_5)) {
			var i32 int32
			var f32 float32
			f32 = float32(p.from.u.dval)
			i32 = int32(math.Float32bits(f32))
			literal = fmt.Sprintf("$f32.%08ux", uint32(i32))
			s = linklookup(ctxt, string(literal), 0)
			if s.typ == 0 {
				s.typ = int(SRODATA)
				adduint32(ctxt, s, uint32(i32))
				s.reachable = 0
			}
			p.from.typ = int(D_OREG_5)
			p.from.sym = s
			p.from.name = int(D_EXTERN_8)
			p.from.offset = 0
		}
		break
	case AMOVD_5:
		if p.from.typ == int(D_FCONST_8) && chipfloat5(ctxt, p.from.u.dval) < 0 && (chipzero5(ctxt, p.from.u.dval) < 0 || (p.scond&int(C_SCOND_5)) != int(C_SCOND_NONE_5)) {
			var i64 int64
			i64 = int64(math.Float64bits(p.from.u.dval))
			literal = fmt.Sprintf("$f64.%016llux", uint64(i64))
			s = linklookup(ctxt, string(literal), 0)
			if s.typ == 0 {
				s.typ = int(SRODATA)
				adduint64(ctxt, s, uint64(i64))
				s.reachable = 0
			}
			p.from.typ = int(D_OREG_5)
			p.from.sym = s
			p.from.name = int(D_EXTERN_8)
			p.from.offset = 0
		}
		break
	}
	if ctxt.flag_shared != 0 {
		// Shared libraries use R_ARM_TLS_IE32 instead of
		// R_ARM_TLS_LE32, replacing the link time constant TLS offset in
		// runtime.tlsg with an address to a GOT entry containing the
		// offset. Rewrite $runtime.tlsg(SB) to runtime.tlsg(SB) to
		// compensate.
		if ctxt.tlsg == nil {
			ctxt.tlsg = linklookup(ctxt, "runtime.tlsg", 0)
		}
		if p.from.typ == int(D_CONST_8) && p.from.name == int(D_EXTERN_8) && p.from.sym == ctxt.tlsg {
			p.from.typ = int(D_OREG_5)
		}
		if p.to.typ == int(D_CONST_8) && p.to.name == int(D_EXTERN_8) && p.to.sym == ctxt.tlsg {
			p.to.typ = int(D_OREG_5)
		}
	}
}

func prg_obj5() *Prog {
	var p *Prog
	p = new(Prog)
	*p = zprg_obj5
	return p
}

func stacksplit_obj5(ctxt *Link, p *Prog, framesize int, noctxt bool) *Prog {
	var arg int
	// MOVW			g_stackguard(g), R1
	p = appendp(ctxt, p)
	p.as = int(AMOVW_8)
	p.from.typ = int(D_OREG_5)
	p.from.reg = int(REGG_5)
	p.to.typ = int(D_REG_5)
	p.to.reg = 1
	if framesize <= int(StackSmall_stack) {
		// small stack: SP < stackguard
		//	CMP	stackguard, SP
		p = appendp(ctxt, p)
		p.as = int(ACMP_5)
		p.from.typ = int(D_REG_5)
		p.from.reg = 1
		p.reg = int(REGSP_8)
	} else {
		if framesize <= int(StackBig_stack) {
			// large stack: SP-framesize < stackguard-StackSmall
			//	MOVW $-framesize(SP), R2
			//	CMP stackguard, R2
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.from.typ = int(D_CONST_8)
			p.from.reg = int(REGSP_8)
			p.from.offset = int64(-framesize)
			p.to.typ = int(D_REG_5)
			p.to.reg = 2
			p = appendp(ctxt, p)
			p.as = int(ACMP_5)
			p.from.typ = int(D_REG_5)
			p.from.reg = 1
			p.reg = 2
		} else {
			// Such a large stack we need to protect against wraparound
			// if SP is close to zero.
			//	SP-stackguard+StackGuard < framesize + (StackGuard-StackSmall)
			// The +StackGuard on both sides is required to keep the left side positive:
			// SP is allowed to be slightly below stackguard. See stack.h.
			//	CMP $StackPreempt, R1
			//	MOVW.NE $StackGuard(SP), R2
			//	SUB.NE R1, R2
			//	MOVW.NE $(framesize+(StackGuard-StackSmall)), R3
			//	CMP.NE R3, R2
			p = appendp(ctxt, p)
			p.as = int(ACMP_5)
			p.from.typ = int(D_CONST_8)
			p.from.offset = int64(uint32(StackPreempt_stack & 0xFFFFFFFF))
			p.reg = 1
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.from.typ = int(D_CONST_8)
			p.from.reg = int(REGSP_8)
			p.from.offset = int64(StackGuard_stack)
			p.to.typ = int(D_REG_5)
			p.to.reg = 2
			p.scond = int(C_SCOND_NE_5)
			p = appendp(ctxt, p)
			p.as = int(ASUB_5)
			p.from.typ = int(D_REG_5)
			p.from.reg = 1
			p.to.typ = int(D_REG_5)
			p.to.reg = 2
			p.scond = int(C_SCOND_NE_5)
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.from.typ = int(D_CONST_8)
			p.from.offset = int64(framesize) + (int64(StackGuard_stack) - int64(StackSmall_stack))
			p.to.typ = int(D_REG_5)
			p.to.reg = 3
			p.scond = int(C_SCOND_NE_5)
			p = appendp(ctxt, p)
			p.as = int(ACMP_5)
			p.from.typ = int(D_REG_5)
			p.from.reg = 3
			p.reg = 2
			p.scond = int(C_SCOND_NE_5)
		}
	}
	// MOVW.LS		$framesize, R1
	p = appendp(ctxt, p)
	p.as = int(AMOVW_8)
	p.scond = int(C_SCOND_LS_5)
	p.from.typ = int(D_CONST_8)
	p.from.offset = int64(framesize)
	p.to.typ = int(D_REG_5)
	p.to.reg = 1
	// MOVW.LS		$args, R2
	p = appendp(ctxt, p)
	p.as = int(AMOVW_8)
	p.scond = int(C_SCOND_LS_5)
	p.from.typ = int(D_CONST_8)
	arg = ctxt.cursym.text.to.offset2
	if arg == 1 { // special marker for known 0
		arg = 0
	}
	if arg&3 != 0 /*untyped*/ {
		ctxt.diag("misaligned argument size in stack split")
	}
	p.from.offset = int64(arg)
	p.to.typ = int(D_REG_5)
	p.to.reg = 2
	// MOVW.LS	R14, R3
	p = appendp(ctxt, p)
	p.as = int(AMOVW_8)
	p.scond = int(C_SCOND_LS_5)
	p.from.typ = int(D_REG_5)
	p.from.reg = int(REGLINK_5)
	p.to.typ = int(D_REG_5)
	p.to.reg = 3
	// BL.LS		runtime.morestack(SB) // modifies LR, returns with LO still asserted
	p = appendp(ctxt, p)
	p.as = int(ABL_5)
	p.scond = int(C_SCOND_LS_5)
	p.to.typ = int(D_BRANCH_8)
	p.to.sym = ctxt.symmorestack[bool2int(noctxt)]
	// BLS	start
	p = appendp(ctxt, p)
	p.as = int(ABLS_5)
	p.to.typ = int(D_BRANCH_8)
	p.pcond = ctxt.cursym.text.link
	return p
}

func initdiv_obj5(ctxt *Link) {
	if ctxt.sym_div != nil {
		return
	}
	ctxt.sym_div = linklookup(ctxt, "_div", 0)
	ctxt.sym_divu = linklookup(ctxt, "_divu", 0)
	ctxt.sym_mod = linklookup(ctxt, "_mod", 0)
	ctxt.sym_modu = linklookup(ctxt, "_modu", 0)
}

func softfloat_obj5(ctxt *Link, cursym *LSym) {
	var p *Prog
	var next *Prog
	var symsfloat *LSym
	var wasfloat int
	if ctxt.goarm > 5 {
		return
	}
	symsfloat = linklookup(ctxt, "_sfloat", 0)
	wasfloat = 0
	for p = cursym.text; p != nil; p = p.link {
		if p.pcond != nil {
			p.pcond.mark |= int(LABEL_obj5)
		}
	}
	for p = cursym.text; p != nil; p = p.link {
		switch p.as {
		case AMOVW_8:
			if p.to.typ == int(D_FREG_5) || p.from.typ == int(D_FREG_5) {
				break
			}
			wasfloat = 0
			continue
		case AMOVWD_5:
		case AMOVWF_5:
		case AMOVDW_5:
		case AMOVFW_5:
		case AMOVFD_5:
		case AMOVDF_5:
		case AMOVF_5:
		case AMOVD_5:
		case ACMPF_5:
		case ACMPD_5:
		case AADDF_5:
		case AADDD_5:
		case ASUBF_5:
		case ASUBD_5:
		case AMULF_5:
		case AMULD_5:
		case ADIVF_5:
		case ADIVD_5:
		case ASQRTF_5:
		case ASQRTD_5:
		case AABSF_5:
		case AABSD_5:
			break
		default:
			wasfloat = 0
			continue
		}
		if !(wasfloat != 0) || (p.mark&int(LABEL_obj5) != 0) {
			next = ctxt.arch.prg()
			*next = *p
			// BL _sfloat(SB)
			*p = zprg_obj5
			p.link = next
			p.as = int(ABL_5)
			p.to.typ = int(D_BRANCH_8)
			p.to.sym = symsfloat
			p.lineno = next.lineno
			p = next
			wasfloat = 1
		}
	}
}

// Prog.mark
const (
	FOLL_obj5  = 1 << 0
	LABEL_obj5 = 1 << 1
	LEAF_obj5  = 1 << 2
)

func linkcase_obj5(casep *Prog) {
	var p *Prog
	for p = casep; p != nil; p = p.link {
		if p.as == int(ABCASE_5) {
			for ; p != nil && p.as == int(ABCASE_5); p = p.link {
				p.pcrel = casep
			}
			break
		}
	}
}

func nocache_obj5(p *Prog) {
	p.optab = 0
	p.from.class = 0
	p.to.class = 0
}

func addstacksplit_obj5(ctxt *Link, cursym *LSym) {
	var p *Prog
	var pl *Prog
	var q *Prog
	var q1 *Prog
	var q2 *Prog
	var o int
	var autosize int
	var autoffset int
	autosize = 0
	if ctxt.symmorestack[0] == nil {
		ctxt.symmorestack[0] = linklookup(ctxt, "runtime.morestack", 0)
		ctxt.symmorestack[1] = linklookup(ctxt, "runtime.morestack_noctxt", 0)
	}
	q = (*Prog)(nil)
	ctxt.cursym = cursym
	if cursym.text == nil || cursym.text.link == nil {
		return
	}
	softfloat_obj5(ctxt, cursym)
	p = cursym.text
	autoffset = int(p.to.offset)
	if autoffset < 0 {
		autoffset = 0
	}
	cursym.locals = autoffset
	cursym.args = p.to.offset2
	if ctxt.debugzerostack != 0 {
		if autoffset != 0 && !(p.reg&int(NOSPLIT_textflag) != 0) {
			// MOVW $4(R13), R1
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.from.typ = int(D_CONST_8)
			p.from.reg = 13
			p.from.offset = 4
			p.to.typ = int(D_REG_5)
			p.to.reg = 1
			// MOVW $n(R13), R2
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.from.typ = int(D_CONST_8)
			p.from.reg = 13
			p.from.offset = 4 + int64(autoffset)
			p.to.typ = int(D_REG_5)
			p.to.reg = 2
			// MOVW $0, R3
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.from.typ = int(D_CONST_8)
			p.from.offset = 0
			p.to.typ = int(D_REG_5)
			p.to.reg = 3
			// L:
			//	MOVW.nil R3, 0(R1) +4
			//	CMP R1, R2
			//	BNE L
			pl = appendp(ctxt, p)
			p = pl
			p.as = int(AMOVW_8)
			p.from.typ = int(D_REG_5)
			p.from.reg = 3
			p.to.typ = int(D_OREG_5)
			p.to.reg = 1
			p.to.offset = 4
			p.scond |= int(C_PBIT_5)
			p = appendp(ctxt, p)
			p.as = int(ACMP_5)
			p.from.typ = int(D_REG_5)
			p.from.reg = 1
			p.reg = 2
			p = appendp(ctxt, p)
			p.as = int(ABNE_5)
			p.to.typ = int(D_BRANCH_8)
			p.pcond = pl
		}
	}
	/*
	 * find leaf subroutines
	 * strip NOPs
	 * expand RET
	 * expand BECOME pseudo
	 */
	for p = cursym.text; p != nil; p = p.link {
		switch p.as {
		case ACASE_5:
			if ctxt.flag_shared != 0 {
				linkcase_obj5(p)
			}
			break
		case ATEXT_8:
			p.mark |= int(LEAF_obj5)
			break
		case ARET_8:
			break
		case ADIV_5:
		case ADIVU_5:
		case AMOD_5:
		case AMODU_5:
			q = p
			if ctxt.sym_div == nil {
				initdiv_obj5(ctxt)
			}
			cursym.text.mark &^= int(LEAF_obj5)
			continue
		case ANOP_8:
			q1 = p.link
			q.link = q1 /* q is non-nop */
			if q1 != nil {
				q1.mark |= p.mark
			}
			continue
		case ABL_5:
		case ABX_5:
		case ADUFFZERO_8:
		case ADUFFCOPY_8:
			cursym.text.mark &^= int(LEAF_obj5)
		case ABCASE_5:
		case AB_5:
		case ABEQ_5:
		case ABNE_5:
		case ABCS_5:
		case ABHS_5:
		case ABCC_5:
		case ABLO_5:
		case ABMI_5:
		case ABPL_5:
		case ABVS_5:
		case ABVC_5:
		case ABHI_5:
		case ABLS_5:
		case ABGE_5:
		case ABLT_5:
		case ABGT_5:
		case ABLE_5:
			q1 = p.pcond
			if q1 != nil {
				for q1.as == int(ANOP_8) {
					q1 = q1.link
					p.pcond = q1
				}
			}
			break
		}
		q = p
	}
	for p = cursym.text; p != nil; p = p.link {
		o = p.as
		switch o {
		case ATEXT_8:
			autosize = int(p.to.offset + 4)
			if autosize <= 4 {
				if cursym.text.mark&int(LEAF_obj5) != 0 {
					p.to.offset = -4
					autosize = 0
				}
			}
			if !(autosize != 0) && !(cursym.text.mark&int(LEAF_obj5) != 0) {
				if ctxt.debugvlog != 0 {
					Bprint(ctxt.bso, "save suppressed in: %s\n", cursym.name)
					Bflush(ctxt.bso)
				}
				cursym.text.mark |= int(LEAF_obj5)
			}
			if cursym.text.mark&int(LEAF_obj5) != 0 {
				cursym.leaf = 1
				if !(autosize != 0) {
					break
				}
			}
			if !(p.reg&int(NOSPLIT_textflag) != 0) {
				p = stacksplit_obj5(ctxt, p, autosize, !(cursym.text.reg&int(NEEDCTXT_textflag) != 0)) // emit split check
			}
			// MOVW.W		R14,$-autosize(SP)
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.scond |= int(C_WBIT_5)
			p.from.typ = int(D_REG_5)
			p.from.reg = int(REGLINK_5)
			p.to.typ = int(D_OREG_5)
			p.to.offset = int64(-autosize)
			p.to.reg = int(REGSP_8)
			p.spadj = autosize
			if cursym.text.reg&int(WRAPPER_textflag) != 0 {
				// g->panicwrap += autosize;
				// MOVW panicwrap_offset(g), R3
				// ADD $autosize, R3
				// MOVW R3 panicwrap_offset(g)
				p = appendp(ctxt, p)
				p.as = int(AMOVW_8)
				p.from.typ = int(D_OREG_5)
				p.from.reg = int(REGG_5)
				p.from.offset = 2 * int64(ctxt.arch.ptrsize)
				p.to.typ = int(D_REG_5)
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = int(AADD_5)
				p.from.typ = int(D_CONST_8)
				p.from.offset = int64(autosize)
				p.to.typ = int(D_REG_5)
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = int(AMOVW_8)
				p.from.typ = int(D_REG_5)
				p.from.reg = 3
				p.to.typ = int(D_OREG_5)
				p.to.reg = int(REGG_5)
				p.to.offset = 2 * int64(ctxt.arch.ptrsize)
			}
			break
		case ARET_8:
			nocache_obj5(p)
			if cursym.text.mark&int(LEAF_obj5) != 0 {
				if !(autosize != 0) {
					p.as = int(AB_5)
					p.from = zprg_obj5.from
					if p.to.sym != nil { // retjmp
						p.to.typ = int(D_BRANCH_8)
					} else {
						p.to.typ = int(D_OREG_5)
						p.to.offset = 0
						p.to.reg = int(REGLINK_5)
					}
					break
				}
			}
			if cursym.text.reg&int(WRAPPER_textflag) != 0 {
				var scond int
				// Preserve original RET's cond, to allow RET.EQ
				// in the implementation of reflect.call.
				scond = p.scond
				p.scond = int(C_SCOND_NONE_5)
				// g->panicwrap -= autosize;
				// MOVW panicwrap_offset(g), R3
				// SUB $autosize, R3
				// MOVW R3 panicwrap_offset(g)
				p.as = int(AMOVW_8)
				p.from.typ = int(D_OREG_5)
				p.from.reg = int(REGG_5)
				p.from.offset = 2 * int64(ctxt.arch.ptrsize)
				p.to.typ = int(D_REG_5)
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = int(ASUB_5)
				p.from.typ = int(D_CONST_8)
				p.from.offset = int64(autosize)
				p.to.typ = int(D_REG_5)
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = int(AMOVW_8)
				p.from.typ = int(D_REG_5)
				p.from.reg = 3
				p.to.typ = int(D_OREG_5)
				p.to.reg = int(REGG_5)
				p.to.offset = 2 * int64(ctxt.arch.ptrsize)
				p = appendp(ctxt, p)
				p.scond = scond
			}
			p.as = int(AMOVW_8)
			p.scond |= int(C_PBIT_5)
			p.from.typ = int(D_OREG_5)
			p.from.offset = int64(autosize)
			p.from.reg = int(REGSP_8)
			p.to.typ = int(D_REG_5)
			p.to.reg = int(REGPC_5)
			// If there are instructions following
			// this ARET, they come from a branch
			// with the same stackframe, so no spadj.
			if p.to.sym != nil { // retjmp
				p.to.reg = int(REGLINK_5)
				q2 = appendp(ctxt, p)
				q2.as = int(AB_5)
				q2.to.typ = int(D_BRANCH_8)
				q2.to.sym = p.to.sym
				p.to.sym = (*LSym)(nil)
				p = q2
			}
			break
		case AADD_5:
			if p.from.typ == int(D_CONST_8) && p.from.reg == int(NREG_5) && p.to.typ == int(D_REG_5) && p.to.reg == int(REGSP_8) {
				p.spadj = int(-p.from.offset)
			}
			break
		case ASUB_5:
			if p.from.typ == int(D_CONST_8) && p.from.reg == int(NREG_5) && p.to.typ == int(D_REG_5) && p.to.reg == int(REGSP_8) {
				p.spadj = int(p.from.offset)
			}
			break
		case ADIV_5:
		case ADIVU_5:
		case AMOD_5:
		case AMODU_5:
			if ctxt.debugdivmod != 0 {
				break
			}
			if p.from.typ != int(D_REG_5) {
				break
			}
			if p.to.typ != int(D_REG_5) {
				break
			}
			q1 = p
			/* MOV a,4(SP) */
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.lineno = q1.lineno
			p.from.typ = int(D_REG_5)
			p.from.reg = q1.from.reg
			p.to.typ = int(D_OREG_5)
			p.to.reg = int(REGSP_8)
			p.to.offset = 4
			/* MOV b,REGTMP */
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.lineno = q1.lineno
			p.from.typ = int(D_REG_5)
			p.from.reg = q1.reg
			if q1.reg == int(NREG_5) {
				p.from.reg = q1.to.reg
			}
			p.to.typ = int(D_REG_5)
			p.to.reg = int(REGTMP_8)
			p.to.offset = 0
			/* CALL appropriate */
			p = appendp(ctxt, p)
			p.as = int(ABL_5)
			p.lineno = q1.lineno
			p.to.typ = int(D_BRANCH_8)
			switch o {
			case ADIV_5:
				p.to.sym = ctxt.sym_div
				break
			case ADIVU_5:
				p.to.sym = ctxt.sym_divu
				break
			case AMOD_5:
				p.to.sym = ctxt.sym_mod
				break
			case AMODU_5:
				p.to.sym = ctxt.sym_modu
				break
			}
			/* MOV REGTMP, b */
			p = appendp(ctxt, p)
			p.as = int(AMOVW_8)
			p.lineno = q1.lineno
			p.from.typ = int(D_REG_5)
			p.from.reg = int(REGTMP_8)
			p.from.offset = 0
			p.to.typ = int(D_REG_5)
			p.to.reg = q1.to.reg
			/* ADD $8,SP */
			p = appendp(ctxt, p)
			p.as = int(AADD_5)
			p.lineno = q1.lineno
			p.from.typ = int(D_CONST_8)
			p.from.reg = int(NREG_5)
			p.from.offset = 8
			p.reg = int(NREG_5)
			p.to.typ = int(D_REG_5)
			p.to.reg = int(REGSP_8)
			p.spadj = -8
			/* Keep saved LR at 0(SP) after SP change. */
			/* MOVW 0(SP), REGTMP; MOVW REGTMP, -8!(SP) */
			/* TODO: Remove SP adjustments; see issue 6699. */
			q1.as = int(AMOVW_8)
			q1.from.typ = int(D_OREG_5)
			q1.from.reg = int(REGSP_8)
			q1.from.offset = 0
			q1.reg = int(NREG_5)
			q1.to.typ = int(D_REG_5)
			q1.to.reg = int(REGTMP_8)
			/* SUB $8,SP */
			q1 = appendp(ctxt, q1)
			q1.as = int(AMOVW_8)
			q1.from.typ = int(D_REG_5)
			q1.from.reg = int(REGTMP_8)
			q1.reg = int(NREG_5)
			q1.to.typ = int(D_OREG_5)
			q1.to.reg = int(REGSP_8)
			q1.to.offset = -8
			q1.scond |= int(C_WBIT_5)
			q1.spadj = 8
			break
		case AMOVW_8:
			if (p.scond&int(C_WBIT_5) != 0) && p.to.typ == int(D_OREG_5) && p.to.reg == int(REGSP_8) {
				p.spadj = int(-p.to.offset)
			}
			if (p.scond&int(C_PBIT_5) != 0) && p.from.typ == int(D_OREG_5) && p.from.reg == int(REGSP_8) && p.to.reg != int(REGPC_5) {
				p.spadj = int(-p.from.offset)
			}
			if p.from.typ == int(D_CONST_8) && p.from.reg == int(REGSP_8) && p.to.typ == int(D_REG_5) && p.to.reg == int(REGSP_8) {
				p.spadj = int(-p.from.offset)
			}
			break
		}
	}
}

func xfol_obj5(ctxt *Link, p *Prog, last **Prog) {
	var q *Prog
	var r *Prog
	var a int
	var i int
loop:
	if p == nil {
		return
	}
	a = p.as
	if a == int(AB_5) {
		q = p.pcond
		if q != nil && q.as != int(ATEXT_8) {
			p.mark |= int(FOLL_obj5)
			p = q
			if !(p.mark&int(FOLL_obj5) != 0) {
				goto loop
			}
		}
	}
	if p.mark&int(FOLL_obj5) != 0 {
		i = 0
		q = p
		for ; i < 4; (func() { i++; q = q.link })() {
			if q == *last || q == nil {
				break
			}
			a = q.as
			if a == int(ANOP_8) {
				i--
				continue
			}
			if a == int(AB_5) || (a == int(ARET_8) && q.scond == int(C_SCOND_NONE_5)) || a == int(ARFE_5) || a == int(AUNDEF_8) {
				goto copy
			}
			if q.pcond == nil || (q.pcond.mark&int(FOLL_obj5) != 0) {
				continue
			}
			if a != int(ABEQ_5) && a != int(ABNE_5) {
				continue
			}
		copy:
			for {
				r = ctxt.arch.prg()
				*r = *p
				if !(r.mark&int(FOLL_obj5) != 0) {
					print("can't happen 1\n")
				}
				r.mark |= int(FOLL_obj5)
				if p != q {
					p = p.link
					(*last).link = r
					*last = r
					continue
				}
				(*last).link = r
				*last = r
				if a == int(AB_5) || (a == int(ARET_8) && q.scond == int(C_SCOND_NONE_5)) || a == int(ARFE_5) || a == int(AUNDEF_8) {
					return
				}
				r.as = int(ABNE_5)
				if a == int(ABNE_5) {
					r.as = int(ABEQ_5)
				}
				r.pcond = p.link
				r.link = p.pcond
				if !(r.link.mark&int(FOLL_obj5) != 0) {
					xfol_obj5(ctxt, r.link, last)
				}
				if !(r.pcond.mark&int(FOLL_obj5) != 0) {
					print("can't happen 2\n")
				}
				return
			}
		}
		a = int(AB_5)
		q = ctxt.arch.prg()
		q.as = a
		q.lineno = p.lineno
		q.to.typ = int(D_BRANCH_8)
		q.to.offset = int64(p.pc)
		q.pcond = p
		p = q
	}
	p.mark |= int(FOLL_obj5)
	(*last).link = p
	*last = p
	if a == int(AB_5) || (a == int(ARET_8) && p.scond == int(C_SCOND_NONE_5)) || a == int(ARFE_5) || a == int(AUNDEF_8) {
		return
	}
	if p.pcond != nil {
		if a != int(ABL_5) && a != int(ABX_5) && p.link != nil {
			q = brchain(ctxt, p.link)
			if a != int(ATEXT_8) && a != int(ABCASE_5) {
				if q != nil && (q.mark&int(FOLL_obj5) != 0) {
					p.as = relinv_obj5(a)
					p.link = p.pcond
					p.pcond = q
				}
			}
			xfol_obj5(ctxt, p.link, last)
			q = brchain(ctxt, p.pcond)
			if q == nil {
				q = p.pcond
			}
			if q.mark&int(FOLL_obj5) != 0 {
				p.pcond = q
				return
			}
			p = q
			goto loop
		}
	}
	p = p.link
	goto loop
}

func follow_obj5(ctxt *Link, s *LSym) {
	var firstp *Prog
	var lastp *Prog
	ctxt.cursym = s
	firstp = ctxt.arch.prg()
	lastp = firstp
	xfol_obj5(ctxt, s.text, &lastp)
	lastp.link = (*Prog)(nil)
	s.text = firstp.link
}

func relinv_obj5(a int) int {
	switch a {
	case ABEQ_5:
		return int(ABNE_5)
	case ABNE_5:
		return int(ABEQ_5)
	case ABCS_5:
		return int(ABCC_5)
	case ABHS_5:
		return int(ABLO_5)
	case ABCC_5:
		return int(ABCS_5)
	case ABLO_5:
		return int(ABHS_5)
	case ABMI_5:
		return int(ABPL_5)
	case ABPL_5:
		return int(ABMI_5)
	case ABVS_5:
		return int(ABVC_5)
	case ABVC_5:
		return int(ABVS_5)
	case ABHI_5:
		return int(ABLS_5)
	case ABLS_5:
		return int(ABHI_5)
	case ABGE_5:
		return int(ABLT_5)
	case ABLT_5:
		return int(ABGE_5)
	case ABGT_5:
		return int(ABLE_5)
	case ABLE_5:
		return int(ABGT_5)
	}
	sysfatal("unknown relation: %s", anames5[a])
	return 0
}
