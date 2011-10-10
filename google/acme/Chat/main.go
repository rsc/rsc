// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
TODO:
 - Del of main window should move to other window.
 - Editing main window should update status on \n or something like that.
 - Make use of full names from roster
*/

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	
	"goplan9.googlecode.com/hg/plan9/acme"
	"rsc.googlecode.com/hg/google"
	"rsc.googlecode.com/hg/xmpp"
)

var acmeDebug = flag.Bool("acmedebug", false, "print acme debugging")

type Window struct {
	*acme.Win  // acme window
	*acme.Event  // most recent event received
	err os.Error  // error reading event
	typ string  // kind of window "main", "chat"
	name string  // acme window title
	remote string  // for typ=="chat", remote address
	dirty bool  // window is dirty
	blinky bool  // window's dirty box is blinking
	
	lastTime int64
}

type Msg struct {
	w *Window  // window where message belongs
	*xmpp.Chat  // recently received chat
	err os.Error  // error reading chat message
}

var (
	client *xmpp.Client  // current xmpp client (can reconnect)
	acct google.Account  // google acct info
	statusCache = make(map[string][]*xmpp.Presence)
	active = make(map[string]*Window)  // active windows
	acmeChan = make(chan *Window)  // acme events
	msgChan = make(chan *Msg)  // chat events
	mainWin *Window
	status = xmpp.Available
	statusMsg = ""
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: Chat [-a acct] name...\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var acctName = flag.String("a", "", "account to use")

func main() {
	flag.Usage = usage
	flag.Parse()

	acct = google.Acct(*acctName)

	aw, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	aw.Name("Chat/" + acct.Nick + "/")

	client, err = xmpp.NewClient("talk.google.com:443", acct.Email, acct.Password)
	if err != nil {
		log.Fatal(err)
	}
	
	w := &Window{Win: aw, typ: "main", name: "Chat/" + acct.Nick + "/"}
	data, err := ioutil.ReadFile(google.Dir() + "/chat." + acct.Nick)
	if err != nil {
log.Fatal(err)
	}
	if err == nil {
		w.Write("body", data)
	}
	mainWin = w
	active[w.name] = w
	go w.readAcme()
	client.Roster()
	client.Status(status, statusMsg)
	go w.readChat()

	tick := time.Tick(0.5e9)
Loop:
	for len(active) > 0 {
		select {
		case w := <-acmeChan:
			if w == nil {
				// Sync with reader.
				continue
			}
			if w.err != nil {
				if active[w.name] == nil {
					continue
				}
				log.Fatal(w.err)
			}
			if *acmeDebug {
				fmt.Fprintf(os.Stderr, "%s %c%c %d,%d %q\n", w.name, w.C1, w.C2, w.Q0, w.Q1, w.Text)
			}
			if (w.C2 == 'x' || w.C2 == 'X') && string(w.Text) == "Del" {
				// TODO: Hangup connection for w.typ == "acct"?
				active[w.name] = nil, false
				w.Del(true)
				continue Loop
			}

			switch w.typ {
			case "main":
				switch w.C2 {
				case 'L':  // Button 3 in body: load chat window for contact.
					w.expand()
					fallthrough
				case 'l':  // Button 3 in tag
					arg := string(w.Text)
					showContact(arg)
					continue Loop
				}					
			case "chat":
				if w.C1 == 'F' && w.C2 == 'I' {
					continue Loop
				}
				if w.C1 != 'M' && w.C1 != 'K' {
					break
				}
				if w.blinky {
					w.blinky = false
					w.Printf("ctl", "dirty\n")
				}
				switch w.C2 {
				case 'X', 'x':
					if string(w.Text) == "Ack" {
						w.Printf("ctl", "clean\n")
					}
				case 'I':
					w.sendMsg()
					continue Loop
				}
			}
			w.WriteEvent(w.Event)
		
		case msg := <-msgChan:
			w := msg.w
			if msg.err != nil {
				w.Printf("body", "ERROR: %s\n", msg.err)
				continue Loop
			}
			you := msg.Remote
			if i := strings.Index(you, "/"); i >= 0 {
				you = you[:i]
			}
			switch msg.Type {
			case "chat":
				w := showContact(you)
				text := strings.TrimSpace(msg.Text)
				if text == "" {
					// Probably a composing notification.
					continue
				}
				w.message("> %s\n", text)
				w.blinky = true
				w.dirty = true

			case "presence":
				pr := msg.Presence
				pr, new := savePresence(pr, you)
				if !new {
					continue
				}
				w := lookContact(you)
				if w != nil {
					w.status(pr)
				}
				mainStatus(pr, you)
			}

		case t := <-tick:
			_ = t
			for _, w := range active {
				if w.blinky {
					w.dirty = !w.dirty
					if w.dirty {
						w.Printf("ctl", "dirty\n")
					} else {
						w.Printf("ctl", "clean\n")
					}
				}
			}
		}
	}
}

func savePresence(pr *xmpp.Presence, you string) (pr1 *xmpp.Presence, new bool) {
	old := cachedPresence(you)

	pr.StatusMsg = strings.TrimSpace(pr.StatusMsg)
	c := statusCache[you]
	for i, p := range c {
		if p.Remote == pr.Remote {
			c[i] = pr
			c[0], c[i] = c[i], c[0]
			goto Best
		}
	}
	c = append(c, pr)
	c[0], c[len(c)-1] = c[len(c)-1], c[0]
	statusCache[you] = c

Best:
	best := cachedPresence(you)
	return best, old==nil || old.Status!=best.Status || old.StatusMsg != best.StatusMsg
}

func cachedPresence(you string) *xmpp.Presence {
	c := statusCache[you]
	if len(c) == 0 {
		return nil
	}
	best := c[0]
	for _, p := range c {
		if p.Status > best.Status {
			best = p
		}
	}
	return best
}

func short(st xmpp.Status) string {
	switch st {
	case xmpp.Unavailable:
		return "?"
	case xmpp.ExtendedAway:
		return "x"
	case xmpp.Away:
		return "-"
	case xmpp.Available:
		return "+"
	case xmpp.DoNotDisturb:
		return "!"
	}
	return st.String()
}

func long(st xmpp.Status) string {
	switch st {
	case xmpp.Unavailable:
		return "unavailable"
	case xmpp.ExtendedAway:
		return "offline"
	case xmpp.Away:
		return "away"
	case xmpp.Available:
		return "available"
	case xmpp.DoNotDisturb:
		return "busy"
	}
	return st.String()
}

func (w *Window) time() string {
/*
Auto-date chat windows:
 
 	Show date and time on first message.
 	Show time if minute is different from last message.
 	Show date if day is different from last message.
 	
 	Oct 10 12:01 > hi
 	12:03 hello there
 	12:05 > what's up?

	12:10 [Away]
*/
	now := time.Seconds()
	t2 := time.SecondsToLocalTime(now)
	t1 := time.SecondsToLocalTime(w.lastTime)
	w.lastTime = now
	if t1.Month != t2.Month || t1.Day != t2.Day || t1.Year != t2.Year {
		return t2.Format("Jan 2 15:04 ")
	}
	return t2.Format("15:04 ")
}


func (w *Window) status(pr *xmpp.Presence) {
	msg := ""
	if pr.StatusMsg != "" {
		msg = ": " + pr.StatusMsg
	}
	w.message("[%s%s]\n", long(pr.Status), msg)

	data, err := w.ReadAll("tag")
	if err != nil {
		log.Printf("read tag: %v", err)
		return
	}
//log.Printf("tag1: %s\n", data)
	i := bytes.IndexByte(data, '|')
	if i >= 0 {
		data = data[i+1:]
	} else {
		data = nil
	}
//log.Printf("tag2: %s\n", data)
	j := bytes.IndexByte(data, '|')
	if j >= 0 {
		data = data[j+1:]
	}
//log.Printf("tag3: %s\n", data)

	msg = ""
	if pr.StatusMsg != "" {
		msg = " " + pr.StatusMsg
	}
	w.Ctl("cleartag\n")
	w.Write("tag", []byte(" " + short(pr.Status) + msg + " |" + string(data)))
}

func mainStatus(pr *xmpp.Presence, you string) {
	w := mainWin
	if err := w.Addr("#0/^(.[ \t]+)?"+regexp.QuoteMeta(you)+"([ \t]*|$)/"); err != nil {
		return
	}
	q0, q1, err := w.ReadAddr()
	if err != nil {
		log.Printf("ReadAddr: %s\n", err)
		return
	}
	if err := w.Addr("#%d/"+regexp.QuoteMeta(you)+"/", q0); err != nil {
		log.Printf("Addr2: %s\n", err)
	}
	q2, q3, err := w.ReadAddr()
	if err != nil {
		log.Printf("ReadAddr2: %s\n", err)
		return
	}

	space := " "
	if q1 > q3 || pr.StatusMsg == "" { // already have or don't need space
		space = ""
	}
	if err := w.Addr("#%d/.*/", q1); err != nil {
		log.Printf("Addr3: %s\n", err)
	}
	w.Printf("data", "%s%s", space, pr.StatusMsg)

	space = ""
	if q0 == q2 {
		w.Addr("#%d,#%d", q0, q0)
		space = " "
	} else {
		w.Addr("#%d,#%d", q0, q0+1)
	}
	w.Printf("data", "%s%s", short(pr.Status), space)
}

func (w *Window) expand() {
	// Use selection if any.
	w.Printf("ctl", "addr=dot\n")
	q0, q1, err := w.ReadAddr()
	if err == nil && q0 <= w.Q0 && w.Q0 <= q1 {
		goto Read
	}
	if err = w.Addr("#%d-/[a-zA-Z0-9_@.\\-]*/,#%d+/[a-zA-Z0-9_@.\\-]*/", w.Q0, w.Q1); err != nil {
		log.Printf("expand: %v", err)
		return
	}
	q0, q1, err = w.ReadAddr()
	if err != nil {
		log.Printf("expand: %v", err)
		return
	}

Read:
	data, err := w.ReadAll("xdata")
	if err != nil {
		log.Printf("read: %v", err)
		return
	}
	w.Text = data
	w.Q0 = q0
	w.Q1 = q1
	return
}

// Invariant: in chat windows, the acme addr corresponds to the
// empty string just before the input being typed.  Text before addr
// is the chat history (usually ending in a blank line).

func (w *Window) message(format string, args ...interface{}) {
	if *acmeDebug {
		q0, q1, _ := w.ReadAddr()
		log.Printf("message; addr=%d,%d", q0, q1)
	}
	if err := w.Addr(".-/\\n?\\n?/"); err != nil && *acmeDebug {
		log.Printf("set addr: %s", err)
	}
	q0, _, _ := w.ReadAddr()
	nl := ""
	if q0 > 0 {
		nl = "\n"
	}
	if *acmeDebug {
		q0, q1, _ := w.ReadAddr()
		log.Printf("inserting; addr=%d,%d", q0, q1)
	}
	w.Printf("data", nl + w.time() + format + "\n", args...)
	if *acmeDebug {
		q0, q1, _ := w.ReadAddr()
		log.Printf("wrote; addr=%d,%d", q0, q1)
	}
}

func (w *Window) sendMsg() {
	if *acmeDebug {
		q0, q1, _ := w.ReadAddr()
		log.Printf("sendMsg; addr=%d,%d", q0, q1)
	}
	if err := w.Addr(`.,./(.|\n)*\n/`); err != nil {
		if *acmeDebug {
			q0, q1, _ := w.ReadAddr()
			log.Printf("no text (%s); addr=%d,%d", err, q0, q1)
		}
		return
	}
	q0, q1, _ := w.ReadAddr()
	if *acmeDebug {
		log.Printf("found msg; addr=%d,%d", q0, q1)
	}
	line, _ := w.ReadAll("xdata")
	trim := string(bytes.TrimSpace(line))
	if len(trim) > 0 {
		err := client.Send(xmpp.Chat{Remote: w.remote, Type: "chat", Text: trim})
		
		// Select blank line before input (if any) and input.
		w.Addr("#%d-/\\n?\\n?/,#%d", q0, q1)
		if *acmeDebug {
			q0, q1, _ := w.ReadAddr()
			log.Printf("selected text; addr=%d,%d", q0, q1)
		}
		q0, _, _ := w.ReadAddr()
		
		// Overwrite with \nmsg\n\n.
		// Leaves addr after final \n, which is where we want it.
		nl := ""
		if q0 > 0 {
			nl = "\n"
		}
		errstr := ""
		if err != nil {
			errstr = fmt.Sprintf("\n%s", errstr)
		}
		w.Printf("data", "%s%s%s%s\n\n", nl, w.time(), trim, errstr)
		if *acmeDebug {
			q0, q1, _ := w.ReadAddr()
			log.Printf("wrote; addr=%d,%d", q0, q1)
		}
		w.Printf("ctl", "clean\n")
	}
}

func (w *Window) readAcme() {
	for {
		e, err := w.ReadEvent()
		if err != nil {
			w.err = err
			acmeChan <- w
			break
		}
//fmt.Printf("%c%c %d,%d %d,%d %#x %#q %#q %#q\n", e.C1, e.C2, e.Q0, e.Q1, e.OrigQ0, e.OrigQ1, e.Flag, e.Text, e.Arg, e.Loc)
		w.Event = e
		acmeChan <- w
		acmeChan <- nil
	}
}

func (w *Window) readChat() {
	for {
		msg, err := client.Recv()
		if err != nil {
			msgChan <- &Msg{w: w, err: err}
			break
		}
//fmt.Printf("%s\n", *msg)
		msgChan <- &Msg{w: w, Chat: &msg}
	}
}

func lookContact(you string) *Window {
	return active["Chat/" + acct.Nick + "/" + you]
}

func showContact(you string) *Window {
	w := lookContact(you)
	if w != nil {
		w.Ctl("show\n")
		return w
	}
	
	ww, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	
	name := "Chat/" + acct.Nick + "/" + you
	ww.Name(name)
	w = &Window{Win: ww, typ: "chat", name: name, remote: you}
	w.Printf("body", "\n")
	w.Addr("#1")
	w.OpenEvent()
	w.Printf("ctl", "cleartag\n")
	w.Printf("tag", " Ack")
	if p := cachedPresence(you); p != nil {
		w.status(p)
	}
	active[name] = w
	go w.readAcme()
	return w
}

func randid() string {
	return fmt.Sprint(time.Nanoseconds())
}
