#!/bin/sh
#
# Copyright (c) 2015 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test private network feature"

. lib/test-lib.sh

test_init_ipfs

export LIBP2P_FORCE_PNET=1

test_expect_success "daemon won't start with force pnet env but with no key" '
	test_must_fail go-timeout 5 ipfs daemon > stdout 2>&1
'

unset LIBP2P_FORCE_PNET

test_expect_success "daemon output incudes info about the reason" '
	grep "private network was not configured but is enforced by the environment" stdout ||
	test_fsh cat stdout
'

pnet_key() {
	echo '/key/swarm/psk/1.0.0/'
	echo '/bin/'
	random 16
}

pnet_key > $IPFS_PATH/swarm.key

LIBP2P_FORCE_PNET=1 test_launch_ipfs_daemon

check_file_fetch() {
	node=$1
	fhash=$2
	fname=$3

	test_expect_success "can fetch file" '
		ipfsi $node cat $fhash > fetch_out
	'

	test_expect_success "file looks good" '
		test_cmp $fname fetch_out
	'
}

test_expect_success "set up tcp testbed" '
	iptb init -n 5 -p 0 -f --bootstrap=none
'

test_done
