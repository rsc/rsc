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
	src        = flag.String("src", "", "src of search")
	dst        = flag.String("dst", "", "dst of search")
	out        = flag.String("o", "/tmp/c2go", "output directory")
	strip      = flag.String("", "", "strip from input paths when writing in output directory")
	inc        = flag.String("I", "", "include directory")
	showGroups = flag.Bool("groups", false, "show groups")
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	if *inc != "" {
		cc.AddInclude(*inc)
	}
	if len(args) == 0 {
		flag.Usage()
	}
	var r []io.Reader
	files := args
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, f)
		defer f.Close()
	}
	prog, err := cc.ReadMany(files, r)
	if err != nil {
		log.Fatal(err)
	}
	inferTypes(prog)
	if *showGroups {
		return
	}
	fix(prog)
	write(prog, files)
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

func inferTypes(prog *cc.Prog) {
	add(prog)

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

	if *showGroups && *src == "" && *dst == "" {
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
		if *showGroups && *src == "" && *dst == "" {
			fmt.Printf("group(%d)\n", len(list))
		}
		var types []string
		var haveType = map[string]bool{}
		var ops []string
		var haveOp = map[string]bool{"": true}
		var oneType *cc.Type
		for _, tv := range list {
			x, ok := tv.Src.(*cc.Expr)
			if ok && x.XType != nil {
				if oneType == nil {
					oneType = x.XType
				}
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

		if len(types) > 0 {
			best := ""
			for _, typ := range types {
				switch typ {
				case "enum", "int", "uchar", "short", "int8":
					if best == "" {
						best = "int"
					}
				case "int32":
					if best != "int64" && best != "uint64" && best != "uint32" {
						best = "int32"
					}
				case "long", "longlong", "vlong":
					if best != "uint64" {
						best = "int64"
					}
				case "char[]", "char*":
					best = "string"
				}
			}
			var target *cc.Type
			for _, op := range ops {
				if op == "ptr+" && oneType != nil && oneType.Kind == cc.Ptr {
					t := *oneType
					t.Kind = c2go.Slice
					target = &t
				}
			}
			if best != "" && target == nil {
				target = &cc.Type{Kind: cc.TypedefType, Name: best}
			}
			if target != nil {
				for _, tv := range list {
					if tv.TypePtr != nil {
						typ := *tv.TypePtr
						if typ != nil && typ.Kind == cc.Func {
							typ.Base = target
							continue
						}
					}
					tv.TargetType = target
				}
			}
		}

		if *showGroups && *src == "" && *dst == "" {
			fmt.Printf("types: %s\n", strings.Join(types, ", "))
			fmt.Printf("ops: %s\n", strings.Join(ops, ", "))
			for _, tv := range list {
				fmt.Printf("\t%s (%d):", tv.Src, len(tv.Link))
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

	TargetType *cc.Type
	TypePtr    **cc.Type
}

var typeVars = map[cc.Syntax]*TypeVar{}
var typeVarsIndir = map[cc.Syntax]*TypeVar{}
var typeVarsAddr = map[cc.Syntax]*TypeVar{}

func add(x cc.Syntax) {
	var curFn *cc.Decl
	before := func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Decl:
			// Note:
			// For function declarations, the typeVar represents
			// the of the function's result (the type of a call to the function).
			// For other declarations, the typeVar represents the
			// type of the declared variable.
			if x.Storage&cc.Typedef == 0 {
				tv := &TypeVar{Src: x, TypePtr: &x.Type}
				switch x.Name {
				case "nil", "N", "L", "S", "T", "C", "...", "bval", "vval", "mpgetfix", "smprint", "namebuf":
					tv.NoLink = true
				}
				if x.Type != nil && x.Type.Is(cc.Ptr) && x.Type.Base.Is(cc.Void) {
					tv.NoLink = true
				}
				typeVars[x] = tv
			}

			if x.Storage&cc.Static != 0 && curFn != nil {
				fmt.Printf("func-static variable %s\n", x.Name)
			}

			if x.Body != nil {
				curFn = x
			}

			if x.Type != nil {
				for _, decl := range x.Type.Decls {
					decl.XOuter = x
				}
			}
		}
	}
	after := func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Decl:
			addDecl(x)

			//	if x.Type != nil && x.Type.Kind == cc.Enum {
			//		addLink(x, x.Type)
			//	}
			if x.Body != nil {
				curFn = nil
			}

		case *cc.Expr:
			addDecl(x.XDecl)
			addType(x.XType)
			addExpr(x)

		case *cc.Stmt:
			addStmt(x, curFn)

		case *cc.Type:
			//	if x.Kind == cc.Enum {
			//		typeVars[x] = &TypeVar{Src: x}
			//	}

		}
	}
	cc.Walk(x, before, after)
}

func addDecl(x *cc.Decl) {
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x, TypePtr: &x.Type}
	typeVars[x] = tv
	addType(x.Type)
}

