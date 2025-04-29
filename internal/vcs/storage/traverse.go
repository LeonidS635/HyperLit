package storage

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/storage/hashtraverser"
)

func (s ObjectsStorage) Traverse(ctx context.Context, rootHash string) (*trie.Node[info.Section], <-chan error) {
	t := hashtraverser.NewHashTraverser(s.LoadEntry)
	go t.Traverse(ctx, rootHash)
	return t.GetOutputs()
}
