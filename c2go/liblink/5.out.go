package main

// Inferno utils/5l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/5l/span.c
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
// Instruction layout.
// Inferno utils/5c/5.out.h
// http://code.google.com/p/inferno-os/source/browse/utils/5c/5.out.h
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
	NSNAME_5 = 8
	NSYM_5   = 50
	NREG_5   = 16
)

/* -1 disables use of REGARG */
const (
	REGARG_5 = -1
)

const (
	REGRET_5  = 0
	REGEXT_5  = 10
	REGG_5    = REGEXT_5 - 0
	REGM_5    = REGEXT_5 - 1
	REGTMP_5  = 11
	REGSP_5   = 13
	REGLINK_5 = 14
	REGPC_5   = 15
	NFREG_5   = 16
	FREGRET_5 = 0
	FREGEXT_5 = 7
	FREGTMP_5 = 15
)

/* compiler allocates register variables F0 up */
/* compiler allocates external registers F7 down */
const (
	AXXX_5 = iota
	AAND_5
	AEOR_5
	ASUB_5
	ARSB_5
	AADD_5
	AADC_5
	ASBC_5
	ARSC_5
	ATST_5
	ATEQ_5
	ACMP_5
	ACMN_5
	AORR_5
	ABIC_5
	AMVN_5
	AB_5
	ABL_5
	ABEQ_5
	ABNE_5
	ABCS_5
	ABHS_5
	ABCC_5
	ABLO_5
	ABMI_5
	ABPL_5
	ABVS_5
	ABVC_5
	ABHI_5
	ABLS_5
	ABGE_5
	ABLT_5
	ABGT_5
	ABLE_5
	AMOVWD_5
	AMOVWF_5
	AMOVDW_5
	AMOVFW_5
	AMOVFD_5
	AMOVDF_5
	AMOVF_5
	AMOVD_5
	ACMPF_5
	ACMPD_5
	AADDF_5
	AADDD_5
	ASUBF_5
	ASUBD_5
	AMULF_5
	AMULD_5
	ADIVF_5
	ADIVD_5
	ASQRTF_5
	ASQRTD_5
	AABSF_5
	AABSD_5
	ASRL_5
	ASRA_5
	ASLL_5
	AMULU_5
	ADIVU_5
	AMUL_5
	ADIV_5
	AMOD_5
	AMODU_5
	AMOVB_5
	AMOVBS_5
	AMOVBU_5
	AMOVH_5
	AMOVHS_5
	AMOVHU_5
	AMOVW_5
	AMOVM_5
	ASWPBU_5
	ASWPW_5
	ANOP_5
	ARFE_5
	ASWI_5
	AMULA_5
	ADATA_5
	AGLOBL_5
	AGOK_5
	AHISTORY_5
	ANAME_5
	ARET_5
	ATEXT_5
	AWORD_5
	ADYNT__5
	AINIT__5
	ABCASE_5
	ACASE_5
	AEND_5
	AMULL_5
	AMULAL_5
	AMULLU_5
	AMULALU_5
	ABX_5
	ABXRET_5
	ADWORD_5
	ASIGNAME_5
	ALDREX_5
	ASTREX_5
	ALDREXD_5
	ASTREXD_5
	APLD_5
	AUNDEF_5
	ACLZ_5
	AMULWT_5
	AMULWB_5
	AMULAWT_5
	AMULAWB_5
	AUSEFIELD_5
	ATYPE_5
	AFUNCDATA_5
	APCDATA_5
	ACHECKNIL_5
	AVARDEF_5
	AVARKILL_5
	ADUFFCOPY_5
	ADUFFZERO_5
	ADATABUNDLE_5
	ADATABUNDLEEND_5
	AMRC_5
	ALAST_5
)

/* scond byte */
const (
	C_SCOND_5      = (1 << 4) - 1
	C_SBIT_5       = 1 << 4
	C_PBIT_5       = 1 << 5
	C_WBIT_5       = 1 << 6
	C_FBIT_5       = 1 << 7
	C_UBIT_5       = 1 << 7
	C_SCOND_EQ_5   = 0
	C_SCOND_NE_5   = 1
	C_SCOND_HS_5   = 2
	C_SCOND_LO_5   = 3
	C_SCOND_MI_5   = 4
	C_SCOND_PL_5   = 5
	C_SCOND_VS_5   = 6
	C_SCOND_VC_5   = 7
	C_SCOND_HI_5   = 8
	C_SCOND_LS_5   = 9
	C_SCOND_GE_5   = 10
	C_SCOND_LT_5   = 11
	C_SCOND_GT_5   = 12
	C_SCOND_LE_5   = 13
	C_SCOND_NONE_5 = 14
	C_SCOND_NV_5   = 15
	SHIFT_LL_5     = 0 << 5
	SHIFT_LR_5     = 1 << 5
	SHIFT_AR_5     = 2 << 5
	SHIFT_RR_5     = 3 << 5
)

const (
	D_GOK_5     = 0
	D_NONE_5    = 1
	D_BRANCH_5  = D_NONE_5 + 1
	D_OREG_5    = D_NONE_5 + 2
	D_CONST_5   = D_NONE_5 + 7
	D_FCONST_5  = D_NONE_5 + 8
	D_SCONST_5  = D_NONE_5 + 9
	D_PSR_5     = D_NONE_5 + 10
	D_REG_5     = D_NONE_5 + 12
	D_FREG_5    = D_NONE_5 + 13
	D_FILE_5    = D_NONE_5 + 16
	D_OCONST_5  = D_NONE_5 + 17
	D_FILE1_5   = D_NONE_5 + 18
	D_SHIFT_5   = D_NONE_5 + 19
	D_FPCR_5    = D_NONE_5 + 20
	D_REGREG_5  = D_NONE_5 + 21
	D_ADDR_5    = D_NONE_5 + 22
	D_SBIG_5    = D_NONE_5 + 23
	D_CONST2_5  = D_NONE_5 + 24
	D_REGREG2_5 = D_NONE_5 + 25
	D_EXTERN_5  = D_NONE_5 + 3
	D_STATIC_5  = D_NONE_5 + 4
	D_AUTO_5    = D_NONE_5 + 5
	D_PARAM_5   = D_NONE_5 + 6
)

/*
 * this is the ranlib header
 */
var SYMDEF []int8
