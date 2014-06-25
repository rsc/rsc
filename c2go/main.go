// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"code.google.com/p/rsc/c2go"
	"code.google.com/p/rsc/cc"
)

var (
	src   = flag.String("src", "", "src of search")
	dst   = flag.String("dst", "", "dst of search")
	out   = flag.String("o", "/tmp/c2go", "output directory")
	strip = flag.String("", "", "strip from input paths when writing in output directory")
	inc   = flag.String("I", "", "include directory")
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	if *inc != "" {
		cc.AddInclude(*inc)
	}
	if len(args) == 0 {
		cc.Read("<stdin>", os.Stdin)
	} else {
		var r []io.Reader
		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			r = append(r, f)
			defer f.Close()
		}
		prog, err := cc.ReadMany(args, r)
		if err != nil {
			log.Fatal(err)
		}
		do(prog, args)
		write(prog, args)
	}
}

type byLen [][]*TypeVar

func (x byLen) Len() int      { return len(x) }
func (x byLen) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byLen) Less(i, j int) bool {
	if len(x[i]) != len(x[j]) {
		return len(x[i]) < len(x[j])
	}
	return fmt.Sprint(x[i][0].Src) < fmt.Sprint(x[j][0].Src)
}

func do(prog *cc.Prog, files []string) {
	for _, decl := range prog.Decls {
		addDecl(decl)
	}

	groups := map[*TypeVar][]*TypeVar{}
	for _, tv := range typeVars {
		upwalk(tv)
		groups[tv.parent] = append(groups[tv.parent], tv)
	}

	changed := true
	for changed {
		changed = false
		for _, group := range groups {
			var indir *TypeVar
			for _, tv := range group {
				if tv.Indir != nil {
					if indir == nil {
						indir = tv.Indir
					} else {
						changed = merge(indir, tv.Indir) || changed
					}
				}
				/*
					if tv.Addr != nil {
						if addr == nil {
							addr = tv.Addr
						} else {
							changed = merge(addr, tv.Addr) || changed
						}
					}
				*/
			}
		}
	}

	groups = map[*TypeVar][]*TypeVar{}
	for _, tv := range typeVars {
		upwalk(tv)
		groups[tv.parent] = append(groups[tv.parent], tv)
	}

	if false {
		fmt.Printf("groups:\n")
	}
	var grps [][]*TypeVar
	for _, list := range groups {
		grps = append(grps, list)
	}

	sort.Sort(byLen(grps))

	if true {
		for _, list := range grps {
			if len(list) == 1 {
				continue
			}
			decls := map[*cc.Decl]bool{}
			var ptr *cc.Decl
			var ptrBad = false
			for _, tv := range list {
				if d, ok := tv.Src.(*cc.Decl); ok {
					decls[d] = true
					if !ptrBad && d.Type != nil && d.Type.Kind == cc.Ptr && d.Type.Base.Kind == cc.TypedefType {
						if ptr == nil || ptr.Type.Base.Name == d.Type.Base.Name {
							ptr = d
						} else {
							ptrBad = true
							fmt.Fprintf(os.Stderr, "grouped %v and %v\n", ptr, d)
							findPath(typeVars[ptr], typeVars[d])
						}
					}
				}
			}
			var v1, v2 *TypeVar
			for _, tv := range list {
				if x, ok := tv.Src.(*cc.Expr); ok {
					if x.Op == cc.Call && decls[x.Left.XDecl] || x.Op == cc.Eq {
						continue
					}
				}
				str := fmt.Sprint(tv.Src)
				if src != nil && str == *src {
					v1 = tv
				}
				if dst != nil && str == *dst {
					v2 = tv
				}
			}
			if v1 != nil && v2 != nil {
				findPath(v1, v2)
			}
		}
	}

	for _, list := range grps {
		if false {
			fmt.Printf("group(%d)\n", len(list))
		}
		var types []string
		var haveType = map[string]bool{}
		var ops []string
		var haveOp = map[string]bool{"": true}
		for _, tv := range list {
			x, ok := tv.Src.(*cc.Expr)
			if ok && x.XType != nil {
				str := x.XType.String()
				if !haveType[str] {
					haveType[str] = true
					types = append(types, str)
				}
			}
			for _, ctxt1 := range tv.Context {
				ctxt, ok := ctxt1.(*cc.Expr)
				if !ok {
					continue
				}
				op := ""
				switch ctxt.Op {
				case cc.PreInc, cc.PostInc, cc.PreDec, cc.PostDec:
					op = "++"
					if ctxt.Op == cc.PreDec || ctxt.Op == cc.PostDec {
						op = "--"
					}
					if x != nil && x.XType.Is(cc.Ptr) {
						op = "ptr" + op
						if op == "ptr--" && false {
							fmt.Println("ptr--: %s\n", x.GetSpan())
						}
					}
				case cc.Add, cc.AddEq:
					if x != nil && x.XType.Is(cc.Ptr) {
						op = "ptr+"
					} else {
						op = "+"
					}
				case cc.Sub, cc.SubEq:
					if x != nil && x.XType.Is(cc.Ptr) {
						op = "ptr-"
					} else {
						op = "-"
					}
				case cc.Indir:
					op = "*"
				case cc.Arrow:
					op = "->"
				case cc.Dot:
					op = "."
				case cc.Addr:
					op = "&"
				case cc.Index:
					op = "[i]"
				}
				if !haveOp[op] {
					haveOp[op] = true
					ops = append(ops, op)
				}
			}
		}
		sort.Strings(types)
		sort.Strings(ops)
		if false {
			fmt.Printf("types: %s\n", strings.Join(types, ", "))
			fmt.Printf("ops: %s\n", strings.Join(ops, ", "))
			for _, tv := range list {
				fmt.Printf("\t%s:", tv.Src)
				if tv.UsedAsBool {
					fmt.Printf(" bool")
				}
				if tv.StmtExpr {
					fmt.Printf(" stmt")
				}
				for _, x := range tv.Context {
					fmt.Printf(" %s", x)
				}
				fmt.Printf("\n")
			}
		}
	}

	if false {
		for x, tv := range typeVars {
			fmt.Printf("%s: typevar for %T %s:", x.GetSpan(), x, x)
			if tv.UsedAsBool {
				fmt.Printf(" bool")
			}
			fmt.Printf("\n")
		}
		return
	}
}

