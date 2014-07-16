package main

import (
	"fmt"
	"math"
)

var linkamd64 = LinkArch{
	name:    "amd64",
	thechar: '6',

	addstacksplit: addstacksplit_obj6,
	assemble:      span6,
	datasize:      datasize_obj6,
	follow:        follow_obj6,
	iscall:        iscall_obj6,
	isdata:        isdata_obj6,
	prg:           prg_obj6,
	progedit:      progedit_obj6,
	settextflag:   settextflag_obj6,
	symtype:       symtype_obj6,
	textflag:      textflag_obj6,
	Pconv:         Pconv_list6,

	minlc:   1,
	ptrsize: 8,
	regsize: 8,

	D_ADDR:   D_ADDR_6,
	D_AUTO:   D_AUTO_6,
	D_BRANCH: D_BRANCH_6,
	D_CONST:  D_CONST_6,
	D_EXTERN: D_EXTERN_6,
	D_FCONST: D_FCONST_6,
	D_NONE:   D_NONE_6,
	D_PARAM:  D_PARAM_6,
	D_SCONST: D_SCONST_6,
	D_STATIC: D_STATIC_6,

	ACALL:     ACALL_6,
	ADATA:     ADATA_6,
	AEND:      AEND_6,
	AFUNCDATA: AFUNCDATA_6,
	AGLOBL:    AGLOBL_6,
	AJMP:      AJMP_6,
	ANOP:      ANOP_6,
	APCDATA:   APCDATA_6,
	ARET:      ARET_6,
	ATEXT:     ATEXT_6,
	ATYPE:     ATYPE_6,
	AUSEFIELD: AUSEFIELD_6,
}

var linkamd64p32 LinkArch

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
var zprg_obj6 = Prog{
	back: 2,
	as:   AGOK_6,
	from: Addr{
		typ:   D_NONE_6,
		index: D_NONE_6,
	},
	to: Addr{
		typ:   D_NONE_6,
		index: D_NONE_6,
	},
}

func nopout_obj6(p *Prog) {
	p.as = int(ANOP_6)
	p.from.typ = int(D_NONE_6)
	p.to.typ = int(D_NONE_6)
}

func symtype_obj6(a *Addr) int {
	var t int
	t = a.typ
	if t == int(D_ADDR_6) {
		t = a.index
	}
	return t
}

func isdata_obj6(p *Prog) bool {
	return p.as == int(ADATA_6) || p.as == int(AGLOBL_6)
}

func iscall_obj6(p *Prog) bool {
	return p.as == int(ACALL_6)
}

func datasize_obj6(p *Prog) int {
	return p.from.scale
}

func textflag_obj6(p *Prog) int {
	return p.from.scale
}

func settextflag_obj6(p *Prog, f int) {
	p.from.scale = f
}

func nacladdr_obj6(ctxt *Link, p *Prog, a *Addr) {
	if p.as == int(ALEAL_6) || p.as == int(ALEAQ_6) {
		return
	}
	if a.typ == int(D_BP_6) || a.typ == int(D_INDIR_6+D_BP_6) {
		ctxt.diag("invalid address: %v", ctxt.Pconv(p))
		return
	}
	if a.typ == int(D_INDIR_6+D_TLS_6) {
		a.typ = int(D_INDIR_6 + D_BP_6)
	} else {
		if a.typ == int(D_TLS_6) {
			a.typ = int(D_BP_6)
		}
	}
	if D_INDIR_6 <= int(a.typ) && a.typ <= int(D_INDIR_6+D_INDIR_6) {
		switch a.typ {
		// all ok
		case D_INDIR_6 + D_BP_6:
		case D_INDIR_6 + D_SP_6:
		case D_INDIR_6 + int(D_R15_6):
			break
		default:
			if a.index != int(D_NONE_6) {
				ctxt.diag("invalid address %v", ctxt.Pconv(p))
			}
			a.index = a.typ - int(D_INDIR_6)
			if a.index != int(D_NONE_6) {
				a.scale = 1
			}
			a.typ = int(D_INDIR_6 + int(D_R15_6))
			break
		}
	}
}

func canuselocaltls_obj6(ctxt *Link) int {
	switch ctxt.headtype {
	case Hplan9:
	case Hwindows:
		return 0
	}
	return 1
}

