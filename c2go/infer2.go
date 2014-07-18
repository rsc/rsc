// +build ignore

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"code.google.com/p/rsc/cc"
)

// A flowSyntax is a node representing the flow of a value
// through a piece of syntax in the program.
// The syntax can be a *cc.Expr or a *cc.Decl.
type flowSyntax struct {
	syntax      cc.Syntax     // original syntax
	adj         []*flowSyntax // adjacent syntax connected by assignment
	returnValue *flowSyntax
	isStmtExpr  bool
	usedAsBool  bool
	group       *flowGroup
	stopFlow    bool
	ptrAdd      bool
	ptrIndex    bool
}

type flowGroup struct {
	syntax    []*flowSyntax
	decls     []*cc.Decl
	exprs     []*cc.Expr
	goType    *cc.Type
	goKind    cc.TypeKind
	goKey     string
	goFlow    *flowSyntax
	canon     *cc.Type
	canonDecl *cc.Decl
}

var flowCache = map[cc.Syntax]*flowSyntax{}
var flowGroups []*flowGroup

func inferTypes(prog *cc.Prog) {
	cc.Postorder(prog, func(x cc.Syntax) {
		if t, ok := x.(*cc.Type); ok {
			if t.Kind == cc.Struct || t.Kind == cc.Enum {
				for _, d := range t.Decls {
					d.OuterType = t
				}
			}
		}
	})

	addFlow(prog)

	for _, f := range flowCache {
		if f.group != nil || f.stopFlow {
			continue
		}
		g := &flowGroup{}
		exploreGroup(g, f)
		if len(g.decls) == 0 {
			continue
		}
		flowGroups = append(flowGroups, g)
	}
	sort.Sort(flowGroupsBySize(flowGroups))

	for _, g := range flowGroups {
		var typ *cc.Type
		var typDecl *cc.Decl
		for _, d := range g.decls {
			if d.Type == nil {
				continue
			}
			dt := d.Type.Def()
			if typ == nil || typ.Kind == cc.Ptr && typ.Base.Is(cc.Void) {
				typ = dt
				typDecl = d
			}
			if !inferCompatible(dt, typ) {
				fmt.Printf("BAD INFER: mixing %v (%v) and %v (%v)\n", typ, declKey(typDecl), d.Type, declKey(d))
				findFlowPath(flowCache[typDecl], flowCache[d])
				os.Exit(1)
			}
			if isNumericCType(typ) && isNumericCType(dt) && typ.Kind == cc.Enum || dt.Kind != cc.Enum && typ.Kind < dt.Kind {
				typ = dt
				typDecl = d
			}
		}
		g.canon = typ
		g.canonDecl = typDecl
	}

	if *src != "" && *dst != "" {
		var fsrc, fdst *flowSyntax
		for _, f := range flowCache {
			d, ok := f.syntax.(*cc.Decl)
			if ok {
				key := declKey(d)
				if key == "" {
					continue
				}
				if *src == key {
					fsrc = f
					fmt.Printf("%s in %p %p\n", key, f, f.group)
				}
				if strings.HasSuffix(*src, key) {
					fmt.Printf("near: %s\n", key)
				}
				if *dst == key {
					fdst = f
					fmt.Printf("%s in %p %p\n", key, f, f.group)
				}
				if strings.HasSuffix(*dst, key) {
					fmt.Printf("near: %s\n", key)
				}
			}
		}
		if fsrc != nil && fdst != nil {
			findFlowPath(fsrc, fdst)
			os.Exit(0)
		}
		fmt.Printf("%s and %s are not in the same group\n", *src, *dst)
		os.Exit(0)
	}
}

func inferCompatible(t1, t2 *cc.Type) bool {
	t1 = t1.Def()
	t2 = t2.Def()
	if isNumericCType(t1) && isNumericCType(t2) {
		return true
	}
	if t1.Kind == cc.Ptr && t1.Base.Kind == cc.Func {
		t1 = t1.Base
	}
	if t2.Kind == cc.Ptr && t2.Base.Kind == cc.Func {
		t2 = t2.Base
	}
	if sameType(t1, t2) {
		return true
	}
	if t1.Kind > t2.Kind {
		t1, t2 = t2, t1
	}
	if t1.Kind == cc.Ptr && t2.Kind == cc.Array && (t1.Base.Is(cc.Void) || sameType(t1.Base, t2.Base)) {
		return true
	}
	if t1.Kind == cc.Ptr && t2.Kind == cc.Ptr && (t1.Base.Is(cc.Void) || t2.Base.Is(cc.Void)) {
		return true
	}
	return false
}

func isNumericCType(t *cc.Type) bool {
	return t != nil && cc.Char <= t.Kind && t.Kind <= cc.Enum
}

type flowGroupsBySize []*flowGroup

