package hyperlit

import (
	"context"
	"fmt"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
)

func (h *HyperLit) getSectionsStatuses(ctx context.Context) error {
	rootHash, err := h.vcs.GetRootHash()
	if err != nil {
		return err
	}

	var filesRoot *trie.Node[info.File]
	var sectionsRoot *trie.Node[info.Section]
	var filesErr, sectionsErr error

	wg := sync.WaitGroup{}
	done := make(chan struct{})

	statusCtx, statusCtxCancel := context.WithCancel(ctx)
	defer statusCtxCancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		filesRoot, filesErr = h.parser.Traverse(statusCtx, h.projectPath)
		if filesErr != nil {
			fmt.Println(filesErr)
			statusCtxCancel()
		}
	}()

	if rootHash != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sectionsRoot, sectionsErr = h.vcs.Read(statusCtx, rootHash)
			if sectionsErr != nil {
				statusCtxCancel()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
	}

	fmt.Println(filesRoot)
	h.sectionsStatuses = info.Compare(ctx, filesRoot, sectionsRoot, h.projectPath)
	return nil
}
