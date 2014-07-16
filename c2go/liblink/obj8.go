package main

import (
	"fmt"
	"math"
)

var link386 LinkArch

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
	var t int
	t = a.typ
	if t == int(D_ADDR_8) {
		t = a.index
	}
	return t
}

func isdata_obj8(p *Prog) bool {
	return p.as == int(ADATA_8) || p.as == int(AGLOBL_8)
}

func iscall_obj8(p *Prog) bool {
	return p.as == int(ACALL_8)
}

func datasize_obj8(p *Prog) int {
	return p.from.scale
}

func textflag_obj8(p *Prog) int {
	return p.from.scale
}

func settextflag_obj8(p *Prog, f int) {
	p.from.scale = f
}

func canuselocaltls_obj8(ctxt *Link) int {
	switch ctxt.headtype {
	case Hlinux:
	case Hnacl:
	case Hplan9:
	case Hwindows:
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
		if p.as == int(AMOVL_8) && p.from.typ == int(D_TLS_8) && D_AX_8 <= int(p.to.typ) && p.to.typ <= int(D_DI_8) {
			p.as = int(ANOP_8)
			p.from.typ = int(D_NONE_8)
			p.to.typ = int(D_NONE_8)
		}
		if p.from.index == int(D_TLS_8) && D_INDIR_8+D_AX_8 <= int(p.from.typ) && p.from.typ <= int(D_INDIR_8+D_DI_8) {
			p.from.typ = int(D_INDIR_8 + D_TLS_8)
			p.from.scale = 0
			p.from.index = int(D_NONE_8)
		}
		if p.to.index == int(D_TLS_8) && D_INDIR_8+D_AX_8 <= int(p.to.typ) && p.to.typ <= int(D_INDIR_8+D_DI_8) {
			p.to.typ = int(D_INDIR_8 + D_TLS_8)
			p.to.scale = 0
			p.to.index = int(D_NONE_8)
		}
	} else {
		// As a courtesy to the C compilers, rewrite TLS local exec load as TLS initial exec load.
		// The instruction
		//	MOVL off(TLS), BX
		// becomes the sequence
		//	MOVL TLS, BX
		//	MOVL off(BX)(TLS*1), BX
		// This allows the C compilers to emit references to m and g using the direct off(TLS) form.
		if p.as == int(AMOVL_8) && p.from.typ == int(D_INDIR_8+D_TLS_8) && D_AX_8 <= int(p.to.typ) && p.to.typ <= int(D_DI_8) {
			q = appendp(ctxt, p)
			q.as = p.as
			q.from = p.from
			q.from.typ = int(D_INDIR_8 + int(p.to.typ))
			q.from.index = int(D_TLS_8)
			q.from.scale = 2 // TODO: use 1
			q.to = p.to
			p.from.typ = int(D_TLS_8)
			p.from.index = int(D_NONE_8)
			p.from.offset = 0
		}
	}
	// TODO: Remove.
	if ctxt.headtype == int(Hplan9) {
		if p.from.scale == 1 && p.from.index == int(D_TLS_8) {
			p.from.scale = 2
		}
		if p.to.scale == 1 && p.to.index == int(D_TLS_8) {
			p.to.scale = 2
		}
	}
	// Rewrite CALL/JMP/RET to symbol as D_BRANCH.
	switch p.as {
	case ACALL_8:
	case AJMP_8:
	case ARET_8:
		if (p.to.typ == int(D_EXTERN_8) || p.to.typ == int(D_STATIC_8)) && p.to.sym != nil {
			p.to.typ = int(D_BRANCH_8)
		}
		break
	}
	// Rewrite float constants to values stored in memory.
	switch p.as {
	case AFMOVF_8:
	case AFADDF_8:
	case AFSUBF_8:
	case AFSUBRF_8:
	case AFMULF_8:
	case AFDIVF_8:
	case AFDIVRF_8:
	case AFCOMF_8:
	case AFCOMFP_8:
	case AMOVSS_8:
	case AADDSS_8:
	case ASUBSS_8:
	case AMULSS_8:
	case ADIVSS_8:
	case ACOMISS_8:
	case AUCOMISS_8:
		if p.from.typ == int(D_FCONST_8) {
			var i32 int32
			var f32 float32
			f32 = float32(p.from.u.dval)
			i32 = int32(math.Float32bits(f32))
			literal = fmt.Sprintf("$f32.%08x", uint32(i32))
			s = linklookup(ctxt, string(literal), 0)
			if s.typ == 0 {
				s.typ = int(SRODATA)
				adduint32(ctxt, s, uint32(i32))
				s.reachable = 0
			}
			p.from.typ = int(D_EXTERN_8)
			p.from.sym = s
			p.from.offset = 0
		}
		break
	case AFMOVD_8:
	case AFADDD_8:
	case AFSUBD_8:
	case AFSUBRD_8:
	case AFMULD_8:
	case AFDIVD_8:
	case AFDIVRD_8:
	case AFCOMD_8:
	case AFCOMDP_8:
	case AMOVSD_8:
	case AADDSD_8:
	case ASUBSD_8:
	case AMULSD_8:
	case ADIVSD_8:
	case ACOMISD_8:
	case AUCOMISD_8:
		if p.from.typ == int(D_FCONST_8) {
			var i64 int64
			i64 = int64(math.Float64bits(p.from.u.dval))
			literal = fmt.Sprintf("$f64.%016x", uint64(i64))
			s = linklookup(ctxt, string(literal), 0)
			if s.typ == 0 {
				s.typ = int(SRODATA)
				adduint64(ctxt, s, uint64(i64))
				s.reachable = 0
			}
			p.from.typ = int(D_EXTERN_8)
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

// Append code to p to load g into cx.
// Overwrites p with the first instruction (no first appendp).
// Overwriting p is unusual but it lets use this in both the
// prologue (caller must call appendp first) and in the epilogue.
// Returns last new instruction.
func load_g_cx_obj8(ctxt *Link, p *Prog) *Prog {
	var next *Prog
	p.as = int(AMOVL_8)
	p.from.typ = int(D_INDIR_8 + D_TLS_8)
	p.from.offset = 0
	p.to.typ = int(D_CX_8)
	next = p.link
	progedit_obj8(ctxt, p)
	for p.link != next {
		p = p.link
	}
	if p.from.index == int(D_TLS_8) {
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
func stacksplit_obj8(ctxt *Link, p *Prog, framesize int, noctxt bool, jmpok **Prog) *Prog {
	var q *Prog
	var q1 *Prog
	var arg int
	if ctxt.debugstack != 0 {
		// 8l -K means check not only for stack
		// overflow but stack underflow.
		// On underflow, INT 3 (breakpoint).
		// Underflow itself is rare but this also
		// catches out-of-sync stack guard info.
		p = appendp(ctxt, p)
		p.as = int(ACMPL_8)
		p.from.typ = int(D_INDIR_8 + D_CX_8)
		p.from.offset = 4
		p.to.typ = int(D_SP_8)
		p = appendp(ctxt, p)
		p.as = int(AJCC_8)
		p.to.typ = int(D_BRANCH_8)
		p.to.offset = 4
		q1 = p
		p = appendp(ctxt, p)
		p.as = int(AINT_8)
		p.from.typ = int(D_CONST_8)
		p.from.offset = 3
		p = appendp(ctxt, p)
		p.as = int(ANOP_8)
		q1.pcond = p
	}
	q1 = (*Prog)(nil)
	if framesize <= int(StackSmall_stack) {
		// small stack: SP <= stackguard
		//	CMPL SP, stackguard
		p = appendp(ctxt, p)
		p.as = int(ACMPL_8)
		p.from.typ = int(D_SP_8)
		p.to.typ = int(D_INDIR_8 + D_CX_8)
	} else {
		if framesize <= int(StackBig_stack) {
			// large stack: SP-framesize <= stackguard-StackSmall
			//	LEAL -(framesize-StackSmall)(SP), AX
			//	CMPL AX, stackguard
			p = appendp(ctxt, p)
			p.as = int(ALEAL_8)
			p.from.typ = int(D_INDIR_8 + D_SP_8)
			p.from.offset = -(int64(framesize) - int64(StackSmall_stack))
			p.to.typ = int(D_AX_8)
			p = appendp(ctxt, p)
			p.as = int(ACMPL_8)
			p.from.typ = int(D_AX_8)
			p.to.typ = int(D_INDIR_8 + D_CX_8)
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
			p.as = int(AMOVL_8)
			p.from.typ = int(D_INDIR_8 + D_CX_8)
			p.from.offset = 0
			p.to.typ = int(D_SI_8)
			p = appendp(ctxt, p)
			p.as = int(ACMPL_8)
			p.from.typ = int(D_SI_8)
			p.to.typ = int(D_CONST_8)
			p.to.offset = int64(uint32(StackPreempt_stack & 0xFFFFFFFF))
			p = appendp(ctxt, p)
			p.as = int(AJEQ_8)
			p.to.typ = int(D_BRANCH_8)
			q1 = p
			p = appendp(ctxt, p)
			p.as = int(ALEAL_8)
			p.from.typ = int(D_INDIR_8 + D_SP_8)
			p.from.offset = int64(StackGuard_stack)
			p.to.typ = int(D_AX_8)
			p = appendp(ctxt, p)
			p.as = int(ASUBL_8)
			p.from.typ = int(D_SI_8)
			p.from.offset = 0
			p.to.typ = int(D_AX_8)
			p = appendp(ctxt, p)
			p.as = int(ACMPL_8)
			p.from.typ = int(D_AX_8)
			p.to.typ = int(D_CONST_8)
			p.to.offset = int64(framesize) + (int64(StackGuard_stack) - int64(StackSmall_stack))
		}
	}
	// common
	p = appendp(ctxt, p)
	p.as = int(AJHI_8)
	p.to.typ = int(D_BRANCH_8)
	p.to.offset = 4
	q = p
	p = appendp(ctxt, p) // save frame size in DI
	p.as = int(AMOVL_8)
	p.to.typ = int(D_DI_8)
	p.from.typ = int(D_CONST_8)
	// If we ask for more stack, we'll get a minimum of StackMin bytes.
	// We need a stack frame large enough to hold the top-of-stack data,
	// the function arguments+results, our caller's PC, our frame,
	// a word for the return PC of the next call, and then the StackLimit bytes
	// that must be available on entry to any function called from a function
	// that did a stack check.  If StackMin is enough, don't ask for a specific
	// amount: then we can use the custom functions and save a few
	// instructions.
	if StackTop_stack+int(ctxt.cursym.text.to.offset2)+int(ctxt.arch.ptrsize)+int(framesize)+int(ctxt.arch.ptrsize)+StackLimit_stack >= StackMin_stack {
		p.from.offset = (int64(framesize) + 7) &^ 7
	}
	arg = ctxt.cursym.text.to.offset2
	if arg == 1 { // special marker for known 0
		arg = 0
	}
	if arg&3 != 0 /*untyped*/ {
		ctxt.diag("misaligned argument size in stack split")
	}
	p = appendp(ctxt, p) // save arg size in AX
	p.as = int(AMOVL_8)
	p.to.typ = int(D_AX_8)
	p.from.typ = int(D_CONST_8)
	p.from.offset = int64(arg)
	p = appendp(ctxt, p)
	p.as = int(ACALL_8)
	p.to.typ = int(D_BRANCH_8)
	p.to.sym = ctxt.symmorestack[bool2int(noctxt)]
	p = appendp(ctxt, p)
	p.as = int(AJMP_8)
	p.to.typ = int(D_BRANCH_8)
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

func addstacksplit_obj8(ctxt *Link, cursym *LSym) {
	var p *Prog
	var q *Prog
	var autoffset int
	var deltasp int
	var a int
	if ctxt.symmorestack[0] == nil {
		ctxt.symmorestack[0] = linklookup(ctxt, "runtime.morestack", 0)
		ctxt.symmorestack[1] = linklookup(ctxt, "runtime.morestack_noctxt", 0)
	}
	if ctxt.headtype == int(Hplan9) && ctxt.plan9privates == nil {
		ctxt.plan9privates = linklookup(ctxt, "_privates", 0)
	}
	ctxt.cursym = cursym
	if cursym.text == nil || cursym.text.link == nil {
		return
	}
	p = cursym.text
	autoffset = int(p.to.offset)
	if autoffset < 0 {
		autoffset = 0
	}
	cursym.locals = autoffset
	cursym.args = p.to.offset2
	q = (*Prog)(nil)
	if !(p.from.scale&int(NOSPLIT_textflag) != 0) || (p.from.scale&int(WRAPPER_textflag) != 0) {
		p = appendp(ctxt, p)
		p = load_g_cx_obj8(ctxt, p) // load g into CX
	}
	if !(cursym.text.from.scale&int(NOSPLIT_textflag) != 0) {
		p = stacksplit_obj8(ctxt, p, autoffset, !(cursym.text.from.scale&int(NEEDCTXT_textflag) != 0), &q) // emit split check
	}
	if autoffset != 0 {
		p = appendp(ctxt, p)
		p.as = int(AADJSP_8)
		p.from.typ = int(D_CONST_8)
		p.from.offset = int64(autoffset)
		p.spadj = autoffset
	} else {
		// zero-byte stack adjustment.
		// Insert a fake non-zero adjustment so that stkcheck can
		// recognize the end of the stack-splitting prolog.
		p = appendp(ctxt, p)
		p.as = int(ANOP_8)
		p.spadj = -ctxt.arch.ptrsize
		p = appendp(ctxt, p)
		p.as = int(ANOP_8)
		p.spadj = ctxt.arch.ptrsize
	}
	if q != nil {
		q.pcond = p
	}
	deltasp = autoffset
	if cursym.text.from.scale&int(WRAPPER_textflag) != 0 {
		// g->panicwrap += autoffset + ctxt->arch->ptrsize;
		p = appendp(ctxt, p)
		p.as = int(AADDL_8)
		p.from.typ = int(D_CONST_8)
		p.from.offset = int64(autoffset) + int64(ctxt.arch.ptrsize)
		p.to.typ = int(D_INDIR_8 + D_CX_8)
		p.to.offset = 2 * int64(ctxt.arch.ptrsize)
	}
	if ctxt.debugzerostack != 0 && autoffset != 0 && !(cursym.text.from.scale&int(NOSPLIT_textflag) != 0) {
		// 8l -Z means zero the stack frame on entry.
		// This slows down function calls but can help avoid
		// false positives in garbage collection.
		p = appendp(ctxt, p)
		p.as = int(AMOVL_8)
		p.from.typ = int(D_SP_8)
		p.to.typ = int(D_DI_8)
		p = appendp(ctxt, p)
		p.as = int(AMOVL_8)
		p.from.typ = int(D_CONST_8)
		p.from.offset = int64(autoffset) / 4
		p.to.typ = int(D_CX_8)
		p = appendp(ctxt, p)
		p.as = int(AMOVL_8)
		p.from.typ = int(D_CONST_8)
		p.from.offset = 0
		p.to.typ = int(D_AX_8)
		p = appendp(ctxt, p)
		p.as = int(AREP_8)
		p = appendp(ctxt, p)
		p.as = int(ASTOSL_8)
	}
	for ; p != nil; p = p.link {
		a = p.from.typ
		if a == int(D_AUTO_8) {
			p.from.offset += int64(deltasp)
		}
		if a == int(D_PARAM_8) {
			p.from.offset += int64(deltasp) + 4
		}
		a = p.to.typ
		if a == int(D_AUTO_8) {
			p.to.offset += int64(deltasp)
		}
		if a == int(D_PARAM_8) {
			p.to.offset += int64(deltasp) + 4
		}
		switch p.as {
		default:
			continue
		case APUSHL_8:
		case APUSHFL_8:
			deltasp += 4
			p.spadj = 4
			continue
		case APUSHW_8:
		case APUSHFW_8:
			deltasp += 2
			p.spadj = 2
			continue
		case APOPL_8:
		case APOPFL_8:
			deltasp -= 4
			p.spadj = -4
			continue
		case APOPW_8:
		case APOPFW_8:
			deltasp -= 2
			p.spadj = -2
			continue
		case ARET_8:
			break
		}
		if autoffset != deltasp {
			ctxt.diag("unbalanced PUSH/POP")
		}
		if cursym.text.from.scale&int(WRAPPER_textflag) != 0 {
			p = load_g_cx_obj8(ctxt, p)
			p = appendp(ctxt, p)
			// g->panicwrap -= autoffset + ctxt->arch->ptrsize;
			p.as = int(ASUBL_8)
			p.from.typ = int(D_CONST_8)
			p.from.offset = int64(autoffset) + int64(ctxt.arch.ptrsize)
			p.to.typ = int(D_INDIR_8 + D_CX_8)
			p.to.offset = 2 * int64(ctxt.arch.ptrsize)
			p = appendp(ctxt, p)
			p.as = int(ARET_8)
		}
		if autoffset != 0 {
			p.as = int(AADJSP_8)
			p.from.typ = int(D_CONST_8)
			p.from.offset = int64(-autoffset)
			p.spadj = -autoffset
			p = appendp(ctxt, p)
			p.as = int(ARET_8)
			// If there are instructions following
			// this ARET, they come from a branch
			// with the same stackframe, so undo
			// the cleanup.
			p.spadj = +autoffset
		}
		if p.to.sym != nil { // retjmp
			p.as = int(AJMP_8)
		}
	}
}

func xfol_obj8(ctxt *Link, p *Prog, last **Prog) {
	var q *Prog
	var i int
	var a int
loop:
	if p == nil {
		return
	}
	if p.as == int(AJMP_8) {
		q = p.pcond
		if (q) != nil && q.as != int(ATEXT_8) {
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
			if a == int(ANOP_8) {
				i--
				continue
			}
			if nofollow_obj8(a) != 0 || pushpop_obj8(a) != 0 {
				break // NOTE(rsc): arm does goto copy
			}
			if q.pcond == nil || q.pcond.mark != 0 {
				continue
			}
			if a == int(ACALL_8) || a == int(ALOOP_8) {
				continue
			}
			for {
				if p.as == int(ANOP_8) {
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
		q = ctxt.arch.prg()
		q.as = int(AJMP_8)
		q.lineno = p.lineno
		q.to.typ = int(D_BRANCH_8)
		q.to.offset = int64(p.pc)
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
	if p.pcond != nil && a != int(ACALL_8) {
		/*
		 * some kind of conditional branch.
		 * recurse to follow one path.
		 * continue loop on the other.
		 */
		q = brchain(ctxt, p.pcond)
		if (q) != nil {
			p.pcond = q
		}
		q = brchain(ctxt, p.link)
		if (q) != nil {
			p.link = q
		}
		if p.from.typ == int(D_CONST_8) {
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
				if a != int(ALOOP_8) {
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

func follow_obj8(ctxt *Link, s *LSym) {
	var firstp *Prog
	var lastp *Prog
	ctxt.cursym = s
	firstp = ctxt.arch.prg()
	lastp = firstp
	xfol_obj8(ctxt, s.text, &lastp)
	lastp.link = (*Prog)(nil)
	s.text = firstp.link
}

func nofollow_obj8(a int) int {
	switch a {
	case AJMP_8:
	case ARET_8:
	case AIRETL_8:
	case AIRETW_8:
	case AUNDEF_8:
		return 1
	}
	return 0
}

func pushpop_obj8(a int) int {
	switch a {
	case APUSHL_8:
	case APUSHFL_8:
	case APUSHW_8:
	case APUSHFW_8:
	case APOPL_8:
	case APOPFL_8:
	case APOPW_8:
	case APOPFW_8:
		return 1
	}
	return 0
}

func relinv_obj8(a int) int {
	switch a {
	case AJEQ_8:
		return int(AJNE_8)
	case AJNE_8:
		return int(AJEQ_8)
	case AJLE_8:
		return int(AJGT_8)
	case AJLS_8:
		return int(AJHI_8)
	case AJLT_8:
		return int(AJGE_8)
	case AJMI_8:
		return int(AJPL_8)
	case AJGE_8:
		return int(AJLT_8)
	case AJPL_8:
		return int(AJMI_8)
	case AJGT_8:
		return int(AJLE_8)
	case AJHI_8:
		return int(AJLS_8)
	case AJCS_8:
		return int(AJCC_8)
	case AJCC_8:
		return int(AJCS_8)
	case AJPS_8:
		return int(AJPC_8)
	case AJPC_8:
		return int(AJPS_8)
	case AJOS_8:
		return int(AJOC_8)
	case AJOC_8:
		return int(AJOS_8)
	}
	sysfatal("unknown relation: %s", anames8[a])
	return 0
}
