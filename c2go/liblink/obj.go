package main

import "fmt"

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
const (
	HISTSZ_obj = 10
	NSYM_obj   = 50
)

func linklinefmt(ctxt *Link, lno int, showAll, showFullPath bool) string {
	var a [HISTSZ_obj]struct {
		incl *Hist
		idel int
		line *Hist
		ldel int32
	}
	lno1 := lno
	var d int32
	var i int
	var n int
	var h *Hist
	n = 0
	var fp string
	for h = ctxt.hist; h != nil; h = h.link {
		if h.offset < 0 {
			continue
		}
		if lno < int(h.line) {
			break
		}
		if h.name != "XXXXXXX" {
			if h.offset > 0 {
				// #line directive
				if n > 0 && n < int(HISTSZ_obj) {
					a[n-1].line = h
					a[n-1].ldel = int32(h.line) - h.offset + 1
				}
			} else {
				// beginning of file
				if n < int(HISTSZ_obj) {
					a[n].incl = h
					a[n].idel = int(h.line)
					a[n].line = nil
				}
				n++
			}
			continue
		}
		n--
		if n > 0 && n < int(HISTSZ_obj) {
			d = int32(h.line) - int32(a[n].incl.line)
			a[n-1].ldel += d
			a[n-1].idel += int(d)
		}
	}
	if n > int(HISTSZ_obj) {
		n = int(HISTSZ_obj)
	}
	for i = n - 1; i >= 0; i-- {
		if i != n-1 {
			if !showAll {
				break
			}
			fp += " "
		}
		if ctxt.debugline != 0 || showFullPath {
			fp += fmt.Sprintf("%s/", ctxt.pathname)
		}
		if a[i].line != nil {
			fp += fmt.Sprintf("%s:%d[%s:%d]", a[i].line.name, int32(lno)-a[i].ldel+1, a[i].incl.name, lno-a[i].idel+1)
		} else {
			fp += fmt.Sprintf("%s:%d", a[i].incl.name, lno-a[i].idel+1)
		}
		lno = int(a[i].incl.line - 1) // now print out start of this file
	}
	if n == 0 {
		fp += fmt.Sprintf("<unknown line number %d %d %d %s>", lno1, ctxt.hist.offset, ctxt.hist.line, ctxt.hist.name)
	}
	return fp
}

// Does s have t as a path prefix?
// That is, does s == t or does s begin with t followed by a slash?
// For portability, we allow ASCII case folding, so that haspathprefix("a/b/c", "A/B") is true.
// Similarly, we allow slash folding, so that haspathprefix("a/b/c", "a\\b") is true.
func haspathprefix_obj(s string, t string) bool {
	var i int
	var cs int
	var ct int
	if len(t) > len(s) {
		return false
	}
	for i = 0; i < len(t); i++ {
		cs = int(s[i])
		ct = int(t[i])
		if 'A' <= cs && cs <= 'Z' {
			cs += 'a' - 'A'
		}
		if 'A' <= ct && ct <= 'Z' {
			ct += 'a' - 'A'
		}
		if cs == '\\' {
			cs = '/'
		}
		if ct == '\\' {
			ct = '/'
		}
		if cs != ct {
			return false
		}
	}
	return i == len(s) || s[i] == '/' || s[i] == '\\'
}