func progedit_obj6(ctxt *Link, p *Prog) {
	var literal string
	var s *LSym
	var q *Prog
	// Thread-local storage references use the TLS pseudo-register.
	// As a register, TLS refers to the thread-local storage base, and it
	// can only be loaded into another register:
	//
	//         MOVQ TLS, AX
	//
	// An offset from the thread-local storage base is written off(reg)(TLS*1).
	// Semantically it is off(reg), but the (TLS*1) annotation marks this as
	// indexing from the loaded TLS base. This emits a relocation so that
	// if the linker needs to adjust the offset, it can. For example:
	//
	//         MOVQ TLS, AX
	//         MOVQ 8(AX)(TLS*1), CX // load m into CX
	//
	// On systems that support direct access to the TLS memory, this
	// pair of instructions can be reduced to a direct TLS memory reference:
	//
	//         MOVQ 8(TLS), CX // load m into CX
	//
	// The 2-instruction and 1-instruction forms correspond roughly to
	// ELF TLS initial exec mode and ELF TLS local exec mode, respectively.
	//
	// We applies this rewrite on systems that support the 1-instruction form.
	// The decision is made using only the operating system (and probably
	// the -shared flag, eventually), not the link mode. If some link modes
	// on a particular operating system require the 2-instruction form,
	// then all builds for that operating system will use the 2-instruction
	// form, so that the link mode decision can be delayed to link time.
	//
	// In this way, all supported systems use identical instructions to
	// access TLS, and they are rewritten appropriately first here in
	// liblink and then finally using relocations in the linker.
	if canuselocaltls_obj6(ctxt) != 0 {
		// Reduce TLS initial exec model to TLS local exec model.
		// Sequences like
		//	MOVQ TLS, BX
		//	... off(BX)(TLS*1) ...
		// become
		//	NOP
		//	... off(TLS) ...
		//
		// TODO(rsc): Remove the Hsolaris special case. It exists only to
		// guarantee we are producing byte-identical binaries as before this code.
		// But it should be unnecessary.
		if (p.as == int(AMOVQ_6) || p.as == int(AMOVL_6)) && p.from.typ == int(D_TLS_6) && D_AX_6 <= int(p.to.typ) && p.to.typ <= int(D_R15_6) && ctxt.headtype != int(Hsolaris) {
			nopout_obj6(p)
		}
		if p.from.index == int(D_TLS_6) && D_INDIR_6+D_AX_6 <= int(p.from.typ) && p.from.typ <= int(D_INDIR_6+int(D_R15_6)) {
			p.from.typ = int(D_INDIR_6 + D_TLS_6)
			p.from.scale = 0
			p.from.index = int(D_NONE_6)
		}
		if p.to.index == int(D_TLS_6) && D_INDIR_6+D_AX_6 <= int(p.to.typ) && p.to.typ <= int(D_INDIR_6+int(D_R15_6)) {
			p.to.typ = int(D_INDIR_6 + D_TLS_6)
			p.to.scale = 0
			p.to.index = int(D_NONE_6)
		}
	} else {
		// As a courtesy to the C compilers, rewrite TLS local exec load as TLS initial exec load.
		// The instruction
		//	MOVQ off(TLS), BX
		// becomes the sequence
		//	MOVQ TLS, BX
		//	MOVQ off(BX)(TLS*1), BX
		// This allows the C compilers to emit references to m and g using the direct off(TLS) form.
		if (p.as == int(AMOVQ_6) || p.as == int(AMOVL_6)) && p.from.typ == int(D_INDIR_6+D_TLS_6) && D_AX_6 <= int(p.to.typ) && p.to.typ <= int(D_R15_6) {
			q = appendp(ctxt, p)
			q.as = p.as
			q.from = p.from
			q.from.typ = int(D_INDIR_6 + int(p.to.typ))
			q.from.index = int(D_TLS_6)
			q.from.scale = 2 // TODO: use 1
			q.to = p.to
			p.from.typ = int(D_TLS_6)
			p.from.index = int(D_NONE_6)
			p.from.offset = 0
		}
	}
	// TODO: Remove.
	if ctxt.headtype == int(Hwindows) || ctxt.headtype == int(Hplan9) {
		if p.from.scale == 1 && p.from.index == int(D_TLS_6) {
			p.from.scale = 2
		}
		if p.to.scale == 1 && p.to.index == int(D_TLS_6) {
			p.to.scale = 2
		}
	}
	if ctxt.headtype == int(Hnacl) {
		nacladdr_obj6(ctxt, p, &p.from)
		nacladdr_obj6(ctxt, p, &p.to)
	}
	// Maintain information about code generation mode.
	if ctxt.mode == 0 {
		ctxt.mode = 64
	}
	p.mode = ctxt.mode
	switch p.as {
	case AMODE_6:
		if p.from.typ == int(D_CONST_6) || p.from.typ == int(D_INDIR_6+D_NONE_6) {
			switch int(p.from.offset) {
			case 16:
			case 32:
			case 64:
				ctxt.mode = int(p.from.offset)
				break
			}
		}
		nopout_obj6(p)
		break
	}
	// Rewrite CALL/JMP/RET to symbol as D_BRANCH.
	switch p.as {
	case ACALL_6:
	case AJMP_6:
	case ARET_6:
		if (p.to.typ == int(D_EXTERN_6) || p.to.typ == int(D_STATIC_6)) && p.to.sym != nil {
			p.to.typ = int(D_BRANCH_6)
		}
		break
	}
	// Rewrite float constants to values stored in memory.
	switch p.as {
	case AFMOVF_6:
	case AFADDF_6:
	case AFSUBF_6:
	case AFSUBRF_6:
	case AFMULF_6:
	case AFDIVF_6:
	case AFDIVRF_6:
	case AFCOMF_6:
	case AFCOMFP_6:
	case AMOVSS_6:
	case AADDSS_6:
	case ASUBSS_6:
	case AMULSS_6:
	case ADIVSS_6:
	case ACOMISS_6:
	case AUCOMISS_6:
		if p.from.typ == int(D_FCONST_6) {
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
			p.from.typ = int(D_EXTERN_6)
			p.from.sym = s
			p.from.offset = 0
		}
		break
	case AFMOVD_6:
	case AFADDD_6:
	case AFSUBD_6:
	case AFSUBRD_6:
	case AFMULD_6:
	case AFDIVD_6:
	case AFDIVRD_6:
	case AFCOMD_6:
	case AFCOMDP_6:
	case AMOVSD_6:
	case AADDSD_6:
	case ASUBSD_6:
	case AMULSD_6:
	case ADIVSD_6:
	case ACOMISD_6:
	case AUCOMISD_6:
		if p.from.typ == int(D_FCONST_6) {
			var i64 int64
			i64 = int64(math.Float64bits(p.from.u.dval))
			literal = fmt.Sprintf("$f64.%016x", uint64(i64))
			s = linklookup(ctxt, string(literal), 0)
			if s.typ == 0 {
				s.typ = int(SRODATA)
				adduint64(ctxt, s, uint64(i64))
				s.reachable = 0
			}
			p.from.typ = int(D_EXTERN_6)
			p.from.sym = s
			p.from.offset = 0
		}
		break
	}
}

