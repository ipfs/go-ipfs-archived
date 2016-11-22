package posinfo

import (
	"os"

	node "gx/ipfs/QmNbpVWj1LwXR74hc3fuxBSmJoXtpgoKcnd1N7J6e88bRU/go-ipld-node"
)

type PosInfo struct {
	Offset   uint64
	FullPath string
	Stat     os.FileInfo // can be nil
}

type FilestoreNode struct {
	node.Node
	PosInfo *PosInfo
}
