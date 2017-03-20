package cfgds

import (
	"encoding/json"
	"testing"

	ds "gx/ipfs/QmRWDav6mzWseLWeYfVd5fvUKiVe9xNH29YfMF438fG364/go-datastore"
)

func TestTrampolineToMem(t *testing.T) {
	in := []byte(`{ "type": "mem" }`)
	tr := &CoreCtor{}
	err := json.Unmarshal(in, tr)
	if err != nil {
		t.Fatal(err)
	}

	uds, err := tr.Create()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := uds.(*ds.MapDatastore)
	if !ok {
		t.Fatal("wrong datastore type")
	}

}
