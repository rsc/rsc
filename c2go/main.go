// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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
	forceTypes(prog)
	inferTypes(prog)
	if *showGroups {
		return
	}
	rewriteSyntax(prog)
	rewriteTypes(prog)
	fixGoTypes(prog)
	write(prog, files)
}

func forceTypes(x cc.Syntax) {
	cc.Postorder(x, func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Decl:
			switch x.Name {
			case "o1", "h":
				x.Type = &cc.Type{Kind: cc.Ulong}
			case "oprrr", "opbra", "olr", "olhr", "olrr", "olhrr", "osr", "oshr", "ofsr", "osrr", "oshrr", "omvl", "ocmp":
				x.Type.Base = &cc.Type{Kind: cc.Ulong}
			case "out":
				if x.Type != nil && x.Type.Base != nil && (x.Type.Base.Kind == cc.Long || x.Type.Base.Kind == cc.Int) {
					x.Type.Base.Kind++
				}
			}

		case *cc.Expr:
			switch x.Op {
			case cc.Name:
				x.XType = x.XDecl.Type
			}
		}
	})
}

// Rewrite C types to be Go types.
func rewriteTypes(x cc.Syntax) {
	cache := make(map[*cc.Type]*cc.Type)
	cc.Postorder(x, func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Type:
			if len(x.Decls) > 0 {
				last := x.Decls[len(x.Decls)-1]
				if last.Name == "" && last.Type.Is(cc.Void) {
					x.Decls = x.Decls[:len(x.Decls)-1]
				}
			}
		}
	})

	cc.Preorder(x, func(x cc.Syntax) {
		switch x := x.(type) {
		case *cc.Decl:
			x.Type = toGoType(x, x.Type, cache)
			if x.Type != nil && x.Type.Kind == cc.Func && !x.Type.Base.Is(cc.Void) {
				x.Type.Base = toGoType(x, x.Type.Base, cache)
			}
			if x.Name == "andptr" {
				x.Type.Kind = c2go.Slice
			}

		case *cc.Expr:
			x.XType = toGoType(x, x.XType, cache)
			x.Type = toGoType(x, x.Type, cache)
		}
	})
}

var c2goKind = map[cc.TypeKind]cc.TypeKind{
	cc.Char:      c2go.Int8,
	cc.Uchar:     c2go.Uint8,
	cc.Short:     c2go.Int16,
	cc.Ushort:    c2go.Uint16,
	cc.Int:       c2go.Int,
	cc.Uint:      c2go.Uint,
	cc.Long:      c2go.Int32,
	cc.Ulong:     c2go.Uint32,
	cc.Longlong:  c2go.Int64,
	cc.Ulonglong: c2go.Uint64,
	cc.Float:     c2go.Float32,
	cc.Double:    c2go.Float64,
	cc.Enum:      c2go.Int,
}

