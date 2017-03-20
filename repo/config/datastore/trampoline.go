package cfgds

import (
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-ipfs/repo"
)

// CoreCtor is entrypoint and trampoline for other datastore Ctors
type CoreCtor struct {
	sub Ctor
}

func (tr *CoreCtor) UnmarshalJSON(data []byte) error {
	typeStruct := &struct {
		Type string
	}{}

	err := json.Unmarshal(data, typeStruct)
	if err != nil {
		return err
	}

	ctrctr, ok := registeredCtors[typeStruct.Type]
	if !ok {
		return fmt.Errorf(fErrUnknownDatastoreType, typeStruct.Type)
	}
	tr.sub = ctrctr()

	return json.Unmarshal(data, tr.sub)
}

func (tr *CoreCtor) Create(r repo.Repo) (repo.Datastore, error) {
	return tr.sub.Create(r)
}

var _ json.Unmarshaler = (*CoreCtor)(nil)
var _ Ctor = (*CoreCtor)(nil)
