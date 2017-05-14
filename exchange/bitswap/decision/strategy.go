package decision

import (
	pq "github.com/ipfs/go-ipfs/thirdparty/pq"

	peer "gx/ipfs/QmdS9KpbDyPrieswibZhkod1oXqRwZJrUPzxCofAMWpFGq/go-libp2p-peer"
)

func newStrategyPRQ(strategy Strategy) *prq {
	return &prq{
		taskMap:  make(map[string]*peerRequestTask),
		partners: make(map[peer.ID]*activePartner),
		frozen:   make(map[peer.ID]*activePartner),
		pQueue:   pq.New(getPartnerComparator(strategy)),
	}
}

// a Bitswap Strategy is implemented via the function that orders partners in
// the peerRequestQueue. a Strategy function returns true if peer 'a'
// (represented by activePartner `pa`) has higher priority than peer 'b' (`pb`)
type Strategy func(pa, pb *activePartner) bool

func DefaultStrategy(pa, pb *activePartner) bool {
	// having no blocks in their wantlist means lowest priority
	// having both of these checks ensures stability of the sort
	if pa.requests == 0 {
		return false
	}
	if pb.requests == 0 {
		return true
	}

	if pa.freezeVal > pb.freezeVal {
		return false
	}
	if pa.freezeVal < pb.freezeVal {
		return true
	}

	if pa.active == pb.active {
		// sorting by taskQueue.Len() aids in cleaning out trash entries faster
		// if we sorted instead by requests, one peer could potentially build up
		// a huge number of cancelled entries in the queue resulting in a memory leak
		return pa.taskQueue.Len() > pb.taskQueue.Len()
	}
	return pa.active < pb.active
}

// getPartnerComparator takes in a Strategy function and returns an
// implementation of pq.ElemComparator. This is an Engine function due to the
// required access to peers' ledgers when those peers are compared
func getPartnerComparator(strategy Strategy) pq.ElemComparator {
	return func(a, b pq.Elem) bool {
		pa := a.(*activePartner)
		pb := b.(*activePartner)

		return strategy(pa, pb)
	}
}
