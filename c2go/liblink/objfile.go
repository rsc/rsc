package main

import (
	"bytes"
	"strconv"
	"strings"
)

func ldobjfile(ctxt *Link, f *Biobuf, pkg string, len int64, pn string) {
	var c int
	var buf [8]uint8
	var start int64
	var lib string
	start = Boffset(f)
	ctxt.version++
	buf = [8]uint8{}
	Bread(f, buf[:])
	if string(buf[:]) != startmagic_objfile {
		sysfatal("%s: invalid file start %x %x %x %x %x %x %x %x", pn, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7])
	}
	c = Bgetc(f)
	if (c) != 1 {
		sysfatal("%s: invalid file version number %d", pn, c)
	}
	for {
		lib = rdstring_objfile(f)
		if lib[0] == 0 {
			break
		}
		addlib(ctxt, pkg, pn, lib)
	}
	for {
		c = Bgetc(f)
		Bungetc(f)
		if c == 0xff {
			break
		}
		readsym_objfile(ctxt, f, pkg, pn)
	}
	buf = [8]uint8{}
	Bread(f, buf[:])
	if string(buf[:]) != endmagic_objfile {
		sysfatal("%s: invalid file end", pn)
	}
	if Boffset(f) != start+len {
		sysfatal("%s: unexpected end at %lld, want %lld", pn, int64(Boffset(f)), int64(start+len))
	}
}

// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Writing and reading of Go object files.
//
// Originally, Go object files were Plan 9 object files, but no longer.
// Now they are more like standard object files, in that each symbol is defined
// by an associated memory image (bytes) and a list of relocations to apply
// during linking. We do not (yet?) use a standard file format, however.
// For now, the format is chosen to be as simple as possible to read and write.
// It may change for reasons of efficiency, or we may even switch to a
// standard file format if there are compelling benefits to doing so.
// See golang.org/s/go13linker for more background.
//
// The file format is:
//
//	- magic header: "\x00\x00go13ld"
//	- byte 1 - version number
//	- sequence of strings giving dependencies (imported packages)
//	- empty string (marks end of sequence)
//	- sequence of defined symbols
//	- byte 0xff (marks end of sequence)
//	- magic footer: "\xff\xffgo13ld"
//
// All integers are stored in a zigzag varint format.
// See golang.org/s/go12symtab for a definition.
//
// Data blocks and strings are both stored as an integer
// followed by that many bytes.
//
// A symbol reference is a string name followed by a version.
// An empty name corresponds to a nil LSym* pointer.
//
// Each symbol is laid out as the following fields (taken from LSym*):
//
//	- byte 0xfe (sanity check for synchronization)
//	- type [int]
//	- name [string]
//	- version [int]
//	- dupok [int]
//	- size [int]
//	- gotype [symbol reference]
//	- p [data block]
//	- nr [int]
//	- r [nr relocations, sorted by off]
//
// If type == STEXT, there are a few more fields:
//
//	- args [int]
//	- locals [int]
//	- nosplit [int]
//	- leaf [int]
//	- nlocal [int]
//	- local [nlocal automatics]
//	- pcln [pcln table]
//
// Each relocation has the encoding:
//
//	- off [int]
//	- siz [int]
//	- type [int]
//	- add [int]
//	- xadd [int]
//	- sym [symbol reference]
//	- xsym [symbol reference]
//
// Each local has the encoding:
//
//	- asym [symbol reference]
//	- offset [int]
//	- type [int]
//	- gotype [symbol reference]
//
// The pcln table has the encoding:
//
//	- pcsp [data block]
//	- pcfile [data block]
//	- pcline [data block]
//	- npcdata [int]
//	- pcdata [npcdata data blocks]
//	- nfuncdata [int]
//	- funcdata [nfuncdata symbol references]
//	- funcdatasym [nfuncdata ints]
//	- nfile [int]
//	- file [nfile symbol references]
//
// The file layout and meaning of type integers are architecture-independent.
//
// TODO(rsc): The file format is good for a first pass but needs work.
//	- There are SymID in the object file that should really just be strings.
//	- The actual symbol memory images are interlaced with the symbol
//	  metadata. They should be separated, to reduce the I/O required to
//	  load just the metadata.
//	- The symbol references should be shortened, either with a symbol
//	  table or by using a simple backward index to an earlier mentioned symbol.
// The Go and C compilers, and the assembler, call writeobj to write
// out a Go object file.  The linker does not call this; the linker
// does not write out object files.
func writeobj(ctxt *Link, b *Biobuf) {
	var flag int
	var h *Hist
	var s *LSym
	var text *LSym
	var etext *LSym
	var curtext *LSym
	var data *LSym
	var edata *LSym
	var pl *Plist
	var p *Prog
	var plink *Prog
	var a *Auto
	// Build list of symbols, and assign instructions to lists.
	// Ignore ctxt->plist boundaries. There are no guarantees there,
	// and the C compilers and assemblers just use one big list.
	text = (*LSym)(nil)
	curtext = (*LSym)(nil)
	data = (*LSym)(nil)
	etext = (*LSym)(nil)
	edata = (*LSym)(nil)
	for pl = ctxt.plist; pl != nil; pl = pl.link {
		println("writeobj", pl)
		for p = pl.firstpc; p != nil; p = plink {
			plink = p.link
			p.link = (*Prog)(nil)
			if p.as == ctxt.arch.AEND {
				continue
			}
			if p.as == ctxt.arch.ATYPE {
				// Assume each TYPE instruction describes
				// a different local variable or parameter,
				// so no dedup.
				// Using only the TYPE instructions means
				// that we discard location information about local variables
				// in C and assembly functions; that information is inferred
				// from ordinary references, because there are no TYPE
				// instructions there. Without the type information, gdb can't
				// use the locations, so we don't bother to save them.
				// If something else could use them, we could arrange to
				// preserve them.
				if curtext == nil {
					continue
				}
				a = new(Auto)
				a.asym = p.from.sym
				a.aoffset = int32(p.from.offset)
				a.typ = ctxt.arch.symtype(&p.from)
				a.gotype = p.from.gotype
				a.link = curtext.autom
				curtext.autom = a
				continue
			}
			if p.as == ctxt.arch.AGLOBL {
				s = p.from.sym
				var tmp uint8 = s.seenglobl
				s.seenglobl++
				if tmp != 0 {
					print("duplicate %P\n", p)
				}
				if s.onlist != 0 {
					sysfatal("symbol %s listed multiple times", s.name)
				}
				s.onlist = 1
				if data == nil {
					data = s
				} else {
					edata.next = s
				}
				s.next = (*LSym)(nil)
				s.size = int(p.to.offset)
				if s.typ == 0 || s.typ == int(SXREF) {
					s.typ = int(SBSS)
				}
				if ctxt.arch.thechar == '5' {
					flag = p.reg
				} else {
					flag = p.from.scale
				}
				if flag&int(DUPOK_textflag) != 0 {
					s.dupok = 1
				}
				if flag&int(RODATA_textflag) != 0 {
					s.typ = int(SRODATA)
				} else {
					if flag&int(NOPTR_textflag) != 0 {
						s.typ = int(SNOPTRBSS)
					}
				}
				edata = s
				continue
			}
			if p.as == ctxt.arch.ADATA {
				savedata(ctxt, p.from.sym, p, "<input>")
				continue
			}
			if p.as == ctxt.arch.ATEXT {
				s = p.from.sym
				if s == nil {
					// func _() { }
					curtext = (*LSym)(nil)
					continue
				}
				if s.text != nil {
					sysfatal("duplicate TEXT for %s", s.name)
				}
				if s.onlist != 0 {
					sysfatal("symbol %s listed multiple times", s.name)
				}
				s.onlist = 1
				if text == nil {
					text = s
				} else {
					etext.next = s
				}
				etext = s
				if ctxt.arch.thechar == '5' {
					flag = p.reg
				} else {
					flag = p.from.scale
				}
				if flag&int(DUPOK_textflag) != 0 {
					s.dupok = 1
				}
				if flag&int(NOSPLIT_textflag) != 0 {
					s.nosplit = 1
				}
				s.next = (*LSym)(nil)
				s.typ = int(STEXT)
				s.text = p
				s.etext = p
				curtext = s
				continue
			}
			if curtext == nil {
				continue
			}
			s = curtext
			s.etext.link = p
			s.etext = p
		}
	}
	// Turn functions into machine code images.
	for s = text; s != nil; s = s.next {
		mkfwd(s)
		linkpatch(ctxt, s)
		ctxt.arch.follow(ctxt, s)
		ctxt.arch.addstacksplit(ctxt, s)
		ctxt.arch.assemble(ctxt, s)
		linkpcln(ctxt, s)
	}
	// Emit header.
	Bputc(b, 0)
	Bputc(b, 0)
	Bprint(b, "go13ld")
	Bputc(b, 1) // version
	// Emit autolib.
	for h = ctxt.hist; h != nil; h = h.link {
		if h.offset < 0 {
			wrstring_objfile(b, h.name)
		}
	}
	wrstring_objfile(b, "")
	// Emit symbols.
	for s = text; s != nil; s = s.next {
		writesym_objfile(ctxt, b, s)
	}
	for s = data; s != nil; s = s.next {
		writesym_objfile(ctxt, b, s)
	}
	// Emit footer.
	Bputc(b, 0xff)
	Bputc(b, 0xff)
	Bprint(b, "go13ld")
}

