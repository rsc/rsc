// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"fmt"
	"strings"
)

var printf = fmt.Printf

type Type struct {
	Span  Span
	Kind  TypeKind
	Qual  TypeQual
	Base  *Type
	Tag   string
	Decls []*Decl
	Width *Expr
	Name  string
}

type TypeKind int

const (
	_ TypeKind = iota
	Void
	Char
	Uchar
	Short
	Ushort
	Int
	Uint
	Long
	Ulong
	Longlong
	Ulonglong
	Float
	Double
	Enum
	Ptr
	Struct
	Union
	Array
	Func
	TypedefType
)

var typeKindString = []string{
	Void:        "void",
	Char:        "char",
	Uchar:       "uchar",
	Short:       "short",
	Ushort:      "ushort",
	Int:         "int",
	Uint:        "uint",
	Long:        "long",
	Ulong:       "ulong",
	Longlong:    "longlong",
	Ulonglong:   "ulonglong",
	Float:       "float",
	Double:      "double",
	Ptr:         "pointer",
	Struct:      "struct",
	Union:       "union",
	Enum:        "enum",
	Array:       "array",
	Func:        "func",
	TypedefType: "<typedef>",
}

func (k TypeKind) String() string {
	if 0 <= int(k) && int(k) <= len(typeKindString) {
		return typeKindString[k]
	}
	return fmt.Sprintf("TypeKind(%d)", k)
}

type TypeQual int

const (
	Const TypeQual = 1 << iota
	Volatile
)

func (q TypeQual) String() string {
	s := ""
	if q&Const != 0 {
		s += "const "
	}
	if q&Volatile != 0 {
		s += "volatile "
	}
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}

type Storage int

const (
	Auto Storage = 1 << iota
	Static
	Extern
	Typedef
	Register
	Inline
)

func (c Storage) String() string {
	s := ""
	if c&Auto != 0 {
		s += "auto "
	}
	if c&Static != 0 {
		s += "static "
	}
	if c&Extern != 0 {
		s += "extern "
	}
	if c&Typedef != 0 {
		s += "typedef "
	}
	if c&Register != 0 {
		s += "register "
	}
	if c&Inline != 0 {
		s += "inline "
	}
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}

var (
	typeChar      = newType(Char)
	typeUchar     = newType(Uchar)
	typeShort     = newType(Short)
	typeUshort    = newType(Ushort)
	typeInt       = newType(Int)
	typeUint      = newType(Uint)
	typeLong      = newType(Long)
	typeUlong     = newType(Ulong)
	typeLonglong  = newType(Longlong)
	typeUlonglong = newType(Ulonglong)
	typeFloat     = newType(Float)
	typeDouble    = newType(Double)
	typeVoid      = newType(Void)
	typeBool = &Type{Kind: TypedefType, Name: "bool", Base: typeInt}
)

type typeOp int

const (
	tChar typeOp = 1 << iota
	tShort
	tInt
	tLong
	tSigned
	tUnsigned
	tFloat
	tDouble
	tVoid
	tLonglong
)

var builtinTypes = map[typeOp]*Type{
	tChar:                     typeChar,
	tChar | tSigned:           typeChar,
	tChar | tUnsigned:         typeUchar,
	tShort:                    typeShort,
	tShort | tSigned:          typeShort,
	tShort | tUnsigned:        typeUshort,
	tShort | tInt:             typeShort,
	tShort | tSigned | tInt:   typeShort,
	tShort | tUnsigned | tInt: typeUshort,
	tInt:                         typeInt,
	tInt | tSigned:               typeInt,
	tInt | tUnsigned:             typeUint,
	tLong:                        typeLong,
	tLong | tSigned:              typeLong,
	tLong | tUnsigned:            typeUlong,
	tLong | tInt:                 typeLong,
	tLong | tSigned | tInt:       typeLong,
	tLong | tUnsigned | tInt:     typeUlong,
	tLonglong:                    typeLonglong,
	tLonglong | tSigned:          typeLonglong,
	tLonglong | tUnsigned:        typeUlonglong,
	tLonglong | tInt:             typeLonglong,
	tLonglong | tSigned | tInt:   typeLonglong,
	tLonglong | tUnsigned | tInt: typeUlonglong,
	tFloat:  typeFloat,
	tDouble: typeDouble,
	tVoid:   typeVoid,
}

func splitTypeWords(ws []string) (c Storage, q TypeQual, ty *Type) {
	// Could check for doubled words in general,
	// like const const, but no one cares.
	var t typeOp
	var ts []string
	for _, w := range ws {
		switch w {
		case "const":
			q |= Const
		case "volatile":
			q |= Volatile
		case "auto":
			c |= Auto
		case "static":
			c |= Static
		case "extern":
			c |= Extern
		case "typedef":
			c |= Typedef
		case "register":
			c |= Register
		case "inline":
			c |= Inline
		case "char":
			t |= tChar
			ts = append(ts, w)
		case "short":
			t |= tShort
			ts = append(ts, w)
		case "int":
			t |= tInt
			ts = append(ts, w)
		case "long":
			if t&tLong != 0 {
				t ^= tLonglong | tLong
			} else {
				t |= tLong
			}
			ts = append(ts, w)
		case "signed":
			t |= tSigned
			ts = append(ts, w)
		case "unsigned":
			t |= tUnsigned
			ts = append(ts, w)
		case "float":
			t |= tFloat
			ts = append(ts, w)
		case "double":
			t |= tDouble
			ts = append(ts, w)
		case "void":
			t |= tVoid
			ts = append(ts, w)
		}
	}

	if t == 0 {
		t |= tInt
	}

	ty = builtinTypes[t]
	if ty == nil {
		fmt.Printf("unsupported type %q\n", strings.Join(ts, " "))
	}

	return c, q, builtinTypes[t]
}

func newType(k TypeKind) *Type {
	return &Type{Kind: k}
}

func (t *Type) String() string {
	if t == nil {
		return "<nil>"
	}
	switch t.Kind {
	default:
		return t.Kind.String()
	case TypedefType:
		return t.Name
	case Ptr:
		return t.Base.String() + "*"
	case Struct, Union, Enum:
		if t.Tag == "" {
			return t.Kind.String()
		}
		return t.Kind.String() + " " + t.Tag
	case Array:
		return t.Base.String() + "[]"
	case Func:
		return "func() " + t.Base.String()
	}
}

type Decl struct {
	Span    Span
	Name    string
	Type    *Type
	Storage Storage
	Init    *Init
	Body    *Stmt
}
