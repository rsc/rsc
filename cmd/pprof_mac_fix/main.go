// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pprof_mac_fix applies a binary patch to the OS X kernel in order to make
// pprof profiling report accurate values.
//
// Warning Warning Warning
//
// This program is meant to modify the operating system kernel, the program
// that runs your computer and makes it safe for all the other programs to run.
// If you damage the kernel, your computer will not be able to boot.
//
// Before using this program, ensure you can boot into ``recovery mode.''
// Many recent Macs make this possible by holding down Alt/Option when you
// hear the boot chime and selecting the ``Recovery HD.'' Otherwise, you can boot
// to the opening screen of an install DVD or thumb drive.
//
// You have been warned.
//
// Compatibility
//
// This program has been used successfully on the following systems:
//
//	OS X 10.6 Snow Leopard   / Darwin 10.8 / i386 only
//	OS X 10.7 Lion           / Darwin 11.4 / x86_64 only
//	OS X 10.8 Mountain Lion  / Darwin 12.4 / x86_64 only
//
// Snow Leopard x86_64 may work too but is untried.
//
// Installation
//
// First, read the warning above.
//
// Next, install this program and run it to create a modified kernel in /tmp:
//
//	go get code.google.com/p/rsc/cmd/pprof_mac_fix
//	pprof_mac_fix /mach_kernel /tmp/kernel
//
// Next, as root (sudo sh), make a backup of the standard kernel and then
// install the new one.
//
//	cp /mach_kernel /mach_kernel0 # only the first time!
//	cat /tmp/kernel >/mach_kernel
//
// Using cat instead of cp to install the kernel preserves the extended
// attributes of the original kernel, on the off chance that they matter.
//
// Finally, cross your fingers and reboot.
//
// If all goes well, running ``uname -a'' will report the time at which you
// ran pprof_mac_fix as the kernel build time.
//
// If you have a Go tree built at tip,
//
//	go test -v runtime/pprof
//
// should now say that the CPU profiling tests pass, whereas before they
// printed failure messages and were marked as skipped.
//
// Recovery
//
// If something goes wrong, you will need to restore the original kernel.
// To do this, boot into recovery mode.
// If you are using FileVault whole-disk encryption, start Disk Utility, unlock the disk,
// and then quit disk utility.
// (Disk Utility may be an option shown on the recovery mode screen or you may have to
// select it from the Utilities menu in the top-of-screen menu bar.)
// Start Terminal and then run these commands:
//
//	cd /Volumes/Mac*HD*
//	cp mach_kernel0 mach_kernel
//	bless /Volumes/Mac*HD*/System/Library/CoreServices
//
// The Mac*HD* pattern matches "Macintosh HD" and "Macintosh HD " [sic].
// If you have changed your disk's volume name you may need to use a
// different pattern (run "mount" to see the mounted disks).
//
// I am not sure whether the bless command is strictly necessary.
//
// Reboot. You should be back to the original, unmodified kernel.
// Either way, you need to be able to
// start Terminal and, if you are using FileVault whole-disk encryption, Disk Utility.
//
// For details on creating a bootable recovery disk or bootable installation disk,
// see http://support.apple.com/kb/HT4848 and http://lifehacker.com/5928780/.
//
// Theory of Operation
//
// The program rewrites the kernel code that delivers the profiling signals
// SIGPROF and SIGVTALRM in response to setitimer(2) calls.
// Instead of delivering the signal to the process as a whole,
// the new code delivers the signal to the thread whose execution
// triggered the signal; that is, it delivers the signal to the thread
// that is actually running and should be profiled.
//
// The rewrite only edits code in the function named bsd_ast, which is
// in charge of little more than delivering these signals.
// It is therefore unlikely to cause problems in programs not using the
// signals. Of course, there are no safety nets when changing an operating
// system kernel; caution is warranted.
//
package main

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var _ time.Time

var dumpFlag = flag.Bool("dump", false, "kernel dump")
var arch = flag.String("arch", getArch(), "arch to modify")

