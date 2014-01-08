// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"bytes"
	"fmt"
)

type Printer struct {
	buf bytes.Buffer
}

func (p *Printer) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *Printer) String() string {
	return p.buf.String()
}

type exprPrec struct {
	expr *Expr
	prec int
}

func (p *Printer) Print(args ...interface{}) {
	for _, arg := range args {
		switch arg := arg.(type) {
		default:
			fmt.Fprintf(&p.buf, "(?%T)", arg)
		case string:
			p.buf.WriteString(arg)
		case exprPrec:
			p.printExpr(arg.expr, arg.prec)
		case *Expr:
			p.printExpr(arg, precLow)
		case *Prefix:
			p.printPrefix(arg)
		case *Init:
			p.printInit(arg)
		case *Prog:
			p.printProg(arg)
		//case *Stmt:
		//	p.printStmt(arg)
		case *Label:
			p.printLabel(arg)
		case *Type:
			p.printType(arg)
		case *Decl:
			p.printDecl(arg)
		}
	}
}

const (
	precNone = iota
	precArrow
	precAddr
	precMul
	precAdd
	precLsh
	precLt
	precEqEq
	precAnd
	precXor
	precOr
	precAndAnd
	precOrOr
	precCond
	precEq
	precComma
	precLow
)

var opPrec = []int{
	Add:        precAdd,
	AddEq:      precEq,
	Addr:       precAddr,
	And:        precAnd,
	AndAnd:     precAndAnd,
	AndEq:      precEq,
	Arrow:      precArrow,
	Call:       precArrow,
	Cast:       precAddr,
	CastInit:   precAddr,
	Comma:      precComma,
	Cond:       precCond,
	Div:        precMul,
	DivEq:      precEq,
	Dot:        precArrow,
	Eq:         precEq,
	EqEq:       precEqEq,
	Gt:         precLt,
	GtEq:       precLt,
	Index:      precArrow,
	Indir:      precAddr,
	Lsh:        precLsh,
	LshEq:      precEq,
	Lt:         precLt,
	LtEq:       precEq,
	Minus:      precAddr,
	Mod:        precMul,
	ModEq:      precEq,
	Mul:        precMul,
	MulEq:      precEq,
	Name:       precNone,
	Not:        precAddr,
	NotEq:      precEqEq,
	Number:     precNone,
	Offsetof:   precAddr,
	Or:         precOr,
	OrEq:       precEq,
	OrOr:       precOrOr,
	Paren:      precLow,
	Plus:       precAddr,
	PostDec:    precAddr,
	PostInc:    precAddr,
	PreDec:     precAddr,
	PreInc:     precAddr,
	Rsh:        precLsh,
	RshEq:      precEq,
	SizeofExpr: precAddr,
	SizeofType: precAddr,
	String:     precNone,
	Sub:        precAdd,
	SubEq:      precEq,
	Twid:       precAddr,
	VaArg:      precAddr,
	Xor:        precXor,
	XorEq:      precEq,
}

var opStr = []string{
	Add:        "+",
	AddEq:      "+=",
	Addr:       "&",
	And:        "&",
	AndAnd:     "&&",
	AndEq:      "&=",
	Arrow:      "->",
	Div:        "/",
	DivEq:      "/=",
	Dot:        ".",
	Eq:         "=",
	EqEq:       "==",
	Gt:         ">",
	GtEq:       ">=",
	Indir:      "*",
	Lsh:        "<<",
	LshEq:      "<<=",
	Lt:         "<",
	LtEq:       "<=",
	Minus:      "-",
	Mod:        "%",
	ModEq:      "%=",
	Mul:        "*",
	MulEq:      "*=",
	Not:        "!",
	NotEq:      "!=",
	Or:         "|",
	OrEq:       "|=",
	OrOr:       "||",
	Plus:       "+",
	PreDec:     "--",
	PreInc:     "++",
	Rsh:        ">>",
	RshEq:      ">>=",
	Sub:        "-",
	SubEq:      "-=",
	Twid:       "~",
	Xor:        "^",
	XorEq:      "^=",
	SizeofExpr: "sizeof ",
}

