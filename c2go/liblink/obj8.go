package main

import (
	"fmt"
	"log"
	"math"
)

// Inferno utils/8l/pass.c
// http://code.google.com/p/inferno-os/source/browse/utils/8l/pass.c
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
var zprg_obj8 = Prog{
	back: 2,
	as:   AGOK_8,
	from: Addr{
		typ:   D_NONE_8,
		index: D_NONE_8,
		scale: 1,
	},
	to: Addr{
		typ:   D_NONE_8,
		index: D_NONE_8,
		scale: 1,
	},
}

func symtype_obj8(a *Addr) int {
	var t int64
	t = a.typ
	if t == D_ADDR_8 {
		t = a.index
	}
	return int(t)
}

func isdata_obj8(p *Prog) int {
	return bool2int(p.as == ADATA_8 || p.as == AGLOBL_8)
}

func iscall_obj8(p *Prog) int {
	return bool2int(p.as == ACALL_8)
}

func datasize_obj8(p *Prog) int {
	return int(p.from.scale)
}

func textflag_obj8(p *Prog) int {
	return int(p.from.scale)
}

func settextflag_obj8(p *Prog, f int) {
	p.from.scale = int64(f)
}

func canuselocaltls_obj8(ctxt *Link) int {
	switch ctxt.headtype {
	case Hlinux,
		Hnacl,
		Hplan9,
		Hwindows:
		return 0
	}
	return 1
}

func progedit_obj8(ctxt *Link, p *Prog) {
	var literal string
	var s *LSym
	var q *Prog
	// See obj6.c for discussion of TLS.
	if canuselocaltls_obj8(ctxt) != 0 {
		// Reduce TLS initial exec model to TLS local exec model.
		// Sequences like
		//	MOVL TLS, BX
		//	... off(BX)(TLS*1) ...
		// become
		//	NOP
		//	... off(TLS) ...
		if p.as == AMOVL_8 && p.from.typ == D_TLS_8 && D_AX_8 <= p.to.typ && p.to.typ <= D_DI_8 {
			p.as = ANOP_8
			p.from.typ = D_NONE_8
			p.to.typ = D_NONE_8
		}
		if p.from.index == D_TLS_8 && D_INDIR_8+D_AX_8 <= p.from.typ && p.from.typ <= D_INDIR_8+D_DI_8 {
			p.from.typ = D_INDIR_8 + D_TLS_8
			p.from.scale = 0
			p.from.index = D_NONE_8
		}
		if p.to.index == D_TLS_8 && D_INDIR_8+D_AX_8 <= p.to.typ && p.to.typ <= D_INDIR_8+D_DI_8 {
			p.to.typ = D_INDIR_8 + D_TLS_8
			p.to.scale = 0
			p.to.index = D_NONE_8
		}
	} else {
		// As a courtesy to the C compilers, rewrite TLS local exec load as TLS initial exec load.
		// The instruction
		//	MOVL off(TLS), BX
		// becomes the sequence
		//	MOVL TLS, BX
		//	MOVL off(BX)(TLS*1), BX
		// This allows the C compilers to emit references to m and g using the direct off(TLS) form.
		if p.as == AMOVL_8 && p.from.typ == D_INDIR_8+D_TLS_8 && D_AX_8 <= p.to.typ && p.to.typ <= D_DI_8 {
			q = appendp(ctxt, p)
			q.as = p.as
			q.from = p.from
			q.from.typ = D_INDIR_8 + p.to.typ
			q.from.index = D_TLS_8
			q.from.scale = 2 // TODO: use 1
			q.to = p.to
			p.from.typ = D_TLS_8
			p.from.index = D_NONE_8
			p.from.offset = 0
		}
	}
	// TODO: Remove.
	if ctxt.headtype == Hplan9 {
		if p.from.scale == 1 && p.from.index == D_TLS_8 {
			p.from.scale = 2
		}
		if p.to.scale == 1 && p.to.index == D_TLS_8 {
			p.to.scale = 2
		}
	}
	// Rewrite CALL/JMP/RET to symbol as D_BRANCH.
	switch p.as {
	case ACALL_8,
		AJMP_8,
		ARET_8:
		if (p.to.typ == D_EXTERN_8 || p.to.typ == D_STATIC_8) && p.to.sym != nil {
			p.to.typ = D_BRANCH_8
		}
		break
	}
	// Rewrite float constants to values stored in memory.
	switch p.as {
	case AFMOVF_8,
		AFADDF_8,
		AFSUBF_8,
		AFSUBRF_8,
		AFMULF_8,
		AFDIVF_8,
		AFDIVRF_8,
		AFCOMF_8,
		AFCOMFP_8,
		AMOVSS_8,
		AADDSS_8,
		ASUBSS_8,
		AMULSS_8,
		ADIVSS_8,
		ACOMISS_8,
		AUCOMISS_8:
		if p.from.typ == D_FCONST_8 {
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
			p.from.typ = D_EXTERN_8
			p.from.sym = s
			p.from.offset = 0
		}
	case AFMOVD_8,
		AFADDD_8,
		AFSUBD_8,
		AFSUBRD_8,
		AFMULD_8,
		AFDIVD_8,
		AFDIVRD_8,
		AFCOMD_8,
		AFCOMDP_8,
		AMOVSD_8,
		AADDSD_8,
		ASUBSD_8,
		AMULSD_8,
		ADIVSD_8,
		ACOMISD_8,
		AUCOMISD_8:
		if p.from.typ == D_FCONST_8 {
			var i64 uint64
			i64 = math.Float64bits(p.from.u.dval)
			literal = fmt.Sprintf("$f64.%016x", uint64(i64))
			s = linklookup(ctxt, literal, 0)
			if s.typ == 0 {
				s.typ = SRODATA
				adduint64(ctxt, s, i64)
				s.reachable = 0
			}
			p.from.typ = D_EXTERN_8
			p.from.sym = s
			p.from.offset = 0
		}
		break
	}
}

