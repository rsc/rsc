package main

func pciterinit(ctxt *Link, it *Pciter, d *Pcdata) {
	it.d = *d
	it.p = it.d.p
	it.pc = 0
	it.nextpc = 0
	it.value = -1
	it.start = 1
	it.done = 0
	it.pcscale = ctxt.arch.minlc
	pciternext(it)
}

func pciternext(it *Pciter) {
	var v uint32
	var dv int32
	it.pc = it.nextpc
	if it.done != 0 {
		return
	}
	if len(it.p) == 0 {
		it.done = 1
		return
	}
	// value delta
	v = getvarint_pcln(&it.p)
	if v == 0 && !(it.start != 0) {
		it.done = 1
		return
	}
	it.start = 0
	dv = int32(v>>1) ^ (int32(v<<31) >> 31)
	it.value += dv
	// pc delta
	v = getvarint_pcln(&it.p)
	it.nextpc = it.pc + v*uint32(it.pcscale)
}

func linkpcln(ctxt *Link, cursym *LSym) {
	var p *Prog
	var pcln *Pcln
	var i int
	var npcdata int
	var nfuncdata int
	var n int32
	var havepc []uint32
	var havefunc []uint32
	ctxt.cursym = cursym
	pcln = new(Pcln)
	cursym.pcln = pcln
	npcdata = 0
	nfuncdata = 0
	for p = cursym.text; p != nil; p = p.link {
		if p.as == ctxt.arch.APCDATA && p.from.offset >= int64(npcdata) {
			npcdata = int(p.from.offset + 1)
		}
		if p.as == ctxt.arch.AFUNCDATA && p.from.offset >= int64(nfuncdata) {
			nfuncdata = int(p.from.offset + 1)
		}
	}
	pcln.pcdata = make([]Pcdata, npcdata)
	pcln.funcdata = make([]*LSym, nfuncdata)
	pcln.funcdataoff = make([]int64, nfuncdata)
	funcpctab_pcln(ctxt, &pcln.pcsp, cursym, "pctospadj", pctospadj_pcln, nil)
	funcpctab_pcln(ctxt, &pcln.pcfile, cursym, "pctofile", pctofileline_pcln, pcln)
	funcpctab_pcln(ctxt, &pcln.pcline, cursym, "pctoline", pctofileline_pcln, nil)
	// tabulate which pc and func data we have.
	n = ((int32(npcdata)+31)/32 + (int32(nfuncdata)+31)/32) * 4
	havepc = make([]uint32, n/4)
	havefunc = havepc[(npcdata+31)/32:]
	for p = cursym.text; p != nil; p = p.link {
		if p.as == ctxt.arch.AFUNCDATA {
			if (havefunc[p.from.offset/32]>>uint(p.from.offset%32))&1 != 0 /*untyped*/ {
				ctxt.diag("multiple definitions for FUNCDATA $%d", p.from.offset)
			}
			havefunc[p.from.offset/32] |= 1 << uint(p.from.offset%32)
		}
		if p.as == ctxt.arch.APCDATA {
			havepc[p.from.offset/32] |= 1 << uint(p.from.offset%32)
		}
	}
	// pcdata.
	for i = 0; i < npcdata; i++ {
		if (havepc[i/32]>>uint(i%32))&1 != 0 {
			continue
		}
		funcpctab_pcln(ctxt, &pcln.pcdata[i], cursym, "pctopcdata", pctopcdata_pcln, uint32(i))
	}
	// funcdata
	if nfuncdata > 0 {
		for p = cursym.text; p != nil; p = p.link {
			if p.as == ctxt.arch.AFUNCDATA {
				i = int(p.from.offset)
				pcln.funcdataoff[i] = p.to.offset
				if p.to.typ != ctxt.arch.D_CONST {
					// TODO: Dedup.
					//funcdata_bytes += p->to.sym->size;
					pcln.funcdata[i] = p.to.sym
				}
			}
		}
	}
}

// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func addvarint_pcln(ctxt *Link, d *Pcdata, val uint32) {
	var v uint32
	for v = val; v >= 0x80; v >>= 7 {
		d.p = append(d.p, uint8(v|0x80))
	}
	d.p = append(d.p, uint8(v))
}

// funcpctab writes to dst a pc-value table mapping the code in func to the values
// returned by valfunc parameterized by arg. The invocation of valfunc to update the
// current value is, for each p,
//
//	val = valfunc(func, val, p, 0, arg);
//	record val as value at p->pc;
//	val = valfunc(func, val, p, 1, arg);
//
// where func is the function, val is the current value, p is the instruction being
// considered, and arg can be used to further parameterize valfunc.
func funcpctab_pcln(ctxt *Link, dst *Pcdata, fun *LSym, desc string, valfunc func(*Link, *LSym, int32, *Prog, int32, interface{}) int32, arg interface{}) {
	var dbg int
	var i int
	var oldval int32
	var val int32
	var started int32
	var delta uint32
	var pc int
	var p *Prog
	// To debug a specific function, uncomment second line and change name.
	dbg = 0
	//dbg = strcmp(func->name, "main.main") == 0;
	//dbg = strcmp(desc, "pctofile") == 0;
	ctxt.debugpcln += int32(dbg)
	dst.p = dst.p[:0]
	if ctxt.debugpcln != 0 {
		Bprint(ctxt.bso, "funcpctab %s [valfunc=%s]\n", fun.name, desc)
	}
	val = -1
	oldval = val
	if fun.text == nil {
		ctxt.debugpcln -= int32(dbg)
		return
	}
	pc = fun.text.pc
	if ctxt.debugpcln != 0 {
		Bprint(ctxt.bso, "%6llux %6d %P\n", pc, val, fun.text)
	}
	started = 0
	for p = fun.text; p != nil; p = p.link {
		// Update val. If it's not changing, keep going.
		val = valfunc(ctxt, fun, val, p, 0, arg)
		if val == oldval && started != 0 {
			val = valfunc(ctxt, fun, val, p, 1, arg)
			if ctxt.debugpcln != 0 {
				Bprint(ctxt.bso, "%6llux %6s %P\n", int64(p.pc), "", p)
			}
			continue
		}
		// If the pc of the next instruction is the same as the
		// pc of this instruction, this instruction is not a real
		// instruction. Keep going, so that we only emit a delta
		// for a true instruction boundary in the program.
		if p.link != nil && p.link.pc == p.pc {
			val = valfunc(ctxt, fun, val, p, 1, arg)
			if ctxt.debugpcln != 0 {
				Bprint(ctxt.bso, "%6llux %6s %P\n", int64(p.pc), "", p)
			}
			continue
		}
		// The table is a sequence of (value, pc) pairs, where each
		// pair states that the given value is in effect from the current position
		// up to the given pc, which becomes the new current position.
		// To generate the table as we scan over the program instructions,
		// we emit a "(value" when pc == func->value, and then
		// each time we observe a change in value we emit ", pc) (value".
		// When the scan is over, we emit the closing ", pc)".
		//
		// The table is delta-encoded. The value deltas are signed and
		// transmitted in zig-zag form, where a complement bit is placed in bit 0,
		// and the pc deltas are unsigned. Both kinds of deltas are sent
		// as variable-length little-endian base-128 integers,
		// where the 0x80 bit indicates that the integer continues.
		if ctxt.debugpcln != 0 {
			Bprint(ctxt.bso, "%6llux %6d %P\n", int64(p.pc), val, p)
		}
		if started != 0 {
			addvarint_pcln(ctxt, dst, (uint32(p.pc)-uint32(pc))/uint32(ctxt.arch.minlc))
			pc = p.pc
		}
		delta = uint32(val) - uint32(oldval)
		if delta>>31 != 0 {
			delta = 1 | ^(delta << 1)
		} else {
			delta <<= 1
		}
		addvarint_pcln(ctxt, dst, delta)
		oldval = val
		started = 1
		val = valfunc(ctxt, fun, val, p, 1, arg)
	}
	if started != 0 {
		if ctxt.debugpcln != 0 {
			Bprint(ctxt.bso, "%6llux done\n", int64(fun.text.pc)+int64(fun.size))
		}
		addvarint_pcln(ctxt, dst, uint32((fun.value+int64(fun.size)-int64(pc))/int64(ctxt.arch.minlc)))
		addvarint_pcln(ctxt, dst, 0) // terminator
	}
	if ctxt.debugpcln != 0 {
		Bprint(ctxt.bso, "wrote %d bytes to %p\n", len(dst.p), dst)
		for i = range dst.p {
			Bprint(ctxt.bso, " %02ux", dst.p[i])
		}
		Bprint(ctxt.bso, "\n")
	}
	ctxt.debugpcln -= int32(dbg)
}

