// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On-disk data structures

package arq

import (
	"fmt"
	"time"
)

// plist data structures

type computerInfo struct {
	UserName string `plist:"userName"`
	ComputerName string `plist:"computerName"`
}

type folderInfo struct {
	BucketUUID string
	BucketName string
	ComputerUUID string
	LocalPath string
	LocalMountPoint string
	// don't care about IgnoredRelativePaths or Excludes
}

type reflog struct {
	OldHeadSHA1 string `plist:"oldHeadSHA1"`
	NewHeadSHA1 string `plist:"newHeadSHA1"`
}

// binary data structures

type score [20]byte

type sscore struct {
	Score      score `arq:"HexScore"`
	StretchKey bool  // v4+
}

type tag string

type commit struct {
	Tag tag                 `arq:"CommitV005"`
	Author              string
	Comment             string
	ParentCommits       []sscore
	Tree                sscore
	Location            string
	MergeCommonAncestor sscore
	CreateTime          time.Time
	Failed              []failed // v3+
	BucketXML           []byte   // v5+
}

type tree struct {
	Tag tag           `arq:"TreeV015"`
	CompressXattr bool
	CompressACL   bool
	Xattr         sscore
	XattrSize     uint64
	ACL           sscore
	Uid           int32
	Gid           int32
	Mode          int32
	Mtime         unixTime
	Flags         int64
	FinderFlags   int32
	XFinderFlags  int32
	StDev         int32
	StIno         int32
	StNlink       uint32
	StRdev        int32
	Ctime         unixTime
	StBlocks      int64
	StBlksize     uint32
	AggrSize      uint64
	WTF           [16]byte   // what is this?
	Nodes         []nameNode `arq:"count32"`
}

type nameNode struct {
	Name string
	Node node
}

type node struct {
	IsTree            bool
	CompressData      bool
	CompressXattr     bool
	CompressACL       bool
	Blob              []sscore `arq:"count32"`
	UncompressedSize  uint64
	Thumbnail         sscore
	Preview           sscore
	Xattr             sscore
	XattrSize         uint64
	ACL               sscore
	Uid               int32
	Gid               int32
	Mode              int32
	Mtime             unixTime
	Flags             int64
	FinderFlags       int32
	XFinderFlags      int32
	FinderFileType    string
	FinderFileCreator string
	IsExtHidden       bool
	StDev             int32
	StIno             int32
	StNlink           uint32
	StRdev            int32
	Ctime             unixTime
	CreateTime        unixTime
	StBlocks          int64
	StBlksize         uint32
}

type unixTime struct {
	Sec  int64
	Nsec int64
}

type failed struct {
	Path  string
	Error string
}

type ientry struct {
	File   string
	Offset int64
	Size   int64
	Score  score
}

func (s score) Equal(t score) bool {
	for i := range s {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}

func (s score) String() string {
	return fmt.Sprintf("%x", s[:])
}

func binaryScore(b []byte) score {
	if len(b) < 20 {
		panic("BinaryScore: not enough data")
	}
	var sc score
	copy(sc[:], b)
	return sc
}

func hexScore(b string) score {
	if len(b) < 40 {
		panic("HexScore: not enough data")
	}
	var sc score
	for i := 0; i < 40; i++ {
		ch := b[i]
		if '0' <= ch && ch <= '9' {
			ch -= '0'
		} else if 'a' <= ch && ch <= 'f' {
			ch -= 'a' - 10
		} else {
			panic("HexScore: invalid lower hex digit")
		}
		if i%2 == 0 {
			ch <<= 4
		}
		sc[i/2] |= ch
	}
	return sc
}

func (ss sscore) String() string {
	str := ss.Score.String()
	if ss.StretchKey {
		str += "Y"
	}
	return str
}
