package coreapi

import (
	"context"

	core "github.com/ipfs/go-ipfs/core"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	dag "github.com/ipfs/go-ipfs/merkledag"
	path "github.com/ipfs/go-ipfs/path"
)

func resolve(ctx context.Context, n *core.IpfsNode, p string) (*dag.Node, error) {
	pp, err := path.ParsePath(p)
	if err != nil {
		return nil, err
	}

	dagnode, err := core.Resolve(ctx, n, pp)
	if err == core.ErrNoNamesys {
		return nil, coreiface.ErrOffline
	} else if err != nil {
		return nil, err
	}
	return dagnode, nil
}
