package main

import (
	"bufio"
	"bytes"
	"exec"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strings"
	"time"

	"rsc.googlecode.com/hg/google"
	"rsc.googlecode.com/hg/imap"
)

var cmdtab = []struct{
	Name string
	Args int
	F func(*Cmd, *imap.MsgPart) *imap.MsgPart
	Help string
}{
//	{ "a",	1,	acmd,	"a        reply to sender and recipients" },
//	{ "A",	1,	acmd,	"A        reply to sender and recipients with copy" },
	{ "b",	0,	bcmd,	"b        print the next 10 headers" },
	{ "d",	0,	dcmd,	"d        mark for deletion" },
//	{ "f",	0,	fcmd,	"f        file message by from address" },
	{ "h",	0,	hcmd,	"h        print elided message summary (,h for all)" },
	{ "help", 0,	nil, "help     print this info" },
	{ "H",	0,	Hcmd,	"H        print message's MIME structure " },
	{ "i",	0,	icmd,	"i        incorporate new mail" },
//	{ "k",	0,	kcmd,	"k        kill (mute) mail" },
//	{ "m",	1,	mcmd,	"m addr   forward mail" },
//	{ "M",	1,	mcmd,	"M addr   forward mail with message" },
	{ "p",	0,	pcmd,	"p        print the processed message" },
//	{ "p+",	0,	pcmd,	"p        print the processed message, showing all quoted text" },
	{ "P",	0,	Pcmd,	"P        print the raw message" },
	{ `"`,	0,	quotecmd, "\"        print a quoted version of msg" },
	{ "q",	0,	qcmd,	"q        exit and remove all deleted mail" },
	{ "r",	1,	rcmd,	"r [addr] reply to sender plus any addrs specified" },
//	{ "rf",	1,	rcmd,	"rf [addr]file message and reply" },
//	{ "R",	1,	rcmd,	"R [addr] reply including copy of message" },
//	{ "Rf",	1,	rcmd,	"Rf [addr]file message and reply with copy" },
//	{ "s",	1,	scmd,	"s file   append raw message to file" },
	{ "u",	0,	ucmd,	"u        remove deletion mark" },
//	{ "w",	1,	wcmd,	"w file   store message contents as file" },
	{ "W",	0,	Wcmd,	"W	open in web browser" },
	{ "x",	0,	xcmd,	"x        exit without flushing deleted messages" },
	{ "y",	0,	ycmd,	"y        synchronize with mail box" },
	{ "=",	1,	eqcmd,	"=        print current message number" },
//	{ "|",	1,	pipecmd, "|cmd     pipe message body to a command" },
//	{ "||",	1,	rpipecmd, "||cmd     pipe raw message to a command" },
//	{ "!",	1,	bangcmd, "!cmd     run a command" },
}

func init() {
	// Have to insert helpcmd by hand because it refers to cmdtab,
	// so it would cause an init loop above.
	for i := range cmdtab {
		if cmdtab[i].Name == "help" {
			cmdtab[i].F = helpcmd
		}
	}
}

type Cmd struct {
	Name string
	Args []string
	F func(*Cmd, *imap.MsgPart) *imap.MsgPart
	Delete bool
	Thread bool
	Targ *imap.MsgPart
	Targs []*imap.Msg
	A1, A2 int
}

var (
	bin = bufio.NewReader(os.Stdin)
	bout = bufio.NewWriter(os.Stdout)
	
	acctName = flag.String("a", "", "account to use")
	
	dot *imap.MsgPart  // Selected messages
	
	inbox *imap.Box
	msgs []*imap.Msg
	msgNum = make(map[*imap.Msg]int)
	deleted = make(map[*imap.Msg]bool)
	isGmail = false
	acct google.Account
	threaded bool
	interrupted bool

	maxfrom int
	subjlen int
)

func nextMsg(m *imap.Msg) *imap.Msg {
	i := msgNum[m]
	i++
	if i >= len(msgs) {
		return nil
	}
	return msgs[i]
}

