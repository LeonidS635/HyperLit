package hyperlit

import (
	"context"
	"fmt"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/info"
)

func (h *HyperLit) commitSections(ctx context.Context) error {
	saveCtx, saveCtxCancel := context.WithCancel(ctx)
	defer saveCtxCancel()

	sectionsCh, sectionErrCh := h.parser.InitChannels()
	done, errCh := make(chan struct{}), make(chan error)

	go func() {
		defer close(done)

		if err := h.vcs.Save(saveCtx, sectionsCh); err != nil {
			helpers.SendCtx(saveCtx, errCh, err)
		}
		saveCtxCancel()
	}()

	wg := &sync.WaitGroup{}
	for _, sectionStatus := range h.sectionsStatuses.Get(info.StatusProbablyModified) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			rootHash, rootNode := h.parser.Parse(saveCtx, sectionStatus.Path)
			if rootHash != nil && sectionStatus.Path == h.projectPath {
				fmt.Println(rootHash)

				if err := h.vcs.SaveRootHash(rootHash); err != nil {
					helpers.SendCtx(saveCtx, errCh, err)
					saveCtxCancel()
				}
			}
			if rootNode != nil {
				info.CompareSectionsInOneFile(ctx, rootNode, sectionStatus.Trie, sectionStatus.Path, h.sectionsStatuses)
			}
		}()
	}
	for _, sectionStatus := range h.sectionsStatuses.Get(info.StatusCreated) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			rootHash, _ := h.parser.Parse(saveCtx, sectionStatus.Path)
			if rootHash != nil && sectionStatus.Path == h.projectPath {
				fmt.Println(rootHash)

				if err := h.vcs.SaveRootHash(rootHash); err != nil {
					helpers.SendCtx(saveCtx, errCh, err)
					saveCtxCancel()
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		h.parser.CloseChannels()
	}()

	select {
	case <-ctx.Done():
	case <-done:
	case err, ok := <-errCh:
		if ok {
			return err
		}
	case err, ok := <-sectionErrCh:
		if ok {
			return err
		}
	}

	h.sectionsStatuses.Remove(info.StatusProbablyModified)

	return nil
}
