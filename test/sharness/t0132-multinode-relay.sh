#!/bin/sh
#
# Copyright (c) 2015 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test client mode dht"

. lib/test-lib.sh

test_expect_success "set up testbed" '
	iptb init -n 3 -f --bootstrap=none
	ipfsi 0 config --json Addresses.Swarm '"'"'["/ip4/127.0.0.1/tcp/0"]'"'"'
	ipfsi 1 config --json Addresses.Swarm '"'"'["/ip4/127.0.0.1/tcp/0", "/ip6/::1/tcp/0"]'"'"'
	ipfsi 2 config --json Addresses.Swarm '"'"'["/ip6/::1/tcp/0"]'"'"'
	# defuse the fallback listener
	ipfsi 0 config --json Swarm.AddrFilters '"'"'["/ip6/::/ipcidr/0"]'"'"'
	ipfsi 2 config --json Swarm.AddrFilters '"'"'["/ip4/0.0.0.0/ipcidr/0"]'"'"'
'

test_expect_success "start up nodes" '
	iptb start [0-2] --wait
'

test_expect_success "connect up nodes" '
	iptb connect 0 1
	iptb connect 1 2
'

test_expect_success "open the relay and try pinging" '
	peerid0=`ipfsi 0 config Identity.PeerID`
	peerid2=`ipfsi 2 config Identity.PeerID`
	ipfsi 0 swarm connect /exp-relay/$peerid2
	ipfsi 2 ping -n 3 $peerid0
	ipfsi 0 ping -n 3 $peerid2
'

test_expect_success "shut down nodes" '
	iptb stop
'

test_done
