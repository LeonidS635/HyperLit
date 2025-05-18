package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (s *ObjectsStorage) Traverse(ctx context.Context, rootHash string) (*trie.Node[info.Section], error) {
	t := NewHashTraverser(s)
	return t.Traverse(ctx, rootHash)
}

// Traverser

type HashTraverser struct {
	*ObjectsStorage

	sema *resourceslimiter.Semaphore
	wg   sync.WaitGroup

	sectionsTrie *trie.Node[info.Section]
	errCh        chan error
}

func NewHashTraverser(storage *ObjectsStorage) *HashTraverser {
	return &HashTraverser{
		ObjectsStorage: storage,
		sema:           resourceslimiter.NewSemaphore(),
		wg:             sync.WaitGroup{},
		sectionsTrie:   trie.NewNode[info.Section](),
		errCh:          make(chan error),
	}
}

func (t *HashTraverser) Traverse(ctx context.Context, rootHash string) (*trie.Node[info.Section], error) {
	sectionsTrieRootNode := trie.NewNode[info.Section]()

	traverseCtx, traverseCancel := context.WithCancel(ctx)
	defer traverseCancel()

	// Start traversing
	t.wg.Add(1)
	go t.traverse(traverseCtx, rootHash, t.projectPath, sectionsTrieRootNode)

	err := helpers.WaitCtx(traverseCtx, &t.wg, t.errCh)
	return sectionsTrieRootNode, err
}

func (t *HashTraverser) traverse(ctx context.Context, hash, path string, curNode *trie.Node[info.Section]) {
	defer t.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	e, err := t.LoadEntry(hash)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}
	if eType := e.GetType(); eType != format.TreeType {
		helpers.SendCtx(ctx, t.errCh, fmt.Errorf("expected tree entry, got %v", eType))
		return
	}

	_, filePath := getDirAndFilePathByHash(t.workingDir, hash)
	sectionInfo, err := os.Stat(filePath)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}

	//fmt.Println("Parsing section", path, ", hash", hash, "tree hash", hasher.ConvertToHex(tr.GetHash()))
	section := info.Section{
		Path:  path,
		Hash:  hash,
		MTime: sectionInfo.ModTime(),
		This:  e.(*tree.Tree),
	}

	childEntries, err := tree.Parse(e.GetContent())
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}

	for _, childEntry := range childEntries {
		switch childEntry.Type {
		case format.TreeType:
			t.wg.Add(1)
			next := curNode.Insert(childEntry.Name)
			go t.traverse(ctx, hasher.ConvertToHex(childEntry.Hash), filepath.Join(path, childEntry.Name), next)
		case format.CodeType:
			section.CodeHash = hasher.ConvertToHex(childEntry.Hash)
		case format.DocsType:
			section.DocsHash = hasher.ConvertToHex(childEntry.Hash)
		default:
			helpers.SendCtx(ctx, t.errCh, fmt.Errorf("unknown entry type: %v", childEntry.Type))
			return
		}
	}

	curNode.Data = section
}

func (t *HashTraverser) getChildren(ctx context.Context, content []byte) []entry.Entry {
	ok := t.sema.Acquire(ctx)
	if !ok {
		return nil
	}
	defer t.sema.Release()

	entries, err := tree.Parse(content)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return nil
	}
	return entries
}