func prg_obj8() *Prog {
	var p *Prog
	p = new(Prog)
	*p = zprg_obj8
	return p
}

func addstacksplit_obj8(ctxt *Link, cursym *LSym) {
	var p *Prog
	var q *Prog
	var autoffset int64
	var deltasp int64
	var a int64
	if ctxt.symmorestack[0] == nil {
		ctxt.symmorestack[0] = linklookup(ctxt, "runtime.morestack", 0)
		ctxt.symmorestack[1] = linklookup(ctxt, "runtime.morestack_noctxt", 0)
	}
	if ctxt.headtype == Hplan9 && ctxt.plan9privates == nil {
		ctxt.plan9privates = linklookup(ctxt, "_privates", 0)
	}
	ctxt.cursym = cursym
	if cursym.text == nil || cursym.text.link == nil {
		return
	}
	p = cursym.text
	autoffset = p.to.offset
	if autoffset < 0 {
		autoffset = 0
	}
	cursym.locals = autoffset
	cursym.args = p.to.offset2
	q = nil
	if !(p.from.scale&NOSPLIT_textflag != 0) || (p.from.scale&WRAPPER_textflag != 0) {
		p = appendp(ctxt, p)
		p = load_g_cx_obj8(ctxt, p) // load g into CX
	}
	if !(cursym.text.from.scale&NOSPLIT_textflag != 0) {
		p = stacksplit_obj8(ctxt, p, autoffset, bool2int(!(cursym.text.from.scale&NEEDCTXT_textflag != 0)), &q) // emit split check
	}
	if autoffset != 0 {
		p = appendp(ctxt, p)
		p.as = AADJSP_8
		p.from.typ = D_CONST_8
		p.from.offset = autoffset
		p.spadj = autoffset
	} else {
		// zero-byte stack adjustment.
		// Insert a fake non-zero adjustment so that stkcheck can
		// recognize the end of the stack-splitting prolog.
		p = appendp(ctxt, p)
		p.as = ANOP_8
		p.spadj = -ctxt.arch.ptrsize
		p = appendp(ctxt, p)
		p.as = ANOP_8
		p.spadj = ctxt.arch.ptrsize
	}
	if q != nil {
		q.pcond = p
	}
	deltasp = autoffset
	if cursym.text.from.scale&WRAPPER_textflag != 0 {
		// g->panicwrap += autoffset + ctxt->arch->ptrsize;
		p = appendp(ctxt, p)
		p.as = AADDL_8
		p.from.typ = D_CONST_8
		p.from.offset = autoffset + ctxt.arch.ptrsize
		p.to.typ = D_INDIR_8 + D_CX_8
		p.to.offset = 2 * ctxt.arch.ptrsize
	}
	if ctxt.debugzerostack != 0 && autoffset != 0 && !(cursym.text.from.scale&NOSPLIT_textflag != 0) {
		// 8l -Z means zero the stack frame on entry.
		// This slows down function calls but can help avoid
		// false positives in garbage collection.
		p = appendp(ctxt, p)
		p.as = AMOVL_8
		p.from.typ = D_SP_8
		p.to.typ = D_DI_8
		p = appendp(ctxt, p)
		p.as = AMOVL_8
		p.from.typ = D_CONST_8
		p.from.offset = autoffset / 4
		p.to.typ = D_CX_8
		p = appendp(ctxt, p)
		p.as = AMOVL_8
		p.from.typ = D_CONST_8
		p.from.offset = 0
		p.to.typ = D_AX_8
		p = appendp(ctxt, p)
		p.as = AREP_8
		p = appendp(ctxt, p)
		p.as = ASTOSL_8
	}
	for ; p != nil; p = p.link {
		a = p.from.typ
		if a == D_AUTO_8 {
			p.from.offset += deltasp
		}
		if a == D_PARAM_8 {
			p.from.offset += deltasp + 4
		}
		a = p.to.typ
		if a == D_AUTO_8 {
			p.to.offset += deltasp
		}
		if a == D_PARAM_8 {
			p.to.offset += deltasp + 4
		}
		switch p.as {
		default:
			continue
		case APUSHL_8,
			APUSHFL_8:
			deltasp += 4
			p.spadj = 4
			continue
		case APUSHW_8,
			APUSHFW_8:
			deltasp += 2
			p.spadj = 2
			continue
		case APOPL_8,
			APOPFL_8:
			deltasp -= 4
			p.spadj = -4
			continue
		case APOPW_8,
			APOPFW_8:
			deltasp -= 2
			p.spadj = -2
			continue
		case ARET_8:
			break
		}
		if autoffset != deltasp {
			ctxt.diag("unbalanced PUSH/POP")
		}
		if cursym.text.from.scale&WRAPPER_textflag != 0 {
			p = load_g_cx_obj8(ctxt, p)
			p = appendp(ctxt, p)
			// g->panicwrap -= autoffset + ctxt->arch->ptrsize;
			p.as = ASUBL_8
			p.from.typ = D_CONST_8
			p.from.offset = autoffset + ctxt.arch.ptrsize
			p.to.typ = D_INDIR_8 + D_CX_8
			p.to.offset = 2 * ctxt.arch.ptrsize
			p = appendp(ctxt, p)
			p.as = ARET_8
		}
		if autoffset != 0 {
			p.as = AADJSP_8
			p.from.typ = D_CONST_8
			p.from.offset = -autoffset
			p.spadj = -autoffset
			p = appendp(ctxt, p)
			p.as = ARET_8
			// If there are instructions following
			// this ARET, they come from a branch
			// with the same stackframe, so undo
			// the cleanup.
			p.spadj = +autoffset
		}
		if p.to.sym != nil { // retjmp
			p.as = AJMP_8
		}
	}
}