func (p *Printer) printExpr(x *Expr, prec int) {
	var newPrec int
	if 0 <= int(x.Op) && int(x.Op) < len(opPrec) {
		newPrec = opPrec[x.Op]
	}
	if prec < newPrec {
		p.buf.WriteString("(")
		defer p.buf.WriteString(")")
	}
	prec = newPrec

	var str string
	if 0 <= int(x.Op) && int(x.Op) < len(opStr) {
		str = opStr[x.Op]
	}
	if str != "" {
		if x.Right != nil {
			// binary operator
			if prec == precEq {
				// right associative
				p.Print(exprPrec{x.Left, prec - 1}, " ", str, " ", exprPrec{x.Right, prec})
			} else {
				// left associative
				p.Print(exprPrec{x.Left, prec}, " ", str, " ", exprPrec{x.Right, prec - 1})
			}
		} else {
			// unary operator
			if (x.Op == Plus || x.Op == Minus || x.Op == Addr) && x.Left.Op == x.Op ||
				x.Op == Plus && x.Left.Op == PreInc ||
				x.Op == Minus && x.Left.Op == PreDec {
				prec-- // force parenthesization +(+x) not ++x
			}
			p.Print(str, exprPrec{x.Left, prec})
		}
		return
	}

	// special cases
	switch x.Op {
	default:
		panic(fmt.Sprintf("printExpr missing case for %v", x.Op))

	case Arrow:
		p.Print(exprPrec{x.Left, prec}, "->", x.Text)

	case Call:
		p.Print(exprPrec{x.Left, precAddr}, "(")
		for i, y := range x.List {
			if i > 0 {
				p.Print(", ")
			}
			p.printExpr(y, prec-1)
		}
		p.Print(")")

	case Cast:
		p.Print("(", x.Type, ")", exprPrec{x.Left, prec})

	case CastInit:
		p.Print("(", x.Type, ")", x.Init)

	case Comma:
		for i, y := range x.List {
			if i > 0 {
				p.Print(", ")
			}
			p.printExpr(y, prec-1)
		}

	case Cond:
		p.Print(exprPrec{x.List[0], prec - 1}, " ? ", exprPrec{x.List[1], prec}, " : ", exprPrec{x.List[2], prec})

	case Dot:
		p.Print(exprPrec{x.Left, prec}, ".", x.Text)

	case Index:
		p.Print(exprPrec{x.Left, prec}, "[", exprPrec{x.Right, precLow}, "]")

	case Name, Number:
		p.Print(x.Text)

	case String:
		for i, str := range x.Texts {
			if i > 0 {
				p.buf.WriteString(" ")
			}
			p.buf.WriteString(str)
		}

	case Offsetof:
		p.Print("offsetof(", x.Type, ", ", exprPrec{x.Left, precComma}, ")")

	case Paren:
		p.Print("(", exprPrec{x.Left, prec}, ")")

	case PostDec:
		p.Print(exprPrec{x.Left, prec}, "--")

	case PostInc:
		p.Print(exprPrec{x.Left, prec}, "++")

	case SizeofType:
		p.Print("sizeof(", x.Type, ")")

	case VaArg:
		p.Print("va_arg(", exprPrec{x.Left, precComma}, ", ", x.Type, ")")
	}
}

func (p *Printer) printPrefix(x *Prefix) {
}

func (p *Printer) printInit(x *Init) {
	for _, pre := range x.Prefix {
		p.Print(pre)
	}
	if x.Expr != nil {
		p.printExpr(x.Expr, precComma)
	} else {
		p.Print("{")
		for i, y := range x.Braced {
			if i > 0 {
				p.Print(", ")
			}
			p.Print(y)
		}
		p.Print("}")
	}
}

func (p *Printer) printProg(x *Prog) {
}

/*
func (p *Printer) printStmt(x *Stmt) {
	// TODO labels

	switch x.Op {
	case ARGBEGIN:
		p.print("ARGBEGIN{", indent, x.Body, unindent, "\n", "}ARGEND")

	case Block:
		p.print("{", indent)
		for _, b := range x.Stmt {
			p.print("\n", b)
		}
		p.print(unindent, "\n", "}")

	case Break:
		p.print("break;")

	case Continue:
		p.print("continue;")

	case StmtDecl:
		//xxx

	case Do:
		p.print("do", nestBlock{p.Body}, " while(", p.Expr, ");")

	case Empty:
		p.print(";")

	case StmtExpr:
		p.print(x.Expr)

	case For:
		p.print("for(", x.Pre, "; ", x.Expr, "; ", x.Post, ")", nestBlock{p.Body})

	case If:
		p.print("if(", x.Expr, ")", nestBlock{p.Body})
		if x.Else != nil {
			if p.Body.Op == Block && p.Body.Labels == nil {
				p.print(" else")
			} else {
				p.print("\n", "else")
			}
			p.print(nestBlock{p.Else})
		}

	case Goto:
		p.print("goto ", x.Text, ";")

	case Return:
		if x.Expr == nil {
			p.print("return;")
		} else {
			p.print("return ", x.Expr, ";")
		}

	case Switch:
		p.print("switch(", x.Expr, ")", nestBlock{p.Body})

	case While:
		p.print("while(", x.Expr, ")", nestBlock{p.Body})
	}
}
*/

func (p *Printer) printLabel(x *Label) {
}

func (p *Printer) printType(x *Type) {
	switch x.Kind {
	default:
		p.Print(x.Kind.String())
		//case Ptr:
		//case Struct:
		//case Union:
		//case Enum:
		//case Array:
		//case Func:
	}
}

func (p *Printer) printDecl(x *Decl) {
}
