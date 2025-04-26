package dirtraverser

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

type FileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
}

type DirTraverser struct {
	sema *resourceslimiter.Semaphore
	wg   *sync.WaitGroup

	filesTrie *trie.Node[FileInfo]
	errCh     chan error
}

func NewDirTraverser() *DirTraverser {
	return &DirTraverser{
		sema:      resourceslimiter.NewSemaphore(),
		wg:        &sync.WaitGroup{},
		filesTrie: trie.NewNode[FileInfo](),
		errCh:     make(chan error),
	}
}

func (t *DirTraverser) GetOutputs() (*trie.Node[FileInfo], <-chan error) {
	return t.filesTrie, t.errCh
}

func (t *DirTraverser) Traverse(ctx context.Context, path string) {
	t.wg.Add(1)
	t.traverse(ctx, path, t.filesTrie)
	t.wg.Wait()

	close(t.errCh)
}

func (t *DirTraverser) traverse(ctx context.Context, path string, curNode *trie.Node[FileInfo]) {
	defer t.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return
	}

	curNode.Data = FileInfo{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	if info.IsDir() {
		for _, entry := range t.getDirEntries(ctx, path) {
			next := curNode.Insert(entry.Name())
			t.wg.Add(1)
			go t.traverse(ctx, filepath.Join(path, entry.Name()), next)
		}
	}
}

func (t *DirTraverser) getDirEntries(ctx context.Context, dir string) []os.DirEntry {
	ok := t.sema.Acquire(ctx)
	if !ok {
		return nil
	}
	defer t.sema.Release()

	entries, err := os.ReadDir(dir)
	if err != nil {
		helpers.SendCtx(ctx, t.errCh, err)
		return nil
	}
	return entries
}