func writesym_objfile(ctxt *Link, b *Biobuf, s *LSym) {
	var r *Reloc
	var i int
	var j int
	var c int
	var n int
	var pc *Pcln
	var p *Prog
	var a *Auto
	var name string
	if ctxt.debugasm != 0 {
		Bprint(ctxt.bso, "%s ", s.name)
		if s.version != 0 {
			Bprint(ctxt.bso, "v=%d ", s.version)
		}
		if s.typ != 0 {
			Bprint(ctxt.bso, "t=%d ", s.typ)
		}
		if s.dupok != 0 {
			Bprint(ctxt.bso, "dupok ")
		}
		if s.nosplit != 0 {
			Bprint(ctxt.bso, "nosplit ")
		}
		Bprint(ctxt.bso, "size=%lld value=%lld", int64(s.size), int64(s.value))
		if s.typ == int(STEXT) {
			Bprint(ctxt.bso, " args=%#llux locals=%#llux", uint64(s.args), uint64(s.locals))
			if s.leaf != 0 {
				Bprint(ctxt.bso, " leaf")
			}
		}
		Bprint(ctxt.bso, "\n")
		for p = s.text; p != nil; p = p.link {
			Bprint(ctxt.bso, "\t%#06ux %P\n", int(p.pc), p)
		}
		for i = 0; i < len(s.p); {
			Bprint(ctxt.bso, "\t%#06ux", i)
			for j = i; j < i+16 && j < len(s.p); j++ {
				Bprint(ctxt.bso, " %02ux", s.p[j])
			}
			for ; j < i+16; j++ {
				Bprint(ctxt.bso, "   ")
			}
			Bprint(ctxt.bso, "  ")
			for j = i; j < i+16 && j < len(s.p); j++ {
				c = int(s.p[j])
				if ' ' <= c && c <= 0x7e {
					Bprint(ctxt.bso, "%c", c)
				} else {
					Bprint(ctxt.bso, ".")
				}
			}
			Bprint(ctxt.bso, "\n")
			i += 16
		}
		for i = range s.r {
			r = &s.r[i]
			name = ""
			if r.sym != nil {
				name = r.sym.name
			}
			Bprint(ctxt.bso, "\trel %d+%d t=%d %s+%lld\n", int(r.off), r.siz, r.typ, name, int64(r.add))
		}
	}
	Bputc(b, 0xfe)
	wrint_objfile(b, int64(s.typ))
	wrstring_objfile(b, s.name)
	wrint_objfile(b, int64(s.version))
	wrint_objfile(b, int64(s.dupok))
	wrint_objfile(b, int64(s.size))
	wrsym_objfile(b, s.gotype)
	wrdata_objfile(b, s.p)
	wrint_objfile(b, int64(len(s.r)))
	for i = range s.r {
		r = &s.r[i]
		wrint_objfile(b, int64(r.off))
		wrint_objfile(b, int64(r.siz))
		wrint_objfile(b, int64(r.typ))
		wrint_objfile(b, r.add)
		wrint_objfile(b, r.xadd)
		wrsym_objfile(b, r.sym)
		wrsym_objfile(b, r.xsym)
	}
	if s.typ == int(STEXT) {
		wrint_objfile(b, int64(s.args))
		wrint_objfile(b, int64(s.locals))
		wrint_objfile(b, int64(s.nosplit))
		wrint_objfile(b, int64(s.leaf))
		n = 0
		for a = s.autom; a != nil; a = a.link {
			n++
		}
		wrint_objfile(b, int64(n))
		for a = s.autom; a != nil; a = a.link {
			wrsym_objfile(b, a.asym)
			wrint_objfile(b, int64(a.aoffset))
			if a.typ == ctxt.arch.D_AUTO {
				wrint_objfile(b, int64(A_AUTO))
			} else {
				if a.typ == ctxt.arch.D_PARAM {
					wrint_objfile(b, int64(A_PARAM))
				} else {
					sysfatal("%s: invalid local variable type %d", s.name, a.typ)
				}
			}
			wrsym_objfile(b, a.gotype)
		}
		pc = s.pcln
		wrdata_objfile(b, pc.pcsp.p)
		wrdata_objfile(b, pc.pcfile.p)
		wrdata_objfile(b, pc.pcline.p)
		wrint_objfile(b, int64(len(pc.pcdata)))
		for i = range pc.pcdata {
			wrdata_objfile(b, pc.pcdata[i].p)
		}
		wrint_objfile(b, int64(len(pc.funcdata)))
		for i = range pc.funcdata {
			wrsym_objfile(b, pc.funcdata[i])
		}
		for i = range pc.funcdata {
			wrint_objfile(b, pc.funcdataoff[i])
		}
		wrint_objfile(b, int64(len(pc.file)))
		for i = range pc.file {
			wrpathsym_objfile(ctxt, b, pc.file[i])
		}
	}
}

