package fsrepo

import (
	"errors"
	"path/filepath"

	"github.com/ipfs/go-ipfs/repo"
	dcfg "github.com/ipfs/go-ipfs/repo/config/datastore"

	flatfs "gx/ipfs/QmXZEfbEv9sXG9JnLoMNhREDMDgkq5Jd7uWJ7d77VJ4pxn/go-ds-flatfs"
	leveldb "gx/ipfs/QmaHHmfEozrrotyhyN44omJouyuEtx6ahddqV6W5yRaUSQ/go-ds-leveldb"
	ldbopts "gx/ipfs/QmbBhyDKsY4mbY6xsKt3qu9Y7FPvMJ6qbD8AMjYYvPRw1g/goleveldb/leveldb/opt"
)

var (
	errNotFSRepo = errors.New("used repo is not FSRepo")
)

func fsRepoOrError(r repo.Repo) (*FSRepo, error) {
	fsrepo, ok := r.(*FSRepo)
	if !ok {
		return nil, errNotFSRepo
	}
	return fsrepo, nil
}

func relPath(fsrepo *FSRepo, p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(fsrepo.path, p)
	}
	return p
}

func registerDStoreCtors() {
	dcfg.RegisterCtor("leveldb", func() dcfg.Ctor { return &levelDBCtor{} })
	dcfg.RegisterCtor("flatfs", func() dcfg.Ctor { return &flatfsCtor{} })
}

type flatfsCtor struct {
	Path      string
	ShardFunc string
	NoSync    bool
}

func (ctr *flatfsCtor) Create(r repo.Repo) (repo.Datastore, error) {
	fsrepo, err := fsRepoOrError(r)
	if err != nil {
		return nil, err
	}

	p := relPath(fsrepo, ctr.Path)
	shardFunc, err := flatfs.ParseShardFunc(ctr.ShardFunc)
	if err != nil {
		return nil, err
	}

	return flatfs.CreateOrOpen(p, shardFunc, ctr.NoSync)
}

var _ dcfg.Ctor = (*flatfsCtor)(nil)

type levelDBCtor struct {
	Path        string
	Compression string
}

func (ctr *levelDBCtor) Create(r repo.Repo) (repo.Datastore, error) {
	fsrepo, err := fsRepoOrError(r)
	if err != nil {
		return nil, err
	}

	p := relPath(fsrepo, ctr.Path)
	c := ldbopts.DefaultCompression
	switch ctr.Compression {
	case "none":
		c = ldbopts.NoCompression
	case "snappy":
		c = ldbopts.SnappyCompression
	}

	return leveldb.NewDatastore(p, &leveldb.Options{
		Compression: c,
	})
}

var _ dcfg.Ctor = (*levelDBCtor)(nil)
