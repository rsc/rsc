package main

import (
	"encoding/binary"
	"fmt"
)

// Derived from Inferno utils/6l/l.h and related files.
// http://code.google.com/p/inferno-os/source/browse/utils/6l/l.h
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
type Addr struct {
	offset int64
	u      struct {
		sval   string
		dval   float64
		branch *Prog
	}
	sym     *LSym
	gotype  *LSym
	typ     int64
	index   int64
	scale   int64
	reg     int64
	name    int64
	class   int
	etype   uint8
	offset2 int64
	node    *struct{}
	width   int64
}

type Prog struct {
	ctxt     *Link
	pc       int64
	lineno   int32
	link     *Prog
	as       int
	reg      int64
	scond    int
	from     Addr
	to       Addr
	opt      *struct{}
	forwd    *Prog
	pcond    *Prog
	comefrom *Prog
	pcrel    *Prog
	spadj    int64
	mark     int64
	back     int
	ft       uint8
	tt       uint8
	optab    int
	isize    int64
	width    int8
	mode     int64
	TEXTFLAG uint8
}

func (p *Prog) Line() string {
	return linklinefmt(p.ctxt, int(p.lineno), false, false)
}

type LSym struct {
	name        string
	extname     string
	typ         int64
	version     int64
	dupok       int64
	external    uint8
	nosplit     int64
	reachable   uint8
	cgoexport   uint8
	special     uint8
	stkcheck    uint8
	hide        uint8
	leaf        int64
	fnptr       uint8
	seenglobl   uint8
	onlist      uint8
	symid       int16
	dynid       int32
	sig         int32
	plt         int32
	got         int32
	align       int32
	elfsym      int32
	args        int64
	locals      int64
	value       int64
	size        int64
	hash        *LSym
	allsym      *LSym
	next        *LSym
	sub         *LSym
	outer       *LSym
	gotype      *LSym
	reachparent *LSym
	queue       *LSym
	file        string
	dynimplib   string
	dynimpvers  string
	sect        *struct{}
	autom       *Auto
	text        *Prog
	etext       *Prog
	pcln        *Pcln
	p           []uint8
	r           []Reloc
}

type Reloc struct {
	off  int64
	siz  int64
	done uint8
	typ  int64
	add  int64
	xadd int64
	sym  *LSym
	xsym *LSym
}

type Auto struct {
	asym    *LSym
	link    *Auto
	aoffset int64
	typ     int64
	gotype  *LSym
}

type Hist struct {
	link   *Hist
	name   string
	line   int32
	offset int32
}

type Link struct {
	thechar        int32
	thestring      string
	goarm          int32
	headtype       int
	arch           *LinkArch
	ignore         func(string) int32
	debugasm       int32
	debugline      int32
	debughist      int32
	debugread      int32
	debugvlog      int32
	debugstack     int32
	debugzerostack int32
	debugdivmod    int32
	debugfloat     int32
	debugpcln      int32
	flag_shared    int32
	iself          int32
	bso            *Biobuf
	pathname       string
	windows        int32
	trimpath       string
	goroot         string
	goroot_final   string
	hash           [LINKHASH]*LSym
	allsym         *LSym
	nsymbol        int32
	hist           *Hist
	ehist          *Hist
	plist          *Plist
	plast          *Plist
	sym_div        *LSym
	sym_divu       *LSym
	sym_mod        *LSym
	sym_modu       *LSym
	symmorestack   [20]*LSym
	tlsg           *LSym
	plan9privates  *LSym
	curp           *Prog
	printp         *Prog
	blitrl         *Prog
	elitrl         *Prog
	rexflag        int
	rep            int
	repn           int
	lock           int
	asmode         int64
	andptr         []uint8
	and            [100]uint8
	instoffset     int32
	autosize       int32
	armsize        int32
	pc             int64
	libdir         []string
	library        []Library
	tlsoffset      int
	diag           func(string, ...interface{})
	mode           int64
	curauto        *Auto
	curhist        *Auto
	cursym         *LSym
	version        int64
	textp          *LSym
	etextp         *LSym
	histdepth      int32
	nhistfile      int32
	filesyms       *LSym
}

type Plist struct {
	name    *LSym
	firstpc *Prog
	recur   int
	link    *Plist
}

