// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "github.com/rsc/rsc/blog/post"
	"github.com/rsc/rsc/devweb/slave"
)

func main() {
	slave.Main()
}