func getArch() string {
	out, _ := exec.Command("uname", "-m").CombinedOutput()
	switch s := strings.TrimSpace(string(out)); s {
	case "x86_64", "i386":
		return s
	}
	return "x86_64"
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [-arch ARCH] oldkernel newkernel\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "The -arch flag controls which kernel in a fat binary is modified.\n")
		fmt.Fprintf(os.Stderr, "The default setting is the architecture reported by `uname -m`,\n")
		fmt.Fprintf(os.Stderr, "on this machine %s.\n", getArch())
		os.Exit(2)
	}
	flag.Parse()
	if *dumpFlag {
		if flag.NArg() != 1 {
			fmt.Fprintf(os.Stderr, "usage: %s -dump oldkernel\n", os.Args[0])
			os.Exit(2)
		}
		dump(loadKernel(flag.Arg(0)))
		return
	}
	if flag.NArg() != 2 {
		flag.Usage()
	}

	k := loadKernel(flag.Arg(0))
	fmt.Printf("old: %s\n", k.version)

	errs := fixAnyVersion(k)
	if errs != nil {
		fmt.Fprintf(os.Stderr, "unrecognized kernel code.\n")
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		fmt.Fprintf(os.Stderr, updateText, os.Args[0], k.file)
		os.Exit(2)
	}

	// Update version string as displayed by uname -a.
	copy(k.timestamp, []byte(time.Now().Format(time.UnixDate)))
	fmt.Printf("new: %s\n", string(k.version))

	if err := ioutil.WriteFile(flag.Arg(1), k.data, 0666); err != nil {
		log.Fatal(err)
	}
}

func fixAnyVersion(k *kernel) []error {
	var errs []error
	for _, f := range fixes {
		err := f.apply(k.current_thread, k.bsd_ast)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %v", f.version, err))
			continue
		}
		return nil
	}
	return errs
}

var updateText = `
For an update, mail rsc@golang.org with the output printed by:
	%s -dump %s
`

type kernel struct {
	file           string
	data           []byte
	version        []byte
	timestamp      []byte
	current_thread []byte
	bsd_ast        []byte
}

type byValue []*macho.Symbol

func (x byValue) Len() int           { return len(x) }
func (x byValue) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byValue) Less(i, j int) bool { return x[i].Value < x[j].Value }

var versionRE = regexp.MustCompile(
	`Darwin Kernel Version [0-9]+\.[0-9]+\.[0-9]+: ` +
		`([A-Z][a-z][a-z] [A-Z][a-z][a-z] [ 1-9][0-9] \d{2}:\d{2}:\d{2} [A-Z]{3} \d{4});[^\0]*`,
)

type fatHeader struct {
	Magic   uint32
	NumArch uint32
	Entry   [4]struct {
		CPUType    uint32
		CPUSubType uint32
		Offset     uint32
		Size       uint32
		AlignBits  uint32
	}
}

func loadKernel(file string) *kernel {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	k := &kernel{
		file: file,
		data: data,
	}

	kdata := data

	var fat fatHeader
	binary.Read(bytes.NewReader(data), binary.BigEndian, &fat)
	if fat.Magic == 0xcafebabe {
		// It is a fat binary.
		n := int(fat.NumArch)
		if n > len(fat.Entry) {
			n = len(fat.Entry)
		}
		for i := range fat.Entry[:n] {
			e := &fat.Entry[i]
			switch {
			case *arch == "x86_64" && e.CPUType == 0x01000007 && e.CPUSubType == 0x00000003,
				*arch == "i386" && e.CPUType == 0x00000007 && e.CPUSubType == 0x00000003:
				fmt.Printf("%s(%s) is at offset %d, size %d\n", file, *arch, e.Offset, e.Size)
				kdata = data[e.Offset : e.Offset+e.Size]
				goto HaveKdata
			}	
		}
		log.Fatalf("cannot find %s kernel in fat kernel binary", *arch)
	HaveKdata:
	}

	if n := len(versionRE.FindAll(kdata, -1)); n == 0 {
		log.Fatalf("cannot find kernel version string")
	} else if n > 1 {
		log.Printf("warning: found multiple kernel version strings")
	}

	m := versionRE.FindSubmatchIndex(kdata)
	k.version = kdata[m[0]:m[1]]
	k.timestamp = kdata[m[2]:m[3]]

	// Look for current_thread body to make sure our inlining
	// of it is correct.
	f, err := macho.NewFile(bytes.NewReader(kdata))
	if err != nil {
		log.Fatal(err)
	}

	var syms []*macho.Symbol
	for i := range f.Symtab.Syms {
		syms = append(syms, &f.Symtab.Syms[i])
	}
	sort.Sort(byValue(syms))

	for i, sym := range syms {
		var save *[]byte
		switch sym.Name {
		case "_current_thread":
			save = &k.current_thread
		case "_bsd_ast":
			save = &k.bsd_ast
		}
		if save == nil {
			continue
		}
		sect := f.Sections[sym.Sect]
		off := int(sect.Offset) + int(sym.Value-sect.Addr)
		var n int
		if i == len(syms)-1 {
			n = int(sect.Addr + sect.Size - sym.Value)
		} else {
			n = int(syms[i+1].Value - sym.Value)
		}
		if off >= len(kdata) || off+n < off || off+n >= len(kdata) {
			log.Fatalf("invalid address [%d:%d] for %s in data [:%d]", off, off+n, sym.Name, len(kdata))
		}
		*save = kdata[off : off+n]
	}

	if k.current_thread == nil {
		log.Fatalf("cannot find current_thread in kernel")
	}
	if k.bsd_ast == nil {
		log.Fatalf("cannot find bsd_ast in kernel")
	}

	return k
}

