package main

import (
	"fmt"
	"log"
)

/*
 * this is the ranlib header
 */
const (
	MaxAlign_asm8  = 32
	FuncAlign_asm8 = 16
)

type Optab_asm8 struct {
	as     int
	ytab   []uint8
	prefix int
	op     [13]uint8
}

const (
	Yxxx_asm8 = 0 + iota
	Ynone_asm8
	Yi0_asm8
	Yi1_asm8
	Yi8_asm8
	Yi32_asm8
	Yiauto_asm8
	Yal_asm8
	Ycl_asm8
	Yax_asm8
	Ycx_asm8
	Yrb_asm8
	Yrl_asm8
	Yrf_asm8
	Yf0_asm8
	Yrx_asm8
	Ymb_asm8
	Yml_asm8
	Ym_asm8
	Ybr_asm8
	Ycol_asm8
	Ytls_asm8
	Ycs_asm8
	Yss_asm8
	Yds_asm8
	Yes_asm8
	Yfs_asm8
	Ygs_asm8
	Ygdtr_asm8
	Yidtr_asm8
	Yldtr_asm8
	Ymsw_asm8
	Ytask_asm8
	Ycr0_asm8
	Ycr1_asm8
	Ycr2_asm8
	Ycr3_asm8
	Ycr4_asm8
	Ycr5_asm8
	Ycr6_asm8
	Ycr7_asm8
	Ydr0_asm8
	Ydr1_asm8
	Ydr2_asm8
	Ydr3_asm8
	Ydr4_asm8
	Ydr5_asm8
	Ydr6_asm8
	Ydr7_asm8
	Ytr0_asm8
	Ytr1_asm8
	Ytr2_asm8
	Ytr3_asm8
	Ytr4_asm8
	Ytr5_asm8
	Ytr6_asm8
	Ytr7_asm8
	Ymr_asm8
	Ymm_asm8
	Yxr_asm8
	Yxm_asm8
	Ymax_asm8
	Zxxx_asm8 = 0 + iota - 62
	Zlit_asm8
	Zlitm_r_asm8
	Z_rp_asm8
	Zbr_asm8
	Zcall_asm8
	Zcallcon_asm8
	Zcallind_asm8
	Zcallindreg_asm8
	Zib__asm8
	Zib_rp_asm8
	Zibo_m_asm8
	Zil__asm8
	Zil_rp_asm8
	Zilo_m_asm8
	Zjmp_asm8
	Zjmpcon_asm8
	Zloop_asm8
	Zm_o_asm8
	Zm_r_asm8
	Zm2_r_asm8
	Zm_r_xm_asm8
	Zm_r_i_xm_asm8
	Zaut_r_asm8
	Zo_m_asm8
	Zpseudo_asm8
	Zr_m_asm8
	Zr_m_xm_asm8
	Zr_m_i_xm_asm8
	Zrp__asm8
	Z_ib_asm8
	Z_il_asm8
	Zm_ibo_asm8
	Zm_ilo_asm8
	Zib_rr_asm8
	Zil_rr_asm8
	Zclr_asm8
	Zibm_r_asm8
	Zbyte_asm8
	Zmov_asm8
	Zmax_asm8
	Px_asm8  = 0
	Pe_asm8  = 0x66
	Pm_asm8  = 0x0f
	Pq_asm8  = 0xff
	Pb_asm8  = 0xfe
	Pf2_asm8 = 0xf2
	Pf3_asm8 = 0xf3
)

var ycover_asm8 [Ymax_asm8 * Ymax_asm8]uint8

var reg_asm8 [D_NONE_8]int

var ynone_asm8 = []uint8{
	Ynone_asm8,
	Ynone_asm8,
	Zlit_asm8,
	1,
	0,
}

var ytext_asm8 = []uint8{
	Ymb_asm8,
	Yi32_asm8,
	Zpseudo_asm8,
	1,
	0,
}

var ynop_asm8 = []uint8{
	Ynone_asm8,
	Ynone_asm8,
	Zpseudo_asm8,
	0,
	Ynone_asm8,
	Yiauto_asm8,
	Zpseudo_asm8,
	0,
	Ynone_asm8,
	Yml_asm8,
	Zpseudo_asm8,
	0,
	Ynone_asm8,
	Yrf_asm8,
	Zpseudo_asm8,
	0,
	Yiauto_asm8,
	Ynone_asm8,
	Zpseudo_asm8,
	0,
	Ynone_asm8,
	Yxr_asm8,
	Zpseudo_asm8,
	0,
	Yml_asm8,
	Ynone_asm8,
	Zpseudo_asm8,
	0,
	Yrf_asm8,
	Ynone_asm8,
	Zpseudo_asm8,
	0,
	Yxr_asm8,
	Ynone_asm8,
	Zpseudo_asm8,
	1,
	0,
}

var yfuncdata_asm8 = []uint8{
	Yi32_asm8,
	Ym_asm8,
	Zpseudo_asm8,
	0,
	0,
}

var ypcdata_asm8 = []uint8{
	Yi32_asm8,
	Yi32_asm8,
	Zpseudo_asm8,
	0,
	0,
}

