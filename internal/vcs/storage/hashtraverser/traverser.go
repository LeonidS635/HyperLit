package hashtraverser

import (
	"context"
	"fmt"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

type HashTraverser struct {
	sema *resourceslimiter.Semaphore
	wg   *sync.WaitGroup

	loadEntryFn func(hash string) (entry.Entry, error)

	sectionsTrie *trie.Node[objects.Section]
	errCh        chan error
}

func NewHashTraverser(loadEntryFn func(hash string) (entry.Entry, error)) *HashTraverser {
	return &HashTraverser{
		sema:         resourceslimiter.NewSemaphore(),
		wg:           &sync.WaitGroup{},
		loadEntryFn:  loadEntryFn,
		sectionsTrie: trie.NewNode[objects.Section](),
		errCh:        make(chan error),
	}
}

func (t *HashTraverser) GetOutputs() (*trie.Node[objects.Section], <-chan error) {
	return t.sectionsTrie, t.errCh
}

func (t *HashTraverser) Traverse(ctx context.Context, rootHash string) {
	t.wg.Add(1)
	t.traverse(ctx, rootHash, t.sectionsTrie)
	t.wg.Wait()

	close(t.errCh)
}

func (t *HashTraverser) traverse(ctx context.Context, hash string, curNode *trie.Node[objects.Section]) {
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

	section := objects.Section{Path: e.Path}

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
			go t.traverse(ctx, hasher.ConvertToHex(childEntry.GetHash()), next)
		case format.CodeType:
			section.Code = e
		case format.DocsType:
			section.Docs = e
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
