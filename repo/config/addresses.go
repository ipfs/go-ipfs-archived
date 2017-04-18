package config

// Addresses stores the (string) multiaddr addresses for the node.
type Addresses struct {
	Swarm      []string // addresses the swarm should listen on
	Announce   []string // addresses the swarm should announce to the network
	NoAnnounce []string // addresses the swarm should not announce to the network
	API        string   // address for the local API (RPC)
	Gateway    string   // address to listen on for IPFS HTTP object gateway
}