// Append code to p to load g into cx.
// Overwrites p with the first instruction (no first appendp).
// Overwriting p is unusual but it lets use this in both the
// prologue (caller must call appendp first) and in the epilogue.
// Returns last new instruction.
func load_g_cx_obj8(ctxt *Link, p *Prog) *Prog {
	var next *Prog
	p.as = AMOVL_8
	p.from.typ = D_INDIR_8 + D_TLS_8
	p.from.offset = 0
	p.to.typ = D_CX_8
	next = p.link
	progedit_obj8(ctxt, p)
	for p.link != next {
		p = p.link
	}
	if p.from.index == D_TLS_8 {
		p.from.scale = 2
	}
	return p
}

// Append code to p to check for stack split.
// Appends to (does not overwrite) p.
// Assumes g is in CX.
// Returns last new instruction.
// On return, *jmpok is the instruction that should jump
// to the stack frame allocation if no split is needed.
func stacksplit_obj8(ctxt *Link, p *Prog, framesize int64, noctxt int, jmpok **Prog) *Prog {
	var q *Prog
	var q1 *Prog
	var arg int64
	if ctxt.debugstack != 0 {
		// 8l -K means check not only for stack
		// overflow but stack underflow.
		// On underflow, INT 3 (breakpoint).
		// Underflow itself is rare but this also
		// catches out-of-sync stack guard info.
		p = appendp(ctxt, p)
		p.as = ACMPL_8
		p.from.typ = D_INDIR_8 + D_CX_8
		p.from.offset = 4
		p.to.typ = D_SP_8
		p = appendp(ctxt, p)
		p.as = AJCC_8
		p.to.typ = D_BRANCH_8
		p.to.offset = 4
		q1 = p
		p = appendp(ctxt, p)
		p.as = AINT_8
		p.from.typ = D_CONST_8
		p.from.offset = 3
		p = appendp(ctxt, p)
		p.as = ANOP_8
		q1.pcond = p
	}
	q1 = nil
	if framesize <= StackSmall_stack {
		// small stack: SP <= stackguard
		//	CMPL SP, stackguard
		p = appendp(ctxt, p)
		p.as = ACMPL_8
		p.from.typ = D_SP_8
		p.to.typ = D_INDIR_8 + D_CX_8
	} else if framesize <= StackBig_stack {
		// large stack: SP-framesize <= stackguard-StackSmall
		//	LEAL -(framesize-StackSmall)(SP), AX
		//	CMPL AX, stackguard
		p = appendp(ctxt, p)
		p.as = ALEAL_8
		p.from.typ = D_INDIR_8 + D_SP_8
		p.from.offset = -(framesize - StackSmall_stack)
		p.to.typ = D_AX_8
		p = appendp(ctxt, p)
		p.as = ACMPL_8
		p.from.typ = D_AX_8
		p.to.typ = D_INDIR_8 + D_CX_8
	} else {
		// Such a large stack we need to protect against wraparound
		// if SP is close to zero.
		//	SP-stackguard+StackGuard <= framesize + (StackGuard-StackSmall)
		// The +StackGuard on both sides is required to keep the left side positive:
		// SP is allowed to be slightly below stackguard. See stack.h.
		//
		// Preemption sets stackguard to StackPreempt, a very large value.
		// That breaks the math above, so we have to check for that explicitly.
		//	MOVL	stackguard, CX
		//	CMPL	CX, $StackPreempt
		//	JEQ	label-of-call-to-morestack
		//	LEAL	StackGuard(SP), AX
		//	SUBL	stackguard, AX
		//	CMPL	AX, $(framesize+(StackGuard-StackSmall))
		p = appendp(ctxt, p)
		p.as = AMOVL_8
		p.from.typ = D_INDIR_8 + D_CX_8
		p.from.offset = 0
		p.to.typ = D_SI_8
		p = appendp(ctxt, p)
		p.as = ACMPL_8
		p.from.typ = D_SI_8
		p.to.typ = D_CONST_8
		p.to.offset = int64(uint32(StackPreempt_stack & 0xFFFFFFFF))
		p = appendp(ctxt, p)
		p.as = AJEQ_8
		p.to.typ = D_BRANCH_8
		q1 = p
		p = appendp(ctxt, p)
		p.as = ALEAL_8
		p.from.typ = D_INDIR_8 + D_SP_8
		p.from.offset = StackGuard_stack
		p.to.typ = D_AX_8
		p = appendp(ctxt, p)
		p.as = ASUBL_8
		p.from.typ = D_SI_8
		p.from.offset = 0
		p.to.typ = D_AX_8
		p = appendp(ctxt, p)
		p.as = ACMPL_8
		p.from.typ = D_AX_8
		p.to.typ = D_CONST_8
		p.to.offset = framesize + (StackGuard_stack - StackSmall_stack)
	}
	// common
	p = appendp(ctxt, p)
	p.as = AJHI_8
	p.to.typ = D_BRANCH_8
	p.to.offset = 4
	q = p
	p = appendp(ctxt, p) // save frame size in DI
	p.as = AMOVL_8
	p.to.typ = D_DI_8
	p.from.typ = D_CONST_8
	// If we ask for more stack, we'll get a minimum of StackMin bytes.
	// We need a stack frame large enough to hold the top-of-stack data,
	// the function arguments+results, our caller's PC, our frame,
	// a word for the return PC of the next call, and then the StackLimit bytes
	// that must be available on entry to any function called from a function
	// that did a stack check.  If StackMin is enough, don't ask for a specific
	// amount: then we can use the custom functions and save a few
	// instructions.
	if StackTop_stack+ctxt.cursym.text.to.offset2+ctxt.arch.ptrsize+framesize+ctxt.arch.ptrsize+StackLimit_stack >= StackMin_stack {
		p.from.offset = (framesize + 7) &^ 7
	}
	arg = ctxt.cursym.text.to.offset2
	if arg == 1 { // special marker for known 0
		arg = 0
	}
	if arg&3 != 0 {
		ctxt.diag("misaligned argument size in stack split")
	}
	p = appendp(ctxt, p) // save arg size in AX
	p.as = AMOVL_8
	p.to.typ = D_AX_8
	p.from.typ = D_CONST_8
	p.from.offset = arg
	p = appendp(ctxt, p)
	p.as = ACALL_8
	p.to.typ = D_BRANCH_8
	p.to.sym = ctxt.symmorestack[noctxt]
	p = appendp(ctxt, p)
	p.as = AJMP_8
	p.to.typ = D_BRANCH_8
	p.pcond = ctxt.cursym.text.link
	if q != nil {
		q.pcond = p.link
	}
	if q1 != nil {
		q1.pcond = q.link
	}
	*jmpok = q
	return p
}