type TypeVar struct {
	Src        cc.Syntax
	UsedAsBool bool
	StmtExpr   bool
	NoLink     bool
	Link       []*TypeVar
	Indir      *TypeVar
	Addr       *TypeVar
	Context    []cc.Syntax

	parent *TypeVar
}

var typeVars = map[cc.Syntax]*TypeVar{}
var typeVarsIndir = map[cc.Syntax]*TypeVar{}
var typeVarsAddr = map[cc.Syntax]*TypeVar{}

var curFn *cc.Decl
var seen = map[cc.Syntax]bool{}

func add(x cc.Syntax) {
	switch x := x.(type) {
	case *cc.Decl:
		addDecl(x)
	case *cc.Init:
		addInit(x)
	case *cc.Type:
		addType(x, nil)
	case *cc.Expr:
		addExpr(x)
	case *cc.Stmt:
		addStmt(x)
	default:
		panic(fmt.Sprintf("unexpected type %T in add", x))
	}
}

func addDecl(x *cc.Decl) {
	if x == nil || typeVars[x] != nil {
		return
	}

	// Note:
	// For function declarations, the typeVar represents
	// the of the function's result (the type of a call to the function).
	// For other declarations, the typeVar represents the
	// type of the declared variable.
	if x.Storage&cc.Typedef == 0 {
		tv := &TypeVar{Src: x}
		switch x.Name {
		case "nil", "N", "L", "S", "T", "C", "...", "bval", "vval", "mpgetfix", "smprint", "namebuf":
			tv.NoLink = true
		}
		if x.Type != nil && x.Type.Is(cc.Ptr) && x.Type.Base.Is(cc.Void) {
			tv.NoLink = true
		}
		typeVars[x] = tv
	}

	addInit(x.Init)
	addType(x.Type, x)

	//	if x.Type != nil && x.Type.Kind == cc.Enum {
	//		addLink(x, x.Type)
	//	}

	if x.Body != nil {
		curFn = x
		addStmt(x.Body)
		curFn = nil
	}
}

