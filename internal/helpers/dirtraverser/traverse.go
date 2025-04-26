package dirtraverser

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
)

type DirTraverser struct {
	sema *resourceslimiter.Semaphore
	wg   *sync.WaitGroup

	filesInfoCh chan FileInfo
	errCh       chan error
}

func NewDirTraverser() *DirTraverser {
	return &DirTraverser{
		sema: resourceslimiter.NewSemaphore(),
		wg:   &sync.WaitGroup{},
	}
}

func (t *DirTraverser) GetChannels() (<-chan FileInfo, <-chan error) {
	t.filesInfoCh = make(chan FileInfo, resourceslimiter.MaxOpenedEntries)
	t.errCh = make(chan error)

	return t.filesInfoCh, t.errCh
}

func (t *DirTraverser) Traverse(ctx context.Context, path string) {
	t.wg.Add(1)
	t.traverse(ctx, path)
	t.wg.Wait()

	close(t.filesInfoCh)
}

func (t *DirTraverser) traverse(ctx context.Context, dir string) {
	defer t.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	for _, entry := range t.getDirEntries(ctx, dir) {
		if entry.IsDir() {
			t.wg.Add(1)
			go t.traverse(ctx, filepath.Join(dir, entry.Name()))
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				helpers.SendCtx(ctx, t.errCh, err)
				return
			}

			helpers.SendCtx(
				ctx, t.filesInfoCh, FileInfo{
					Path:    filepath.Join(dir, entry.Name()),
					Mode:    0,
					ModTime: fileInfo.ModTime(),
				},
			)
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
