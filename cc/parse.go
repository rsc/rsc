// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func Read(name string, r io.Reader) (*Prog, error) {
	return ReadMany([]string{name}, []io.Reader{r})
}

func ReadMany(names []string, readers []io.Reader) (*Prog, error) {
	lx := &lexer{}
	var prog *Prog
	for i, name := range names {
		r := readers[i]
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		data = append(data, '\n')
		lx.start = startProg
		lx.lexInput = lexInput{
			input:  string(data),
			file:   name,
			lineno: 1,
		}
		lx.parse()
		if lx.errors != nil {
			return nil, fmt.Errorf("%v", lx.errors[0])
		}
		if prog == nil {
			prog = lx.prog
		} else {
			prog.Span.End = lx.prog.Span.End
			prog.Decls = append(prog.Decls, lx.prog.Decls...)
		}
		lx.prog = nil
	}
	lx.prog = prog
	lx.assignComments()
	lx.typecheck(lx.prog)
	if lx.errors != nil {
		return nil, fmt.Errorf("%v", strings.Join(lx.errors, "\n"))
	}
	return lx.prog, nil
}

func ParseExpr(str string) (*Expr, error) {
	lx := &lexer{
		start: startExpr,
		lexInput: lexInput{
			input:  str + "\n",
			file:   "<string>",
			lineno: 1,
		},
	}
	lx.parse()
	if lx.errors != nil {
		return nil, fmt.Errorf("parsing expression %#q: %v", str, lx.errors[0])
	}
	return lx.expr, nil
}

type Prog struct {
	SyntaxInfo
	Decls []*Decl
}