var yxorb_asm8 = []uint8{
	Yi32_asm8,
	Yal_asm8,
	Zib__asm8,
	1,
	Yi32_asm8,
	Ymb_asm8,
	Zibo_m_asm8,
	2,
	Yrb_asm8,
	Ymb_asm8,
	Zr_m_asm8,
	1,
	Ymb_asm8,
	Yrb_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yxorl_asm8 = []uint8{
	Yi8_asm8,
	Yml_asm8,
	Zibo_m_asm8,
	2,
	Yi32_asm8,
	Yax_asm8,
	Zil__asm8,
	1,
	Yi32_asm8,
	Yml_asm8,
	Zilo_m_asm8,
	2,
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yaddl_asm8 = []uint8{
	Yi8_asm8,
	Yml_asm8,
	Zibo_m_asm8,
	2,
	Yi32_asm8,
	Yax_asm8,
	Zil__asm8,
	1,
	Yi32_asm8,
	Yml_asm8,
	Zilo_m_asm8,
	2,
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yincb_asm8 = []uint8{
	Ynone_asm8,
	Ymb_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yincl_asm8 = []uint8{
	Ynone_asm8,
	Yrl_asm8,
	Z_rp_asm8,
	1,
	Ynone_asm8,
	Yml_asm8,
	Zo_m_asm8,
	2,
	0,
}

var ycmpb_asm8 = []uint8{
	Yal_asm8,
	Yi32_asm8,
	Z_ib_asm8,
	1,
	Ymb_asm8,
	Yi32_asm8,
	Zm_ibo_asm8,
	2,
	Ymb_asm8,
	Yrb_asm8,
	Zm_r_asm8,
	1,
	Yrb_asm8,
	Ymb_asm8,
	Zr_m_asm8,
	1,
	0,
}

var ycmpl_asm8 = []uint8{
	Yml_asm8,
	Yi8_asm8,
	Zm_ibo_asm8,
	2,
	Yax_asm8,
	Yi32_asm8,
	Z_il_asm8,
	1,
	Yml_asm8,
	Yi32_asm8,
	Zm_ilo_asm8,
	2,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	0,
}

var yshb_asm8 = []uint8{
	Yi1_asm8,
	Ymb_asm8,
	Zo_m_asm8,
	2,
	Yi32_asm8,
	Ymb_asm8,
	Zibo_m_asm8,
	2,
	Ycx_asm8,
	Ymb_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yshl_asm8 = []uint8{
	Yi1_asm8,
	Yml_asm8,
	Zo_m_asm8,
	2,
	Yi32_asm8,
	Yml_asm8,
	Zibo_m_asm8,
	2,
	Ycl_asm8,
	Yml_asm8,
	Zo_m_asm8,
	2,
	Ycx_asm8,
	Yml_asm8,
	Zo_m_asm8,
	2,
	0,
}

var ytestb_asm8 = []uint8{
	Yi32_asm8,
	Yal_asm8,
	Zib__asm8,
	1,
	Yi32_asm8,
	Ymb_asm8,
	Zibo_m_asm8,
	2,
	Yrb_asm8,
	Ymb_asm8,
	Zr_m_asm8,
	1,
	Ymb_asm8,
	Yrb_asm8,
	Zm_r_asm8,
	1,
	0,
}

var ytestl_asm8 = []uint8{
	Yi32_asm8,
	Yax_asm8,
	Zil__asm8,
	1,
	Yi32_asm8,
	Yml_asm8,
	Zilo_m_asm8,
	2,
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var ymovb_asm8 = []uint8{
	Yrb_asm8,
	Ymb_asm8,
	Zr_m_asm8,
	1,
	Ymb_asm8,
	Yrb_asm8,
	Zm_r_asm8,
	1,
	Yi32_asm8,
	Yrb_asm8,
	Zib_rp_asm8,
	1,
	Yi32_asm8,
	Ymb_asm8,
	Zibo_m_asm8,
	2,
	0,
}

var ymovw_asm8 = []uint8{
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	Yi0_asm8,
	Yrl_asm8,
	Zclr_asm8,
	1 + 2,
	//	Yi0,	Yml,	Zibo_m,	2,	// shorter but slower AND $0,dst
	Yi32_asm8,
	Yrl_asm8,
	Zil_rp_asm8,
	1,
	Yi32_asm8,
	Yml_asm8,
	Zilo_m_asm8,
	2,
	Yiauto_asm8,
	Yrl_asm8,
	Zaut_r_asm8,
	1,
	0,
}

var ymovl_asm8 = []uint8{
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	Yi0_asm8,
	Yrl_asm8,
	Zclr_asm8,
	1 + 2,
	//	Yi0,	Yml,	Zibo_m,	2,	// shorter but slower AND $0,dst
	Yi32_asm8,
	Yrl_asm8,
	Zil_rp_asm8,
	1,
	Yi32_asm8,
	Yml_asm8,
	Zilo_m_asm8,
	2,
	Yml_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	2, // XMM MOVD (32 bit)
	Yxr_asm8,
	Yml_asm8,
	Zr_m_xm_asm8,
	2, // XMM MOVD (32 bit)
	Yiauto_asm8,
	Yrl_asm8,
	Zaut_r_asm8,
	1,
	0,
}

var ymovq_asm8 = []uint8{
	Yml_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	2,
	0,
}

var ym_rl_asm8 = []uint8{
	Ym_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yrl_m_asm8 = []uint8{
	Yrl_asm8,
	Ym_asm8,
	Zr_m_asm8,
	1,
	0,
}

var ymb_rl_asm8 = []uint8{
	Ymb_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yml_rl_asm8 = []uint8{
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yrb_mb_asm8 = []uint8{
	Yrb_asm8,
	Ymb_asm8,
	Zr_m_asm8,
	1,
	0,
}

var yrl_ml_asm8 = []uint8{
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	0,
}

var yml_mb_asm8 = []uint8{
	Yrb_asm8,
	Ymb_asm8,
	Zr_m_asm8,
	1,
	Ymb_asm8,
	Yrb_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yxchg_asm8 = []uint8{
	Yax_asm8,
	Yrl_asm8,
	Z_rp_asm8,
	1,
	Yrl_asm8,
	Yax_asm8,
	Zrp__asm8,
	1,
	Yrl_asm8,
	Yml_asm8,
	Zr_m_asm8,
	1,
	Yml_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var ydivl_asm8 = []uint8{
	Yml_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	0,
}

var ydivb_asm8 = []uint8{
	Ymb_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	0,
}

var yimul_asm8 = []uint8{
	Yml_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	Yi8_asm8,
	Yrl_asm8,
	Zib_rr_asm8,
	1,
	Yi32_asm8,
	Yrl_asm8,
	Zil_rr_asm8,
	1,
	0,
}

var ybyte_asm8 = []uint8{
	Yi32_asm8,
	Ynone_asm8,
	Zbyte_asm8,
	1,
	0,
}

var yin_asm8 = []uint8{
	Yi32_asm8,
	Ynone_asm8,
	Zib__asm8,
	1,
	Ynone_asm8,
	Ynone_asm8,
	Zlit_asm8,
	1,
	0,
}

var yint_asm8 = []uint8{
	Yi32_asm8,
	Ynone_asm8,
	Zib__asm8,
	1,
	0,
}

var ypushl_asm8 = []uint8{
	Yrl_asm8,
	Ynone_asm8,
	Zrp__asm8,
	1,
	Ym_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	Yi8_asm8,
	Ynone_asm8,
	Zib__asm8,
	1,
	Yi32_asm8,
	Ynone_asm8,
	Zil__asm8,
	1,
	0,
}

var ypopl_asm8 = []uint8{
	Ynone_asm8,
	Yrl_asm8,
	Z_rp_asm8,
	1,
	Ynone_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	0,
}

var ybswap_asm8 = []uint8{
	Ynone_asm8,
	Yrl_asm8,
	Z_rp_asm8,
	1,
	0,
}

var yscond_asm8 = []uint8{
	Ynone_asm8,
	Ymb_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yjcond_asm8 = []uint8{
	Ynone_asm8,
	Ybr_asm8,
	Zbr_asm8,
	0,
	Yi0_asm8,
	Ybr_asm8,
	Zbr_asm8,
	0,
	Yi1_asm8,
	Ybr_asm8,
	Zbr_asm8,
	1,
	0,
}

var yloop_asm8 = []uint8{
	Ynone_asm8,
	Ybr_asm8,
	Zloop_asm8,
	1,
	0,
}

var ycall_asm8 = []uint8{
	Ynone_asm8,
	Yml_asm8,
	Zcallindreg_asm8,
	0,
	Yrx_asm8,
	Yrx_asm8,
	Zcallindreg_asm8,
	2,
	Ynone_asm8,
	Ycol_asm8,
	Zcallind_asm8,
	2,
	Ynone_asm8,
	Ybr_asm8,
	Zcall_asm8,
	0,
	Ynone_asm8,
	Yi32_asm8,
	Zcallcon_asm8,
	1,
	0,
}

var yduff_asm8 = []uint8{
	Ynone_asm8,
	Yi32_asm8,
	Zcall_asm8,
	1,
	0,
}

var yjmp_asm8 = []uint8{
	Ynone_asm8,
	Yml_asm8,
	Zo_m_asm8,
	2,
	Ynone_asm8,
	Ybr_asm8,
	Zjmp_asm8,
	0,
	Ynone_asm8,
	Yi32_asm8,
	Zjmpcon_asm8,
	1,
	0,
}

var yfmvd_asm8 = []uint8{
	Ym_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	Yf0_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	Yrf_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	Yf0_asm8,
	Yrf_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yfmvdp_asm8 = []uint8{
	Yf0_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	Yf0_asm8,
	Yrf_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yfmvf_asm8 = []uint8{
	Ym_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	Yf0_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yfmvx_asm8 = []uint8{
	Ym_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	0,
}

var yfmvp_asm8 = []uint8{
	Yf0_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yfcmv_asm8 = []uint8{
	Yrf_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	0,
}

var yfadd_asm8 = []uint8{
	Ym_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	Yrf_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	Yf0_asm8,
	Yrf_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yfaddp_asm8 = []uint8{
	Yf0_asm8,
	Yrf_asm8,
	Zo_m_asm8,
	2,
	0,
}

var yfxch_asm8 = []uint8{
	Yf0_asm8,
	Yrf_asm8,
	Zo_m_asm8,
	2,
	Yrf_asm8,
	Yf0_asm8,
	Zm_o_asm8,
	2,
	0,
}

var ycompp_asm8 = []uint8{
	Yf0_asm8,
	Yrf_asm8,
	Zo_m_asm8,
	2, /* botch is really f0,f1 */
	0,
}

var ystsw_asm8 = []uint8{
	Ynone_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	Ynone_asm8,
	Yax_asm8,
	Zlit_asm8,
	1,
	0,
}

var ystcw_asm8 = []uint8{
	Ynone_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	Ym_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	0,
}

var ysvrs_asm8 = []uint8{
	Ynone_asm8,
	Ym_asm8,
	Zo_m_asm8,
	2,
	Ym_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	0,
}

var ymskb_asm8 = []uint8{
	Yxr_asm8,
	Yrl_asm8,
	Zm_r_xm_asm8,
	2,
	Ymr_asm8,
	Yrl_asm8,
	Zm_r_xm_asm8,
	1,
	0,
}

var yxm_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	1,
	0,
}

var yxcvm1_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	2,
	Yxm_asm8,
	Ymr_asm8,
	Zm_r_xm_asm8,
	2,
	0,
}

var yxcvm2_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	2,
	Ymm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	2,
	0,
}

var yxmq_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	2,
	0,
}

var yxr_asm8 = []uint8{
	Yxr_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	1,
	0,
}

var yxr_ml_asm8 = []uint8{
	Yxr_asm8,
	Yml_asm8,
	Zr_m_xm_asm8,
	1,
	0,
}

var yxcmp_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	1,
	0,
}

var yxcmpi_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_i_xm_asm8,
	2,
	0,
}

var yxmov_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	1,
	Yxr_asm8,
	Yxm_asm8,
	Zr_m_xm_asm8,
	1,
	0,
}

var yxcvfl_asm8 = []uint8{
	Yxm_asm8,
	Yrl_asm8,
	Zm_r_xm_asm8,
	1,
	0,
}

var yxcvlf_asm8 = []uint8{
	Yml_asm8,
	Yxr_asm8,
	Zm_r_xm_asm8,
	1,
	0,
}

/*
static uchar	yxcvfq[] =
{
	Yxm,	Yrl,	Zm_r_xm,	2,
	0
};
static uchar	yxcvqf[] =
{
	Yml,	Yxr,	Zm_r_xm,	2,
	0
};
*/
var yxrrl_asm8 = []uint8{
	Yxr_asm8,
	Yrl_asm8,
	Zm_r_asm8,
	1,
	0,
}

var yprefetch_asm8 = []uint8{
	Ym_asm8,
	Ynone_asm8,
	Zm_o_asm8,
	2,
	0,
}

var yaes_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zlitm_r_asm8,
	2,
	0,
}

var yinsrd_asm8 = []uint8{
	Yml_asm8,
	Yxr_asm8,
	Zibm_r_asm8,
	2,
	0,
}

var ymshufb_asm8 = []uint8{
	Yxm_asm8,
	Yxr_asm8,
	Zm2_r_asm8,
	2,
	0,
}

var optab_asm8 = /*	as, ytab, andproto, opcode */
[]Optab_asm8{
	{AXXX_8, nil, 0, [13]uint8{}},
	{AAAA_8, ynone_asm8, Px_asm8, [13]uint8{0x37}},
	{AAAD_8, ynone_asm8, Px_asm8, [13]uint8{0xd5, 0x0a}},
	{AAAM_8, ynone_asm8, Px_asm8, [13]uint8{0xd4, 0x0a}},
	{AAAS_8, ynone_asm8, Px_asm8, [13]uint8{0x3f}},
	{AADCB_8, yxorb_asm8, Pb_asm8, [13]uint8{0x14, 0x80, 02, 0x10, 0x10}},
	{AADCL_8, yxorl_asm8, Px_asm8, [13]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCW_8, yxorl_asm8, Pe_asm8, [13]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADDB_8, yxorb_asm8, Px_asm8, [13]uint8{0x04, 0x80, 00, 0x00, 0x02}},
	{AADDL_8, yaddl_asm8, Px_asm8, [13]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDW_8, yaddl_asm8, Pe_asm8, [13]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADJSP_8, nil, 0, [13]uint8{}},
	{AANDB_8, yxorb_asm8, Pb_asm8, [13]uint8{0x24, 0x80, 04, 0x20, 0x22}},
	{AANDL_8, yxorl_asm8, Px_asm8, [13]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDW_8, yxorl_asm8, Pe_asm8, [13]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AARPL_8, yrl_ml_asm8, Px_asm8, [13]uint8{0x63}},
	{ABOUNDL_8, yrl_m_asm8, Px_asm8, [13]uint8{0x62}},
	{ABOUNDW_8, yrl_m_asm8, Pe_asm8, [13]uint8{0x62}},
	{ABSFL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xbc}},
	{ABSFW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0xbc}},
	{ABSRL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xbd}},
	{ABSRW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0xbd}},
	{ABTL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xa3}},
	{ABTW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0xa3}},
	{ABTCL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xbb}},
	{ABTCW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0xbb}},
	{ABTRL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xb3}},
	{ABTRW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0xb3}},
	{ABTSL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xab}},
	{ABTSW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0xab}},
	{ABYTE_8, ybyte_asm8, Px_asm8, [13]uint8{1}},
	{ACALL_8, ycall_asm8, Px_asm8, [13]uint8{0xff, 02, 0xff, 0x15, 0xe8}},
	{ACLC_8, ynone_asm8, Px_asm8, [13]uint8{0xf8}},
	{ACLD_8, ynone_asm8, Px_asm8, [13]uint8{0xfc}},
	{ACLI_8, ynone_asm8, Px_asm8, [13]uint8{0xfa}},
	{ACLTS_8, ynone_asm8, Pm_asm8, [13]uint8{0x06}},
	{ACMC_8, ynone_asm8, Px_asm8, [13]uint8{0xf5}},
	{ACMPB_8, ycmpb_asm8, Pb_asm8, [13]uint8{0x3c, 0x80, 07, 0x38, 0x3a}},
	{ACMPL_8, ycmpl_asm8, Px_asm8, [13]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPW_8, ycmpl_asm8, Pe_asm8, [13]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPSB_8, ynone_asm8, Pb_asm8, [13]uint8{0xa6}},
	{ACMPSL_8, ynone_asm8, Px_asm8, [13]uint8{0xa7}},
	{ACMPSW_8, ynone_asm8, Pe_asm8, [13]uint8{0xa7}},
	{ADAA_8, ynone_asm8, Px_asm8, [13]uint8{0x27}},
	{ADAS_8, ynone_asm8, Px_asm8, [13]uint8{0x2f}},
	{ADATA_8, nil, 0, [13]uint8{}},
	{ADECB_8, yincb_asm8, Pb_asm8, [13]uint8{0xfe, 01}},
	{ADECL_8, yincl_asm8, Px_asm8, [13]uint8{0x48, 0xff, 01}},
	{ADECW_8, yincl_asm8, Pe_asm8, [13]uint8{0x48, 0xff, 01}},
	{ADIVB_8, ydivb_asm8, Pb_asm8, [13]uint8{0xf6, 06}},
	{ADIVL_8, ydivl_asm8, Px_asm8, [13]uint8{0xf7, 06}},
	{ADIVW_8, ydivl_asm8, Pe_asm8, [13]uint8{0xf7, 06}},
	{AENTER_8, nil, 0, [13]uint8{}}, /* botch */
	{AGLOBL_8, nil, 0, [13]uint8{}},
	{AGOK_8, nil, 0, [13]uint8{}},
	{AHISTORY_8, nil, 0, [13]uint8{}},
	{AHLT_8, ynone_asm8, Px_asm8, [13]uint8{0xf4}},
	{AIDIVB_8, ydivb_asm8, Pb_asm8, [13]uint8{0xf6, 07}},
	{AIDIVL_8, ydivl_asm8, Px_asm8, [13]uint8{0xf7, 07}},
	{AIDIVW_8, ydivl_asm8, Pe_asm8, [13]uint8{0xf7, 07}},
	{AIMULB_8, ydivb_asm8, Pb_asm8, [13]uint8{0xf6, 05}},
	{AIMULL_8, yimul_asm8, Px_asm8, [13]uint8{0xf7, 05, 0x6b, 0x69}},
	{AIMULW_8, yimul_asm8, Pe_asm8, [13]uint8{0xf7, 05, 0x6b, 0x69}},
	{AINB_8, yin_asm8, Pb_asm8, [13]uint8{0xe4, 0xec}},
	{AINL_8, yin_asm8, Px_asm8, [13]uint8{0xe5, 0xed}},
	{AINW_8, yin_asm8, Pe_asm8, [13]uint8{0xe5, 0xed}},
	{AINCB_8, yincb_asm8, Pb_asm8, [13]uint8{0xfe, 00}},
	{AINCL_8, yincl_asm8, Px_asm8, [13]uint8{0x40, 0xff, 00}},
	{AINCW_8, yincl_asm8, Pe_asm8, [13]uint8{0x40, 0xff, 00}},
	{AINSB_8, ynone_asm8, Pb_asm8, [13]uint8{0x6c}},
	{AINSL_8, ynone_asm8, Px_asm8, [13]uint8{0x6d}},
	{AINSW_8, ynone_asm8, Pe_asm8, [13]uint8{0x6d}},
	{AINT_8, yint_asm8, Px_asm8, [13]uint8{0xcd}},
	{AINTO_8, ynone_asm8, Px_asm8, [13]uint8{0xce}},
	{AIRETL_8, ynone_asm8, Px_asm8, [13]uint8{0xcf}},
	{AIRETW_8, ynone_asm8, Pe_asm8, [13]uint8{0xcf}},
	{AJCC_8, yjcond_asm8, Px_asm8, [13]uint8{0x73, 0x83, 00}},
	{AJCS_8, yjcond_asm8, Px_asm8, [13]uint8{0x72, 0x82}},
	{AJCXZL_8, yloop_asm8, Px_asm8, [13]uint8{0xe3}},
	{AJCXZW_8, yloop_asm8, Px_asm8, [13]uint8{0xe3}},
	{AJEQ_8, yjcond_asm8, Px_asm8, [13]uint8{0x74, 0x84}},
	{AJGE_8, yjcond_asm8, Px_asm8, [13]uint8{0x7d, 0x8d}},
	{AJGT_8, yjcond_asm8, Px_asm8, [13]uint8{0x7f, 0x8f}},
	{AJHI_8, yjcond_asm8, Px_asm8, [13]uint8{0x77, 0x87}},
	{AJLE_8, yjcond_asm8, Px_asm8, [13]uint8{0x7e, 0x8e}},
	{AJLS_8, yjcond_asm8, Px_asm8, [13]uint8{0x76, 0x86}},
	{AJLT_8, yjcond_asm8, Px_asm8, [13]uint8{0x7c, 0x8c}},
	{AJMI_8, yjcond_asm8, Px_asm8, [13]uint8{0x78, 0x88}},
	{AJMP_8, yjmp_asm8, Px_asm8, [13]uint8{0xff, 04, 0xeb, 0xe9}},
	{AJNE_8, yjcond_asm8, Px_asm8, [13]uint8{0x75, 0x85}},
	{AJOC_8, yjcond_asm8, Px_asm8, [13]uint8{0x71, 0x81, 00}},
	{AJOS_8, yjcond_asm8, Px_asm8, [13]uint8{0x70, 0x80, 00}},
	{AJPC_8, yjcond_asm8, Px_asm8, [13]uint8{0x7b, 0x8b}},
	{AJPL_8, yjcond_asm8, Px_asm8, [13]uint8{0x79, 0x89}},
	{AJPS_8, yjcond_asm8, Px_asm8, [13]uint8{0x7a, 0x8a}},
	{ALAHF_8, ynone_asm8, Px_asm8, [13]uint8{0x9f}},
	{ALARL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x02}},
	{ALARW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x02}},
	{ALEAL_8, ym_rl_asm8, Px_asm8, [13]uint8{0x8d}},
	{ALEAW_8, ym_rl_asm8, Pe_asm8, [13]uint8{0x8d}},
	{ALEAVEL_8, ynone_asm8, Px_asm8, [13]uint8{0xc9}},
	{ALEAVEW_8, ynone_asm8, Pe_asm8, [13]uint8{0xc9}},
	{ALOCK_8, ynone_asm8, Px_asm8, [13]uint8{0xf0}},
	{ALODSB_8, ynone_asm8, Pb_asm8, [13]uint8{0xac}},
	{ALODSL_8, ynone_asm8, Px_asm8, [13]uint8{0xad}},
	{ALODSW_8, ynone_asm8, Pe_asm8, [13]uint8{0xad}},
	{ALONG_8, ybyte_asm8, Px_asm8, [13]uint8{4}},
	{ALOOP_8, yloop_asm8, Px_asm8, [13]uint8{0xe2}},
	{ALOOPEQ_8, yloop_asm8, Px_asm8, [13]uint8{0xe1}},
	{ALOOPNE_8, yloop_asm8, Px_asm8, [13]uint8{0xe0}},
	{ALSLL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x03}},
	{ALSLW_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x03}},
	{AMOVB_8, ymovb_asm8, Pb_asm8, [13]uint8{0x88, 0x8a, 0xb0, 0xc6, 00}},
	{AMOVL_8, ymovl_asm8, Px_asm8, [13]uint8{0x89, 0x8b, 0x31, 0x83, 04, 0xb8, 0xc7, 00, Pe_asm8, 0x6e, Pe_asm8, 0x7e, 0}},
	{AMOVW_8, ymovw_asm8, Pe_asm8, [13]uint8{0x89, 0x8b, 0x31, 0x83, 04, 0xb8, 0xc7, 00, 0}},
	{AMOVQ_8, ymovq_asm8, Pf3_asm8, [13]uint8{0x7e}},
	{AMOVBLSX_8, ymb_rl_asm8, Pm_asm8, [13]uint8{0xbe}},
	{AMOVBLZX_8, ymb_rl_asm8, Pm_asm8, [13]uint8{0xb6}},
	{AMOVBWSX_8, ymb_rl_asm8, Pq_asm8, [13]uint8{0xbe}},
	{AMOVBWZX_8, ymb_rl_asm8, Pq_asm8, [13]uint8{0xb6}},
	{AMOVWLSX_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xbf}},
	{AMOVWLZX_8, yml_rl_asm8, Pm_asm8, [13]uint8{0xb7}},
	{AMOVSB_8, ynone_asm8, Pb_asm8, [13]uint8{0xa4}},
	{AMOVSL_8, ynone_asm8, Px_asm8, [13]uint8{0xa5}},
	{AMOVSW_8, ynone_asm8, Pe_asm8, [13]uint8{0xa5}},
	{AMULB_8, ydivb_asm8, Pb_asm8, [13]uint8{0xf6, 04}},
	{AMULL_8, ydivl_asm8, Px_asm8, [13]uint8{0xf7, 04}},
	{AMULW_8, ydivl_asm8, Pe_asm8, [13]uint8{0xf7, 04}},
	{ANAME_8, nil, 0, [13]uint8{}},
	{ANEGB_8, yscond_asm8, Px_asm8, [13]uint8{0xf6, 03}},
	{ANEGL_8, yscond_asm8, Px_asm8, [13]uint8{0xf7, 03}},
	{ANEGW_8, yscond_asm8, Pe_asm8, [13]uint8{0xf7, 03}},
	{ANOP_8, ynop_asm8, Px_asm8, [13]uint8{0, 0}},
	{ANOTB_8, yscond_asm8, Px_asm8, [13]uint8{0xf6, 02}},
	{ANOTL_8, yscond_asm8, Px_asm8, [13]uint8{0xf7, 02}},
	{ANOTW_8, yscond_asm8, Pe_asm8, [13]uint8{0xf7, 02}},
	{AORB_8, yxorb_asm8, Pb_asm8, [13]uint8{0x0c, 0x80, 01, 0x08, 0x0a}},
	{AORL_8, yxorl_asm8, Px_asm8, [13]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORW_8, yxorl_asm8, Pe_asm8, [13]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AOUTB_8, yin_asm8, Pb_asm8, [13]uint8{0xe6, 0xee}},
	{AOUTL_8, yin_asm8, Px_asm8, [13]uint8{0xe7, 0xef}},
	{AOUTW_8, yin_asm8, Pe_asm8, [13]uint8{0xe7, 0xef}},
	{AOUTSB_8, ynone_asm8, Pb_asm8, [13]uint8{0x6e}},
	{AOUTSL_8, ynone_asm8, Px_asm8, [13]uint8{0x6f}},
	{AOUTSW_8, ynone_asm8, Pe_asm8, [13]uint8{0x6f}},
	{APAUSE_8, ynone_asm8, Px_asm8, [13]uint8{0xf3, 0x90}},
	{APOPAL_8, ynone_asm8, Px_asm8, [13]uint8{0x61}},
	{APOPAW_8, ynone_asm8, Pe_asm8, [13]uint8{0x61}},
	{APOPFL_8, ynone_asm8, Px_asm8, [13]uint8{0x9d}},
	{APOPFW_8, ynone_asm8, Pe_asm8, [13]uint8{0x9d}},
	{APOPL_8, ypopl_asm8, Px_asm8, [13]uint8{0x58, 0x8f, 00}},
	{APOPW_8, ypopl_asm8, Pe_asm8, [13]uint8{0x58, 0x8f, 00}},
	{APUSHAL_8, ynone_asm8, Px_asm8, [13]uint8{0x60}},
	{APUSHAW_8, ynone_asm8, Pe_asm8, [13]uint8{0x60}},
	{APUSHFL_8, ynone_asm8, Px_asm8, [13]uint8{0x9c}},
	{APUSHFW_8, ynone_asm8, Pe_asm8, [13]uint8{0x9c}},
	{APUSHL_8, ypushl_asm8, Px_asm8, [13]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHW_8, ypushl_asm8, Pe_asm8, [13]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{ARCLB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 02, 0xc0, 02, 0xd2, 02}},
	{ARCLL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCRB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 03, 0xc0, 03, 0xd2, 03}},
	{ARCRL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{AREP_8, ynone_asm8, Px_asm8, [13]uint8{0xf3}},
	{AREPN_8, ynone_asm8, Px_asm8, [13]uint8{0xf2}},
	{ARET_8, ynone_asm8, Px_asm8, [13]uint8{0xc3}},
	{AROLB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 00, 0xc0, 00, 0xd2, 00}},
	{AROLL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{ARORB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 01, 0xc0, 01, 0xd2, 01}},
	{ARORL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ASAHF_8, ynone_asm8, Px_asm8, [13]uint8{0x9e}},
	{ASALB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASALL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASARB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 07, 0xc0, 07, 0xd2, 07}},
	{ASARL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASBBB_8, yxorb_asm8, Pb_asm8, [13]uint8{0x1c, 0x80, 03, 0x18, 0x1a}},
	{ASBBL_8, yxorl_asm8, Px_asm8, [13]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBW_8, yxorl_asm8, Pe_asm8, [13]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASCASB_8, ynone_asm8, Pb_asm8, [13]uint8{0xae}},
	{ASCASL_8, ynone_asm8, Px_asm8, [13]uint8{0xaf}},
	{ASCASW_8, ynone_asm8, Pe_asm8, [13]uint8{0xaf}},
	{ASETCC_8, yscond_asm8, Pm_asm8, [13]uint8{0x93, 00}},
	{ASETCS_8, yscond_asm8, Pm_asm8, [13]uint8{0x92, 00}},
	{ASETEQ_8, yscond_asm8, Pm_asm8, [13]uint8{0x94, 00}},
	{ASETGE_8, yscond_asm8, Pm_asm8, [13]uint8{0x9d, 00}},
	{ASETGT_8, yscond_asm8, Pm_asm8, [13]uint8{0x9f, 00}},
	{ASETHI_8, yscond_asm8, Pm_asm8, [13]uint8{0x97, 00}},
	{ASETLE_8, yscond_asm8, Pm_asm8, [13]uint8{0x9e, 00}},
	{ASETLS_8, yscond_asm8, Pm_asm8, [13]uint8{0x96, 00}},
	{ASETLT_8, yscond_asm8, Pm_asm8, [13]uint8{0x9c, 00}},
	{ASETMI_8, yscond_asm8, Pm_asm8, [13]uint8{0x98, 00}},
	{ASETNE_8, yscond_asm8, Pm_asm8, [13]uint8{0x95, 00}},
	{ASETOC_8, yscond_asm8, Pm_asm8, [13]uint8{0x91, 00}},
	{ASETOS_8, yscond_asm8, Pm_asm8, [13]uint8{0x90, 00}},
	{ASETPC_8, yscond_asm8, Pm_asm8, [13]uint8{0x96, 00}},
	{ASETPL_8, yscond_asm8, Pm_asm8, [13]uint8{0x99, 00}},
	{ASETPS_8, yscond_asm8, Pm_asm8, [13]uint8{0x9a, 00}},
	{ACDQ_8, ynone_asm8, Px_asm8, [13]uint8{0x99}},
	{ACWD_8, ynone_asm8, Pe_asm8, [13]uint8{0x99}},
	{ASHLB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASHLL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHRB_8, yshb_asm8, Pb_asm8, [13]uint8{0xd0, 05, 0xc0, 05, 0xd2, 05}},
	{ASHRL_8, yshl_asm8, Px_asm8, [13]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRW_8, yshl_asm8, Pe_asm8, [13]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASTC_8, ynone_asm8, Px_asm8, [13]uint8{0xf9}},
	{ASTD_8, ynone_asm8, Px_asm8, [13]uint8{0xfd}},
	{ASTI_8, ynone_asm8, Px_asm8, [13]uint8{0xfb}},
	{ASTOSB_8, ynone_asm8, Pb_asm8, [13]uint8{0xaa}},
	{ASTOSL_8, ynone_asm8, Px_asm8, [13]uint8{0xab}},
	{ASTOSW_8, ynone_asm8, Pe_asm8, [13]uint8{0xab}},
	{ASUBB_8, yxorb_asm8, Pb_asm8, [13]uint8{0x2c, 0x80, 05, 0x28, 0x2a}},
	{ASUBL_8, yaddl_asm8, Px_asm8, [13]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBW_8, yaddl_asm8, Pe_asm8, [13]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASYSCALL_8, ynone_asm8, Px_asm8, [13]uint8{0xcd, 100}},
	{ATESTB_8, ytestb_asm8, Pb_asm8, [13]uint8{0xa8, 0xf6, 00, 0x84, 0x84}},
	{ATESTL_8, ytestl_asm8, Px_asm8, [13]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTW_8, ytestl_asm8, Pe_asm8, [13]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATEXT_8, ytext_asm8, Px_asm8, [13]uint8{}},
	{AVERR_8, ydivl_asm8, Pm_asm8, [13]uint8{0x00, 04}},
	{AVERW_8, ydivl_asm8, Pm_asm8, [13]uint8{0x00, 05}},
	{AWAIT_8, ynone_asm8, Px_asm8, [13]uint8{0x9b}},
	{AWORD_8, ybyte_asm8, Px_asm8, [13]uint8{2}},
	{AXCHGB_8, yml_mb_asm8, Pb_asm8, [13]uint8{0x86, 0x86}},
	{AXCHGL_8, yxchg_asm8, Px_asm8, [13]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGW_8, yxchg_asm8, Pe_asm8, [13]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXLAT_8, ynone_asm8, Px_asm8, [13]uint8{0xd7}},
	{AXORB_8, yxorb_asm8, Pb_asm8, [13]uint8{0x34, 0x80, 06, 0x30, 0x32}},
	{AXORL_8, yxorl_asm8, Px_asm8, [13]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORW_8, yxorl_asm8, Pe_asm8, [13]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AFMOVB_8, yfmvx_asm8, Px_asm8, [13]uint8{0xdf, 04}},
	{AFMOVBP_8, yfmvp_asm8, Px_asm8, [13]uint8{0xdf, 06}},
	{AFMOVD_8, yfmvd_asm8, Px_asm8, [13]uint8{0xdd, 00, 0xdd, 02, 0xd9, 00, 0xdd, 02}},
	{AFMOVDP_8, yfmvdp_asm8, Px_asm8, [13]uint8{0xdd, 03, 0xdd, 03}},
	{AFMOVF_8, yfmvf_asm8, Px_asm8, [13]uint8{0xd9, 00, 0xd9, 02}},
	{AFMOVFP_8, yfmvp_asm8, Px_asm8, [13]uint8{0xd9, 03}},
	{AFMOVL_8, yfmvf_asm8, Px_asm8, [13]uint8{0xdb, 00, 0xdb, 02}},
	{AFMOVLP_8, yfmvp_asm8, Px_asm8, [13]uint8{0xdb, 03}},
	{AFMOVV_8, yfmvx_asm8, Px_asm8, [13]uint8{0xdf, 05}},
	{AFMOVVP_8, yfmvp_asm8, Px_asm8, [13]uint8{0xdf, 07}},
	{AFMOVW_8, yfmvf_asm8, Px_asm8, [13]uint8{0xdf, 00, 0xdf, 02}},
	{AFMOVWP_8, yfmvp_asm8, Px_asm8, [13]uint8{0xdf, 03}},
	{AFMOVX_8, yfmvx_asm8, Px_asm8, [13]uint8{0xdb, 05}},
	{AFMOVXP_8, yfmvp_asm8, Px_asm8, [13]uint8{0xdb, 07}},
	{AFCOMB_8, nil, 0, [13]uint8{}},
	{AFCOMBP_8, nil, 0, [13]uint8{}},
	{AFCOMD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 02, 0xd8, 02, 0xdc, 02}},  /* botch */
	{AFCOMDP_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 03, 0xd8, 03, 0xdc, 03}}, /* botch */
	{AFCOMDPP_8, ycompp_asm8, Px_asm8, [13]uint8{0xde, 03}},
	{AFCOMF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 02}},
	{AFCOMFP_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 03}},
	{AFCOMI_8, yfmvx_asm8, Px_asm8, [13]uint8{0xdb, 06}},
	{AFCOMIP_8, yfmvx_asm8, Px_asm8, [13]uint8{0xdf, 06}},
	{AFCOML_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 02}},
	{AFCOMLP_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 03}},
	{AFCOMW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 02}},
	{AFCOMWP_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 03}},
	{AFUCOM_8, ycompp_asm8, Px_asm8, [13]uint8{0xdd, 04}},
	{AFUCOMI_8, ycompp_asm8, Px_asm8, [13]uint8{0xdb, 05}},
	{AFUCOMIP_8, ycompp_asm8, Px_asm8, [13]uint8{0xdf, 05}},
	{AFUCOMP_8, ycompp_asm8, Px_asm8, [13]uint8{0xdd, 05}},
	{AFUCOMPP_8, ycompp_asm8, Px_asm8, [13]uint8{0xda, 13}},
	{AFADDDP_8, yfaddp_asm8, Px_asm8, [13]uint8{0xde, 00}},
	{AFADDW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 00}},
	{AFADDL_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 00}},
	{AFADDF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 00}},
	{AFADDD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 00, 0xd8, 00, 0xdc, 00}},
	{AFMULDP_8, yfaddp_asm8, Px_asm8, [13]uint8{0xde, 01}},
	{AFMULW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 01}},
	{AFMULL_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 01}},
	{AFMULF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 01}},
	{AFMULD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 01, 0xd8, 01, 0xdc, 01}},
	{AFSUBDP_8, yfaddp_asm8, Px_asm8, [13]uint8{0xde, 05}},
	{AFSUBW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 04}},
	{AFSUBL_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 04}},
	{AFSUBF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 04}},
	{AFSUBD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 04, 0xd8, 04, 0xdc, 05}},
	{AFSUBRDP_8, yfaddp_asm8, Px_asm8, [13]uint8{0xde, 04}},
	{AFSUBRW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 05}},
	{AFSUBRL_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 05}},
	{AFSUBRF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 05}},
	{AFSUBRD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 05, 0xd8, 05, 0xdc, 04}},
	{AFDIVDP_8, yfaddp_asm8, Px_asm8, [13]uint8{0xde, 07}},
	{AFDIVW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 06}},
	{AFDIVL_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 06}},
	{AFDIVF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 06}},
	{AFDIVD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 06, 0xd8, 06, 0xdc, 07}},
	{AFDIVRDP_8, yfaddp_asm8, Px_asm8, [13]uint8{0xde, 06}},
	{AFDIVRW_8, yfmvx_asm8, Px_asm8, [13]uint8{0xde, 07}},
	{AFDIVRL_8, yfmvx_asm8, Px_asm8, [13]uint8{0xda, 07}},
	{AFDIVRF_8, yfmvx_asm8, Px_asm8, [13]uint8{0xd8, 07}},
	{AFDIVRD_8, yfadd_asm8, Px_asm8, [13]uint8{0xdc, 07, 0xd8, 07, 0xdc, 06}},
	{AFXCHD_8, yfxch_asm8, Px_asm8, [13]uint8{0xd9, 01, 0xd9, 01}},
	{AFFREE_8, nil, 0, [13]uint8{}},
	{AFLDCW_8, ystcw_asm8, Px_asm8, [13]uint8{0xd9, 05, 0xd9, 05}},
	{AFLDENV_8, ystcw_asm8, Px_asm8, [13]uint8{0xd9, 04, 0xd9, 04}},
	{AFRSTOR_8, ysvrs_asm8, Px_asm8, [13]uint8{0xdd, 04, 0xdd, 04}},
	{AFSAVE_8, ysvrs_asm8, Px_asm8, [13]uint8{0xdd, 06, 0xdd, 06}},
	{AFSTCW_8, ystcw_asm8, Px_asm8, [13]uint8{0xd9, 07, 0xd9, 07}},
	{AFSTENV_8, ystcw_asm8, Px_asm8, [13]uint8{0xd9, 06, 0xd9, 06}},
	{AFSTSW_8, ystsw_asm8, Px_asm8, [13]uint8{0xdd, 07, 0xdf, 0xe0}},
	{AF2XM1_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf0}},
	{AFABS_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xe1}},
	{AFCHS_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xe0}},
	{AFCLEX_8, ynone_asm8, Px_asm8, [13]uint8{0xdb, 0xe2}},
	{AFCOS_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xff}},
	{AFDECSTP_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf6}},
	{AFINCSTP_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf7}},
	{AFINIT_8, ynone_asm8, Px_asm8, [13]uint8{0xdb, 0xe3}},
	{AFLD1_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xe8}},
	{AFLDL2E_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xea}},
	{AFLDL2T_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xe9}},
	{AFLDLG2_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xec}},
	{AFLDLN2_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xed}},
	{AFLDPI_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xeb}},
	{AFLDZ_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xee}},
	{AFNOP_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xd0}},
	{AFPATAN_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf3}},
	{AFPREM_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf8}},
	{AFPREM1_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf5}},
	{AFPTAN_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf2}},
	{AFRNDINT_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xfc}},
	{AFSCALE_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xfd}},
	{AFSIN_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xfe}},
	{AFSINCOS_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xfb}},
	{AFSQRT_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xfa}},
	{AFTST_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xe4}},
	{AFXAM_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xe5}},
	{AFXTRACT_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf4}},
	{AFYL2X_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf1}},
	{AFYL2XP1_8, ynone_asm8, Px_asm8, [13]uint8{0xd9, 0xf9}},
	{AEND_8, nil, 0, [13]uint8{}},
	{ADYNT__8, nil, 0, [13]uint8{}},
	{AINIT__8, nil, 0, [13]uint8{}},
	{ASIGNAME_8, nil, 0, [13]uint8{}},
	{ACMPXCHGB_8, yrb_mb_asm8, Pm_asm8, [13]uint8{0xb0}},
	{ACMPXCHGL_8, yrl_ml_asm8, Pm_asm8, [13]uint8{0xb1}},
	{ACMPXCHGW_8, yrl_ml_asm8, Pm_asm8, [13]uint8{0xb1}},
	{ACMPXCHG8B_8, yscond_asm8, Pm_asm8, [13]uint8{0xc7, 01}},
	{ACPUID_8, ynone_asm8, Pm_asm8, [13]uint8{0xa2}},
	{ARDTSC_8, ynone_asm8, Pm_asm8, [13]uint8{0x31}},
	{AXADDB_8, yrb_mb_asm8, Pb_asm8, [13]uint8{0x0f, 0xc0}},
	{AXADDL_8, yrl_ml_asm8, Pm_asm8, [13]uint8{0xc1}},
	{AXADDW_8, yrl_ml_asm8, Pe_asm8, [13]uint8{0x0f, 0xc1}},
	{ACMOVLCC_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x43}},
	{ACMOVLCS_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x42}},
	{ACMOVLEQ_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x44}},
	{ACMOVLGE_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x4d}},
	{ACMOVLGT_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x4f}},
	{ACMOVLHI_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x47}},
	{ACMOVLLE_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x4e}},
	{ACMOVLLS_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x46}},
	{ACMOVLLT_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x4c}},
	{ACMOVLMI_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x48}},
	{ACMOVLNE_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x45}},
	{ACMOVLOC_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x41}},
	{ACMOVLOS_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x40}},
	{ACMOVLPC_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x4b}},
	{ACMOVLPL_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x49}},
	{ACMOVLPS_8, yml_rl_asm8, Pm_asm8, [13]uint8{0x4a}},
	{ACMOVWCC_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x43}},
	{ACMOVWCS_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x42}},
	{ACMOVWEQ_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x44}},
	{ACMOVWGE_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x4d}},
	{ACMOVWGT_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x4f}},
	{ACMOVWHI_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x47}},
	{ACMOVWLE_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x4e}},
	{ACMOVWLS_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x46}},
	{ACMOVWLT_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x4c}},
	{ACMOVWMI_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x48}},
	{ACMOVWNE_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x45}},
	{ACMOVWOC_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x41}},
	{ACMOVWOS_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x40}},
	{ACMOVWPC_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x4b}},
	{ACMOVWPL_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x49}},
	{ACMOVWPS_8, yml_rl_asm8, Pq_asm8, [13]uint8{0x4a}},
	{AFCMOVCC_8, yfcmv_asm8, Px_asm8, [13]uint8{0xdb, 00}},
	{AFCMOVCS_8, yfcmv_asm8, Px_asm8, [13]uint8{0xda, 00}},
	{AFCMOVEQ_8, yfcmv_asm8, Px_asm8, [13]uint8{0xda, 01}},
	{AFCMOVHI_8, yfcmv_asm8, Px_asm8, [13]uint8{0xdb, 02}},
	{AFCMOVLS_8, yfcmv_asm8, Px_asm8, [13]uint8{0xda, 02}},
	{AFCMOVNE_8, yfcmv_asm8, Px_asm8, [13]uint8{0xdb, 01}},
	{AFCMOVNU_8, yfcmv_asm8, Px_asm8, [13]uint8{0xdb, 03}},
	{AFCMOVUN_8, yfcmv_asm8, Px_asm8, [13]uint8{0xda, 03}},
	{ALFENCE_8, ynone_asm8, Pm_asm8, [13]uint8{0xae, 0xe8}},
	{AMFENCE_8, ynone_asm8, Pm_asm8, [13]uint8{0xae, 0xf0}},
	{ASFENCE_8, ynone_asm8, Pm_asm8, [13]uint8{0xae, 0xf8}},
	{AEMMS_8, ynone_asm8, Pm_asm8, [13]uint8{0x77}},
	{APREFETCHT0_8, yprefetch_asm8, Pm_asm8, [13]uint8{0x18, 01}},
	{APREFETCHT1_8, yprefetch_asm8, Pm_asm8, [13]uint8{0x18, 02}},
	{APREFETCHT2_8, yprefetch_asm8, Pm_asm8, [13]uint8{0x18, 03}},
	{APREFETCHNTA_8, yprefetch_asm8, Pm_asm8, [13]uint8{0x18, 00}},
	{ABSWAPL_8, ybswap_asm8, Pm_asm8, [13]uint8{0xc8}},
	{AUNDEF_8, ynone_asm8, Px_asm8, [13]uint8{0x0f, 0x0b}},
	{AADDPD_8, yxm_asm8, Pq_asm8, [13]uint8{0x58}},
	{AADDPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x58}},
	{AADDSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x58}},
	{AADDSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x58}},
	{AANDNPD_8, yxm_asm8, Pq_asm8, [13]uint8{0x55}},
	{AANDNPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x55}},
	{AANDPD_8, yxm_asm8, Pq_asm8, [13]uint8{0x54}},
	{AANDPS_8, yxm_asm8, Pq_asm8, [13]uint8{0x54}},
	{ACMPPD_8, yxcmpi_asm8, Px_asm8, [13]uint8{Pe_asm8, 0xc2}},
	{ACMPPS_8, yxcmpi_asm8, Pm_asm8, [13]uint8{0xc2, 0}},
	{ACMPSD_8, yxcmpi_asm8, Px_asm8, [13]uint8{Pf2_asm8, 0xc2}},
	{ACMPSS_8, yxcmpi_asm8, Px_asm8, [13]uint8{Pf3_asm8, 0xc2}},
	{ACOMISD_8, yxcmp_asm8, Pe_asm8, [13]uint8{0x2f}},
	{ACOMISS_8, yxcmp_asm8, Pm_asm8, [13]uint8{0x2f}},
	{ACVTPL2PD_8, yxcvm2_asm8, Px_asm8, [13]uint8{Pf3_asm8, 0xe6, Pe_asm8, 0x2a}},
	{ACVTPL2PS_8, yxcvm2_asm8, Pm_asm8, [13]uint8{0x5b, 0, 0x2a, 0}},
	{ACVTPD2PL_8, yxcvm1_asm8, Px_asm8, [13]uint8{Pf2_asm8, 0xe6, Pe_asm8, 0x2d}},
	{ACVTPD2PS_8, yxm_asm8, Pe_asm8, [13]uint8{0x5a}},
	{ACVTPS2PL_8, yxcvm1_asm8, Px_asm8, [13]uint8{Pe_asm8, 0x5b, Pm_asm8, 0x2d}},
	{ACVTPS2PD_8, yxm_asm8, Pm_asm8, [13]uint8{0x5a}},
	{ACVTSD2SL_8, yxcvfl_asm8, Pf2_asm8, [13]uint8{0x2d}},
	{ACVTSD2SS_8, yxm_asm8, Pf2_asm8, [13]uint8{0x5a}},
	{ACVTSL2SD_8, yxcvlf_asm8, Pf2_asm8, [13]uint8{0x2a}},
	{ACVTSL2SS_8, yxcvlf_asm8, Pf3_asm8, [13]uint8{0x2a}},
	{ACVTSS2SD_8, yxm_asm8, Pf3_asm8, [13]uint8{0x5a}},
	{ACVTSS2SL_8, yxcvfl_asm8, Pf3_asm8, [13]uint8{0x2d}},
	{ACVTTPD2PL_8, yxcvm1_asm8, Px_asm8, [13]uint8{Pe_asm8, 0xe6, Pe_asm8, 0x2c}},
	{ACVTTPS2PL_8, yxcvm1_asm8, Px_asm8, [13]uint8{Pf3_asm8, 0x5b, Pm_asm8, 0x2c}},
	{ACVTTSD2SL_8, yxcvfl_asm8, Pf2_asm8, [13]uint8{0x2c}},
	{ACVTTSS2SL_8, yxcvfl_asm8, Pf3_asm8, [13]uint8{0x2c}},
	{ADIVPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x5e}},
	{ADIVPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x5e}},
	{ADIVSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x5e}},
	{ADIVSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x5e}},
	{AMASKMOVOU_8, yxr_asm8, Pe_asm8, [13]uint8{0xf7}},
	{AMAXPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x5f}},
	{AMAXPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x5f}},
	{AMAXSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x5f}},
	{AMAXSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x5f}},
	{AMINPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x5d}},
	{AMINPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x5d}},
	{AMINSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x5d}},
	{AMINSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x5d}},
	{AMOVAPD_8, yxmov_asm8, Pe_asm8, [13]uint8{0x28, 0x29}},
	{AMOVAPS_8, yxmov_asm8, Pm_asm8, [13]uint8{0x28, 0x29}},
	{AMOVO_8, yxmov_asm8, Pe_asm8, [13]uint8{0x6f, 0x7f}},
	{AMOVOU_8, yxmov_asm8, Pf3_asm8, [13]uint8{0x6f, 0x7f}},
	{AMOVHLPS_8, yxr_asm8, Pm_asm8, [13]uint8{0x12}},
	{AMOVHPD_8, yxmov_asm8, Pe_asm8, [13]uint8{0x16, 0x17}},
	{AMOVHPS_8, yxmov_asm8, Pm_asm8, [13]uint8{0x16, 0x17}},
	{AMOVLHPS_8, yxr_asm8, Pm_asm8, [13]uint8{0x16}},
	{AMOVLPD_8, yxmov_asm8, Pe_asm8, [13]uint8{0x12, 0x13}},
	{AMOVLPS_8, yxmov_asm8, Pm_asm8, [13]uint8{0x12, 0x13}},
	{AMOVMSKPD_8, yxrrl_asm8, Pq_asm8, [13]uint8{0x50}},
	{AMOVMSKPS_8, yxrrl_asm8, Pm_asm8, [13]uint8{0x50}},
	{AMOVNTO_8, yxr_ml_asm8, Pe_asm8, [13]uint8{0xe7}},
	{AMOVNTPD_8, yxr_ml_asm8, Pe_asm8, [13]uint8{0x2b}},
	{AMOVNTPS_8, yxr_ml_asm8, Pm_asm8, [13]uint8{0x2b}},
	{AMOVSD_8, yxmov_asm8, Pf2_asm8, [13]uint8{0x10, 0x11}},
	{AMOVSS_8, yxmov_asm8, Pf3_asm8, [13]uint8{0x10, 0x11}},
	{AMOVUPD_8, yxmov_asm8, Pe_asm8, [13]uint8{0x10, 0x11}},
	{AMOVUPS_8, yxmov_asm8, Pm_asm8, [13]uint8{0x10, 0x11}},
	{AMULPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x59}},
	{AMULPS_8, yxm_asm8, Ym_asm8, [13]uint8{0x59}},
	{AMULSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x59}},
	{AMULSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x59}},
	{AORPD_8, yxm_asm8, Pq_asm8, [13]uint8{0x56}},
	{AORPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x56}},
	{APADDQ_8, yxm_asm8, Pe_asm8, [13]uint8{0xd4}},
	{APAND_8, yxm_asm8, Pe_asm8, [13]uint8{0xdb}},
	{APCMPEQB_8, yxmq_asm8, Pe_asm8, [13]uint8{0x74}},
	{APMAXSW_8, yxm_asm8, Pe_asm8, [13]uint8{0xee}},
	{APMAXUB_8, yxm_asm8, Pe_asm8, [13]uint8{0xde}},
	{APMINSW_8, yxm_asm8, Pe_asm8, [13]uint8{0xea}},
	{APMINUB_8, yxm_asm8, Pe_asm8, [13]uint8{0xda}},
	{APMOVMSKB_8, ymskb_asm8, Px_asm8, [13]uint8{Pe_asm8, 0xd7, 0xd7}},
	{APSADBW_8, yxm_asm8, Pq_asm8, [13]uint8{0xf6}},
	{APSUBB_8, yxm_asm8, Pe_asm8, [13]uint8{0xf8}},
	{APSUBL_8, yxm_asm8, Pe_asm8, [13]uint8{0xfa}},
	{APSUBQ_8, yxm_asm8, Pe_asm8, [13]uint8{0xfb}},
	{APSUBSB_8, yxm_asm8, Pe_asm8, [13]uint8{0xe8}},
	{APSUBSW_8, yxm_asm8, Pe_asm8, [13]uint8{0xe9}},
	{APSUBUSB_8, yxm_asm8, Pe_asm8, [13]uint8{0xd8}},
	{APSUBUSW_8, yxm_asm8, Pe_asm8, [13]uint8{0xd9}},
	{APSUBW_8, yxm_asm8, Pe_asm8, [13]uint8{0xf9}},
	{APUNPCKHQDQ_8, yxm_asm8, Pe_asm8, [13]uint8{0x6d}},
	{APUNPCKLQDQ_8, yxm_asm8, Pe_asm8, [13]uint8{0x6c}},
	{APXOR_8, yxm_asm8, Pe_asm8, [13]uint8{0xef}},
	{ARCPPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x53}},
	{ARCPSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x53}},
	{ARSQRTPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x52}},
	{ARSQRTSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x52}},
	{ASQRTPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x51}},
	{ASQRTPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x51}},
	{ASQRTSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x51}},
	{ASQRTSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x51}},
	{ASUBPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x5c}},
	{ASUBPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x5c}},
	{ASUBSD_8, yxm_asm8, Pf2_asm8, [13]uint8{0x5c}},
	{ASUBSS_8, yxm_asm8, Pf3_asm8, [13]uint8{0x5c}},
	{AUCOMISD_8, yxcmp_asm8, Pe_asm8, [13]uint8{0x2e}},
	{AUCOMISS_8, yxcmp_asm8, Pm_asm8, [13]uint8{0x2e}},
	{AUNPCKHPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x15}},
	{AUNPCKHPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x15}},
	{AUNPCKLPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x14}},
	{AUNPCKLPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x14}},
	{AXORPD_8, yxm_asm8, Pe_asm8, [13]uint8{0x57}},
	{AXORPS_8, yxm_asm8, Pm_asm8, [13]uint8{0x57}},
	{AAESENC_8, yaes_asm8, Pq_asm8, [13]uint8{0x38, 0xdc, 0}},
	{APINSRD_8, yinsrd_asm8, Pq_asm8, [13]uint8{0x3a, 0x22, 00}},
	{APSHUFB_8, ymshufb_asm8, Pq_asm8, [13]uint8{0x38, 0x00}},
	{AUSEFIELD_8, ynop_asm8, Px_asm8, [13]uint8{0, 0}},
	{ATYPE_8, nil, 0, [13]uint8{}},
	{AFUNCDATA_8, yfuncdata_asm8, Px_asm8, [13]uint8{0, 0}},
	{APCDATA_8, ypcdata_asm8, Px_asm8, [13]uint8{0, 0}},
	{ACHECKNIL_8, nil, 0, [13]uint8{}},
	{AVARDEF_8, nil, 0, [13]uint8{}},
	{AVARKILL_8, nil, 0, [13]uint8{}},
	{ADUFFCOPY_8, yduff_asm8, Px_asm8, [13]uint8{0xe8}},
	{ADUFFZERO_8, yduff_asm8, Px_asm8, [13]uint8{0xe8}},
	{0, nil, 0, [13]uint8{}},
}