func addType(x *cc.Type) {
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x}
	typeVars[x] = tv
	for _, d := range x.Decls {
		addDecl(d)
	}
}

func addExpr(x *cc.Expr) {
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x, TypePtr: &x.Type}
	typeVars[x] = tv

	if x.XType != nil && x.XType.Kind == cc.Ptr && x.XType.Base.Is(cc.Void) {
		tv.NoLink = true
	}

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
			switch x.XDecl.Name {
			case "andptr", "add", "offset", "siz":
				return
			}
			addLink(x, x.XDecl)
		}
	case cc.Call:
		if x.Left.Op == cc.Name {
			switch x.Left.Text {
			case "duintptr", "nodconst", "duintxx", "duint32", "duint8", "duint16", "duint64", "nodintconst", "mpmovecfix", "strcmp", "memcmp", "memmove", "strcpy", "strlen", "strncmp", "strstr", "strchr", "strrchr", "erealloc", "malloc", "emalloc", "mal", "wrint", "rdint", "mediaop":
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

	case cc.Eq, cc.AndEq, cc.XorEq:
		addLink(x, x.Left)
		addLink(x.Left, x.Right)

	case cc.EqEq, cc.NotEq, cc.And, cc.Xor:
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
			switch x.XDecl.Name {
			case "o1", "o2", "o3", "o4", "op":
				return
			}
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
	tx := typeVars[x]
	ty := typeVars[y]
	if tx == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", x.GetSpan(), x, x, x)
		panic("missing")
		return
	}
	if ty == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", y.GetSpan(), y, y, y)
		panic("missing")
		return
	}

	tx.Indir = ty
}

func addAddrLink(x, y cc.Syntax) {
	tx := typeVars[x]
	ty := typeVars[y]
	if tx == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", x.GetSpan(), x, x, x)
		panic("missing")
		return
	}
	if ty == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", y.GetSpan(), y, y, y)
		panic("missing")
		return
	}

	tx.Addr = ty
}

func addLink(x, y cc.Syntax) {
	tx := typeVars[x]
	ty := typeVars[y]
	if tx == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", x.GetSpan(), x, x, x)
		panic("missing")
		return
	}
	if ty == nil {
		fmt.Fprintf(os.Stderr, "addLink: missing typeVar for %s %T %s %p\n", y.GetSpan(), y, y, y)
		panic("missing")
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

func addStmt(x *cc.Stmt, curFn *cc.Decl) {
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

// Rewrite from C constructs to Go constructs.

func fix(x cc.Syntax) {
	cc.Preorder(x, func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Stmt:
			fixStmt(x)

		case *cc.Expr:
			switch x.Op {
			case cc.Number:
				// Rewrite char literal \0 to \x00.
				// In general we'd need to rewrite all string and char literals
				// but this is the only form that comes up.
				if x.Text == `'\0'` {
					x.Text = `'\x00'`
				}
			}

		case *cc.Decl:
			// Rewrite declaration names to avoid Go keywords.
			switch x.Name {
			case "type":
				x.Name = "typ"
			case "func":
				x.Name = "fun"
			}

			// Add file name to file-static variables to avoid conflicts.
			// TODO: Don't do this when there's no conflict?
			if x.Storage&cc.Static != 0 || x.Name != "" && x.Type != nil && x.Type.Kind == cc.Enum && !strings.Contains(x.Span.Start.File, "/include/") {
				file := filepath.Base(x.Span.Start.File)
				if i := strings.Index(file, "."); i >= 0 {
					file = file[:i]
				}
				x.Name += "_" + file
			}
		}
	})

	cc.Preorder(x, func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Stmt:
			switch x.Op {
			case cc.If, cc.For:
				x.Expr = forceBool(x.Expr)
			}

		case *cc.Expr:
			switch x.Op {
			case cc.Eq, cc.AddEq, cc.SubEq, cc.MulEq, cc.DivEq, cc.ModEq, cc.XorEq, cc.OrEq, cc.AndEq, cc.LshEq, cc.RshEq:
				if x.Left.XType != nil && x.Right.XType != nil && x.Left.XType.Def().Kind != x.Right.XType.Def().Kind {
					x.Right = &cc.Expr{Op: cc.Cast, Type: x.Left.XType, Left: x.Right}
				}

			case cc.AndAnd, cc.OrOr:
				x.Left = forceBool(x.Left)
				x.Right = forceBool(x.Right)
			}

			x.Type = fixType(x, x.Type)

		case *cc.Decl:
			x.Type = fixType(x, x.Type)

		case *cc.Type:
			x.Base = fixType(x.Base, x.Base)
		}
	})
}