func toGoType(x cc.Syntax, typ *cc.Type, cache map[*cc.Type]*cc.Type) (r *cc.Type) {
	if typ == nil {
		return nil
	}

	// Look in cache first. This cuts off recursion for self-referential types.
	// The cache only contains aggregate types - numeric types are shared
	// by many expressions in the program and we might want to translate
	// them differently in different contexts.
	if cache[typ] != nil {
		return cache[typ]
	}

	tv := typeVars[x]

	switch typ.Kind {
	default:
		panic(fmt.Sprintf("unexpected C type %s", typ))

	case cc.Void:
		return &cc.Type{Kind: cc.Struct} // struct{}

	case cc.Char, cc.Uchar, cc.Short, cc.Ushort, cc.Int, cc.Uint, cc.Long, cc.Ulong, cc.Longlong, cc.Ulonglong, cc.Float, cc.Double, cc.Enum:
		// Use type group to decide.
		var t *cc.Type
		if tv != nil && tv.Group != nil {
			g := tv.Group
			if g.Target == nil {
				var kind cc.TypeKind
				for _, tv := range g.Vars {
					if tv.Type == nil {
						continue
					}
					if cc.Char <= tv.Type.Kind && tv.Type.Kind <= cc.Enum {
						if k := c2goKind[tv.Type.Kind]; kind < k {
							kind = k
						}
					}
				}
				if kind != 0 {
					g.Target = &cc.Type{Kind: kind}
				}
			}
			if g.Target != nil {
				t = g.Target
			}
		}
		if t == nil {
			t = &cc.Type{Kind: c2goKind[typ.Kind]}
		}
		if typ.Decls != nil {
			tt := *t
			t = &tt
			t.Decls = typ.Decls
		}
		return t

	case cc.Ptr:
		t := &cc.Type{Kind: cc.Ptr}
		cache[typ] = t

		if typ.Base.Kind == cc.Char {
			t.Kind = c2go.String
			return t
		}

		// Use type group to decide slice vs string vs ptr, if available.
		if tv != nil && tv.Group != nil {
			g := tv.Group
			if g.Target != nil {
				t = g.Target
				cache[typ] = t
				return t
			}

		PtrOps:
			for _, op := range tv.Group.Ops {
				switch op {
				case "ptr+", "ptr++", "[i]":
					t.Kind = c2go.Slice
					break PtrOps
				}
			}
		}

		t.Base = toGoType(typ.Base, typ.Base, cache)
		return t

	case cc.Struct, cc.Func:
		// For structs or funcs, and we rewrite the Decls in place.
		cache[typ] = typ
		//if typ.Kind == cc.Func && !isDecl(x) && !typ.Base.Is(cc.Void) {
		//	typ.Base = toGoType(typ.Base, typ.Base, cache)
		//	fmt.Printf("now %s\n", typ)
		//}
		return typ

	case cc.Array:
		t := &cc.Type{Kind: cc.Array, Width: typ.Width}
		if t.Width == nil {
			t.Kind = c2go.Slice
		}
		cache[typ] = t
		t.Base = toGoType(typ.Base, typ.Base, cache)
		return t

	case cc.TypedefType:
		k := typ.Base.Kind
		if cc.Void <= k && k <= cc.Enum {
			return toGoType(x, typ.Base, cache)
		}
		t := &cc.Type{Kind: cc.TypedefType, Name: typ.Name}
		cache[typ] = t
		t.Base = toGoType(typ.Base, typ.Base, cache)
		return t
	}
}

func isDecl(x cc.Syntax) bool {
	_, ok := x.(*cc.Decl)
	return ok
}

// fixGoTypes fixes all the Go type mismatches.
func fixGoTypes(prog *cc.Prog) {
	did := make(map[*cc.Decl]bool)
	for _, decl := range prog.Decls {
		if did[decl] {
			continue
		}
		did[decl] = true
		if decl.Init != nil {
			fixGoTypesInit(decl, decl.Init)
		}
		if decl.Body != nil {
			fixGoTypesStmt(decl, decl.Body)
		}
	}
}

func fixGoTypesInit(decl *cc.Decl, x *cc.Init) {
	if x.Expr != nil {
		fixGoTypesExpr(nil, x.Expr, x.XType)
	}
	for _, init := range x.Braced {
		fixGoTypesInit(decl, init)
	}
}

func fixGoTypesStmt(fn *cc.Decl, x *cc.Stmt) {
	if x == nil {
		return
	}

	switch x.Op {
	case cc.StmtDecl, cc.StmtExpr:
		fixGoTypesExpr(fn, x.Expr, nil)

	case cc.If, cc.For:
		fixGoTypesExpr(fn, x.Pre, nil)
		fixGoTypesExpr(fn, x.Post, nil)
		fixGoTypesExpr(fn, x.Expr, boolType)

	case cc.Return:
		if x.Expr != nil {
			forceGoType(fn, x.Expr, fn.Type.Base)
		}
	}
	for _, stmt := range x.Block {
		fixGoTypesStmt(fn, stmt)
	}
	if len(x.Block) > 0 && x.Body != nil {
		panic("block and body")
	}
	fixGoTypesStmt(fn, x.Body)
	fixGoTypesStmt(fn, x.Else)

	for _, lab := range x.Labels {
		// TODO: use correct type
		fixGoTypesExpr(fn, lab.Expr, nil)
	}
}