func dump(k *kernel) {
	fmt.Printf("%s\nversion: %s\n", k.file, k.version)

	dumpDisas(k, k.current_thread, "current_thread")
	dumpDisas(k, k.bsd_ast, "bsd_ast")
}

var disasRE = regexp.MustCompile(`0x[0-9a-f]+\s+<\w+\+(\d+)>:`)

func dumpDisas(k *kernel, code []byte, name string) {
	cmd := exec.Command("gdb", "-arch", *arch, k.file)
	cmd.Stdin = strings.NewReader("disas " + name + "\n")
	disas, err := cmd.CombinedOutput()
	fmt.Printf("$ gdb %s # disas %s\n", k.file, name)
	if err != nil {
		fmt.Printf("running gdb 'disas %s': %v\n", name, err)
	}
	lines := strings.Split(string(disas), "\n")
	lastOff := -1
	flush := func(off int) {
		off1 := off
		if lastOff >= 0 && off1 < 0 {
			off1 = len(code)
		}
		if lastOff >= 0 && off1 > lastOff && off1 <= len(code) {
			n := off1 - lastOff
			if n > 20 {
				n = 20
			}
			fmt.Printf("\t% x\n", code[lastOff:lastOff+n])
		}
		lastOff = off
	}
	for _, line := range lines {
		m := disasRE.FindStringSubmatch(line)
		if m == nil {
			flush(-1)
		} else {
			n, _ := strconv.Atoi(m[1])
			flush(n)
		}
		fmt.Printf("%s\n", line)
	}
	flush(-1)
}

type pattern struct {
	mark    []int
	mask    []byte
	value   []byte
	leading []byte
}

var commentRE = regexp.MustCompile(`//[^\n]*`)

func mustCompile(text string) *pattern {
	p := new(pattern)
	text = commentRE.ReplaceAllString(text, "")
	for _, f := range strings.Fields(text) {
		if f == "*" {
			p.mark = append(p.mark, len(p.value))
			continue
		}
		val := f
		mask := "0xff"
		if i := strings.Index(f, "/"); i >= 0 {
			val, mask = f[:i], f[i+1:]
		}
		v, err := strconv.ParseUint(val, 0, 8)
		if err != nil {
			log.Fatalf("invalid value %s", f)
		}
		m, err := strconv.ParseUint(mask, 0, 8)
		if err != nil {
			log.Fatalf("invalid value %s", f)
		}
		p.value = append(p.value, byte(v))
		p.mask = append(p.mask, byte(m))
	}
	i := 0
	for i < len(p.mask) && p.mask[i] == 0xff {
		i++
	}
	p.leading = p.value[:i]
	return p
}

func (p *pattern) findAll(data []byte) []int {
	var out []int
	for i := 0; i < len(data); i++ {
		j := p.find(data[i:])
		if j < 0 {
			break
		}
		i += j
		out = append(out, i)
	}
	return out
}