// pctofileline computes either the file number (arg == 0)
// or the line number (arg == 1) to use at p.
// Because p->lineno applies to p, phase == 0 (before p)
// takes care of the update.
func pctofileline_pcln(ctxt *Link, sym *LSym, oldval int32, p *Prog, phase int32, arg interface{}) int32 {
	var i int
	var l int32
	var f *LSym
	var pcln *Pcln
	if p.as == ctxt.arch.ATEXT || p.as == ctxt.arch.ANOP || p.as == ctxt.arch.AUSEFIELD || p.lineno == 0 || phase == 1 {
		return oldval
	}
	linkgetline(ctxt, p.lineno, &f, (*int32)(&l))
	if f == nil {
		//	print("getline failed for %s %P\n", ctxt->cursym->name, p);
		return oldval
	}
	if arg == nil {
		return l
	}
	pcln = arg.(*Pcln)
	if f == pcln.lastfile {
		return pcln.lastindex
	}
	for i = range pcln.file {
		if pcln.file[i] == f {
			pcln.lastfile = f
			pcln.lastindex = int32(i)
			return int32(i)
		}
	}
	pcln.file = append(pcln.file, f)
	pcln.lastfile = f
	pcln.lastindex = int32(i)
	return int32(i)
}

// pctospadj computes the sp adjustment in effect.
// It is oldval plus any adjustment made by p itself.
// The adjustment by p takes effect only after p, so we
// apply the change during phase == 1.
func pctospadj_pcln(ctxt *Link, sym *LSym, oldval int32, p *Prog, phase int32, arg interface{}) int32 {
	if oldval == -1 { // starting
		oldval = 0
	}
	if phase == 0 {
		return oldval
	}
	if oldval+int32(p.spadj) < -10000 || oldval+int32(p.spadj) > 1100000000 {
		ctxt.diag("overflow in spadj: %d + %d = %d", oldval, p.spadj, oldval+int32(p.spadj))
		sysfatal("bad code")
	}
	return oldval + int32(p.spadj)
}

// pctopcdata computes the pcdata value in effect at p.
// A PCDATA instruction sets the value in effect at future
// non-PCDATA instructions.
// Since PCDATA instructions have no width in the final code,
// it does not matter which phase we use for the update.
func pctopcdata_pcln(ctxt *Link, sym *LSym, oldval int32, p *Prog, phase int32, arg interface{}) int32 {
	if phase == 0 || p.as != ctxt.arch.APCDATA || p.from.offset != int64(arg.(uint32)) {
		return oldval
	}
	if int64(int32(p.to.offset)) != p.to.offset {
		ctxt.diag("overflow in PCDATA instruction: %P", p)
		sysfatal("bad code")
	}
	return int32(p.to.offset)
}

// iteration over encoded pcdata tables.
func getvarint_pcln(pp *[]uint8) uint32 {
	var p []uint8
	var shift int
	var v uint32
	v = 0
	p = *pp
	for shift = 0; ; shift += 7 {
		v |= uint32(p[0]&0x7F) << uint(shift)
		var tmp []uint8 = p
		p = p[1:]
		if !(tmp[0]&0x80 != 0 /*untyped*/) {
			break
		}
	}
	*pp = p
	return v
}
