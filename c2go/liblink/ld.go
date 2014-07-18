package main

import (
	"os"
	"path"
	"strconv"
	"strings"
)

func addlib(ctxt *Link, src string, obj string, pathname string) {
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
		for _, dir := range ctxt.libdir {
			pname = dir + "/" + name
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

const (
	LOG_ld = 5
)

func mkfwd(sym *LSym) {
	var p *Prog
	var i int
	var dwn [LOG_ld]int32
	var cnt [LOG_ld]int32
	var lst [LOG_ld]*Prog
	for i = 0; i < LOG_ld; i++ {
		if i == 0 {
			cnt[i] = 1
		} else {
			cnt[i] = LOG_ld * cnt[i-1]
		}
		dwn[i] = 1
		lst[i] = nil
	}
	i = 0
	for p = sym.text; p != nil && p.link != nil; p = p.link {
		i--
		if i < 0 {
			i = LOG_ld - 1
		}
		p.forwd = nil
		dwn[i]--
		if dwn[i] <= 0 {
			dwn[i] = cnt[i]
			if lst[i] != nil {
				lst[i].forwd = p
			}
			lst[i] = p
		}
	}
}

func copyp(ctxt *Link, q *Prog) *Prog {
	var p *Prog
	p = ctxt.prg()
	*p = *q
	return p
}

func appendp(ctxt *Link, q *Prog) *Prog {
	var p *Prog
	p = ctxt.prg()
	p.link = q.link
	q.link = p
	p.lineno = q.lineno
	p.mode = q.mode
	return p
}

func atolwhex(s string) int64 {
	x, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		x = 0
	}
	return x
}