func follow_obj8(ctxt *Link, s *LSym) {
	var firstp *Prog
	var lastp *Prog
	ctxt.cursym = s
	firstp = ctxt.prg()
	lastp = firstp
	xfol_obj8(ctxt, s.text, &lastp)
	lastp.link = nil
	s.text = firstp.link
}

func nofollow_obj8(a int) int {
	switch a {
	case AJMP_8,
		ARET_8,
		AIRETL_8,
		AIRETW_8,
		AUNDEF_8:
		return 1
	}
	return 0
}

func pushpop_obj8(a int) int {
	switch a {
	case APUSHL_8,
		APUSHFL_8,
		APUSHW_8,
		APUSHFW_8,
		APOPL_8,
		APOPFL_8,
		APOPW_8,
		APOPFW_8:
		return 1
	}
	return 0
}

func relinv_obj8(a int) int {
	switch a {
	case AJEQ_8:
		return AJNE_8
	case AJNE_8:
		return AJEQ_8
	case AJLE_8:
		return AJGT_8
	case AJLS_8:
		return AJHI_8
	case AJLT_8:
		return AJGE_8
	case AJMI_8:
		return AJPL_8
	case AJGE_8:
		return AJLT_8
	case AJPL_8:
		return AJMI_8
	case AJGT_8:
		return AJLE_8
	case AJHI_8:
		return AJLS_8
	case AJCS_8:
		return AJCC_8
	case AJCC_8:
		return AJCS_8
	case AJPS_8:
		return AJPC_8
	case AJPC_8:
		return AJPS_8
	case AJOS_8:
		return AJOC_8
	case AJOC_8:
		return AJOS_8
	}
	log.Fatalf("unknown relation: %s", anames8[a])
	return 0
}