func zeroFor(targ *cc.Type) *cc.Expr {
	if targ != nil {
		switch targ.Def().Kind {
		case c2go.String:
			return &cc.Expr{Op: cc.String, Texts: []string{`""`}}

		case c2go.Slice, cc.Ptr:
			return &cc.Expr{Op: cc.Name, Text: "nil"}

		case cc.Struct, cc.Array:
			return &cc.Expr{Op: cc.CastInit, Type: targ, Init: &cc.Init{}}
		}
		return &cc.Expr{Op: cc.Number, Text: "0 /*" + targ.String() + "*/"}
	}

	return &cc.Expr{Op: cc.Number, Text: "0 /*untyped*/"}
}

func fixGoTypesExpr(fn *cc.Decl, x *cc.Expr, targ *cc.Type) (ret *cc.Type) {

	if x == nil {
		return nil
	}

	defer func() {
		x.XType = ret
	}()

	if x.Op == cc.Paren {
		return fixGoTypesExpr(fn, x.Left, targ)
	}

	// Make explicit C's implicit conversions from boolean to non-boolean and vice versa.
	switch x.Op {
	case cc.AndAnd, cc.OrOr, cc.Not, cc.EqEq, cc.Lt, cc.LtEq, cc.Gt, cc.GtEq, cc.NotEq:
		if targ != nil && targ.Kind != c2go.Bool {
			old := copyExpr(x)
			x.Op = cc.Call
			x.Left = &cc.Expr{Op: cc.Name, Text: targ.String()}
			x.Right = old
			fixGoTypesExpr(fn, old, boolType)
			return targ
		}
	default:
		if targ != nil && targ.Kind == c2go.Bool {
			old := copyExpr(x)
			x.Op = cc.NotEq
			x.Left = old
			left := fixGoTypesExpr(fn, old, nil)
			x.Right = zeroFor(left)
			return targ
		}
	}

	switch x.Op {
	default:
		panic(fmt.Sprintf("unexpected construct %v in fixGoTypesExpr - %v - %v", x, x.Op, x.Span))

	case cc.Comma:
		for i, y := range x.List {
			t := targ
			if i+1 < len(x.List) {
				t = nil
			}
			fixGoTypesExpr(fn, y, t)
		}
		return nil

	case c2go.ExprBlock:
		for _, stmt := range x.Block {
			fixGoTypesStmt(fn, stmt)
		}
		return nil

	case cc.Add, cc.And, cc.Div, cc.Mod, cc.Mul, cc.Or, cc.Sub, cc.Xor:
		left := fixGoTypesExpr(fn, x.Left, targ)

		if x.Op == cc.And && x.Right.Op == cc.Twid {
			x.Op = c2go.AndNot
			x.Right = x.Right.Left
		}

		right := fixGoTypesExpr(fn, x.Right, targ)

		if x.Op == cc.Add && isSliceOrString(left) {
			x.Op = c2go.ExprSlice
			x.List = []*cc.Expr{x.Left, x.Right, nil}
			x.Left = nil
			x.Right = nil
			return left
		}

		return fixBinary(fn, x, left, right, targ)

	case cc.AddEq, cc.AndEq, cc.DivEq, cc.Eq, cc.ModEq, cc.MulEq, cc.OrEq, cc.SubEq, cc.XorEq:
		left := fixGoTypesExpr(fn, x.Left, nil)

		if x.Op == cc.AndEq && x.Right.Op == cc.Twid {
			x.Op = c2go.AndNotEq
			x.Right = x.Right.Left
		}

		forceGoType(fn, x.Right, left)

		if x.Op == cc.AddEq && isSliceOrString(left) {
			old := copyExpr(x.Left)
			x.Op = cc.Eq
			x.Right = &cc.Expr{Op: c2go.ExprSlice, List: []*cc.Expr{old, x.Right, nil}}
		}

		return left

	case cc.Addr:
		fixGoTypesExpr(fn, x.Left, nil)
		return nil

	case cc.AndAnd, cc.OrOr, cc.Not:
		fixGoTypesExpr(fn, x.Left, boolType)
		if x.Right != nil {
			fixGoTypesExpr(fn, x.Right, boolType)
		}
		return boolType

	case cc.Arrow, cc.Dot:
		left := fixGoTypesExpr(fn, x.Left, nil)

		if x.Op == cc.Arrow && isSliceOrString(left) {
			x.Left = &cc.Expr{Op: cc.Index, Left: x.Left, Right: &cc.Expr{Op: cc.Number, Text: "0"}}
		}

		return x.XDecl.Type

	case cc.Call:
		if fixSpecialCall(fn, x) {
			return x.XType
		}
		left := fixGoTypesExpr(fn, x.Left, nil)
		for i, y := range x.List {
			if left != nil && left.Kind == cc.Func && i < len(left.Decls) {
				forceGoType(fn, y, left.Decls[i].Type)
			} else {
				fixGoTypesExpr(fn, y, nil)
			}
		}
		if left != nil && left.Kind == cc.Func {
			return left.Base
		}
		return nil

	case cc.Cast:
		fixGoTypesExpr(fn, x.Left, nil)
		return x.Type

	case cc.CastInit:
		fixGoTypesInit(nil, x.Init)
		return x.Type

	case cc.EqEq, cc.Gt, cc.GtEq, cc.Lt, cc.LtEq, cc.NotEq:
		left := fixGoTypesExpr(fn, x.Left, nil)
		right := fixGoTypesExpr(fn, x.Right, nil)
		fixBinary(fn, x, left, right, nil)
		return boolType

	case cc.Index, cc.Indir:
		left := fixGoTypesExpr(fn, x.Left, nil)
		if x.Right != nil {
			fixGoTypesExpr(fn, x.Right, nil)
		}
		if left == nil {
			return nil
		}

		if isSliceOrString(left) && x.Op == cc.Indir {
			x.Op = cc.Index
			x.Right = &cc.Expr{Op: cc.Number, Text: "0"}
		}

		switch left.Kind {
		case c2go.Slice, cc.Array:
			return left.Base

		case c2go.String:
			return byteType
		}
		return nil

	case cc.Lsh, cc.Rsh:
		left := fixGoTypesExpr(fn, x.Left, targ)
		fixShiftCount(fn, x.Right)
		return left

	case cc.LshEq, cc.RshEq:
		left := fixGoTypesExpr(fn, x.Left, nil)
		fixShiftCount(fn, x.Right)
		return left

	case cc.Name:
		if x.Text == "nelem" {
			x.Text = "len"
			x.XDecl = nil
			return intType
		}

		if x.XDecl == nil {
			return nil
		}
		return x.XDecl.Type

	case cc.Number:
		if targ == nil {
			return nil
		}
		if targ.Kind <= c2go.Int8 && targ.Kind <= c2go.Float64 {
			return targ
		}
		return intType

	case cc.Minus, cc.Plus, cc.Twid:
		return fixGoTypesExpr(fn, x.Left, targ)

	case cc.Offsetof:
		// TODO
		return nil

	case cc.Paren:
		return fixGoTypesExpr(fn, x.Left, targ)

	case cc.PostDec, cc.PostInc:
		left := fixGoTypesExpr(fn, x.Left, nil)

		if x.Op == cc.PostInc && isSliceOrString(left) {
			old := copyExpr(x.Left)
			x.Op = cc.Eq
			x.Right = &cc.Expr{Op: c2go.ExprSlice, List: []*cc.Expr{old, &cc.Expr{Op: cc.Number, Text: "1"}, nil}}
		}

		return nil

	case cc.SizeofExpr, cc.SizeofType:
		// TODO
		return nil

	case cc.String:
		return &cc.Type{Kind: c2go.String}

	case cc.VaArg:
		// TODO
		return nil
	}
}

