// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: Something about redialing.

package google

import (
//	"flag"
	"log"
	"os"
	"rpc"
	"net"
	"time"
	"exec"
	"syscall"
)

func Dir() string {
	dir := os.Getenv("HOME") + "/.goog"
	st, err := os.Stat(dir)
	if err != nil {
		if err := os.Mkdir(dir, 0700); err != nil {
			log.Fatal(err)
		}
		st, err = os.Stat(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	if !st.IsDirectory() {
		log.Fatalf("%s exists but is not a directory", dir)
	}
	if st.Mode&0077 != 0 {
		log.Fatalf("%s exists but allows group or other permissions: %#o", dir, st.Mode&0777)
	}
	return dir
}

func Dial() (*Client, os.Error) {
	socket := Dir()+"/socket"
	c, err := net.Dial("unix", socket)
	if err == nil {
		return &Client{rpc.NewClient(c)}, nil
	}
	log.Print("starting server")
	os.Remove(socket)
	runServer()
	for i := 0; i < 50; i++ {
		c, err = net.Dial("unix", socket)
		if err == nil {
			return &Client{rpc.NewClient(c)}, nil
		}
		time.Sleep(200e6)
		if i == 0 {
			log.Print("waiting for server...")
		}
	}
	return nil, err
}

type Client struct {
	client *rpc.Client
}

type Empty struct {}

func (g *Client) Ping() os.Error {
	return g.client.Call("goog.Ping", &Empty{}, &Empty{})
}

func (g *Client) Accounts() ([]string, os.Error) {
	var out []string
	if err := g.client.Call("goog.Accounts", &Empty{}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func runServer() {
	cmd := exec.Command("googleserver", "serve")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}

