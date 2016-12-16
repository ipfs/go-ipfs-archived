package corerepo

import (
	"context"
	"errors"
	"time"

	bstore "github.com/ipfs/go-ipfs/blocks/blockstore"
	dag "github.com/ipfs/go-ipfs/merkledag"
	mfs "github.com/ipfs/go-ipfs/mfs"
	pin "github.com/ipfs/go-ipfs/pin"
	gc "github.com/ipfs/go-ipfs/pin/gc"
	repo "github.com/ipfs/go-ipfs/repo"
	config "github.com/ipfs/go-ipfs/repo/config"

	humanize "gx/ipfs/QmPSBJL4momYnE7DcUyk2DVhD6rH488ZmHBGLbxNdhU44K/go-humanize"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	cid "gx/ipfs/QmcTcsTvfaeEBRFo1TkFgT8sRmgi1n1LTZpecfVP8fzpGD/go-cid"
)

var log = logging.Logger("corerepo")

const (
	defaultGCCapacity  = "10G"
	defaultGCWatermark = 90
	defaultGCPeriod    = "1h"
)

type GC struct {
	Repo        repo.Repo
	Blockstore  bstore.GCBlockstore
	LinkService dag.LinkService
	Pinning     pin.Pinner
	MfsRoot     mfs.Root
	Capacity    uint64
	Watermark   uint64
	Period      time.Period
}

func NewGC(r repo.Repo, bs bstore.GCBlockstore, ls dag.LinkService, pn pin.Pinner, mfsr mfs.Root) *GC {
	return &GC{
		Repo:        r,
		Blockstore:  bs,
		LinkService: ls.GetOfflineLinkService(),
		Pinning:     pn,
		MfsRoot:     mfsr,
		Capacity:    defaultGCCapacity,
		Watermark:   defaultGCWatermark,
		Period:      defaultGCPeriod,
	}
}

func (gc *GC) SetConfig(cfg config.Datastore) error {
	max := cfg.StorageMax
	if cfg.StorageMax == "" {
		max = defaultGCCapacity
	}
	capacity, err := humanize.ParseBytes(max)
	if err != nil {
		return nil, err
	}

	wm := cfg.StorageGCWatermark
	if wm == 0 {
		wm = defaultGCWatermark
	}
	watermark := capacity * wm / 100

	prd := cfg.GCPeriod
	if prd == "" {
		prd = defaultGCPeriod
	}
	period, err := time.ParseDuration(prd)
	if err != nil {
		return err
	}

	gc.Capacity = capacity
	gc.Watermark = watermark
	gc.Period = period
	return nil
}

func (gc *GC) AboveWatermark() (bool, error) {
	stored, err := gc.Repo.GetStorageUsage()
	if err != nil {
		return false, err
	}

	return (stored >= gc.Watermark), nil
}

// Run GC periodically, until the passed Context expires or finishes.
// The period between GC runs is determined by cfg.GCPeriod.
//
// Best used with a goroutine.
func (gc *GC) Run(ctx context.Context, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(gc.Period):
			above, err := gc.AboveWatermark()
			if err != nil {
				errCh <- err
				continue
			}

			if above {
				err = gc.run(ctx, nil)
				if err != nil {
					errCh <- err
				}
			}
		}
	}
}

// Run GC instantly without any regard to GC.Period, and return an
// iterator function covering the CIDs of DAG nodes deleted by the GC run.
// The iterator will return an error if the passed Context expires or finishes
// before the iterator is exhausted.
func (gc *GC) RunOnce(ctx context.Context, force bool) (func(*cid.Cid, bool, error), error) {
	above, err := gc.AboveWatermark()
	if err != nil {
		return nil, err
	}
	if !force && !above {
		return nil, nil
	}

	// TODO: implement

	iter := func(*cid.Cid, bool, error) {
		// TODO: implement
	}
	return iter, nil
}

func (gc *GC) run(ctx context.Context, cidCh chan *cid.Cid) error {
	defer close(cidCh)
	// TODO: implement
	return nil
}