func addInit(x *cc.Init) {
	if x == nil {
		return
	}

	addExpr(x.Expr)
	for _, init := range x.Braced {
		addInit(init)
	}
	for _, pre := range x.Prefix {
		addExpr(pre.Index)
	}
}

func addType(x *cc.Type, d *cc.Decl) {
	if x == nil || seen[x] {
		return
	}
	seen[x] = true

	//	if x.Kind == cc.Enum {
	//		typeVars[x] = &TypeVar{Src: x}
	//	}

	addType(x.Base, nil)

	addExpr(x.Width)
	for _, decl := range x.Decls {
		decl.XOuter = d
		addDecl(decl)
	}
}

func addExpr(x *cc.Expr) {
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x}
	typeVars[x] = tv

	if x.XType != nil && x.XType.Kind == cc.Ptr && x.XType.Base.Is(cc.Void) {
		tv.NoLink = true
	}

	addExpr(x.Left)
	addExpr(x.Right)
	for _, y := range x.List {
		addExpr(y)
	}
	addType(x.Type, nil)
	addInit(x.Init)

	if x.Left != nil {
		addContext(x.Left, x)
	}
	if x.Right != nil {
		addContext(x.Right, x)
	}

	switch x.Op {
	case cc.Add:
	case cc.AddEq:
		addLink(x, x.Left)
	case cc.Addr:
		addAddrLink(x.Left, x)
		addIndirLink(x, x.Left)

	case cc.AndAnd, cc.OrOr, cc.Not:
		addBool(x.Left)
		addBool(x.Right)
	case cc.Arrow, cc.Dot:
		if x.XDecl != nil {
			addLink(x, x.XDecl)
		}
	case cc.Call:
		if x.Left.Op == cc.Name {
			switch x.Left.Text {
			case "duintptr", "nodconst", "duintxx", "duint32", "duint8", "duint16", "duint64", "nodintconst", "mpmovecfix", "strcmp", "memcmp", "memmove", "strcpy", "strlen", "strncmp", "strstr", "strchr", "strrchr":
				return
			}
		}
		if x.Left.Op == cc.Name {
			addLink(x, x.Left.XDecl)
		}
		typ := x.Left.XType
		if typ == nil {
			println(x.GetSpan().String(), "no type for", x.Left.String(), "in", x.String())
			break
		}
		for i := 0; i < len(typ.Decls) && i < len(x.List); i++ {
			addLink(typ.Decls[i], x.List[i])
		}
	case cc.Cast:
	case cc.CastInit:
	case cc.Comma:
	case cc.Cond:
		addBool(x.List[0])
		addContext(x.List[0], x)
		addContext(x.List[1], x)
		addContext(x.List[2], x)
		addLink(x.List[1], x.List[2])

	case cc.Div, cc.Mul:
	case cc.DivEq, cc.MulEq:

	case cc.Eq, cc.AndEq, cc.OrEq, cc.XorEq:
		addLink(x, x.Left)
		addLink(x.Left, x.Right)

	case cc.EqEq, cc.NotEq, cc.And, cc.Or, cc.Xor:
		addLink(x.Left, x.Right)

	case cc.Gt, cc.GtEq, cc.Lt, cc.LtEq:
		addLink(x.Left, x.Right)

	case cc.Index:

	case cc.Indir:
		addIndirLink(x.Left, x)
		addAddrLink(x, x.Left)

	case cc.Lsh, cc.Rsh:
	case cc.LshEq, cc.RshEq:
	case cc.Minus, cc.Plus:
	case cc.Mod:
	case cc.ModEq:

	case cc.Name:
		if x.XDecl != nil && x.XType != nil && x.XType.Kind != cc.Func {
			addLink(x, x.XDecl)
		}

	case cc.Number:
	case cc.Offsetof:
	case cc.Paren:

	case cc.PostDec, cc.PostInc, cc.PreDec, cc.PreInc:
		addLink(x, x.Left)

	case cc.SizeofExpr:
	case cc.SizeofType:
	case cc.String:
	case cc.Sub:
	case cc.SubEq:
	case cc.Twid:
	case cc.VaArg:
	}
}

