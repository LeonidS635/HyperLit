package vcs

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/vcs/index"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/storage"
)

type VCS struct {
	index   index.Index
	storage storage.ObjectsStorage
}

func NewVCS(path string) VCS {
	return VCS{
		index:   index.NewIndex(path),
		storage: storage.NewObjectsStorage(path),
	}
}

func (v VCS) Init() error {
	if err := v.storage.Init(); err != nil {
		return err
	}
	return nil
}

func (v VCS) LoadEntry(ctx context.Context, hash string) (entry.Entry, error) {
	return v.storage.LoadEntry(hash)
}

func (v VCS) Save(ctx context.Context, entriesCh <-chan entry.Interface) error {
	var entriesInfo []entry.Info
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e, ok := <-entriesCh:
			if !ok {
				return nil
			}

			entriesInfo = append(entriesInfo, entry.GetInfo(e))
			if err := v.storage.SaveEntry(e); err != nil {
				return err
			}
		}
	}
}

func (v VCS) Read(ctx context.Context, rootHash string) (*trie.Node[objects.Section], error) {
	readCtx, readCancel := context.WithCancel(ctx)
	defer readCancel()

	root, errCh := v.storage.Traverse(readCtx, rootHash)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err, ok := <-errCh:
		if ok {
			return nil, err
		}
		return root, nil
	}
}

func (v VCS) Close() error {
	return nil
}