// This is a simplified copy of linklinefmt above.
// It doesn't allow printing the full stack, and it returns the file name and line number separately.
// TODO: Unify with linklinefmt somehow.
func linkgetline(ctxt *Link, line int32, f **LSym, l *int32) {
	var a [HISTSZ_obj]struct {
		incl *Hist
		idel int
		line *Hist
		ldel int32
	}
	var lno int
	var d int32
	var dlno int32
	var n int
	var h *Hist
	var buf string
	var file string
	lno = int(line)
	lno0 := lno
	n = 0
	for h = ctxt.hist; h != nil; h = h.link {
		if h.offset < 0 {
			continue
		}
		if lno < int(h.line) {
			break
		}
		if h.name != "XXXXXXX" {
			if h.offset > 0 {
				// #line directive
				if n > 0 && n < int(HISTSZ_obj) {
					a[n-1].line = h
					a[n-1].ldel = int32(h.line) - h.offset + 1
				}
			} else {
				// beginning of file
				if n < int(HISTSZ_obj) {
					if n < 0 {
						println("linkgetline", lno0, n, h)
						for h := ctxt.hist; h != nil; h = h.link {
							fmt.Printf("%p %v\n", h, *h)
						}
						panic("linkgetline")
					}
					a[n].incl = h
					a[n].idel = int(h.line)
					a[n].line = nil
				}
				n++
			}
			continue
		}
		n--
		if n > 0 && n < int(HISTSZ_obj) {
			d = int32(h.line) - int32(a[n].incl.line)
			a[n-1].ldel += d
			a[n-1].idel += int(d)
		}
	}
	if n > int(HISTSZ_obj) {
		n = int(HISTSZ_obj)
	}
	if n <= 0 {
		*f = linklookup(ctxt, "??", HistVersion)
		*l = 0
		return
	}
	n--
	if a[n].line != nil {
		file = a[n].line.name
		dlno = a[n].ldel - 1
	} else {
		file = a[n].incl.name
		dlno = int32(a[n].idel) - 1
	}
	if (!(ctxt.windows != 0) && file[0] == '/') || (ctxt.windows != 0 && file[1] == ':') || file[0] == '<' {
		buf = file
	} else {
		buf = fmt.Sprintf("%s/%s", ctxt.pathname, file)
	}
	// Remove leading ctxt->trimpath, or else rewrite $GOROOT to $GOROOT_FINAL.
	if ctxt.trimpath != "" && haspathprefix_obj(buf, ctxt.trimpath) {
		if len(buf) == len(ctxt.trimpath) {
			buf = "??"
		} else {
			buf1 := buf[len(ctxt.trimpath)+1:]
			if buf1 == "" {
				buf1 = "??"
			}
			buf = buf1
		}
	} else {
		if ctxt.goroot_final != "" && haspathprefix_obj(buf, ctxt.goroot) {
			buf = ctxt.goroot_final + buf[len(ctxt.goroot):]
		}
	}
	lno -= int(dlno)
	*f = linklookup(ctxt, buf, HistVersion)
	*l = int32(lno)
}

func linklinehist(ctxt *Link, lineno int32, f string, offset int32) {
	var h *Hist
	if false { // debug['f']
		if f != "" {
			if offset != 0 {
				fmt.Printf("%4d: %s (#line %d)\n", lineno, f, offset)
			} else {
				fmt.Printf("%4d: %s\n", lineno, f)
			}
		} else {
			fmt.Printf("%4d: <pop>\n", lineno)
		}
	}
	h = new(Hist)
	*h = Hist{}
	h.name = f
	h.line = lineno
	h.offset = offset
	h.link = nil
	if ctxt.ehist == nil {
		ctxt.hist = h
		ctxt.ehist = h
		return
	}
	ctxt.ehist.link = h
	ctxt.ehist = h
}

func linkprfile(ctxt *Link, l int32) {
	var i int
	var n int
	var a [HISTSZ_obj]Hist
	var h *Hist
	var d int32
	n = 0
	for h = ctxt.hist; h != nil; h = h.link {
		if l < h.line {
			break
		}
		if h.name != "XXXXXXX" {
			if h.offset == 0 {
				if n >= 0 && n < int(HISTSZ_obj) {
					a[n] = *h
				}
				n++
				continue
			}
			if n > 0 && n < int(HISTSZ_obj) {
				if a[n-1].offset == 0 {
					a[n] = *h
					n++
				} else {
					a[n-1] = *h
				}
			}
			continue
		}
		n--
		if n >= 0 && n < int(HISTSZ_obj) {
			d = int32(h.line) - int32(a[n].line)
			for i = 0; i < n; i++ {
				a[i].line += d
			}
		}
	}
	if n > int(HISTSZ_obj) {
		n = int(HISTSZ_obj)
	}
	for i = 0; i < n; i++ {
		print("%s:%d ", a[i].name, int32(int32(l-a[i].line)+a[i].offset+1))
	}
}

/*
 * start a new Prog list.
 */
func linknewplist(ctxt *Link) *Plist {
	var pl *Plist
	pl = new(Plist)
	*pl = Plist{}
	if ctxt.plist == nil {
		ctxt.plist = pl
	} else {
		ctxt.plast.link = pl
	}
	ctxt.plast = pl
	return pl
}
