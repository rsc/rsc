package main

func span6(ctxt *Link, s *LSym) {
	var p *Prog
	var q *Prog
	var c int
	var v int32
	var loop int32
	var bp []uint8
	var n int
	var m int
	var i int
	ctxt.cursym = s
	if s.p != nil {
		return
	}
	if ycover_asm6[0] == 0 {
		instinit_asm6()
	}
	for p = ctxt.cursym.text; p != nil; p = p.link {
		n = 0
		if p.to.typ == int(D_BRANCH_6) {
			if p.pcond == nil {
				p.pcond = p
			}
		}
		q = p.pcond
		if (q) != nil {
			if q.back != 2 {
				n = 1
			}
		}
		p.back = n
		if p.as == int(AADJSP_6) {
			p.to.typ = int(D_SP_6)
			v = int32(-p.from.offset)
			p.from.offset = int64(v)
			p.as = spadjop_asm6(ctxt, p, int(AADDL_6), int(AADDQ_6))
			if v < 0 {
				p.as = spadjop_asm6(ctxt, p, int(ASUBL_6), int(ASUBQ_6))
				v = -v
				p.from.offset = int64(v)
			}
			if v == 0 {
				p.as = int(ANOP_6)
			}
		}
	}
	for p = s.text; p != nil; p = p.link {
		p.back = 2 // use short branches first time through
		q = p.pcond
		if (q) != nil && (q.back&2 != 0 /*untyped*/) {
			p.back |= 1 // backward jump
			q.back |= 4 // loop head
		}
		if p.as == int(AADJSP_6) {
			p.to.typ = int(D_SP_6)
			v = int32(-p.from.offset)
			p.from.offset = int64(v)
			p.as = spadjop_asm6(ctxt, p, int(AADDL_6), int(AADDQ_6))
			if v < 0 {
				p.as = spadjop_asm6(ctxt, p, int(ASUBL_6), int(ASUBQ_6))
				v = -v
				p.from.offset = int64(v)
			}
			if v == 0 {
				p.as = int(ANOP_6)
			}
		}
	}
	n = 0
	for {
		loop = 0
		s.r = s.r[:0]
		s.p = s.p[:0]
		c = 0
		for p = s.text; p != nil; p = p.link {
			if ctxt.headtype == int(Hnacl) && p.isize > 0 {
				var deferreturn_asm6 *LSym
				if deferreturn_asm6 == nil {
					deferreturn_asm6 = linklookup(ctxt, "runtime.deferreturn", 0)
				}
				// pad everything to avoid crossing 32-byte boundary
				if (c >> 5) != ((c + p.isize - 1) >> 5) {
					c = int(naclpad_asm6(ctxt, s, int32(c), -c&31))
				}
				// pad call deferreturn to start at 32-byte boundary
				// so that subtracting 5 in jmpdefer will jump back
				// to that boundary and rerun the call.
				if p.as == int(ACALL_6) && p.to.sym == deferreturn_asm6 {
					c = int(naclpad_asm6(ctxt, s, int32(c), -c&31))
				}
				// pad call to end at 32-byte boundary
				if p.as == int(ACALL_6) {
					c = int(naclpad_asm6(ctxt, s, int32(c), -(c+p.isize)&31))
				}
				// the linker treats REP and STOSQ as different instructions
				// but in fact the REP is a prefix on the STOSQ.
				// make sure REP has room for 2 more bytes, so that
				// padding will not be inserted before the next instruction.
				if (p.as == int(AREP_6) || p.as == int(AREPN_6)) && (c>>5) != ((c+3-1)>>5) {
					c = int(naclpad_asm6(ctxt, s, int32(c), -c&31))
				}
				// same for LOCK.
				// various instructions follow; the longest is 4 bytes.
				// give ourselves 8 bytes so as to avoid surprises.
				if p.as == int(ALOCK_6) && (c>>5) != ((c+8-1)>>5) {
					c = int(naclpad_asm6(ctxt, s, int32(c), -c&31))
				}
			}
			if (p.back&4 != 0 /*untyped*/) && (c&(LoopAlign_asm6-1)) != 0 {
				// pad with NOPs
				v = int32(-c) & (int32(LoopAlign_asm6) - 1)
				if v <= int32(MaxLoopPad_asm6) {
					symgrow(ctxt, s, int64(c)+int64(v))
					fillnop_asm6(s.p[c:], int(v))
					c += int(v)
				}
			}
			p.pc = c
			// process forward jumps to p
			for q = p.comefrom; q != nil; q = q.forwd {
				v = int32(p.pc) - (int32(q.pc) + int32(q.mark))
				if q.back&2 != 0 /*untyped*/ { // short
					if v > 127 {
						loop++
						q.back ^= 2
					}
					if q.as == int(AJCXZL_6) {
						s.p[q.pc+2] = uint8(v)
					} else {
						s.p[q.pc+1] = uint8(v)
					}
				} else {
					bp = s.p[q.pc+q.mark-4:]
					bp[0] = uint8(v)
					bp = bp[1:]
					bp[0] = uint8(v >> 8)
					bp = bp[1:]
					bp[0] = uint8(v >> 16)
					bp = bp[1:]
					bp[0] = uint8(v >> 24)
				}
			}
			p.comefrom = (*Prog)(nil)
			p.pc = c
			asmins_asm6(ctxt, p)
			m = -cap(ctxt.andptr) + cap(ctxt.and)
			if p.isize != m {
				p.isize = m
				loop++
			}
			symgrow(ctxt, s, int64(p.pc)+int64(m))
			copy(s.p[p.pc:], ctxt.and[:m])
			p.mark = m
			c += m
		}
		n++
		if n > 20 {
			ctxt.diag("span must be looping")
			sysfatal("loop")
		}
		if !(loop != 0) {
			break
		}
	}
	if ctxt.headtype == int(Hnacl) {
		c = int(naclpad_asm6(ctxt, s, int32(c), -c&31))
	}
	c += -c & int(FuncAlign_asm6-1)
	s.size = c
	if false { /* debug['a'] > 1 */
		print("span1 %s %lld (%d tries)\n %.6ux", s.name, s.size, n, 0)
		for i = range s.p {
			print(" %.2ux", s.p[i])
			if i%16 == 15 {
				print("\n  %.6ux", i+1)
			}
		}
		if i%16 != 0 /*untyped*/ {
			print("\n")
		}
		for i = range s.r {
			var r *Reloc
			r = &s.r[i]
			print(" rel %#.4ux/%d %s%+lld\n", r.off, r.siz, r.sym.name, r.add)
		}
	}
}

/*
 * this is the ranlib header
 */
const (
	MaxAlign_asm6   = 32
	LoopAlign_asm6  = 16
	MaxLoopPad_asm6 = 0
	FuncAlign_asm6  = 16
)

type Optab_asm6 struct {
	as     int
	ytab   []uint8
	prefix int
	op     [23]uint8
}

type Movtab_asm6 struct {
	as   int
	ft   uint8
	tt   uint8
	code uint8
	op   [4]uint8
}

const (
	Yxxx_asm6 = 0 + iota
	Ynone_asm6
	Yi0_asm6
	Yi1_asm6
	Yi8_asm6
	Ys32_asm6
	Yi32_asm6
	Yi64_asm6
	Yiauto_asm6
	Yal_asm6
	Ycl_asm6
	Yax_asm6
	Ycx_asm6
	Yrb_asm6
	Yrl_asm6
	Yrf_asm6
	Yf0_asm6
	Yrx_asm6
	Ymb_asm6
	Yml_asm6
	Ym_asm6
	Ybr_asm6
	Ycol_asm6
	Ycs_asm6
	Yss_asm6
	Yds_asm6
	Yes_asm6
	Yfs_asm6
	Ygs_asm6
	Ygdtr_asm6
	Yidtr_asm6
	Yldtr_asm6
	Ymsw_asm6
	Ytask_asm6
	Ycr0_asm6
	Ycr1_asm6
	Ycr2_asm6
	Ycr3_asm6
	Ycr4_asm6
	Ycr5_asm6
	Ycr6_asm6
	Ycr7_asm6
	Ycr8_asm6
	Ydr0_asm6
	Ydr1_asm6
	Ydr2_asm6
	Ydr3_asm6
	Ydr4_asm6
	Ydr5_asm6
	Ydr6_asm6
	Ydr7_asm6
	Ytr0_asm6
	Ytr1_asm6
	Ytr2_asm6
	Ytr3_asm6
	Ytr4_asm6
	Ytr5_asm6
	Ytr6_asm6
	Ytr7_asm6
	Yrl32_asm6
	Yrl64_asm6
	Ymr_asm6
	Ymm_asm6
	Yxr_asm6
	Yxm_asm6
	Ytls_asm6
	Ymax_asm6
	Zxxx_asm6 = 0 + iota - 67
	Zlit_asm6
	Zlitm_r_asm6
	Z_rp_asm6
	Zbr_asm6
	Zcall_asm6
	Zcallindreg_asm6
	Zib__asm6
	Zib_rp_asm6
	Zibo_m_asm6
	Zibo_m_xm_asm6
	Zil__asm6
	Zil_rp_asm6
	Ziq_rp_asm6
	Zilo_m_asm6
	Ziqo_m_asm6
	Zjmp_asm6
	Zloop_asm6
	Zo_iw_asm6
	Zm_o_asm6
	Zm_r_asm6
	Zm2_r_asm6
	Zm_r_xm_asm6
	Zm_r_i_xm_asm6
	Zm_r_3d_asm6
	Zm_r_xm_nr_asm6
	Zr_m_xm_nr_asm6
	Zibm_r_asm6
	Zmb_r_asm6
	Zaut_r_asm6
	Zo_m_asm6
	Zo_m64_asm6
	Zpseudo_asm6
	Zr_m_asm6
	Zr_m_xm_asm6
	Zr_m_i_xm_asm6
	Zrp__asm6
	Z_ib_asm6
	Z_il_asm6
	Zm_ibo_asm6
	Zm_ilo_asm6
	Zib_rr_asm6
	Zil_rr_asm6
	Zclr_asm6
	Zbyte_asm6
	Zmax_asm6
	Px_asm6     = 0
	P32_asm6    = 0x32
	Pe_asm6     = 0x66
	Pm_asm6     = 0x0f
	Pq_asm6     = 0xff
	Pb_asm6     = 0xfe
	Pf2_asm6    = 0xf2
	Pf3_asm6    = 0xf3
	Pq3_asm6    = 0x67
	Pw_asm6     = 0x48
	Py_asm6     = 0x80
	Rxf_asm6    = 1 << 9
	Rxt_asm6    = 1 << 8
	Rxw_asm6    = 1 << 3
	Rxr_asm6    = 1 << 2
	Rxx_asm6    = 1 << 1
	Rxb_asm6    = 1 << 0
	Maxand_asm6 = 10
)

var ycover_asm6 [Ymax_asm6 * Ymax_asm6]int8

var reg_asm6 [D_NONE_6]int

var regrex_asm6 [D_NONE_6 + 1]int

func asmins_asm6(ctxt *Link, p *Prog) {
	var n int
	var np int
	var c int
	var and0 []uint8
	var r *Reloc
	ctxt.andptr = ctxt.and[:]
	ctxt.asmode = p.mode
	if p.as == int(AUSEFIELD_6) {
		r = addrel(ctxt.cursym)
		r.off = 0
		r.siz = 0
		r.sym = p.from.sym
		r.typ = int(R_USEFIELD)
		return
	}
	if ctxt.headtype == int(Hnacl) {
		if p.as == int(AREP_6) {
			ctxt.rep++
			return
		}
		if p.as == int(AREPN_6) {
			ctxt.repn++
			return
		}
		if p.as == int(ALOCK_6) {
			ctxt.lock++
			return
		}
		if p.as != int(ALEAQ_6) && p.as != int(ALEAL_6) {
			if p.from.index != int(D_NONE_6) && p.from.scale > 0 {
				nacltrunc_asm6(ctxt, p.from.index)
			}
			if p.to.index != int(D_NONE_6) && p.to.scale > 0 {
				nacltrunc_asm6(ctxt, p.to.index)
			}
		}
		switch p.as {
		case ARET_6:
			copy(ctxt.andptr, naclret_asm6)
			ctxt.andptr = ctxt.andptr[len(naclret_asm6):]
			return
		case ACALL_6, AJMP_6:
			if D_AX_6 <= int(p.to.typ) && p.to.typ <= int(D_DI_6) {
				// ANDL $~31, reg
				ctxt.andptr[0] = 0x83
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(0xe0 | (p.to.typ - int(D_AX_6)))
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0xe0
				ctxt.andptr = ctxt.andptr[1:]
				// ADDQ R15, reg
				ctxt.andptr[0] = 0x4c
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0x01
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(0xf8 | (p.to.typ - int(D_AX_6)))
				ctxt.andptr = ctxt.andptr[1:]
			}
			if D_R8_6 <= int(p.to.typ) && p.to.typ <= int(D_R15_6) {
				// ANDL $~31, reg
				ctxt.andptr[0] = 0x41
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0x83
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(0xe0 | (p.to.typ - int(D_R8_6)))
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0xe0
				ctxt.andptr = ctxt.andptr[1:]
				// ADDQ R15, reg
				ctxt.andptr[0] = 0x4d
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0x01
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(0xf8 | (p.to.typ - int(D_R8_6)))
				ctxt.andptr = ctxt.andptr[1:]
			}
			break
		case AINT_6:
			ctxt.andptr[0] = 0xf4
			ctxt.andptr = ctxt.andptr[1:]
			return
		case ASCASB_6, ASCASW_6, ASCASL_6, ASCASQ_6, ASTOSB_6, ASTOSW_6, ASTOSL_6, ASTOSQ_6:
			copy(ctxt.andptr, naclstos_asm6)
			ctxt.andptr = ctxt.andptr[len(naclstos_asm6):]
			break
		case AMOVSB_6, AMOVSW_6, AMOVSL_6, AMOVSQ_6:
			copy(ctxt.andptr, naclmovs_asm6)
			ctxt.andptr = ctxt.andptr[len(naclmovs_asm6):]
			break
		}
		if ctxt.rep != 0 {
			ctxt.andptr[0] = 0xf3
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.rep = 0
		}
		if ctxt.repn != 0 {
			ctxt.andptr[0] = 0xf2
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.repn = 0
		}
		if ctxt.lock != 0 {
			ctxt.andptr[0] = 0xf0
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.lock = 0
		}
	}
	ctxt.rexflag = 0
	and0 = ctxt.andptr
	ctxt.asmode = p.mode
	doasm_asm6(ctxt, p)
	if ctxt.rexflag != 0 {
		/*
		 * as befits the whole approach of the architecture,
		 * the rex prefix must appear before the first opcode byte
		 * (and thus after any 66/67/f2/f3/26/2e/3e prefix bytes, but
		 * before the 0f opcode escape!), or it might be ignored.
		 * note that the handbook often misleadingly shows 66/f2/f3 in `opcode'.
		 */
		if p.mode != 64 {
			ctxt.diag("asmins: illegal in mode %d: %P", p.mode, p)
		}
		n = -cap(ctxt.andptr) + cap(and0)
		for np = 0; np < n; np++ {
			c = int(and0[np])
			if c != 0xf2 && c != 0xf3 && (c < 0x64 || c > 0x67) && c != 0x2e && c != 0x3e && c != 0x26 {
				break
			}
		}
		copy(and0[np+1:], and0[np:n])
		and0[np] = uint8(0x40 | ctxt.rexflag)
		ctxt.andptr = ctxt.andptr[1:]
	}
	n = -cap(ctxt.andptr) + cap(ctxt.and)
	for i := len(ctxt.cursym.r) - 1; i >= 0; i-- {
		r := &ctxt.cursym.r[i]
		if r.off < p.pc {
			break
		}
		if ctxt.rexflag != 0 {
			r.off++
		}
		if r.typ == int(R_PCREL) || r.typ == int(R_CALL) {
			r.add -= int64(p.pc) + int64(n) - (int64(r.off) + int64(r.siz))
		}
	}
	if ctxt.headtype == int(Hnacl) && p.as != int(ACMPL_6) && p.as != int(ACMPQ_6) {
		switch p.to.typ {
		case D_SP_6:
			copy(ctxt.andptr, naclspfix_asm6)
			ctxt.andptr = ctxt.andptr[len(naclspfix_asm6):]
			break
		case D_BP_6:
			copy(ctxt.andptr, naclbpfix_asm6)
			ctxt.andptr = ctxt.andptr[len(naclbpfix_asm6):]
			break
		}
	}
}

var ynone_asm6 = []uint8{
	Ynone_asm6,
	Ynone_asm6,
	Zlit_asm6,
	1,
	0,
}

var ytext_asm6 = []uint8{
	Ymb_asm6,
	Yi64_asm6,
	Zpseudo_asm6,
	1,
	0,
}

var ynop_asm6 = []uint8{
	Ynone_asm6,
	Ynone_asm6,
	Zpseudo_asm6,
	0,
	Ynone_asm6,
	Yiauto_asm6,
	Zpseudo_asm6,
	0,
	Ynone_asm6,
	Yml_asm6,
	Zpseudo_asm6,
	0,
	Ynone_asm6,
	Yrf_asm6,
	Zpseudo_asm6,
	0,
	Ynone_asm6,
	Yxr_asm6,
	Zpseudo_asm6,
	0,
	Yiauto_asm6,
	Ynone_asm6,
	Zpseudo_asm6,
	0,
	Yml_asm6,
	Ynone_asm6,
	Zpseudo_asm6,
	0,
	Yrf_asm6,
	Ynone_asm6,
	Zpseudo_asm6,
	0,
	Yxr_asm6,
	Ynone_asm6,
	Zpseudo_asm6,
	1,
	0,
}

var yfuncdata_asm6 = []uint8{
	Yi32_asm6,
	Ym_asm6,
	Zpseudo_asm6,
	0,
	0,
}

var ypcdata_asm6 = []uint8{
	Yi32_asm6,
	Yi32_asm6,
	Zpseudo_asm6,
	0,
	0,
}

var yxorb_asm6 = []uint8{
	Yi32_asm6,
	Yal_asm6,
	Zib__asm6,
	1,
	Yi32_asm6,
	Ymb_asm6,
	Zibo_m_asm6,
	2,
	Yrb_asm6,
	Ymb_asm6,
	Zr_m_asm6,
	1,
	Ymb_asm6,
	Yrb_asm6,
	Zm_r_asm6,
	1,
	0,
}