func upwalk(tv *TypeVar) *TypeVar {
	if tv.parent == nil {
		tv.parent = tv
	}
	if tv.parent != tv {
		tv.parent = upwalk(tv.parent)
	}
	return tv.parent
}

func addContext(x, y cc.Syntax) {
	tx := typeVars[x]
	tx.Context = append(tx.Context, y)
}

func addIndirLink(x, y cc.Syntax) {
	add(x)
	add(y)

	tx := typeVars[x]
	ty := typeVars[y]
	if tx == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", x.GetSpan(), x, x, x)
		return
	}
	if ty == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", y.GetSpan(), y, y, y)
		return
	}

	tx.Indir = ty
}

func addAddrLink(x, y cc.Syntax) {
	add(x)
	add(y)

	tx := typeVars[x]
	ty := typeVars[y]
	if tx == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", x.GetSpan(), x, x, x)
		return
	}
	if ty == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", y.GetSpan(), y, y, y)
		return
	}

	tx.Addr = ty
}

func addLink(x, y cc.Syntax) {
	add(x)
	add(y)

	tx := typeVars[x]
	ty := typeVars[y]
	if tx == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", x.GetSpan(), x, x, x)
		return
	}
	if ty == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", y.GetSpan(), y, y, y)
		return
	}
	if tx.NoLink || ty.NoLink {
		return
	}

	tx.Link = append(tx.Link, ty)
	ty.Link = append(ty.Link, tx)

	merge(tx, ty)
}

func merge(tx, ty *TypeVar) bool {
	upwalk(tx)
	upwalk(ty)
	if tx.parent == ty.parent {
		return false
	}
	tx.parent.parent = upwalk(ty)
	return true
}

func addStmt(x *cc.Stmt) {
	if x == nil || seen[x] {
		return
	}
	seen[x] = true

	addExpr(x.Pre)
	addExpr(x.Expr)
	addExpr(x.Post)
	addDecl(x.Decl)
	addStmt(x.Body)
	addStmt(x.Else)
	for _, stmt := range x.Block {
		addStmt(stmt)
	}
	for _, label := range x.Labels {
		addExpr(label.Expr)
	}

	switch x.Op {
	case cc.StmtDecl:

	case cc.StmtExpr:
		typeVars[x.Expr].StmtExpr = true

	case cc.Empty:
	case cc.Block:
	case cc.ARGBEGIN:
	case cc.Break:
	case cc.Continue:
	case cc.Goto:

	case cc.Do, cc.For, cc.If, cc.While:
		if x.Pre != nil {
			typeVars[x.Pre].StmtExpr = true
		}
		if x.Post != nil {
			typeVars[x.Post].StmtExpr = true
		}
		if x.Expr != nil {
			addBool(x.Expr)
		}

	case cc.Return:
		if x.Expr != nil && curFn != nil {
			addLink(x.Expr, curFn)
		}

	case cc.Switch:
		// comparison of x.Expr with case list expressions.
		for _, stmt := range x.Body.Block {
			for _, lab := range stmt.Labels {
				if lab.Op == cc.Case && lab.Expr != nil {
					addLink(x.Expr, lab.Expr)
				}
			}
		}
	}
}

func addBool(x *cc.Expr) {
	if x == nil {
		return
	}
	tv := typeVars[x]
	tv.UsedAsBool = true
}

func fix(decl *cc.Decl) {
	fixDecl(decl)
	fixStmt(decl.Body)
}

func fixStmt(stmt *cc.Stmt) {
	if stmt == nil {
		return
	}
	fixExpr(stmt.Pre)
	switch stmt.Op {
	case cc.If, cc.For, cc.Do, cc.While:
		fixCond(stmt)
	}
	if stmt.Op == cc.StmtExpr || stmt.Op == cc.Return {
		fixStmtExpr(stmt)
	}
	fixExpr(stmt.Expr)
	fixExpr(stmt.Post)
	fixStmt(stmt.Body)
	fixStmt(stmt.Else)
	for _, s := range stmt.Block {
		fixStmt(s)
	}
}