func (x flowGroupsBySize) Len() int           { return len(x) }
func (x flowGroupsBySize) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x flowGroupsBySize) Less(i, j int) bool { return len(x[i].decls) > len(x[j].decls) }

func exploreGroup(g *flowGroup, f *flowSyntax) {
	if f == nil || f.group == g || f.stopFlow {
		return
	}
	if f.group != nil {
		panic("mixed groups")
	}
	f.group = g
	g.syntax = append(g.syntax, f)
	switch x := f.syntax.(type) {
	case *cc.Decl:
		g.decls = append(g.decls, x)
	case *cc.Expr:
		g.exprs = append(g.exprs, x)
	default:
		panic(fmt.Sprintf("unexpected syntax %T", x))
	}
	for _, ff := range f.adj {
		exploreGroup(g, ff)
	}
}

func addFlow(prog *cc.Prog) {
	for _, d := range prog.Decls {
		addFlowDecl(nil, d)
	}

	// Mop up the rest.
	cc.Preorder(prog, func(x cc.Syntax) {
		if d, ok := x.(*cc.Decl); ok {
			addFlowDecl(nil, d)
		}
	})
}

func addFlowDecl(curfn, d *cc.Decl) {
	if d == nil || d.Type == nil || flowCache[d] != nil {
		return
	}

	f := &flowSyntax{syntax: d}
	flowCache[d] = f
	if d.Type.IsPtrVoid() {
		f.stopFlow = true
	}

	switch d.Name {
	case "nil":
		f.stopFlow = true
	}
	switch declKey(d) {
	case "Link.instoffset":
		f.stopFlow = true
	}

	if d.Type.Kind == cc.Func {
		rv := &flowSyntax{syntax: &cc.Decl{Name: "return", Type: d.Type.Base, CurFn: d}}
		if d.Type.Base.IsPtrVoid() {
			rv.stopFlow = true
		}
		f.returnValue = rv
		for _, dd := range d.Type.Decls {
			dd.CurFn = d
		}
	}

	for _, dd := range d.Type.Decls {
		addFlowDecl(curfn, dd)
	}

	if d.Init != nil {
		addFlowInit(d, d.Init)
	}
	if d.Body != nil {
		addFlowStmt(d, d.Body)
	}
}

func addFlowInit(d *cc.Decl, init *cc.Init) {
	if init == nil {
		return
	}

	addFlowExpr(nil, init.Expr)

	last := d
	for _, pre := range init.Prefix {
		last = pre.XDecl
		addFlowDecl(nil, last)
		addFlowExpr(nil, pre.Index)
	}
	if init.Expr != nil && last != nil {
		addFlowExpr(nil, init.Expr)
		addFlowEdge(flowCache[init.Expr], flowCache[last])
	}

	typ := init.XType
	for i, br := range init.Braced {
		var field *cc.Decl
		if typ != nil && i < len(typ.Decls) {
			field = typ.Decls[i]
			addFlowDecl(nil, field)
		}
		addFlowInit(field, br)
	}
}

func addFlowStmt(curfn *cc.Decl, x *cc.Stmt) {
	if x == nil {
		return
	}
	addFlowExpr(curfn, x.Pre)
	addFlowExpr(curfn, x.Post)
	addFlowExpr(curfn, x.Expr)
	addFlowDecl(curfn, x.Decl)
	addFlowStmt(curfn, x.Else)
	addFlowStmt(curfn, x.Body)
	for _, stmt := range x.Block {
		addFlowStmt(curfn, stmt)
	}
	for _, lab := range x.Labels {
		addFlowExpr(curfn, lab.Expr)
	}

	switch x.Op {
	case cc.Return:
		if x.Expr != nil {
			f := flowCache[x.Expr]
			addFlowEdge(f, flowCache[curfn].returnValue)
		}
		return

	case cc.StmtDecl:
		x.Decl.CurFn = curfn

	case cc.StmtExpr:
		flowCache[x.Expr].isStmtExpr = true

	case cc.For:
		if x.Pre != nil {
			flowCache[x.Pre].isStmtExpr = true
		}
		if x.Post != nil {
			flowCache[x.Post].isStmtExpr = true
		}
		if x.Expr != nil {
			flowCache[x.Expr].usedAsBool = true
		}

	case cc.If, cc.Do, cc.While:
		if x.Expr != nil {
			flowCache[x.Expr].usedAsBool = true
		}

	case cc.Switch:
		f := flowCache[x.Expr]
		for _, stmt := range x.Body.Block {
			for _, lab := range stmt.Labels {
				if lab.Op == cc.Case && lab.Expr != nil {
					addFlowExpr(curfn, lab.Expr)
					addFlowEdge(f, flowCache[lab.Expr])
				}
			}
		}
	}
}

