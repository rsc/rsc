// +build ignore

package main

import (
	"fmt"
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
func write(prog *cc.Prog, files []string) {
	for _, file := range files {
		writeFile(prog, file, "")
	}
	writeFile(prog, "/Users/rsc/g/go/include/fmt.h", "liblink/fmt_h.go")
	writeFile(prog, "/Users/rsc/g/go/include/bio.h", "liblink/bio_h.go")
	writeFile(prog, "/Users/rsc/g/go/include/link.h", "liblink/link_h.go")
	writeFile(prog, "/Users/rsc/g/go/src/cmd/5l/5.out.h", "liblink/5.out.go")
	writeFile(prog, "/Users/rsc/g/go/src/cmd/6l/6.out.h", "liblink/6.out.go")
	writeFile(prog, "/Users/rsc/g/go/src/cmd/8l/8.out.h", "liblink/8.out.go")

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

func writeFile(prog *cc.Prog, file, dstfile string) {
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
		off := len(p.Bytes())
		p.Print(decl)
		if len(p.Bytes()) > off {
			p.Print(c2go.Newline)
			p.Print(c2go.Newline)
		}
		if err := os.MkdirAll(filepath.Dir(dstfile), 0777); err != nil {
			log.Print(err)
		}
		if err := ioutil.WriteFile(dstfile, p.Bytes(), 0666); err != nil {
			log.Print(err)
		}
	}
}