func main() {
	flag.BoolVar(&imap.Debug, "imapdebug", false, "imap debugging trace")
	flag.Parse()

	acct = google.Acct(*acctName)
	c, err := imap.NewClient(imap.TLS, "imap.gmail.com", acct.Email, acct.Password, "")
	if err != nil {
		log.Fatal(err)
	}
	isGmail = c.IsGmail()
	threaded = isGmail
	
	inbox = c.Inbox()
	if err := inbox.Check(); err != nil {
		log.Fatal(err)
	}

	msgs = inbox.Msgs()
	maxfrom = 12
	for i, m := range msgs {
		msgNum[m] = i
		if n := len(from(m.Hdr)); n > maxfrom {
			maxfrom = n
		}
	}
	if maxfrom > 20 {
		maxfrom = 20
	}
	subjlen = 80 - maxfrom
	
	rethread()
	
	go func() {
		for sig := range signal.Incoming {
			if sig == os.SIGINT {
				fmt.Fprintf(os.Stderr, "!interrupt\n")
				interrupted = true
				continue
			}
			if sig == os.SIGCHLD {
				continue
			}
			fmt.Fprintf(os.Stderr, "!%s\n", sig)
		}
	}()
			
	for {
		if dot != nil {
			fmt.Fprintf(bout, "%d", msgNum[dot.Msg]+1)
			if dot != &dot.Msg.Root {
				fmt.Fprintf(bout, ".%s", dot.ID)
			}
		}
		fmt.Fprintf(bout, ": ")
		bout.Flush()
		
		line, err := bin.ReadString('\n')
		if err != nil {
			break
		}
		
		cmd, err := parsecmd(line)
		if err != nil {
			fmt.Fprintf(bout, "!%s\n", err)
			continue
		}

		if cmd.Targ != nil || cmd.Targs == nil && cmd.A2 == 0 {
			x := cmd.F(cmd, cmd.Targ)
			if x != nil {
				dot = x
			}
		} else {
			targs := cmd.Targs
			if targs == nil {
				delta := +1
				if cmd.A1 > cmd.A2 {
					delta = -1
				}
				for i := cmd.A1; i <= cmd.A2; i += delta {
					if i < 1 || i > len(msgs) {
						continue
					}
					targs = append(targs, msgs[i-1])
				}
			}
			if cmd.Thread {
				if !isGmail {
					fmt.Fprintf(bout, "!need gmail for threaded command\n")
					continue
				}
				byThread := make(map[uint64][]*imap.Msg)
				for _, m := range msgs {
					t := m.GmailThread
					byThread[t] = append(byThread[t], m)
				}
				for _, m := range targs {
					t := m.GmailThread
					for _, mm := range byThread[t] {
						x := cmd.F(nil, &mm.Root)
						if x != nil {
							dot = x
						}
					}
					byThread[t] = nil, false
				}
				continue
			}
			for _, m := range targs {
				if cmd.Delete {
					dcmd(cmd, &m.Root)
					if cmd.Name == "p" {
						// dp is a special case: it advances to the next message before the p.
						next := nextMsg(m)
						if next == nil {
							fmt.Fprintf(bout, "!address\n")
							dot = &m.Root
							break
						}
						m = next
					}
				}
				x := cmd.F(nil, &m.Root)
				if x != nil {
					dot = x
				}
				// TODO: Break loop on interrupt.
			}
		}
	}
	qcmd(nil, nil)
}