func addFlowExpr(curfn *cc.Decl, x *cc.Expr) {
	if x == nil || flowCache[x] != nil {
		return
	}
	f := &flowSyntax{syntax: x}
	flowCache[x] = f

	addFlowExpr(curfn, x.Left)
	addFlowExpr(curfn, x.Right)
	for _, expr := range x.List {
		addFlowExpr(curfn, expr)
	}
	for _, stmt := range x.Block {
		addFlowStmt(curfn, stmt)
	}
	addFlowInit(nil, x.Init)
	addFlowDecl(curfn, x.XDecl)

	switch x.Op {
	case cc.Add, cc.Sub:
		if x.XType.Is(cc.Ptr) {
			if x.Left.XType.Is(cc.Ptr) {
				f.ptrAdd = true
				addFlowEdge(f, flowCache[x.Left])
			}
			if x.Right.XType.Is(cc.Ptr) {
				f.ptrAdd = true
				addFlowEdge(f, flowCache[x.Right])
			}
		} else if x.Op == cc.Sub && x.Left.XType.Is(cc.Ptr) && x.Right.XType.Is(cc.Ptr) {
			f1 := flowCache[x.Left]
			f2 := flowCache[x.Right]
			f1.ptrAdd = true
			f2.ptrAdd = true
			addFlowEdge(f1, f2)
		}

	case cc.Lt, cc.LtEq, cc.Gt, cc.GtEq:
		if x.Left.XType.Is(cc.Ptr) {
			f1 := flowCache[x.Left]
			f2 := flowCache[x.Right]
			f1.ptrAdd = true
			f2.ptrAdd = true
			addFlowEdge(f1, f2)
		}
	}

	switch x.Op {
	case cc.AddEq, cc.PostDec, cc.PostInc, cc.PreDec, cc.PreInc:
		// no flow to right - may be pointer += int
		addFlowEdge(f, flowCache[x.Left])
		if x.Left.XType.Is(cc.Ptr) {
			f.ptrAdd = true
		}

	case cc.AndAnd, cc.OrOr, cc.Not:
		flowCache[x.Left].usedAsBool = true

	case cc.AndEq, cc.OrEq, cc.XorEq, cc.MulEq, cc.DivEq:
		// no flow to right - a bit too fussy
		addFlowEdge(f, flowCache[x.Left])

	case cc.Arrow, cc.Dot:
		if x.XDecl != nil {
			addFlowEdge(f, flowCache[x.XDecl])
		}

	case cc.Call:
		if x.Left.Op == cc.Name && x.Left.XDecl != nil {
			d := x.Left.XDecl
			addFlowDecl(nil, d)
			if fd := flowCache[d]; fd != nil && fd.returnValue != nil {
				addFlowEdge(f, fd.returnValue)
			}
			for i := 0; i < len(d.Type.Decls) && i < len(x.List); i++ {
				dd := d.Type.Decls[i]
				if dd.Type != nil {
					addFlowDecl(nil, dd)
					addFlowEdge(flowCache[dd], flowCache[x.List[i]])
				}
			}
		}

	case cc.Comma:
		if len(x.List) > 0 {
			addFlowEdge(f, flowCache[x.List[len(x.List)-1]])
		}

	case cc.Cond:
		flowCache[x.List[0]].usedAsBool = true
		addFlowEdge(f, flowCache[x.List[1]])
		addFlowEdge(f, flowCache[x.List[2]])

	case cc.Eq:
		// flow to left and right
		addFlowEdge(f, flowCache[x.Left])
		addFlowEdge(f, flowCache[x.Right])

	case cc.EqEq, cc.NotEq, cc.Gt, cc.GtEq, cc.Lt, cc.LtEq:
		addFlowEdge(flowCache[x.Left], flowCache[x.Right])

	case cc.Index:
		if x.Left.XType.Is(cc.Ptr) {
			ff := flowCache[x.Left]
			ff.ptrAdd = true
		}

	case cc.Minus, cc.Plus, cc.Twid:
		addFlowEdge(f, flowCache[x.Left])

	case cc.Name:
		if x.XDecl != nil {
			addFlowEdge(f, flowCache[x.XDecl])
		}
	}
}

func addFlowEdge(f, g *flowSyntax) {
	f.adj = append(f.adj, g)
	g.adj = append(g.adj, f)
}

func findFlowPath(src, dst *flowSyntax) {
	next := map[*flowSyntax]*flowSyntax{dst: dst}
	q := []*flowSyntax{dst}
Search:
	for i := 0; i < len(q); i++ {
		f := q[i]
		for _, ff := range f.adj {
			if ff.group == f.group && next[ff] == nil {
				next[ff] = f
				q = append(q, ff)
				if ff == src {
					break Search
				}
			}
		}
	}
	if next[src] == nil {
		fmt.Printf("no path from %s to %s\n", src.syntax, dst.syntax)
		return
	}
	for f := src; ; f = next[f] {
		fmt.Printf("%s %s [stop=%v]\n", f.syntax.GetSpan(), f.syntax, f.stopFlow)
		if f == dst {
			break
		}
	}
}