var yxorl_asm6 = []uint8{
	Yi8_asm6,
	Yml_asm6,
	Zibo_m_asm6,
	2,
	Yi32_asm6,
	Yax_asm6,
	Zil__asm6,
	1,
	Yi32_asm6,
	Yml_asm6,
	Zilo_m_asm6,
	2,
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var yaddl_asm6 = []uint8{
	Yi8_asm6,
	Yml_asm6,
	Zibo_m_asm6,
	2,
	Yi32_asm6,
	Yax_asm6,
	Zil__asm6,
	1,
	Yi32_asm6,
	Yml_asm6,
	Zilo_m_asm6,
	2,
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var yincb_asm6 = []uint8{
	Ynone_asm6,
	Ymb_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yincw_asm6 = []uint8{
	Ynone_asm6,
	Yml_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yincl_asm6 = []uint8{
	Ynone_asm6,
	Yml_asm6,
	Zo_m_asm6,
	2,
	0,
}

var ycmpb_asm6 = []uint8{
	Yal_asm6,
	Yi32_asm6,
	Z_ib_asm6,
	1,
	Ymb_asm6,
	Yi32_asm6,
	Zm_ibo_asm6,
	2,
	Ymb_asm6,
	Yrb_asm6,
	Zm_r_asm6,
	1,
	Yrb_asm6,
	Ymb_asm6,
	Zr_m_asm6,
	1,
	0,
}

var ycmpl_asm6 = []uint8{
	Yml_asm6,
	Yi8_asm6,
	Zm_ibo_asm6,
	2,
	Yax_asm6,
	Yi32_asm6,
	Z_il_asm6,
	1,
	Yml_asm6,
	Yi32_asm6,
	Zm_ilo_asm6,
	2,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	0,
}

var yshb_asm6 = []uint8{
	Yi1_asm6,
	Ymb_asm6,
	Zo_m_asm6,
	2,
	Yi32_asm6,
	Ymb_asm6,
	Zibo_m_asm6,
	2,
	Ycx_asm6,
	Ymb_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yshl_asm6 = []uint8{
	Yi1_asm6,
	Yml_asm6,
	Zo_m_asm6,
	2,
	Yi32_asm6,
	Yml_asm6,
	Zibo_m_asm6,
	2,
	Ycl_asm6,
	Yml_asm6,
	Zo_m_asm6,
	2,
	Ycx_asm6,
	Yml_asm6,
	Zo_m_asm6,
	2,
	0,
}

var ytestb_asm6 = []uint8{
	Yi32_asm6,
	Yal_asm6,
	Zib__asm6,
	1,
	Yi32_asm6,
	Ymb_asm6,
	Zibo_m_asm6,
	2,
	Yrb_asm6,
	Ymb_asm6,
	Zr_m_asm6,
	1,
	Ymb_asm6,
	Yrb_asm6,
	Zm_r_asm6,
	1,
	0,
}

var ytestl_asm6 = []uint8{
	Yi32_asm6,
	Yax_asm6,
	Zil__asm6,
	1,
	Yi32_asm6,
	Yml_asm6,
	Zilo_m_asm6,
	2,
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var ymovb_asm6 = []uint8{
	Yrb_asm6,
	Ymb_asm6,
	Zr_m_asm6,
	1,
	Ymb_asm6,
	Yrb_asm6,
	Zm_r_asm6,
	1,
	Yi32_asm6,
	Yrb_asm6,
	Zib_rp_asm6,
	1,
	Yi32_asm6,
	Ymb_asm6,
	Zibo_m_asm6,
	2,
	0,
}

var ymbs_asm6 = []uint8{
	Ymb_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	0,
}

var ybtl_asm6 = []uint8{
	Yi8_asm6,
	Yml_asm6,
	Zibo_m_asm6,
	2,
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	0,
}

var ymovw_asm6 = []uint8{
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	Yi0_asm6,
	Yrl_asm6,
	Zclr_asm6,
	1,
	Yi32_asm6,
	Yrl_asm6,
	Zil_rp_asm6,
	1,
	Yi32_asm6,
	Yml_asm6,
	Zilo_m_asm6,
	2,
	Yiauto_asm6,
	Yrl_asm6,
	Zaut_r_asm6,
	2,
	0,
}

var ymovl_asm6 = []uint8{
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	Yi0_asm6,
	Yrl_asm6,
	Zclr_asm6,
	1,
	Yi32_asm6,
	Yrl_asm6,
	Zil_rp_asm6,
	1,
	Yi32_asm6,
	Yml_asm6,
	Zilo_m_asm6,
	2,
	Yml_asm6,
	Ymr_asm6,
	Zm_r_xm_asm6,
	1, // MMX MOVD
	Ymr_asm6,
	Yml_asm6,
	Zr_m_xm_asm6,
	1, // MMX MOVD
	Yml_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2, // XMM MOVD (32 bit)
	Yxr_asm6,
	Yml_asm6,
	Zr_m_xm_asm6,
	2, // XMM MOVD (32 bit)
	Yiauto_asm6,
	Yrl_asm6,
	Zaut_r_asm6,
	2,
	0,
}

var yret_asm6 = []uint8{
	Ynone_asm6,
	Ynone_asm6,
	Zo_iw_asm6,
	1,
	Yi32_asm6,
	Ynone_asm6,
	Zo_iw_asm6,
	1,
	0,
}

var ymovq_asm6 = []uint8{
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1, // 0x89
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1, // 0x8b
	Yi0_asm6,
	Yrl_asm6,
	Zclr_asm6,
	1, // 0x31
	Ys32_asm6,
	Yrl_asm6,
	Zilo_m_asm6,
	2, // 32 bit signed 0xc7,(0)
	Yi64_asm6,
	Yrl_asm6,
	Ziq_rp_asm6,
	1, // 0xb8 -- 32/64 bit immediate
	Yi32_asm6,
	Yml_asm6,
	Zilo_m_asm6,
	2, // 0xc7,(0)
	Ym_asm6,
	Ymr_asm6,
	Zm_r_xm_nr_asm6,
	1, // MMX MOVQ (shorter encoding)
	Ymr_asm6,
	Ym_asm6,
	Zr_m_xm_nr_asm6,
	1, // MMX MOVQ
	Ymm_asm6,
	Ymr_asm6,
	Zm_r_xm_asm6,
	1, // MMX MOVD
	Ymr_asm6,
	Ymm_asm6,
	Zr_m_xm_asm6,
	1, // MMX MOVD
	Yxr_asm6,
	Ymr_asm6,
	Zm_r_xm_nr_asm6,
	2, // MOVDQ2Q
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_nr_asm6,
	2, // MOVQ xmm1/m64 -> xmm2
	Yxr_asm6,
	Yxm_asm6,
	Zr_m_xm_nr_asm6,
	2, // MOVQ xmm1 -> xmm2/m64
	Yml_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2, // MOVD xmm load
	Yxr_asm6,
	Yml_asm6,
	Zr_m_xm_asm6,
	2, // MOVD xmm store
	Yiauto_asm6,
	Yrl_asm6,
	Zaut_r_asm6,
	2, // built-in LEAQ
	0,
}

var ym_rl_asm6 = []uint8{
	Ym_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var yrl_m_asm6 = []uint8{
	Yrl_asm6,
	Ym_asm6,
	Zr_m_asm6,
	1,
	0,
}

var ymb_rl_asm6 = []uint8{
	Ymb_asm6,
	Yrl_asm6,
	Zmb_r_asm6,
	1,
	0,
}

var yml_rl_asm6 = []uint8{
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var yrl_ml_asm6 = []uint8{
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	0,
}

var yml_mb_asm6 = []uint8{
	Yrb_asm6,
	Ymb_asm6,
	Zr_m_asm6,
	1,
	Ymb_asm6,
	Yrb_asm6,
	Zm_r_asm6,
	1,
	0,
}

var yrb_mb_asm6 = []uint8{
	Yrb_asm6,
	Ymb_asm6,
	Zr_m_asm6,
	1,
	0,
}

var yxchg_asm6 = []uint8{
	Yax_asm6,
	Yrl_asm6,
	Z_rp_asm6,
	1,
	Yrl_asm6,
	Yax_asm6,
	Zrp__asm6,
	1,
	Yrl_asm6,
	Yml_asm6,
	Zr_m_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var ydivl_asm6 = []uint8{
	Yml_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	0,
}

var ydivb_asm6 = []uint8{
	Ymb_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	0,
}

var yimul_asm6 = []uint8{
	Yml_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	Yi8_asm6,
	Yrl_asm6,
	Zib_rr_asm6,
	1,
	Yi32_asm6,
	Yrl_asm6,
	Zil_rr_asm6,
	1,
	Yml_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	2,
	0,
}

var yimul3_asm6 = []uint8{
	Yml_asm6,
	Yrl_asm6,
	Zibm_r_asm6,
	2,
	0,
}

var ybyte_asm6 = []uint8{
	Yi64_asm6,
	Ynone_asm6,
	Zbyte_asm6,
	1,
	0,
}

var yin_asm6 = []uint8{
	Yi32_asm6,
	Ynone_asm6,
	Zib__asm6,
	1,
	Ynone_asm6,
	Ynone_asm6,
	Zlit_asm6,
	1,
	0,
}

var yint_asm6 = []uint8{
	Yi32_asm6,
	Ynone_asm6,
	Zib__asm6,
	1,
	0,
}

var ypushl_asm6 = []uint8{
	Yrl_asm6,
	Ynone_asm6,
	Zrp__asm6,
	1,
	Ym_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	Yi8_asm6,
	Ynone_asm6,
	Zib__asm6,
	1,
	Yi32_asm6,
	Ynone_asm6,
	Zil__asm6,
	1,
	0,
}

var ypopl_asm6 = []uint8{
	Ynone_asm6,
	Yrl_asm6,
	Z_rp_asm6,
	1,
	Ynone_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	0,
}

var ybswap_asm6 = []uint8{
	Ynone_asm6,
	Yrl_asm6,
	Z_rp_asm6,
	2,
	0,
}

var yscond_asm6 = []uint8{
	Ynone_asm6,
	Ymb_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yjcond_asm6 = []uint8{
	Ynone_asm6,
	Ybr_asm6,
	Zbr_asm6,
	0,
	Yi0_asm6,
	Ybr_asm6,
	Zbr_asm6,
	0,
	Yi1_asm6,
	Ybr_asm6,
	Zbr_asm6,
	1,
	0,
}

var yloop_asm6 = []uint8{
	Ynone_asm6,
	Ybr_asm6,
	Zloop_asm6,
	1,
	0,
}

var ycall_asm6 = []uint8{
	Ynone_asm6,
	Yml_asm6,
	Zcallindreg_asm6,
	0,
	Yrx_asm6,
	Yrx_asm6,
	Zcallindreg_asm6,
	2,
	Ynone_asm6,
	Ybr_asm6,
	Zcall_asm6,
	1,
	0,
}

var yduff_asm6 = []uint8{
	Ynone_asm6,
	Yi32_asm6,
	Zcall_asm6,
	1,
	0,
}

var yjmp_asm6 = []uint8{
	Ynone_asm6,
	Yml_asm6,
	Zo_m64_asm6,
	2,
	Ynone_asm6,
	Ybr_asm6,
	Zjmp_asm6,
	1,
	0,
}

var yfmvd_asm6 = []uint8{
	Ym_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	Yf0_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	Yrf_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	Yf0_asm6,
	Yrf_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yfmvdp_asm6 = []uint8{
	Yf0_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	Yf0_asm6,
	Yrf_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yfmvf_asm6 = []uint8{
	Ym_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	Yf0_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yfmvx_asm6 = []uint8{
	Ym_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	0,
}

var yfmvp_asm6 = []uint8{
	Yf0_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yfadd_asm6 = []uint8{
	Ym_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	Yrf_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	Yf0_asm6,
	Yrf_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yfaddp_asm6 = []uint8{
	Yf0_asm6,
	Yrf_asm6,
	Zo_m_asm6,
	2,
	0,
}

var yfxch_asm6 = []uint8{
	Yf0_asm6,
	Yrf_asm6,
	Zo_m_asm6,
	2,
	Yrf_asm6,
	Yf0_asm6,
	Zm_o_asm6,
	2,
	0,
}

var ycompp_asm6 = []uint8{
	Yf0_asm6,
	Yrf_asm6,
	Zo_m_asm6,
	2, /* botch is really f0,f1 */
	0,
}

var ystsw_asm6 = []uint8{
	Ynone_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	Ynone_asm6,
	Yax_asm6,
	Zlit_asm6,
	1,
	0,
}

var ystcw_asm6 = []uint8{
	Ynone_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	Ym_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	0,
}

var ysvrs_asm6 = []uint8{
	Ynone_asm6,
	Ym_asm6,
	Zo_m_asm6,
	2,
	Ym_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	0,
}

var ymm_asm6 = []uint8{
	Ymm_asm6,
	Ymr_asm6,
	Zm_r_xm_asm6,
	1,
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2,
	0,
}

var yxm_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var yxcvm1_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2,
	Yxm_asm6,
	Ymr_asm6,
	Zm_r_xm_asm6,
	2,
	0,
}

var yxcvm2_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2,
	Ymm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2,
	0,
}

/*
static uchar	yxmq[] =
{
	Yxm,	Yxr,	Zm_r_xm,	2,
	0
};
*/
var yxr_asm6 = []uint8{
	Yxr_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var yxr_ml_asm6 = []uint8{
	Yxr_asm6,
	Yml_asm6,
	Zr_m_xm_asm6,
	1,
	0,
}

var ymr_asm6 = []uint8{
	Ymr_asm6,
	Ymr_asm6,
	Zm_r_asm6,
	1,
	0,
}

var ymr_ml_asm6 = []uint8{
	Ymr_asm6,
	Yml_asm6,
	Zr_m_xm_asm6,
	1,
	0,
}

var yxcmp_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var yxcmpi_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_i_xm_asm6,
	2,
	0,
}

var yxmov_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	1,
	Yxr_asm6,
	Yxm_asm6,
	Zr_m_xm_asm6,
	1,
	0,
}

var yxcvfl_asm6 = []uint8{
	Yxm_asm6,
	Yrl_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var yxcvlf_asm6 = []uint8{
	Yml_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var yxcvfq_asm6 = []uint8{
	Yxm_asm6,
	Yrl_asm6,
	Zm_r_xm_asm6,
	2,
	0,
}

var yxcvqf_asm6 = []uint8{
	Yml_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2,
	0,
}

var yps_asm6 = []uint8{
	Ymm_asm6,
	Ymr_asm6,
	Zm_r_xm_asm6,
	1,
	Yi8_asm6,
	Ymr_asm6,
	Zibo_m_xm_asm6,
	2,
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	2,
	Yi8_asm6,
	Yxr_asm6,
	Zibo_m_xm_asm6,
	3,
	0,
}

var yxrrl_asm6 = []uint8{
	Yxr_asm6,
	Yrl_asm6,
	Zm_r_asm6,
	1,
	0,
}

var ymfp_asm6 = []uint8{
	Ymm_asm6,
	Ymr_asm6,
	Zm_r_3d_asm6,
	1,
	0,
}

var ymrxr_asm6 = []uint8{
	Ymr_asm6,
	Yxr_asm6,
	Zm_r_asm6,
	1,
	Yxm_asm6,
	Yxr_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var ymshuf_asm6 = []uint8{
	Ymm_asm6,
	Ymr_asm6,
	Zibm_r_asm6,
	2,
	0,
}

var ymshufb_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zm2_r_asm6,
	2,
	0,
}

var yxshuf_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zibm_r_asm6,
	2,
	0,
}

var yextrw_asm6 = []uint8{
	Yxr_asm6,
	Yrl_asm6,
	Zibm_r_asm6,
	2,
	0,
}

var yinsrw_asm6 = []uint8{
	Yml_asm6,
	Yxr_asm6,
	Zibm_r_asm6,
	2,
	0,
}

var yinsr_asm6 = []uint8{
	Ymm_asm6,
	Yxr_asm6,
	Zibm_r_asm6,
	3,
	0,
}

var ypsdq_asm6 = []uint8{
	Yi8_asm6,
	Yxr_asm6,
	Zibo_m_asm6,
	2,
	0,
}

var ymskb_asm6 = []uint8{
	Yxr_asm6,
	Yrl_asm6,
	Zm_r_xm_asm6,
	2,
	Ymr_asm6,
	Yrl_asm6,
	Zm_r_xm_asm6,
	1,
	0,
}

var ycrc32l_asm6 = []uint8{Yml_asm6, Yrl_asm6, Zlitm_r_asm6, 0}

var yprefetch_asm6 = []uint8{
	Ym_asm6,
	Ynone_asm6,
	Zm_o_asm6,
	2,
	0,
}

var yaes_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zlitm_r_asm6,
	2,
	0,
}

var yaes2_asm6 = []uint8{
	Yxm_asm6,
	Yxr_asm6,
	Zibm_r_asm6,
	2,
	0,
}

/*
 * You are doasm, holding in your hand a Prog* with p->as set to, say, ACRC32,
 * and p->from and p->to as operands (Addr*).  The linker scans optab to find
 * the entry with the given p->as and then looks through the ytable for that
 * instruction (the second field in the optab struct) for a line whose first
 * two values match the Ytypes of the p->from and p->to operands.  The function
 * oclass in span.c computes the specific Ytype of an operand and then the set
 * of more general Ytypes that it satisfies is implied by the ycover table, set
 * up in instinit.  For example, oclass distinguishes the constants 0 and 1
 * from the more general 8-bit constants, but instinit says
 *
 *        ycover[Yi0*Ymax + Ys32] = 1;
 *        ycover[Yi1*Ymax + Ys32] = 1;
 *        ycover[Yi8*Ymax + Ys32] = 1;
 *
 * which means that Yi0, Yi1, and Yi8 all count as Ys32 (signed 32)
 * if that's what an instruction can handle.
 *
 * In parallel with the scan through the ytable for the appropriate line, there
 * is a z pointer that starts out pointing at the strange magic byte list in
 * the Optab struct.  With each step past a non-matching ytable line, z
 * advances by the 4th entry in the line.  When a matching line is found, that
 * z pointer has the extra data to use in laying down the instruction bytes.
 * The actual bytes laid down are a function of the 3rd entry in the line (that
 * is, the Ztype) and the z bytes.
 *
 * For example, let's look at AADDL.  The optab line says:
 *        { AADDL,        yaddl,  Px, 0x83,(00),0x05,0x81,(00),0x01,0x03 },
 *
 * and yaddl says
 *        uchar   yaddl[] =
 *        {
 *                Yi8,    Yml,    Zibo_m, 2,
 *                Yi32,   Yax,    Zil_,   1,
 *                Yi32,   Yml,    Zilo_m, 2,
 *                Yrl,    Yml,    Zr_m,   1,
 *                Yml,    Yrl,    Zm_r,   1,
 *                0
 *        };
 *
 * so there are 5 possible types of ADDL instruction that can be laid down, and
 * possible states used to lay them down (Ztype and z pointer, assuming z
 * points at {0x83,(00),0x05,0x81,(00),0x01,0x03}) are:
 *
 *        Yi8, Yml -> Zibo_m, z (0x83, 00)
 *        Yi32, Yax -> Zil_, z+2 (0x05)
 *        Yi32, Yml -> Zilo_m, z+2+1 (0x81, 0x00)
 *        Yrl, Yml -> Zr_m, z+2+1+2 (0x01)
 *        Yml, Yrl -> Zm_r, z+2+1+2+1 (0x03)
 *
 * The Pconstant in the optab line controls the prefix bytes to emit.  That's
 * relatively straightforward as this program goes.
 *
 * The switch on t[2] in doasm implements the various Z cases.  Zibo_m, for
 * example, is an opcode byte (z[0]) then an asmando (which is some kind of
 * encoded addressing mode for the Yml arg), and then a single immediate byte.
 * Zilo_m is the same but a long (32-bit) immediate.
 */
var optab_asm6 = /*	as, ytab, andproto, opcode */
[]Optab_asm6{
	{AXXX_6, nil, 0, [23]uint8{}},
	{AAAA_6, ynone_asm6, P32_asm6, [23]uint8{0x37}},
	{AAAD_6, ynone_asm6, P32_asm6, [23]uint8{0xd5, 0x0a}},
	{AAAM_6, ynone_asm6, P32_asm6, [23]uint8{0xd4, 0x0a}},
	{AAAS_6, ynone_asm6, P32_asm6, [23]uint8{0x3f}},
	{AADCB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x14, 0x80, (02), 0x10, 0x10}},
	{AADCL_6, yxorl_asm6, Px_asm6, [23]uint8{0x83, (02), 0x15, 0x81, (02), 0x11, 0x13}},
	{AADCQ_6, yxorl_asm6, Pw_asm6, [23]uint8{0x83, (02), 0x15, 0x81, (02), 0x11, 0x13}},
	{AADCW_6, yxorl_asm6, Pe_asm6, [23]uint8{0x83, (02), 0x15, 0x81, (02), 0x11, 0x13}},
	{AADDB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x04, 0x80, (00), 0x00, 0x02}},
	{AADDL_6, yaddl_asm6, Px_asm6, [23]uint8{0x83, (00), 0x05, 0x81, (00), 0x01, 0x03}},
	{AADDPD_6, yxm_asm6, Pq_asm6, [23]uint8{0x58}},
	{AADDPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x58}},
	{AADDQ_6, yaddl_asm6, Pw_asm6, [23]uint8{0x83, (00), 0x05, 0x81, (00), 0x01, 0x03}},
	{AADDSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x58}},
	{AADDSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x58}},
	{AADDW_6, yaddl_asm6, Pe_asm6, [23]uint8{0x83, (00), 0x05, 0x81, (00), 0x01, 0x03}},
	{AADJSP_6, nil, 0, [23]uint8{}},
	{AANDB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x24, 0x80, (04), 0x20, 0x22}},
	{AANDL_6, yxorl_asm6, Px_asm6, [23]uint8{0x83, (04), 0x25, 0x81, (04), 0x21, 0x23}},
	{AANDNPD_6, yxm_asm6, Pq_asm6, [23]uint8{0x55}},
	{AANDNPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x55}},
	{AANDPD_6, yxm_asm6, Pq_asm6, [23]uint8{0x54}},
	{AANDPS_6, yxm_asm6, Pq_asm6, [23]uint8{0x54}},
	{AANDQ_6, yxorl_asm6, Pw_asm6, [23]uint8{0x83, (04), 0x25, 0x81, (04), 0x21, 0x23}},
	{AANDW_6, yxorl_asm6, Pe_asm6, [23]uint8{0x83, (04), 0x25, 0x81, (04), 0x21, 0x23}},
	{AARPL_6, yrl_ml_asm6, P32_asm6, [23]uint8{0x63}},
	{ABOUNDL_6, yrl_m_asm6, P32_asm6, [23]uint8{0x62}},
	{ABOUNDW_6, yrl_m_asm6, Pe_asm6, [23]uint8{0x62}},
	{ABSFL_6, yml_rl_asm6, Pm_asm6, [23]uint8{0xbc}},
	{ABSFQ_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0xbc}},
	{ABSFW_6, yml_rl_asm6, Pq_asm6, [23]uint8{0xbc}},
	{ABSRL_6, yml_rl_asm6, Pm_asm6, [23]uint8{0xbd}},
	{ABSRQ_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0xbd}},
	{ABSRW_6, yml_rl_asm6, Pq_asm6, [23]uint8{0xbd}},
	{ABSWAPL_6, ybswap_asm6, Px_asm6, [23]uint8{0x0f, 0xc8}},
	{ABSWAPQ_6, ybswap_asm6, Pw_asm6, [23]uint8{0x0f, 0xc8}},
	{ABTCL_6, ybtl_asm6, Pm_asm6, [23]uint8{0xba, (07), 0xbb}},
	{ABTCQ_6, ybtl_asm6, Pw_asm6, [23]uint8{0x0f, 0xba, (07), 0x0f, 0xbb}},
	{ABTCW_6, ybtl_asm6, Pq_asm6, [23]uint8{0xba, (07), 0xbb}},
	{ABTL_6, ybtl_asm6, Pm_asm6, [23]uint8{0xba, (04), 0xa3}},
	{ABTQ_6, ybtl_asm6, Pw_asm6, [23]uint8{0x0f, 0xba, (04), 0x0f, 0xa3}},
	{ABTRL_6, ybtl_asm6, Pm_asm6, [23]uint8{0xba, (06), 0xb3}},
	{ABTRQ_6, ybtl_asm6, Pw_asm6, [23]uint8{0x0f, 0xba, (06), 0x0f, 0xb3}},
	{ABTRW_6, ybtl_asm6, Pq_asm6, [23]uint8{0xba, (06), 0xb3}},
	{ABTSL_6, ybtl_asm6, Pm_asm6, [23]uint8{0xba, (05), 0xab}},
	{ABTSQ_6, ybtl_asm6, Pw_asm6, [23]uint8{0x0f, 0xba, (05), 0x0f, 0xab}},
	{ABTSW_6, ybtl_asm6, Pq_asm6, [23]uint8{0xba, (05), 0xab}},
	{ABTW_6, ybtl_asm6, Pq_asm6, [23]uint8{0xba, (04), 0xa3}},
	{ABYTE_6, ybyte_asm6, Px_asm6, [23]uint8{1}},
	{ACALL_6, ycall_asm6, Px_asm6, [23]uint8{0xff, (02), 0xe8}},
	{ACDQ_6, ynone_asm6, Px_asm6, [23]uint8{0x99}},
	{ACLC_6, ynone_asm6, Px_asm6, [23]uint8{0xf8}},
	{ACLD_6, ynone_asm6, Px_asm6, [23]uint8{0xfc}},
	{ACLI_6, ynone_asm6, Px_asm6, [23]uint8{0xfa}},
	{ACLTS_6, ynone_asm6, Pm_asm6, [23]uint8{0x06}},
	{ACMC_6, ynone_asm6, Px_asm6, [23]uint8{0xf5}},
	{ACMOVLCC_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x43}},
	{ACMOVLCS_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x42}},
	{ACMOVLEQ_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x44}},
	{ACMOVLGE_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x4d}},
	{ACMOVLGT_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x4f}},
	{ACMOVLHI_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x47}},
	{ACMOVLLE_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x4e}},
	{ACMOVLLS_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x46}},
	{ACMOVLLT_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x4c}},
	{ACMOVLMI_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x48}},
	{ACMOVLNE_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x45}},
	{ACMOVLOC_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x41}},
	{ACMOVLOS_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x40}},
	{ACMOVLPC_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x4b}},
	{ACMOVLPL_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x49}},
	{ACMOVLPS_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x4a}},
	{ACMOVQCC_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x43}},
	{ACMOVQCS_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x42}},
	{ACMOVQEQ_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x44}},
	{ACMOVQGE_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x4d}},
	{ACMOVQGT_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x4f}},
	{ACMOVQHI_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x47}},
	{ACMOVQLE_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x4e}},
	{ACMOVQLS_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x46}},
	{ACMOVQLT_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x4c}},
	{ACMOVQMI_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x48}},
	{ACMOVQNE_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x45}},
	{ACMOVQOC_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x41}},
	{ACMOVQOS_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x40}},
	{ACMOVQPC_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x4b}},
	{ACMOVQPL_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x49}},
	{ACMOVQPS_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0x4a}},
	{ACMOVWCC_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x43}},
	{ACMOVWCS_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x42}},
	{ACMOVWEQ_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x44}},
	{ACMOVWGE_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x4d}},
	{ACMOVWGT_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x4f}},
	{ACMOVWHI_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x47}},
	{ACMOVWLE_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x4e}},
	{ACMOVWLS_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x46}},
	{ACMOVWLT_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x4c}},
	{ACMOVWMI_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x48}},
	{ACMOVWNE_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x45}},
	{ACMOVWOC_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x41}},
	{ACMOVWOS_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x40}},
	{ACMOVWPC_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x4b}},
	{ACMOVWPL_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x49}},
	{ACMOVWPS_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x4a}},
	{ACMPB_6, ycmpb_asm6, Pb_asm6, [23]uint8{0x3c, 0x80, (07), 0x38, 0x3a}},
	{ACMPL_6, ycmpl_asm6, Px_asm6, [23]uint8{0x83, (07), 0x3d, 0x81, (07), 0x39, 0x3b}},
	{ACMPPD_6, yxcmpi_asm6, Px_asm6, [23]uint8{Pe_asm6, 0xc2}},
	{ACMPPS_6, yxcmpi_asm6, Pm_asm6, [23]uint8{0xc2, 0}},
	{ACMPQ_6, ycmpl_asm6, Pw_asm6, [23]uint8{0x83, (07), 0x3d, 0x81, (07), 0x39, 0x3b}},
	{ACMPSB_6, ynone_asm6, Pb_asm6, [23]uint8{0xa6}},
	{ACMPSD_6, yxcmpi_asm6, Px_asm6, [23]uint8{Pf2_asm6, 0xc2}},
	{ACMPSL_6, ynone_asm6, Px_asm6, [23]uint8{0xa7}},
	{ACMPSQ_6, ynone_asm6, Pw_asm6, [23]uint8{0xa7}},
	{ACMPSS_6, yxcmpi_asm6, Px_asm6, [23]uint8{Pf3_asm6, 0xc2}},
	{ACMPSW_6, ynone_asm6, Pe_asm6, [23]uint8{0xa7}},
	{ACMPW_6, ycmpl_asm6, Pe_asm6, [23]uint8{0x83, (07), 0x3d, 0x81, (07), 0x39, 0x3b}},
	{ACOMISD_6, yxcmp_asm6, Pe_asm6, [23]uint8{0x2f}},
	{ACOMISS_6, yxcmp_asm6, Pm_asm6, [23]uint8{0x2f}},
	{ACPUID_6, ynone_asm6, Pm_asm6, [23]uint8{0xa2}},
	{ACVTPL2PD_6, yxcvm2_asm6, Px_asm6, [23]uint8{Pf3_asm6, 0xe6, Pe_asm6, 0x2a}},
	{ACVTPL2PS_6, yxcvm2_asm6, Pm_asm6, [23]uint8{0x5b, 0, 0x2a, 0}},
	{ACVTPD2PL_6, yxcvm1_asm6, Px_asm6, [23]uint8{Pf2_asm6, 0xe6, Pe_asm6, 0x2d}},
	{ACVTPD2PS_6, yxm_asm6, Pe_asm6, [23]uint8{0x5a}},
	{ACVTPS2PL_6, yxcvm1_asm6, Px_asm6, [23]uint8{Pe_asm6, 0x5b, Pm_asm6, 0x2d}},
	{ACVTPS2PD_6, yxm_asm6, Pm_asm6, [23]uint8{0x5a}},
	{API2FW_6, ymfp_asm6, Px_asm6, [23]uint8{0x0c}},
	{ACVTSD2SL_6, yxcvfl_asm6, Pf2_asm6, [23]uint8{0x2d}},
	{ACVTSD2SQ_6, yxcvfq_asm6, Pw_asm6, [23]uint8{Pf2_asm6, 0x2d}},
	{ACVTSD2SS_6, yxm_asm6, Pf2_asm6, [23]uint8{0x5a}},
	{ACVTSL2SD_6, yxcvlf_asm6, Pf2_asm6, [23]uint8{0x2a}},
	{ACVTSQ2SD_6, yxcvqf_asm6, Pw_asm6, [23]uint8{Pf2_asm6, 0x2a}},
	{ACVTSL2SS_6, yxcvlf_asm6, Pf3_asm6, [23]uint8{0x2a}},
	{ACVTSQ2SS_6, yxcvqf_asm6, Pw_asm6, [23]uint8{Pf3_asm6, 0x2a}},
	{ACVTSS2SD_6, yxm_asm6, Pf3_asm6, [23]uint8{0x5a}},
	{ACVTSS2SL_6, yxcvfl_asm6, Pf3_asm6, [23]uint8{0x2d}},
	{ACVTSS2SQ_6, yxcvfq_asm6, Pw_asm6, [23]uint8{Pf3_asm6, 0x2d}},
	{ACVTTPD2PL_6, yxcvm1_asm6, Px_asm6, [23]uint8{Pe_asm6, 0xe6, Pe_asm6, 0x2c}},
	{ACVTTPS2PL_6, yxcvm1_asm6, Px_asm6, [23]uint8{Pf3_asm6, 0x5b, Pm_asm6, 0x2c}},
	{ACVTTSD2SL_6, yxcvfl_asm6, Pf2_asm6, [23]uint8{0x2c}},
	{ACVTTSD2SQ_6, yxcvfq_asm6, Pw_asm6, [23]uint8{Pf2_asm6, 0x2c}},
	{ACVTTSS2SL_6, yxcvfl_asm6, Pf3_asm6, [23]uint8{0x2c}},
	{ACVTTSS2SQ_6, yxcvfq_asm6, Pw_asm6, [23]uint8{Pf3_asm6, 0x2c}},
	{ACWD_6, ynone_asm6, Pe_asm6, [23]uint8{0x99}},
	{ACQO_6, ynone_asm6, Pw_asm6, [23]uint8{0x99}},
	{ADAA_6, ynone_asm6, P32_asm6, [23]uint8{0x27}},
	{ADAS_6, ynone_asm6, P32_asm6, [23]uint8{0x2f}},
	{ADATA_6, nil, 0, [23]uint8{}},
	{ADECB_6, yincb_asm6, Pb_asm6, [23]uint8{0xfe, (01)}},
	{ADECL_6, yincl_asm6, Px_asm6, [23]uint8{0xff, (01)}},
	{ADECQ_6, yincl_asm6, Pw_asm6, [23]uint8{0xff, (01)}},
	{ADECW_6, yincw_asm6, Pe_asm6, [23]uint8{0xff, (01)}},
	{ADIVB_6, ydivb_asm6, Pb_asm6, [23]uint8{0xf6, (06)}},
	{ADIVL_6, ydivl_asm6, Px_asm6, [23]uint8{0xf7, (06)}},
	{ADIVPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x5e}},
	{ADIVPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x5e}},
	{ADIVQ_6, ydivl_asm6, Pw_asm6, [23]uint8{0xf7, (06)}},
	{ADIVSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x5e}},
	{ADIVSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x5e}},
	{ADIVW_6, ydivl_asm6, Pe_asm6, [23]uint8{0xf7, (06)}},
	{AEMMS_6, ynone_asm6, Pm_asm6, [23]uint8{0x77}},
	{AENTER_6, nil, 0, [23]uint8{}}, /* botch */
	{AFXRSTOR_6, ysvrs_asm6, Pm_asm6, [23]uint8{0xae, (01), 0xae, (01)}},
	{AFXSAVE_6, ysvrs_asm6, Pm_asm6, [23]uint8{0xae, (00), 0xae, (00)}},
	{AFXRSTOR64_6, ysvrs_asm6, Pw_asm6, [23]uint8{0x0f, 0xae, (01), 0x0f, 0xae, (01)}},
	{AFXSAVE64_6, ysvrs_asm6, Pw_asm6, [23]uint8{0x0f, 0xae, (00), 0x0f, 0xae, (00)}},
	{AGLOBL_6, nil, 0, [23]uint8{}},
	{AGOK_6, nil, 0, [23]uint8{}},
	{AHISTORY_6, nil, 0, [23]uint8{}},
	{AHLT_6, ynone_asm6, Px_asm6, [23]uint8{0xf4}},
	{AIDIVB_6, ydivb_asm6, Pb_asm6, [23]uint8{0xf6, (07)}},
	{AIDIVL_6, ydivl_asm6, Px_asm6, [23]uint8{0xf7, (07)}},
	{AIDIVQ_6, ydivl_asm6, Pw_asm6, [23]uint8{0xf7, (07)}},
	{AIDIVW_6, ydivl_asm6, Pe_asm6, [23]uint8{0xf7, (07)}},
	{AIMULB_6, ydivb_asm6, Pb_asm6, [23]uint8{0xf6, (05)}},
	{AIMULL_6, yimul_asm6, Px_asm6, [23]uint8{0xf7, (05), 0x6b, 0x69, Pm_asm6, 0xaf}},
	{AIMULQ_6, yimul_asm6, Pw_asm6, [23]uint8{0xf7, (05), 0x6b, 0x69, Pm_asm6, 0xaf}},
	{AIMULW_6, yimul_asm6, Pe_asm6, [23]uint8{0xf7, (05), 0x6b, 0x69, Pm_asm6, 0xaf}},
	{AIMUL3Q_6, yimul3_asm6, Pw_asm6, [23]uint8{0x6b, (00)}},
	{AINB_6, yin_asm6, Pb_asm6, [23]uint8{0xe4, 0xec}},
	{AINCB_6, yincb_asm6, Pb_asm6, [23]uint8{0xfe, (00)}},
	{AINCL_6, yincl_asm6, Px_asm6, [23]uint8{0xff, (00)}},
	{AINCQ_6, yincl_asm6, Pw_asm6, [23]uint8{0xff, (00)}},
	{AINCW_6, yincw_asm6, Pe_asm6, [23]uint8{0xff, (00)}},
	{AINL_6, yin_asm6, Px_asm6, [23]uint8{0xe5, 0xed}},
	{AINSB_6, ynone_asm6, Pb_asm6, [23]uint8{0x6c}},
	{AINSL_6, ynone_asm6, Px_asm6, [23]uint8{0x6d}},
	{AINSW_6, ynone_asm6, Pe_asm6, [23]uint8{0x6d}},
	{AINT_6, yint_asm6, Px_asm6, [23]uint8{0xcd}},
	{AINTO_6, ynone_asm6, P32_asm6, [23]uint8{0xce}},
	{AINW_6, yin_asm6, Pe_asm6, [23]uint8{0xe5, 0xed}},
	{AIRETL_6, ynone_asm6, Px_asm6, [23]uint8{0xcf}},
	{AIRETQ_6, ynone_asm6, Pw_asm6, [23]uint8{0xcf}},
	{AIRETW_6, ynone_asm6, Pe_asm6, [23]uint8{0xcf}},
	{AJCC_6, yjcond_asm6, Px_asm6, [23]uint8{0x73, 0x83, (00)}},
	{AJCS_6, yjcond_asm6, Px_asm6, [23]uint8{0x72, 0x82}},
	{AJCXZL_6, yloop_asm6, Px_asm6, [23]uint8{0xe3}},
	{AJCXZQ_6, yloop_asm6, Px_asm6, [23]uint8{0xe3}},
	{AJEQ_6, yjcond_asm6, Px_asm6, [23]uint8{0x74, 0x84}},
	{AJGE_6, yjcond_asm6, Px_asm6, [23]uint8{0x7d, 0x8d}},
	{AJGT_6, yjcond_asm6, Px_asm6, [23]uint8{0x7f, 0x8f}},
	{AJHI_6, yjcond_asm6, Px_asm6, [23]uint8{0x77, 0x87}},
	{AJLE_6, yjcond_asm6, Px_asm6, [23]uint8{0x7e, 0x8e}},
	{AJLS_6, yjcond_asm6, Px_asm6, [23]uint8{0x76, 0x86}},
	{AJLT_6, yjcond_asm6, Px_asm6, [23]uint8{0x7c, 0x8c}},
	{AJMI_6, yjcond_asm6, Px_asm6, [23]uint8{0x78, 0x88}},
	{AJMP_6, yjmp_asm6, Px_asm6, [23]uint8{0xff, (04), 0xeb, 0xe9}},
	{AJNE_6, yjcond_asm6, Px_asm6, [23]uint8{0x75, 0x85}},
	{AJOC_6, yjcond_asm6, Px_asm6, [23]uint8{0x71, 0x81, (00)}},
	{AJOS_6, yjcond_asm6, Px_asm6, [23]uint8{0x70, 0x80, (00)}},
	{AJPC_6, yjcond_asm6, Px_asm6, [23]uint8{0x7b, 0x8b}},
	{AJPL_6, yjcond_asm6, Px_asm6, [23]uint8{0x79, 0x89}},
	{AJPS_6, yjcond_asm6, Px_asm6, [23]uint8{0x7a, 0x8a}},
	{ALAHF_6, ynone_asm6, Px_asm6, [23]uint8{0x9f}},
	{ALARL_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x02}},
	{ALARW_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x02}},
	{ALDMXCSR_6, ysvrs_asm6, Pm_asm6, [23]uint8{0xae, (02), 0xae, (02)}},
	{ALEAL_6, ym_rl_asm6, Px_asm6, [23]uint8{0x8d}},
	{ALEAQ_6, ym_rl_asm6, Pw_asm6, [23]uint8{0x8d}},
	{ALEAVEL_6, ynone_asm6, P32_asm6, [23]uint8{0xc9}},
	{ALEAVEQ_6, ynone_asm6, Py_asm6, [23]uint8{0xc9}},
	{ALEAVEW_6, ynone_asm6, Pe_asm6, [23]uint8{0xc9}},
	{ALEAW_6, ym_rl_asm6, Pe_asm6, [23]uint8{0x8d}},
	{ALOCK_6, ynone_asm6, Px_asm6, [23]uint8{0xf0}},
	{ALODSB_6, ynone_asm6, Pb_asm6, [23]uint8{0xac}},
	{ALODSL_6, ynone_asm6, Px_asm6, [23]uint8{0xad}},
	{ALODSQ_6, ynone_asm6, Pw_asm6, [23]uint8{0xad}},
	{ALODSW_6, ynone_asm6, Pe_asm6, [23]uint8{0xad}},
	{ALONG_6, ybyte_asm6, Px_asm6, [23]uint8{4}},
	{ALOOP_6, yloop_asm6, Px_asm6, [23]uint8{0xe2}},
	{ALOOPEQ_6, yloop_asm6, Px_asm6, [23]uint8{0xe1}},
	{ALOOPNE_6, yloop_asm6, Px_asm6, [23]uint8{0xe0}},
	{ALSLL_6, yml_rl_asm6, Pm_asm6, [23]uint8{0x03}},
	{ALSLW_6, yml_rl_asm6, Pq_asm6, [23]uint8{0x03}},
	{AMASKMOVOU_6, yxr_asm6, Pe_asm6, [23]uint8{0xf7}},
	{AMASKMOVQ_6, ymr_asm6, Pm_asm6, [23]uint8{0xf7}},
	{AMAXPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x5f}},
	{AMAXPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x5f}},
	{AMAXSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x5f}},
	{AMAXSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x5f}},
	{AMINPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x5d}},
	{AMINPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x5d}},
	{AMINSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x5d}},
	{AMINSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x5d}},
	{AMOVAPD_6, yxmov_asm6, Pe_asm6, [23]uint8{0x28, 0x29}},
	{AMOVAPS_6, yxmov_asm6, Pm_asm6, [23]uint8{0x28, 0x29}},
	{AMOVB_6, ymovb_asm6, Pb_asm6, [23]uint8{0x88, 0x8a, 0xb0, 0xc6, (00)}},
	{AMOVBLSX_6, ymb_rl_asm6, Pm_asm6, [23]uint8{0xbe}},
	{AMOVBLZX_6, ymb_rl_asm6, Pm_asm6, [23]uint8{0xb6}},
	{AMOVBQSX_6, ymb_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0xbe}},
	{AMOVBQZX_6, ymb_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0xb6}},
	{AMOVBWSX_6, ymb_rl_asm6, Pq_asm6, [23]uint8{0xbe}},
	{AMOVBWZX_6, ymb_rl_asm6, Pq_asm6, [23]uint8{0xb6}},
	{AMOVO_6, yxmov_asm6, Pe_asm6, [23]uint8{0x6f, 0x7f}},
	{AMOVOU_6, yxmov_asm6, Pf3_asm6, [23]uint8{0x6f, 0x7f}},
	{AMOVHLPS_6, yxr_asm6, Pm_asm6, [23]uint8{0x12}},
	{AMOVHPD_6, yxmov_asm6, Pe_asm6, [23]uint8{0x16, 0x17}},
	{AMOVHPS_6, yxmov_asm6, Pm_asm6, [23]uint8{0x16, 0x17}},
	{AMOVL_6, ymovl_asm6, Px_asm6, [23]uint8{0x89, 0x8b, 0x31, 0xb8, 0xc7, (00), 0x6e, 0x7e, Pe_asm6, 0x6e, Pe_asm6, 0x7e, 0}},
	{AMOVLHPS_6, yxr_asm6, Pm_asm6, [23]uint8{0x16}},
	{AMOVLPD_6, yxmov_asm6, Pe_asm6, [23]uint8{0x12, 0x13}},
	{AMOVLPS_6, yxmov_asm6, Pm_asm6, [23]uint8{0x12, 0x13}},
	{AMOVLQSX_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x63}},
	{AMOVLQZX_6, yml_rl_asm6, Px_asm6, [23]uint8{0x8b}},
	{AMOVMSKPD_6, yxrrl_asm6, Pq_asm6, [23]uint8{0x50}},
	{AMOVMSKPS_6, yxrrl_asm6, Pm_asm6, [23]uint8{0x50}},
	{AMOVNTO_6, yxr_ml_asm6, Pe_asm6, [23]uint8{0xe7}},
	{AMOVNTPD_6, yxr_ml_asm6, Pe_asm6, [23]uint8{0x2b}},
	{AMOVNTPS_6, yxr_ml_asm6, Pm_asm6, [23]uint8{0x2b}},
	{AMOVNTQ_6, ymr_ml_asm6, Pm_asm6, [23]uint8{0xe7}},
	{AMOVQ_6, ymovq_asm6, Pw_asm6, [23]uint8{0x89, 0x8b, 0x31, 0xc7, (00), 0xb8, 0xc7, (00), 0x6f, 0x7f, 0x6e, 0x7e, Pf2_asm6, 0xd6, Pf3_asm6, 0x7e, Pe_asm6, 0xd6, Pe_asm6, 0x6e, Pe_asm6, 0x7e, 0}},
	{AMOVQOZX_6, ymrxr_asm6, Pf3_asm6, [23]uint8{0xd6, 0x7e}},
	{AMOVSB_6, ynone_asm6, Pb_asm6, [23]uint8{0xa4}},
	{AMOVSD_6, yxmov_asm6, Pf2_asm6, [23]uint8{0x10, 0x11}},
	{AMOVSL_6, ynone_asm6, Px_asm6, [23]uint8{0xa5}},
	{AMOVSQ_6, ynone_asm6, Pw_asm6, [23]uint8{0xa5}},
	{AMOVSS_6, yxmov_asm6, Pf3_asm6, [23]uint8{0x10, 0x11}},
	{AMOVSW_6, ynone_asm6, Pe_asm6, [23]uint8{0xa5}},
	{AMOVUPD_6, yxmov_asm6, Pe_asm6, [23]uint8{0x10, 0x11}},
	{AMOVUPS_6, yxmov_asm6, Pm_asm6, [23]uint8{0x10, 0x11}},
	{AMOVW_6, ymovw_asm6, Pe_asm6, [23]uint8{0x89, 0x8b, 0x31, 0xb8, 0xc7, (00), 0}},
	{AMOVWLSX_6, yml_rl_asm6, Pm_asm6, [23]uint8{0xbf}},
	{AMOVWLZX_6, yml_rl_asm6, Pm_asm6, [23]uint8{0xb7}},
	{AMOVWQSX_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0xbf}},
	{AMOVWQZX_6, yml_rl_asm6, Pw_asm6, [23]uint8{0x0f, 0xb7}},
	{AMULB_6, ydivb_asm6, Pb_asm6, [23]uint8{0xf6, (04)}},
	{AMULL_6, ydivl_asm6, Px_asm6, [23]uint8{0xf7, (04)}},
	{AMULPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x59}},
	{AMULPS_6, yxm_asm6, Ym_asm6, [23]uint8{0x59}},
	{AMULQ_6, ydivl_asm6, Pw_asm6, [23]uint8{0xf7, (04)}},
	{AMULSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x59}},
	{AMULSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x59}},
	{AMULW_6, ydivl_asm6, Pe_asm6, [23]uint8{0xf7, (04)}},
	{ANAME_6, nil, 0, [23]uint8{}},
	{ANEGB_6, yscond_asm6, Pb_asm6, [23]uint8{0xf6, (03)}},
	{ANEGL_6, yscond_asm6, Px_asm6, [23]uint8{0xf7, (03)}},
	{ANEGQ_6, yscond_asm6, Pw_asm6, [23]uint8{0xf7, (03)}},
	{ANEGW_6, yscond_asm6, Pe_asm6, [23]uint8{0xf7, (03)}},
	{ANOP_6, ynop_asm6, Px_asm6, [23]uint8{0, 0}},
	{ANOTB_6, yscond_asm6, Pb_asm6, [23]uint8{0xf6, (02)}},
	{ANOTL_6, yscond_asm6, Px_asm6, [23]uint8{0xf7, (02)}},
	{ANOTQ_6, yscond_asm6, Pw_asm6, [23]uint8{0xf7, (02)}},
	{ANOTW_6, yscond_asm6, Pe_asm6, [23]uint8{0xf7, (02)}},
	{AORB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x0c, 0x80, (01), 0x08, 0x0a}},
	{AORL_6, yxorl_asm6, Px_asm6, [23]uint8{0x83, (01), 0x0d, 0x81, (01), 0x09, 0x0b}},
	{AORPD_6, yxm_asm6, Pq_asm6, [23]uint8{0x56}},
	{AORPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x56}},
	{AORQ_6, yxorl_asm6, Pw_asm6, [23]uint8{0x83, (01), 0x0d, 0x81, (01), 0x09, 0x0b}},
	{AORW_6, yxorl_asm6, Pe_asm6, [23]uint8{0x83, (01), 0x0d, 0x81, (01), 0x09, 0x0b}},
	{AOUTB_6, yin_asm6, Pb_asm6, [23]uint8{0xe6, 0xee}},
	{AOUTL_6, yin_asm6, Px_asm6, [23]uint8{0xe7, 0xef}},
	{AOUTSB_6, ynone_asm6, Pb_asm6, [23]uint8{0x6e}},
	{AOUTSL_6, ynone_asm6, Px_asm6, [23]uint8{0x6f}},
	{AOUTSW_6, ynone_asm6, Pe_asm6, [23]uint8{0x6f}},
	{AOUTW_6, yin_asm6, Pe_asm6, [23]uint8{0xe7, 0xef}},
	{APACKSSLW_6, ymm_asm6, Py_asm6, [23]uint8{0x6b, Pe_asm6, 0x6b}},
	{APACKSSWB_6, ymm_asm6, Py_asm6, [23]uint8{0x63, Pe_asm6, 0x63}},
	{APACKUSWB_6, ymm_asm6, Py_asm6, [23]uint8{0x67, Pe_asm6, 0x67}},
	{APADDB_6, ymm_asm6, Py_asm6, [23]uint8{0xfc, Pe_asm6, 0xfc}},
	{APADDL_6, ymm_asm6, Py_asm6, [23]uint8{0xfe, Pe_asm6, 0xfe}},
	{APADDQ_6, yxm_asm6, Pe_asm6, [23]uint8{0xd4}},
	{APADDSB_6, ymm_asm6, Py_asm6, [23]uint8{0xec, Pe_asm6, 0xec}},
	{APADDSW_6, ymm_asm6, Py_asm6, [23]uint8{0xed, Pe_asm6, 0xed}},
	{APADDUSB_6, ymm_asm6, Py_asm6, [23]uint8{0xdc, Pe_asm6, 0xdc}},
	{APADDUSW_6, ymm_asm6, Py_asm6, [23]uint8{0xdd, Pe_asm6, 0xdd}},
	{APADDW_6, ymm_asm6, Py_asm6, [23]uint8{0xfd, Pe_asm6, 0xfd}},
	{APAND_6, ymm_asm6, Py_asm6, [23]uint8{0xdb, Pe_asm6, 0xdb}},
	{APANDN_6, ymm_asm6, Py_asm6, [23]uint8{0xdf, Pe_asm6, 0xdf}},
	{APAUSE_6, ynone_asm6, Px_asm6, [23]uint8{0xf3, 0x90}},
	{APAVGB_6, ymm_asm6, Py_asm6, [23]uint8{0xe0, Pe_asm6, 0xe0}},
	{APAVGW_6, ymm_asm6, Py_asm6, [23]uint8{0xe3, Pe_asm6, 0xe3}},
	{APCMPEQB_6, ymm_asm6, Py_asm6, [23]uint8{0x74, Pe_asm6, 0x74}},
	{APCMPEQL_6, ymm_asm6, Py_asm6, [23]uint8{0x76, Pe_asm6, 0x76}},
	{APCMPEQW_6, ymm_asm6, Py_asm6, [23]uint8{0x75, Pe_asm6, 0x75}},
	{APCMPGTB_6, ymm_asm6, Py_asm6, [23]uint8{0x64, Pe_asm6, 0x64}},
	{APCMPGTL_6, ymm_asm6, Py_asm6, [23]uint8{0x66, Pe_asm6, 0x66}},
	{APCMPGTW_6, ymm_asm6, Py_asm6, [23]uint8{0x65, Pe_asm6, 0x65}},
	{APEXTRW_6, yextrw_asm6, Pq_asm6, [23]uint8{0xc5, (00)}},
	{APF2IL_6, ymfp_asm6, Px_asm6, [23]uint8{0x1d}},
	{APF2IW_6, ymfp_asm6, Px_asm6, [23]uint8{0x1c}},
	{API2FL_6, ymfp_asm6, Px_asm6, [23]uint8{0x0d}},
	{APFACC_6, ymfp_asm6, Px_asm6, [23]uint8{0xae}},
	{APFADD_6, ymfp_asm6, Px_asm6, [23]uint8{0x9e}},
	{APFCMPEQ_6, ymfp_asm6, Px_asm6, [23]uint8{0xb0}},
	{APFCMPGE_6, ymfp_asm6, Px_asm6, [23]uint8{0x90}},
	{APFCMPGT_6, ymfp_asm6, Px_asm6, [23]uint8{0xa0}},
	{APFMAX_6, ymfp_asm6, Px_asm6, [23]uint8{0xa4}},
	{APFMIN_6, ymfp_asm6, Px_asm6, [23]uint8{0x94}},
	{APFMUL_6, ymfp_asm6, Px_asm6, [23]uint8{0xb4}},
	{APFNACC_6, ymfp_asm6, Px_asm6, [23]uint8{0x8a}},
	{APFPNACC_6, ymfp_asm6, Px_asm6, [23]uint8{0x8e}},
	{APFRCP_6, ymfp_asm6, Px_asm6, [23]uint8{0x96}},
	{APFRCPIT1_6, ymfp_asm6, Px_asm6, [23]uint8{0xa6}},
	{APFRCPI2T_6, ymfp_asm6, Px_asm6, [23]uint8{0xb6}},
	{APFRSQIT1_6, ymfp_asm6, Px_asm6, [23]uint8{0xa7}},
	{APFRSQRT_6, ymfp_asm6, Px_asm6, [23]uint8{0x97}},
	{APFSUB_6, ymfp_asm6, Px_asm6, [23]uint8{0x9a}},
	{APFSUBR_6, ymfp_asm6, Px_asm6, [23]uint8{0xaa}},
	{APINSRW_6, yinsrw_asm6, Pq_asm6, [23]uint8{0xc4, (00)}},
	{APINSRD_6, yinsr_asm6, Pq_asm6, [23]uint8{0x3a, 0x22, (00)}},
	{APINSRQ_6, yinsr_asm6, Pq3_asm6, [23]uint8{0x3a, 0x22, (00)}},
	{APMADDWL_6, ymm_asm6, Py_asm6, [23]uint8{0xf5, Pe_asm6, 0xf5}},
	{APMAXSW_6, yxm_asm6, Pe_asm6, [23]uint8{0xee}},
	{APMAXUB_6, yxm_asm6, Pe_asm6, [23]uint8{0xde}},
	{APMINSW_6, yxm_asm6, Pe_asm6, [23]uint8{0xea}},
	{APMINUB_6, yxm_asm6, Pe_asm6, [23]uint8{0xda}},
	{APMOVMSKB_6, ymskb_asm6, Px_asm6, [23]uint8{Pe_asm6, 0xd7, 0xd7}},
	{APMULHRW_6, ymfp_asm6, Px_asm6, [23]uint8{0xb7}},
	{APMULHUW_6, ymm_asm6, Py_asm6, [23]uint8{0xe4, Pe_asm6, 0xe4}},
	{APMULHW_6, ymm_asm6, Py_asm6, [23]uint8{0xe5, Pe_asm6, 0xe5}},
	{APMULLW_6, ymm_asm6, Py_asm6, [23]uint8{0xd5, Pe_asm6, 0xd5}},
	{APMULULQ_6, ymm_asm6, Py_asm6, [23]uint8{0xf4, Pe_asm6, 0xf4}},
	{APOPAL_6, ynone_asm6, P32_asm6, [23]uint8{0x61}},
	{APOPAW_6, ynone_asm6, Pe_asm6, [23]uint8{0x61}},
	{APOPFL_6, ynone_asm6, P32_asm6, [23]uint8{0x9d}},
	{APOPFQ_6, ynone_asm6, Py_asm6, [23]uint8{0x9d}},
	{APOPFW_6, ynone_asm6, Pe_asm6, [23]uint8{0x9d}},
	{APOPL_6, ypopl_asm6, P32_asm6, [23]uint8{0x58, 0x8f, (00)}},
	{APOPQ_6, ypopl_asm6, Py_asm6, [23]uint8{0x58, 0x8f, (00)}},
	{APOPW_6, ypopl_asm6, Pe_asm6, [23]uint8{0x58, 0x8f, (00)}},
	{APOR_6, ymm_asm6, Py_asm6, [23]uint8{0xeb, Pe_asm6, 0xeb}},
	{APSADBW_6, yxm_asm6, Pq_asm6, [23]uint8{0xf6}},
	{APSHUFHW_6, yxshuf_asm6, Pf3_asm6, [23]uint8{0x70, (00)}},
	{APSHUFL_6, yxshuf_asm6, Pq_asm6, [23]uint8{0x70, (00)}},
	{APSHUFLW_6, yxshuf_asm6, Pf2_asm6, [23]uint8{0x70, (00)}},
	{APSHUFW_6, ymshuf_asm6, Pm_asm6, [23]uint8{0x70, (00)}},
	{APSHUFB_6, ymshufb_asm6, Pq_asm6, [23]uint8{0x38, 0x00}},
	{APSLLO_6, ypsdq_asm6, Pq_asm6, [23]uint8{0x73, (07)}},
	{APSLLL_6, yps_asm6, Py_asm6, [23]uint8{0xf2, 0x72, (06), Pe_asm6, 0xf2, Pe_asm6, 0x72, (06)}},
	{APSLLQ_6, yps_asm6, Py_asm6, [23]uint8{0xf3, 0x73, (06), Pe_asm6, 0xf3, Pe_asm6, 0x73, (06)}},
	{APSLLW_6, yps_asm6, Py_asm6, [23]uint8{0xf1, 0x71, (06), Pe_asm6, 0xf1, Pe_asm6, 0x71, (06)}},
	{APSRAL_6, yps_asm6, Py_asm6, [23]uint8{0xe2, 0x72, (04), Pe_asm6, 0xe2, Pe_asm6, 0x72, (04)}},
	{APSRAW_6, yps_asm6, Py_asm6, [23]uint8{0xe1, 0x71, (04), Pe_asm6, 0xe1, Pe_asm6, 0x71, (04)}},
	{APSRLO_6, ypsdq_asm6, Pq_asm6, [23]uint8{0x73, (03)}},
	{APSRLL_6, yps_asm6, Py_asm6, [23]uint8{0xd2, 0x72, (02), Pe_asm6, 0xd2, Pe_asm6, 0x72, (02)}},
	{APSRLQ_6, yps_asm6, Py_asm6, [23]uint8{0xd3, 0x73, (02), Pe_asm6, 0xd3, Pe_asm6, 0x73, (02)}},
	{APSRLW_6, yps_asm6, Py_asm6, [23]uint8{0xd1, 0x71, (02), Pe_asm6, 0xe1, Pe_asm6, 0x71, (02)}},
	{APSUBB_6, yxm_asm6, Pe_asm6, [23]uint8{0xf8}},
	{APSUBL_6, yxm_asm6, Pe_asm6, [23]uint8{0xfa}},
	{APSUBQ_6, yxm_asm6, Pe_asm6, [23]uint8{0xfb}},
	{APSUBSB_6, yxm_asm6, Pe_asm6, [23]uint8{0xe8}},
	{APSUBSW_6, yxm_asm6, Pe_asm6, [23]uint8{0xe9}},
	{APSUBUSB_6, yxm_asm6, Pe_asm6, [23]uint8{0xd8}},
	{APSUBUSW_6, yxm_asm6, Pe_asm6, [23]uint8{0xd9}},
	{APSUBW_6, yxm_asm6, Pe_asm6, [23]uint8{0xf9}},
	{APSWAPL_6, ymfp_asm6, Px_asm6, [23]uint8{0xbb}},
	{APUNPCKHBW_6, ymm_asm6, Py_asm6, [23]uint8{0x68, Pe_asm6, 0x68}},
	{APUNPCKHLQ_6, ymm_asm6, Py_asm6, [23]uint8{0x6a, Pe_asm6, 0x6a}},
	{APUNPCKHQDQ_6, yxm_asm6, Pe_asm6, [23]uint8{0x6d}},
	{APUNPCKHWL_6, ymm_asm6, Py_asm6, [23]uint8{0x69, Pe_asm6, 0x69}},
	{APUNPCKLBW_6, ymm_asm6, Py_asm6, [23]uint8{0x60, Pe_asm6, 0x60}},
	{APUNPCKLLQ_6, ymm_asm6, Py_asm6, [23]uint8{0x62, Pe_asm6, 0x62}},
	{APUNPCKLQDQ_6, yxm_asm6, Pe_asm6, [23]uint8{0x6c}},
	{APUNPCKLWL_6, ymm_asm6, Py_asm6, [23]uint8{0x61, Pe_asm6, 0x61}},
	{APUSHAL_6, ynone_asm6, P32_asm6, [23]uint8{0x60}},
	{APUSHAW_6, ynone_asm6, Pe_asm6, [23]uint8{0x60}},
	{APUSHFL_6, ynone_asm6, P32_asm6, [23]uint8{0x9c}},
	{APUSHFQ_6, ynone_asm6, Py_asm6, [23]uint8{0x9c}},
	{APUSHFW_6, ynone_asm6, Pe_asm6, [23]uint8{0x9c}},
	{APUSHL_6, ypushl_asm6, P32_asm6, [23]uint8{0x50, 0xff, (06), 0x6a, 0x68}},
	{APUSHQ_6, ypushl_asm6, Py_asm6, [23]uint8{0x50, 0xff, (06), 0x6a, 0x68}},
	{APUSHW_6, ypushl_asm6, Pe_asm6, [23]uint8{0x50, 0xff, (06), 0x6a, 0x68}},
	{APXOR_6, ymm_asm6, Py_asm6, [23]uint8{0xef, Pe_asm6, 0xef}},
	{AQUAD_6, ybyte_asm6, Px_asm6, [23]uint8{8}},
	{ARCLB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (02), 0xc0, (02), 0xd2, (02)}},
	{ARCLL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (02), 0xc1, (02), 0xd3, (02), 0xd3, (02)}},
	{ARCLQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (02), 0xc1, (02), 0xd3, (02), 0xd3, (02)}},
	{ARCLW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (02), 0xc1, (02), 0xd3, (02), 0xd3, (02)}},
	{ARCPPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x53}},
	{ARCPSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x53}},
	{ARCRB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (03), 0xc0, (03), 0xd2, (03)}},
	{ARCRL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (03), 0xc1, (03), 0xd3, (03), 0xd3, (03)}},
	{ARCRQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (03), 0xc1, (03), 0xd3, (03), 0xd3, (03)}},
	{ARCRW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (03), 0xc1, (03), 0xd3, (03), 0xd3, (03)}},
	{AREP_6, ynone_asm6, Px_asm6, [23]uint8{0xf3}},
	{AREPN_6, ynone_asm6, Px_asm6, [23]uint8{0xf2}},
	{ARET_6, ynone_asm6, Px_asm6, [23]uint8{0xc3}},
	{ARETFW_6, yret_asm6, Pe_asm6, [23]uint8{0xcb, 0xca}},
	{ARETFL_6, yret_asm6, Px_asm6, [23]uint8{0xcb, 0xca}},
	{ARETFQ_6, yret_asm6, Pw_asm6, [23]uint8{0xcb, 0xca}},
	{AROLB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (00), 0xc0, (00), 0xd2, (00)}},
	{AROLL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (00), 0xc1, (00), 0xd3, (00), 0xd3, (00)}},
	{AROLQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (00), 0xc1, (00), 0xd3, (00), 0xd3, (00)}},
	{AROLW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (00), 0xc1, (00), 0xd3, (00), 0xd3, (00)}},
	{ARORB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (01), 0xc0, (01), 0xd2, (01)}},
	{ARORL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (01), 0xc1, (01), 0xd3, (01), 0xd3, (01)}},
	{ARORQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (01), 0xc1, (01), 0xd3, (01), 0xd3, (01)}},
	{ARORW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (01), 0xc1, (01), 0xd3, (01), 0xd3, (01)}},
	{ARSQRTPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x52}},
	{ARSQRTSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x52}},
	{ASAHF_6, ynone_asm6, Px_asm6, [23]uint8{0x86, 0xe0, 0x50, 0x9d}}, /* XCHGB AH,AL; PUSH AX; POPFL */
	{ASALB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (04), 0xc0, (04), 0xd2, (04)}},
	{ASALL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (04), 0xc1, (04), 0xd3, (04), 0xd3, (04)}},
	{ASALQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (04), 0xc1, (04), 0xd3, (04), 0xd3, (04)}},
	{ASALW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (04), 0xc1, (04), 0xd3, (04), 0xd3, (04)}},
	{ASARB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (07), 0xc0, (07), 0xd2, (07)}},
	{ASARL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (07), 0xc1, (07), 0xd3, (07), 0xd3, (07)}},
	{ASARQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (07), 0xc1, (07), 0xd3, (07), 0xd3, (07)}},
	{ASARW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (07), 0xc1, (07), 0xd3, (07), 0xd3, (07)}},
	{ASBBB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x1c, 0x80, (03), 0x18, 0x1a}},
	{ASBBL_6, yxorl_asm6, Px_asm6, [23]uint8{0x83, (03), 0x1d, 0x81, (03), 0x19, 0x1b}},
	{ASBBQ_6, yxorl_asm6, Pw_asm6, [23]uint8{0x83, (03), 0x1d, 0x81, (03), 0x19, 0x1b}},
	{ASBBW_6, yxorl_asm6, Pe_asm6, [23]uint8{0x83, (03), 0x1d, 0x81, (03), 0x19, 0x1b}},
	{ASCASB_6, ynone_asm6, Pb_asm6, [23]uint8{0xae}},
	{ASCASL_6, ynone_asm6, Px_asm6, [23]uint8{0xaf}},
	{ASCASQ_6, ynone_asm6, Pw_asm6, [23]uint8{0xaf}},
	{ASCASW_6, ynone_asm6, Pe_asm6, [23]uint8{0xaf}},
	{ASETCC_6, yscond_asm6, Pm_asm6, [23]uint8{0x93, (00)}},
	{ASETCS_6, yscond_asm6, Pm_asm6, [23]uint8{0x92, (00)}},
	{ASETEQ_6, yscond_asm6, Pm_asm6, [23]uint8{0x94, (00)}},
	{ASETGE_6, yscond_asm6, Pm_asm6, [23]uint8{0x9d, (00)}},
	{ASETGT_6, yscond_asm6, Pm_asm6, [23]uint8{0x9f, (00)}},
	{ASETHI_6, yscond_asm6, Pm_asm6, [23]uint8{0x97, (00)}},
	{ASETLE_6, yscond_asm6, Pm_asm6, [23]uint8{0x9e, (00)}},
	{ASETLS_6, yscond_asm6, Pm_asm6, [23]uint8{0x96, (00)}},
	{ASETLT_6, yscond_asm6, Pm_asm6, [23]uint8{0x9c, (00)}},
	{ASETMI_6, yscond_asm6, Pm_asm6, [23]uint8{0x98, (00)}},
	{ASETNE_6, yscond_asm6, Pm_asm6, [23]uint8{0x95, (00)}},
	{ASETOC_6, yscond_asm6, Pm_asm6, [23]uint8{0x91, (00)}},
	{ASETOS_6, yscond_asm6, Pm_asm6, [23]uint8{0x90, (00)}},
	{ASETPC_6, yscond_asm6, Pm_asm6, [23]uint8{0x96, (00)}},
	{ASETPL_6, yscond_asm6, Pm_asm6, [23]uint8{0x99, (00)}},
	{ASETPS_6, yscond_asm6, Pm_asm6, [23]uint8{0x9a, (00)}},
	{ASHLB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (04), 0xc0, (04), 0xd2, (04)}},
	{ASHLL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (04), 0xc1, (04), 0xd3, (04), 0xd3, (04)}},
	{ASHLQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (04), 0xc1, (04), 0xd3, (04), 0xd3, (04)}},
	{ASHLW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (04), 0xc1, (04), 0xd3, (04), 0xd3, (04)}},
	{ASHRB_6, yshb_asm6, Pb_asm6, [23]uint8{0xd0, (05), 0xc0, (05), 0xd2, (05)}},
	{ASHRL_6, yshl_asm6, Px_asm6, [23]uint8{0xd1, (05), 0xc1, (05), 0xd3, (05), 0xd3, (05)}},
	{ASHRQ_6, yshl_asm6, Pw_asm6, [23]uint8{0xd1, (05), 0xc1, (05), 0xd3, (05), 0xd3, (05)}},
	{ASHRW_6, yshl_asm6, Pe_asm6, [23]uint8{0xd1, (05), 0xc1, (05), 0xd3, (05), 0xd3, (05)}},
	{ASHUFPD_6, yxshuf_asm6, Pq_asm6, [23]uint8{0xc6, (00)}},
	{ASHUFPS_6, yxshuf_asm6, Pm_asm6, [23]uint8{0xc6, (00)}},
	{ASQRTPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x51}},
	{ASQRTPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x51}},
	{ASQRTSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x51}},
	{ASQRTSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x51}},
	{ASTC_6, ynone_asm6, Px_asm6, [23]uint8{0xf9}},
	{ASTD_6, ynone_asm6, Px_asm6, [23]uint8{0xfd}},
	{ASTI_6, ynone_asm6, Px_asm6, [23]uint8{0xfb}},
	{ASTMXCSR_6, ysvrs_asm6, Pm_asm6, [23]uint8{0xae, (03), 0xae, (03)}},
	{ASTOSB_6, ynone_asm6, Pb_asm6, [23]uint8{0xaa}},
	{ASTOSL_6, ynone_asm6, Px_asm6, [23]uint8{0xab}},
	{ASTOSQ_6, ynone_asm6, Pw_asm6, [23]uint8{0xab}},
	{ASTOSW_6, ynone_asm6, Pe_asm6, [23]uint8{0xab}},
	{ASUBB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x2c, 0x80, (05), 0x28, 0x2a}},
	{ASUBL_6, yaddl_asm6, Px_asm6, [23]uint8{0x83, (05), 0x2d, 0x81, (05), 0x29, 0x2b}},
	{ASUBPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x5c}},
	{ASUBPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x5c}},
	{ASUBQ_6, yaddl_asm6, Pw_asm6, [23]uint8{0x83, (05), 0x2d, 0x81, (05), 0x29, 0x2b}},
	{ASUBSD_6, yxm_asm6, Pf2_asm6, [23]uint8{0x5c}},
	{ASUBSS_6, yxm_asm6, Pf3_asm6, [23]uint8{0x5c}},
	{ASUBW_6, yaddl_asm6, Pe_asm6, [23]uint8{0x83, (05), 0x2d, 0x81, (05), 0x29, 0x2b}},
	{ASWAPGS_6, ynone_asm6, Pm_asm6, [23]uint8{0x01, 0xf8}},
	{ASYSCALL_6, ynone_asm6, Px_asm6, [23]uint8{0x0f, 0x05}}, /* fast syscall */
	{ATESTB_6, ytestb_asm6, Pb_asm6, [23]uint8{0xa8, 0xf6, (00), 0x84, 0x84}},
	{ATESTL_6, ytestl_asm6, Px_asm6, [23]uint8{0xa9, 0xf7, (00), 0x85, 0x85}},
	{ATESTQ_6, ytestl_asm6, Pw_asm6, [23]uint8{0xa9, 0xf7, (00), 0x85, 0x85}},
	{ATESTW_6, ytestl_asm6, Pe_asm6, [23]uint8{0xa9, 0xf7, (00), 0x85, 0x85}},
	{ATEXT_6, ytext_asm6, Px_asm6, [23]uint8{}},
	{AUCOMISD_6, yxcmp_asm6, Pe_asm6, [23]uint8{0x2e}},
	{AUCOMISS_6, yxcmp_asm6, Pm_asm6, [23]uint8{0x2e}},
	{AUNPCKHPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x15}},
	{AUNPCKHPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x15}},
	{AUNPCKLPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x14}},
	{AUNPCKLPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x14}},
	{AVERR_6, ydivl_asm6, Pm_asm6, [23]uint8{0x00, (04)}},
	{AVERW_6, ydivl_asm6, Pm_asm6, [23]uint8{0x00, (05)}},
	{AWAIT_6, ynone_asm6, Px_asm6, [23]uint8{0x9b}},
	{AWORD_6, ybyte_asm6, Px_asm6, [23]uint8{2}},
	{AXCHGB_6, yml_mb_asm6, Pb_asm6, [23]uint8{0x86, 0x86}},
	{AXCHGL_6, yxchg_asm6, Px_asm6, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGQ_6, yxchg_asm6, Pw_asm6, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGW_6, yxchg_asm6, Pe_asm6, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXLAT_6, ynone_asm6, Px_asm6, [23]uint8{0xd7}},
	{AXORB_6, yxorb_asm6, Pb_asm6, [23]uint8{0x34, 0x80, (06), 0x30, 0x32}},
	{AXORL_6, yxorl_asm6, Px_asm6, [23]uint8{0x83, (06), 0x35, 0x81, (06), 0x31, 0x33}},
	{AXORPD_6, yxm_asm6, Pe_asm6, [23]uint8{0x57}},
	{AXORPS_6, yxm_asm6, Pm_asm6, [23]uint8{0x57}},
	{AXORQ_6, yxorl_asm6, Pw_asm6, [23]uint8{0x83, (06), 0x35, 0x81, (06), 0x31, 0x33}},
	{AXORW_6, yxorl_asm6, Pe_asm6, [23]uint8{0x83, (06), 0x35, 0x81, (06), 0x31, 0x33}},
	{AFMOVB_6, yfmvx_asm6, Px_asm6, [23]uint8{0xdf, (04)}},
	{AFMOVBP_6, yfmvp_asm6, Px_asm6, [23]uint8{0xdf, (06)}},
	{AFMOVD_6, yfmvd_asm6, Px_asm6, [23]uint8{0xdd, (00), 0xdd, (02), 0xd9, (00), 0xdd, (02)}},
	{AFMOVDP_6, yfmvdp_asm6, Px_asm6, [23]uint8{0xdd, (03), 0xdd, (03)}},
	{AFMOVF_6, yfmvf_asm6, Px_asm6, [23]uint8{0xd9, (00), 0xd9, (02)}},
	{AFMOVFP_6, yfmvp_asm6, Px_asm6, [23]uint8{0xd9, (03)}},
	{AFMOVL_6, yfmvf_asm6, Px_asm6, [23]uint8{0xdb, (00), 0xdb, (02)}},
	{AFMOVLP_6, yfmvp_asm6, Px_asm6, [23]uint8{0xdb, (03)}},
	{AFMOVV_6, yfmvx_asm6, Px_asm6, [23]uint8{0xdf, (05)}},
	{AFMOVVP_6, yfmvp_asm6, Px_asm6, [23]uint8{0xdf, (07)}},
	{AFMOVW_6, yfmvf_asm6, Px_asm6, [23]uint8{0xdf, (00), 0xdf, (02)}},
	{AFMOVWP_6, yfmvp_asm6, Px_asm6, [23]uint8{0xdf, (03)}},
	{AFMOVX_6, yfmvx_asm6, Px_asm6, [23]uint8{0xdb, (05)}},
	{AFMOVXP_6, yfmvp_asm6, Px_asm6, [23]uint8{0xdb, (07)}},
	{AFCOMB_6, nil, 0, [23]uint8{}},
	{AFCOMBP_6, nil, 0, [23]uint8{}},
	{AFCOMD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (02), 0xd8, (02), 0xdc, (02)}},  /* botch */
	{AFCOMDP_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (03), 0xd8, (03), 0xdc, (03)}}, /* botch */
	{AFCOMDPP_6, ycompp_asm6, Px_asm6, [23]uint8{0xde, (03)}},
	{AFCOMF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (02)}},
	{AFCOMFP_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (03)}},
	{AFCOML_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (02)}},
	{AFCOMLP_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (03)}},
	{AFCOMW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (02)}},
	{AFCOMWP_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (03)}},
	{AFUCOM_6, ycompp_asm6, Px_asm6, [23]uint8{0xdd, (04)}},
	{AFUCOMP_6, ycompp_asm6, Px_asm6, [23]uint8{0xdd, (05)}},
	{AFUCOMPP_6, ycompp_asm6, Px_asm6, [23]uint8{0xda, (13)}},
	{AFADDDP_6, yfaddp_asm6, Px_asm6, [23]uint8{0xde, (00)}},
	{AFADDW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (00)}},
	{AFADDL_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (00)}},
	{AFADDF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (00)}},
	{AFADDD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (00), 0xd8, (00), 0xdc, (00)}},
	{AFMULDP_6, yfaddp_asm6, Px_asm6, [23]uint8{0xde, (01)}},
	{AFMULW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (01)}},
	{AFMULL_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (01)}},
	{AFMULF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (01)}},
	{AFMULD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (01), 0xd8, (01), 0xdc, (01)}},
	{AFSUBDP_6, yfaddp_asm6, Px_asm6, [23]uint8{0xde, (05)}},
	{AFSUBW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (04)}},
	{AFSUBL_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (04)}},
	{AFSUBF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (04)}},
	{AFSUBD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (04), 0xd8, (04), 0xdc, (05)}},
	{AFSUBRDP_6, yfaddp_asm6, Px_asm6, [23]uint8{0xde, (04)}},
	{AFSUBRW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (05)}},
	{AFSUBRL_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (05)}},
	{AFSUBRF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (05)}},
	{AFSUBRD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (05), 0xd8, (05), 0xdc, (04)}},
	{AFDIVDP_6, yfaddp_asm6, Px_asm6, [23]uint8{0xde, (07)}},
	{AFDIVW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (06)}},
	{AFDIVL_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (06)}},
	{AFDIVF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (06)}},
	{AFDIVD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (06), 0xd8, (06), 0xdc, (07)}},
	{AFDIVRDP_6, yfaddp_asm6, Px_asm6, [23]uint8{0xde, (06)}},
	{AFDIVRW_6, yfmvx_asm6, Px_asm6, [23]uint8{0xde, (07)}},
	{AFDIVRL_6, yfmvx_asm6, Px_asm6, [23]uint8{0xda, (07)}},
	{AFDIVRF_6, yfmvx_asm6, Px_asm6, [23]uint8{0xd8, (07)}},
	{AFDIVRD_6, yfadd_asm6, Px_asm6, [23]uint8{0xdc, (07), 0xd8, (07), 0xdc, (06)}},
	{AFXCHD_6, yfxch_asm6, Px_asm6, [23]uint8{0xd9, (01), 0xd9, (01)}},
	{AFFREE_6, nil, 0, [23]uint8{}},
	{AFLDCW_6, ystcw_asm6, Px_asm6, [23]uint8{0xd9, (05), 0xd9, (05)}},
	{AFLDENV_6, ystcw_asm6, Px_asm6, [23]uint8{0xd9, (04), 0xd9, (04)}},
	{AFRSTOR_6, ysvrs_asm6, Px_asm6, [23]uint8{0xdd, (04), 0xdd, (04)}},
	{AFSAVE_6, ysvrs_asm6, Px_asm6, [23]uint8{0xdd, (06), 0xdd, (06)}},
	{AFSTCW_6, ystcw_asm6, Px_asm6, [23]uint8{0xd9, (07), 0xd9, (07)}},
	{AFSTENV_6, ystcw_asm6, Px_asm6, [23]uint8{0xd9, (06), 0xd9, (06)}},
	{AFSTSW_6, ystsw_asm6, Px_asm6, [23]uint8{0xdd, (07), 0xdf, 0xe0}},
	{AF2XM1_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf0}},
	{AFABS_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xe1}},
	{AFCHS_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xe0}},
	{AFCLEX_6, ynone_asm6, Px_asm6, [23]uint8{0xdb, 0xe2}},
	{AFCOS_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xff}},
	{AFDECSTP_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf6}},
	{AFINCSTP_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf7}},
	{AFINIT_6, ynone_asm6, Px_asm6, [23]uint8{0xdb, 0xe3}},
	{AFLD1_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xe8}},
	{AFLDL2E_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xea}},
	{AFLDL2T_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xe9}},
	{AFLDLG2_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xec}},
	{AFLDLN2_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xed}},
	{AFLDPI_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xeb}},
	{AFLDZ_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xee}},
	{AFNOP_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xd0}},
	{AFPATAN_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf3}},
	{AFPREM_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf8}},
	{AFPREM1_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf5}},
	{AFPTAN_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf2}},
	{AFRNDINT_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xfc}},
	{AFSCALE_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xfd}},
	{AFSIN_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xfe}},
	{AFSINCOS_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xfb}},
	{AFSQRT_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xfa}},
	{AFTST_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xe4}},
	{AFXAM_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xe5}},
	{AFXTRACT_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf4}},
	{AFYL2X_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf1}},
	{AFYL2XP1_6, ynone_asm6, Px_asm6, [23]uint8{0xd9, 0xf9}},
	{ACMPXCHGB_6, yrb_mb_asm6, Pb_asm6, [23]uint8{0x0f, 0xb0}},
	{ACMPXCHGL_6, yrl_ml_asm6, Px_asm6, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGW_6, yrl_ml_asm6, Pe_asm6, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGQ_6, yrl_ml_asm6, Pw_asm6, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHG8B_6, yscond_asm6, Pm_asm6, [23]uint8{0xc7, (01)}},
	{AINVD_6, ynone_asm6, Pm_asm6, [23]uint8{0x08}},
	{AINVLPG_6, ymbs_asm6, Pm_asm6, [23]uint8{0x01, (07)}},
	{ALFENCE_6, ynone_asm6, Pm_asm6, [23]uint8{0xae, 0xe8}},
	{AMFENCE_6, ynone_asm6, Pm_asm6, [23]uint8{0xae, 0xf0}},
	{AMOVNTIL_6, yrl_ml_asm6, Pm_asm6, [23]uint8{0xc3}},
	{AMOVNTIQ_6, yrl_ml_asm6, Pw_asm6, [23]uint8{0x0f, 0xc3}},
	{ARDMSR_6, ynone_asm6, Pm_asm6, [23]uint8{0x32}},
	{ARDPMC_6, ynone_asm6, Pm_asm6, [23]uint8{0x33}},
	{ARDTSC_6, ynone_asm6, Pm_asm6, [23]uint8{0x31}},
	{ARSM_6, ynone_asm6, Pm_asm6, [23]uint8{0xaa}},
	{ASFENCE_6, ynone_asm6, Pm_asm6, [23]uint8{0xae, 0xf8}},
	{ASYSRET_6, ynone_asm6, Pm_asm6, [23]uint8{0x07}},
	{AWBINVD_6, ynone_asm6, Pm_asm6, [23]uint8{0x09}},
	{AWRMSR_6, ynone_asm6, Pm_asm6, [23]uint8{0x30}},
	{AXADDB_6, yrb_mb_asm6, Pb_asm6, [23]uint8{0x0f, 0xc0}},
	{AXADDL_6, yrl_ml_asm6, Px_asm6, [23]uint8{0x0f, 0xc1}},
	{AXADDQ_6, yrl_ml_asm6, Pw_asm6, [23]uint8{0x0f, 0xc1}},
	{AXADDW_6, yrl_ml_asm6, Pe_asm6, [23]uint8{0x0f, 0xc1}},
	{ACRC32B_6, ycrc32l_asm6, Px_asm6, [23]uint8{0xf2, 0x0f, 0x38, 0xf0, 0}},
	{ACRC32Q_6, ycrc32l_asm6, Pw_asm6, [23]uint8{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{APREFETCHT0_6, yprefetch_asm6, Pm_asm6, [23]uint8{0x18, (01)}},
	{APREFETCHT1_6, yprefetch_asm6, Pm_asm6, [23]uint8{0x18, (02)}},
	{APREFETCHT2_6, yprefetch_asm6, Pm_asm6, [23]uint8{0x18, (03)}},
	{APREFETCHNTA_6, yprefetch_asm6, Pm_asm6, [23]uint8{0x18, (00)}},
	{AMOVQL_6, yrl_ml_asm6, Px_asm6, [23]uint8{0x89}},
	{AUNDEF_6, ynone_asm6, Px_asm6, [23]uint8{0x0f, 0x0b}},
	{AAESENC_6, yaes_asm6, Pq_asm6, [23]uint8{0x38, 0xdc, (0)}},
	{AAESENCLAST_6, yaes_asm6, Pq_asm6, [23]uint8{0x38, 0xdd, (0)}},
	{AAESDEC_6, yaes_asm6, Pq_asm6, [23]uint8{0x38, 0xde, (0)}},
	{AAESDECLAST_6, yaes_asm6, Pq_asm6, [23]uint8{0x38, 0xdf, (0)}},
	{AAESIMC_6, yaes_asm6, Pq_asm6, [23]uint8{0x38, 0xdb, (0)}},
	{AAESKEYGENASSIST_6, yaes2_asm6, Pq_asm6, [23]uint8{0x3a, 0xdf, (0)}},
	{APSHUFD_6, yaes2_asm6, Pq_asm6, [23]uint8{0x70, (0)}},
	{APCLMULQDQ_6, yxshuf_asm6, Pq_asm6, [23]uint8{0x3a, 0x44, 0}},
	{AUSEFIELD_6, ynop_asm6, Px_asm6, [23]uint8{0, 0}},
	{ATYPE_6, nil, 0, [23]uint8{}},
	{AFUNCDATA_6, yfuncdata_asm6, Px_asm6, [23]uint8{0, 0}},
	{APCDATA_6, ypcdata_asm6, Px_asm6, [23]uint8{0, 0}},
	{ACHECKNIL_6, nil, 0, [23]uint8{}},
	{AVARDEF_6, nil, 0, [23]uint8{}},
	{AVARKILL_6, nil, 0, [23]uint8{}},
	{ADUFFCOPY_6, yduff_asm6, Px_asm6, [23]uint8{0xe8}},
	{ADUFFZERO_6, yduff_asm6, Px_asm6, [23]uint8{0xe8}},
	{AEND_6, nil, 0, [23]uint8{}},
	{0, nil, 0, [23]uint8{}},
}

var opindex_asm6 [ALAST_6 + 1]*Optab_asm6

/*
static void
relput8(Prog *p, Addr *a)
{
	vlong v;
	Reloc rel, *r;

	v = vaddr(ctxt, a, &rel);
	if(rel.siz != 0) {
		r = addrel(ctxt->cursym);
		*r = rel;
		r->siz = 8;
		r->off = p->pc + ctxt->andptr - ctxt->and;
	}
	put8(ctxt, v);
}
*/
func vaddr_asm6(ctxt *Link, a *Addr, r *Reloc) int64 {
	var t int
	var v int64
	var s *LSym
	if r != nil {
		*r = Reloc{}
	}
	t = a.typ
	v = a.offset
	if t == int(D_ADDR_6) {
		t = a.index
	}
	switch t {
	case D_STATIC_6, D_EXTERN_6:
		s = a.sym
		if r == nil {
			ctxt.diag("need reloc for %D", a)
			sysfatal("reloc")
		}
		r.siz = 4  // TODO: 8 for external symbols
		r.off = -1 // caller must fill in
		r.sym = s
		r.add = v
		v = 0
		if ctxt.flag_shared != 0 || ctxt.headtype == int(Hnacl) {
			if s.typ == int(STLSBSS) {
				r.xadd = r.add - int64(r.siz)
				r.typ = int(R_TLS)
				r.xsym = s
			} else {
				r.typ = int(R_PCREL)
			}
		} else {
			r.typ = int(R_ADDR)
		}
		break
	case D_INDIR_6 + D_TLS_6:
		if r == nil {
			ctxt.diag("need reloc for %D", a)
			sysfatal("reloc")
		}
		r.typ = int(R_TLS_LE)
		r.siz = 4
		r.off = -1 // caller must fill in
		r.add = v
		v = 0
		break
	}
	return v
}

// single-instruction no-ops of various lengths.
// constructed by hand and disassembled with gdb to verify.
// see http://www.agner.org/optimize/optimizing_assembly.pdf for discussion.
var nop_asm6 = [][16]uint8{
	{0x90},
	{0x66, 0x90},
	{0x0F, 0x1F, 0x00},
	{0x0F, 0x1F, 0x40, 0x00},
	{0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x0F, 0x1F, 0x80, 0x00, 0x00, 0x00, 0x00},
	{0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
}

// Native Client rejects the repeated 0x66 prefix.
// {0x66, 0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
func fillnop_asm6(p []uint8, n int) {
	var m int
	for n > 0 {
		m = n
		if m > len(nop_asm6) {
			m = len(nop_asm6)
		}
		copy(p, nop_asm6[m-1][:m])
		p = p[m:]
		n -= m
	}
}

func instinit_asm6() {
	var c int
	var i int
	for i = 1; optab_asm6[i].as != 0; i++ {
		c = optab_asm6[i].as
		if opindex_asm6[c] != nil {
			sysfatal("phase error in optab: %d (%A)", i, c)
		}
		opindex_asm6[c] = &optab_asm6[i]
	}
	for i = 0; i < int(Ymax_asm6); i++ {
		ycover_asm6[i*int(Ymax_asm6)+i] = 1
	}
	ycover_asm6[Yi0_asm6*Ymax_asm6+Yi8_asm6] = 1
	ycover_asm6[Yi1_asm6*Ymax_asm6+Yi8_asm6] = 1
	ycover_asm6[Yi0_asm6*Ymax_asm6+Ys32_asm6] = 1
	ycover_asm6[Yi1_asm6*Ymax_asm6+Ys32_asm6] = 1
	ycover_asm6[Yi8_asm6*Ymax_asm6+Ys32_asm6] = 1
	ycover_asm6[Yi0_asm6*Ymax_asm6+Yi32_asm6] = 1
	ycover_asm6[Yi1_asm6*Ymax_asm6+Yi32_asm6] = 1
	ycover_asm6[Yi8_asm6*Ymax_asm6+Yi32_asm6] = 1
	ycover_asm6[Ys32_asm6*Ymax_asm6+Yi32_asm6] = 1
	ycover_asm6[Yi0_asm6*Ymax_asm6+Yi64_asm6] = 1
	ycover_asm6[Yi1_asm6*Ymax_asm6+Yi64_asm6] = 1
	ycover_asm6[Yi8_asm6*Ymax_asm6+Yi64_asm6] = 1
	ycover_asm6[Ys32_asm6*Ymax_asm6+Yi64_asm6] = 1
	ycover_asm6[Yi32_asm6*Ymax_asm6+Yi64_asm6] = 1
	ycover_asm6[Yal_asm6*Ymax_asm6+Yrb_asm6] = 1
	ycover_asm6[Ycl_asm6*Ymax_asm6+Yrb_asm6] = 1
	ycover_asm6[Yax_asm6*Ymax_asm6+Yrb_asm6] = 1
	ycover_asm6[Ycx_asm6*Ymax_asm6+Yrb_asm6] = 1
	ycover_asm6[Yrx_asm6*Ymax_asm6+Yrb_asm6] = 1
	ycover_asm6[Yrl_asm6*Ymax_asm6+Yrb_asm6] = 1
	ycover_asm6[Ycl_asm6*Ymax_asm6+Ycx_asm6] = 1
	ycover_asm6[Yax_asm6*Ymax_asm6+Yrx_asm6] = 1
	ycover_asm6[Ycx_asm6*Ymax_asm6+Yrx_asm6] = 1
	ycover_asm6[Yax_asm6*Ymax_asm6+Yrl_asm6] = 1
	ycover_asm6[Ycx_asm6*Ymax_asm6+Yrl_asm6] = 1
	ycover_asm6[Yrx_asm6*Ymax_asm6+Yrl_asm6] = 1
	ycover_asm6[Yf0_asm6*Ymax_asm6+Yrf_asm6] = 1
	ycover_asm6[Yal_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Ycl_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Yax_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Ycx_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Yrx_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Yrb_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Yrl_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Ym_asm6*Ymax_asm6+Ymb_asm6] = 1
	ycover_asm6[Yax_asm6*Ymax_asm6+Yml_asm6] = 1
	ycover_asm6[Ycx_asm6*Ymax_asm6+Yml_asm6] = 1
	ycover_asm6[Yrx_asm6*Ymax_asm6+Yml_asm6] = 1
	ycover_asm6[Yrl_asm6*Ymax_asm6+Yml_asm6] = 1
	ycover_asm6[Ym_asm6*Ymax_asm6+Yml_asm6] = 1
	ycover_asm6[Yax_asm6*Ymax_asm6+Ymm_asm6] = 1
	ycover_asm6[Ycx_asm6*Ymax_asm6+Ymm_asm6] = 1
	ycover_asm6[Yrx_asm6*Ymax_asm6+Ymm_asm6] = 1
	ycover_asm6[Yrl_asm6*Ymax_asm6+Ymm_asm6] = 1
	ycover_asm6[Ym_asm6*Ymax_asm6+Ymm_asm6] = 1
	ycover_asm6[Ymr_asm6*Ymax_asm6+Ymm_asm6] = 1
	ycover_asm6[Ym_asm6*Ymax_asm6+Yxm_asm6] = 1
	ycover_asm6[Yxr_asm6*Ymax_asm6+Yxm_asm6] = 1
	for i = 0; i < int(D_NONE_6); i++ {
		reg_asm6[i] = -1
		if i >= int(D_AL_6) && i <= int(D_R15B_6) {
			reg_asm6[i] = (i - int(D_AL_6)) & 7
			if i >= int(D_SPB_6) && i <= int(D_DIB_6) {
				regrex_asm6[i] = 0x40
			}
			if i >= int(D_R8B_6) && i <= int(D_R15B_6) {
				regrex_asm6[i] = int(Rxr_asm6 | Rxx_asm6 | Rxb_asm6)
			}
		}
		if i >= int(D_AH_6) && i <= int(D_BH_6) {
			reg_asm6[i] = 4 + ((i - int(D_AH_6)) & 7)
		}
		if i >= int(D_AX_6) && i <= int(D_R15_6) {
			reg_asm6[i] = (i - int(D_AX_6)) & 7
			if i >= int(D_R8_6) {
				regrex_asm6[i] = int(Rxr_asm6 | Rxx_asm6 | Rxb_asm6)
			}
		}
		if i >= int(D_F0_6) && i <= D_F0_6+7 {
			reg_asm6[i] = (i - int(D_F0_6)) & 7
		}
		if i >= int(D_M0_6) && i <= D_M0_6+7 {
			reg_asm6[i] = (i - int(D_M0_6)) & 7
		}
		if i >= int(D_X0_6) && i <= D_X0_6+15 {
			reg_asm6[i] = (i - int(D_X0_6)) & 7
			if i >= D_X0_6+8 {
				regrex_asm6[i] = int(Rxr_asm6 | Rxx_asm6 | Rxb_asm6)
			}
		}
		if i >= D_CR_6+8 && i <= D_CR_6+15 {
			regrex_asm6[i] = int(Rxr_asm6)
		}
	}
}

func naclpad_asm6(ctxt *Link, s *LSym, c int32, pad int) int32 {
	symgrow(ctxt, s, int64(c)+int64(pad))
	fillnop_asm6(s.p[c:], pad)
	return c + int32(pad)
}

func spadjop_asm6(ctxt *Link, p *Prog, l int, q int) int {
	if p.mode != 64 || ctxt.arch.ptrsize == 4 {
		return l
	}
	return q
}

func prefixof_asm6(ctxt *Link, a *Addr) int {
	switch a.typ {
	case D_INDIR_6 + D_CS_6:
		return 0x2e
	case D_INDIR_6 + D_DS_6:
		return 0x3e
	case D_INDIR_6 + D_ES_6:
		return 0x26
	case D_INDIR_6 + D_FS_6:
		return 0x64
	case D_INDIR_6 + D_GS_6:
		return 0x65
	// NOTE: Systems listed here should be only systems that
	// support direct TLS references like 8(TLS) implemented as
	// direct references from FS or GS. Systems that require
	// the initial-exec model, where you load the TLS base into
	// a register and then index from that register, do not reach
	// this code and should not be listed.
	case D_INDIR_6 + D_TLS_6:
		switch ctxt.headtype {
		default:
			sysfatal("unknown TLS base register for %s", headstr(ctxt.headtype))
		case Hdragonfly, Hfreebsd, Hlinux, Hnetbsd, Hopenbsd, Hsolaris:
			return 0x64 // FS
		case Hdarwin:
			return 0x65 // GS
		}
	}
	switch a.index {
	case D_CS_6:
		return 0x2e
	case D_DS_6:
		return 0x3e
	case D_ES_6:
		return 0x26
	case D_FS_6:
		return 0x64
	case D_GS_6:
		return 0x65
	}
	return 0
}

func oclass_asm6(ctxt *Link, a *Addr) int {
	var v int64
	var l int32
	if a.typ >= int(D_INDIR_6) || a.index != int(D_NONE_6) {
		if a.index != int(D_NONE_6) && a.scale == 0 {
			if a.typ == int(D_ADDR_6) {
				switch a.index {
				case D_EXTERN_6, D_STATIC_6:
					if ctxt.flag_shared != 0 || ctxt.headtype == int(Hnacl) {
						return int(Yiauto_asm6)
					} else {
						return int(Yi32_asm6) /* TO DO: Yi64 */
					}
				case D_AUTO_6, D_PARAM_6:
					return int(Yiauto_asm6)
				}
				return int(Yxxx_asm6)
			}
			return int(Ycol_asm6)
		}
		return int(Ym_asm6)
	}
	switch a.typ {
	case D_AL_6:
		return int(Yal_asm6)
	case D_AX_6:
		return int(Yax_asm6)
	/*
		case D_SPB:
	*/
	case D_BPB_6, D_SIB_6, D_DIB_6, D_R8B_6, D_R9B_6, D_R10B_6, D_R11B_6, D_R12B_6, D_R13B_6, D_R14B_6, D_R15B_6:
		if ctxt.asmode != 64 {
			return int(Yxxx_asm6)
		}
	case D_DL_6, D_BL_6, D_AH_6, D_CH_6, D_DH_6, D_BH_6:
		return int(Yrb_asm6)
	case D_CL_6:
		return int(Ycl_asm6)
	case D_CX_6:
		return int(Ycx_asm6)
	case D_DX_6, D_BX_6:
		return int(Yrx_asm6)
	case D_R8_6: /* not really Yrl */
	case D_R9_6, D_R10_6, D_R11_6, D_R12_6, D_R13_6, D_R14_6, D_R15_6:
		if ctxt.asmode != 64 {
			return int(Yxxx_asm6)
		}
	case D_SP_6, D_BP_6, D_SI_6, D_DI_6:
		return int(Yrl_asm6)
	case D_F0_6 + 0:
		return int(Yf0_asm6)
	case D_F0_6 + 1, D_F0_6 + 2, D_F0_6 + 3, D_F0_6 + 4, D_F0_6 + 5, D_F0_6 + 6, D_F0_6 + 7:
		return int(Yrf_asm6)
	case D_M0_6 + 0, D_M0_6 + 1, D_M0_6 + 2, D_M0_6 + 3, D_M0_6 + 4, D_M0_6 + 5, D_M0_6 + 6, D_M0_6 + 7:
		return int(Ymr_asm6)
	case D_X0_6 + 0, D_X0_6 + 1, D_X0_6 + 2, D_X0_6 + 3, D_X0_6 + 4, D_X0_6 + 5, D_X0_6 + 6, D_X0_6 + 7, D_X0_6 + 8, D_X0_6 + 9, D_X0_6 + 10, D_X0_6 + 11, D_X0_6 + 12, D_X0_6 + 13, D_X0_6 + 14, D_X0_6 + 15:
		return int(Yxr_asm6)
	case D_NONE_6:
		return int(Ynone_asm6)
	case D_CS_6:
		return int(Ycs_asm6)
	case D_SS_6:
		return int(Yss_asm6)
	case D_DS_6:
		return int(Yds_asm6)
	case D_ES_6:
		return int(Yes_asm6)
	case D_FS_6:
		return int(Yfs_asm6)
	case D_GS_6:
		return int(Ygs_asm6)
	case D_TLS_6:
		return int(Ytls_asm6)
	case D_GDTR_6:
		return int(Ygdtr_asm6)
	case D_IDTR_6:
		return int(Yidtr_asm6)
	case D_LDTR_6:
		return int(Yldtr_asm6)
	case D_MSW_6:
		return int(Ymsw_asm6)
	case D_TASK_6:
		return int(Ytask_asm6)
	case D_CR_6 + 0:
		return int(Ycr0_asm6)
	case D_CR_6 + 1:
		return int(Ycr1_asm6)
	case D_CR_6 + 2:
		return int(Ycr2_asm6)
	case D_CR_6 + 3:
		return int(Ycr3_asm6)
	case D_CR_6 + 4:
		return int(Ycr4_asm6)
	case D_CR_6 + 5:
		return int(Ycr5_asm6)
	case D_CR_6 + 6:
		return int(Ycr6_asm6)
	case D_CR_6 + 7:
		return int(Ycr7_asm6)
	case D_CR_6 + 8:
		return int(Ycr8_asm6)
	case D_DR_6 + 0:
		return int(Ydr0_asm6)
	case D_DR_6 + 1:
		return int(Ydr1_asm6)
	case D_DR_6 + 2:
		return int(Ydr2_asm6)
	case D_DR_6 + 3:
		return int(Ydr3_asm6)
	case D_DR_6 + 4:
		return int(Ydr4_asm6)
	case D_DR_6 + 5:
		return int(Ydr5_asm6)
	case D_DR_6 + 6:
		return int(Ydr6_asm6)
	case D_DR_6 + 7:
		return int(Ydr7_asm6)
	case D_TR_6 + 0:
		return int(Ytr0_asm6)
	case D_TR_6 + 1:
		return int(Ytr1_asm6)
	case D_TR_6 + 2:
		return int(Ytr2_asm6)
	case D_TR_6 + 3:
		return int(Ytr3_asm6)
	case D_TR_6 + 4:
		return int(Ytr4_asm6)
	case D_TR_6 + 5:
		return int(Ytr5_asm6)
	case D_TR_6 + 6:
		return int(Ytr6_asm6)
	case D_TR_6 + 7:
		return int(Ytr7_asm6)
	case D_EXTERN_6, D_STATIC_6, D_AUTO_6, D_PARAM_6:
		return int(Ym_asm6)
	case D_CONST_6, D_ADDR_6:
		if a.sym == nil {
			v = a.offset
			if v == 0 {
				return int(Yi0_asm6)
			}
			if v == 1 {
				return int(Yi1_asm6)
			}
			if v >= -128 && v <= 127 {
				return int(Yi8_asm6)
			}
			l = int32(v)
			if int64(l) == v {
				return int(Ys32_asm6) /* can sign extend */
			}
			if (v >> 32) == 0 {
				return int(Yi32_asm6) /* unsigned */
			}
			return int(Yi64_asm6)
		}
		return int(Yi32_asm6) /* TO DO: D_ADDR as Yi64 */
	case D_BRANCH_6:
		return int(Ybr_asm6)
	}
	return int(Yxxx_asm6)
}

func asmidx_asm6(ctxt *Link, scale int, index int, base int) {
	var i int
	switch index {
	default:
		goto bad
	case D_NONE_6:
		i = 4 << 3
		goto bas
	case D_R8_6, D_R9_6, D_R10_6, D_R11_6, D_R12_6, D_R13_6, D_R14_6, D_R15_6:
		if ctxt.asmode != 64 {
			goto bad
		}
	case D_AX_6, D_CX_6, D_DX_6, D_BX_6, D_BP_6, D_SI_6, D_DI_6:
		i = reg_asm6[index] << 3
		break
	}
	switch scale {
	default:
		goto bad
	case 1:
		break
	case 2:
		i |= (1 << 6)
		break
	case 4:
		i |= (2 << 6)
		break
	case 8:
		i |= (3 << 6)
		break
	}
bas:
	switch base {
	default:
		goto bad
	case D_NONE_6: /* must be mod=00 */
		i |= 5
		break
	case D_R8_6, D_R9_6, D_R10_6, D_R11_6, D_R12_6, D_R13_6, D_R14_6, D_R15_6:
		if ctxt.asmode != 64 {
			goto bad
		}
	case D_AX_6, D_CX_6, D_DX_6, D_BX_6, D_SP_6, D_BP_6, D_SI_6, D_DI_6:
		i |= reg_asm6[base]
		break
	}
	ctxt.andptr[0] = uint8(i)
	ctxt.andptr = ctxt.andptr[1:]
	return
bad:
	ctxt.diag("asmidx: bad address %d/%d/%d", scale, index, base)
	ctxt.andptr[0] = 0
	ctxt.andptr = ctxt.andptr[1:]
	return
}

func put4_asm6(ctxt *Link, v int32) {
	ctxt.andptr[0] = uint8(v)
	ctxt.andptr[1] = uint8(v >> 8)
	ctxt.andptr[2] = uint8(v >> 16)
	ctxt.andptr[3] = uint8(v >> 24)
	ctxt.andptr = ctxt.andptr[4:]
}

func relput4_asm6(ctxt *Link, p *Prog, a *Addr) {
	var v int64
	var rel Reloc
	var r *Reloc
	v = vaddr_asm6(ctxt, a, &rel)
	if rel.siz != 0 {
		if rel.siz != 4 {
			ctxt.diag("bad reloc")
		}
		r = addrel(ctxt.cursym)
		*r = rel
		r.off = p.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
	}
	put4_asm6(ctxt, int32(v))
}

func put8_asm6(ctxt *Link, v int64) {
	ctxt.andptr[0] = uint8(v)
	ctxt.andptr[1] = uint8(v >> 8)
	ctxt.andptr[2] = uint8(v >> 16)
	ctxt.andptr[3] = uint8(v >> 24)
	ctxt.andptr[4] = uint8(v >> 32)
	ctxt.andptr[5] = uint8(v >> 40)
	ctxt.andptr[6] = uint8(v >> 48)
	ctxt.andptr[7] = uint8(v >> 56)
	ctxt.andptr = ctxt.andptr[8:]
}

func asmandsz_asm6(ctxt *Link, a *Addr, r int, rex int, m64 int) {
	var v int32
	var t int
	var scale int
	var rel Reloc
	rex &= (0x40 | int(Rxr_asm6))
	v = int32(a.offset)
	t = a.typ
	rel.siz = 0
	if a.index != int(D_NONE_6) && a.index != int(D_TLS_6) {
		if t < int(D_INDIR_6) {
			switch t {
			default:
				goto bad
			case D_STATIC_6, D_EXTERN_6:
				if ctxt.flag_shared != 0 || ctxt.headtype == int(Hnacl) {
					goto bad
				}
				t = int(D_NONE_6)
				v = int32(vaddr_asm6(ctxt, a, &rel))
				break
			case D_AUTO_6, D_PARAM_6:
				t = int(D_SP_6)
				break
			}
		} else {
			t -= int(D_INDIR_6)
		}
		ctxt.rexflag |= (regrex_asm6[int(a.index)] & int(Rxx_asm6)) | (regrex_asm6[t] & int(Rxb_asm6)) | rex
		if t == int(D_NONE_6) {
			ctxt.andptr[0] = uint8((0 << 6) | (4 << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm6(ctxt, a.scale, a.index, t)
			goto putrelv
		}
		if v == 0 && rel.siz == 0 && t != int(D_BP_6) && t != int(D_R13_6) {
			ctxt.andptr[0] = uint8((0 << 6) | (4 << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm6(ctxt, a.scale, a.index, t)
			return
		}
		if v >= -128 && v < 128 && rel.siz == 0 {
			ctxt.andptr[0] = uint8((1 << 6) | (4 << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm6(ctxt, a.scale, a.index, t)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
		ctxt.andptr[0] = uint8((2 << 6) | (4 << 0) | (r << 3))
		ctxt.andptr = ctxt.andptr[1:]
		asmidx_asm6(ctxt, a.scale, a.index, t)
		goto putrelv
	}
	if t >= int(D_AL_6) && t <= D_X0_6+15 {
		if v != 0 {
			goto bad
		}
		ctxt.andptr[0] = uint8((3 << 6) | (reg_asm6[t] << 0) | (r << 3))
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.rexflag |= (regrex_asm6[t] & (0x40 | int(Rxb_asm6))) | rex
		return
	}
	scale = a.scale
	if t < int(D_INDIR_6) {
		switch a.typ {
		default:
			goto bad
		case D_STATIC_6, D_EXTERN_6:
			t = int(D_NONE_6)
			v = int32(vaddr_asm6(ctxt, a, &rel))
			break
		case D_AUTO_6, D_PARAM_6:
			t = int(D_SP_6)
			break
		}
		scale = 1
	} else {
		t -= int(D_INDIR_6)
	}
	if t == int(D_TLS_6) {
		v = int32(vaddr_asm6(ctxt, a, &rel))
	}
	ctxt.rexflag |= (regrex_asm6[t] & int(Rxb_asm6)) | rex
	if t == int(D_NONE_6) || (D_CS_6 <= int(t) && t <= int(D_GS_6)) || t == int(D_TLS_6) {
		if (ctxt.flag_shared != 0 || ctxt.headtype == int(Hnacl)) && t == int(D_NONE_6) && (a.typ == int(D_STATIC_6) || a.typ == int(D_EXTERN_6)) || ctxt.asmode != 64 {
			ctxt.andptr[0] = uint8((0 << 6) | (5 << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			goto putrelv
		}
		/* temporary */
		ctxt.andptr[0] = uint8((0 << 6) | (4 << 0) | (r << 3))
		ctxt.andptr = ctxt.andptr[1:] /* sib present */
		ctxt.andptr[0] = uint8((0 << 6) | (4 << 3) | (5 << 0))
		ctxt.andptr = ctxt.andptr[1:] /* DS:d32 */
		goto putrelv
	}
	if t == int(D_SP_6) || t == int(D_R12_6) {
		if v == 0 {
			ctxt.andptr[0] = uint8((0 << 6) | (reg_asm6[t] << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm6(ctxt, scale, int(D_NONE_6), t)
			return
		}
		if v >= -128 && v < 128 {
			ctxt.andptr[0] = uint8((1 << 6) | (reg_asm6[t] << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm6(ctxt, scale, int(D_NONE_6), t)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
		ctxt.andptr[0] = uint8((2 << 6) | (reg_asm6[t] << 0) | (r << 3))
		ctxt.andptr = ctxt.andptr[1:]
		asmidx_asm6(ctxt, scale, int(D_NONE_6), t)
		goto putrelv
	}
	if t >= int(D_AX_6) && t <= int(D_R15_6) {
		if a.index == int(D_TLS_6) {
			rel = Reloc{}
			rel.typ = int(R_TLS_IE)
			rel.siz = 4
			rel.sym = (*LSym)(nil)
			rel.add = int64(v)
			v = 0
		}
		if v == 0 && rel.siz == 0 && t != int(D_BP_6) && t != int(D_R13_6) {
			ctxt.andptr[0] = uint8((0 << 6) | (reg_asm6[t] << 0) | (r << 3))
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
		if v >= -128 && v < 128 && rel.siz == 0 {
			ctxt.andptr[0] = uint8((1 << 6) | (reg_asm6[t] << 0) | (r << 3))
			ctxt.andptr[1] = uint8(v)
			ctxt.andptr = ctxt.andptr[2:]
			return
		}
		ctxt.andptr[0] = uint8((2 << 6) | (reg_asm6[t] << 0) | (r << 3))
		ctxt.andptr = ctxt.andptr[1:]
		goto putrelv
	}
	goto bad
putrelv:
	if rel.siz != 0 {
		var r *Reloc
		if rel.siz != 4 {
			ctxt.diag("bad rel")
			goto bad
		}
		r = addrel(ctxt.cursym)
		*r = rel
		r.off = ctxt.curp.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
	}
	put4_asm6(ctxt, v)
	return
bad:
	ctxt.diag("asmand: bad address %D", a)
	return
}

func asmand_asm6(ctxt *Link, a *Addr, ra *Addr) {
	asmandsz_asm6(ctxt, a, reg_asm6[ra.typ], regrex_asm6[ra.typ], 0)
}

func asmando_asm6(ctxt *Link, a *Addr, o int) {
	asmandsz_asm6(ctxt, a, o, 0, 0)
}

func bytereg_asm6(a *Addr, t *int) {
	if a.index == int(D_NONE_6) && (a.typ >= int(D_AX_6) && a.typ <= int(D_R15_6)) {
		a.typ = int(D_AL_6 + int(a.typ-int(D_AX_6)))
		*t = 0
	}
}

const (
	E_asm6 = 0xff
)

var ymovtab_asm6 = []Movtab_asm6{
	/* push */
	{APUSHL_6, Ycs_asm6, Ynone_asm6, 0, [4]uint8{0x0e, E_asm6, 0, 0}},
	{APUSHL_6, Yss_asm6, Ynone_asm6, 0, [4]uint8{0x16, E_asm6, 0, 0}},
	{APUSHL_6, Yds_asm6, Ynone_asm6, 0, [4]uint8{0x1e, E_asm6, 0, 0}},
	{APUSHL_6, Yes_asm6, Ynone_asm6, 0, [4]uint8{0x06, E_asm6, 0, 0}},
	{APUSHL_6, Yfs_asm6, Ynone_asm6, 0, [4]uint8{0x0f, 0xa0, E_asm6, 0}},
	{APUSHL_6, Ygs_asm6, Ynone_asm6, 0, [4]uint8{0x0f, 0xa8, E_asm6, 0}},
	{APUSHQ_6, Yfs_asm6, Ynone_asm6, 0, [4]uint8{0x0f, 0xa0, E_asm6, 0}},
	{APUSHQ_6, Ygs_asm6, Ynone_asm6, 0, [4]uint8{0x0f, 0xa8, E_asm6, 0}},
	{APUSHW_6, Ycs_asm6, Ynone_asm6, 0, [4]uint8{Pe_asm6, 0x0e, E_asm6, 0}},
	{APUSHW_6, Yss_asm6, Ynone_asm6, 0, [4]uint8{Pe_asm6, 0x16, E_asm6, 0}},
	{APUSHW_6, Yds_asm6, Ynone_asm6, 0, [4]uint8{Pe_asm6, 0x1e, E_asm6, 0}},
	{APUSHW_6, Yes_asm6, Ynone_asm6, 0, [4]uint8{Pe_asm6, 0x06, E_asm6, 0}},
	{APUSHW_6, Yfs_asm6, Ynone_asm6, 0, [4]uint8{Pe_asm6, 0x0f, 0xa0, E_asm6}},
	{APUSHW_6, Ygs_asm6, Ynone_asm6, 0, [4]uint8{Pe_asm6, 0x0f, 0xa8, E_asm6}},
	/* pop */
	{APOPL_6, Ynone_asm6, Yds_asm6, 0, [4]uint8{0x1f, E_asm6, 0, 0}},
	{APOPL_6, Ynone_asm6, Yes_asm6, 0, [4]uint8{0x07, E_asm6, 0, 0}},
	{APOPL_6, Ynone_asm6, Yss_asm6, 0, [4]uint8{0x17, E_asm6, 0, 0}},
	{APOPL_6, Ynone_asm6, Yfs_asm6, 0, [4]uint8{0x0f, 0xa1, E_asm6, 0}},
	{APOPL_6, Ynone_asm6, Ygs_asm6, 0, [4]uint8{0x0f, 0xa9, E_asm6, 0}},
	{APOPQ_6, Ynone_asm6, Yfs_asm6, 0, [4]uint8{0x0f, 0xa1, E_asm6, 0}},
	{APOPQ_6, Ynone_asm6, Ygs_asm6, 0, [4]uint8{0x0f, 0xa9, E_asm6, 0}},
	{APOPW_6, Ynone_asm6, Yds_asm6, 0, [4]uint8{Pe_asm6, 0x1f, E_asm6, 0}},
	{APOPW_6, Ynone_asm6, Yes_asm6, 0, [4]uint8{Pe_asm6, 0x07, E_asm6, 0}},
	{APOPW_6, Ynone_asm6, Yss_asm6, 0, [4]uint8{Pe_asm6, 0x17, E_asm6, 0}},
	{APOPW_6, Ynone_asm6, Yfs_asm6, 0, [4]uint8{Pe_asm6, 0x0f, 0xa1, E_asm6}},
	{APOPW_6, Ynone_asm6, Ygs_asm6, 0, [4]uint8{Pe_asm6, 0x0f, 0xa9, E_asm6}},
	/* mov seg */
	{AMOVW_6, Yes_asm6, Yml_asm6, 1, [4]uint8{0x8c, 0, 0, 0}},
	{AMOVW_6, Ycs_asm6, Yml_asm6, 1, [4]uint8{0x8c, 1, 0, 0}},
	{AMOVW_6, Yss_asm6, Yml_asm6, 1, [4]uint8{0x8c, 2, 0, 0}},
	{AMOVW_6, Yds_asm6, Yml_asm6, 1, [4]uint8{0x8c, 3, 0, 0}},
	{AMOVW_6, Yfs_asm6, Yml_asm6, 1, [4]uint8{0x8c, 4, 0, 0}},
	{AMOVW_6, Ygs_asm6, Yml_asm6, 1, [4]uint8{0x8c, 5, 0, 0}},
	{AMOVW_6, Yml_asm6, Yes_asm6, 2, [4]uint8{0x8e, 0, 0, 0}},
	{AMOVW_6, Yml_asm6, Ycs_asm6, 2, [4]uint8{0x8e, 1, 0, 0}},
	{AMOVW_6, Yml_asm6, Yss_asm6, 2, [4]uint8{0x8e, 2, 0, 0}},
	{AMOVW_6, Yml_asm6, Yds_asm6, 2, [4]uint8{0x8e, 3, 0, 0}},
	{AMOVW_6, Yml_asm6, Yfs_asm6, 2, [4]uint8{0x8e, 4, 0, 0}},
	{AMOVW_6, Yml_asm6, Ygs_asm6, 2, [4]uint8{0x8e, 5, 0, 0}},
	/* mov cr */
	{AMOVL_6, Ycr0_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVL_6, Ycr2_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVL_6, Ycr3_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVL_6, Ycr4_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVL_6, Ycr8_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVQ_6, Ycr0_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVQ_6, Ycr2_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVQ_6, Ycr3_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVQ_6, Ycr4_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVQ_6, Ycr8_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVL_6, Yml_asm6, Ycr0_asm6, 4, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVL_6, Yml_asm6, Ycr2_asm6, 4, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVL_6, Yml_asm6, Ycr3_asm6, 4, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVL_6, Yml_asm6, Ycr4_asm6, 4, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVL_6, Yml_asm6, Ycr8_asm6, 4, [4]uint8{0x0f, 0x22, 8, 0}},
	{AMOVQ_6, Yml_asm6, Ycr0_asm6, 4, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVQ_6, Yml_asm6, Ycr2_asm6, 4, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVQ_6, Yml_asm6, Ycr3_asm6, 4, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVQ_6, Yml_asm6, Ycr4_asm6, 4, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVQ_6, Yml_asm6, Ycr8_asm6, 4, [4]uint8{0x0f, 0x22, 8, 0}},
	/* mov dr */
	{AMOVL_6, Ydr0_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVL_6, Ydr6_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVL_6, Ydr7_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVQ_6, Ydr0_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVQ_6, Ydr6_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVQ_6, Ydr7_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVL_6, Yml_asm6, Ydr0_asm6, 4, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVL_6, Yml_asm6, Ydr6_asm6, 4, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVL_6, Yml_asm6, Ydr7_asm6, 4, [4]uint8{0x0f, 0x23, 7, 0}},
	{AMOVQ_6, Yml_asm6, Ydr0_asm6, 4, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVQ_6, Yml_asm6, Ydr6_asm6, 4, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVQ_6, Yml_asm6, Ydr7_asm6, 4, [4]uint8{0x0f, 0x23, 7, 0}},
	/* mov tr */
	{AMOVL_6, Ytr6_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x24, 6, 0}},
	{AMOVL_6, Ytr7_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x24, 7, 0}},
	{AMOVL_6, Yml_asm6, Ytr6_asm6, 4, [4]uint8{0x0f, 0x26, 6, E_asm6}},
	{AMOVL_6, Yml_asm6, Ytr7_asm6, 4, [4]uint8{0x0f, 0x26, 7, E_asm6}},
	/* lgdt, sgdt, lidt, sidt */
	{AMOVL_6, Ym_asm6, Ygdtr_asm6, 4, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVL_6, Ygdtr_asm6, Ym_asm6, 3, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVL_6, Ym_asm6, Yidtr_asm6, 4, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVL_6, Yidtr_asm6, Ym_asm6, 3, [4]uint8{0x0f, 0x01, 1, 0}},
	{AMOVQ_6, Ym_asm6, Ygdtr_asm6, 4, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVQ_6, Ygdtr_asm6, Ym_asm6, 3, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVQ_6, Ym_asm6, Yidtr_asm6, 4, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVQ_6, Yidtr_asm6, Ym_asm6, 3, [4]uint8{0x0f, 0x01, 1, 0}},
	/* lldt, sldt */
	{AMOVW_6, Yml_asm6, Yldtr_asm6, 4, [4]uint8{0x0f, 0x00, 2, 0}},
	{AMOVW_6, Yldtr_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x00, 0, 0}},
	/* lmsw, smsw */
	{AMOVW_6, Yml_asm6, Ymsw_asm6, 4, [4]uint8{0x0f, 0x01, 6, 0}},
	{AMOVW_6, Ymsw_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x01, 4, 0}},
	/* ltr, str */
	{AMOVW_6, Yml_asm6, Ytask_asm6, 4, [4]uint8{0x0f, 0x00, 3, 0}},
	{AMOVW_6, Ytask_asm6, Yml_asm6, 3, [4]uint8{0x0f, 0x00, 1, 0}},
	/* load full pointer */
	{AMOVL_6, Yml_asm6, Ycol_asm6, 5, [4]uint8{0, 0, 0, 0}},
	{AMOVW_6, Yml_asm6, Ycol_asm6, 5, [4]uint8{Pe_asm6, 0, 0, 0}},
	/* double shift */
	{ASHLL_6, Ycol_asm6, Yml_asm6, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHRL_6, Ycol_asm6, Yml_asm6, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHLQ_6, Ycol_asm6, Yml_asm6, 6, [4]uint8{Pw_asm6, 0xa4, 0xa5, 0}},
	{ASHRQ_6, Ycol_asm6, Yml_asm6, 6, [4]uint8{Pw_asm6, 0xac, 0xad, 0}},
	{ASHLW_6, Ycol_asm6, Yml_asm6, 6, [4]uint8{Pe_asm6, 0xa4, 0xa5, 0}},
	{ASHRW_6, Ycol_asm6, Yml_asm6, 6, [4]uint8{Pe_asm6, 0xac, 0xad, 0}},
	/* load TLS base */
	{AMOVQ_6, Ytls_asm6, Yrl_asm6, 7, [4]uint8{0, 0, 0, 0}},
	{0, 0, 0, 0, [4]uint8{}},
}

func isax_asm6(a *Addr) int {
	switch a.typ {
	case D_AX_6, D_AL_6, D_AH_6, D_INDIR_6 + D_AX_6:
		return 1
	}
	if a.index == int(D_AX_6) {
		return 1
	}
	return 0
}

func subreg_asm6(p *Prog, from int, to int) {
	if false { /*debug['Q']*/
		print("\n%P	s/%R/%R/\n", p, from, to)
	}
	if p.from.typ == from {
		p.from.typ = to
	}
	if p.to.typ == from {
		p.to.typ = to
	}
	if p.from.index == from {
		p.from.index = to
	}
	if p.to.index == from {
		p.to.index = to
	}
	from += int(D_INDIR_6)
	if p.from.typ == from {
		p.from.typ = to + int(D_INDIR_6)
	}
	if p.to.typ == from {
		p.to.typ = to + int(D_INDIR_6)
	}
	if false { /*debug['Q']*/
		print("%P\n", p)
	}
}

func mediaop_asm6(ctxt *Link, o *Optab_asm6, op int, osize int, z int) int {
	switch op {
	case Pm_asm6, Pe_asm6, Pf2_asm6, Pf3_asm6:
		if osize != 1 {
			if op != int(Pm_asm6) {
				ctxt.andptr[0] = uint8(op)
				ctxt.andptr = ctxt.andptr[1:]
			}
			ctxt.andptr[0] = uint8(Pm_asm6)
			ctxt.andptr = ctxt.andptr[1:]
			z++
			op = int(o.op[z])
			break
		}
	default:
		if cap(ctxt.andptr) == cap(ctxt.and) || int(ctxt.and[len(ctxt.andptr)-1]) != Pm_asm6 {
			ctxt.andptr[0] = uint8(Pm_asm6)
			ctxt.andptr = ctxt.andptr[1:]
		}
		break
	}
	ctxt.andptr[0] = uint8(op)
	ctxt.andptr = ctxt.andptr[1:]
	return z
}

func doasm_asm6(ctxt *Link, p *Prog) {
	var o *Optab_asm6
	var q *Prog
	var pp Prog
	var t []uint8
	var mo []Movtab_asm6
	var z int
	var op int
	var ft int
	var tt int
	var xo bool
	var l int
	var pre int
	var v int64
	var rel Reloc
	var r *Reloc
	var a *Addr
	ctxt.curp = p // TODO
	o = opindex_asm6[p.as]
	if o == nil {
		ctxt.diag("asmins: missing op %P", p)
		return
	}
	pre = prefixof_asm6(ctxt, &p.from)
	if pre != 0 {
		ctxt.andptr[0] = uint8(pre)
		ctxt.andptr = ctxt.andptr[1:]
	}
	pre = prefixof_asm6(ctxt, &p.to)
	if pre != 0 {
		ctxt.andptr[0] = uint8(pre)
		ctxt.andptr = ctxt.andptr[1:]
	}
	if p.ft == 0 {
		p.ft = oclass_asm6(ctxt, &p.from)
	}
	if p.tt == 0 {
		p.tt = oclass_asm6(ctxt, &p.to)
	}
	ft = p.ft * int(Ymax_asm6)
	tt = p.tt * int(Ymax_asm6)
	t = o.ytab
	if t == nil {
		ctxt.diag("asmins: noproto %P", p)
		return
	}
	xo = o.op[0] == 0x0f
	for z = 0; t[0] != 0; (func() { z += int(t[3]) + int(bool2int(xo)); t = t[4:] })() {
		if ycover_asm6[ft+int(t[0])] != 0 {
			if ycover_asm6[tt+int(t[1])] != 0 {
				goto found
			}
		}
	}
	goto domov
found:
	switch o.prefix {
	case Pq_asm6: /* 16 bit escape and opcode escape */
		ctxt.andptr[0] = uint8(Pe_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(Pm_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Pq3_asm6: /* 16 bit escape, Rex.w, and opcode escape */
		ctxt.andptr[0] = uint8(Pe_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(Pw_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(Pm_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Pf2_asm6: /* xmm opcode escape */
	case Pf3_asm6:
		ctxt.andptr[0] = uint8(o.prefix)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(Pm_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Pm_asm6: /* opcode escape */
		ctxt.andptr[0] = uint8(Pm_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Pe_asm6: /* 16 bit escape */
		ctxt.andptr[0] = uint8(Pe_asm6)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Pw_asm6: /* 64-bit escape */
		if p.mode != 64 {
			ctxt.diag("asmins: illegal 64: %P", p)
		}
		ctxt.rexflag |= int(Pw_asm6)
		break
	case Pb_asm6: /* botch */
		bytereg_asm6(&p.from, &p.ft)
		bytereg_asm6(&p.to, &p.tt)
		break
	case P32_asm6: /* 32 bit but illegal if 64-bit mode */
		if p.mode == 64 {
			ctxt.diag("asmins: illegal in 64-bit mode: %P", p)
		}
		break
	case Py_asm6: /* 64-bit only, no prefix */
		if p.mode != 64 {
			ctxt.diag("asmins: illegal in %d-bit mode: %P", p.mode, p)
		}
		break
	}
	if z >= len(o.op) {
		sysfatal("asmins bad table %P", p)
	}
	op = int(o.op[z])
	if op == 0x0f {
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		z++
		op = int(o.op[z])
	}
	switch t[2] {
	default:
		ctxt.diag("asmins: unknown z %d %P", t[2], p)
		return
	case Zpseudo_asm6:
		break
	case Zlit_asm6:
		for ; ; z++ {
			op = int(o.op[z])
			if !(op != 0) {
				break
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		}
		break
	case Zlitm_r_asm6:
		for ; ; z++ {
			op = int(o.op[z])
			if !(op != 0) {
				break
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		}
		asmand_asm6(ctxt, &p.from, &p.to)
		break
	case Zmb_r_asm6:
		bytereg_asm6(&p.from, &p.ft)
	/* fall through */
	case Zm_r_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.from, &p.to)
		break
	case Zm2_r_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = o.op[z+1]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.from, &p.to)
		break
	case Zm_r_xm_asm6:
		mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmand_asm6(ctxt, &p.from, &p.to)
		break
	case Zm_r_xm_nr_asm6:
		ctxt.rexflag = 0
		mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmand_asm6(ctxt, &p.from, &p.to)
		break
	case Zm_r_i_xm_asm6:
		mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmand_asm6(ctxt, &p.from, &p.to)
		ctxt.andptr[0] = uint8(p.to.offset)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zm_r_3d_asm6:
		ctxt.andptr[0] = 0x0f
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = 0x0f
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.from, &p.to)
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zibm_r_asm6:
		for {
			var tmp int = z
			z++
			op = int(o.op[tmp])
			if !((op) != 0) {
				break
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		}
		asmand_asm6(ctxt, &p.from, &p.to)
		ctxt.andptr[0] = uint8(p.to.offset)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zaut_r_asm6:
		ctxt.andptr[0] = 0x8d
		ctxt.andptr = ctxt.andptr[1:] /* leal */
		if p.from.typ != int(D_ADDR_6) {
			ctxt.diag("asmins: Zaut sb type ADDR")
		}
		p.from.typ = p.from.index
		p.from.index = int(D_NONE_6)
		asmand_asm6(ctxt, &p.from, &p.to)
		p.from.index = p.from.typ
		p.from.typ = int(D_ADDR_6)
		break
	case Zm_o_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.from, int(o.op[z+1]))
		break
	case Zr_m_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.to, &p.from)
		break
	case Zr_m_xm_asm6:
		mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmand_asm6(ctxt, &p.to, &p.from)
		break
	case Zr_m_xm_nr_asm6:
		ctxt.rexflag = 0
		mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmand_asm6(ctxt, &p.to, &p.from)
		break
	case Zr_m_i_xm_asm6:
		mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmand_asm6(ctxt, &p.to, &p.from)
		ctxt.andptr[0] = uint8(p.from.offset)
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zo_m_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.to, int(o.op[z+1]))
		break
	case Zcallindreg_asm6:
		r = addrel(ctxt.cursym)
		r.off = p.pc
		r.typ = int(R_CALLIND)
		r.siz = 0
		fallthrough
	case Zo_m64_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmandsz_asm6(ctxt, &p.to, int(o.op[z+1]), 0, 1)
		break
	case Zm_ibo_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.from, int(o.op[z+1]))
		ctxt.andptr[0] = uint8(vaddr_asm6(ctxt, &p.to, (*Reloc)(nil)))
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zibo_m_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.to, int(o.op[z+1]))
		ctxt.andptr[0] = uint8(vaddr_asm6(ctxt, &p.from, (*Reloc)(nil)))
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zibo_m_xm_asm6:
		z = mediaop_asm6(ctxt, o, op, int(t[3]), z)
		asmando_asm6(ctxt, &p.to, int(o.op[z+1]))
		ctxt.andptr[0] = uint8(vaddr_asm6(ctxt, &p.from, (*Reloc)(nil)))
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Z_ib_asm6, Zib__asm6:
		if int(t[2]) == Zib__asm6 {
			a = &p.from
		} else {
			a = &p.to
		}
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(vaddr_asm6(ctxt, a, (*Reloc)(nil)))
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zib_rp_asm6:
		ctxt.rexflag |= regrex_asm6[p.to.typ] & int(Rxb_asm6|0x40)
		ctxt.andptr[0] = uint8(op + reg_asm6[p.to.typ])
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(vaddr_asm6(ctxt, &p.from, (*Reloc)(nil)))
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zil_rp_asm6:
		ctxt.rexflag |= regrex_asm6[p.to.typ] & int(Rxb_asm6)
		ctxt.andptr[0] = uint8(op + reg_asm6[p.to.typ])
		ctxt.andptr = ctxt.andptr[1:]
		if o.prefix == int(Pe_asm6) {
			v = vaddr_asm6(ctxt, &p.from, (*Reloc)(nil))
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm6(ctxt, p, &p.from)
		}
		break
	case Zo_iw_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		if p.from.typ != int(D_NONE_6) {
			v = vaddr_asm6(ctxt, &p.from, (*Reloc)(nil))
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		}
		break
	case Ziq_rp_asm6:
		v = vaddr_asm6(ctxt, &p.from, &rel)
		l = int(v >> 32)
		if l == 0 && rel.siz != 8 {
			//p->mark |= 0100;
			//print("zero: %llux %P\n", v, p);
			ctxt.rexflag &^= (0x40 | int(Rxw_asm6))
			ctxt.rexflag |= regrex_asm6[p.to.typ] & int(Rxb_asm6)
			ctxt.andptr[0] = uint8(0xb8 + reg_asm6[p.to.typ])
			ctxt.andptr = ctxt.andptr[1:]
			if rel.typ != 0 {
				r = addrel(ctxt.cursym)
				*r = rel
				r.off = p.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
			}
			put4_asm6(ctxt, int32(v))
		} else {
			if l == -1 && (uint64(v)&(uint64(1)<<31)) != 0 { /* sign extend */
				//p->mark |= 0100;
				//print("sign: %llux %P\n", v, p);
				ctxt.andptr[0] = 0xc7
				ctxt.andptr = ctxt.andptr[1:]
				asmando_asm6(ctxt, &p.to, 0)
				put4_asm6(ctxt, int32(v)) /* need all 8 */
			} else {
				//print("all: %llux %P\n", v, p);
				ctxt.rexflag |= regrex_asm6[p.to.typ] & int(Rxb_asm6)
				ctxt.andptr[0] = uint8(op + reg_asm6[p.to.typ])
				ctxt.andptr = ctxt.andptr[1:]
				if rel.typ != 0 {
					r = addrel(ctxt.cursym)
					*r = rel
					r.off = p.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
				}
				put8_asm6(ctxt, v)
			}
		}
		break
	case Zib_rr_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.to, &p.to)
		ctxt.andptr[0] = uint8(vaddr_asm6(ctxt, &p.from, (*Reloc)(nil)))
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Z_il_asm6, Zil__asm6:
		if int(t[2]) == Zil__asm6 {
			a = &p.from
		} else {
			a = &p.to
		}
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		if o.prefix == int(Pe_asm6) {
			v = vaddr_asm6(ctxt, a, (*Reloc)(nil))
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm6(ctxt, p, a)
		}
		break
	case Zm_ilo_asm6, Zilo_m_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		if int(t[2]) == Zilo_m_asm6 {
			a = &p.from
			asmando_asm6(ctxt, &p.to, int(o.op[z+1]))
		} else {
			a = &p.to
			asmando_asm6(ctxt, &p.from, int(o.op[z+1]))
		}
		if o.prefix == int(Pe_asm6) {
			v = vaddr_asm6(ctxt, a, (*Reloc)(nil))
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm6(ctxt, p, a)
		}
		break
	case Zil_rr_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.to, &p.to)
		if o.prefix == int(Pe_asm6) {
			v = vaddr_asm6(ctxt, &p.from, (*Reloc)(nil))
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm6(ctxt, p, &p.from)
		}
		break
	case Z_rp_asm6:
		ctxt.rexflag |= regrex_asm6[p.to.typ] & int(Rxb_asm6|0x40)
		ctxt.andptr[0] = uint8(op + reg_asm6[p.to.typ])
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zrp__asm6:
		ctxt.rexflag |= regrex_asm6[p.from.typ] & int(Rxb_asm6|0x40)
		ctxt.andptr[0] = uint8(op + reg_asm6[p.from.typ])
		ctxt.andptr = ctxt.andptr[1:]
		break
	case Zclr_asm6:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm6(ctxt, &p.to, &p.to)
		break
	case Zcall_asm6:
		if p.to.sym == nil {
			ctxt.diag("call without target")
			sysfatal("bad code")
		}
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		r = addrel(ctxt.cursym)
		r.off = p.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
		r.sym = p.to.sym
		r.add = p.to.offset
		r.typ = int(R_CALL)
		r.siz = 4
		put4_asm6(ctxt, 0)
		break
	// TODO: jump across functions needs reloc
	case Zbr_asm6, Zjmp_asm6, Zloop_asm6:
		if p.to.sym != nil {
			if int(t[2]) != Zjmp_asm6 {
				ctxt.diag("branch to ATEXT")
				sysfatal("bad code")
			}
			ctxt.andptr[0] = o.op[z+1]
			ctxt.andptr = ctxt.andptr[1:]
			r = addrel(ctxt.cursym)
			r.off = p.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
			r.sym = p.to.sym
			r.typ = int(R_PCREL)
			r.siz = 4
			put4_asm6(ctxt, 0)
			break
		}
		// Assumes q is in this function.
		// TODO: Check in input, preserve in brchain.
		// Fill in backward jump now.
		q = p.pcond
		if q == nil {
			ctxt.diag("jmp/branch/loop without target")
			sysfatal("bad code")
		}
		if p.back&1 != 0 /*untyped*/ {
			v = int64(q.pc) - (int64(p.pc) + 2)
			if v >= -128 {
				if p.as == int(AJCXZL_6) {
					ctxt.andptr[0] = 0x67
					ctxt.andptr = ctxt.andptr[1:]
				}
				ctxt.andptr[0] = uint8(op)
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(v)
				ctxt.andptr = ctxt.andptr[1:]
			} else {
				if int(t[2]) == Zloop_asm6 {
					ctxt.diag("loop too far: %P", p)
				} else {
					v -= 5 - 2
					if int(t[2]) == Zbr_asm6 {
						ctxt.andptr[0] = 0x0f
						ctxt.andptr = ctxt.andptr[1:]
						v--
					}
					ctxt.andptr[0] = o.op[z+1]
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v)
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v >> 8)
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v >> 16)
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v >> 24)
					ctxt.andptr = ctxt.andptr[1:]
				}
			}
			break
		}
		// Annotate target; will fill in later.
		p.forwd = q.comefrom
		q.comefrom = p
		if p.back&2 != 0 /*untyped*/ { // short
			if p.as == int(AJCXZL_6) {
				ctxt.andptr[0] = 0x67
				ctxt.andptr = ctxt.andptr[1:]
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			if int(t[2]) == Zloop_asm6 {
				ctxt.diag("loop too far: %P", p)
			} else {
				if int(t[2]) == Zbr_asm6 {
					ctxt.andptr[0] = 0x0f
					ctxt.andptr = ctxt.andptr[1:]
				}
				ctxt.andptr[0] = o.op[z+1]
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0
				ctxt.andptr = ctxt.andptr[1:]
			}
		}
		break
		/*
			v = q->pc - p->pc - 2;
			if((v >= -128 && v <= 127) || p->pc == -1 || q->pc == -1) {
				*ctxt->andptr++ = op;
				*ctxt->andptr++ = v;
			} else {
				v -= 5-2;
				if(t[2] == Zbr) {
					*ctxt->andptr++ = 0x0f;
					v--;
				}
				*ctxt->andptr++ = o->op[z+1];
				*ctxt->andptr++ = v;
				*ctxt->andptr++ = v>>8;
				*ctxt->andptr++ = v>>16;
				*ctxt->andptr++ = v>>24;
			}
		*/
		break
	case Zbyte_asm6:
		v = vaddr_asm6(ctxt, &p.from, &rel)
		if rel.siz != 0 {
			rel.siz = uint8(op)
			r = addrel(ctxt.cursym)
			*r = rel
			r.off = p.pc - (-cap(ctxt.andptr) + cap(ctxt.and))
		}
		ctxt.andptr[0] = uint8(v)
		ctxt.andptr = ctxt.andptr[1:]
		if op > 1 {
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
			if op > 2 {
				ctxt.andptr[0] = uint8(v >> 16)
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(v >> 24)
				ctxt.andptr = ctxt.andptr[1:]
				if op > 4 {
					ctxt.andptr[0] = uint8(v >> 32)
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v >> 40)
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v >> 48)
					ctxt.andptr = ctxt.andptr[1:]
					ctxt.andptr[0] = uint8(v >> 56)
					ctxt.andptr = ctxt.andptr[1:]
				}
			}
		}
		break
	}
	return
domov:
	for mo = ymovtab_asm6; mo[0].as != 0; mo = mo[1:] {
		if p.as == mo[0].as {
			if ycover_asm6[ft+int(mo[0].ft)] != 0 {
				if ycover_asm6[tt+int(mo[0].tt)] != 0 {
					t = mo[0].op[:]
					goto mfound
				}
			}
		}
	}
bad:
	if p.mode != 64 {
		/*
		 * here, the assembly has failed.
		 * if its a byte instruction that has
		 * unaddressable registers, try to
		 * exchange registers and reissue the
		 * instruction with the operands renamed.
		 */
		pp = *p
		z = p.from.typ
		if z >= int(D_BP_6) && z <= int(D_DI_6) {
			if isax_asm6(&p.to) != 0 || p.to.typ == int(D_NONE_6) {
				// We certainly don't want to exchange
				// with AX if the op is MUL or DIV.
				ctxt.andptr[0] = 0x87
				ctxt.andptr = ctxt.andptr[1:] /* xchg lhs,bx */
				asmando_asm6(ctxt, &p.from, reg_asm6[D_BX_6])
				subreg_asm6(&pp, z, int(D_BX_6))
				doasm_asm6(ctxt, &pp)
				ctxt.andptr[0] = 0x87
				ctxt.andptr = ctxt.andptr[1:] /* xchg lhs,bx */
				asmando_asm6(ctxt, &p.from, reg_asm6[D_BX_6])
			} else {
				ctxt.andptr[0] = uint8(0x90 + reg_asm6[z])
				ctxt.andptr = ctxt.andptr[1:] /* xchg lsh,ax */
				subreg_asm6(&pp, z, int(D_AX_6))
				doasm_asm6(ctxt, &pp)
				ctxt.andptr[0] = uint8(0x90 + reg_asm6[z])
				ctxt.andptr = ctxt.andptr[1:] /* xchg lsh,ax */
			}
			return
		}
		z = p.to.typ
		if z >= int(D_BP_6) && z <= int(D_DI_6) {
			if isax_asm6(&p.from) != 0 {
				ctxt.andptr[0] = 0x87
				ctxt.andptr = ctxt.andptr[1:] /* xchg rhs,bx */
				asmando_asm6(ctxt, &p.to, reg_asm6[D_BX_6])
				subreg_asm6(&pp, z, int(D_BX_6))
				doasm_asm6(ctxt, &pp)
				ctxt.andptr[0] = 0x87
				ctxt.andptr = ctxt.andptr[1:] /* xchg rhs,bx */
				asmando_asm6(ctxt, &p.to, reg_asm6[D_BX_6])
			} else {
				ctxt.andptr[0] = uint8(0x90 + reg_asm6[z])
				ctxt.andptr = ctxt.andptr[1:] /* xchg rsh,ax */
				subreg_asm6(&pp, z, int(D_AX_6))
				doasm_asm6(ctxt, &pp)
				ctxt.andptr[0] = uint8(0x90 + reg_asm6[z])
				ctxt.andptr = ctxt.andptr[1:] /* xchg rsh,ax */
			}
			return
		}
	}
	ctxt.diag("doasm: notfound from=%#x to=%#x ft=%d tt=%d %v", p.from.typ, p.to.typ, ft, tt, Pconv_list6(ctxt, p))
	return
mfound:
	switch mo[0].code {
	default:
		ctxt.diag("asmins: unknown mov %d %P", mo[0].code, p)
		break
	case 0: /* lit */
		for z = 0; int(t[z]) != E_asm6; z++ {
			ctxt.andptr[0] = t[z]
			ctxt.andptr = ctxt.andptr[1:]
		}
		break
	case 1: /* r,m */
		ctxt.andptr[0] = t[0]
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.to, int(t[1]))
		break
	case 2: /* m,r */
		ctxt.andptr[0] = t[0]
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.from, int(t[1]))
		break
	case 3: /* r,m - 2op */
		ctxt.andptr[0] = t[0]
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = t[1]
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.to, int(t[2]))
		ctxt.rexflag |= regrex_asm6[p.from.typ] & int(Rxr_asm6|0x40)
		break
	case 4: /* m,r - 2op */
		ctxt.andptr[0] = t[0]
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = t[1]
		ctxt.andptr = ctxt.andptr[1:]
		asmando_asm6(ctxt, &p.from, int(t[2]))
		ctxt.rexflag |= regrex_asm6[p.to.typ] & int(Rxr_asm6|0x40)
		break
	case 5: /* load full pointer, trash heap */
		if t[0] != 0 {
			ctxt.andptr[0] = t[0]
			ctxt.andptr = ctxt.andptr[1:]
		}
		switch p.to.index {
		default:
			goto bad
		case D_DS_6:
			ctxt.andptr[0] = 0xc5
			ctxt.andptr = ctxt.andptr[1:]
			break
		case D_SS_6:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0xb2
			ctxt.andptr = ctxt.andptr[1:]
			break
		case D_ES_6:
			ctxt.andptr[0] = 0xc4
			ctxt.andptr = ctxt.andptr[1:]
			break
		case D_FS_6:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0xb4
			ctxt.andptr = ctxt.andptr[1:]
			break
		case D_GS_6:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0xb5
			ctxt.andptr = ctxt.andptr[1:]
			break
		}
		asmand_asm6(ctxt, &p.from, &p.to)
		break
	case 6: /* double shift */
		if int(t[0]) == Pw_asm6 {
			if p.mode != 64 {
				ctxt.diag("asmins: illegal 64: %P", p)
			}
			ctxt.rexflag |= int(Pw_asm6)
			t = t[1:]
		} else {
			if int(t[0]) == Pe_asm6 {
				ctxt.andptr[0] = uint8(Pe_asm6)
				ctxt.andptr = ctxt.andptr[1:]
				t = t[1:]
			}
		}
		z = p.from.typ
		switch z {
		default:
			goto bad
		case D_CONST_6:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = t[0]
			ctxt.andptr = ctxt.andptr[1:]
			asmandsz_asm6(ctxt, &p.to, reg_asm6[int(p.from.index)], regrex_asm6[int(p.from.index)], 0)
			ctxt.andptr[0] = uint8(p.from.offset)
			ctxt.andptr = ctxt.andptr[1:]
			break
		case D_CL_6, D_CX_6:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = t[1]
			ctxt.andptr = ctxt.andptr[1:]
			asmandsz_asm6(ctxt, &p.to, reg_asm6[int(p.from.index)], regrex_asm6[int(p.from.index)], 0)
			break
		}
		break
	// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
	// where you load the TLS base register into a register and then index off that
	// register to access the actual TLS variables. Systems that allow direct TLS access
	// are handled in prefixof above and should not be listed here.
	case 7: /* mov tls, r */
		switch ctxt.headtype {
		default:
			sysfatal("unknown TLS base location for %s", headstr(ctxt.headtype))
		case Hplan9:
			if ctxt.plan9privates == nil {
				ctxt.plan9privates = linklookup(ctxt, "_privates", 0)
			}
			pp.from = Addr{}
			pp.from.typ = int(D_EXTERN_6)
			pp.from.sym = ctxt.plan9privates
			pp.from.offset = 0
			pp.from.index = int(D_NONE_6)
			ctxt.rexflag |= int(Pw_asm6)
			ctxt.andptr[0] = 0x8B
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm6(ctxt, &pp.from, &p.to)
			break
		// TLS base is 0(FS).
		case Hsolaris: // TODO(rsc): Delete Hsolaris from list. Should not use this code. See progedit in obj6.c.
			pp.from = p.from
			pp.from.typ = int(D_INDIR_6 + D_NONE_6)
			pp.from.offset = 0
			pp.from.index = int(D_NONE_6)
			pp.from.scale = 0
			ctxt.rexflag |= int(Pw_asm6)
			ctxt.andptr[0] = 0x64
			ctxt.andptr = ctxt.andptr[1:] // FS
			ctxt.andptr[0] = 0x8B
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm6(ctxt, &pp.from, &p.to)
			break
		// Windows TLS base is always 0x28(GS).
		case Hwindows:
			pp.from = p.from
			pp.from.typ = int(D_INDIR_6 + D_GS_6)
			pp.from.offset = 0x28
			pp.from.index = int(D_NONE_6)
			pp.from.scale = 0
			ctxt.rexflag |= int(Pw_asm6)
			ctxt.andptr[0] = 0x65
			ctxt.andptr = ctxt.andptr[1:] // GS
			ctxt.andptr[0] = 0x8B
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm6(ctxt, &pp.from, &p.to)
			break
		}
		break
	}
}

var naclret_asm6 = []uint8{
	0x5e, // POPL SI
	// 0x8b, 0x7d, 0x00, // MOVL (BP), DI - catch return to invalid address, for debugging
	0x83,
	0xe6,
	0xe0, // ANDL $~31, SI
	0x4c,
	0x01,
	0xfe, // ADDQ R15, SI
	0xff,
	0xe6, // JMP SI
}

var naclspfix_asm6 = []uint8{0x4c, 0x01, 0xfc} // ADDQ R15, SP

var naclbpfix_asm6 = []uint8{0x4c, 0x01, 0xfd} // ADDQ R15, BP

var naclmovs_asm6 = []uint8{
	0x89,
	0xf6, // MOVL SI, SI
	0x49,
	0x8d,
	0x34,
	0x37, // LEAQ (R15)(SI*1), SI
	0x89,
	0xff, // MOVL DI, DI
	0x49,
	0x8d,
	0x3c,
	0x3f, // LEAQ (R15)(DI*1), DI
}

var naclstos_asm6 = []uint8{
	0x89,
	0xff, // MOVL DI, DI
	0x49,
	0x8d,
	0x3c,
	0x3f, // LEAQ (R15)(DI*1), DI
}

func nacltrunc_asm6(ctxt *Link, reg int) {
	if reg >= int(D_R8_6) {
		ctxt.andptr[0] = 0x45
		ctxt.andptr = ctxt.andptr[1:]
	}
	reg = (reg - int(D_AX_6)) & 7
	ctxt.andptr[0] = 0x89
	ctxt.andptr = ctxt.andptr[1:]
	ctxt.andptr[0] = uint8((3 << 6) | (reg << 3) | reg)
	ctxt.andptr = ctxt.andptr[1:]
}
