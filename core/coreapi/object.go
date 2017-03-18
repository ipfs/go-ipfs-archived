package coreapi

import (
	"context"
	"io"

	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"

	cid "gx/ipfs/QmV5gPoRsjN1Gid3LMdNZTyfCtP2DsvqEbMAmz82RmmiGk/go-cid"
)

type ObjectAPI CoreAPI

// set-data /ipfs/Qmfoo/a/b hello
// => /ipfs/Qmbar/a/b
// add-link /ipfs/Qmfoo/x y/z /ipns/example.com
// => /ipfs/Qmbaz/x/y/z

func (api *ObjectAPI) AddLink(ctx context.Context, p coreiface.Path, name string, pl coreiface.Path) (coreiface.Path, error) {
	linkee, err := api.core().ResolveNode(pl)
	if err != nil {
		return nil, err
	}


}

func (api *ObjectAPI) RmLink(ctx context.Context, p coreiface.Path, name string) (coreiface.Path, error) {
}

func (api *ObjectAPI) core() coreiface.CoreAPI {
	return (*CoreAPI)(api)
}


p, err := api.Put(ctx, r.Body)
if err != nil {
  return err
}

pr, err := api.AddLink()