func (p *pattern) find(data []byte) int {
	for i := 0; i < len(data); i++ {
		j := bytes.Index(data[i:], p.leading)
		if j < 0 {
			return -1
		}
		i += j
		if p.matchStart(data, i) != nil {
			return i
		}
	}
	return -1
}

func (p *pattern) matchStart(data []byte, off int) []int {
	sub := data[off:]
	if len(p.value) > len(sub) {
		return nil
	}
	for i := range p.value {
		if sub[i]&p.mask[i] != p.value[i] {
			return nil
		}
	}

	out := []int{}
	for _, m := range p.mark {
		out = append(out, off+m)
	}
	return out
}

type fix struct {
	version        string
	current_thread *pattern
	bsd_ast        []*pattern
}

var le = binary.LittleEndian

func (f *fix) apply(current_thread []byte, bsd_ast []byte) error {
	m := f.current_thread.matchStart(current_thread, 0)
	if m == nil {
		return fmt.Errorf("cannot match current_thread")
	}
	tlsOff := binary.LittleEndian.Uint32(current_thread[m[0]:])

	total := 0
	var timers [][]int
	for _, p := range f.bsd_ast {
		m := p.findAll(bsd_ast)
		total += len(m)
		timers = append(timers, m)
	}

	if total != 2 {
		if total == 0 {
			return fmt.Errorf("cannot match bsd_ast timer call")
		}
		if total == 1 {
			return fmt.Errorf("1 match for bsd_ast timer call %v, want 2", timers)
		}
		return fmt.Errorf("%d matches for bsd_ast timer call %v, want 2", total, timers)
	}

	var replace [][]byte
	for i, timer1 := range timers {
		for _, timer := range timer1 {
			p := f.bsd_ast[i]
			old := bsd_ast[timer:]
			m = p.matchStart(old, 0)
			if m == nil {
				// shouldn't happen - we found the offset above
				return fmt.Errorf("cannot match bsd_ast timer")
			}
			if !bytes.Equal(old[m[0]:m[1]], old[m[2]:m[3]]) {
				return fmt.Errorf("bsd_ast timer sequences differ")
			}
			if old[m[0]-2]&0xF8 != 0x70 {
				return fmt.Errorf("bsd_ast timer sequence missing conditional jump %x", old[m[0]-2])
			}
			if old[m[2]-2] != 0xeb {
				return fmt.Errorf("bsd_ast timer sequence missing unconditional jump %x", old[m[2]-2])
			}

			var new []byte
			new = append(new, old[m[0]:m[1]]...)
			new = append(new, old[:m[0]]...)
			// Last instruction is cond jump over call sequence.
			// We moved old[m[0]:m[1]] out,
			// so the jump must be shortened.
			new[len(new)-1] -= byte(m[1] - m[0])
			// "If" body.
			// The call instruction hasn't moved, so it's still correct.
			// The jmp at the end skips the else body,
			// so it must be shortened.
			new = append(new, old[m[1]:m[2]]...)
			new[len(new)-1] -= byte(m[1] - m[0])
			// "Else" body.
			// The call instruction has moved, so the offset must be adjusted.
			new = append(new, old[m[3]:m[4]]...)
			le.PutUint32(new[len(new)-4:], le.Uint32(new[len(new)-4:])-uint32(len(new)-m[4]))
			// Set up arguments to psignal_internal.
			if strings.Contains(f.version, "i386") {
				new = append(new,
					// xor %eax, %eax
					0x31, 0xc0,
					// xor %edx, %edx
					0x31, 0xd2,
					// mov %gs:threadTLS, %ecx
					0x65, 0x8b, 0x0d,
					byte(tlsOff), byte(tlsOff>>8), byte(tlsOff>>16), byte(tlsOff>>24),
					// mov $4, (%esp)
					0xc7, 0x04, 0x24, 0x04, 0x00, 0x00, 0x00,
					// mov $0x1a (or $0x1b), 4(%esp)
					0xc7, 0x44, 0x24, 0x04, old[m[5]], 0x00, 0x00, 0x00,
				)
			} else { // x86_64
				new = append(new,
					// xor %edi, %edi
					0x31, 0xff,
					// xor %esi, %esi
					0x31, 0xf6,
					// mov %gs:threadTLS, %rdx
					0x65, 0x48, 0x8b, 0x14, 0x25,
					byte(tlsOff), byte(tlsOff>>8), byte(tlsOff>>16), byte(tlsOff>>24),
					// mov $4, %ecx
					0xb9, 0x04, 0x00, 0x00, 0x00,
					// mov $0x1a (or $0x1b), %r8d
					0x41, 0xb8, old[m[5]], 0x00, 0x00, 0x00,
				)
			}
			for len(new) < m[6] {
				new = append(new, 0x90) // nop
			}
			if len(new) > m[6] {
				return fmt.Errorf("bsd_ast timer sequence rewrite too long")
			}
			replace = append(replace, new)
		}
	}

	// Commit rewrite.
	n := 0
	for _, timer1 := range timers {
		for _, timer := range timer1 {
			copy(bsd_ast[timer:], replace[n])
			n++
		}
	}

	return nil
}

