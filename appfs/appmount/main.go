// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// appmount mounts an appfs file system.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"sync"
	"runtime"

	"code.google.com/p/rsc/appfs/client"
	"code.google.com/p/rsc/appfs/proto"
	"code.google.com/p/rsc/fuse"
	"code.google.com/p/rsc/keychain"
)

var usageMessage = `usage: appmount [-h host] /mnt

Appmount mounts the appfs file system on the named mount point.

The default host is localhost:8080.
`

var cl client.Client
var fc *fuse.Conn

func init() {
	flag.StringVar(&cl.Host, "h", "localhost:8080", "app serving host")
	flag.StringVar(&cl.User, "u", "", "user name")
	flag.StringVar(&cl.Password, "p", "", "password")
}

func usage() {
	fmt.Fprint(os.Stderr, usageMessage)
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
	}
	mtpt := args[0]

	if cl.Password == "" {
		var err error
		cl.User, cl.Password, err = keychain.UserPasswd(cl.Host, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to obtain user and password: %s\n", err)
			os.Exit(2)
		}
	}

	if _, err := cl.Stat("/"); err != nil {
		log.Fatal(err)
	}

	fc, err := fuse.Mount(mtpt)
	if err != nil {
		log.Fatal(err)
	}

	fuse.Debugf = log.Printf
	fmt.Fprintf(os.Stderr, "serving %s\n", mtpt)
	err = fc.Serve(FS{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(2)
	}
}

type FS struct{}

func (FS) Root() (fuse.Node, fuse.Error) {
	return file("/")
}

type File struct {
	Name     string
	FileInfo *proto.FileInfo
	Data     []byte
}

type statEntry struct {
	fi *proto.FileInfo
	err error
	t time.Time
}

var statCache struct {
	mu sync.Mutex
	m map[string] statEntry
}

func stat(name string) (*proto.FileInfo, error) {
	if runtime.GOOS == "darwin" && strings.Contains(name, "/._") {
		// Mac resource forks
		return nil, fmt.Errorf("file not found")
	}
	statCache.mu.Lock()
	e, ok := statCache.m[name]
	statCache.mu.Unlock()
	if ok && time.Since(e.t) < 2*time.Minute {
println("usestat", name)
		return e.fi, e.err
	}
	fi, err := cl.Stat(name)
	saveStat(name, fi, err)
	return fi, err	
}

func saveStat(name string, fi *proto.FileInfo, err error) {
if fi != nil {
	fmt.Fprintf(os.Stderr, "savestat %s %+v\n", name, *fi)
} else {
	fmt.Fprintf(os.Stderr, "savestat %s %v\n", name, err)
}	
	statCache.mu.Lock()
	if statCache.m == nil {
		statCache.m = make(map[string]statEntry)
	}
	statCache.m[name] = statEntry{fi, err, time.Now()}
	statCache.mu.Unlock()
}

func delStat(name string) {
println("delStat", name)
	statCache.mu.Lock()
	if statCache.m != nil {
		delete(statCache.m, name)
	}
	statCache.mu.Unlock()
}

func file(name string) (fuse.Node, fuse.Error) {
	fi, err := stat(name)
	if err != nil {
		if strings.Contains(err.Error(), "no such entity") {
			return nil, fuse.ENOENT
		}
		log.Printf("stat %s: %v", name, err)
		return nil, fuse.EIO
	}
	return &File{name, fi, nil}, nil
}

func (f *File) Attr() (attr fuse.Attr) {
	fi := f.FileInfo
	attr.Mode = 0666
	if fi.IsDir {
		attr.Mode |= 0111 | os.ModeDir
	}
	attr.Mtime =  fi.ModTime
	attr.Size = uint64(fi.Size)
	return
}

func (f *File) Lookup(name string, intr fuse.Intr) (fuse.Node, fuse.Error) {
	return file(path.Join(f.Name, name))
}

func (f *File) ReadAll(intr fuse.Intr) ([]byte, fuse.Error) {
	data, err := cl.Read(f.Name)
	if err != nil {
		log.Printf("read %s: %v", f.Name, err)
		return nil, fuse.EIO
	}
	return data, nil
}

func (f *File) ReadDir(intr fuse.Intr) ([]fuse.Dirent, fuse.Error) {
	fis, err := cl.ReadDir(f.Name)
	if err != nil {
		log.Printf("read %s: %v", f.Name, err)
		return nil, fuse.EIO
	}
	var dirs []fuse.Dirent
	for _, fi := range fis {
		saveStat(path.Join(f.Name, fi.Name), fi, nil)
		dirs = append(dirs, fuse.Dirent{Name: fi.Name})
	}
	return dirs, nil
}

func (f *File) WriteAll(data []byte, intr fuse.Intr) fuse.Error {
	defer delStat(f.Name)
	if err := cl.Write(f.Name[1:], data); err != nil {
		log.Printf("write %s: %v", f.Name, err)
		return fuse.EIO
	}
	return nil
}

func (f *File) Mkdir(req *fuse.MkdirRequest, intr fuse.Intr) (fuse.Node, fuse.Error) {
	defer delStat(f.Name)
	p := path.Join(f.Name, req.Name)
	if err := cl.Create(p[1:], true); err != nil {
		log.Printf("mkdir %s: %v", p, err)
		return nil, fuse.EIO
	}
	delStat(p)
	return file(p)
}

func (f *File) Create(req *fuse.CreateRequest, resp *fuse.CreateResponse, intr fuse.Intr) (fuse.Node, fuse.Handle, fuse.Error) {
	defer delStat(f.Name)
	p := path.Join(f.Name, req.Name)
	if err := cl.Create(p[1:], false); err != nil {
		log.Printf("create %s: %v", p, err)
		return nil, nil, fuse.EIO
	}
	delStat(p)
	n, err := file(p)
	if err != nil {
		return nil, nil, err
	}
	return n, n, nil
}
