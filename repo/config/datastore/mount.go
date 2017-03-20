package cfgds

import "github.com/ipfs/go-ipfs/repo"

type mountCtor struct {
	Mounts map[string]*CoreCtor
}

func (ctr mountCtor) Create(r repo.Repo) (repo.Datastore, error) {
	return nil, nil
}

var _ Ctor = (*mountCtor)(nil)