func fixType(x cc.Syntax, typ *cc.Type) *cc.Type {
	tv := typeVars[x]
	if tv == nil || tv.TargetType == nil {
		if typ.String() == "char*" || typ.String() == "char[]" {
			return &cc.Type{Kind: cc.TypedefType, Name: "string"}
		}
		return typ
	}
	return tv.TargetType
}

func forceBool(x *cc.Expr) *cc.Expr {
	if x != nil && !isBoolOp(x.Op) {
		t := x.XType.Def()
		if t.Kind == cc.Ptr {
			x = &cc.Expr{Op: cc.NotEq, Left: x, Right: &cc.Expr{Op: cc.Name, Text: "nil"}}
		} else {
			x = &cc.Expr{Op: cc.NotEq, Left: x, Right: &cc.Expr{Op: cc.Name, Text: "0"}}
		}
	}
	return x
}

func isBoolOp(op cc.ExprOp) bool {
	switch op {
	case cc.Not, cc.AndAnd, cc.OrOr, cc.Lt, cc.LtEq, cc.Gt, cc.GtEq, cc.EqEq, cc.NotEq:
		return true
	}
	return false
}

func fixStmt(stmt *cc.Stmt) {
	// TODO: Double-check stmt.Labels

	switch stmt.Op {
	case cc.ARGBEGIN:
		panic(fmt.Sprintf("unexpected ARGBEGIN"))

	case cc.Do:
		// Rewrite do { ... } while(x)
		// to for(;;) { ... if(!x) break }
		// Since fixStmt is called in a preorder traversal,
		// the recursion into the children will clean up x
		// in the if condition as needed.
		stmt.Op = cc.For
		x := stmt.Expr
		stmt.Expr = nil
		stmt.Body = forceBlock(stmt.Body)
		stmt.Body.Block = append(stmt.Body.Block, &cc.Stmt{
			Op:   cc.If,
			Expr: &cc.Expr{Op: cc.Not, Left: x},
			Body: &cc.Stmt{Op: cc.Break},
		})

	case cc.While:
		stmt.Op = cc.For
		fallthrough

	case cc.For:
		before1, _ := extractSideEffects(stmt.Pre, sideStmt|sideNoAfter)
		before2, _ := extractSideEffects(stmt.Expr, sideNoAfter)
		if len(before2) > 0 {
			x := stmt.Expr
			stmt.Expr = nil
			stmt.Body = forceBlock(stmt.Body)
			top := &cc.Stmt{
				Op:   cc.If,
				Expr: &cc.Expr{Op: cc.Not, Left: x},
				Body: &cc.Stmt{Op: cc.Break},
			}
			stmt.Body.Block = append(append(before2, top), stmt.Body.Block...)
		}
		if len(before1) > 0 {
			old := copyStmt(stmt)
			stmt.Pre = nil
			stmt.Expr = nil
			stmt.Post = nil
			stmt.Body = nil
			stmt.Op = c2go.BlockNoBrace
			stmt.Block = append(before1, old)
		}
		before, after := extractSideEffects(stmt.Post, sideStmt)
		if len(before)+len(after) > 0 {
			all := append(append(before, &cc.Stmt{Op: cc.StmtExpr, Expr: stmt.Post}), after...)
			stmt.Post = &cc.Expr{Op: c2go.ExprBlock, Block: all}
		}

	case cc.If, cc.Return:
		before, _ := extractSideEffects(stmt.Expr, sideNoAfter)
		if len(before) > 0 {
			old := copyStmt(stmt)
			stmt.Expr = nil
			stmt.Body = nil
			stmt.Else = nil
			stmt.Op = c2go.BlockNoBrace
			stmt.Block = append(before, old)
		}

	case cc.StmtExpr:
		before, after := extractSideEffects(stmt.Expr, sideStmt)
		if len(before)+len(after) > 0 {
			old := copyStmt(stmt)
			stmt.Expr = nil
			stmt.Op = c2go.BlockNoBrace
			stmt.Block = append(append(before, old), after...)
		}

	case cc.Goto:
		// TODO: Figure out where the goto goes and maybe rewrite
		// to labeled break/continue.
		// Otherwise move code or something.

	case cc.Switch:
		// TODO: Change default fallthrough to default break.
	}
}