func parsecmd(line string) (cmd *Cmd, err os.Error) {
	cmd = &Cmd{}
	line = strings.TrimSpace(line)
	if line == "" {
		// Empty command is a special case: advance and print.
		cmd.F = pcmd
		if dot == nil {
			cmd.A1 = 1
			cmd.A2 = 1
		} else {
			n := msgNum[dot.Msg]+2
			if n > len(msgs) {
				return nil, fmt.Errorf("out of messages")
			}
			cmd.A1 = n
			cmd.A2 = n
		}
		return cmd, nil
	}
	
	// Global search?
	if line[0] == 'g' {
		line = line[1:]
		if line == "" || line[0] != '/' {
			// No search string means all messages.
			cmd.A1 = 1
			cmd.A2 = len(msgs)
		} else if line[0] == '/' {
			re, rest, err := parsere(line)
			if err != nil {
				return nil, err
			}
			line = rest
			// Find all messages matching this search string.
			var targ []*imap.Msg
			for _, m := range msgs {
				if re.MatchString(header(m)) {
					targ = append(targ, m)
				}
			}
			if len(targ) == 0 {
				return nil, fmt.Errorf("no matches")
			}
			cmd.Targs = targ
		}
	} else {
		// Parse an address.
		a1, targ, rest, err := parseaddr(line, 1)
		if err != nil {
			return nil, err
		}
		if targ != nil {
			cmd.Targ = targ
			line = rest
		} else {
			if a1 < 1 || a1 > len(msgs) {
				return nil, fmt.Errorf("message number %d out of range", a1)
			}
			cmd.A1 = a1
			cmd.A2 = a1
			a2 := a1
			if rest != "" && rest[0] == ',' {
				// This is an address range.
				a2, targ, rest, err = parseaddr(rest[1:], len(msgs))
				if err != nil {
					return nil, err
				}
				if a2 < 1 || a2 > len(msgs) {
					return nil, fmt.Errorf("message number %d out of range", a2)
				}
				cmd.A2 = a2
			} else if rest == line {
				// There was no address.
				if dot == nil {
					cmd.A1 = 1
					cmd.A2 = 0
				} else {
					if dot != nil {
						if dot == &dot.Msg.Root {
							// If dot is a plain msg, use a range so that dp works.
							cmd.A1 = msgNum[dot.Msg]+1
							cmd.A2 = cmd.A1
						} else {
							cmd.Targ = dot
						}
					}
				}
			}
			line = rest
		}
	}

	// Insert space after ! or | for tokenization.
	for j := 0; j < len(line); j++ {
		if line[j] == '!' || line[j] == '|' {
			line = line[:j+1] + " " + line[j+1:]
			break
		}
	}
	av := strings.Fields(strings.TrimSpace(line))
	cmd.Args = av
	if len(av) == 0 || av[0] == "" {
		// Default is to print.
		cmd.F = pcmd
		return cmd, nil
	}

	name := av[0]

	// Hack to allow t prefix on all commands.
	if len(name) >= 2 && name[0] == 't' {
		cmd.Thread = true
		name = name[1:]
	}

	// Hack to allow d prefix on all commands.
	if len(name) >= 2 && name[0] == 'd' {
		cmd.Delete = true
		name = name[1:]
	}
	cmd.Name = name

	// Search command table.
	for _, ct := range cmdtab {
		if ct.Name == name {
			if ct.Args == 0 && len(av) > 1 {
				return nil, fmt.Errorf("%s doesn't take an argument", name)
			}
			cmd.F = ct.F
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("unknown command %s", name)
}

func parseaddr(addr string, deflt int) (n int, targ *imap.MsgPart, rest string, err os.Error) {
	dot := dot
	n = deflt
	for {
		old := addr
		n, targ, rest, err = parseaddr1(addr, n, dot)
		if targ != nil || rest == old || err != nil {
			break
		}
		if n < 1 || n > len(msgs) {
			return 0, nil, "", fmt.Errorf("message number %d out of range", n)
		}
		dot = &msgs[n-1].Root
		addr = rest
	}
	return
}

func parseaddr1(addr string, deflt int, dot *imap.MsgPart) (n int, targ *imap.MsgPart, rest string, err os.Error) {
	base := 0
	if dot != nil {
		base = msgNum[dot.Msg] + 1
	}
	if addr == "" {
		return deflt, nil, addr, nil
	}
	var i int
	sign := 0
	switch c := addr[0]; c {
	case '+':
		sign = +1
		addr = addr[1:]
	case '-':
		sign = -1
		addr = addr[1:]
	case '.':
		if base == 0 {
			return 0, nil, "", fmt.Errorf("no message selected")
		}
		n = base
		i = 1
		goto HaveNumber
	case '$':
		if len(msgs) == 0 {
			return 0, nil, "", fmt.Errorf("no messages")
		}
		n = len(msgs)
		i = 1
		goto HaveNumber
	case '/', '?':
		var re *regexp.Regexp
		re, addr, err = parsere(addr)
		if err != nil {
			return
		}
		var delta int
		if c == '/' {
			delta = +1
		} else {
			delta = -1
		}
		for j := base+delta; 1 <= j && j <= len(msgs); j += delta {
			if re.MatchString(header(msgs[j-1])) {
				n = j
				i = 0  // already cut addr
				goto HaveNumber
			}
		}
		err = fmt.Errorf("search")
		return
	// TODO case '%'
	}
	for i = 0; i < len(addr) && '0' <= addr[i] && addr[i] <= '9'; i++ {
		n = 10*n + int(addr[i]) - '0'
	}
	if sign != 0 {
		if n == 0 {
			n = 1
		}
		n = base+n*sign
		goto HaveNumber
	}
	if i == 0 {
		return deflt, nil, addr, nil
	}
HaveNumber:
	rest = addr[i:]
	if i < len(addr) && addr[i] == '.' {
		if n < 1 || n > len(msgs) {
			err = fmt.Errorf("message number %d out of range", n)
			return
		}
		targ = &msgs[n-1].Root
		for i < len(addr) && addr[i] == '.' {
			i++
			var j int
			n = 0
			for j = i; j < len(addr) && '0' <= addr[j] && addr[j] <= '9'; j++ {
				n = 10*n + int(addr[j]) - '0'
			}
			if j == i {
				err = fmt.Errorf("malformed message number %s", addr[:j])
				return
			}
			if n < 1 || n > len(targ.Child) {
				err = fmt.Errorf("message number %s out of range", addr[:j])
				return
			}
			targ = targ.Child[n-1]
			i = j
		}
		n = 0
		rest = addr[i:]
		return
	}
	return
}

func parsere(addr string) (re *regexp.Regexp, rest string, err os.Error) {
	prog, rest, err := parseprog(addr)
	if err != nil {
		return
	}
	re, err = regexp.Compile(prog)
	return
}

var lastProg string

func parseprog(addr string) (prog string, rest string, err os.Error) {
	if len(addr) == 1 {
		if lastProg != "" {
			return lastProg, "", nil
		}
		err = fmt.Errorf("no search")
		return
	}
	i := strings.Index(addr[1:], addr[:1])
	if i < 0 {
		prog = addr[1:]
		rest = ""
	} else {
		i += 1  // adjust for slice in IndexByte arg
		prog, rest = addr[1:i], addr[i+1:]
	}
	lastProg = prog
	return
}

func bcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	var m *imap.Msg
	if dot == nil {
		if len(msgs) == 0 {
			return nil
		}
		m = msgs[0]
	} else {
		m = dot.Msg
	}
	for i := 0; i < 10; i++ {
		hcmd(c, &m.Root)
		next := nextMsg(m)
		if next == nil {
			break
		}
		m = next
	}
	return &m.Root
}

func dcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		fmt.Fprintf(bout, "!address\n")
		return nil
	}
	deleted[dot.Msg] = true
	return &dot.Msg.Root
}

func ucmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		fmt.Fprintf(bout, "!address\n")
		return nil
	}
	deleted[dot.Msg] = false, false
	return &dot.Msg.Root
}

func eqcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		fmt.Fprintf(bout, "0")
	} else {
		fmt.Fprintf(bout, "%d", msgNum[dot.Msg]+1)
		if dot != &dot.Msg.Root {
			fmt.Fprintf(bout, ".%s", dot.ID)
		}
	}
	fmt.Fprintf(bout, "\n")
	return nil
}

func from(h *imap.MsgHdr) string {
	if len(h.From) < 1 {
		return "?"
	}
	if name := h.From[0].Name; name != "" {
		return name
	}
	return h.From[0].Email
}

func header(m *imap.Msg) string {
	var t string
	if m.Date >= time.Seconds()-86400*365 {
		t = time.SecondsToLocalTime(m.Date).Format("01/02 15:04")
	} else {
		t = time.SecondsToLocalTime(m.Date).Format("01/02 2006 ")
	}
	ch := ' '
	if len(m.Root.Child) > 1 || len(m.Root.Child) == 1 && len(m.Root.Child[0].Child) > 0 {
		ch = 'H'
	}
	del := ' '
	if deleted[m] {
		del = 'd'
	}
	return fmt.Sprintf("%-3d %c%c %s %-*.*s %.*s",
		msgNum[m]+1, ch, del, t,
		maxfrom, maxfrom, from(m.Hdr),
		subjlen, m.Hdr.Subject)
}

func hcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot != nil {
		fmt.Fprintf(bout, "%s\n", header(dot.Msg))
	}
	return nil
}

func helpcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	fmt.Fprint(bout, "Commands are of the form [<range>] <command> [args]\n");
	fmt.Fprint(bout, "<range> := <addr> | <addr>','<addr>| 'g'<search>\n");
	fmt.Fprint(bout, "<addr> := '.' | '$' | '^' | <number> | <search> | <addr>'+'<addr> | <addr>'-'<addr>\n");
	fmt.Fprint(bout, "<search> := '/'<gmail search>'/' | '?'<gmail search>'?'\n");
	fmt.Fprint(bout, "<command> :=\n");
	for _, ct := range cmdtab {
		fmt.Fprintf(bout, "%s\n", ct.Help);
	}
	return dot
}

func Hcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot != nil {
		H(fmt.Sprint(msgNum[dot.Msg]+1), dot)
	}
	return nil
}

func H(id string, p *imap.MsgPart) {
	if p.ID != "" {
		id = id + "." + p.ID
	}
	fmt.Fprintf(bout, "%s %s %s %#q %d\n", id, p.Type, p.Encoding+"/"+p.Charset, p.Name, p.Bytes)
	for _, child := range p.Child {
		H(id, child)
	}
}

func icmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	sync(false)
	return nil
}

func ycmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	sync(true)
	return nil
}

func addrlist(x []imap.Addr) string {
	var b bytes.Buffer
	for i, a := range x {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(a.String())
	}
	return b.String()
}

func pcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		return nil
	}
	if dot == &dot.Msg.Root {
		h := dot.Msg.Hdr
		if len(h.From) > 0 {
			fmt.Fprintf(bout, "From: %s\n", addrlist(h.From))
		}
		fmt.Fprintf(bout, "Date: %s\n", time.SecondsToLocalTime(dot.Msg.Date))
		if len(h.From) > 0 {
			fmt.Fprintf(bout, "To: %s\n", addrlist(h.To))
		}
		if len(h.CC) > 0 {
			fmt.Fprintf(bout, "CC: %s\n", addrlist(h.CC))
		}
		if len(h.BCC) > 0 {
			fmt.Fprintf(bout, "BCC: %s\n", addrlist(h.BCC))
		}
		if len(h.Subject) > 0 {
			fmt.Fprintf(bout, "Subject: %s\n", h.Subject)
		}
		fmt.Fprintf(bout, "\n")
	}
	printMIME(dot, true, false)
	fmt.Fprintf(bout, "\n")
	return dot
}

func unixfrom(h *imap.MsgHdr) string {
	if len(h.From) == 0 {
		return ""
	}
	return h.From[0].Email
}

func unixtime(m *imap.Msg) string {
	return time.SecondsToLocalTime(dot.Msg.Date).Format("Mon Jan _2 15:04:05 MST 2006")
}


func Pcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		return nil
	}
	if dot == &dot.Msg.Root {
		fmt.Fprintf(bout, "From %s %s\n",
			unixfrom(dot.Msg.Hdr),
			unixtime(dot.Msg))
	}
	bout.Write(dot.Raw())
	return dot
}

func printMIME(p *imap.MsgPart, top, quote bool) {
	if top && strings.HasPrefix(p.Type, "text/") {
		text := p.Text()
		if quote {
			text = append([]byte("> "), bytes.Replace(text, []byte("\n"), []byte("\n> "), -1)...)
			if bytes.HasSuffix(text, []byte("\n> ")) {
				text = text[:len(text)-2]
			}
		}
		bout.Write(text)
		return
	}
	switch p.Type {
	case "text/plain":
		if top {
			panic("printMIME loop")
		}
		printMIME(p, true, quote)
	case "multipart/alternative":
		for _, pp := range p.Child {
			if pp.Type == "text/plain" {
				printMIME(pp, false, quote)
				return
			}
		}
		if len(p.Child) > 0 {
			printMIME(p.Child[0], false, quote)
		}
	case "multipart/mixed":
		for _, pp := range p.Child {
			printMIME(pp, false, quote)
		}
	default:
		fmt.Fprintf(bout, "%d.%s !%s %s %s\n", msgNum[p.Msg]+1, p.ID, p.Type, p.Desc, p.Name)
	}
}

func qcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	flushDeleted()
	xcmd(c, dot)
	panic("not reached")
}

func quotecmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		return nil
	}
	m := dot.Msg
	if len(m.Hdr.From) != 0 {
		a := m.Hdr.From[0]
		name := a.Name
		if name == "" {
			name = a.Email
		}
		date := time.SecondsToLocalTime(m.Date).Format("Jan 2, 2006 at 15:04")
		fmt.Fprintf(bout, "On %s, %s wrote:\n", date, name)
	}
	printMIME(dot, true, true)
	return dot
}

func addre(s string) string {
	if len(s) < 4 || !strings.EqualFold(s[:4], "re: ") {
		return "Re: " + s
	}
	return s
}

func rcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		fmt.Fprintf(bout, "!nothing to reply to\n")
		return nil
	}

	h := dot.Msg.Hdr
	reply := h.ReplyTo
	if len(reply) == 0 {
		reply = h.From
	}
	if h == nil || len(reply) == 0 {
		fmt.Fprintf(bout, "!no one to reply to\n")
		return dot
	}
	
	args := []string{"-a", acct.Email, "-s", addre(h.Subject), "-in-reply-to", h.MessageID}
	for _, a := range h.ReplyTo {
		args = append(args, "-to", a.String())
	}
	bout.Flush()
	cmd := exec.Command("gmailsend", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(bout, "!%s\n", err)
	}
	return dot
}

func rethread() {
	if !threaded {
		sort.Sort(byUIDRev(msgs))
	} else {
		byThread := make(map[uint64][]*imap.Msg)
		for _, m := range msgs {
			t := m.GmailThread
			byThread[t] = append(byThread[t], m)
		}
		
		var threadList [][]*imap.Msg
		for _, t := range byThread {
			sort.Sort(byUID(t))
			threadList = append(threadList, t)
		}
		sort.Sort(byUIDList(threadList))
		
		msgs = msgs[:0]
		for _, t := range threadList {
			msgs = append(msgs, t...)
		}
	}
	for i, m := range msgs {
		msgNum[m] = i
	}
}

type byUID []*imap.Msg

func (l byUID) Less(i, j int) bool { return l[i].UID < l[j].UID }
func (l byUID) Len() int { return len(l) }
func (l byUID) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

type byUIDRev []*imap.Msg

func (l byUIDRev) Less(i, j int) bool { return l[i].UID > l[j].UID }
func (l byUIDRev) Len() int { return len(l) }
func (l byUIDRev) Swap(i, j int) { l[i], l[j] = l[j], l[i] }


type byUIDList [][]*imap.Msg

func (l byUIDList) Less(i, j int) bool { return l[i][len(l[i])-1].UID > l[j][len(l[j])-1].UID }
func (l byUIDList) Len() int { return len(l) }
func (l byUIDList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func subj(m *imap.Msg) string {
	s := m.Hdr.Subject
	for strings.HasPrefix(s, "Re: ") || strings.HasPrefix(s, "RE: ") {
		s = s[4:]
	}
	return s
}

func Wcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	if dot == nil {
		return nil
	}
	if !isGmail {
		fmt.Fprintf(bout, "!cmd W requires gmail\n")
		return dot
	}
	url := fmt.Sprintf("https://mail.google.com/mail/b/%s/?shva=1#inbox/%x", acct.Email, dot.Msg.GmailThread)
	if err := exec.Command("open", url).Run(); err != nil {
		fmt.Fprintf(bout, "!%s\n", err)
	}
	return dot
}

func xcmd(c *Cmd, dot *imap.MsgPart) *imap.MsgPart {
	// TODO: remove saved attachments?
	os.Exit(0)
	panic("not reached")
}

func flushDeleted() {
	var toDelete []*imap.Msg
	for m := range deleted {
		toDelete = append(toDelete, m)
	}
	if len(toDelete) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "!deleting %d\n", len(toDelete))
	err := inbox.Delete(toDelete)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!deleting: %s\n", err)
	}
}

func loadNew() {
	if err := inbox.Check(); err != nil {
		fmt.Fprintf(os.Stderr, "!inbox: %s\n", err)
	}
	
	old := make(map[*imap.Msg]bool)
	for _, m := range msgs {
		old[m] = true
	}
	
	nnew := 0
	new := inbox.Msgs()
	for _, m := range new {
		if old[m] {
			old[m] = false, false
		} else {
			msgs = append(msgs, m)
			nnew++
		}
	}
	if nnew > 0 {
		fmt.Fprintf(os.Stderr, "!%d new messages\n", nnew)
	}
	for m := range old {
		// Deleted
		m.Flags |= imap.FlagDeleted
		deleted[m] = nil, false
	}
}

func sync(delete bool) {
	if delete {
		flushDeleted()
	}
	loadNew()
	if delete {
		w := 0
		for _, m := range msgs {
			if !m.Deleted() {
				msgs[w] = m
				w++
			}
		}
		msgs = msgs[:w]
	}
	rethread()
}