// single-instruction no-ops of various lengths.
// constructed by hand and disassembled with gdb to verify.
// see http://www.agner.org/optimize/optimizing_assembly.pdf for discussion.
var nop_asm8 = [][16]uint8{
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
func fillnop_asm8(p []uint8, n int) {
	var m int
	for n > 0 {
		m = n
		if m > len(nop_asm8) {
			m = len(nop_asm8)
		}
		copy(p, nop_asm8[m-1][:m])
		p = p[m:]
		n -= m
	}
}

func naclpad_asm8(ctxt *Link, s *LSym, c int64, pad int) int64 {
	symgrow(ctxt, s, c+int64(pad))
	fillnop_asm8(s.p[c:], pad)
	return c + int64(pad)
}

func span8(ctxt *Link, s *LSym) {
	var p *Prog
	var q *Prog
	var c int64
	var v int
	var loop int
	var bp []uint8
	var n int
	var m int
	var i int
	ctxt.cursym = s
	if s.text == nil || s.text.link == nil {
		return
	}
	if ycover_asm8[0] == 0 {
		instinit_asm8()
	}
	for p = s.text; p != nil; p = p.link {
		n = 0
		if p.to.typ == D_BRANCH_8 {
			if p.pcond == nil {
				p.pcond = p
			}
		}
		q = p.pcond
		if q != nil {
			if q.back != 2 {
				n = 1
			}
		}
		p.back = n
		if p.as == AADJSP_8 {
			p.to.typ = D_SP_8
			v = int(-p.from.offset)
			p.from.offset = int64(v)
			p.as = AADDL_8
			if v < 0 {
				p.as = ASUBL_8
				v = -v
				p.from.offset = int64(v)
			}
			if v == 0 {
				p.as = ANOP_8
			}
		}
	}
	for p = s.text; p != nil; p = p.link {
		p.back = 2 // use short branches first time through
		q = p.pcond
		if q != nil && (q.back&2 != 0) {
			p.back |= 1 // backward jump
		}
		if p.as == AADJSP_8 {
			p.to.typ = D_SP_8
			v = int(-p.from.offset)
			p.from.offset = int64(v)
			p.as = AADDL_8
			if v < 0 {
				p.as = ASUBL_8
				v = -v
				p.from.offset = int64(v)
			}
			if v == 0 {
				p.as = ANOP_8
			}
		}
	}
	n = 0
	for {
		loop = 0
		for i = 0; i < len(s.r); i++ {
			s.r[i] = Reloc{}
		}
		s.r = s.r[:0]
		s.p = s.p[:0]
		c = 0
		for p = s.text; p != nil; p = p.link {
			if ctxt.headtype == Hnacl && p.isize > 0 {
				var deferreturn_asm8 *LSym
				if deferreturn_asm8 == nil {
					deferreturn_asm8 = linklookup(ctxt, "runtime.deferreturn", 0)
				}
				// pad everything to avoid crossing 32-byte boundary
				if c>>5 != (c+int64(p.isize)-1)>>5 {
					c = naclpad_asm8(ctxt, s, c, int(-c&31))
				}
				// pad call deferreturn to start at 32-byte boundary
				// so that subtracting 5 in jmpdefer will jump back
				// to that boundary and rerun the call.
				if p.as == ACALL_8 && p.to.sym == deferreturn_asm8 {
					c = naclpad_asm8(ctxt, s, c, int(-c&31))
				}
				// pad call to end at 32-byte boundary
				if p.as == ACALL_8 {
					c = naclpad_asm8(ctxt, s, c, int(-(c+int64(p.isize))&31))
				}
				// the linker treats REP and STOSQ as different instructions
				// but in fact the REP is a prefix on the STOSQ.
				// make sure REP has room for 2 more bytes, so that
				// padding will not be inserted before the next instruction.
				if p.as == AREP_8 && c>>5 != (c+3-1)>>5 {
					c = naclpad_asm8(ctxt, s, c, int(-c&31))
				}
				// same for LOCK.
				// various instructions follow; the longest is 4 bytes.
				// give ourselves 8 bytes so as to avoid surprises.
				if p.as == ALOCK_8 && c>>5 != (c+8-1)>>5 {
					c = naclpad_asm8(ctxt, s, c, int(-c&31))
				}
			}
			p.pc = c
			// process forward jumps to p
			for q = p.comefrom; q != nil; q = q.forwd {
				v = int(p.pc - (q.pc + int64(q.mark)))
				if q.back&2 != 0 { // short
					if v > 127 {
						loop++
						q.back ^= 2
					}
					if q.as == AJCXZW_8 {
						s.p[q.pc+2] = uint8(v)
					} else {
						s.p[q.pc+1] = uint8(v)
					}
				} else {
					bp = s.p[q.pc+int64(q.mark)-4:]
					bp[0] = uint8(v)
					bp = bp[1:]
					bp[0] = uint8(v >> 8)
					bp = bp[1:]
					bp[0] = uint8(v >> 16)
					bp = bp[1:]
					bp[0] = uint8(v >> 24)
				}
			}
			p.comefrom = nil
			p.pc = c
			asmins_asm8(ctxt, p)
			m = -cap(ctxt.andptr) + cap(ctxt.and[:])
			if p.isize != m {
				p.isize = m
				loop++
			}
			symgrow(ctxt, s, p.pc+int64(m))
			copy(s.p[p.pc:], ctxt.and[:m])
			p.mark = m
			c += int64(m)
		}
		n++
		if n > 20 {
			ctxt.diag("span must be looping")
			log.Fatalf("bad code")
		}
		if loop == 0 {
			break
		}
	}
	if ctxt.headtype == Hnacl {
		c = naclpad_asm8(ctxt, s, c, int(-c&31))
	}
	c += -c & (FuncAlign_asm8 - 1)
	s.size = c
	if false { /* debug['a'] > 1 */
		fmt.Printf("span1 %s %d (%d tries)\n %.6x", s.name, s.size, n, 0)
		for i = 0; i < len(s.p); i++ {
			fmt.Printf(" %.2x", s.p[i])
			if i%16 == 15 {
				fmt.Printf("\n  %.6x", uint(i+1))
			}
		}
		if i%16 != 0 {
			fmt.Printf("\n")
		}
		for i = 0; i < len(s.r); i++ {
			var r *Reloc
			r = &s.r[i]
			fmt.Printf(" rel %#.4x/%d %s%+d\n", uint64(r.off), r.siz, r.sym.name, r.add)
		}
	}
}

func instinit_asm8() {
	var i int
	for i = 1; optab_asm8[i].as != 0; i++ {
		if i != optab_asm8[i].as {
			log.Fatalf("phase error in optab: at %v found %v", Aconv_list8(i), Aconv_list8(optab_asm8[i].as))
		}
	}
	for i = 0; i < Ymax_asm8; i++ {
		ycover_asm8[i*Ymax_asm8+i] = 1
	}
	ycover_asm8[Yi0_asm8*Ymax_asm8+Yi8_asm8] = 1
	ycover_asm8[Yi1_asm8*Ymax_asm8+Yi8_asm8] = 1
	ycover_asm8[Yi0_asm8*Ymax_asm8+Yi32_asm8] = 1
	ycover_asm8[Yi1_asm8*Ymax_asm8+Yi32_asm8] = 1
	ycover_asm8[Yi8_asm8*Ymax_asm8+Yi32_asm8] = 1
	ycover_asm8[Yal_asm8*Ymax_asm8+Yrb_asm8] = 1
	ycover_asm8[Ycl_asm8*Ymax_asm8+Yrb_asm8] = 1
	ycover_asm8[Yax_asm8*Ymax_asm8+Yrb_asm8] = 1
	ycover_asm8[Ycx_asm8*Ymax_asm8+Yrb_asm8] = 1
	ycover_asm8[Yrx_asm8*Ymax_asm8+Yrb_asm8] = 1
	ycover_asm8[Yax_asm8*Ymax_asm8+Yrx_asm8] = 1
	ycover_asm8[Ycx_asm8*Ymax_asm8+Yrx_asm8] = 1
	ycover_asm8[Yax_asm8*Ymax_asm8+Yrl_asm8] = 1
	ycover_asm8[Ycx_asm8*Ymax_asm8+Yrl_asm8] = 1
	ycover_asm8[Yrx_asm8*Ymax_asm8+Yrl_asm8] = 1
	ycover_asm8[Yf0_asm8*Ymax_asm8+Yrf_asm8] = 1
	ycover_asm8[Yal_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Ycl_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Yax_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Ycx_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Yrx_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Yrb_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Ym_asm8*Ymax_asm8+Ymb_asm8] = 1
	ycover_asm8[Yax_asm8*Ymax_asm8+Yml_asm8] = 1
	ycover_asm8[Ycx_asm8*Ymax_asm8+Yml_asm8] = 1
	ycover_asm8[Yrx_asm8*Ymax_asm8+Yml_asm8] = 1
	ycover_asm8[Yrl_asm8*Ymax_asm8+Yml_asm8] = 1
	ycover_asm8[Ym_asm8*Ymax_asm8+Yml_asm8] = 1
	ycover_asm8[Yax_asm8*Ymax_asm8+Ymm_asm8] = 1
	ycover_asm8[Ycx_asm8*Ymax_asm8+Ymm_asm8] = 1
	ycover_asm8[Yrx_asm8*Ymax_asm8+Ymm_asm8] = 1
	ycover_asm8[Yrl_asm8*Ymax_asm8+Ymm_asm8] = 1
	ycover_asm8[Ym_asm8*Ymax_asm8+Ymm_asm8] = 1
	ycover_asm8[Ymr_asm8*Ymax_asm8+Ymm_asm8] = 1
	ycover_asm8[Ym_asm8*Ymax_asm8+Yxm_asm8] = 1
	ycover_asm8[Yxr_asm8*Ymax_asm8+Yxm_asm8] = 1
	for i = 0; i < D_NONE_8; i++ {
		reg_asm8[i] = -1
		if i >= D_AL_8 && i <= D_BH_8 {
			reg_asm8[i] = (i - D_AL_8) & 7
		}
		if i >= D_AX_8 && i <= D_DI_8 {
			reg_asm8[i] = (i - D_AX_8) & 7
		}
		if i >= D_F0_8 && i <= D_F0_8+7 {
			reg_asm8[i] = (i - D_F0_8) & 7
		}
		if i >= D_X0_8 && i <= D_X0_8+7 {
			reg_asm8[i] = (i - D_X0_8) & 7
		}
	}
}

func prefixof_asm8(ctxt *Link, a *Addr) int {
	switch a.typ {
	case D_INDIR_8 + D_CS_8:
		return 0x2e
	case D_INDIR_8 + D_DS_8:
		return 0x3e
	case D_INDIR_8 + D_ES_8:
		return 0x26
	case D_INDIR_8 + D_FS_8:
		return 0x64
	case D_INDIR_8 + D_GS_8:
		return 0x65
	// NOTE: Systems listed here should be only systems that
	// support direct TLS references like 8(TLS) implemented as
	// direct references from FS or GS. Systems that require
	// the initial-exec model, where you load the TLS base into
	// a register and then index from that register, do not reach
	// this code and should not be listed.
	case D_INDIR_8 + D_TLS_8:
		switch ctxt.headtype {
		default:
			log.Fatalf("unknown TLS base register for %s", headstr(ctxt.headtype))
		case Hdarwin,
			Hdragonfly,
			Hfreebsd,
			Hnetbsd,
			Hopenbsd:
			return 0x65 // GS
		}
	}
	return 0
}

func oclass_asm8(a *Addr) int {
	var v int
	if (a.typ >= D_INDIR_8 && a.typ < 2*D_INDIR_8) || a.index != D_NONE_8 {
		if a.index != D_NONE_8 && a.scale == 0 {
			if a.typ == D_ADDR_8 {
				switch a.index {
				case D_EXTERN_8,
					D_STATIC_8:
					return Yi32_asm8
				case D_AUTO_8,
					D_PARAM_8:
					return Yiauto_asm8
				}
				return Yxxx_asm8
			}
			//if(a->type == D_INDIR+D_ADDR)
			//	print("*Ycol\n");
			return Ycol_asm8
		}
		return Ym_asm8
	}
	switch a.typ {
	case D_AL_8:
		return Yal_asm8
	case D_AX_8:
		return Yax_asm8
	case D_CL_8,
		D_DL_8,
		D_BL_8,
		D_AH_8,
		D_CH_8,
		D_DH_8,
		D_BH_8:
		return Yrb_asm8
	case D_CX_8:
		return Ycx_asm8
	case D_DX_8,
		D_BX_8:
		return Yrx_asm8
	case D_SP_8,
		D_BP_8,
		D_SI_8,
		D_DI_8:
		return Yrl_asm8
	case D_F0_8 + 0:
		return Yf0_asm8
	case D_F0_8 + 1,
		D_F0_8 + 2,
		D_F0_8 + 3,
		D_F0_8 + 4,
		D_F0_8 + 5,
		D_F0_8 + 6,
		D_F0_8 + 7:
		return Yrf_asm8
	case D_X0_8 + 0,
		D_X0_8 + 1,
		D_X0_8 + 2,
		D_X0_8 + 3,
		D_X0_8 + 4,
		D_X0_8 + 5,
		D_X0_8 + 6,
		D_X0_8 + 7:
		return Yxr_asm8
	case D_NONE_8:
		return Ynone_asm8
	case D_CS_8:
		return Ycs_asm8
	case D_SS_8:
		return Yss_asm8
	case D_DS_8:
		return Yds_asm8
	case D_ES_8:
		return Yes_asm8
	case D_FS_8:
		return Yfs_asm8
	case D_GS_8:
		return Ygs_asm8
	case D_TLS_8:
		return Ytls_asm8
	case D_GDTR_8:
		return Ygdtr_asm8
	case D_IDTR_8:
		return Yidtr_asm8
	case D_LDTR_8:
		return Yldtr_asm8
	case D_MSW_8:
		return Ymsw_asm8
	case D_TASK_8:
		return Ytask_asm8
	case D_CR_8 + 0:
		return Ycr0_asm8
	case D_CR_8 + 1:
		return Ycr1_asm8
	case D_CR_8 + 2:
		return Ycr2_asm8
	case D_CR_8 + 3:
		return Ycr3_asm8
	case D_CR_8 + 4:
		return Ycr4_asm8
	case D_CR_8 + 5:
		return Ycr5_asm8
	case D_CR_8 + 6:
		return Ycr6_asm8
	case D_CR_8 + 7:
		return Ycr7_asm8
	case D_DR_8 + 0:
		return Ydr0_asm8
	case D_DR_8 + 1:
		return Ydr1_asm8
	case D_DR_8 + 2:
		return Ydr2_asm8
	case D_DR_8 + 3:
		return Ydr3_asm8
	case D_DR_8 + 4:
		return Ydr4_asm8
	case D_DR_8 + 5:
		return Ydr5_asm8
	case D_DR_8 + 6:
		return Ydr6_asm8
	case D_DR_8 + 7:
		return Ydr7_asm8
	case D_TR_8 + 0:
		return Ytr0_asm8
	case D_TR_8 + 1:
		return Ytr1_asm8
	case D_TR_8 + 2:
		return Ytr2_asm8
	case D_TR_8 + 3:
		return Ytr3_asm8
	case D_TR_8 + 4:
		return Ytr4_asm8
	case D_TR_8 + 5:
		return Ytr5_asm8
	case D_TR_8 + 6:
		return Ytr6_asm8
	case D_TR_8 + 7:
		return Ytr7_asm8
	case D_EXTERN_8,
		D_STATIC_8,
		D_AUTO_8,
		D_PARAM_8:
		return Ym_asm8
	case D_CONST_8,
		D_CONST2_8,
		D_ADDR_8:
		if a.sym == nil {
			v = int(int32(a.offset))
			if v == 0 {
				return Yi0_asm8
			}
			if v == 1 {
				return Yi1_asm8
			}
			if v >= -128 && v <= 127 {
				return Yi8_asm8
			}
		}
		return Yi32_asm8
	case D_BRANCH_8:
		return Ybr_asm8
	}
	return Yxxx_asm8
}

func asmidx_asm8(ctxt *Link, scale int, index int, base int) {
	var i int
	switch index {
	default:
		goto bad
	case D_NONE_8:
		i = 4 << 3
		goto bas
	case D_AX_8,
		D_CX_8,
		D_DX_8,
		D_BX_8,
		D_BP_8,
		D_SI_8,
		D_DI_8:
		i = reg_asm8[index] << 3
		break
	}
	switch scale {
	default:
		goto bad
	case 1:
		break
	case 2:
		i |= 1 << 6
	case 4:
		i |= 2 << 6
	case 8:
		i |= 3 << 6
		break
	}
bas:
	switch base {
	default:
		goto bad
	case D_NONE_8: /* must be mod=00 */
		i |= 5
	case D_AX_8,
		D_CX_8,
		D_DX_8,
		D_BX_8,
		D_SP_8,
		D_BP_8,
		D_SI_8,
		D_DI_8:
		i |= reg_asm8[base]
		break
	}
	ctxt.andptr[0] = uint8(i)
	ctxt.andptr = ctxt.andptr[1:]
	return
bad:
	ctxt.diag("asmidx: bad address %d,%d,%d", scale, index, base)
	ctxt.andptr[0] = 0
	ctxt.andptr = ctxt.andptr[1:]
	return
}

func put4_asm8(ctxt *Link, v int64) {
	ctxt.andptr[0] = uint8(v)
	ctxt.andptr[1] = uint8(v >> 8)
	ctxt.andptr[2] = uint8(v >> 16)
	ctxt.andptr[3] = uint8(v >> 24)
	ctxt.andptr = ctxt.andptr[4:]
}

func relput4_asm8(ctxt *Link, p *Prog, a *Addr) {
	var v int64
	var rel Reloc
	var r *Reloc
	v = vaddr_asm8(ctxt, a, &rel)
	if rel.siz != 0 {
		if rel.siz != 4 {
			ctxt.diag("bad reloc")
		}
		r = addrel(ctxt.cursym)
		*r = rel
		r.off = p.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
	}
	put4_asm8(ctxt, v)
}

func vaddr_asm8(ctxt *Link, a *Addr, r *Reloc) int64 {
	var t int
	var v int64
	var s *LSym
	if r != nil {
		*r = Reloc{}
	}
	t = a.typ
	v = a.offset
	if t == D_ADDR_8 {
		t = a.index
	}
	switch t {
	case D_STATIC_8,
		D_EXTERN_8:
		s = a.sym
		if s != nil {
			if r == nil {
				ctxt.diag("need reloc for %D", a)
				log.Fatalf("bad code")
			}
			r.typ = R_ADDR
			r.siz = 4
			r.off = -1
			r.sym = s
			r.add = v
			v = 0
		}
	case D_INDIR_8 + D_TLS_8:
		if r == nil {
			ctxt.diag("need reloc for %D", a)
			log.Fatalf("bad code")
		}
		r.typ = R_TLS_LE
		r.siz = 4
		r.off = -1 // caller must fill in
		r.add = v
		v = 0
		break
	}
	return v
}

func asmand_asm8(ctxt *Link, a *Addr, r int) {
	var v int64
	var t int
	var scale int
	var rel Reloc
	v = a.offset
	t = a.typ
	rel.siz = 0
	if a.index != D_NONE_8 && a.index != D_TLS_8 {
		if t < D_INDIR_8 || t >= 2*D_INDIR_8 {
			switch t {
			default:
				goto bad
			case D_STATIC_8,
				D_EXTERN_8:
				t = D_NONE_8
				v = vaddr_asm8(ctxt, a, &rel)
			case D_AUTO_8,
				D_PARAM_8:
				t = D_SP_8
				break
			}
		} else {
			t -= D_INDIR_8
		}
		if t == D_NONE_8 {
			ctxt.andptr[0] = uint8(0<<6 | 4<<0 | r<<3)
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm8(ctxt, int(a.scale), a.index, t)
			goto putrelv
		}
		if v == 0 && rel.siz == 0 && t != D_BP_8 {
			ctxt.andptr[0] = uint8(0<<6 | 4<<0 | r<<3)
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm8(ctxt, int(a.scale), a.index, t)
			return
		}
		if v >= -128 && v < 128 && rel.siz == 0 {
			ctxt.andptr[0] = uint8(1<<6 | 4<<0 | r<<3)
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm8(ctxt, int(a.scale), a.index, t)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
		ctxt.andptr[0] = uint8(2<<6 | 4<<0 | r<<3)
		ctxt.andptr = ctxt.andptr[1:]
		asmidx_asm8(ctxt, int(a.scale), a.index, t)
		goto putrelv
	}
	if t >= D_AL_8 && t <= D_F7_8 || t >= D_X0_8 && t <= D_X7_8 {
		if v != 0 {
			goto bad
		}
		ctxt.andptr[0] = uint8(3<<6 | reg_asm8[t]<<0 | r<<3)
		ctxt.andptr = ctxt.andptr[1:]
		return
	}
	scale = int(a.scale)
	if t < D_INDIR_8 || t >= 2*D_INDIR_8 {
		switch a.typ {
		default:
			goto bad
		case D_STATIC_8,
			D_EXTERN_8:
			t = D_NONE_8
			v = vaddr_asm8(ctxt, a, &rel)
		case D_AUTO_8,
			D_PARAM_8:
			t = D_SP_8
			break
		}
		scale = 1
	} else {
		t -= D_INDIR_8
	}
	if t == D_TLS_8 {
		v = vaddr_asm8(ctxt, a, &rel)
	}
	if t == D_NONE_8 || (D_CS_8 <= t && t <= D_GS_8) || t == D_TLS_8 {
		ctxt.andptr[0] = uint8(0<<6 | 5<<0 | r<<3)
		ctxt.andptr = ctxt.andptr[1:]
		goto putrelv
	}
	if t == D_SP_8 {
		if v == 0 && rel.siz == 0 {
			ctxt.andptr[0] = uint8(0<<6 | 4<<0 | r<<3)
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm8(ctxt, scale, D_NONE_8, t)
			return
		}
		if v >= -128 && v < 128 && rel.siz == 0 {
			ctxt.andptr[0] = uint8(1<<6 | 4<<0 | r<<3)
			ctxt.andptr = ctxt.andptr[1:]
			asmidx_asm8(ctxt, scale, D_NONE_8, t)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
		ctxt.andptr[0] = uint8(2<<6 | 4<<0 | r<<3)
		ctxt.andptr = ctxt.andptr[1:]
		asmidx_asm8(ctxt, scale, D_NONE_8, t)
		goto putrelv
	}
	if t >= D_AX_8 && t <= D_DI_8 {
		if a.index == D_TLS_8 {
			rel = Reloc{}
			rel.typ = R_TLS_IE
			rel.siz = 4
			rel.sym = nil
			rel.add = v
			v = 0
		}
		if v == 0 && rel.siz == 0 && t != D_BP_8 {
			ctxt.andptr[0] = uint8(0<<6 | reg_asm8[t]<<0 | r<<3)
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
		if v >= -128 && v < 128 && rel.siz == 0 {
			ctxt.andptr[0] = uint8(1<<6 | reg_asm8[t]<<0 | r<<3)
			ctxt.andptr[1] = uint8(v)
			ctxt.andptr = ctxt.andptr[2:]
			return
		}
		ctxt.andptr[0] = uint8(2<<6 | reg_asm8[t]<<0 | r<<3)
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
		r.off = ctxt.curp.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
	}
	put4_asm8(ctxt, v)
	return
bad:
	ctxt.diag("asmand: bad address %D", a)
	return
}

const (
	E_asm8 = 0xff
)

var ymovtab_asm8 = []uint8{
	/* push */
	APUSHL_8,
	Ycs_asm8,
	Ynone_asm8,
	0,
	0x0e,
	E_asm8,
	0,
	0,
	APUSHL_8,
	Yss_asm8,
	Ynone_asm8,
	0,
	0x16,
	E_asm8,
	0,
	0,
	APUSHL_8,
	Yds_asm8,
	Ynone_asm8,
	0,
	0x1e,
	E_asm8,
	0,
	0,
	APUSHL_8,
	Yes_asm8,
	Ynone_asm8,
	0,
	0x06,
	E_asm8,
	0,
	0,
	APUSHL_8,
	Yfs_asm8,
	Ynone_asm8,
	0,
	0x0f,
	0xa0,
	E_asm8,
	0,
	APUSHL_8,
	Ygs_asm8,
	Ynone_asm8,
	0,
	0x0f,
	0xa8,
	E_asm8,
	0,
	APUSHW_8,
	Ycs_asm8,
	Ynone_asm8,
	0,
	Pe_asm8,
	0x0e,
	E_asm8,
	0,
	APUSHW_8,
	Yss_asm8,
	Ynone_asm8,
	0,
	Pe_asm8,
	0x16,
	E_asm8,
	0,
	APUSHW_8,
	Yds_asm8,
	Ynone_asm8,
	0,
	Pe_asm8,
	0x1e,
	E_asm8,
	0,
	APUSHW_8,
	Yes_asm8,
	Ynone_asm8,
	0,
	Pe_asm8,
	0x06,
	E_asm8,
	0,
	APUSHW_8,
	Yfs_asm8,
	Ynone_asm8,
	0,
	Pe_asm8,
	0x0f,
	0xa0,
	E_asm8,
	APUSHW_8,
	Ygs_asm8,
	Ynone_asm8,
	0,
	Pe_asm8,
	0x0f,
	0xa8,
	E_asm8,
	/* pop */
	APOPL_8,
	Ynone_asm8,
	Yds_asm8,
	0,
	0x1f,
	E_asm8,
	0,
	0,
	APOPL_8,
	Ynone_asm8,
	Yes_asm8,
	0,
	0x07,
	E_asm8,
	0,
	0,
	APOPL_8,
	Ynone_asm8,
	Yss_asm8,
	0,
	0x17,
	E_asm8,
	0,
	0,
	APOPL_8,
	Ynone_asm8,
	Yfs_asm8,
	0,
	0x0f,
	0xa1,
	E_asm8,
	0,
	APOPL_8,
	Ynone_asm8,
	Ygs_asm8,
	0,
	0x0f,
	0xa9,
	E_asm8,
	0,
	APOPW_8,
	Ynone_asm8,
	Yds_asm8,
	0,
	Pe_asm8,
	0x1f,
	E_asm8,
	0,
	APOPW_8,
	Ynone_asm8,
	Yes_asm8,
	0,
	Pe_asm8,
	0x07,
	E_asm8,
	0,
	APOPW_8,
	Ynone_asm8,
	Yss_asm8,
	0,
	Pe_asm8,
	0x17,
	E_asm8,
	0,
	APOPW_8,
	Ynone_asm8,
	Yfs_asm8,
	0,
	Pe_asm8,
	0x0f,
	0xa1,
	E_asm8,
	APOPW_8,
	Ynone_asm8,
	Ygs_asm8,
	0,
	Pe_asm8,
	0x0f,
	0xa9,
	E_asm8,
	/* mov seg */
	AMOVW_8,
	Yes_asm8,
	Yml_asm8,
	1,
	0x8c,
	0,
	0,
	0,
	AMOVW_8,
	Ycs_asm8,
	Yml_asm8,
	1,
	0x8c,
	1,
	0,
	0,
	AMOVW_8,
	Yss_asm8,
	Yml_asm8,
	1,
	0x8c,
	2,
	0,
	0,
	AMOVW_8,
	Yds_asm8,
	Yml_asm8,
	1,
	0x8c,
	3,
	0,
	0,
	AMOVW_8,
	Yfs_asm8,
	Yml_asm8,
	1,
	0x8c,
	4,
	0,
	0,
	AMOVW_8,
	Ygs_asm8,
	Yml_asm8,
	1,
	0x8c,
	5,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Yes_asm8,
	2,
	0x8e,
	0,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Ycs_asm8,
	2,
	0x8e,
	1,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Yss_asm8,
	2,
	0x8e,
	2,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Yds_asm8,
	2,
	0x8e,
	3,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Yfs_asm8,
	2,
	0x8e,
	4,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Ygs_asm8,
	2,
	0x8e,
	5,
	0,
	0,
	/* mov cr */
	AMOVL_8,
	Ycr0_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x20,
	0,
	0,
	AMOVL_8,
	Ycr2_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x20,
	2,
	0,
	AMOVL_8,
	Ycr3_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x20,
	3,
	0,
	AMOVL_8,
	Ycr4_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x20,
	4,
	0,
	AMOVL_8,
	Yml_asm8,
	Ycr0_asm8,
	4,
	0x0f,
	0x22,
	0,
	0,
	AMOVL_8,
	Yml_asm8,
	Ycr2_asm8,
	4,
	0x0f,
	0x22,
	2,
	0,
	AMOVL_8,
	Yml_asm8,
	Ycr3_asm8,
	4,
	0x0f,
	0x22,
	3,
	0,
	AMOVL_8,
	Yml_asm8,
	Ycr4_asm8,
	4,
	0x0f,
	0x22,
	4,
	0,
	/* mov dr */
	AMOVL_8,
	Ydr0_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x21,
	0,
	0,
	AMOVL_8,
	Ydr6_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x21,
	6,
	0,
	AMOVL_8,
	Ydr7_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x21,
	7,
	0,
	AMOVL_8,
	Yml_asm8,
	Ydr0_asm8,
	4,
	0x0f,
	0x23,
	0,
	0,
	AMOVL_8,
	Yml_asm8,
	Ydr6_asm8,
	4,
	0x0f,
	0x23,
	6,
	0,
	AMOVL_8,
	Yml_asm8,
	Ydr7_asm8,
	4,
	0x0f,
	0x23,
	7,
	0,
	/* mov tr */
	AMOVL_8,
	Ytr6_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x24,
	6,
	0,
	AMOVL_8,
	Ytr7_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x24,
	7,
	0,
	AMOVL_8,
	Yml_asm8,
	Ytr6_asm8,
	4,
	0x0f,
	0x26,
	6,
	E_asm8,
	AMOVL_8,
	Yml_asm8,
	Ytr7_asm8,
	4,
	0x0f,
	0x26,
	7,
	E_asm8,
	/* lgdt, sgdt, lidt, sidt */
	AMOVL_8,
	Ym_asm8,
	Ygdtr_asm8,
	4,
	0x0f,
	0x01,
	2,
	0,
	AMOVL_8,
	Ygdtr_asm8,
	Ym_asm8,
	3,
	0x0f,
	0x01,
	0,
	0,
	AMOVL_8,
	Ym_asm8,
	Yidtr_asm8,
	4,
	0x0f,
	0x01,
	3,
	0,
	AMOVL_8,
	Yidtr_asm8,
	Ym_asm8,
	3,
	0x0f,
	0x01,
	1,
	0,
	/* lldt, sldt */
	AMOVW_8,
	Yml_asm8,
	Yldtr_asm8,
	4,
	0x0f,
	0x00,
	2,
	0,
	AMOVW_8,
	Yldtr_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x00,
	0,
	0,
	/* lmsw, smsw */
	AMOVW_8,
	Yml_asm8,
	Ymsw_asm8,
	4,
	0x0f,
	0x01,
	6,
	0,
	AMOVW_8,
	Ymsw_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x01,
	4,
	0,
	/* ltr, str */
	AMOVW_8,
	Yml_asm8,
	Ytask_asm8,
	4,
	0x0f,
	0x00,
	3,
	0,
	AMOVW_8,
	Ytask_asm8,
	Yml_asm8,
	3,
	0x0f,
	0x00,
	1,
	0,
	/* load full pointer */
	AMOVL_8,
	Yml_asm8,
	Ycol_asm8,
	5,
	0,
	0,
	0,
	0,
	AMOVW_8,
	Yml_asm8,
	Ycol_asm8,
	5,
	Pe_asm8,
	0,
	0,
	0,
	/* double shift */
	ASHLL_8,
	Ycol_asm8,
	Yml_asm8,
	6,
	0xa4,
	0xa5,
	0,
	0,
	ASHRL_8,
	Ycol_asm8,
	Yml_asm8,
	6,
	0xac,
	0xad,
	0,
	0,
	/* extra imul */
	AIMULW_8,
	Yml_asm8,
	Yrl_asm8,
	7,
	Pq_asm8,
	0xaf,
	0,
	0,
	AIMULL_8,
	Yml_asm8,
	Yrl_asm8,
	7,
	Pm_asm8,
	0xaf,
	0,
	0,
	/* load TLS base pointer */
	AMOVL_8,
	Ytls_asm8,
	Yrl_asm8,
	8,
	0,
	0,
	0,
	0,
	0,
}

// byteswapreg returns a byte-addressable register (AX, BX, CX, DX)
// which is not referenced in a->type.
// If a is empty, it returns BX to account for MULB-like instructions
// that might use DX and AX.
func byteswapreg_asm8(ctxt *Link, a *Addr) int {
	var cana int
	var canb int
	var canc int
	var cand int
	cand = 1
	canc = cand
	canb = canc
	cana = canb
	switch a.typ {
	case D_NONE_8:
		cand = 0
		cana = cand
	case D_AX_8,
		D_AL_8,
		D_AH_8,
		D_INDIR_8 + D_AX_8:
		cana = 0
	case D_BX_8,
		D_BL_8,
		D_BH_8,
		D_INDIR_8 + D_BX_8:
		canb = 0
	case D_CX_8,
		D_CL_8,
		D_CH_8,
		D_INDIR_8 + D_CX_8:
		canc = 0
	case D_DX_8,
		D_DL_8,
		D_DH_8,
		D_INDIR_8 + D_DX_8:
		cand = 0
		break
	}
	switch a.index {
	case D_AX_8:
		cana = 0
	case D_BX_8:
		canb = 0
	case D_CX_8:
		canc = 0
	case D_DX_8:
		cand = 0
		break
	}
	if cana != 0 {
		return D_AX_8
	}
	if canb != 0 {
		return D_BX_8
	}
	if canc != 0 {
		return D_CX_8
	}
	if cand != 0 {
		return D_DX_8
	}
	ctxt.diag("impossible byte register")
	log.Fatalf("bad code")
	return 0
}

func subreg_asm8(p *Prog, from int, to int) {
	if false { /* debug['Q'] */
		fmt.Printf("\n%v\ts/%v/%v/\n", p, Rconv_list8(from), Rconv_list8(to))
	}
	if p.from.typ == from {
		p.from.typ = to
		p.ft = 0
	}
	if p.to.typ == from {
		p.to.typ = to
		p.tt = 0
	}
	if p.from.index == from {
		p.from.index = to
		p.ft = 0
	}
	if p.to.index == from {
		p.to.index = to
		p.tt = 0
	}
	from += D_INDIR_8
	if p.from.typ == from {
		p.from.typ = to + D_INDIR_8
		p.ft = 0
	}
	if p.to.typ == from {
		p.to.typ = to + D_INDIR_8
		p.tt = 0
	}
	if false { /* debug['Q'] */
		fmt.Printf("%v\n", p)
	}
}

func mediaop_asm8(ctxt *Link, o *Optab_asm8, op int, osize int, z int) int {
	switch op {
	case Pm_asm8,
		Pe_asm8,
		Pf2_asm8,
		Pf3_asm8:
		if osize != 1 {
			if op != Pm_asm8 {
				ctxt.andptr[0] = uint8(op)
				ctxt.andptr = ctxt.andptr[1:]
			}
			ctxt.andptr[0] = Pm_asm8
			ctxt.andptr = ctxt.andptr[1:]
			z++
			op = int(o.op[z])
			break
		}
		fallthrough
	default:
		if -cap(ctxt.andptr) == -cap(ctxt.and) || ctxt.and[-cap(ctxt.andptr)+cap(ctxt.and[:])-1] != Pm_asm8 {
			ctxt.andptr[0] = Pm_asm8
			ctxt.andptr = ctxt.andptr[1:]
		}
		break
	}
	ctxt.andptr[0] = uint8(op)
	ctxt.andptr = ctxt.andptr[1:]
	return z
}

func doasm_asm8(ctxt *Link, p *Prog) {
	var o *Optab_asm8
	var q *Prog
	var pp Prog
	var t []uint8
	var z int
	var op int
	var ft int
	var tt int
	var breg int
	var v int64
	var pre int
	var rel Reloc
	var r *Reloc
	var a *Addr
	ctxt.curp = p // TODO
	pre = prefixof_asm8(ctxt, &p.from)
	if pre != 0 {
		ctxt.andptr[0] = uint8(pre)
		ctxt.andptr = ctxt.andptr[1:]
	}
	pre = prefixof_asm8(ctxt, &p.to)
	if pre != 0 {
		ctxt.andptr[0] = uint8(pre)
		ctxt.andptr = ctxt.andptr[1:]
	}
	if p.ft == 0 {
		p.ft = uint8(oclass_asm8(&p.from))
	}
	if p.tt == 0 {
		p.tt = uint8(oclass_asm8(&p.to))
	}
	ft = int(p.ft) * Ymax_asm8
	tt = int(p.tt) * Ymax_asm8
	o = &optab_asm8[p.as]
	t = o.ytab
	if t == nil {
		ctxt.diag("asmins: noproto %P", p)
		return
	}
	for z = 0; t[0] != 0; (func() { z += int(t[3]); t = t[4:] })() {
		if ycover_asm8[ft+int(t[0])] != 0 {
			if ycover_asm8[tt+int(t[1])] != 0 {
				goto found
			}
		}
	}
	goto domov
found:
	switch o.prefix {
	case Pq_asm8: /* 16 bit escape and opcode escape */
		ctxt.andptr[0] = Pe_asm8
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = Pm_asm8
		ctxt.andptr = ctxt.andptr[1:]
	case Pf2_asm8, /* xmm opcode escape */
		Pf3_asm8:
		ctxt.andptr[0] = uint8(o.prefix)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = Pm_asm8
		ctxt.andptr = ctxt.andptr[1:]
	case Pm_asm8: /* opcode escape */
		ctxt.andptr[0] = Pm_asm8
		ctxt.andptr = ctxt.andptr[1:]
	case Pe_asm8: /* 16 bit escape */
		ctxt.andptr[0] = Pe_asm8
		ctxt.andptr = ctxt.andptr[1:]
	case Pb_asm8: /* botch */
		break
	}
	op = int(o.op[z])
	switch t[2] {
	default:
		ctxt.diag("asmins: unknown z %d %P", t[2], p)
		return
	case Zpseudo_asm8:
		break
	case Zlit_asm8:
		for ; ; z++ {
			op = int(o.op[z])
			if op == 0 {
				break
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		}
	case Zlitm_r_asm8:
		for ; ; z++ {
			op = int(o.op[z])
			if op == 0 {
				break
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		}
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
	case Zm_r_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
	case Zm2_r_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = o.op[z+1]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
	case Zm_r_xm_asm8:
		mediaop_asm8(ctxt, o, op, int(t[3]), z)
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
	case Zm_r_i_xm_asm8:
		mediaop_asm8(ctxt, o, op, int(t[3]), z)
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
		ctxt.andptr[0] = uint8(p.to.offset)
		ctxt.andptr = ctxt.andptr[1:]
	case Zibm_r_asm8:
		for {
			tmp2 := z
			z++
			op = int(o.op[tmp2])
			if op == 0 {
				break
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		}
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
		ctxt.andptr[0] = uint8(p.to.offset)
		ctxt.andptr = ctxt.andptr[1:]
	case Zaut_r_asm8:
		ctxt.andptr[0] = 0x8d
		ctxt.andptr = ctxt.andptr[1:] /* leal */
		if p.from.typ != D_ADDR_8 {
			ctxt.diag("asmins: Zaut sb type ADDR")
		}
		p.from.typ = p.from.index
		p.from.index = D_NONE_8
		p.ft = 0
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
		p.from.index = p.from.typ
		p.from.typ = D_ADDR_8
		p.ft = 0
	case Zm_o_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, int(o.op[z+1]))
	case Zr_m_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, reg_asm8[p.from.typ])
	case Zr_m_xm_asm8:
		mediaop_asm8(ctxt, o, op, int(t[3]), z)
		asmand_asm8(ctxt, &p.to, reg_asm8[p.from.typ])
	case Zr_m_i_xm_asm8:
		mediaop_asm8(ctxt, o, op, int(t[3]), z)
		asmand_asm8(ctxt, &p.to, reg_asm8[p.from.typ])
		ctxt.andptr[0] = uint8(p.from.offset)
		ctxt.andptr = ctxt.andptr[1:]
	case Zcallindreg_asm8:
		r = addrel(ctxt.cursym)
		r.off = p.pc
		r.typ = R_CALLIND
		r.siz = 0
		fallthrough
	// fallthrough
	case Zo_m_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, int(o.op[z+1]))
	case Zm_ibo_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, int(o.op[z+1]))
		ctxt.andptr[0] = uint8(vaddr_asm8(ctxt, &p.to, nil))
		ctxt.andptr = ctxt.andptr[1:]
	case Zibo_m_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, int(o.op[z+1]))
		ctxt.andptr[0] = uint8(vaddr_asm8(ctxt, &p.from, nil))
		ctxt.andptr = ctxt.andptr[1:]
	case Z_ib_asm8,
		Zib__asm8:
		if t[2] == Zib__asm8 {
			a = &p.from
		} else {
			a = &p.to
		}
		v = vaddr_asm8(ctxt, a, nil)
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(v)
		ctxt.andptr = ctxt.andptr[1:]
	case Zib_rp_asm8:
		ctxt.andptr[0] = uint8(op + reg_asm8[p.to.typ])
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = uint8(vaddr_asm8(ctxt, &p.from, nil))
		ctxt.andptr = ctxt.andptr[1:]
	case Zil_rp_asm8:
		ctxt.andptr[0] = uint8(op + reg_asm8[p.to.typ])
		ctxt.andptr = ctxt.andptr[1:]
		if o.prefix == Pe_asm8 {
			v = vaddr_asm8(ctxt, &p.from, nil)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm8(ctxt, p, &p.from)
		}
	case Zib_rr_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, reg_asm8[p.to.typ])
		ctxt.andptr[0] = uint8(vaddr_asm8(ctxt, &p.from, nil))
		ctxt.andptr = ctxt.andptr[1:]
	case Z_il_asm8,
		Zil__asm8:
		if t[2] == Zil__asm8 {
			a = &p.from
		} else {
			a = &p.to
		}
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		if o.prefix == Pe_asm8 {
			v = vaddr_asm8(ctxt, a, nil)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm8(ctxt, p, a)
		}
	case Zm_ilo_asm8,
		Zilo_m_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		if t[2] == Zilo_m_asm8 {
			a = &p.from
			asmand_asm8(ctxt, &p.to, int(o.op[z+1]))
		} else {
			a = &p.to
			asmand_asm8(ctxt, &p.from, int(o.op[z+1]))
		}
		if o.prefix == Pe_asm8 {
			v = vaddr_asm8(ctxt, a, nil)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm8(ctxt, p, a)
		}
	case Zil_rr_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, reg_asm8[p.to.typ])
		if o.prefix == Pe_asm8 {
			v = vaddr_asm8(ctxt, &p.from, nil)
			ctxt.andptr[0] = uint8(v)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = uint8(v >> 8)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			relput4_asm8(ctxt, p, &p.from)
		}
	case Z_rp_asm8:
		ctxt.andptr[0] = uint8(op + reg_asm8[p.to.typ])
		ctxt.andptr = ctxt.andptr[1:]
	case Zrp__asm8:
		ctxt.andptr[0] = uint8(op + reg_asm8[p.from.typ])
		ctxt.andptr = ctxt.andptr[1:]
	case Zclr_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, reg_asm8[p.to.typ])
	case Zcall_asm8:
		if p.to.sym == nil {
			ctxt.diag("call without target")
			log.Fatalf("bad code")
		}
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		r = addrel(ctxt.cursym)
		r.off = p.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
		r.typ = R_CALL
		r.siz = 4
		r.sym = p.to.sym
		r.add = p.to.offset
		put4_asm8(ctxt, 0)
	case Zbr_asm8,
		Zjmp_asm8,
		Zloop_asm8:
		if p.to.sym != nil {
			if t[2] != Zjmp_asm8 {
				ctxt.diag("branch to ATEXT")
				log.Fatalf("bad code")
			}
			ctxt.andptr[0] = o.op[z+1]
			ctxt.andptr = ctxt.andptr[1:]
			r = addrel(ctxt.cursym)
			r.off = p.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
			r.sym = p.to.sym
			r.typ = R_PCREL
			r.siz = 4
			put4_asm8(ctxt, 0)
			break
		}
		// Assumes q is in this function.
		// Fill in backward jump now.
		q = p.pcond
		if q == nil {
			ctxt.diag("jmp/branch/loop without target")
			log.Fatalf("bad code")
		}
		if p.back&1 != 0 {
			v = q.pc - (p.pc + 2)
			if v >= -128 {
				if p.as == AJCXZW_8 {
					ctxt.andptr[0] = 0x67
					ctxt.andptr = ctxt.andptr[1:]
				}
				ctxt.andptr[0] = uint8(op)
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(v)
				ctxt.andptr = ctxt.andptr[1:]
			} else if t[2] == Zloop_asm8 {
				ctxt.diag("loop too far: %P", p)
			} else {
				v -= 5 - 2
				if t[2] == Zbr_asm8 {
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
			break
		}
		// Annotate target; will fill in later.
		p.forwd = q.comefrom
		q.comefrom = p
		if p.back&2 != 0 { // short
			if p.as == AJCXZW_8 {
				ctxt.andptr[0] = 0x67
				ctxt.andptr = ctxt.andptr[1:]
			}
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0
			ctxt.andptr = ctxt.andptr[1:]
		} else if t[2] == Zloop_asm8 {
			ctxt.diag("loop too far: %P", p)
		} else {
			if t[2] == Zbr_asm8 {
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
	case Zcallcon_asm8,
		Zjmpcon_asm8:
		if t[2] == Zcallcon_asm8 {
			ctxt.andptr[0] = uint8(op)
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			ctxt.andptr[0] = o.op[z+1]
			ctxt.andptr = ctxt.andptr[1:]
		}
		r = addrel(ctxt.cursym)
		r.off = p.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
		r.typ = R_PCREL
		r.siz = 4
		r.add = p.to.offset
		put4_asm8(ctxt, 0)
	case Zcallind_asm8:
		ctxt.andptr[0] = uint8(op)
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = o.op[z+1]
		ctxt.andptr = ctxt.andptr[1:]
		r = addrel(ctxt.cursym)
		r.off = p.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
		r.typ = R_ADDR
		r.siz = 4
		r.add = p.to.offset
		r.sym = p.to.sym
		put4_asm8(ctxt, 0)
	case Zbyte_asm8:
		v = vaddr_asm8(ctxt, &p.from, &rel)
		if rel.siz != 0 {
			rel.siz = uint8(op)
			r = addrel(ctxt.cursym)
			*r = rel
			r.off = p.pc + int64(-cap(ctxt.andptr)+cap(ctxt.and[:]))
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
			}
		}
	case Zmov_asm8:
		goto domov
	}
	return
domov:
	for t = ymovtab_asm8; t[0] != 0; t = t[8:] {
		if p.as == int(t[0]) {
			if ycover_asm8[ft+int(t[1])] != 0 {
				if ycover_asm8[tt+int(t[2])] != 0 {
					goto mfound
				}
			}
		}
	}
	/*
	 * here, the assembly has failed.
	 * if its a byte instruction that has
	 * unaddressable registers, try to
	 * exchange registers and reissue the
	 * instruction with the operands renamed.
	 */
bad:
	pp = *p
	z = p.from.typ
	if z >= D_BP_8 && z <= D_DI_8 {
		breg = byteswapreg_asm8(ctxt, &p.to)
		if breg != D_AX_8 {
			ctxt.andptr[0] = 0x87
			ctxt.andptr = ctxt.andptr[1:] /* xchg lhs,bx */
			asmand_asm8(ctxt, &p.from, reg_asm8[breg])
			subreg_asm8(&pp, z, breg)
			doasm_asm8(ctxt, &pp)
			ctxt.andptr[0] = 0x87
			ctxt.andptr = ctxt.andptr[1:] /* xchg lhs,bx */
			asmand_asm8(ctxt, &p.from, reg_asm8[breg])
		} else {
			ctxt.andptr[0] = uint8(0x90 + reg_asm8[z])
			ctxt.andptr = ctxt.andptr[1:] /* xchg lsh,ax */
			subreg_asm8(&pp, z, D_AX_8)
			doasm_asm8(ctxt, &pp)
			ctxt.andptr[0] = uint8(0x90 + reg_asm8[z])
			ctxt.andptr = ctxt.andptr[1:] /* xchg lsh,ax */
		}
		return
	}
	z = p.to.typ
	if z >= D_BP_8 && z <= D_DI_8 {
		breg = byteswapreg_asm8(ctxt, &p.from)
		if breg != D_AX_8 {
			ctxt.andptr[0] = 0x87
			ctxt.andptr = ctxt.andptr[1:] /* xchg rhs,bx */
			asmand_asm8(ctxt, &p.to, reg_asm8[breg])
			subreg_asm8(&pp, z, breg)
			doasm_asm8(ctxt, &pp)
			ctxt.andptr[0] = 0x87
			ctxt.andptr = ctxt.andptr[1:] /* xchg rhs,bx */
			asmand_asm8(ctxt, &p.to, reg_asm8[breg])
		} else {
			ctxt.andptr[0] = uint8(0x90 + reg_asm8[z])
			ctxt.andptr = ctxt.andptr[1:] /* xchg rsh,ax */
			subreg_asm8(&pp, z, D_AX_8)
			doasm_asm8(ctxt, &pp)
			ctxt.andptr[0] = uint8(0x90 + reg_asm8[z])
			ctxt.andptr = ctxt.andptr[1:] /* xchg rsh,ax */
		}
		return
	}
	ctxt.diag("doasm: notfound t2=%ux from=%ux to=%ux %P", t[2], p.from.typ, p.to.typ, p)
	return
mfound:
	switch t[3] {
	default:
		ctxt.diag("asmins: unknown mov %d %P", t[3], p)
	case 0: /* lit */
		for z = 4; t[z] != E_asm8; z++ {
			ctxt.andptr[0] = t[z]
			ctxt.andptr = ctxt.andptr[1:]
		}
	case 1: /* r,m */
		ctxt.andptr[0] = t[4]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, int(t[5]))
	case 2: /* m,r */
		ctxt.andptr[0] = t[4]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, int(t[5]))
	case 3: /* r,m - 2op */
		ctxt.andptr[0] = t[4]
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = t[5]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.to, int(t[6]))
	case 4: /* m,r - 2op */
		ctxt.andptr[0] = t[4]
		ctxt.andptr = ctxt.andptr[1:]
		ctxt.andptr[0] = t[5]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, int(t[6]))
	case 5: /* load full pointer, trash heap */
		if t[4] != 0 {
			ctxt.andptr[0] = t[4]
			ctxt.andptr = ctxt.andptr[1:]
		}
		switch p.to.index {
		default:
			goto bad
		case D_DS_8:
			ctxt.andptr[0] = 0xc5
			ctxt.andptr = ctxt.andptr[1:]
		case D_SS_8:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0xb2
			ctxt.andptr = ctxt.andptr[1:]
		case D_ES_8:
			ctxt.andptr[0] = 0xc4
			ctxt.andptr = ctxt.andptr[1:]
		case D_FS_8:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0xb4
			ctxt.andptr = ctxt.andptr[1:]
		case D_GS_8:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = 0xb5
			ctxt.andptr = ctxt.andptr[1:]
			break
		}
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
	case 6: /* double shift */
		z = p.from.typ
		switch z {
		default:
			goto bad
		case D_CONST_8:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = t[4]
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm8(ctxt, &p.to, reg_asm8[p.from.index])
			ctxt.andptr[0] = uint8(p.from.offset)
			ctxt.andptr = ctxt.andptr[1:]
		case D_CL_8,
			D_CX_8:
			ctxt.andptr[0] = 0x0f
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = t[5]
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm8(ctxt, &p.to, reg_asm8[p.from.index])
			break
		}
	case 7: /* imul rm,r */
		if t[4] == Pq_asm8 {
			ctxt.andptr[0] = Pe_asm8
			ctxt.andptr = ctxt.andptr[1:]
			ctxt.andptr[0] = Pm_asm8
			ctxt.andptr = ctxt.andptr[1:]
		} else {
			ctxt.andptr[0] = t[4]
			ctxt.andptr = ctxt.andptr[1:]
		}
		ctxt.andptr[0] = t[5]
		ctxt.andptr = ctxt.andptr[1:]
		asmand_asm8(ctxt, &p.from, reg_asm8[p.to.typ])
	// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
	// where you load the TLS base register into a register and then index off that
	// register to access the actual TLS variables. Systems that allow direct TLS access
	// are handled in prefixof above and should not be listed here.
	case 8: /* mov tls, r */
		switch ctxt.headtype {
		default:
			log.Fatalf("unknown TLS base location for %s", headstr(ctxt.headtype))
		// ELF TLS base is 0(GS).
		case Hlinux,
			Hnacl:
			pp.from = p.from
			pp.from.typ = D_INDIR_8 + D_GS_8
			pp.from.offset = 0
			pp.from.index = D_NONE_8
			pp.from.scale = 0
			ctxt.andptr[0] = 0x65
			ctxt.andptr = ctxt.andptr[1:] // GS
			ctxt.andptr[0] = 0x8B
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm8(ctxt, &pp.from, reg_asm8[p.to.typ])
		case Hplan9:
			if ctxt.plan9privates == nil {
				ctxt.plan9privates = linklookup(ctxt, "_privates", 0)
			}
			pp.from = Addr{}
			pp.from.typ = D_EXTERN_8
			pp.from.sym = ctxt.plan9privates
			pp.from.offset = 0
			pp.from.index = D_NONE_8
			ctxt.andptr[0] = 0x8B
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm8(ctxt, &pp.from, reg_asm8[p.to.typ])
		// Windows TLS base is always 0x14(FS).
		case Hwindows:
			pp.from = p.from
			pp.from.typ = D_INDIR_8 + D_FS_8
			pp.from.offset = 0x14
			pp.from.index = D_NONE_8
			pp.from.scale = 0
			ctxt.andptr[0] = 0x64
			ctxt.andptr = ctxt.andptr[1:] // FS
			ctxt.andptr[0] = 0x8B
			ctxt.andptr = ctxt.andptr[1:]
			asmand_asm8(ctxt, &pp.from, reg_asm8[p.to.typ])
			break
		}
		break
	}
}