func xfol_obj8(ctxt *Link, p *Prog, last **Prog) {
	var q *Prog
	var i int
	var a int
loop:
	if p == nil {
		return
	}
	if p.as == AJMP_8 {
		q = p.pcond
		if q != nil && q.as != ATEXT_8 {
			/* mark instruction as done and continue layout at target of jump */
			p.mark = 1
			p = q
			if p.mark == 0 {
				goto loop
			}
		}
	}
	if p.mark != 0 {
		/*
		 * p goes here, but already used it elsewhere.
		 * copy up to 4 instructions or else branch to other copy.
		 */
		i = 0
		q = p
		for ; i < 4; (func() { i++; q = q.link })() {
			if q == nil {
				break
			}
			if q == *last {
				break
			}
			a = q.as
			if a == ANOP_8 {
				i--
				continue
			}
			if nofollow_obj8(a) != 0 || pushpop_obj8(a) != 0 {
				break // NOTE(rsc): arm does goto copy
			}
			if q.pcond == nil || q.pcond.mark != 0 {
				continue
			}
			if a == ACALL_8 || a == ALOOP_8 {
				continue
			}
			for {
				if p.as == ANOP_8 {
					p = p.link
					continue
				}
				q = copyp(ctxt, p)
				p = p.link
				q.mark = 1
				(*last).link = q
				*last = q
				if q.as != a || q.pcond == nil || q.pcond.mark != 0 {
					continue
				}
				q.as = relinv_obj8(q.as)
				p = q.pcond
				q.pcond = q.link
				q.link = p
				xfol_obj8(ctxt, q.link, last)
				p = q.link
				if p.mark != 0 {
					return
				}
				goto loop /* */
			}
		}
		q = ctxt.prg()
		q.as = AJMP_8
		q.lineno = p.lineno
		q.to.typ = D_BRANCH_8
		q.to.offset = p.pc
		q.pcond = p
		p = q
	}
	/* emit p */
	p.mark = 1
	(*last).link = p
	*last = p
	a = p.as
	/* continue loop with what comes after p */
	if nofollow_obj8(a) != 0 {
		return
	}
	if p.pcond != nil && a != ACALL_8 {
		/*
		 * some kind of conditional branch.
		 * recurse to follow one path.
		 * continue loop on the other.
		 */
		q = brchain(ctxt, p.pcond)
		if q != nil {
			p.pcond = q
		}
		q = brchain(ctxt, p.link)
		if q != nil {
			p.link = q
		}
		if p.from.typ == D_CONST_8 {
			if p.from.offset == 1 {
				/*
				 * expect conditional jump to be taken.
				 * rewrite so that's the fall-through case.
				 */
				p.as = relinv_obj8(a)
				q = p.link
				p.link = p.pcond
				p.pcond = q
			}
		} else {
			q = p.link
			if q.mark != 0 {
				if a != ALOOP_8 {
					p.as = relinv_obj8(a)
					p.link = p.pcond
					p.pcond = q
				}
			}
		}
		xfol_obj8(ctxt, p.link, last)
		if p.pcond.mark != 0 {
			return
		}
		p = p.pcond
		goto loop
	}
	p = p.link
	goto loop
}