var (
	boolType   = &cc.Type{Kind: c2go.Bool}
	byteType   = &cc.Type{Kind: c2go.Byte}
	intType    = &cc.Type{Kind: c2go.Int}
	uintType   = &cc.Type{Kind: c2go.Uint}
	uint32Type = &cc.Type{Kind: c2go.Uint32}
	uint64Type = &cc.Type{Kind: c2go.Uint64}
)

func forceGoType(fn *cc.Decl, x *cc.Expr, targ *cc.Type) {
	actual := fixGoTypesExpr(fn, x, targ)
	forceConvert(fn, x, actual, targ)
}

func forceConvert(fn *cc.Decl, x *cc.Expr, actual, targ *cc.Type) {
	if isNumericConst(x) && targ != nil {
		switch targ.Kind {
		case cc.Ptr, c2go.Slice:
			if x.Op == cc.Number && x.Text == "0" {
				x.Op = cc.Name
				x.Text = "nil"
				x.XType = targ
			}
		case c2go.String:
			if x.Op == cc.Number && x.Text == "0" {
				x.Op = cc.Name
				x.Text = `""`
				x.XType = targ
			}

		}
		return
	}

	if actual != nil && targ != nil && !sameType(actual, targ) {
		old := copyExpr(x)
		x.Op = cc.Cast
		x.Left = old
		x.Right = nil
		x.List = nil
		x.Type = targ
		x.XType = targ
		if actual.Kind == cc.Array && targ.Kind == c2go.Slice {
			x.Op = c2go.ExprSlice
			x.List = []*cc.Expr{old, nil, nil}
			x.Left = nil
			x.Type = nil
		}
	}
}

