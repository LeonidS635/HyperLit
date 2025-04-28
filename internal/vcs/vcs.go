package vcs

import (
	"context"
	"fmt"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
	"github.com/LeonidS635/HyperLit/internal/vcs/roothash"
	"github.com/LeonidS635/HyperLit/internal/vcs/storage"
)

type VCS struct {
	rootHash roothash.RootHash
	storage  storage.ObjectsStorage
}

func NewVCS(path string) VCS {
	return VCS{
		rootHash: roothash.NewRoot(path),
		storage:  storage.NewObjectsStorage(path),
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

func (v VCS) SaveEntry(ctx context.Context, entry entry.Interface) error {
	return v.storage.SaveEntryTmp(entry)
}

func (v VCS) Delete(ctx context.Context, hash string) error {
	return v.storage.Delete(hash)
}

func (v VCS) Dump(ctx context.Context) error {
	return v.storage.Dump()
}

func (v VCS) GetDocsAndCodeFromTree(hash string) ([]byte, []byte, error) {
	e, err := v.storage.LoadEntry(hash)
	if err != nil {
		return nil, nil, err
	}

	if e.Type != format.TreeType {
		return nil, nil, fmt.Errorf(
			"error forming documentatoin: invalid entry type for %s (expected tree, got %v)", hash, e.Type,
		)
	}
	children, err := tree.Parse(e.Data)
	if err != nil {
		return nil, nil, err
	}

	var docs, code []byte
	for _, child := range children {
		switch child.Type {
		case format.CodeType:
			childEntry, err := v.storage.LoadEntry(hasher.ConvertToHex(child.Hash))
			if err != nil {
				return nil, nil, err
			}
			code = childEntry.Data
		case format.DocsType:
			childEntry, err := v.storage.LoadEntry(hasher.ConvertToHex(child.Hash))
			if err != nil {
				return nil, nil, err
			}
			docs = childEntry.Data
		default:
		}
	}
	return docs, code, nil
}

func (v VCS) Save(ctx context.Context, sectionsCh <-chan entry.Interface) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case section, ok := <-sectionsCh:
			if !ok {
				return nil
			}

			if err := v.storage.SaveEntryTmp(section); err != nil {
				return err
			}
		}
	}
}

func (v VCS) Read(ctx context.Context, rootHash string) (*trie.Node[info.Section], error) {
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

func (v VCS) SaveRootHash(hash []byte) error {
	return v.rootHash.Save(hash)
}

func (v VCS) GetRootHash() (string, error) {
	return v.rootHash.Get()
}
