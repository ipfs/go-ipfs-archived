#!/bin/sh
#
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test ipfs repo auto gc"

. lib/test-lib.sh

check_ipfs_storage() {
    ipfs config Datastore.StorageMax
}

test_init_ipfs

test_expect_success "generate 2 600 kB files and 2 MB file using go-random" '
    random 600k 41 >600k1 &&
    random 600k 42 >600k2 &&
    random 2M 43 >2M
'

test_expect_success "set ipfs gc watermark, storage max, and gc timeout" '
    test_config_set Datastore.StorageMax "2MB" &&
    test_config_set --json Datastore.StorageGCWatermark 60 &&
    test_config_set Datastore.GCPeriod "20ms"
'

test_launch_ipfs_daemon --enable-gc

test_gc() {
    test_expect_success "adding data below watermark doesn't trigger auto gc" '
        ipfs add -q --pin=false 600k1 >/dev/null &&
        disk_usage "$IPFS_PATH/blocks" >expected &&
        go-sleep 200ms &&
        disk_usage "$IPFS_PATH/blocks" >actual &&
        test_cmp expected actual
    '

    test_expect_success "adding data beyond watermark triggers auto gc" '
        ipfs add -q --pin=false 600k2 >/dev/null &&
        go-sleep 200ms &&
        DU=$(disk_usage "$IPFS_PATH/blocks") &&
        if test $(uname -s) = "Darwin"; then
            test "$DU" -lt 1400  # 60% of 2MB
        else
            test "$DU" -lt 1000000
        fi
    '
}

test_expect_success "periodic auto gc stress test" '
    for i in $(test_seq 1 20)
    do
        test_gc
    done
'

test_kill_ipfs_daemon

test_done