func isNumericConst(x *cc.Expr) bool {
	// TODO: better
	return x.Op == cc.Number
}

func fixShiftCount(fn *cc.Decl, x *cc.Expr) {
	typ := fixGoTypesExpr(fn, x, nil)
	if typ == nil {
		return
	}
	switch typ.Kind {
	case c2go.Uint8, c2go.Uint16, c2go.Uint32, c2go.Uint64, c2go.Uint, c2go.Uintptr, c2go.Byte:
		return
	}
	if typ.Kind == c2go.Int64 {
		forceConvert(fn, x, typ, uint64Type)
		return
	}
	forceConvert(fn, x, typ, uintType)
}

func fixBinary(fn *cc.Decl, x *cc.Expr, left, right, targ *cc.Type) *cc.Type {
	if left == nil || right == nil || left.Kind < c2go.Int8 || left.Kind > c2go.Float64 || right.Kind < c2go.Int8 || right.Kind > c2go.Float64 || targ == nil || targ.Kind < c2go.Int8 || targ.Kind > c2go.Float64 {
		return nil
	}

	// Want to do everything at as high a precision as possible for as long as possible.
	// If target is wider, convert early.
	// If target is narrower, don't convert at all - let next step do it.
	// Must make left and right match.
	// Convert to largest of three.
	t := left
	if t.Kind < right.Kind {
		t = right
	}
	if t.Kind < targ.Kind {
		t = targ
	}
	if !sameType(t, left) {
		forceConvert(fn, x.Left, left, t)
	}
	if !sameType(t, right) {
		forceConvert(fn, x.Right, right, t)
	}
	return t
}

func sameType(t, u *cc.Type) bool {
	if t == u {
		return true
	}
	if t == nil || u == nil {
		return false
	}
	if t.Kind != u.Kind {
		return false
	}
	if t.Name != "" || u.Name != "" {
		return t.Name == u.Name
	}
	if !sameType(t.Base, u.Base) || len(t.Decls) != len(u.Decls) {
		return false
	}
	for i, td := range t.Decls {
		ud := u.Decls[i]
		if !sameType(td.Type, ud.Type) || t.Kind == cc.Struct && td.Name != ud.Name {
			return false
		}
	}
	return true
}

func isSliceOrString(typ *cc.Type) bool {
	return typ != nil && (typ.Kind == c2go.Slice || typ.Kind == c2go.String)
}

