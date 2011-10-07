// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	
	"goplan9.googlecode.com/hg/plan9/acme"
	"rsc.googlecode.com/hg/google"
	"rsc.googlecode.com/hg/xmpp"
)

type Window struct {
	*acme.Win
	*acme.Event
	id *google.ChatID
	typ string
	name string
	remote string
	err os.Error
	blinky bool
	dirty bool
	
	hostpt int
	lastc2 int
}

type Msg struct {
	w *Window
	*xmpp.Chat
	err os.Error
}

var (
	client *google.Client
	account []string
	active = make(map[string]*Window)
	acmeChan = make(chan *Window)
	msgChan = make(chan *Msg)
)

func main() {
	ww, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	ww.Name("Chat/")
	
	client, err = google.Dial()
	if err != nil {
		ww.Printf("body", "%s\n", err)
		return
	}
	
	acct, err := client.Accounts()
	if err != nil {
		ww.Printf("body", "%s\n", err)
		return
	}
	account = acct

	w := &Window{Win: ww, typ: "main", name: "Chat/"}
	active["Chat/"] = w
	for _, a := range acct {
		w.Printf("body", "%s\n", a)
	}
	go w.readAcme()
	mainLoop()
}

func mainLoop() {
	tick := time.Tick(0.5e9)

Loop:
	for len(active) > 0 {
		select {
		case w := <-acmeChan:
			if w == nil {
				// Sync with reader.
				continue
			}
			// Expand clicks, because acme doesn't like . or @ in words.
			if w.err != nil {
				if active[w.name] == nil {
					continue
				}
				log.Fatal(w.err)
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
				case 'L':
					// Button 3 in body: load buddy list for account.
					w.expand()
					arg := string(w.Text)
					for _, a := range account {
						if a == arg {
							showAcct(arg)
							continue Loop
						}
					}
					log.Printf("unknown account %s\n", arg)
					continue Loop
				}					
			case "acct":
				switch w.C2 {
				case 'L':
					// Button 3 in body: load chat window for contact.
					w.expand()
					arg := string(w.Text)
					showContact(w.id.Email, arg)
					continue Loop
				}					
			case "chat":
				lastc2 := w.lastc2
				w.lastc2 = w.C2
				if w.C1 == 'F' && w.C2 == 'I' {
					if lastc2 == 'D' {
						// Something we said.
						w.hostpt = w.Q1
					} else {
						// Something someone else said.
						w.hostpt = w.Q1 + 1
						if w.Q1 == 0 {
							w.hostpt = 0
						}
					}
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
					if w.Q0 < w.hostpt {
						w.hostpt += w.Q1 - w.Q0
					} else {
						w.sendMsg()
					}
					continue Loop
				case 'D':
					if w.Q0 < w.hostpt {
						if w.hostpt < w.Q1 {
							w.hostpt = w.Q0
						} else {
							w.hostpt -= w.Q1 - w.Q0
						}
					}
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
				w := showContact(w.id.Email, you)
				w.fixHostpt()
				text := strings.TrimSpace(msg.Text)
				if text == "" {
					// Probably a composing notification.
					w.blinky = true
					continue
				}
				w.Addr("#%d", w.hostpt-1)
				w.Printf("data", "%s\n", text)
				w.blinky = true
				w.dirty = true

			case "presence":
				pr := msg.Presence
				w1 := lookContact(w.id.Email, you)
				if w1 != nil {
					w1.fixHostpt()
					w1.Addr("#%d", w.hostpt-1)
					w1.Printf("data", "[%s %s]\n", pr.Status, pr.StatusMsg)
				}
				w1 = lookAcct(w.id.Email)
				if w1 != nil {
					w1.Printf("body", "[%#q %#q %#q %#q]\n", you, pr.Remote, pr.Status, pr.StatusMsg)
				}				
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

func (w *Window) fixHostpt() {
	var buf [2]byte
	switch w.hostpt {
	case 0:
		goto Fix
	case 1:
		w.Addr("#%d", w.hostpt-1)
		w.Read("data", buf[:1])
		if buf[0] != '\n' {
			goto Fix
		}
	default:
		w.Addr("#%d", w.hostpt-2)
		w.Read("data", buf[:2])
		if buf[0] != '\n' || buf[1] != '\n' {
			goto Fix
		}
	}
	return

Fix:
	w.Addr("#%d,#%d", w.hostpt, w.hostpt)
	w.Printf("data", "\n")
	w.hostpt++
}

func (w *Window) sendMsg() {
	w.fixHostpt()
	if err := w.Addr("#%d", w.hostpt); err != nil {
		w.Addr("$+#0")
		w.hostpt, _, _ = w.ReadAddr()
		return
	}
	if w.Addr(`.,/(.|\n)*\n/`) != nil {
		return
	}
	q0, q1, _ := w.ReadAddr()
	line, _ := w.ReadAll("xdata")
	trim := string(bytes.TrimSpace(line))
	if len(trim) > 0 {
		err := client.ChatSend(w.id, &xmpp.Chat{Remote: w.remote, Type: "chat", Text: trim})
		w.Addr("#%d,#%d", q0-1, q1)
		errstr := ""
		if err != nil {
			errstr = fmt.Sprintf("%s\n", errstr)
		}
		w.Printf("data", "> %s\n%s\n", strings.Replace(trim, "\n", "\n> ", -1), errstr)
	}
	_, w.hostpt, _ = w.ReadAddr()
	w.Printf("ctl", "clean\n")
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
	client.ChatRoster(w.id)
	for {
		msg, err := client.ChatRecv(w.id)
		if err != nil {
			msgChan <- &Msg{w: w, err: err}
			break
		}
//fmt.Printf("%s\n", *msg)
		msgChan <- &Msg{w: w, Chat: msg}
	}
}

func lookAcct(me string) *Window {
	return active["Chat/" + me + "/"]
}

func lookContact(me, you string) *Window {
	return active["Chat/" + me + "/" + you]
}

func showAcct(me string) *Window {
	w := lookAcct(me)
	if w != nil {
		w.Ctl("show\n")
		return w
	}
	
	ww, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	
	name := "Chat/" + me + "/"
	ww.Name(name)
	w = &Window{Win: ww, typ: "acct", name: name}
	w.id = &google.ChatID{ID: me+"/"+randid(), Email: me, Status: xmpp.Available, StatusMsg: ""}
	active[name] = w
	go w.readChat()
	go w.readAcme()
	return w
}

func showContact(me, you string) *Window {
	w := lookContact(me, you)
	if w != nil {
		w.Ctl("show\n")
		return w
	}
	
	ww, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	
	name := "Chat/" + me + "/" + you
	ww.Name(name)
	w = &Window{Win: ww, id: showAcct(me).id, typ: "chat", name: name, remote: you}
	w.fixHostpt()
	w.Printf("tag", "Ack")
	active[name] = w
	go w.readAcme()
	return w
}

func randid() string {
	return fmt.Sprint(time.Nanoseconds())
}