func wrint_objfile(b *Biobuf, sval int64) {
	var uv uint64
	var v uint64
	var buf [10]uint8
	var p []uint8
	uv = (uint64(sval) << 1) ^ uint64(int64(sval>>63))
	p = buf[:]
	for v = uv; v >= 0x80; v >>= 7 {
		p[0] = uint8(v | 0x80)
		p = p[1:]
	}
	p[0] = uint8(v)
	p = p[1:]
	Bwrite(b, buf[:int32(-cap(p)+cap(buf))])
}

func wrstring_objfile(b *Biobuf, s string) {
	wrdata_objfile(b, []byte(s))
}

// wrpath writes a path just like a string, but on windows, it
// translates '\\' to '/' in the process.
func wrpath_objfile(ctxt *Link, b *Biobuf, p string) {
	var i int
	var n int
	if !(ctxt.windows != 0) || !strings.Contains(p, `\`) {
		wrstring_objfile(b, p)
		return
	} else {
		n = len(p)
		wrint_objfile(b, int64(n))
		for i = 0; i < n; i++ {
			var tmp int
			if p[i] == '\\' {
				tmp = '/'
			} else {
				tmp = int(p[i])
			}
			Bputc(b, tmp)
		}
	}
}

func wrdata_objfile(b *Biobuf, v []byte) {
	wrint_objfile(b, int64(len(v)))
	Bwrite(b, v)
}

func wrsym_objfile(b *Biobuf, s *LSym) {
	if s == nil {
		wrint_objfile(b, 0)
		wrint_objfile(b, 0)
		return
	}
	wrstring_objfile(b, s.name)
	wrint_objfile(b, int64(s.version))
}

func wrpathsym_objfile(ctxt *Link, b *Biobuf, s *LSym) {
	if s == nil {
		wrint_objfile(b, 0)
		wrint_objfile(b, 0)
		return
	}
	wrpath_objfile(ctxt, b, s.name)
	wrint_objfile(b, int64(s.version))
}

func readsym_objfile(ctxt *Link, f *Biobuf, pkg string, pn string) {
	var i int
	var j int
	var c int
	var t int
	var v int
	var n int
	var size int
	var dupok int
	var ndup_objfile int
	var name string
	var r *Reloc
	var s *LSym
	var dup *LSym
	var pc *Pcln
	var a *Auto
	if Bgetc(f) != 0xfe {
		sysfatal("readsym out of sync")
	}
	t = int(rdint_objfile(f))
	name = expandpkg(rdstring_objfile(f), pkg)
	v = int(rdint_objfile(f))
	if v != 0 && v != 1 {
		sysfatal("invalid symbol version %d", v)
	}
	dupok = int(rdint_objfile(f))
	size = int(rdint_objfile(f))
	if v != 0 {
		v = ctxt.version
	}
	s = linklookup(ctxt, name, v)
	dup = (*LSym)(nil)
	if s.typ != 0 && s.typ != int(SXREF) {
		if s.typ != int(SBSS) && s.typ != int(SNOPTRBSS) && !(dupok != 0) && !(s.dupok != 0) {
			sysfatal("duplicate symbol %s (types %d and %d) in %s and %s", s.name, s.typ, t, s.file, pn)
		}
		if len(s.p) > 0 {
			dup = s
			s = linknewsym(ctxt, ".dup", ndup_objfile)
			ndup_objfile++ // scratch
		}
	}
	s.file = pkg
	s.dupok = dupok
	if t == int(SXREF) {
		sysfatal("bad sxref")
	}
	if t == 0 {
		sysfatal("missing type for %s in %s", name, pn)
	}
	s.typ = t
	if s.size < size {
		s.size = size
	}
	s.gotype = rdsym_objfile(ctxt, f, pkg)
	rddata_objfile(f, &s.p)
	n = int(rdint_objfile(f))
	if n > 0 {
		s.r = make([]Reloc, n)
		for i = 0; i < n; i++ {
			r = &s.r[i]
			r.off = int(rdint_objfile(f))
			r.siz = uint8(rdint_objfile(f))
			r.typ = int(rdint_objfile(f))
			r.add = rdint_objfile(f)
			r.xadd = rdint_objfile(f)
			r.sym = rdsym_objfile(ctxt, f, pkg)
			r.xsym = rdsym_objfile(ctxt, f, pkg)
		}
	}
	if len(s.p) > 0 && dup != nil && len(dup.p) > 0 && !strings.HasPrefix(s.name, "gclocals·") {
		// content-addressed garbage collection liveness bitmap symbol.
		// double check for hash collisions.
		if !bytes.Equal(s.p, dup.p) {
			sysfatal("dupok hash collision for %s in %s and %s", s.name, s.file, pn)
		}
	}
	if s.typ == int(STEXT) {
		s.args = int(rdint_objfile(f))
		s.locals = int(rdint_objfile(f))
		s.nosplit = uint8(rdint_objfile(f))
		s.leaf = uint8(rdint_objfile(f))
		n = int(rdint_objfile(f))
		for i = 0; i < n; i++ {
			a = new(Auto)
			a.asym = rdsym_objfile(ctxt, f, pkg)
			a.aoffset = int32(rdint_objfile(f))
			a.typ = int(rdint_objfile(f))
			a.gotype = rdsym_objfile(ctxt, f, pkg)
			a.link = s.autom
			s.autom = a
		}
		s.pcln = new(Pcln)
		pc = s.pcln
		rddata_objfile(f, &pc.pcsp.p)
		rddata_objfile(f, &pc.pcfile.p)
		rddata_objfile(f, &pc.pcline.p)
		n = int(rdint_objfile(f))
		pc.pcdata = make([]Pcdata, n)
		for i = 0; i < n; i++ {
			rddata_objfile(f, &pc.pcdata[i].p)
		}
		n = int(rdint_objfile(f))
		pc.funcdata = make([]*LSym, n)
		pc.funcdataoff = make([]int64, n)
		for i = 0; i < n; i++ {
			pc.funcdata[i] = rdsym_objfile(ctxt, f, pkg)
		}
		for i = 0; i < n; i++ {
			pc.funcdataoff[i] = rdint_objfile(f)
		}
		n = int(rdint_objfile(f))
		pc.file = make([]*LSym, n)
		for i = 0; i < n; i++ {
			pc.file[i] = rdsym_objfile(ctxt, f, pkg)
		}
		if dup == nil {
			if s.onlist != 0 {
				sysfatal("symbol %s listed multiple times", s.name)
			}
			s.onlist = 1
			if ctxt.etextp != nil {
				ctxt.etextp.next = s
			} else {
				ctxt.textp = s
			}
			ctxt.etextp = s
		}
	}
	if ctxt.debugasm != 0 {
		Bprint(ctxt.bso, "%s ", s.name)
		if s.version != 0 {
			Bprint(ctxt.bso, "v=%d ", s.version)
		}
		if s.typ != 0 {
			Bprint(ctxt.bso, "t=%d ", s.typ)
		}
		if s.dupok != 0 {
			Bprint(ctxt.bso, "dupok ")
		}
		if s.nosplit != 0 {
			Bprint(ctxt.bso, "nosplit ")
		}
		Bprint(ctxt.bso, "size=%lld value=%lld", int64(s.size), int64(s.value))
		if s.typ == int(STEXT) {
			Bprint(ctxt.bso, " args=%#llux locals=%#llux", uint64(s.args), uint64(s.locals))
		}
		Bprint(ctxt.bso, "\n")
		for i = 0; i < len(s.p); {
			Bprint(ctxt.bso, "\t%#06ux", i)
			for j = i; j < i+16 && j < len(s.p); j++ {
				Bprint(ctxt.bso, " %02ux", s.p[j])
			}
			for ; j < i+16; j++ {
				Bprint(ctxt.bso, "   ")
			}
			Bprint(ctxt.bso, "  ")
			for j = i; j < i+16 && j < len(s.p); j++ {
				c = int(s.p[j])
				if ' ' <= c && c <= 0x7e {
					Bprint(ctxt.bso, "%c", c)
				} else {
					Bprint(ctxt.bso, ".")
				}
			}
			Bprint(ctxt.bso, "\n")
			i += 16
		}
		for i = range s.r {
			r = &s.r[i]
			Bprint(ctxt.bso, "\trel %d+%d t=%d %s+%lld\n", int(r.off), r.siz, r.typ, r.sym.name, int64(r.add))
		}
	}
}

func rdint_objfile(f *Biobuf) int64 {
	var c int
	var uv uint64
	var shift int
	uv = 0
	for shift = 0; ; shift += 7 {
		if shift >= 64 {
			sysfatal("corrupt input")
		}
		c = Bgetc(f)
		uv |= uint64(c&0x7F) << uint(shift)
		if !(c&0x80 != 0 /*untyped*/) {
			break
		}
	}
	return int64(uv>>1) ^ (int64(uint64(uv)<<63) >> 63)
}

func rdstring_objfile(f *Biobuf) string {
	var n int32
	var p []byte
	n = int32(rdint_objfile(f))
	p = make([]byte, n)
	Bread(f, p)
	return string(p)
}

func rddata_objfile(f *Biobuf, pp *[]byte) {
	n := rdint_objfile(f)
	*pp = make([]byte, n)
	Bread(f, *pp)
}

func rdsym_objfile(ctxt *Link, f *Biobuf, pkg string) *LSym {
	var n int32
	var v int
	var p []byte
	var s *LSym
	n = int32(rdint_objfile(f))
	if n == 0 {
		rdint_objfile(f)
		return (*LSym)(nil)
	}
	p = make([]byte, n)
	Bread(f, p)
	v = int(rdint_objfile(f))
	if v != 0 {
		v = ctxt.version
	}
	s = linklookup(ctxt, expandpkg(string(p), pkg), v)
	if v == 0 && s.name[0] == '$' && s.typ == 0 {
		if strings.HasPrefix(s.name, "$f32.") {
			u64, _ := strconv.ParseUint(s.name[5:], 16, 32)
			i32 := int32(u64)
			s.typ = int(SRODATA)
			adduint32(ctxt, s, uint32(i32))
			s.reachable = 0
		} else {
			if strings.HasPrefix(s.name, "$f64.") {
				u64, _ := strconv.ParseUint(s.name[5:], 16, 64)
				i64 := int64(u64)
				s.typ = int(SRODATA)
				adduint64(ctxt, s, uint64(i64))
				s.reachable = 0
			}
		}
	}
	return s
}

var startmagic_objfile = "\x00\x00go13ld"

var endmagic_objfile = "\xff\xffgo13ld"
