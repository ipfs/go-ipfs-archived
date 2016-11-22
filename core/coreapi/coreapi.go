package coreapi

import (
	"context"
	"strings"

	core "github.com/ipfs/go-ipfs/core"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	path "github.com/ipfs/go-ipfs/path"

	ipld "gx/ipfs/QmNbpVWj1LwXR74hc3fuxBSmJoXtpgoKcnd1N7J6e88bRU/go-ipld-node"
	cid "gx/ipfs/QmX4hxL9LDFVpYtNfBEBgVSynRGsooVf4F8nrvJiCZuxqq/go-cid"
)

type CoreAPI struct {
	node *core.IpfsNode
}

func NewCoreAPI(n *core.IpfsNode) coreiface.CoreAPI {
	api := &CoreAPI{n}
	return api
}

func (api *CoreAPI) Unixfs() coreiface.UnixfsAPI {
	return (*UnixfsAPI)(api)
}

// TODO: this func does an unneccessary Cid -> Path conversion,
//       the way to solve this is something like core.ResolveCid().
func resolve(ctx context.Context, n *core.IpfsNode, ref coreiface.Ref) (ipld.Node, error) {
	var p path.Path
	sref, ok := ref.(string)
	if ok {
		pp, err := path.ParsePath(sref)
		if err != nil {
			return nil, err
		}
		p = pp
	} else {
		c, err := cid.Parse(ref)
		if err != nil {
			return nil, err
		}
		p = path.FromCid(c)
	}

	dagnode, err := core.Resolve(ctx, n.Namesys, n.Resolver, p)
	if err == core.ErrNoNamesys {
		return nil, coreiface.ErrOffline
	} else if err != nil {
		return nil, err
	}
	return dagnode, nil
}
