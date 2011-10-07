// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
//	"flag"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"json"
	"rpc"
	"net"
	"syscall"
	"bufio"
	"strings"

	"rsc.googlecode.com/hg/google"
	"rsc.googlecode.com/hg/xmpp"
)


type Config struct {
	Account []*Account
}

type Account struct {
	Email string
	Password string
}

func (cfg *Config) AccountByEmail(email string) *Account {
	for _, a := range cfg.Account {
		if a.Email == email {
			return a
		}
	}
	return nil
}

var cfg Config

func readConfig() {
	file := google.Dir()+"/config"
	st, err := os.Stat(file)
	if err != nil {
		return
	}
	if st.Mode&0077 != 0 {
		log.Fatalf("%s exists but allows group or other permissions: %#o", file, st.Mode&0777)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	cfg = Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}
}

func writeConfig() {
	file := google.Dir()+"/config"
	st, err := os.Stat(file)
	if err != nil {
		if err := ioutil.WriteFile(file, nil, 0600); err != nil {
			log.Fatal(err)
		}
		st, err = os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	if st.Mode&0077 != 0 {
		log.Fatalf("%s exists but allows group or other permissions: %#o", file, st.Mode&0777)
	}
	data, err := json.MarshalIndent(&cfg, "", "\t");
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(file, data, 0600); err != nil {
		log.Fatal(err)
	}
	st, err = os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	if st.Mode&0077 != 0 {
		log.Fatalf("%s allows group or other permissions after writing: %#o", file, st.Mode&0777)
	}
}

func main() {
	readConfig()
	switch os.Args[1] {
	case "add":
		cfg.Account = append(cfg.Account, &Account{Email: os.Args[2], Password: os.Args[3]})
		writeConfig()
	case "serve":
		serve()
	case "accounts":
		c, err := google.Dial()
		if err != nil {
			log.Fatal(err)
		}
		out, err := c.Accounts()
		if err != nil {
			log.Fatal(err)
		}
		for _, email := range out {
			fmt.Printf("%s\n", email)
		}
	case "ping":
		c, err := google.Dial()
		if err != nil {
			log.Fatal(err)
		}
		if err := c.Ping(); err != nil {
			log.Fatal(err)
		}
	case "chat":
		c, err := google.Dial()
		if err != nil {
			log.Fatal(err)
		}
		cid := &google.ChatID{ID: "1", Email: os.Args[2], Status: xmpp.Available, StatusMsg: ""}
		go chatRecv(c, cid)
		c.ChatRoster(cid)
		b := bufio.NewReader(os.Stdin)
		for {
			line, err := b.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			line = line[:len(line)-1]
			i := strings.Index(line, ": ")
			if i < 0 {
				log.Printf("<who>: <msg>, please")
				continue
			}
			who, msg := line[:i], line[i+2:]
			if err := c.ChatSend(cid, &xmpp.Chat{Remote: who, Type: "chat", Text: msg}); err != nil {
				log.Fatal(err)
			}
		}
	}
}


func chatRecv(c *google.Client, cid *google.ChatID) {
	for {
		msg, err := c.ChatRecv(cid)
		if err != nil {
			log.Fatal(err)
		}
		switch msg.Type {
		case "roster":
			for _, contact := range msg.Roster {
				fmt.Printf("%v\n", contact)
			}
		case "presence":
			fmt.Printf("%v\n", msg.Presence)
		case "chat":
			fmt.Printf("%s: %s\n", msg.Remote, msg.Text)
		default:
			fmt.Printf("<%s>\n", msg.Type)
		}
	}
}

func listen() net.Listener {
	socket := google.Dir()+"/socket"
	os.Remove(socket)
	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func serve() {
	f, err := os.OpenFile(google.Dir() + "/log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	syscall.Dup2(f.Fd(), 2)
	os.Stdout = f
	os.Stderr = f
	l := listen()
	rpc.RegisterName("goog", &Server{})
	rpc.Accept(l)
	log.Fatal("rpc.Accept finished: server exiting")
}

type Server struct {}

type Empty google.Empty

func (*Server) Ping(*Empty, *Empty) os.Error {
	return nil
}

func (*Server) Accounts(_ *Empty, out *[]string) os.Error {
	var email []string
	for _, a := range cfg.Account {
		email = append(email, a.Email)
	}
	*out = email
	return nil
}

