package testutil

import (
	ds2 "github.com/ipfs/go-ipfs/thirdparty/datastore2"
	"gx/ipfs/QmVnFgtzxgPB24nyqLzpWRwjtCDjuda1DipFZcek7TFmmQ/go-datastore"
	syncds "gx/ipfs/QmVnFgtzxgPB24nyqLzpWRwjtCDjuda1DipFZcek7TFmmQ/go-datastore/sync"
)

func ThreadSafeCloserMapDatastore() ds2.ThreadSafeDatastoreCloser {
	return ds2.CloserWrap(syncds.MutexWrap(datastore.NewMapDatastore()))
}