var link386 = LinkArch{
	name:          "386",
	thechar:       '8',
	addstacksplit: addstacksplit_obj8,
	assemble:      span8,
	datasize:      datasize_obj8,
	follow:        follow_obj8,
	iscall:        iscall_obj8,
	isdata:        isdata_obj8,
	prg:           prg_obj8,
	progedit:      progedit_obj8,
	settextflag:   settextflag_obj8,
	symtype:       symtype_obj8,
	textflag:      textflag_obj8,
	Pconv:         Pconv_list8,
	minlc:         1,
	ptrsize:       4,
	regsize:       4,
	D_ADDR:        D_ADDR_8,
	D_AUTO:        D_AUTO_8,
	D_BRANCH:      D_BRANCH_8,
	D_CONST:       D_CONST_8,
	D_EXTERN:      D_EXTERN_8,
	D_FCONST:      D_FCONST_8,
	D_NONE:        D_NONE_8,
	D_PARAM:       D_PARAM_8,
	D_SCONST:      D_SCONST_8,
	D_STATIC:      D_STATIC_8,
	ACALL:         ACALL_8,
	ADATA:         ADATA_8,
	AEND:          AEND_8,
	AFUNCDATA:     AFUNCDATA_8,
	AGLOBL:        AGLOBL_8,
	AJMP:          AJMP_8,
	ANOP:          ANOP_8,
	APCDATA:       APCDATA_8,
	ARET:          ARET_8,
	ATEXT:         ATEXT_8,
	ATYPE:         ATYPE_8,
	AUSEFIELD:     AUSEFIELD_8,
}
