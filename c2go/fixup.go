// +build ignore

import (
	"math"
	"os"
	"path"
	"strconv"
	"strings"
)

func symgrow(ctxt *Link, s *LSym, lsiz int64) {
	var siz int
	siz = int(lsiz)
	if int64(siz) != lsiz {
		sysfatal("symgrow size %d too long", lsiz)
	}
	if len(s.p) >= siz {
		return
	}
	for cap(s.p) < siz {
		s.p = append(s.p[:cap(s.p)], 0)
	}
	s.p = s.p[:siz]
}

func savedata(ctxt *Link, s *LSym, p *Prog, pn string) {
	var off int
	var siz int32
	var i int32
	var o int64
	var r *Reloc
	off = int(p.from.offset)
	siz = int32(ctxt.arch.datasize(p))
	if off < 0 || siz < 0 || off >= 1<<30 || siz >= 100 {
		mangle(pn)
	}
	symgrow(ctxt, s, int64(off)+int64(siz))
	if p.to.typ == ctxt.arch.D_FCONST {
		switch siz {
		default:
		case 4:
			ctxt.arch.byteOrder.PutUint32(s.p[off:], math.Float32bits(float32(p.to.u.dval)))
		case 8:
			ctxt.arch.byteOrder.PutUint64(s.p[off:], math.Float64bits(p.to.u.dval))
		}
	} else {
		if p.to.typ == ctxt.arch.D_SCONST {
			for i = 0; i < siz; i++ {
				s.p[int32(off)+i] = uint8(p.to.u.sval[i])
			}
		} else {
			if p.to.typ == ctxt.arch.D_CONST {
				if p.to.sym != nil {
					r = addrel(s)
					r.off = off
					r.siz = uint8(siz)
					r.sym = p.to.sym
					r.typ = int(R_ADDR)
					r.add = p.to.offset
					goto out
				}
				o = p.to.offset
				switch siz {
				default:
					ctxt.diag("bad nuxi %d\n%v", siz, ctxt.Pconv(p))
					break
				case 1:
					s.p[off] = byte(o)
				case 2:
					ctxt.arch.byteOrder.PutUint16(s.p[off:], uint16(o))
				case 4:
					ctxt.arch.byteOrder.PutUint32(s.p[off:], uint32(o))
				case 8:
					ctxt.arch.byteOrder.PutUint64(s.p[off:], uint64(o))
				}
			} else {
				if p.to.typ == ctxt.arch.D_ADDR {
					r = addrel(s)
					r.off = off
					r.siz = uint8(siz)
					r.sym = p.to.sym
					r.typ = int(R_ADDR)
					r.add = p.to.offset
				} else {
					ctxt.diag("bad data: %v", ctxt.Pconv(p))
				}
			}
		out:
		}
	}
}

func addrel(s *LSym) *Reloc {
	s.r = append(s.r, Reloc{})
	return &s.r[len(s.r)-1]
}

func setuintxx(ctxt *Link, s *LSym, off int64, v uint64, wid int64) int64 {
	if s.typ == 0 {
		s.typ = int(SDATA)
	}
	s.reachable = 1
	if int64(s.size) < off+wid {
		s.size = int(off + wid)
		symgrow(ctxt, s, int64(s.size))
	}
	switch wid {
	case 1:
		s.p[off] = uint8(v)
		break
	case 2:
		ctxt.arch.byteOrder.PutUint16(s.p[off:], uint16(v))
	case 4:
		ctxt.arch.byteOrder.PutUint32(s.p[off:], uint32(v))
	case 8:
		ctxt.arch.byteOrder.PutUint64(s.p[off:], uint64(v))
	}
	return off + wid
}

func expandpkg(t0 string, pkg string) string {
	return strings.Replace(t0, `"".`, pkg+".", -1)
}

func double2ieee(ieee *uint64, f float64) {
	*ieee = math.Float64bits(f)
}

//c2go:drop emallocz
//c2go:drop estrdup
//c2go:drop erealloc

func addlib(ctxt *Link, src string, obj string, pathname string) {
	var i int

	name := path.Clean(pathname)

	// runtime.a -> runtime
	short := strings.TrimSuffix(name, ".a")

	// already loaded?
	for i := range ctxt.library {
		if ctxt.library[i].pkg == short {
			return
		}
	}

	var pname string
	// runtime -> runtime.a for search
	if (!(ctxt.windows != 0) && name[0] == '/') || (ctxt.windows != 0 && name[1] == ':') {
		pname = name
	} else {
		// try dot, -L "libdir", and then goroot.
		for i = 0; i < ctxt.nlibdir; i++ {
			pname = ctxt.libdir[i] + "/" + name
			if _, err := os.Stat(pname); !os.IsNotExist(err) {
				break
			}
		}
	}
	pname = path.Clean(pname)

	// runtime.a -> runtime
	pname = strings.TrimSuffix(pname, ".a")

	if ctxt.debugvlog > 1 && ctxt.bso != nil {
		Bprint(ctxt.bso, "%5.2f addlib: %s %s pulls in %s\n", cputime(), obj, src, pname)
	}
	addlibpath(ctxt, src, obj, pname, name)
}

func addlibpath(ctxt *Link, srcref string, objref string, file string, pkg string) {
	for i := range ctxt.library {
		if file == ctxt.library[i].file {
			return
		}
	}
	if ctxt.debugvlog > 1 && ctxt.bso != nil {
		Bprint(ctxt.bso, "%5.2f addlibpath: srcref: %s objref: %s file: %s pkg: %s\n", cputime(), srcref, objref, file, pkg)
	}
	ctxt.library = append(ctxt.library, Library{
		objref: objref,
		srcref: srcref,
		file:   file,
		pkg:    pkg,
	})
}

//c2go:drop nuxiinit
//c2go:drop fnuxi4
//c2go:drop fnuxi8
//c2go:drop inuxi1
//c2go:drop inuxi2
//c2go:drop inuxi4
//c2go:drop inuxi8

func atolwhex(s string) int64 {
	x, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		x = 0
	}
	return x
}

//c2go:drop listinit5
//c2go:drop listinit6
//c2go:drop listinit8

