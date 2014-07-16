package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Create("x.6")
	if err != nil {
		log.Fatal(err)
	}
	ctxt := linknew(&linkamd64)
	ctxt.diag = log.Fatalf
	obuf := Binitw(f)
	Bprint(obuf, "go object %s %s %s\n", getgoos(), getgoarch(), getgoversion())
	Bprint(obuf, "!\n")
	p1 := &Prog{
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
		as: ARET_6,
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
		as: ARET_6,
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