// Darwin 10.8.0 (Snow Leopard)

var current_thread_leave = mustCompile(`
    0x55                            //  0   push %rbp
    0x48 0x89 0xe5                  //  1   mov %rsp, %rbp
    0x65 0x48 0x8b 0x04 0x25        //  4   mov %gs:0x8 %rax
    * 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00
    0xc9                            // 13   leaveq
    0xc3                            // 14   retq
`)

var current_thread_leave_i386 = mustCompile(`
    0x55                            //  0   push %rbp
    0x89 0xe5                       //  1   mov %rsp, %rbp
    0x65 0xa1                       //  3   mov %gs:0x8, %rax
    * 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00
    0xc9                            // 9   leaveq
    0xc3                            // 14   retq
`)

var bsd_ast_10_8_0_a = mustCompile(`
    0x49 0x83 0xbc 0x24 0x00/0x0f 0x01 0x00 0x00 0x00  //  0 cmpq $0x0,0x1b0(%r12)
    0x75 0x0c                                          //  9 jne +12 [23]
    0x41 0x8b 0x84 0x24 0x08/0x0f 0x01 0x00 0x00       // 11 mov 0x1b8(%r12),%eax
    0x85 0xc0                                          // 19 test %eax,%eax
    0x74 0x11                                          // 21 je +17 [40]
    * 0x49 0x8b 0x7c 0x24 0x18                         // 23 mov 0x18(%r12),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                    // 28 mov $0x1,%esi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00       // 33 call task_vtimer_set
    0xeb 0x0f                                          // 38 jmp +15 [55]
    * 0x49 0x8b 0x7c 0x24 0x18                         // 40 mov 0x18(%r12),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                    // 45 mov $0x1,%esi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00       // 50 call task_vtimer_clear
    * 0x41 0xb8 * 0x1a/0xfe 0x00 0x00 0x00             // 55 mov $0x1a,%r8d
    0x31 0xc9                                          // 61 xor %ecx,%ecx
    0x31 0xd2                                          // 63 xor %edx,%edx
    0x31 0xf6                                          // 65 xor %esi,%esi
    0x4c 0x89 0xe7 *                                   // 67 mov %r12,%rdi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00       // 70 call psignal_internal
`)

var bsd_ast_10_8_0_b = mustCompile(`
    0x49 0x83 0xbc 0x24 0x00/0x0f 0x01 0x00 0x00 0x00  //  0 cmpq $0x0,0x1d0(%r12)
    0x75 0x0d                                          //  9 jne +13
    0x45 0x8b 0x9c 0x24 0x08/0x0f 0x01 0x00 0x00       // 11 mov 0x1d8(%r12),%r11d
    0x45 0x85 0xdb                                     // 19 test %r11d,%r11d
    0x74 0x11                                          // 21 je +17
    * 0x49 0x8b 0x7c 0x24 0x18                         // 23 mov 0x18(%r12),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                    // 28 mov $0x2,%esi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00       // 33 call task_vtimer_set
    0xeb 0x0f                                          // 38 jmp +15
    * 0x49 0x8b 0x7c 0x24 0x18                         // 40 mov 0x18(%r12),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                    // 45 mov $0x2,%esi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00       // 50 call task_vtimer_clear
    * 0x41 0xb8 * 0x1a/0xfe 0x00 0x00 0x00             // 55 mov $0x1b,%r8d
    0x31 0xc9                                          // 61 xor %ecx,%ecx
    0x31 0xd2                                          // 63 xor %edx,%edx
    0x31 0xf6                                          // 65 xor %esi,%esi
    0x4c 0x89 0xe7 *                                   // 67 mov %r12,%rdi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00       // 70 call psignal_internal
`)

