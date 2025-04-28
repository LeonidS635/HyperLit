package hashtraverser

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

type HashTraverser struct {
	sema *resourceslimiter.Semaphore
	wg   *sync.WaitGroup

	loadEntryFn func(hash string) (entry.Entry, error)

	sectionsTrie *trie.Node[info.Section]
	errCh        chan error
}

func NewHashTraverser(loadEntryFn func(hash string) (entry.Entry, error)) *HashTraverser {
	return &HashTraverser{
		sema:         resourceslimiter.NewSemaphore(),
		wg:           &sync.WaitGroup{},
		loadEntryFn:  loadEntryFn,
		sectionsTrie: trie.NewNode[info.Section](),
		errCh:        make(chan error),
	}
}

func (t *HashTraverser) GetOutputs() (*trie.Node[info.Section], <-chan error) {
	return t.sectionsTrie, t.errCh
}

func (t *HashTraverser) Traverse(ctx context.Context, rootHash string) {
	t.wg.Add(1)
	t.traverse(ctx, rootHash, t.sectionsTrie)
	t.wg.Wait()

	close(t.errCh)
}

func (t *HashTraverser) traverse(ctx context.Context, hash string, curNode *trie.Node[info.Section]) {
	defer t.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	e, err := t.loadEntryFn(hash)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}
	if e.Type != format.TreeType {
		helpers.SendCtx(ctx, t.errCh, fmt.Errorf("expected tree entry, got %s", e.Type))
		return
	}

	sectionInfo, err := os.Stat(filepath.Join("hl/objects", hash[:2], hash[2:]))
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}

	tr, err := tree.FromEntry(e)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}

	section := info.Section{
		Hash:  hash,
		MTime: sectionInfo.ModTime(),
		This:  tr,
	}

	childEntries, err := tree.Parse(e.Data)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}

	for _, childEntry := range childEntries {
		switch childEntry.Type {
		case format.TreeType:
			t.wg.Add(1)
			next := curNode.Insert(childEntry.Name)
			go t.traverse(ctx, hasher.ConvertToHex(childEntry.Hash), next)
		case format.CodeType:
			section.CodeHash = hasher.ConvertToHex(childEntry.Hash)
		case format.DocsType:
			section.DocsHash = hasher.ConvertToHex(childEntry.Hash)
		default:
			helpers.SendCtx(ctx, t.errCh, fmt.Errorf("unknown entry type: %v", e.Type))
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
