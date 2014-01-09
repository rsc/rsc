// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type lexer struct {
	// input
	start int
	lexInput
	pushed []lexInput
	forcePos Pos
	c2goComment bool // inside /*c2go ... */ comment

	// type checking state
	scope *Scope

	// output
	errors []string
	prog   *Prog
	expr   *Expr
}

func (lx *lexer) parse() {
	lx.scope = &Scope{}
	yyParse(lx)
}

type lexInput struct {
	input   string
	tok     string
	lastsym string
	file    string
	lineno  int
	byte    int
}

func (lx *lexer) pushInclude(includeLine string) {
	s := strings.TrimSpace(strings.TrimPrefix(includeLine, "#include"))
	if !strings.HasPrefix(s, "<") && !strings.HasPrefix(s, "\"") {
		lx.Errorf("malformed #include")
		return
	}
	sep := ">"
	if s[0] == '"' {
		sep = "\""
	}
	i := strings.Index(s[1:], sep)
	if i < 0 {
		lx.Errorf("malformed #include")
		return
	}
	i++
	file := s[1:i]

	file, data, err := lx.findInclude(file, s[0] == '<')
	if err != nil {
		lx.Errorf("#include %s: %v", s[:i+1], err)
		return
	}

	if data == nil {
		return
	}

	lx.pushed = append(lx.pushed, lx.lexInput)
	lx.lexInput = lexInput{
		input:  string(append(data, '\n')),
		file:   file,
		lineno: 1,
	}
}

var stdMap = map[string]string{
	"u.h":    hdr_u_h,
	"libc.h": hdr_libc_h,
}

func (lx *lexer) findInclude(name string, std bool) (string, []byte, error) {
	if std {
		if redir, ok := stdMap[name]; ok {
			if redir == "" {
				return "", nil, nil
			}
			return "internal/" + name, []byte(redir), nil
		}
		name = "/Users/rsc/g/go/include/" + name
	}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", nil, err
	}
	return name, data, nil
}

func (lx *lexer) pop() bool {
	if len(lx.pushed) == 0 {
		return false
	}
	lx.lexInput = lx.pushed[len(lx.pushed)-1]
	lx.pushed = lx.pushed[:len(lx.pushed)-1]
	return true
}

func (lx *lexer) pos() Pos {
	if lx.forcePos.Line != 0 {
		return lx.forcePos
	}
	return Pos{lx.file, lx.lineno, lx.byte}
}
func (lx *lexer) span() Span {
	p := lx.pos()
	return Span{p, p}
}

func (lx *lexer) setSpan(s Span) {
	lx.forcePos = s.Start
}

func span(l1, l2 Span) Span {
	if l1.Start.Line == 0 {
		return l2
	}
	if l2.Start.Line == 0 {
		return l1
	}
	return Span{l1.Start, l2.End}
}

func (lx *lexer) skip(i int) {
	lx.lineno += strings.Count(lx.input[:i], "\n")
	lx.input = lx.input[i:]
	lx.byte += i
}

func (lx *lexer) token(i int) {
	lx.tok = lx.input[:i]
	lx.skip(i)
}

func (lx *lexer) sym(i int) {
	lx.token(i)
	lx.lastsym = lx.tok
}

func (lx *lexer) comment(i int) {
	com := lx.input[:i]
	_ = com
	lx.skip(i)
}

func isalpha(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || c == '_' || c >= 0x80 || '0' <= c && c <= '9'
}

func isspace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

func (lx *lexer) setEnd(yy *yySymType) {
	yy.span.End = lx.pos()
}

