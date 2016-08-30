package blockstore

import (
	"sync"

	"github.com/ipfs/go-ipfs/blocks"
	key "github.com/ipfs/go-ipfs/blocks/key"

	bloom "gx/ipfs/QmWQ2SJisXwcCLsUXLwYCKSfyExXjFRW2WbBH5sqCUnwX5/bbloom"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

// bloomCached returns Blockstore that caches Has requests using Bloom filter
// Size is size of bloom filter in bytes
func newBloomCachedBS(bs Blockstore, ctx context.Context, bloomSize, hashCount int) (*bloomcache, error) {
	bloom, err := bloom.New(float64(bloomSize), float64(hashCount))
	if err != nil {
		return nil, err
	}
	bc := &bloomcache{blockstore: bs}
	bc.invalidate()
	go bc.build(bloom, ctx)
	bc.initialBuildChan = make(chan struct{})

	return bc, nil
}

type bloomcache struct {
	bloomLock sync.RWMutex
	bloom     *bloom.Bloom

	// This chan is only used for testing to wait for bloom to enable
	initialBuildChan chan struct{}
	blockstore       Blockstore

	// Statistics
	hits   uint64
	misses uint64
}

func (b *bloomcache) invalidate() {
	b.bloomLock.Lock()
	b.bloom = nil
	b.bloomLock.Unlock()
}

func (b *bloomcache) BloomActive() bool {
	return b.getBloom() != nil
}
func (b *bloomcache) getBloom() *bloom.Bloom {
	b.bloomLock.RLock()
	defer b.bloomLock.RUnlock()
	return b.bloom
}

func (b *bloomcache) build(bloom *bloom.Bloom, ctx context.Context) {
	evt := log.EventBegin(ctx, "bloomcache.Rebuild")
	defer evt.Done()

	ch, err := b.blockstore.AllKeysChan(ctx)
	if err != nil {
		log.Errorf("AllKeysChan failed in bloomcache rebuild with: %v", err)
		return
	}
	finish := false
	for !finish {
		select {
		case key, ok := <-ch:
			if ok {
				bloom.AddTS([]byte(key)) // Use binary key, the more compact the better
			} else {
				finish = true
			}
		case <-ctx.Done():
			log.Warning("Cache rebuild closed by context finishing.")
			return
		}
	}
	b.bloomLock.Lock()
	b.bloom = bloom
	b.bloomLock.Unlock()
	close(b.initialBuildChan)
}

func (b *bloomcache) DeleteBlock(k key.Key) error {
	if has, ok := b.hasCached(k); ok && !has {
		return ErrNotFound
	}

	return b.blockstore.DeleteBlock(k)
}

// if ok == false has is inconclusive
// if ok == true then has respons to question: is it contained
func (b *bloomcache) hasCached(k key.Key) (has bool, ok bool) {
	if k == "" {
		// Return cache invalid so call to blockstore
		// in case of invalid key is forwarded deeper
		return false, false
	}

	bloom := b.getBloom()
	// check if bloom filter is active
	if bloom != nil {
		blr := bloom.HasTS([]byte(k))
		if blr == false { // not contained in bloom is only conclusive answer bloom gives
			return false, true
		}
	}
	return false, false
}

func (b *bloomcache) Has(k key.Key) (bool, error) {
	if has, ok := b.hasCached(k); ok {
		return has, nil
	}

	return b.blockstore.Has(k)
}

func (b *bloomcache) Get(k key.Key) (blocks.Block, error) {
	if has, ok := b.hasCached(k); ok && !has {
		return nil, ErrNotFound
	}

	return b.blockstore.Get(k)
}

func (b *bloomcache) Put(bl blocks.Block) error {
	if has, ok := b.hasCached(bl.Key()); ok && has {
		return nil
	}

	err := b.blockstore.Put(bl)
	if err == nil {
		bloom := b.getBloom()
		if bloom != nil {
			b.bloom.AddTS([]byte(bl.Key()))
		}
	}
	return err
}

func (b *bloomcache) PutMany(bs []blocks.Block) error {
	var good []blocks.Block
	for _, block := range bs {
		if has, ok := b.hasCached(block.Key()); !ok || (ok && !has) {
			good = append(good, block)
		}
	}
	err := b.blockstore.PutMany(bs)
	bloom := b.getBloom()
	if err == nil && bloom != nil {
		for _, block := range bs {
			b.bloom.AddTS([]byte(block.Key()))
		}
	}
	return err
}

func (b *bloomcache) AllKeysChan(ctx context.Context) (<-chan key.Key, error) {
	return b.blockstore.AllKeysChan(ctx)
}

func (b *bloomcache) GCLock() Unlocker {
	return b.blockstore.(GCBlockstore).GCLock()
}

func (b *bloomcache) PinLock() Unlocker {
	return b.blockstore.(GCBlockstore).PinLock()
}

func (b *bloomcache) GCRequested() bool {
	return b.blockstore.(GCBlockstore).GCRequested()
}
