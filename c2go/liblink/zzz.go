package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func cputime() float64 {
	return 0
}

func sysfatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func getgoroot() string {
	return "/Users/rsc/g/go"
}

func getgoos() string {
	return "darwin"
}

func getgoarch() string {
	return "amd64"
}

func getgoarm() string {
	return "7"
}

type Biobuf struct {
	r       *bufio.Reader
	w       *bufio.Writer
	written int64
}

func Bprint(b *Biobuf, format string, args ...interface{}) {
	n, _ := fmt.Fprintf(b.w, format, args...)
	b.written += int64(n)
}

func Bflush(b *Biobuf) {
	b.w.Flush()
}

func Boffset(b *Biobuf) int64 {
	return b.written
}

func Bread(b *Biobuf, buf []byte) int {
	n, err := io.ReadFull(b.r, buf)
	if err != nil && n == 0 {
		n = -1
	}
	return n
}

func Bgetc(b *Biobuf) int {
	c, err := b.r.ReadByte()
	if err != nil {
		return -1
	}
	return int(c)
}

func Bungetc(b *Biobuf) {
	b.r.UnreadByte()
}

func Bputc(b *Biobuf, c int) {
	err := b.w.WriteByte(byte(c))
	if err == nil {
		b.written++
	}
}

func Bwrite(b *Biobuf, buf []byte) int {
	n, err := b.w.Write(buf)
	b.written += int64(n)
	if err != nil && n == 0 {
		n = -1
	}
	return n
}

func getgoversion() string {
	return "gorsc"
}

func Binitw(f io.Writer) *Biobuf {
	return &Biobuf{w: bufio.NewWriter(f)}
}
