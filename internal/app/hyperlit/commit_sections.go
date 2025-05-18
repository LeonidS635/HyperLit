package hyperlit

import (
	"context"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/info"
)

func (h *HyperLit) commitSections(ctx context.Context) error {
	commitCtx, commitCtxCancel := context.WithCancel(ctx)
	defer commitCtxCancel()

	errCh := make(chan error)

	var wg sync.WaitGroup

	parseAndCompareModifiedSections := func(status int) {
		states := h.sectionsStates.Get(status)
		h.sectionsStates.Remove(status)

		for _, sectionState := range states {
			wg.Add(1)
			go func() {
				defer wg.Done()

				newSectionsTrie, err := h.parser.HandleParsedSections(commitCtx, sectionState.Path, h.vcs.SaveNewEntry)
				if err != nil {
					helpers.SendCtx(commitCtx, errCh, err)
					return
				}

				updatedSectionsTrie, err := info.CompareSectionsTries(
					ctx, newSectionsTrie, sectionState.CurTrie, h.sectionsStates,
				)
				if err != nil {
					helpers.SendCtx(commitCtx, errCh, err)
					return
				}
				sectionState.ProjectTrieNode.Replace(updatedSectionsTrie)
			}()
		}
	}

	parseAndCompareModifiedSections(info.StatusProbablyModified)
	parseAndCompareModifiedSections(info.StatusCreated)

	return helpers.WaitCtx(ctx, &wg, errCh)
}
