// +build ignore

package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/rsc/c2go"
	"code.google.com/p/rsc/cc"
)

// print an error; fprintf is a bad name but helps go vet.
func fprintf(span cc.Span, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s:%d: %s\n", span.Start.File, span.Start.Line, msg)
}

// write actual output
func write(prog *cc.Prog, files []string, cfg *Config) {
	for _, file := range files {
		writeFile(prog, file, "", cfg)
	}
	writeFile(prog, "/Users/rsc/g/go/include/link.h", "liblink/link_h.go", cfg)
	writeFile(prog, "/Users/rsc/g/go/src/pkg/runtime/stack.h", "liblink/stack_h.go", cfg)
	//	writeFile(prog, "/Users/rsc/g/go/src/cmd/ld/textflag.h", "liblink/textflag_h.go", cfg)
	writeFile(prog, "/Users/rsc/g/go/src/cmd/5l/5.out.h", "liblink/5.out.go", cfg)
	writeFile(prog, "/Users/rsc/g/go/src/cmd/6l/6.out.h", "liblink/6.out.go", cfg)
	writeFile(prog, "/Users/rsc/g/go/src/cmd/8l/8.out.h", "liblink/8.out.go", cfg)

	ioutil.WriteFile(filepath.Join(*out, "liblink/zzz.go"), []byte(zzzExtra), 0666)
}

var zzzExtra = `
package main

type Rune rune
type va_list struct{}
func fmtinstall(rune, func(*Fmt)int)
func sprint([]byte, string, ...interface{}) int
func sysfatal(string, ...interface{})
func snprint([]byte, int, string, ...interface{}) int
func cleanname(string) string
func Bprint(*Biobuf, string, ...interface{})
func sizeof(x interface{}) int

`

func writeFile(prog *cc.Prog, file, dstfile string, cfg *Config) {
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
		if cfg.Delete[decl.Name] {
			p.Print(decl.Comments.Before)
			p.Print(decl.Comments.Suffix, decl.Comments.After)
			continue
		}
		off := len(p.Bytes())
		if f, ok := cfg.Replace[decl.Name]; ok {
			p.Print(decl.Comments.Before)
			p.Print(f)
			p.Print(decl.Comments.Suffix, decl.Comments.After)
		} else {
			p.Print(decl)
		}
		if len(p.Bytes()) > off {
			p.Print(c2go.Newline)
			p.Print(c2go.Newline)
		}
	}
	buf := p.Bytes()
	buf1, err := format.Source(p.Bytes())
	if err == nil {
		buf = buf1
	}
	out := string(buf)
	for i, d := range cfg.Diffs {
		if strings.Contains(out, d.Before) {
			out = strings.Replace(out, d.Before, d.After, -1)
			cfg.Diffs[i].Used++
		}
	}
	if err := os.MkdirAll(filepath.Dir(dstfile), 0777); err != nil {
		log.Print(err)
	}
	if err := ioutil.WriteFile(dstfile, []byte(out), 0666); err != nil {
		log.Print(err)
	}
}
