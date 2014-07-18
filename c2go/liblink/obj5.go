package main

import (
	"fmt"
	"log"
	"math"
)

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
	as:    AGOK_5,
	scond: C_SCOND_NONE_5,
	reg:   NREG_5,
	from: Addr{
		name: D_NONE_5,
		typ:  D_NONE_5,
		reg:  NREG_5,
	},
	to: Addr{
		name: D_NONE_5,
		typ:  D_NONE_5,
		reg:  NREG_5,
	},
}

func symtype_obj5(a *Addr) int {
	return int(a.name)
}

func isdata_obj5(p *Prog) int {
	return bool2int(p.as == ADATA_5 || p.as == AGLOBL_5)
}

func iscall_obj5(p *Prog) int {
	return bool2int(p.as == ABL_5)
}

func datasize_obj5(p *Prog) int {
	return int(p.reg)
}

func textflag_obj5(p *Prog) int {
	return int(p.reg)
}

func settextflag_obj5(p *Prog, f int) {
	p.reg = int64(f)
}

func progedit_obj5(ctxt *Link, p *Prog) {
	var literal string
	var s *LSym
	var tlsfallback *LSym
	p.from.class = 0
	p.to.class = 0
	// Rewrite B/BL to symbol as D_BRANCH.
	switch p.as {
	case AB_5,
		ABL_5,
		ADUFFZERO_5,
		ADUFFCOPY_5:
		if p.to.typ == D_OREG_5 && (p.to.name == D_EXTERN_5 || p.to.name == D_STATIC_5) && p.to.sym != nil {
			p.to.typ = D_BRANCH_5
		}
		break
	}
	// Replace TLS register fetches on older ARM procesors.
	switch p.as {
	// If the instruction matches MRC 15, 0, <reg>, C13, C0, 3, replace it.
	case AMRC_5:
		if ctxt.goarm < 7 && p.to.offset&0xffff0fff == 0xee1d0f70 {
			tlsfallback = linklookup(ctxt, "runtime.read_tls_fallback", 0)
			// BL runtime.read_tls_fallback(SB)
			p.as = ABL_5
			p.to.typ = D_BRANCH_5
			p.to.sym = tlsfallback
			p.to.offset = 0
		} else {
			// Otherwise, MRC/MCR instructions need no further treatment.
			p.as = AWORD_5
		}
		break
	}
	// Rewrite float constants to values stored in memory.
	switch p.as {
	case AMOVF_5:
		if p.from.typ == D_FCONST_5 && chipfloat5(ctxt, p.from.u.dval) < 0 && (chipzero5(ctxt, p.from.u.dval) < 0 || p.scond&C_SCOND_5 != C_SCOND_NONE_5) {
			var i32 uint64
			var f32 float64
			f32 = p.from.u.dval
			i32 = uint64(math.Float32bits(float32(f32)))
			literal = fmt.Sprintf("$f32.%08x", uint32(i32))
			s = linklookup(ctxt, literal, 0)
			if s.typ == 0 {
				s.typ = SRODATA
				adduint32(ctxt, s, i32)
				s.reachable = 0
			}
			p.from.typ = D_OREG_5
			p.from.sym = s
			p.from.name = D_EXTERN_5
			p.from.offset = 0
		}
	case AMOVD_5:
		if p.from.typ == D_FCONST_5 && chipfloat5(ctxt, p.from.u.dval) < 0 && (chipzero5(ctxt, p.from.u.dval) < 0 || p.scond&C_SCOND_5 != C_SCOND_NONE_5) {
			var i64 uint64
			i64 = math.Float64bits(p.from.u.dval)
			literal = fmt.Sprintf("$f64.%016x", uint64(i64))
			s = linklookup(ctxt, literal, 0)
			if s.typ == 0 {
				s.typ = SRODATA
				adduint64(ctxt, s, i64)
				s.reachable = 0
			}
			p.from.typ = D_OREG_5
			p.from.sym = s
			p.from.name = D_EXTERN_5
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
		if p.from.typ == D_CONST_5 && p.from.name == D_EXTERN_5 && p.from.sym == ctxt.tlsg {
			p.from.typ = D_OREG_5
		}
		if p.to.typ == D_CONST_5 && p.to.name == D_EXTERN_5 && p.to.sym == ctxt.tlsg {
			p.to.typ = D_OREG_5
		}
	}
}

func prg_obj5() *Prog {
	var p *Prog
	p = new(Prog)
	*p = zprg_obj5
	return p
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
		if p.as == ABCASE_5 {
			for ; p != nil && p.as == ABCASE_5; p = p.link {
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
	var autosize int64
	var autoffset int64
	autosize = 0
	if ctxt.symmorestack[0] == nil {
		ctxt.symmorestack[0] = linklookup(ctxt, "runtime.morestack", 0)
		ctxt.symmorestack[1] = linklookup(ctxt, "runtime.morestack_noctxt", 0)
	}
	q = nil
	ctxt.cursym = cursym
	if cursym.text == nil || cursym.text.link == nil {
		return
	}
	softfloat_obj5(ctxt, cursym)
	p = cursym.text
	autoffset = p.to.offset
	if autoffset < 0 {
		autoffset = 0
	}
	cursym.locals = autoffset
	cursym.args = p.to.offset2
	if ctxt.debugzerostack != 0 {
		if autoffset != 0 && !(p.reg&NOSPLIT_textflag != 0) {
			// MOVW $4(R13), R1
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.from.typ = D_CONST_5
			p.from.reg = 13
			p.from.offset = 4
			p.to.typ = D_REG_5
			p.to.reg = 1
			// MOVW $n(R13), R2
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.from.typ = D_CONST_5
			p.from.reg = 13
			p.from.offset = 4 + autoffset
			p.to.typ = D_REG_5
			p.to.reg = 2
			// MOVW $0, R3
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.from.typ = D_CONST_5
			p.from.offset = 0
			p.to.typ = D_REG_5
			p.to.reg = 3
			// L:
			//	MOVW.nil R3, 0(R1) +4
			//	CMP R1, R2
			//	BNE L
			pl = appendp(ctxt, p)
			p = pl
			p.as = AMOVW_5
			p.from.typ = D_REG_5
			p.from.reg = 3
			p.to.typ = D_OREG_5
			p.to.reg = 1
			p.to.offset = 4
			p.scond |= C_PBIT_5
			p = appendp(ctxt, p)
			p.as = ACMP_5
			p.from.typ = D_REG_5
			p.from.reg = 1
			p.reg = 2
			p = appendp(ctxt, p)
			p.as = ABNE_5
			p.to.typ = D_BRANCH_5
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
		case ATEXT_5:
			p.mark |= LEAF_obj5
		case ARET_5:
			break
		case ADIV_5,
			ADIVU_5,
			AMOD_5,
			AMODU_5:
			q = p
			if ctxt.sym_div == nil {
				initdiv_obj5(ctxt)
			}
			cursym.text.mark &^= LEAF_obj5
			continue
		case ANOP_5:
			q1 = p.link
			q.link = q1 /* q is non-nop */
			if q1 != nil {
				q1.mark |= p.mark
			}
			continue
		case ABL_5,
			ABX_5,
			ADUFFZERO_5,
			ADUFFCOPY_5:
			cursym.text.mark &^= LEAF_obj5
			fallthrough
		case ABCASE_5,
			AB_5,
			ABEQ_5,
			ABNE_5,
			ABCS_5,
			ABHS_5,
			ABCC_5,
			ABLO_5,
			ABMI_5,
			ABPL_5,
			ABVS_5,
			ABVC_5,
			ABHI_5,
			ABLS_5,
			ABGE_5,
			ABLT_5,
			ABGT_5,
			ABLE_5:
			q1 = p.pcond
			if q1 != nil {
				for q1.as == ANOP_5 {
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
		case ATEXT_5:
			autosize = p.to.offset + 4
			if autosize <= 4 {
				if cursym.text.mark&LEAF_obj5 != 0 {
					p.to.offset = -4
					autosize = 0
				}
			}
			if !(autosize != 0) && !(cursym.text.mark&LEAF_obj5 != 0) {
				if ctxt.debugvlog != 0 {
					Bprint(ctxt.bso, "save suppressed in: %s\n", cursym.name)
					Bflush(ctxt.bso)
				}
				cursym.text.mark |= LEAF_obj5
			}
			if cursym.text.mark&LEAF_obj5 != 0 {
				cursym.leaf = 1
				if !(autosize != 0) {
					break
				}
			}
			if !(p.reg&NOSPLIT_textflag != 0) {
				p = stacksplit_obj5(ctxt, p, autosize, bool2int(!(cursym.text.reg&NEEDCTXT_textflag != 0))) // emit split check
			}
			// MOVW.W		R14,$-autosize(SP)
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.scond |= C_WBIT_5
			p.from.typ = D_REG_5
			p.from.reg = REGLINK_5
			p.to.typ = D_OREG_5
			p.to.offset = -autosize
			p.to.reg = REGSP_5
			p.spadj = autosize
			if cursym.text.reg&WRAPPER_textflag != 0 {
				// g->panicwrap += autosize;
				// MOVW panicwrap_offset(g), R3
				// ADD $autosize, R3
				// MOVW R3 panicwrap_offset(g)
				p = appendp(ctxt, p)
				p.as = AMOVW_5
				p.from.typ = D_OREG_5
				p.from.reg = REGG_5
				p.from.offset = 2 * ctxt.arch.ptrsize
				p.to.typ = D_REG_5
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = AADD_5
				p.from.typ = D_CONST_5
				p.from.offset = autosize
				p.to.typ = D_REG_5
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = AMOVW_5
				p.from.typ = D_REG_5
				p.from.reg = 3
				p.to.typ = D_OREG_5
				p.to.reg = REGG_5
				p.to.offset = 2 * ctxt.arch.ptrsize
			}
		case ARET_5:
			nocache_obj5(p)
			if cursym.text.mark&LEAF_obj5 != 0 {
				if !(autosize != 0) {
					p.as = AB_5
					p.from = zprg_obj5.from
					if p.to.sym != nil { // retjmp
						p.to.typ = D_BRANCH_5
					} else {
						p.to.typ = D_OREG_5
						p.to.offset = 0
						p.to.reg = REGLINK_5
					}
					break
				}
			}
			if cursym.text.reg&WRAPPER_textflag != 0 {
				var scond int
				// Preserve original RET's cond, to allow RET.EQ
				// in the implementation of reflect.call.
				scond = p.scond
				p.scond = C_SCOND_NONE_5
				// g->panicwrap -= autosize;
				// MOVW panicwrap_offset(g), R3
				// SUB $autosize, R3
				// MOVW R3 panicwrap_offset(g)
				p.as = AMOVW_5
				p.from.typ = D_OREG_5
				p.from.reg = REGG_5
				p.from.offset = 2 * ctxt.arch.ptrsize
				p.to.typ = D_REG_5
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = ASUB_5
				p.from.typ = D_CONST_5
				p.from.offset = autosize
				p.to.typ = D_REG_5
				p.to.reg = 3
				p = appendp(ctxt, p)
				p.as = AMOVW_5
				p.from.typ = D_REG_5
				p.from.reg = 3
				p.to.typ = D_OREG_5
				p.to.reg = REGG_5
				p.to.offset = 2 * ctxt.arch.ptrsize
				p = appendp(ctxt, p)
				p.scond = scond
			}
			p.as = AMOVW_5
			p.scond |= C_PBIT_5
			p.from.typ = D_OREG_5
			p.from.offset = autosize
			p.from.reg = REGSP_5
			p.to.typ = D_REG_5
			p.to.reg = REGPC_5
			// If there are instructions following
			// this ARET, they come from a branch
			// with the same stackframe, so no spadj.
			if p.to.sym != nil { // retjmp
				p.to.reg = REGLINK_5
				q2 = appendp(ctxt, p)
				q2.as = AB_5
				q2.to.typ = D_BRANCH_5
				q2.to.sym = p.to.sym
				p.to.sym = nil
				p = q2
			}
		case AADD_5:
			if p.from.typ == D_CONST_5 && p.from.reg == NREG_5 && p.to.typ == D_REG_5 && p.to.reg == REGSP_5 {
				p.spadj = -p.from.offset
			}
		case ASUB_5:
			if p.from.typ == D_CONST_5 && p.from.reg == NREG_5 && p.to.typ == D_REG_5 && p.to.reg == REGSP_5 {
				p.spadj = p.from.offset
			}
		case ADIV_5,
			ADIVU_5,
			AMOD_5,
			AMODU_5:
			if ctxt.debugdivmod != 0 {
				break
			}
			if p.from.typ != D_REG_5 {
				break
			}
			if p.to.typ != D_REG_5 {
				break
			}
			q1 = p
			/* MOV a,4(SP) */
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.lineno = q1.lineno
			p.from.typ = D_REG_5
			p.from.reg = q1.from.reg
			p.to.typ = D_OREG_5
			p.to.reg = REGSP_5
			p.to.offset = 4
			/* MOV b,REGTMP */
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.lineno = q1.lineno
			p.from.typ = D_REG_5
			p.from.reg = q1.reg
			if q1.reg == NREG_5 {
				p.from.reg = q1.to.reg
			}
			p.to.typ = D_REG_5
			p.to.reg = REGTMP_5
			p.to.offset = 0
			/* CALL appropriate */
			p = appendp(ctxt, p)
			p.as = ABL_5
			p.lineno = q1.lineno
			p.to.typ = D_BRANCH_5
			switch o {
			case ADIV_5:
				p.to.sym = ctxt.sym_div
			case ADIVU_5:
				p.to.sym = ctxt.sym_divu
			case AMOD_5:
				p.to.sym = ctxt.sym_mod
			case AMODU_5:
				p.to.sym = ctxt.sym_modu
				break
			}
			/* MOV REGTMP, b */
			p = appendp(ctxt, p)
			p.as = AMOVW_5
			p.lineno = q1.lineno
			p.from.typ = D_REG_5
			p.from.reg = REGTMP_5
			p.from.offset = 0
			p.to.typ = D_REG_5
			p.to.reg = q1.to.reg
			/* ADD $8,SP */
			p = appendp(ctxt, p)
			p.as = AADD_5
			p.lineno = q1.lineno
			p.from.typ = D_CONST_5
			p.from.reg = NREG_5
			p.from.offset = 8
			p.reg = NREG_5
			p.to.typ = D_REG_5
			p.to.reg = REGSP_5
			p.spadj = -8
			/* Keep saved LR at 0(SP) after SP change. */
			/* MOVW 0(SP), REGTMP; MOVW REGTMP, -8!(SP) */
			/* TODO: Remove SP adjustments; see issue 6699. */
			q1.as = AMOVW_5
			q1.from.typ = D_OREG_5
			q1.from.reg = REGSP_5
			q1.from.offset = 0
			q1.reg = NREG_5
			q1.to.typ = D_REG_5
			q1.to.reg = REGTMP_5
			/* SUB $8,SP */
			q1 = appendp(ctxt, q1)
			q1.as = AMOVW_5
			q1.from.typ = D_REG_5
			q1.from.reg = REGTMP_5
			q1.reg = NREG_5
			q1.to.typ = D_OREG_5
			q1.to.reg = REGSP_5
			q1.to.offset = -8
			q1.scond |= C_WBIT_5
			q1.spadj = 8
		case AMOVW_5:
			if (p.scond&C_WBIT_5 != 0) && p.to.typ == D_OREG_5 && p.to.reg == REGSP_5 {
				p.spadj = -p.to.offset
			}
			if (p.scond&C_PBIT_5 != 0) && p.from.typ == D_OREG_5 && p.from.reg == REGSP_5 && p.to.reg != REGPC_5 {
				p.spadj = -p.from.offset
			}
			if p.from.typ == D_CONST_5 && p.from.reg == REGSP_5 && p.to.typ == D_REG_5 && p.to.reg == REGSP_5 {
				p.spadj = -p.from.offset
			}
			break
		}
	}
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
			p.pcond.mark |= LABEL_obj5
		}
	}
	for p = cursym.text; p != nil; p = p.link {
		switch p.as {
		case AMOVW_5:
			if p.to.typ == D_FREG_5 || p.from.typ == D_FREG_5 {
				goto soft
			}
			goto notsoft
		case AMOVWD_5,
			AMOVWF_5,
			AMOVDW_5,
			AMOVFW_5,
			AMOVFD_5,
			AMOVDF_5,
			AMOVF_5,
			AMOVD_5,
			ACMPF_5,
			ACMPD_5,
			AADDF_5,
			AADDD_5,
			ASUBF_5,
			ASUBD_5,
			AMULF_5,
			AMULD_5,
			ADIVF_5,
			ADIVD_5,
			ASQRTF_5,
			ASQRTD_5,
			AABSF_5,
			AABSD_5:
			goto soft
		default:
			goto notsoft
		}
	soft:
		if !(wasfloat != 0) || (p.mark&LABEL_obj5 != 0) {
			next = ctxt.prg()
			*next = *p
			// BL _sfloat(SB)
			*p = zprg_obj5
			p.link = next
			p.as = ABL_5
			p.to.typ = D_BRANCH_5
			p.to.sym = symsfloat
			p.lineno = next.lineno
			p = next
			wasfloat = 1
		}
		continue
	notsoft:
		wasfloat = 0
	}
}

func stacksplit_obj5(ctxt *Link, p *Prog, framesize int64, noctxt int) *Prog {
	var arg int64
	// MOVW			g_stackguard(g), R1
	p = appendp(ctxt, p)
	p.as = AMOVW_5
	p.from.typ = D_OREG_5
	p.from.reg = REGG_5
	p.to.typ = D_REG_5
	p.to.reg = 1
	if framesize <= StackSmall_stack {
		// small stack: SP < stackguard
		//	CMP	stackguard, SP
		p = appendp(ctxt, p)
		p.as = ACMP_5
		p.from.typ = D_REG_5
		p.from.reg = 1
		p.reg = REGSP_5
	} else if framesize <= StackBig_stack {
		// large stack: SP-framesize < stackguard-StackSmall
		//	MOVW $-framesize(SP), R2
		//	CMP stackguard, R2
		p = appendp(ctxt, p)
		p.as = AMOVW_5
		p.from.typ = D_CONST_5
		p.from.reg = REGSP_5
		p.from.offset = -framesize
		p.to.typ = D_REG_5
		p.to.reg = 2
		p = appendp(ctxt, p)
		p.as = ACMP_5
		p.from.typ = D_REG_5
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
		p.as = ACMP_5
		p.from.typ = D_CONST_5
		p.from.offset = int64(uint32(StackPreempt_stack & 0xFFFFFFFF))
		p.reg = 1
		p = appendp(ctxt, p)
		p.as = AMOVW_5
		p.from.typ = D_CONST_5
		p.from.reg = REGSP_5
		p.from.offset = StackGuard_stack
		p.to.typ = D_REG_5
		p.to.reg = 2
		p.scond = C_SCOND_NE_5
		p = appendp(ctxt, p)
		p.as = ASUB_5
		p.from.typ = D_REG_5
		p.from.reg = 1
		p.to.typ = D_REG_5
		p.to.reg = 2
		p.scond = C_SCOND_NE_5
		p = appendp(ctxt, p)
		p.as = AMOVW_5
		p.from.typ = D_CONST_5
		p.from.offset = framesize + (StackGuard_stack - StackSmall_stack)
		p.to.typ = D_REG_5
		p.to.reg = 3
		p.scond = C_SCOND_NE_5
		p = appendp(ctxt, p)
		p.as = ACMP_5
		p.from.typ = D_REG_5
		p.from.reg = 3
		p.reg = 2
		p.scond = C_SCOND_NE_5
	}
	// MOVW.LS		$framesize, R1
	p = appendp(ctxt, p)
	p.as = AMOVW_5
	p.scond = C_SCOND_LS_5
	p.from.typ = D_CONST_5
	p.from.offset = framesize
	p.to.typ = D_REG_5
	p.to.reg = 1
	// MOVW.LS		$args, R2
	p = appendp(ctxt, p)
	p.as = AMOVW_5
	p.scond = C_SCOND_LS_5
	p.from.typ = D_CONST_5
	arg = ctxt.cursym.text.to.offset2
	if arg == 1 { // special marker for known 0
		arg = 0
	}
	if arg&3 != 0 {
		ctxt.diag("misaligned argument size in stack split")
	}
	p.from.offset = arg
	p.to.typ = D_REG_5
	p.to.reg = 2
	// MOVW.LS	R14, R3
	p = appendp(ctxt, p)
	p.as = AMOVW_5
	p.scond = C_SCOND_LS_5
	p.from.typ = D_REG_5
	p.from.reg = REGLINK_5
	p.to.typ = D_REG_5
	p.to.reg = 3
	// BL.LS		runtime.morestack(SB) // modifies LR, returns with LO still asserted
	p = appendp(ctxt, p)
	p.as = ABL_5
	p.scond = C_SCOND_LS_5
	p.to.typ = D_BRANCH_5
	p.to.sym = ctxt.symmorestack[noctxt]
	// BLS	start
	p = appendp(ctxt, p)
	p.as = ABLS_5
	p.to.typ = D_BRANCH_5
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

func follow_obj5(ctxt *Link, s *LSym) {
	var firstp *Prog
	var lastp *Prog
	ctxt.cursym = s
	firstp = ctxt.prg()
	lastp = firstp
	xfol_obj5(ctxt, s.text, &lastp)
	lastp.link = nil
	s.text = firstp.link
}

func relinv_obj5(a int) int {
	switch a {
	case ABEQ_5:
		return ABNE_5
	case ABNE_5:
		return ABEQ_5
	case ABCS_5:
		return ABCC_5
	case ABHS_5:
		return ABLO_5
	case ABCC_5:
		return ABCS_5
	case ABLO_5:
		return ABHS_5
	case ABMI_5:
		return ABPL_5
	case ABPL_5:
		return ABMI_5
	case ABVS_5:
		return ABVC_5
	case ABVC_5:
		return ABVS_5
	case ABHI_5:
		return ABLS_5
	case ABLS_5:
		return ABHI_5
	case ABGE_5:
		return ABLT_5
	case ABLT_5:
		return ABGE_5
	case ABGT_5:
		return ABLE_5
	case ABLE_5:
		return ABGT_5
	}
	log.Fatalf("unknown relation: %s", anames5[a])
	return 0
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
	if a == AB_5 {
		q = p.pcond
		if q != nil && q.as != ATEXT_5 {
			p.mark |= FOLL_obj5
			p = q
			if !(p.mark&FOLL_obj5 != 0) {
				goto loop
			}
		}
	}
	if p.mark&FOLL_obj5 != 0 {
		i = 0
		q = p
		for ; i < 4; (func() { i++; q = q.link })() {
			if q == *last || q == nil {
				break
			}
			a = q.as
			if a == ANOP_5 {
				i--
				continue
			}
			if a == AB_5 || (a == ARET_5 && q.scond == C_SCOND_NONE_5) || a == ARFE_5 || a == AUNDEF_5 {
				goto copy
			}
			if q.pcond == nil || (q.pcond.mark&FOLL_obj5 != 0) {
				continue
			}
			if a != ABEQ_5 && a != ABNE_5 {
				continue
			}
		copy:
			for {
				r = ctxt.prg()
				*r = *p
				if !(r.mark&FOLL_obj5 != 0) {
					fmt.Printf("can't happen 1\n")
				}
				r.mark |= FOLL_obj5
				if p != q {
					p = p.link
					(*last).link = r
					*last = r
					continue
				}
				(*last).link = r
				*last = r
				if a == AB_5 || (a == ARET_5 && q.scond == C_SCOND_NONE_5) || a == ARFE_5 || a == AUNDEF_5 {
					return
				}
				r.as = ABNE_5
				if a == ABNE_5 {
					r.as = ABEQ_5
				}
				r.pcond = p.link
				r.link = p.pcond
				if !(r.link.mark&FOLL_obj5 != 0) {
					xfol_obj5(ctxt, r.link, last)
				}
				if !(r.pcond.mark&FOLL_obj5 != 0) {
					fmt.Printf("can't happen 2\n")
				}
				return
			}
		}
		a = AB_5
		q = ctxt.prg()
		q.as = a
		q.lineno = p.lineno
		q.to.typ = D_BRANCH_5
		q.to.offset = p.pc
		q.pcond = p
		p = q
	}
	p.mark |= FOLL_obj5
	(*last).link = p
	*last = p
	if a == AB_5 || (a == ARET_5 && p.scond == C_SCOND_NONE_5) || a == ARFE_5 || a == AUNDEF_5 {
		return
	}
	if p.pcond != nil {
		if a != ABL_5 && a != ABX_5 && p.link != nil {
			q = brchain(ctxt, p.link)
			if a != ATEXT_5 && a != ABCASE_5 {
				if q != nil && (q.mark&FOLL_obj5 != 0) {
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
			if q.mark&FOLL_obj5 != 0 {
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

var linkarm = LinkArch{
	name:          "arm",
	thechar:       '5',
	addstacksplit: addstacksplit_obj5,
	assemble:      span5,
	datasize:      datasize_obj5,
	follow:        follow_obj5,
	iscall:        iscall_obj5,
	isdata:        isdata_obj5,
	prg:           prg_obj5,
	progedit:      progedit_obj5,
	settextflag:   settextflag_obj5,
	symtype:       symtype_obj5,
	textflag:      textflag_obj5,
	Pconv:         Pconv_list5,
	minlc:         4,
	ptrsize:       4,
	regsize:       4,
	D_ADDR:        D_ADDR_5,
	D_AUTO:        D_AUTO_5,
	D_BRANCH:      D_BRANCH_5,
	D_CONST:       D_CONST_5,
	D_EXTERN:      D_EXTERN_5,
	D_FCONST:      D_FCONST_5,
	D_NONE:        D_NONE_5,
	D_PARAM:       D_PARAM_5,
	D_SCONST:      D_SCONST_5,
	D_STATIC:      D_STATIC_5,
	ACALL:         ABL_5,
	ADATA:         ADATA_5,
	AEND:          AEND_5,
	AFUNCDATA:     AFUNCDATA_5,
	AGLOBL:        AGLOBL_5,
	AJMP:          AB_5,
	ANOP:          ANOP_5,
	APCDATA:       APCDATA_5,
	ARET:          ARET_5,
	ATEXT:         ATEXT_5,
	ATYPE:         ATYPE_5,
	AUSEFIELD:     AUSEFIELD_5,
}