var bsd_ast_10_8_0_i386_a = mustCompile(`
    0x8b 0x8b 0xec 0x00 0x00 0x00                   //  0 mov    0xec(%ebx),%ecx
    0x85 0xc9                                       //  6 test   %ecx,%ecx
    0x75 0x0a                                       //  8 jne    +10
    0x8b 0x93 0xf0 0x00 0x00 0x00                   // 10 mov    0xf0(%ebx),%edx
    0x85 0xd2                                       // 16 test   %edx,%edx
    0x74 0x15                                       // 18 je     +21
    * 0x8b 0x43 0x0c                                // 20 mov    0xc(%ebx),%eax
    0xc7 0x44 0x24 0x04 0x01 0x00 0x00 0x00         // 23 movl   $0x1,0x4(%esp)
    0x89 0x04 0x24 *                                // 31 mov    %eax,(%esp)
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 34 call task_vtimer_set
    0xeb 0x13                                       // 39 jmp    +19
    * 0x8b 0x43 0x0c                                // 41 mov    0xc(%ebx),%eax
    0xc7 0x44 0x24 0x04 0x01 0x00 0x00 0x00         // 44 movl   $0x1,0x4(%esp)
    0x89 0x04 0x24 *                                // 52 mov    %eax,(%esp)
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 55 call task_vtimer_clear
    * 0xc7 0x44 0x24 0x04 * 0x1a 0x00 0x00 0x00     // 60 movl   $0x1a,0x4(%esp)
    0xc7 0x04 0x24 0x00 0x00 0x00 0x00              // 68 movl   $0x0,(%esp)
    0x31 0xc9                                       // 75 xor    %ecx,%ecx
    0x31 0xd2                                       // 77 xor    %edx,%edx
    0x89 0xd8 *                                     // 79 mov    %ebx,%eax
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 81 call psignal_internal
`)

var bsd_ast_10_8_0_i386_b = mustCompile(`
    0x8b 0x83 0xfc 0x00 0x00 0x00                   //  0 mov    0xfc(%ebx),%eax
    0x85 0xc0                                       //  6 test   %eax,%eax
    0x75 0x0a                                       //  8 jne    +10
    0x8b 0x83 0x00 0x01 0x00 0x00                   // 10 mov    0x100(%ebx),%eax
    0x85 0xc0                                       // 16 test   %eax,%eax
    0x74 0x15                                       // 18 je     +21
    * 0x8b 0x43 0x0c                                // 20 mov    0xc(%ebx),%eax
    0xc7 0x44 0x24 0x04 0x02 0x00 0x00 0x00         // 23 movl   $0x2,0x4(%esp)
    0x89 0x04 0x24 *                                // 31 mov    %eax,(%esp)
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 34 call task_vtimer_set
    0xeb 0x13                                       // 39 jmp    +19
    * 0x8b 0x43 0x0c                                // 41 mov    0xc(%ebx),%eax
    0xc7 0x44 0x24 0x04 0x02 0x00 0x00 0x00         // 44 movl   $0x2,0x4(%esp)
    0x89 0x04 0x24 *                                // 52 mov    %eax,(%esp)
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 55 call task_vtimer_clear
    * 0xc7 0x44 0x24 0x04 * 0x1b 0x00 0x00 0x00     // 60 movl   $0x1b,0x4(%esp)
    0xc7 0x04 0x24 0x00 0x00 0x00 0x00              // 68 movl   $0x0,(%esp)
    0x31 0xc9                                       // 75 xor    %ecx,%ecx
    0x31 0xd2                                       // 77 xor    %edx,%edx
    0x89 0xd8 *                                     // 79 mov    %ebx,%eax
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 81 call psignal_internal
`)

var fix_10_8_0 = fix{
	"10.8.0",
	current_thread_leave,
	[]*pattern{bsd_ast_10_8_0_a, bsd_ast_10_8_0_b},
}