func fixCond(stmt *cc.Stmt) {
	cond := stmt.Expr
	if cond == nil {
		return
	}
	var before, after []*cc.Stmt
	fixPlusPlus(cond, &before, &after, false)
	if len(before) > 0 {
		if len(before) == 1 {
			stmt.Pre = before[0].Expr
		} else {
			println("too many before in cond")
		}
	}
	if len(after) > 0 {
		println("too many after in cond")
	}
	if stmt.Pre == nil && cond.Left != nil && cond.Left.Op == cc.PreInc {
		stmt.Pre = &cc.Expr{Op: cc.PostInc, Left: cond.Left.Left}
		cond.Left = cond.Left.Left
	}
}

func fixStmtExpr(stmt *cc.Stmt) {
	if stmt.Expr == nil {
		return
	}
	var before, after []*cc.Stmt
	fixPlusPlus(stmt.Expr, &before, &after, true)
	if stmt.Op == cc.Return && len(after) > 0 {
		println("after in return")
	}
	if len(before)+len(after) > 0 {
		y := *stmt
		stmt.Block = append(append(before, &y), after...)
		stmt.Op = c2go.BlockNoBrace
	}
}

func fixPlusPlus(x *cc.Expr, before, after *[]*cc.Stmt, top bool) {
	if x.Left != nil {
		fixPlusPlus(x.Left, before, after, false)
	}
	if x.Right != nil {
		fixPlusPlus(x.Right, before, after, false)
	}
	for _, y := range x.List {
		fixPlusPlus(y, before, after, false)
	}

	if top {
		switch x.Op {
		case cc.PreInc:
			x.Op = cc.PostInc
			return
		case cc.PreDec:
			x.Op = cc.PostDec
			return
		case cc.PostInc, cc.PostDec:
			return
		case cc.Eq:
			return
		}
	}

	var list *[]*cc.Stmt
	var op cc.ExprOp

	switch x.Op {
	case cc.Eq:
		list = before
		*list = append(*list, &cc.Stmt{Op: cc.StmtExpr, Expr: &cc.Expr{Op: x.Op, Left: x.Left, Right: x.Right}})
		*x = *x.Left
		return
	case cc.PreInc:
		list = before
		op = cc.PostInc
	case cc.PreDec:
		list = before
		op = cc.PostDec
	case cc.PostInc, cc.PostDec:
		list = after
		op = x.Op
	}
	if list != nil {
		*list = append(*list, &cc.Stmt{Op: cc.StmtExpr, Expr: &cc.Expr{Op: op, Left: x.Left}})
		*x = *x.Left
	}
}

func fixDecl(x *cc.Decl) {
	if x == nil || fixed[x] {
		return
	}
	fixed[x] = true
	fixName(&x.Name)
	if x.Storage&cc.Static != 0 {
		file := filepath.Base(x.Span.Start.File)
		if i := strings.Index(file, "."); i >= 0 {
			file = file[:i]
		}
		x.Name += "_" + file
	}
	if x.Init != nil {
		fixInit(x.Init)
	}
	fixType(x.Type)
}

func fixInit(x *cc.Init) {
	for _, pre := range x.Prefix {
		fixName(&pre.Dot)
	}
	for _, b := range x.Braced {
		fixInit(b)
	}
	if x.Expr != nil {
		fixExpr(x.Expr)
	}
}

func fixName(s *string) {
	switch *s {
	case "type":
		*s = "typ"
	case "func":
		*s = "fun"
	}
}

var fixed = map[interface{}]bool{}

func fixType(x *cc.Type) {
	if x == nil || fixed[x] {
		return
	}
	fixed[x] = true
	fixType(x.Base)
	for _, d := range x.Decls {
		fixDecl(d)
	}
	fixExpr(x.Width)
}

