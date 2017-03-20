package cfgds

import (
	"github.com/ipfs/go-ipfs/repo"

	ds "gx/ipfs/QmRWDav6mzWseLWeYfVd5fvUKiVe9xNH29YfMF438fG364/go-datastore"
)

var (
	fErrUnknownDatastoreType = "unknown datastore: %s"
)

// Ctor interface is used to create Datastore from config specification
type Ctor interface {
	Create(repo.Repo) (repo.Datastore, error)
}

type memCtor struct{}

func (_ *memCtor) Create(_ repo.Repo) (repo.Datastore, error) {
	return ds.NewMapDatastore(), nil
}

var _ Ctor = (*memCtor)(nil)