func (lx *lexer) Lex(yy *yySymType) int {
	if lx.start != 0 {
		tok := lx.start
		lx.start = 0
		return tok
	}
	*yy = yySymType{}
	defer lx.setEnd(yy)
Restart:
	yy.span.Start = lx.pos()
	in := lx.input
	if len(in) == 0 {
		if lx.pop() {
			goto Restart
		}
		return tokEOF
	}
	c := in[0]
	if isspace(c) {
		i := 1
		for i < len(in) && isspace(in[i]) {
			i++
		}
		lx.skip(i)
		goto Restart
	}

	i := 0
	switch c {
	case '#':
		i++
		for in[i] != '\n' {
			if in[i] == '\\' && in[i+1] == '\n' && i+2 < len(in) {
				i++
			}
			i++
		}
		str := in[:i]
		lx.skip(i)
		if strings.HasPrefix(str, "#include") {
			lx.pushInclude(str)
		}
		goto Restart

	case 'L':
		if in[1] != '\'' && in[1] != '"' {
			break // goes to alpha case after switch
		}
		i = 1
		fallthrough

	case '"', '\'':
		q := in[i]
		i++ // for the quote
		for ; in[i] != q; i++ {
			if in[i] == '\n' {
				what := "string"
				if q == '\'' {
					what = "character"
				}
				lx.Errorf("unterminated %s constant", what)
			}
			if in[i] == '\\' {
				i++
			}
		}
		i++ // for the quote
		lx.sym(i)
		yy.str = lx.tok
		if q == '"' {
			return tokString
		} else {
			return tokLitChar
		}

	case '.':
		if in[1] < '0' || '9' < in[1] {
			if in[1] == '.' && in[2] == '.' {
				lx.token(3)
				return tokDotDotDot
			}
			lx.token(1)
			return int(c)
		}
		fallthrough

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		for '0' <= in[i] && in[i] <= '9' || in[i] == '.' || 'A' <= in[i] && in[i] <= 'Z' || 'a' <= in[i] && in[i] <= 'z' {
			i++
		}
		lx.sym(i)
		yy.str = lx.tok
		return tokNumber

	case '/':
		switch in[1] {
		case '*':
			if strings.HasPrefix(in, "/*c2go") {
				lx.skip(6)
				lx.c2goComment = true
				goto Restart
			}
			i := 2
			for ; ; i++ {
				if i+2 <= len(in) && in[i] == '*' && in[i+1] == '/' {
					i += 2
					break
				}
				if i >= len(in) {
					lx.Errorf("unterminated /* comment")
					return tokError
				}
			}
			lx.comment(i)
			goto Restart

		case '/':
			for in[i] != '\n' {
				i++
			}
			i++
			lx.comment(i)
			goto Restart
		}
		fallthrough

	case '~', '*', '(', ')', '[', ']', '{', '}', '?', ':', ';', ',', '%', '^', '!', '=', '<', '>', '+', '-', '&', '|':
		if lx.c2goComment && in[0] == '*' && in[1] == '/' {
			lx.c2goComment = false
			lx.skip(2)
			goto Restart
		}
		if c == '-' && in[1] == '>' {
			lx.token(2)
			return tokArrow
		}
		if in[1] == '=' && tokEq[c] != 0 {
			lx.token(2)
			return int(tokEq[c])
		}
		if in[1] == c && tokTok[c] != 0 {
			if in[2] == '=' && tokTokEq[c] != 0 {
				lx.token(3)
				return int(tokTokEq[c])
			}
			lx.token(2)
			return int(tokTok[c])
		}
		lx.token(1)
		return int(c)
	}

	if isalpha(c) {
		for isalpha(in[i]) {
			i++
		}
		lx.sym(i)
		yy.str = lx.tok
		if t := tokId[lx.tok]; t != 0 {
			return int(t)
		}
		yy.decl = lx.lookupDecl(lx.tok)
		if yy.decl != nil && yy.decl.Storage&Typedef != 0 {
			t := yy.decl.Type
			for t.Kind == TypedefType && t.Base != nil {
				t = t.Base
			}
			yy.typ = &Type{Kind: TypedefType, Name: yy.str, Base: t}
			return tokTypeName
		}
		return tokName
	}

	lx.Errorf("unexpected input byte %#02x (%c)", c, c)
	return tokError
}

func (lx *lexer) Error(s string) {
	lx.Errorf("%s near %s", s, lx.lastsym)
}

func (lx *lexer) Errorf(format string, args ...interface{}) {
	lx.errors = append(lx.errors, fmt.Sprintf("%s: %s", lx.span(), fmt.Sprintf(format, args...)))
}

type Pos struct {
	File string
	Line int
	Byte int
}

type Span struct {
	Start Pos
	End   Pos
}

func (l Span) String() string {
	return fmt.Sprintf("%s:%d", l.Start.File, l.Start.Line)
}

var tokEq = [256]int32{
	'*': tokMulEq,
	'/': tokDivEq,
	'+': tokAddEq,
	'-': tokSubEq,
	'%': tokModEq,
	'^': tokXorEq,
	'!': tokNotEq,
	'=': tokEqEq,
	'<': tokLtEq,
	'>': tokGtEq,
	'&': tokAndEq,
	'|': tokOrEq,
}

var tokTok = [256]int32{
	'<': tokLsh,
	'>': tokRsh,
	'=': tokEqEq,
	'+': tokInc,
	'-': tokDec,
	'&': tokAndAnd,
	'|': tokOrOr,
}

var tokTokEq = [256]int32{
	'<': tokLshEq,
	'>': tokRshEq,
}

var tokId = map[string]int32{
	"auto":     tokAuto,
	"break":    tokBreak,
	"case":     tokCase,
	"char":     tokChar,
	"const":    tokConst,
	"continue": tokContinue,
	"default":  tokDefault,
	"do":       tokDo,
	"double":   tokDouble,
	"else":     tokElse,
	"enum":     tokEnum,
	"extern":   tokExtern,
	"float":    tokFloat,
	"for":      tokFor,
	"goto":     tokGoto,
	"if":       tokIf,
	"inline":   tokInline,
	"int":      tokInt,
	"long":     tokLong,
	"offsetof": tokOffsetof,
	"register": tokRegister,
	"return":   tokReturn,
	"short":    tokShort,
	"signed":   tokSigned,
	"sizeof":   tokSizeof,
	"static":   tokStatic,
	"struct":   tokStruct,
	"switch":   tokSwitch,
	"typedef":  tokTypedef,
	"union":    tokUnion,
	"unsigned": tokUnsigned,
	"va_arg":   tokVaArg,
	"void":     tokVoid,
	"volatile": tokVolatile,
	"while":    tokWhile,

	"ARGBEGIN": tokARGBEGIN,
	"ARGEND":   tokARGEND,
	"AUTOLIB":  tokAUTOLIB,
	"USED": tokUSED,
	"SET": tokSET,
}
