package hyperlit

import (
	"context"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
)

func (h *HyperLit) getSectionsStatuses(ctx context.Context) error {
	// Get saved root hash
	rootHash, err := h.vcs.GetRootHash()
	if err != nil {
		return err
	}

	statusCtx, statusCtxCancel := context.WithCancel(ctx)
	defer statusCtxCancel()

	var filesTrie *trie.Node[info.File]
	var sectionsTrie *trie.Node[info.Section]

	var wg sync.WaitGroup
	errCh := make(chan error)

	// Traverse project files
	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		filesTrie, err = h.parser.Traverse(statusCtx, h.projectPath)
		if err != nil {
			helpers.SendCtx(statusCtx, errCh, err)
		}
	}()

	// Traverse saved sections
	if rootHash != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var err error
			sectionsTrie, err = h.vcs.Read(statusCtx, rootHash)
			if err != nil {
				helpers.SendCtx(statusCtx, errCh, err)
			}
		}()
	}

	if err = helpers.WaitCtx(statusCtx, &wg, errCh); err != nil {
		return err
	}

	// Start building full project sections trie
	h.projectTrie, err = info.Compare(statusCtx, filesTrie, sectionsTrie, h.sectionsStates)
	return err
}
