// +build ignore

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"code.google.com/p/rsc/c2go"
	"code.google.com/p/rsc/cc"
)

type byLen [][]*TypeVar

func (x byLen) Len() int      { return len(x) }
func (x byLen) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byLen) Less(i, j int) bool {
	if len(x[i]) != len(x[j]) {
		return len(x[i]) < len(x[j])
	}
	return fmt.Sprint(x[i][0].Src) < fmt.Sprint(x[j][0].Src)
}

type TypeGroup struct {
	Vars       []*TypeVar
	Types      []string
	Ops        []string
	Target     *cc.Type
	TargetKind cc.TypeKind
	Bool       bool // bool where 0=false, 1=true
	NegBool    bool // bool where -1=false, 0=true
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
	}
	for _, tv := range typeVars {
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
		var types []string
		var haveType = map[string]bool{}
		var ops []string
		var haveOp = map[string]bool{"": true}
		var bestType *cc.Type
		g := &TypeGroup{
			Vars: list,
		}
		for _, tv := range list {
			tv.Group = g
			x, ok := tv.Src.(*cc.Expr)
			if ok && x.XType != nil {
				t := x.XType
				if t.Kind == cc.Enum {
					tt := *t
					t = &tt
					t.Kind = cc.Int
				}
				if bestType == nil || t.Kind > bestType.Kind {
					bestType = x.XType
				}
				switch t := x.XType; t.Kind {
				case cc.Array, cc.Ptr:
					if t.Width == nil && t.Base != nil && t.Base.Kind == cc.Char {
						bestType = &cc.Type{Kind: c2go.String}
					}
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
		g.Ops = ops
		g.Types = types

		typeGroups = append(typeGroups, g)
	}

	determineBools()

	if *showGroups && *src == "" && *dst == "" {
		for _, g := range typeGroups {
			if len(g.Vars) <= 1 {
				continue
			}
			fmt.Printf("group(%d) %p\n", len(g.Vars), g)
			if g.Bool {
				fmt.Printf("real type: bool\n")
			}
			if g.NegBool {
				fmt.Printf("real type: negbool\n")
			}
			fmt.Printf("types: %s\n", strings.Join(g.Types, ", "))
			fmt.Printf("ops: %s\n", strings.Join(g.Ops, ", "))
			for _, tv := range g.Vars {
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

func determineBools() {
	for _, g := range typeGroups {
		g.Bool = true
	}
	changed := true
	for changed {
		changed = false
	Groups:
		for _, g := range typeGroups {
			if !g.Bool {
				continue
			}
			hasBool := false
			for _, typ := range g.Types {
				if strings.Contains(typ, "*") {
					g.Bool = false
					changed = true
					continue Groups
				}
				if typ == "bool" {
					hasBool = true
				}
			}
			if !hasBool {
				g.Bool = false
				changed = true
				continue
			}
			for _, tv := range g.Vars {
				if tv.UsedAsBool {
					continue
				}
				switch x := tv.Src.(type) {
				case *cc.Decl:
					typ := x.Type
					if typ != nil && typ.Kind == cc.Func {
						typ = typ.Base
					}
					if typ != nil && cc.Char <= typ.Kind && typ.Kind <= cc.Enum {
						continue
					}
				case *cc.Expr:
					switch x.Op {
					case cc.AndAnd, cc.OrOr, cc.Not, cc.EqEq, cc.NotEq, cc.Gt, cc.Lt, cc.GtEq, cc.LtEq:
						continue
					case cc.Number:
						if x.Text == "0" || x.Text == "1" {
							continue
						}
					case cc.Call:
						tv := typeVars[x.Left.XDecl]
						if tv != nil && tv.Group != nil && tv.Group.Bool {
							continue
						}
					case cc.Name, cc.Eq, cc.Paren:
						continue
					}
				case *cc.Type:
					if cc.Char <= x.Kind && x.Kind <= cc.Enum {
						continue
					}
				default:
					fmt.Printf("unexpected syntax in group: %T\n", tv.Src)
				}
				fmt.Printf("not bool %p: %v\n", g, tv.Src)
				g.Bool = false
				changed = true
			}
		}
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
	Group      *TypeGroup
	Return     *cc.Decl

	parent *TypeVar

	Type *cc.Type
}

func (tv *TypeVar) Name() string {
	if x, ok := tv.Src.(*cc.Decl); ok {
		return x.Name
	}
	return ""
}

var typeGroups []*TypeGroup
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
				tv := &TypeVar{Src: x, Type: x.Type}
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

			x.CurFn = curFn

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
			addType(&x.XType)
			addExpr(x)

		case *cc.Stmt:
			addStmt(x, curFn)
		}
	}
	cc.Walk(x, before, after)
}

func addDecl(x *cc.Decl) {
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x, Type: x.Type}
	typeVars[x] = tv
	addType(&x.Type)
}

func addType(p **cc.Type) {
	x := *p
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x, Type: x}
	typeVars[x] = tv
	for _, d := range x.Decls {
		addDecl(d)
	}
}

func addExpr(x *cc.Expr) {
	if x == nil || typeVars[x] != nil {
		return
	}
	tv := &TypeVar{Src: x, Type: x.Type}
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
			typeVars[x.Expr].Return = curFn
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
