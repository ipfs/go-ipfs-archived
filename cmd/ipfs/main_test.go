package main

import (
	"testing"

	"github.com/ipfs/go-ipfs-cmds/cmdsutil"
)

func TestIsCientErr(t *testing.T) {
	t.Log("Catch both pointers and values")
	if !isClientError(cmdsutil.Error{Code: cmdsutil.ErrClient}) {
		t.Errorf("misidentified value")
	}
	if !isClientError(&cmdsutil.Error{Code: cmdsutil.ErrClient}) {
		t.Errorf("misidentified pointer")
	}
}