type LinkArch struct {
	name          string
	thechar       int
	addstacksplit func(*Link, *LSym)
	assemble      func(*Link, *LSym)
	datasize      func(*Prog) int
	follow        func(*Link, *LSym)
	iscall        func(*Prog) int
	isdata        func(*Prog) int
	prg           func() *Prog
	progedit      func(*Link, *Prog)
	settextflag   func(*Prog, int)
	symtype       func(*Addr) int
	textflag      func(*Prog) int
	minlc         uint32
	ptrsize       int64
	regsize       int
	byteOrder     binary.ByteOrder
	Pconv         func(*Prog) string
	D_ADDR        int64
	D_AUTO        int64
	D_BRANCH      int64
	D_CONST       int64
	D_EXTERN      int64
	D_FCONST      int64
	D_NONE        int64
	D_PARAM       int64
	D_SCONST      int64
	D_STATIC      int64
	ACALL         int
	ADATA         int
	AEND          int
	AFUNCDATA     int
	AGLOBL        int
	AJMP          int
	ANOP          int
	APCDATA       int
	ARET          int
	ATEXT         int
	ATYPE         int
	AUSEFIELD     int
}

type Library struct {
	objref string
	srcref string
	file   string
	pkg    string
}

type Pcln struct {
	pcsp        Pcdata
	pcfile      Pcdata
	pcline      Pcdata
	pcdata      []Pcdata
	npcdata     int64
	funcdata    []*LSym
	funcdataoff []int64
	nfuncdata   int64
	file        []*LSym
	nfile       int64
	mfile       int64
	lastfile    *LSym
	lastindex   int64
}

type Pcdata struct {
	p []uint8
	n int64
	m int64
}

type Pciter struct {
	d       Pcdata
	p       []uint8
	pc      uint32
	nextpc  uint32
	pcscale uint32
	value   int32
	start   int
	done    int
}

// prevent incompatible type signatures between liblink and 8l on Plan 9

// prevent incompatible type signatures between liblink and 8l on Plan 9

// LSym.type
const (
	Sxxx = iota
	STEXT
	SELFRXSECT
	STYPE
	SSTRING
	SGOSTRING
	SGOFUNC
	SRODATA
	SFUNCTAB
	STYPELINK
	SSYMTAB
	SPCLNTAB
	SELFROSECT
	SMACHOPLT
	SELFSECT
	SMACHO
	SMACHOGOT
	SNOPTRDATA
	SINITARR
	SDATA
	SWINDOWS
	SBSS
	SNOPTRBSS
	STLSBSS
	SXREF
	SMACHOSYMSTR
	SMACHOSYMTAB
	SMACHOINDIRECTPLT
	SMACHOINDIRECTGOT
	SFILE
	SFILEPATH
	SCONST
	SDYNIMPORT
	SHOSTOBJ
	SSUB    = 1 << 8
	SMASK   = SSUB - 1
	SHIDDEN = 1 << 9
)

// Reloc.type
const (
	R_ADDR = 1 + iota
	R_SIZE
	R_CALL
	R_CALLARM
	R_CALLIND
	R_CONST
	R_PCREL
	R_TLS
	R_TLS_LE
	R_TLS_IE
	R_GOTOFF
	R_PLT0
	R_PLT1
	R_PLT2
	R_USEFIELD
)

// Auto.type
const (
	A_AUTO = 1 + iota
	A_PARAM
)

const (
	LINKHASH = 100003
)

// Pcdata iterator.
//	for(pciterinit(ctxt, &it, &pcd); !it.done; pciternext(&it)) { it.value holds in [it.pc, it.nextpc) }

// symbol version, incremented each time a file is loaded.
// version==1 is reserved for savehist.
const (
	HistVersion = 1
)

// Link holds the context for writing object code from a compiler
// to be linker input or for reading that input into the linker.

// LinkArch is the definition of a single architecture.

/* executable header types */
const (
	Hunknown = 0 + iota
	Hdarwin
	Hdragonfly
	Helf
	Hfreebsd
	Hlinux
	Hnacl
	Hnetbsd
	Hopenbsd
	Hplan9
	Hsolaris
	Hwindows
)

const (
	LinkAuto = 0 + iota
	LinkInternal
	LinkExternal
)

// asm5.c
// asm6.c
// asm8.c
// data.c
// go.c
// ld.c

const (
	fmtLong = 1 << iota
)

const (
	NOPROF_textflag   = 1
	DUPOK_textflag    = 2
	NOSPLIT_textflag  = 4
	RODATA_textflag   = 8
	NOPTR_textflag    = 16
	WRAPPER_textflag  = 32
	NEEDCTXT_textflag = 64
)

func (ctxt *Link) prg() *Prog {
	p := ctxt.arch.prg()
	p.ctxt = ctxt
	return p
}

func (p *Prog) String() string {
	if p.ctxt == nil {
		return fmt.Sprintf("PROG MISSING CTXT")
	}
	return p.ctxt.arch.Pconv(p)
}
