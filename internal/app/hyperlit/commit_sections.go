package hyperlit

import (
	"context"
	"fmt"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
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
				if err := h.vcs.SaveRootHash(rootHash); err != nil {
					helpers.SendCtx(saveCtx, errCh, err)
					saveCtxCancel()
				}
			}
			if rootNode != nil {
				info.CompareSectionsInOneFile(
					ctx, rootNode, sectionStatus.Trie, sectionStatus.FullTrieNode, sectionStatus.Path,
					h.sectionsStatuses,
				)
			}
		}()
	}

	createdSections := h.sectionsStatuses.Get(info.StatusCreated)
	h.sectionsStatuses.Remove(info.StatusCreated)

	for _, sectionStatus := range createdSections {
		wg.Add(1)
		go func() {
			defer wg.Done()

			rootHash, rootNode := h.parser.Parse(saveCtx, sectionStatus.Path)
			if rootHash != nil && sectionStatus.Path == h.projectPath {
				if err := h.vcs.SaveRootHash(rootHash); err != nil {
					helpers.SendCtx(saveCtx, errCh, err)
					saveCtxCancel()
				}
			}
			if rootNode != nil {
				sectionStatus.Trie = rootNode
				h.sectionsStatuses.Add(info.StatusCreated, sectionStatus)
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

	for _, section := range h.sectionsStatuses.Get(info.StatusModified) {
		//section.FullTrieNode.Data = info.TrieSection{
		//	Section: section.Trie.Data.This,
		//	Status:  info.StatusModified,
		//}
		fmt.Println(hasher.ConvertToHex(section.FullTrieNode.Data.Section.GetHash()))
	}
	for _, section := range h.sectionsStatuses.Get(info.StatusCreated) {
		section.FullTrieNode.Data = info.TrieSection{
			Section: section.Trie.Data.This,
			Status:  info.StatusCreated,
		}
	}

	return nil
}
