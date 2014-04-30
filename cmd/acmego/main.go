// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Acmego watches acme for .go files being written.
// Each time a .go file is written, acmego checks whether the
// import block needs adjustment. If so, it makes the changes
// in the window body but does not write the file.
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"code.google.com/p/goplan9/plan9/acme"
)

func main() {
	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}
		if event.Name != "" && event.Op == "put" && strings.HasSuffix(event.Name, ".go") {
			reformat(event.ID, event.Name)
		}
	}
}

func reformat(id int, name string) {
	w, err := acme.Open(id, nil)
	if err != nil {
		//log.Print(err)
		return
	}
	defer w.CloseFiles()

	old, err := ioutil.ReadFile(name)
	if err != nil {
		//log.Print(err)
		return
	}
	new, err := exec.Command("goimports", name).CombinedOutput()
	if err != nil {
		//log.Print(err)
		return
	}

	oldTop, err := readImports(bytes.NewReader(old), true)
	if err != nil {
		//log.Print(err)
		return
	}
	newTop, err := readImports(bytes.NewReader(new), true)
	if err != nil {
		//log.Print(err)
		return
	}

	if bytes.Equal(oldTop, newTop) {
		return
	}

	w.Addr("#0,#%d", len(oldTop))
	w.Write("data", newTop)
}
