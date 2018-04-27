// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var fuseLoaders = []string{
	"/Library/Filesystems/osxfuse.fs/Contents/Resources/load_osxfuse", // latest

	// old fallbacks
	"/Library/Filesystems/osxfusefs.fs/Support/load_osxfusefs",
}

var fuseMounters = []string{
	"/Library/Filesystems/osxfuse.fs/Contents/Resources/mount_osxfuse", // latest

	// old fallbacks
	"/Library/Filesystems/osxfusefs.fs/Support/mount_osxfusefs",
}

func mount(dir string) (int, string) {
	// Find a fuse loader and run it.
	// No-op if fuse is already loaded.
	for _, loader := range fuseLoaders {
		if _, err := os.Stat(loader); err == nil {
			out, err := exec.Command(loader).CombinedOutput()
			if err != nil {
				return -1, fmt.Sprintf("exec %s: %v\n%s", loader, err, out)
			}
			goto Loaded
		}
	}
	return 0, "cannot find load_osxfuse"
Loaded:

	// Find a fuse mounter.
	var mounter string
	for _, mounter = range fuseMounters {
		if _, err := os.Stat(mounter); err == nil {
			goto HaveMounter
		}
	}
	return 0, "cannot find mount_osxfuse"
HaveMounter:

	var fd int
	var err error
	for i := 0; ; i++ {
		dev := fmt.Sprintf("/dev/osxfuse%d", i)
		fd, err = syscall.Open(dev, os.O_RDWR, 0)
		if err != nil {
			println("TRY", dev, err.Error())
			if err == syscall.ENOENT || i >= 10000 {
				return -1, "no available fuse devices"
			}
			return -1, err.Error()
		}
		break
	}

	cmd := exec.Command(mounter, fmt.Sprint(fd), dir)
	cmd.Env = append(os.Environ(),
		"MOUNT_OSXFUSE_CALL_BY_LIB=",
		"MOUNT_OSXFUSE_DAEMON_PATH="+mounter,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		syscall.Close(fd)
		return -1, fmt.Sprintf("exec mount_osxfuse: %v", err)
	}
	return fd, ""
}
