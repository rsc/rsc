package main

import (
	"bufio"
	"go/build"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var arch *LinkArch

func main() {
	switch build.Default.GOARCH {
	case "amd64":
		arch = &linkamd64
	case "amd64p32":
		arch = &linkamd64p32
	case "386":
		arch = &link386
	case "arm":
		arch = &linkarm
	}
	if len(os.Args) == 3 {
		input()
		return
	}
	f, err := os.Create("x.6")
	if err != nil {
		log.Fatal(err)
	}
	ctxt := linknew(arch)
	ctxt.debugasm = 1
	ctxt.bso = Binitw(os.Stderr)
	defer Bflush(ctxt.bso)
	ctxt.diag = log.Fatalf
	obuf := Binitw(f)
	Bprint(obuf, "go object %s %s %s\n", getgoos(), getgoarch(), getgoversion())
	Bprint(obuf, "!\n")
	p1 := &Prog{
		ctxt:   ctxt,
		as:     ATEXT_6,
		lineno: 1,
		from: Addr{
			typ:   D_EXTERN_6,
			index: D_NONE_6,
			sym:   linklookup(ctxt, "main.main", 0),
			scale: 0,
		},
		to: Addr{
			typ:   D_CONST_6,
			index: D_NONE_6,
		},
	}
	p2 := &Prog{
		ctxt: ctxt,
		as:   ARET_6,
		from: Addr{
			typ:   D_NONE_6,
			index: D_NONE_6,
		},
		to: Addr{
			typ:   D_NONE_6,
			index: D_NONE_6,
		},
	}
	p3 := &Prog{
		ctxt:   ctxt,
		as:     ATEXT_6,
		lineno: 1,
		from: Addr{
			typ:   D_EXTERN_6,
			index: D_NONE_6,
			sym:   linklookup(ctxt, "main.init", 0),
			scale: 0,
		},
		to: Addr{
			typ:   D_CONST_6,
			index: D_NONE_6,
		},
	}
	p4 := &Prog{
		ctxt: ctxt,
		as:   ARET_6,
		from: Addr{
			typ:   D_NONE_6,
			index: D_NONE_6,
		},
		to: Addr{
			typ:   D_NONE_6,
			index: D_NONE_6,
		},
	}
	pl := linknewplist(ctxt)
	pl.firstpc = p1
	p1.link = p2
	p2.link = p3
	p3.link = p4
	writeobj(ctxt, obuf)
	Bflush(obuf)
}

var (
	ctxt   *Link
	plists = map[string]*Plist{}
	syms   = map[string]*LSym{}
	progs  = map[string]*Prog{}
	hists  = map[string]*Hist{}
	undef  = map[interface{}]bool{}
)

func input() {
	ctxt = linknew(arch)
	//ctxt.debugasm = 1
	ctxt.bso = Binitw(os.Stderr)
	defer Bflush(ctxt.bso)
	ctxt.diag = log.Fatalf
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	b := bufio.NewReader(f)
	if rdstring(b) != "ctxt" {
		log.Fatal("invalid input - missing ctxt")
	}
	name := rdstring(b)
	if name != ctxt.arch.name {
		log.Fatalf("bad arch %s - want %s", name, ctxt.arch.name)
	}

	ctxt.debugasm = int32(rdint(b))
	ctxt.trimpath = rdstring(b)
	ctxt.plist = rdplist(b)
	ctxt.plast = rdplist(b)
	ctxt.hist = rdhist(b)
	ctxt.ehist = rdhist(b)
	for {
		i := rdint(b)
		if i < 0 {
			break
		}
		ctxt.hash[i] = rdsym(b)
	}
	last := "ctxt"

Loop:
	for {
		s := rdstring(b)
		switch s {
		default:
			log.Fatalf("unexpected input after %s: %v", s, last)
		case "end":
			break Loop
		case "plist":
			readplist(b, rdplist(b))
		case "sym":
			readsym(b, rdsym(b))
		case "prog":
			readprog(b, rdprog(b))
		case "hist":
			readhist(b, rdhist(b))
		}
		last = s
	}

	if len(undef) > 0 {
		panic("missing definitions")
	}

	ff, err := os.Create(os.Args[2])
	obuf := Binitw(ff)
	writeobj(ctxt, obuf)
	Bflush(obuf)
}

func rdstring(b *bufio.Reader) string {
	s, err := b.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s = strings.TrimSpace(s)
	if s == "<nil>" {
		s = ""
	}
	return s
}

func rdplist(b *bufio.Reader) *Plist {
	id := rdstring(b)
	if id == "0" || id == "0x0" {
		return nil
	}
	pl := plists[id]
	if pl == nil {
		pl = new(Plist)
		plists[id] = pl
		undef[pl] = true
	}
	return pl
}

func rdsym(b *bufio.Reader) *LSym {
	id := rdstring(b)
	if id == "0" || id == "0x0" {
		return nil
	}
	sym := syms[id]
	if sym == nil {
		sym = new(LSym)
		syms[id] = sym
		undef[sym] = true
	}
	return sym
}

func rdprog(b *bufio.Reader) *Prog {
	id := rdstring(b)
	if id == "0" || id == "0x0" {
		return nil
	}
	prog := progs[id]
	if prog == nil {
		prog = new(Prog)
		prog.ctxt = ctxt
		progs[id] = prog
		undef[prog] = true
	}
	return prog
}

func rdhist(b *bufio.Reader) *Hist {
	id := rdstring(b)
	if id == "0" || id == "0x0" {
		return nil
	}
	h := hists[id]
	if h == nil {
		h = new(Hist)
		hists[id] = h
		undef[h] = true
	}
	return h
}

func readplist(b *bufio.Reader, pl *Plist) {
	if !undef[pl] {
		panic("double-def")
	}
	delete(undef, pl)
	pl.recur = int(rdint(b))
	pl.name = rdsym(b)
	pl.firstpc = rdprog(b)
	pl.link = rdplist(b)
}

func readsym(b *bufio.Reader, s *LSym) {
	if !undef[s] {
		panic("double-def")
	}
	delete(undef, s)
	s.name = rdstring(b)
	s.extname = rdstring(b)
	s.typ = rdint(b)
	s.version = rdint(b)
	s.dupok = rdint(b)
	s.external = uint8(rdint(b))
	s.nosplit = rdint(b)
	s.reachable = uint8(rdint(b))
	s.cgoexport = uint8(rdint(b))
	s.special = uint8(rdint(b))
	s.stkcheck = uint8(rdint(b))
	s.hide = uint8(rdint(b))
	s.leaf = rdint(b)
	s.fnptr = uint8(rdint(b))
	s.seenglobl = uint8(rdint(b))
	s.onlist = uint8(rdint(b))
	s.symid = int16(rdint(b))
	s.dynid = int32(rdint(b))
	s.sig = int32(rdint(b))
	s.plt = int32(rdint(b))
	s.got = int32(rdint(b))
	s.align = int32(rdint(b))
	s.elfsym = int32(rdint(b))
	s.args = rdint(b)
	s.locals = rdint(b)
	s.value = rdint(b)
	s.size = rdint(b)
	s.hash = rdsym(b)
	s.allsym = rdsym(b)
	s.next = rdsym(b)
	s.sub = rdsym(b)
	s.outer = rdsym(b)
	s.gotype = rdsym(b)
	s.reachparent = rdsym(b)
	s.queue = rdsym(b)
	s.file = rdstring(b)
	s.dynimplib = rdstring(b)
	s.dynimpvers = rdstring(b)
	s.text = rdprog(b)
	s.etext = rdprog(b)
	n := int(rdint(b))
	if n > 0 {
		s.p = make([]byte, n)
		io.ReadFull(b, s.p)
	}
	s.r = make([]Reloc, int(rdint(b)))
	for i := range s.r {
		r := &s.r[i]
		r.off = rdint(b)
		r.siz = rdint(b)
		r.done = uint8(rdint(b))
		r.typ = rdint(b)
		r.add = rdint(b)
		r.xadd = rdint(b)
		r.sym = rdsym(b)
		r.xsym = rdsym(b)
	}
}

func readprog(b *bufio.Reader, p *Prog) {
	if !undef[p] {
		panic("double-def")
	}
	delete(undef, p)
	p.pc = rdint(b)
	p.lineno = int32(rdint(b))
	p.link = rdprog(b)
	p.as = int(rdint(b))
	p.reg = rdint(b)
	p.scond = int(rdint(b))
	p.width = int8(rdint(b))
	readaddr(b, &p.from)
	readaddr(b, &p.to)
}

func readaddr(b *bufio.Reader, a *Addr) {
	if rdstring(b) != "addr" {
		log.Fatal("out of sync")
	}
	a.offset = rdint(b)
	a.u.dval = rdfloat(b)
	buf := make([]byte, 8)
	for i := 0; i < 8; i++ {
		buf[i] = byte(rdint(b))
	}
	a.u.sval = string(buf)
	a.u.branch = rdprog(b)
	a.sym = rdsym(b)
	a.gotype = rdsym(b)
	a.typ = rdint(b)
	a.index = rdint(b)
	a.scale = rdint(b)
	a.reg = rdint(b)
	a.name = rdint(b)
	a.class = int(rdint(b))
	a.etype = uint8(rdint(b))
	a.offset2 = rdint(b)
	a.width = rdint(b)
}

func readhist(b *bufio.Reader, h *Hist) {
	if !undef[h] {
		panic("double-def")
	}
	delete(undef, h)
	h.link = rdhist(b)
	h.name = rdstring(b)
	h.line = int32(rdint(b))
	h.offset = int32(rdint(b))
}

func rdint(b *bufio.Reader) int64 {
	s := rdstring(b)
	x, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		panic(err)
	}
	return x
}

func rdfloat(b *bufio.Reader) float64 {
	s := rdstring(b)
	x, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		panic(err)
	}
	return math.Float64frombits(x)
}