var fix_10_8_0_i386 = fix{
	"10.8.0 (i386)",
	current_thread_leave_i386,
	[]*pattern{bsd_ast_10_8_0_i386_a, bsd_ast_10_8_0_i386_b},
}

// Darwin 11.4.2 (Lion)

var current_thread_pop = mustCompile(`
    0x55                            //  0   push %rbp
    0x48 0x89 0xe5                  //  1   mov %rsp, %rbp
    0x65 0x48 0x8b 0x04 0x25        //  4   mov %gs:0x8 %rax
    * 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00
    0x5d                            // 13   pop %rbp
    0xc3                            // 14   retq
    0x90                            // 15   nop
`)

var bsd_ast_11_4_2 = mustCompile(`
    0x49 0x83 0xbe 0xc0/0xdf 0x01 0x00 0x00 0x00    //  0 cmpq   $0x0,0x1c0(%r14)
    0x75 0x0a                                       //  8 jne    +10
    0x41 0x83 0xbe 0xc8/0xdf 0x01 0x00 0x00 0x00    // 10 cmpl   $0x0,0x1c8(%r14)
    0x74 0x10                                       // 18 je     +16
    * 0x49 0x8b 0x7e 0x18                           // 20 mov    0x18(%r14),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                 // 24 mov    $0x1,%esi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 29 call task_vtimer_set
    0xeb 0x0e                                       // 34 jmp    +15
    * 0x49 0x8b 0x7e 0x18                           // 36 mov    0x18(%r14),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                 // 40 mov    $0x1,%esi
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 45 call task_vtimer_clear
    * 0x31 0xf6                                     // 50 xor    %esi,%esi
    0x31 0xc9                                       // 52 xor    %ecx,%ecx
    0x41 0xb8 * 0x1a/0xfe 0x00 0x00 0x00            // 54 mov    $0x1a,%r8d
    0x4c 0x89 0xf7                                  // 60 mov    %r14,%rdi
    0x31 0xd2 *                                     // 63 xor    %edx,%edx
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 65 call psignal_internal
`)

var fix_11_4_2 = fix{
	"11.4.2",
	current_thread_pop,
	[]*pattern{bsd_ast_11_4_2},
}

// Darwin 12.4.0 (Mountain Lion)

var fix_12_4_0 = fix{
	"12.4.0",
	current_thread_pop,
	[]*pattern{bsd_ast_12_4_0},
}

var bsd_ast_12_4_0 = mustCompile(`
    0x49 0x83 0xbf 0xc0/0xdf 0x01 0x00 0x00 0x00    //  0   cmpq $0x0, 0x1c0(%r15) [or 0x1e0]
    0x75 0x0a                                       //  8   jne +10 [20]
    0x41 0x83 0xbf 0xc8/0xdf 0x01 0x00 0x00 0x00    // 10   cmpl $0x0, 0x1c8(%r15) [or 0x1e8]
    0x74 0x10                                       // 18   je +16 [36]
    * 0x49 0x8b 0x7f 0x18                           // 20   mov 0x18(%r15),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                 // 24   mov $0x1, %esi [or $0x2]
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 29   call task_vtimer_set
    0xeb 0x0e                                       // 34   jmp +14 [50]
    * 0x49 0x8b 0x7f 0x18                           // 36   mov 0x18(%r15),%rdi
    0xbe 0x00/0xfc 0x00 0x00 0x00 *                 // 40   mov $0x1, %esi [or $0x2]
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 45   call task_vtimer_clear
    * 0x4c 0x89 0xff                                // 50   mov %r15, %rdi
    0x31 0xf6                                       // 53   xor %esi, %esi
    0x31 0xd2                                       // 55   xor %edx, %edx
    0x31 0xc9                                       // 57   xor %ecx, %ecx
    0x41 0xb8 * 0x1a/0xfe 0x00 0x00 0x00 *          // 59   mov $0x1a, %r8d [or $0x1b]
    0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00    // 65   call psignal_internal
`)

var fixes = []*fix{
	&fix_10_8_0,
	&fix_10_8_0_i386,
	&fix_11_4_2,
	&fix_12_4_0,
}