func fixExpr(x *cc.Expr) {
	if x == nil {
		return
	}
	if x.Init != nil {
		fixInit(x.Init)
	}
	fixExpr(x.Left)
	fixExpr(x.Right)
	for _, y := range x.List {
		fixExpr(y)
	}
	
	switch x.Op {
	case cc.Arrow, cc.Dot, cc.Name:
		//fixName(&x.Text)
	case cc.Number:
		if x.Text == `'\0'` {
			x.Text = `'\x00'`
		}
	}

	if x.Op == cc.Call && x.Left.Op == cc.Name && x.Left.Text == "print" {
		// TODO insert fmt import
		x.Left.Text = "fmt.Printf"
	}
	if x.Op == cc.Call && x.Left.Op == cc.Name && x.Left.Text == "exits" {
		// TODO insert os import
		x.Left.Text = "os.Exit"
		if len(x.List) == 1 {
			arg := x.List[0]
			if arg.Op == cc.Name && arg.Text == "nil" {
				arg.Text = "0"
			} else if arg.Op == cc.String && len(arg.Texts) == 1 {
				arg.Op = cc.Name
				if arg.Texts[0] == `""` {
					arg.Text = "0"
				} else {
					arg.Text = "1"
					x.Comments.Suffix = append(x.Comments.Suffix, cc.Comment{Text: " // " + arg.Texts[0]})
				}
			}
		}
	}
}

type varLevel struct {
	t *TypeVar
	n int
}

func findPath(v1, v2 *TypeVar) {
	var q []varLevel
	onq := make(map[varLevel]bool)
	prev := make(map[varLevel]varLevel)
	q = append(q, varLevel{v2, 0})
	onq[varLevel{v2, 0}] = true
Search:
	for i := 0; i < len(q); i++ {
		v := q[i]
		for _, l := range v.t.Link {
			lv := varLevel{l, v.n}
			if !onq[lv] {
				prev[lv] = v
				if lv == (varLevel{v1, 0}) {
					break Search
				}
				q = append(q, lv)
				onq[lv] = true
			}
		}
		if v.t.Indir != nil {
			lv := varLevel{v.t.Indir, v.n - 1}
			if !onq[lv] {
				prev[lv] = v
				if lv == (varLevel{v1, 0}) {
					break Search
				}
				q = append(q, lv)
				onq[lv] = true
			}
		}
		if v.t.Addr != nil {
			lv := varLevel{v.t.Addr, v.n + 1}
			if !onq[lv] {
				prev[lv] = v
				if lv == (varLevel{v1, 0}) {
					break Search
				}
				q = append(q, lv)
				onq[lv] = true
			}
		}
	}
	if prev[varLevel{v1, 0}].t == nil {
		println("no path found")
	}
	for v := (varLevel{v1, 0}); ; v = prev[v] {
		fmt.Fprintf(os.Stderr, "\t%s: %v [%d]\n", v.t.Src.GetSpan(), v.t.Src, v.n)
		if v == (varLevel{v2, 0}) {
			break
		}
	}
}

func write(prog *cc.Prog, files []string) {
	for _, decl := range prog.Decls {
		fix(decl)
	}
	for _, file := range files {
		writeFile(prog, file)
	}
}

func writeFile(prog *cc.Prog, file string) {
	dstfile := strings.TrimSuffix(strings.TrimSuffix(file, ".c"), ".h") + ".go"
	if *strip != "" {
		dstfile = strings.TrimPrefix(dstfile, *strip)
	} else if i := strings.LastIndex(dstfile, "/src/"); i >= 0 {
		dstfile = dstfile[i+len("/src/"):]
	}
	dstfile = filepath.Join(*out, dstfile)

	var p c2go.Printer
	p.Print("package main\n\n")
	for _, decl := range prog.Decls {
		if decl.Span.Start.File != file {
			continue
		}
		off := len(p.Bytes())
		p.Print(decl)
		if len(p.Bytes()) > off {
			p.Print(c2go.Newline)
		}
		if err := os.MkdirAll(filepath.Dir(dstfile), 0777); err != nil {
			log.Print(err)
		}
		if err := ioutil.WriteFile(dstfile, p.Bytes(), 0666); err != nil {
			log.Print(err)
		}
	}
}
