// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	
	"goplan9.googlecode.com/hg/plan9/acme"
	"rsc.googlecode.com/hg/google"
	"rsc.googlecode.com/hg/xmpp"
)

var windows sync.WaitGroup
var client *google.Client

func main() {
	w, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	
	client, err = google.Dial()
	if err != nil {
		w.Name("Chat/+Errors")
		w.Write("body", []byte(err.String()))
		return
	}
	
	windows.Add(1)
	go root(w)
	windows.Wait()
}

func root(w *acme.Win) {
	defer windows.Done()

	w.Name("Chat/")
	acct, err := client.Accounts()
	if err != nil {
		w.Write("body", []byte(err.String()))
		return
	}

	for _, a := range acct {
		w.Write("body", []byte(a + "/\n"))
	}

	w.Ctl("clean")
	for word := range events(w) {
		ww, err := acme.New()
		if err != nil {
			w.Write("body", []byte(err.String()+"\n"))
			continue
		}
		windows.Add(1)
		go account(ww, word)
	}
}

func events(w *acme.Win) <-chan string {
	c := make(chan string, 10)
	go func() {
		for e := range w.EventChan() {
			switch e.C2 {
			case 'x', 'X':	// execute
				if string(e.Text) == "Del" {
					w.Ctl("delete")
				}
				w.WriteEvent(e)
			case 'l', 'L':	// look
				w.Ctl("clean")
				c <- string(e.Text)
			}
		}
		w.CloseFiles()
		close(c)
	}()
	return c
}

type Acct struct {
	Name string
	ID *google.ChatID
	Chat map[string]chan *xmpp.Chat
	Main chan *xmpp.Chat
	Group sync.WaitGroup
}

func (a *Acct) process() {
	for {
		msg, err := client.ChatRecv(a.ID)
		if err != nil {
			log.Fatal(err)
			break
		}
		a.Main <- msg
		if i := strings.Index(msg.Remote, "/"); i >= 0 {
			who := msg.Remote[:i]
			if ch := a.Chat[who]; ch != nil {
				ch <- msg
			}
		}
	}
}

func account(w *acme.Win, name string) {
	defer windows.Done()
	
	a := &Acct{Name: name, Chat: map[string]chan *xmpp.Chat{}}
	a.Group.Add(1)
	a.ID = &google.ChatID{ID: name + "/2", Email: name, Status: xmpp.Available, StatusMsg: ""}
	a.Main = make(chan *xmpp.Chat)
	go a.process()

	w.Name("Chat/" + name + "/")
	ev := events(w)
Loop:
	for {
		select {
		case id, ok := <-ev:
			if !ok {
				break Loop
			}
			ww, err := acme.New()
			if err != nil {
				w.Write("body", []byte(err.String()+"\n"))
				continue
			}
			windows.Add(1)
			go a.chat(ww, id)
		case msg := <-a.Main:
			switch msg.Type {
			default:
				w.Write("body", []byte(fmt.Sprintf("? %v\n", msg)))
			case "presence":
				w.Write("body", []byte(fmt.Sprintf("%v\n", msg.Presence)))
			}
		}
	}
}

func (a *Acct) chat(w *acme.Win, name string) {
	defer windows.Done()
	
	w.Name("Chat/" + a.Name + "/" + name)
	ch := make(chan *xmpp.Chat)
	a.Chat[name] = ch
	for msg := range ch {
		switch msg.Type {
		default:
			w.Write("body", []byte(fmt.Sprintf("? %v\n", msg)))
		case "presence":
			w.Write("body", []byte(fmt.Sprintf("%v\n", msg.Presence)))
		case "chat":
			w.Write("body", []byte(fmt.Sprintf("%v %v\n", msg.Remote, msg.Text)))
		}
	}

}
