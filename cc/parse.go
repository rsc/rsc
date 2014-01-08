// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Read(name string, r io.Reader) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("reading %s: %v", name, err)
	}
	lx := &lexer{
		start:  startProg,
		input:  string(data),
		file:   name,
		lineno: 1,
	}
	yyParse(lx)
	if lx.errors != nil {
		for _, err := range lx.errors {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(1)
	}
	if lx.input != "" {
		log.Fatalf("reading %s: did not consume entire file", name)
		os.Exit(1)
	}
}

func ParseExpr(str string) (*Expr, error) {
	lx := &lexer{
		input:  str + "\n",
		file:   "<string>",
		lineno: 1,
		start:  startExpr,
	}
	yyParse(lx)
	if lx.errors != nil {
		return nil, fmt.Errorf("parsing expression %#q: %v", str, lx.errors[0])
	}
	return lx.expr, nil
}

type Prog struct {
	Decl []*Decl
}