func forceBlock(x *cc.Stmt) *cc.Stmt {
	if x.Op != cc.Block {
		x = &cc.Stmt{Op: cc.Block, Block: []*cc.Stmt{x}}
	}
	return x
}

const (
	sideStmt = 1 << iota
	sideNoAfter
)

func extractSideEffects(x *cc.Expr, mode int) (before, after []*cc.Stmt) {
	doSideEffects(x, &before, &after, mode)
	return
}

func doSideEffects(x *cc.Expr, before, after *[]*cc.Stmt, mode int) {
	if x == nil {
		return
	}

	// Cannot hoist side effects from conditionally evaluated expressions
	// into unconditionally evaluated statement lists.
	// For now, detect but do not handle.
	switch x.Op {
	case cc.Cond:
		doSideEffects(x.List[0], before, after, mode&^sideStmt|sideNoAfter)
		checkNoSideEffects(x.List[1], 0)
		checkNoSideEffects(x.List[2], 0)

	case cc.AndAnd, cc.OrOr:
		doSideEffects(x.Left, before, after, mode&^sideStmt|sideNoAfter)
		checkNoSideEffects(x.Right, 0)

	default:
		doSideEffects(x.Left, before, after, mode&^sideStmt)
		doSideEffects(x.Right, before, after, mode&^sideStmt)
		for _, y := range x.List {
			doSideEffects(y, before, after, mode&^sideStmt)
		}
	}

	if mode&sideStmt != 0 {
		// Expression as statement.
		// Can leave x++ alone, can rewrite ++x to x++, can leave x [op]= y alone.
		switch x.Op {
		case cc.PreInc:
			x.Op = cc.PostInc
			return
		case cc.PreDec:
			x.Op = cc.PostDec
			return
		case cc.PostInc, cc.PostDec:
			return
		case cc.Eq, cc.AddEq, cc.SubEq, cc.MulEq, cc.DivEq, cc.ModEq, cc.XorEq, cc.OrEq, cc.AndEq, cc.LshEq, cc.RshEq:
			return
		}
	}

	switch x.Op {
	case cc.Eq, cc.AddEq, cc.SubEq, cc.MulEq, cc.DivEq, cc.ModEq, cc.XorEq, cc.OrEq, cc.AndEq, cc.LshEq, cc.RshEq:
		x.Left = forceCheap(before, x.Left)
		old := copyExpr(x)
		*before = append(*before, &cc.Stmt{Op: cc.StmtExpr, Expr: old})
		fixMerge(x, x.Left)

	case cc.PreInc, cc.PreDec:
		x.Left = forceCheap(before, x.Left)
		old := copyExpr(x)
		old.SyntaxInfo = cc.SyntaxInfo{}
		if old.Op == cc.PreInc {
			old.Op = cc.PostInc
		} else {
			old.Op = cc.PostDec
		}
		*before = append(*before, &cc.Stmt{Op: cc.StmtExpr, Expr: old})
		fixMerge(x, x.Left)

	case cc.PostInc, cc.PostDec:
		x.Left = forceCheap(before, x.Left)
		if mode&sideNoAfter != 0 {
			// Not allowed to generate fixups afterward.
			d := &cc.Decl{
				Name: "tmp",
				Type: x.XType,
				Init: &cc.Init{Expr: x.Left},
			}
			old := copyExpr(x.Left)
			old.SyntaxInfo = cc.SyntaxInfo{}
			*before = append(*before,
				&cc.Stmt{Op: cc.StmtDecl, Decl: d},
				&cc.Stmt{Op: cc.StmtExpr, Expr: &cc.Expr{Op: x.Op, Left: old}},
			)
			x.Op = cc.Name
			x.Text = d.Name
			x.XDecl = d
			x.Left = nil
			break
		}
		old := copyExpr(x)
		old.SyntaxInfo = cc.SyntaxInfo{}
		*after = append(*after, &cc.Stmt{Op: cc.StmtExpr, Expr: old})
		fixMerge(x, x.Left)

	case cc.Cond:
		// Rewrite c ? y : z into tmp with initialization:
		//	var tmp typeof(c?y:z)
		//	if c {
		//		tmp = y
		//	} else {
		//		tmp = z
		//	}
		d := &cc.Decl{
			Name: "tmp",
			Type: x.XType,
		}
		*before = append(*before,
			&cc.Stmt{Op: cc.StmtDecl, Decl: d},
			&cc.Stmt{Op: cc.If, Expr: x.List[0],
				Body: &cc.Stmt{
					Op: cc.StmtExpr,
					Expr: &cc.Expr{
						Op:    cc.Eq,
						Left:  &cc.Expr{Op: cc.Name, Text: d.Name, XDecl: d},
						Right: x.List[1],
					},
				},
				Else: &cc.Stmt{
					Op: cc.StmtExpr,
					Expr: &cc.Expr{
						Op:    cc.Eq,
						Left:  &cc.Expr{Op: cc.Name, Text: d.Name, XDecl: d},
						Right: x.List[2],
					},
				},
			},
		)
		x.Op = cc.Name
		x.Text = d.Name
		x.XDecl = d
		x.List = nil
	}
}

