package iface

import (
	"context"
	"errors"
	"io"

	ipld "gx/ipfs/QmUsVJ7AEnGyjX8YWnrwq9vmECVGwBQNAKPpgz5KSg8dcq/go-ipld-node"
	cid "gx/ipfs/QmcEcrBAMrwMyhSjXt4yfyPpzgSuV8HLHavnfmiKCSRqZU/go-cid"
)

type Link ipld.Link

type Reader interface {
	io.ReadSeeker
	io.Closer
}

type CoreAPI interface {
	Unixfs() UnixfsAPI
	Object() ObjectAPI
}

type UnixfsAPI interface {
	Add(context.Context, io.Reader) (*cid.Cid, error)
	Cat(context.Context, *cid.Cid) (Reader, error)
	Ls(context.Context, *cid.Cid) ([]*Link, error)
}

type ObjectAPI interface {
	Get(context.Context, *cid.Cid) (*Object, error)
	Put(context.Context, Object) (*cid.Cid, error)
	AddLink(ctx context.Context, root *cid.Cid, path string, target *cid.Cid) (*cid.Cid, error)
	RmLink(ctx context.Context, root *cid.Cid, path string) (*cid.Cid, error)
	// New() (cid.Cid, Object)
	// Links(string) ([]*Link, error)
	// Data(string) (Reader, error)
	// Stat(string) (ObjectStat, error)
	// SetData(string, Reader) (cid.Cid, error)
	// AppendData(string, Data) (cid.Cid, error)
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
var ErrIsNonDag = errors.New("not a merkledag object")
var ErrOffline = errors.New("can't resolve, ipfs node is offline")
