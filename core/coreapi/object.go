package coreapi

import (
	"context"

	core "github.com/ipfs/go-ipfs/core"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	cid "gx/ipfs/QmfSc2xehWmWLnwwYR91Y8QF4xdASypTFVknutoKQS3GHp/go-cid"
)

type ObjectAPI struct {
	node *core.IpfsNode
}

func NewObjectAPI(n *core.IpfsNode) coreiface.ObjectAPI {
	api := &ObjectAPI{n}
	return api
}

func (api *ObjectAPI) Get(context.Context, string) (coreiface.Object, error) {
	obj := &coreiface.Object{}

	return obj, nil
}

func (api *ObjectAPI) Put(context.Context, coreiface.Object) (cid.Cid, error) {
	h := "Qmfoobar"

	return cid.Decode(h), nil
}