func fixSpecialCall(fn *cc.Decl, x *cc.Expr) bool {
	if x.Left.Op != cc.Name {
		return false
	}
	switch x.Left.Text {
	case "memset":
		if len(x.List) != 3 || x.List[1].Op != cc.Number || x.List[1].Text != "0" {
			fprintf(x.Span, "unsupported memset - nonzero")
			return false
		}
		obj, objType := objIndir(fn, x.List[0])
		if !matchSize(fn, obj, objType, x.List[2]) {
			fprintf(x.Span, "unsupported memset - wrong size")
			return true
		}

		x.Op = cc.Eq
		x.Left = obj
		x.Right = zeroFor(objType)
		x.List = nil
		return true

	case "memmove":
		if len(x.List) != 3 {
			fprintf(x.Span, "unsupported %v", x)
			return false
		}
		obj1, obj1Type := objIndir(fn, x.List[0])
		obj2, obj2Type := objIndir(fn, x.List[1])
		if obj1Type == nil || obj2Type == nil {
			fprintf(x.Span, "unsupported %v - missing types", x)
			return true
		}

		siz := x.List[2]
		if siz.Op == cc.Number && siz.Text == "4" {
			if (obj1Type.Kind == c2go.Uint32 || obj1Type.Kind == c2go.Int32) && obj2Type.Kind == c2go.Float32 {
				x.Op = cc.Eq
				x.Left = obj1
				x.Right = &cc.Expr{
					Op: cc.Call,
					Left: &cc.Expr{Op: cc.Name,
						Text: "math.Float32bits",
					},
					List: []*cc.Expr{obj2},
				}
				x.XType = uint32Type
				return true
			}
			fprintf(x.Span, "unsupported %v - size 8 type %v %v", x, obj1Type, obj2Type)
		}
		if siz.Op == cc.Number && siz.Text == "8" {
			if (obj1Type.Kind == c2go.Uint64 || obj1Type.Kind == c2go.Int64) && obj2Type.Kind == c2go.Float64 {
				x.Op = cc.Eq
				x.Left = obj1
				x.Right = &cc.Expr{
					Op: cc.Call,
					Left: &cc.Expr{Op: cc.Name,
						Text: "math.Float64bits",
					},
					List: []*cc.Expr{obj2},
				}
				x.XType = uint64Type
				return true
			}
			fprintf(x.Span, "unsupported %v - size 8 type %v %v", x, obj1Type, obj2Type)
		}
		fprintf(x.Span, "unsupported %v", x)
		return true

	case "memcmp":
		if len(x.List) != 3 {
			fprintf(x.Span, "unsupported %v", x)
			return false
		}
		obj1, obj1Type := objIndir(fn, x.List[0])
		obj2, obj2Type := objIndir(fn, x.List[1])
		if obj1Type == nil || !sameType(obj1Type, obj2Type) {
			fprintf(x.Span, "unsupported %v", x)
			return true
		}

		if !matchSize(fn, obj1, obj1Type, x.List[2]) && !matchSize(fn, obj2, obj2Type, x.List[2]) {
			fprintf(x.Span, "unsupported %v - wrong size", x)
			return true
		}

		x.Op = cc.EqEq
		x.Left = obj1
		x.Right = obj2
		x.List = nil
		x.XType = boolType
		return true
	}

	return false
}

func objIndir(fn *cc.Decl, x *cc.Expr) (*cc.Expr, *cc.Type) {
	objType := fixGoTypesExpr(fn, x, nil)
	obj := x
	if obj.XType != nil && obj.XType.Kind == cc.Array {
		// obj stays as is
	} else if obj.Op == cc.Addr {
		obj = obj.Left
		if objType != nil {
			objType = objType.Base
		}
	} else {
		obj = &cc.Expr{Op: cc.Indir, Left: obj}
		if objType != nil {
			objType = objType.Base
		}
	}
	if objType == nil {
		objType = obj.XType
	}
	return obj, objType
}

func matchSize(fn *cc.Decl, obj *cc.Expr, objType *cc.Type, siz *cc.Expr) bool {
	switch siz.Op {
	default:
		return false

	case cc.SizeofType:
		// ok if sizeof type of first arg
		return sameType(siz.Type, objType)

	case cc.SizeofExpr:
		// ok if sizeof *firstarg
		y := siz.Left
		if y.Op == cc.Paren {
			y = y.Left
		}
		return obj.String() == y.String()
	}
}
