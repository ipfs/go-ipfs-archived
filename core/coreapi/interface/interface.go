package iface

import (
	"context"
	"errors"
	"io"

	cid "gx/ipfs/QmfSc2xehWmWLnwwYR91Y8QF4xdASypTFVknutoKQS3GHp/go-cid"
)

// type CoreAPI interface {
// 	ID() CoreID
// 	Version() CoreVersion
// }

type Object struct {
	Links *[]Link
	Data  Reader
}

type Link struct {
	Name string
	Size uint64
	Cid  *cid.Cid
}

type Reader interface {
	io.ReadSeeker
	io.Closer
}

type UnixfsAPI interface {
	Add(context.Context, io.Reader) (*cid.Cid, error)
	Cat(context.Context, string) (Reader, error)
	Ls(context.Context, string) ([]*Link, error)
}

type ObjectAPI interface {
	// 	New() (cid.Cid, Object)
	Get(context.Context, string) (Object, error)
	// 	Links(string) ([]*Link, error)
	// 	Data(string) (Reader, error)
	// 	Stat(string) (ObjectStat, error)
	Put(context.Context, Object) (cid.Cid, error)
	// 	SetData(string, Reader) (cid.Cid, error)
	// 	AppendData(string, Data) (cid.Cid, error)
	//  AddLink(string, string, string) (cid.Cid, error)
	//  RmLink(string, string) (cid.Cid, error)
}

// type ObjectStat struct {
// 	Cid            cid.Cid
// 	NumLinks       int
// 	BlockSize      int
// 	LinksSize      int
// 	DataSize       int
// 	CumulativeSize int
// }

var ErrIsDir = errors.New("object is a directory")
var ErrOffline = errors.New("can't resolve, ipfs node is offline")