var morename_obj6 = []string{
	"runtime.morestack00",
	"runtime.morestack00_noctxt",
	"runtime.morestack10",
	"runtime.morestack10_noctxt",
	"runtime.morestack01",
	"runtime.morestack01_noctxt",
	"runtime.morestack11",
	"runtime.morestack11_noctxt",
	"runtime.morestack8",
	"runtime.morestack8_noctxt",
	"runtime.morestack16",
	"runtime.morestack16_noctxt",
	"runtime.morestack24",
	"runtime.morestack24_noctxt",
	"runtime.morestack32",
	"runtime.morestack32_noctxt",
	"runtime.morestack40",
	"runtime.morestack40_noctxt",
	"runtime.morestack48",
	"runtime.morestack48_noctxt",
}

// Append code to p to load g into cx.
// Overwrites p with the first instruction (no first appendp).
// Overwriting p is unusual but it lets use this in both the
// prologue (caller must call appendp first) and in the epilogue.
// Returns last new instruction.
func load_g_cx_obj6(ctxt *Link, p *Prog) *Prog {
	var next *Prog
	p.as = int(AMOVQ_6)
	if ctxt.arch.ptrsize == 4 {
		p.as = int(AMOVL_6)
	}
	p.from.typ = int(D_INDIR_6 + D_TLS_6)
	p.from.offset = 0
	p.to.typ = int(D_CX_6)
	next = p.link
	progedit_obj6(ctxt, p)
	for p.link != next {
		p = p.link
	}
	if p.from.index == int(D_TLS_6) {
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
func stacksplit_obj6(ctxt *Link, p *Prog, framesize int, textarg int32, noctxt bool, jmpok **Prog) *Prog {
	var q *Prog
	var q1 *Prog
	var moreconst1 int
	var moreconst2 uint32
	var i uint32
	var cmp int
	var lea int
	var mov int
	var sub int
	cmp = int(ACMPQ_6)
	lea = int(ALEAQ_6)
	mov = int(AMOVQ_6)
	sub = int(ASUBQ_6)
	if ctxt.headtype == int(Hnacl) {
		cmp = int(ACMPL_6)
		lea = int(ALEAL_6)
		mov = int(AMOVL_6)
		sub = int(ASUBL_6)
	}
	if ctxt.debugstack != 0 {
		// 6l -K means check not only for stack
		// overflow but stack underflow.
		// On underflow, INT 3 (breakpoint).
		// Underflow itself is rare but this also
		// catches out-of-sync stack guard info
		p = appendp(ctxt, p)
		p.as = cmp
		indir_cx_obj6(ctxt, &p.from)
		p.from.offset = 8
		p.to.typ = int(D_SP_6)
		p = appendp(ctxt, p)
		p.as = int(AJHI_6)
		p.to.typ = int(D_BRANCH_6)
		p.to.offset = 4
		q1 = p
		p = appendp(ctxt, p)
		p.as = int(AINT_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = 3
		p = appendp(ctxt, p)
		p.as = int(ANOP_6)
		q1.pcond = p
	}
	q1 = (*Prog)(nil)
	if framesize <= int(StackSmall_stack) {
		// small stack: SP <= stackguard
		//	CMPQ SP, stackguard
		p = appendp(ctxt, p)
		p.as = cmp
		p.from.typ = int(D_SP_6)
		indir_cx_obj6(ctxt, &p.to)
	} else {
		if framesize <= int(StackBig_stack) {
			// large stack: SP-framesize <= stackguard-StackSmall
			//	LEAQ -xxx(SP), AX
			//	CMPQ AX, stackguard
			p = appendp(ctxt, p)
			p.as = lea
			p.from.typ = int(D_INDIR_6 + D_SP_6)
			p.from.offset = -(int64(framesize) - int64(StackSmall_stack))
			p.to.typ = int(D_AX_6)
			p = appendp(ctxt, p)
			p.as = cmp
			p.from.typ = int(D_AX_6)
			indir_cx_obj6(ctxt, &p.to)
		} else {
			// Such a large stack we need to protect against wraparound.
			// If SP is close to zero:
			//	SP-stackguard+StackGuard <= framesize + (StackGuard-StackSmall)
			// The +StackGuard on both sides is required to keep the left side positive:
			// SP is allowed to be slightly below stackguard. See stack.h.
			//
			// Preemption sets stackguard to StackPreempt, a very large value.
			// That breaks the math above, so we have to check for that explicitly.
			//	MOVQ	stackguard, CX
			//	CMPQ	CX, $StackPreempt
			//	JEQ	label-of-call-to-morestack
			//	LEAQ	StackGuard(SP), AX
			//	SUBQ	CX, AX
			//	CMPQ	AX, $(framesize+(StackGuard-StackSmall))
			p = appendp(ctxt, p)
			p.as = mov
			indir_cx_obj6(ctxt, &p.from)
			p.from.offset = 0
			p.to.typ = int(D_SI_6)
			p = appendp(ctxt, p)
			p.as = cmp
			p.from.typ = int(D_SI_6)
			p.to.typ = int(D_CONST_6)
			p.to.offset = int64(StackPreempt_stack)
			p = appendp(ctxt, p)
			p.as = int(AJEQ_6)
			p.to.typ = int(D_BRANCH_6)
			q1 = p
			p = appendp(ctxt, p)
			p.as = lea
			p.from.typ = int(D_INDIR_6 + D_SP_6)
			p.from.offset = int64(StackGuard_stack)
			p.to.typ = int(D_AX_6)
			p = appendp(ctxt, p)
			p.as = sub
			p.from.typ = int(D_SI_6)
			p.to.typ = int(D_AX_6)
			p = appendp(ctxt, p)
			p.as = cmp
			p.from.typ = int(D_AX_6)
			p.to.typ = int(D_CONST_6)
			p.to.offset = int64(framesize) + (int64(StackGuard_stack) - int64(StackSmall_stack))
		}
	}
	// common
	p = appendp(ctxt, p)
	p.as = int(AJHI_6)
	p.to.typ = int(D_BRANCH_6)
	q = p
	// If we ask for more stack, we'll get a minimum of StackMin bytes.
	// We need a stack frame large enough to hold the top-of-stack data,
	// the function arguments+results, our caller's PC, our frame,
	// a word for the return PC of the next call, and then the StackLimit bytes
	// that must be available on entry to any function called from a function
	// that did a stack check.  If StackMin is enough, don't ask for a specific
	// amount: then we can use the custom functions and save a few
	// instructions.
	moreconst1 = 0
	if int32(StackTop_stack)+textarg+int32(ctxt.arch.ptrsize)+int32(framesize)+int32(ctxt.arch.ptrsize)+int32(StackLimit_stack) >= int32(StackMin_stack) {
		moreconst1 = framesize
	}
	moreconst2 = uint32(textarg)
	if moreconst2 == 1 { // special marker
		moreconst2 = 0
	}
	if (moreconst2 & 7) != 0 {
		ctxt.diag("misaligned argument size in stack split")
	}
	// 4 varieties varieties (const1==0 cross const2==0)
	// and 6 subvarieties of (const1==0 and const2!=0)
	p = appendp(ctxt, p)
	if moreconst1 == 0 && moreconst2 == 0 {
		p.as = int(ACALL_6)
		p.to.typ = int(D_BRANCH_6)
		p.to.sym = ctxt.symmorestack[0*2+bool2int(noctxt)]
	} else {
		if moreconst1 != 0 && moreconst2 == 0 {
			p.as = int(AMOVL_6)
			p.from.typ = int(D_CONST_6)
			p.from.offset = int64(moreconst1)
			p.to.typ = int(D_AX_6)
			p = appendp(ctxt, p)
			p.as = int(ACALL_6)
			p.to.typ = int(D_BRANCH_6)
			p.to.sym = ctxt.symmorestack[1*2+bool2int(noctxt)]
		} else {
			if moreconst1 == 0 && moreconst2 <= 48 && moreconst2%8 == 0 {
				i = moreconst2/8 + 3
				p.as = int(ACALL_6)
				p.to.typ = int(D_BRANCH_6)
				p.to.sym = ctxt.symmorestack[i*2+uint32(bool2int(noctxt))]
			} else {
				if moreconst1 == 0 && moreconst2 != 0 {
					p.as = int(AMOVL_6)
					p.from.typ = int(D_CONST_6)
					p.from.offset = int64(moreconst2)
					p.to.typ = int(D_AX_6)
					p = appendp(ctxt, p)
					p.as = int(ACALL_6)
					p.to.typ = int(D_BRANCH_6)
					p.to.sym = ctxt.symmorestack[2*2+bool2int(noctxt)]
				} else {
					// Pass framesize and argsize.
					p.as = int(AMOVQ_6)
					p.from.typ = int(D_CONST_6)
					p.from.offset = int64(uint64(moreconst2) << 32)
					p.from.offset |= int64(moreconst1)
					p.to.typ = int(D_AX_6)
					p = appendp(ctxt, p)
					p.as = int(ACALL_6)
					p.to.typ = int(D_BRANCH_6)
					p.to.sym = ctxt.symmorestack[3*2+bool2int(noctxt)]
				}
			}
		}
	}
	p = appendp(ctxt, p)
	p.as = int(AJMP_6)
	p.to.typ = int(D_BRANCH_6)
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

func indir_cx_obj6(ctxt *Link, a *Addr) {
	if ctxt.headtype == int(Hnacl) {
		a.typ = int(D_INDIR_6 + int(D_R15_6))
		a.index = int(D_CX_6)
		a.scale = 1
		return
	}
	a.typ = int(D_INDIR_6 + D_CX_6)
}

func parsetextconst_obj6(arg int64, textstksiz *int64, textarg *int64) {
	*textstksiz = arg & 0xffffffff
	if *textstksiz&0x80000000 != 0 /*untyped*/ {
		*textstksiz = -(-*textstksiz & 0xffffffff)
	}
	*textarg = (arg >> 32) & 0xffffffff
	if *textarg&0x80000000 != 0 /*untyped*/ {
		*textarg = 0
	}
	*textarg = (*textarg + 7) &^ 7
}

func addstacksplit_obj6(ctxt *Link, cursym *LSym) {
	var p *Prog
	var q *Prog
	var q1 *Prog
	var autoffset int
	var deltasp int
	var a int
	var pcsize int
	var i int
	var textstksiz int64
	var textarg int64
	if ctxt.tlsg == nil {
		ctxt.tlsg = linklookup(ctxt, "runtime.tlsg", 0)
	}
	if ctxt.symmorestack[0] == nil {
		if len(morename_obj6) > len(ctxt.symmorestack) {
			sysfatal("Link.symmorestack needs at least %d elements", len(morename_obj6))
		}
		for i = 0; i < len(morename_obj6); i++ {
			ctxt.symmorestack[i] = linklookup(ctxt, morename_obj6[i], 0)
		}
	}
	if ctxt.headtype == int(Hplan9) && ctxt.plan9privates == nil {
		ctxt.plan9privates = linklookup(ctxt, "_privates", 0)
	}
	ctxt.cursym = cursym
	if cursym.text == nil || cursym.text.link == nil {
		return
	}
	p = cursym.text
	parsetextconst_obj6(p.to.offset, (*int64)(&textstksiz), &textarg)
	autoffset = int(textstksiz)
	if autoffset < 0 {
		autoffset = 0
	}
	cursym.args = int(p.to.offset >> 32)
	cursym.locals = int(textstksiz)
	if autoffset < int(StackSmall_stack) && !(p.from.scale&int(NOSPLIT_textflag) != 0) {
		for q = p; q != nil; q = q.link {
			if q.as == int(ACALL_6) {
				goto noleaf
			}
			if (q.as == int(ADUFFCOPY_6) || q.as == int(ADUFFZERO_6)) && autoffset >= StackSmall_stack-8 {
				goto noleaf
			}
		}
		p.from.scale |= int(NOSPLIT_textflag)
	noleaf:
	}
	q = (*Prog)(nil)
	if !(p.from.scale&int(NOSPLIT_textflag) != 0) || (p.from.scale&int(WRAPPER_textflag) != 0) {
		p = appendp(ctxt, p)
		p = load_g_cx_obj6(ctxt, p) // load g into CX
	}
	if !(cursym.text.from.scale&int(NOSPLIT_textflag) != 0) {
		p = stacksplit_obj6(ctxt, p, autoffset, int32(textarg), !(cursym.text.from.scale&int(NEEDCTXT_textflag) != 0), &q) // emit split check
	}
	if autoffset != 0 {
		if autoffset%ctxt.arch.regsize != 0 {
			ctxt.diag("unaligned stack size %d", autoffset)
		}
		p = appendp(ctxt, p)
		p.as = int(AADJSP_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = int64(autoffset)
		p.spadj = autoffset
	} else {
		// zero-byte stack adjustment.
		// Insert a fake non-zero adjustment so that stkcheck can
		// recognize the end of the stack-splitting prolog.
		p = appendp(ctxt, p)
		p.as = int(ANOP_6)
		p.spadj = -ctxt.arch.ptrsize
		p = appendp(ctxt, p)
		p.as = int(ANOP_6)
		p.spadj = ctxt.arch.ptrsize
	}
	if q != nil {
		q.pcond = p
	}
	deltasp = autoffset
	if cursym.text.from.scale&int(WRAPPER_textflag) != 0 {
		// g->panicwrap += autoffset + ctxt->arch->regsize;
		p = appendp(ctxt, p)
		p.as = int(AADDL_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = int64(autoffset) + int64(ctxt.arch.regsize)
		indir_cx_obj6(ctxt, &p.to)
		p.to.offset = 2 * int64(ctxt.arch.ptrsize)
	}
	if ctxt.debugstack > 1 && autoffset != 0 {
		// 6l -K -K means double-check for stack overflow
		// even after calling morestack and even if the
		// function is marked as nosplit.
		p = appendp(ctxt, p)
		p.as = int(AMOVQ_6)
		indir_cx_obj6(ctxt, &p.from)
		p.from.offset = 0
		p.to.typ = int(D_BX_6)
		p = appendp(ctxt, p)
		p.as = int(ASUBQ_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = int64(StackSmall_stack) + 32
		p.to.typ = int(D_BX_6)
		p = appendp(ctxt, p)
		p.as = int(ACMPQ_6)
		p.from.typ = int(D_SP_6)
		p.to.typ = int(D_BX_6)
		p = appendp(ctxt, p)
		p.as = int(AJHI_6)
		p.to.typ = int(D_BRANCH_6)
		q1 = p
		p = appendp(ctxt, p)
		p.as = int(AINT_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = 3
		p = appendp(ctxt, p)
		p.as = int(ANOP_6)
		q1.pcond = p
	}
	if ctxt.debugzerostack != 0 && autoffset != 0 && !(cursym.text.from.scale&int(NOSPLIT_textflag) != 0) {
		// 6l -Z means zero the stack frame on entry.
		// This slows down function calls but can help avoid
		// false positives in garbage collection.
		p = appendp(ctxt, p)
		p.as = int(AMOVQ_6)
		p.from.typ = int(D_SP_6)
		p.to.typ = int(D_DI_6)
		p = appendp(ctxt, p)
		p.as = int(AMOVQ_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = int64(autoffset) / 8
		p.to.typ = int(D_CX_6)
		p = appendp(ctxt, p)
		p.as = int(AMOVQ_6)
		p.from.typ = int(D_CONST_6)
		p.from.offset = 0
		p.to.typ = int(D_AX_6)
		p = appendp(ctxt, p)
		p.as = int(AREP_6)
		p = appendp(ctxt, p)
		p.as = int(ASTOSQ_6)
	}
	for ; p != nil; p = p.link {
		pcsize = p.mode / 8
		a = p.from.typ
		if a == int(D_AUTO_6) {
			p.from.offset += int64(deltasp)
		}
		if a == int(D_PARAM_6) {
			p.from.offset += int64(deltasp) + int64(pcsize)
		}
		a = p.to.typ
		if a == int(D_AUTO_6) {
			p.to.offset += int64(deltasp)
		}
		if a == int(D_PARAM_6) {
			p.to.offset += int64(deltasp) + int64(pcsize)
		}
		switch p.as {
		default:
			continue
		case APUSHL_6:
		case APUSHFL_6:
			deltasp += 4
			p.spadj = 4
			continue
		case APUSHQ_6:
		case APUSHFQ_6:
			deltasp += 8
			p.spadj = 8
			continue
		case APUSHW_6:
		case APUSHFW_6:
			deltasp += 2
			p.spadj = 2
			continue
		case APOPL_6:
		case APOPFL_6:
			deltasp -= 4
			p.spadj = -4
			continue
		case APOPQ_6:
		case APOPFQ_6:
			deltasp -= 8
			p.spadj = -8
			continue
		case APOPW_6:
		case APOPFW_6:
			deltasp -= 2
			p.spadj = -2
			continue
		case ARET_6:
			break
		}
		if autoffset != deltasp {
			ctxt.diag("unbalanced PUSH/POP")
		}
		if cursym.text.from.scale&int(WRAPPER_textflag) != 0 {
			p = load_g_cx_obj6(ctxt, p)
			p = appendp(ctxt, p)
			// g->panicwrap -= autoffset + ctxt->arch->regsize;
			p.as = int(ASUBL_6)
			p.from.typ = int(D_CONST_6)
			p.from.offset = int64(autoffset) + int64(ctxt.arch.regsize)
			indir_cx_obj6(ctxt, &p.to)
			p.to.offset = 2 * int64(ctxt.arch.ptrsize)
			p = appendp(ctxt, p)
			p.as = int(ARET_6)
		}
		if autoffset != 0 {
			p.as = int(AADJSP_6)
			p.from.typ = int(D_CONST_6)
			p.from.offset = int64(-autoffset)
			p.spadj = -autoffset
			p = appendp(ctxt, p)
			p.as = int(ARET_6)
			// If there are instructions following
			// this ARET, they come from a branch
			// with the same stackframe, so undo
			// the cleanup.
			p.spadj = +autoffset
		}
		if p.to.sym != nil { // retjmp
			p.as = int(AJMP_6)
		}
	}
}

func xfol_obj6(ctxt *Link, p *Prog, last **Prog) {
	var q *Prog
	var i int
	var a int
loop:
	if p == nil {
		return
	}
	if p.as == int(AJMP_6) {
		q = p.pcond
		if (q) != nil && q.as != int(ATEXT_6) {
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
			if a == int(ANOP_6) {
				i--
				continue
			}
			if nofollow_obj6(a) != 0 || pushpop_obj6(a) != 0 {
				break // NOTE(rsc): arm does goto copy
			}
			if q.pcond == nil || q.pcond.mark != 0 {
				continue
			}
			if a == int(ACALL_6) || a == int(ALOOP_6) {
				continue
			}
			for {
				if p.as == int(ANOP_6) {
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
				q.as = relinv_obj6(q.as)
				p = q.pcond
				q.pcond = q.link
				q.link = p
				xfol_obj6(ctxt, q.link, last)
				p = q.link
				if p.mark != 0 {
					return
				}
				goto loop /* */
			}
		}
		q = ctxt.arch.prg()
		q.as = int(AJMP_6)
		q.lineno = p.lineno
		q.to.typ = int(D_BRANCH_6)
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
	if nofollow_obj6(a) != 0 {
		return
	}
	if p.pcond != nil && a != int(ACALL_6) {
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
		if p.from.typ == int(D_CONST_6) {
			if p.from.offset == 1 {
				/*
				 * expect conditional jump to be taken.
				 * rewrite so that's the fall-through case.
				 */
				p.as = relinv_obj6(a)
				q = p.link
				p.link = p.pcond
				p.pcond = q
			}
		} else {
			q = p.link
			if q.mark != 0 {
				if a != int(ALOOP_6) {
					p.as = relinv_obj6(a)
					p.link = p.pcond
					p.pcond = q
				}
			}
		}
		xfol_obj6(ctxt, p.link, last)
		if p.pcond.mark != 0 {
			return
		}
		p = p.pcond
		goto loop
	}
	p = p.link
	goto loop
}

func follow_obj6(ctxt *Link, s *LSym) {
	var firstp *Prog
	var lastp *Prog
	ctxt.cursym = s
	firstp = ctxt.arch.prg()
	lastp = firstp
	xfol_obj6(ctxt, s.text, &lastp)
	lastp.link = (*Prog)(nil)
	s.text = firstp.link
}

func nofollow_obj6(a int) int {
	switch a {
	case AJMP_6:
	case ARET_6:
	case AIRETL_6:
	case AIRETQ_6:
	case AIRETW_6:
	case ARETFL_6:
	case ARETFQ_6:
	case ARETFW_6:
	case AUNDEF_6:
		return 1
	}
	return 0
}

func pushpop_obj6(a int) int {
	switch a {
	case APUSHL_6:
	case APUSHFL_6:
	case APUSHQ_6:
	case APUSHFQ_6:
	case APUSHW_6:
	case APUSHFW_6:
	case APOPL_6:
	case APOPFL_6:
	case APOPQ_6:
	case APOPFQ_6:
	case APOPW_6:
	case APOPFW_6:
		return 1
	}
	return 0
}

func relinv_obj6(a int) int {
	switch a {
	case AJEQ_6:
		return int(AJNE_6)
	case AJNE_6:
		return int(AJEQ_6)
	case AJLE_6:
		return int(AJGT_6)
	case AJLS_6:
		return int(AJHI_6)
	case AJLT_6:
		return int(AJGE_6)
	case AJMI_6:
		return int(AJPL_6)
	case AJGE_6:
		return int(AJLT_6)
	case AJPL_6:
		return int(AJMI_6)
	case AJGT_6:
		return int(AJLE_6)
	case AJHI_6:
		return int(AJLS_6)
	case AJCS_6:
		return int(AJCC_6)
	case AJCC_6:
		return int(AJCS_6)
	case AJPS_6:
		return int(AJPC_6)
	case AJPC_6:
		return int(AJPS_6)
	case AJOS_6:
		return int(AJOC_6)
	case AJOC_6:
		return int(AJOS_6)
	}
	sysfatal("unknown relation: %s", anames6[a])
	return 0
}

func prg_obj6() *Prog {
	var p *Prog
	p = new(Prog)
	*p = zprg_obj6
	return p
}