func copyExpr(x *cc.Expr) *cc.Expr {
	old := *x
	old.SyntaxInfo = cc.SyntaxInfo{}
	return &old
}

func copyStmt(x *cc.Stmt) *cc.Stmt {
	old := *x
	old.SyntaxInfo = cc.SyntaxInfo{}
	old.Labels = nil
	return &old
}

func forceCheap(before *[]*cc.Stmt, x *cc.Expr) *cc.Expr {
	// TODO
	return x
}

func fixMerge(dst, src *cc.Expr) {
	syn := dst.SyntaxInfo
	syn.Comments.Before = append(syn.Comments.Before, src.Comments.Before...)
	syn.Comments.After = append(syn.Comments.After, src.Comments.After...)
	syn.Comments.Suffix = append(syn.Comments.Suffix, src.Comments.Suffix...)
	*dst = *src
	dst.SyntaxInfo = syn
}

func checkNoSideEffects(x *cc.Expr, mode int) {
	var before, after []*cc.Stmt
	old := x.String()
	doSideEffects(x, &before, &after, mode)
	if len(before)+len(after) > 0 {
		fmt.Printf("cannot handle side effects in %s\n", old)
	}
}

func write(prog *cc.Prog, files []string) {
	for _, file := range files {
		writeFile(prog, file, "")
	}
	writeFile(prog, "/Users/rsc/g/go/include/fmt.h", "liblink/fmt_h.go")
	writeFile(prog, "/Users/rsc/g/go/include/bio.h", "liblink/bio_h.go")
	writeFile(prog, "/Users/rsc/g/go/include/link.h", "liblink/link_h.go")
	writeFile(prog, "/Users/rsc/g/go/src/cmd/5l/5.out.h", "liblink/5.out.go")
	writeFile(prog, "/Users/rsc/g/go/src/cmd/6l/6.out.h", "liblink/6.out.go")
	writeFile(prog, "/Users/rsc/g/go/src/cmd/8l/8.out.h", "liblink/8.out.go")

	ioutil.WriteFile(filepath.Join(*out, "liblink/zzz.go"), []byte(zzzExtra), 0666)
}

var zzzExtra = `
package main

type Rune rune
type va_list struct{}
`

func writeFile(prog *cc.Prog, file, dstfile string) {
	if dstfile == "" {
		dstfile = strings.TrimSuffix(strings.TrimSuffix(file, ".c"), ".h") + ".go"
		if *strip != "" {
			dstfile = strings.TrimPrefix(dstfile, *strip)
		} else if i := strings.LastIndex(dstfile, "/src/"); i >= 0 {
			dstfile = dstfile[i+len("/src/"):]
		}
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