var naclret_asm8 = []uint8{
	0x5d, // POPL BP
	// 0x8b, 0x7d, 0x00, // MOVL (BP), DI - catch return to invalid address, for debugging
	0x83,
	0xe5,
	0xe0, // ANDL $~31, BP
	0xff,
	0xe5, // JMP BP
}

func asmins_asm8(ctxt *Link, p *Prog) {
	var r *Reloc
	ctxt.andptr = ctxt.and[:]
	if p.as == AUSEFIELD_8 {
		r = addrel(ctxt.cursym)
		r.off = 0
		r.sym = p.from.sym
		r.typ = R_USEFIELD
		r.siz = 0
		return
	}
	if ctxt.headtype == Hnacl {
		switch p.as {
		case ARET_8:
			copy(ctxt.andptr, naclret_asm8)
			ctxt.andptr = ctxt.andptr[len(naclret_asm8):]
			return
		case ACALL_8,
			AJMP_8:
			if D_AX_8 <= p.to.typ && p.to.typ <= D_DI_8 {
				ctxt.andptr[0] = 0x83
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = uint8(0xe0 | (p.to.typ - D_AX_8))
				ctxt.andptr = ctxt.andptr[1:]
				ctxt.andptr[0] = 0xe0
				ctxt.andptr = ctxt.andptr[1:]
			}
		case AINT_8:
			ctxt.andptr[0] = 0xf4
			ctxt.andptr = ctxt.andptr[1:]
			return
		}
	}
	doasm_asm8(ctxt, p)
	if -cap(ctxt.andptr) > -cap(ctxt.and[len(ctxt.and):]) {
		fmt.Printf("and[] is too short - %d byte instruction\n", -cap(ctxt.andptr)+cap(ctxt.and[:]))
		log.Fatalf("bad code")
	}
}
