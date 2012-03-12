// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Arqfs implements a file system interface to a collection of Arq backups.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"code.google.com/p/rsc/arq"
	"code.google.com/p/rsc/fuse"
	"code.google.com/p/rsc/keychain"
	"launchpad.net/goamz/aws"
)

func main() {
	access, secret, err := keychain.UserPasswd("s3.amazonaws.com", "")
	if err != nil {
		log.Fatal(err)
	}
	auth := aws.Auth{access, secret}

	conn, err := arq.Dial(auth)
	if err != nil {
		log.Fatal(err)
	}

	comps, err := conn.Computers()
	if err != nil {
		log.Fatal(err)
	}

	fs := &fuse.Tree{}
	for _, c := range comps {
		// TODO: what?
		_, pw, err := keychain.UserPasswd("arq.swtch.com", c.UUID)
		if err != nil {
			log.Fatal(err)
		}
		c.Unlock(pw)

		folders, err := c.Folders()
		if err != nil {
			log.Fatal(err)
		}

		lastDate := ""
		n := 0
		for _, f := range folders {
			if err := f.Load(); err != nil {
				log.Fatal(err)
			}
			trees, err := f.Trees()
			if err != nil {
				log.Fatal(err)
			}
			for _, t := range trees {
				y, m, d := t.Time.Date()
				date := fmt.Sprintf("%04d/%02d%02d", y, m, d)
				suffix := ""
				if date == lastDate {
					n++
					suffix = fmt.Sprintf(".%d", n)
				} else {
					n = 0
				}
				lastDate = date
				f, err := t.Root()
				if err != nil {
					log.Print(err)
				}
				// TODO: Pass times to fs.Add.
				fmt.Printf("%v %s %x\n", t.Time, c.Name+"/"+date+suffix+"/"+t.Path, t.Score)
				fs.Add(c.Name+"/"+date+suffix+"/"+t.Path, &fuseNode{f})
			}
		}
	}

	c, err := fuse.Mount("/mnt/arq")
	if err != nil {
		log.Fatal(err)
	}
	defer exec.Command("umount", "/mnt/arq").Run()

	fmt.Printf("serving /mnt/arq\n")
	c.Serve(fs)
}

type fuseNode struct {
	arq *arq.File
}

func (f *fuseNode) Attr() fuse.Attr {
	de := f.arq.Stat()
	return fuse.Attr{
		Mode:  de.Mode,
		Mtime: de.ModTime,
		Size:  uint64(de.Size),
	}
}

func (f *fuseNode) Lookup(name string, intr fuse.Intr) (fuse.Node, fuse.Error) {
	ff, err := f.arq.Lookup(name)
	if err != nil {
		return nil, fuse.ENOENT
	}
	return &fuseNode{ff}, nil
}

func (f *fuseNode) ReadDir(intr fuse.Intr) ([]fuse.Dirent, fuse.Error) {
	adir, err := f.arq.ReadDir()
	if err != nil {
		return nil, fuse.EIO
	}
	var dir []fuse.Dirent
	for _, ade := range adir {
		dir = append(dir, fuse.Dirent{
			Name: ade.Name,
		})
	}
	return dir, nil
}

// TODO: Implement Read+Release, not ReadAll, to avoid giant buffer.
func (f *fuseNode) ReadAll(intr fuse.Intr) ([]byte, fuse.Error) {
	rc, err := f.arq.Open()
	if err != nil {
		return nil, fuse.EIO
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return data, fuse.EIO
	}
	return data, nil
}
