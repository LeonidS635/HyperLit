package vcs

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/roothash"
	"github.com/LeonidS635/HyperLit/internal/vcs/storage"
)

type VCS struct {
	rootHash roothash.RootHash
	storage  *storage.ObjectsStorage
}

func NewVCS(projectPath, hlPath string) VCS {
	return VCS{
		rootHash: roothash.NewRoot(hlPath),
		storage:  storage.NewObjectsStorage(projectPath, hlPath),
	}
}

func (v *VCS) Init() error {
	return v.storage.Init()
}

func (v *VCS) SaveNewEntry(ctx context.Context, entry entry.Interface) error {
	return v.storage.SaveNewEntry(entry)
}

func (v *VCS) SaveOldEntry(ctx context.Context, hash string) error {
	return v.storage.SaveOldEntry(hash)
}

func (v *VCS) LoadEntry(hash string) (entry.Interface, error) {
	return v.storage.LoadEntry(hash)
}

func (v *VCS) LoadEntryData(hash string) ([]byte, error) {
	e, err := v.storage.LoadEntry(hash)
	if err != nil {
		return nil, err
	}
	return e.GetContent(), nil
}

func (v *VCS) Dump(ctx context.Context) error {
	return v.storage.Dump()
}

func (v *VCS) Read(ctx context.Context, rootHash string) (*trie.Node[info.Section], error) {
	return v.storage.Traverse(ctx, rootHash)
}

func (v *VCS) SaveRootHash(hash []byte) error {
	return v.rootHash.Save(hash)
}

func (v *VCS) GetRootHash() (string, error) {
	return v.rootHash.Get()
}

func (v *VCS) Clear() error {
	return v.storage.Clear()
}
