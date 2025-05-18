package hyperlit

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
)

// TODO: may be situation when error sends in closed errCh (check that IN ALL PROGRAM!!!!!!!!!!!!)
// Perhaps should make errCh buffered (with size sema.recourseLimiter or something like that)

// TODO: Make proper tempdir deletion in case of an error (now renamed files stay empty)

func (h *HyperLit) saveSections(ctx context.Context) error {
	saveCtx, saveCancel := context.WithCancel(ctx)
	defer saveCancel()

	needUpdateRoot := false
	errCh := make(chan error)

	go func() {
		needUpdateRoot = h.saveSection(saveCtx, filepath.Base(h.projectPath), h.projectTrie, errCh)
		close(errCh)
	}()

	if err := helpers.WaitCtx(saveCtx, nil, errCh); err != nil {
		return err
	}

	if needUpdateRoot {
		return h.vcs.SaveRootHash(h.projectTrie.Data.This.GetHash())
	}
	return nil
}

func (h *HyperLit) saveSection(
	ctx context.Context, name string, projectTrieNode *trie.Node[info.Section], errCh chan error,
) bool {
	if helpers.IsCtxCancelled(ctx) {
		return false
	}

	//fmt.Println("Saving", name, "and children", projectTrieNode.GetAll())
	//projectTrieNode.Print()

	var wg sync.WaitGroup
	childrenErrCh := make(chan error)

	var isMySectionOutdated atomic.Bool
	for childName, childSection := range projectTrieNode.GetAll() {
		if childSection.Data.Status == info.StatusDeleted {
			isMySectionOutdated.Store(true)
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				isMySectionOutdated.CompareAndSwap(false, h.saveSection(ctx, childName, childSection, childrenErrCh))
			}()
		}
	}

	if err := helpers.WaitCtx(ctx, &wg, childrenErrCh); err != nil {
		helpers.SendCtx(ctx, errCh, err)
		return false
	}

	// If a section is outdated (some of the child sections have been modified), its contents must be replaced
	// with new child tree hashes. Outdated trees only have trees as children because modified sections with blobs
	// have already been saved during the traversal and comparison of ProbablyModified and Created sections in previous steps.
	if isMySectionOutdated.Load() {
		var err error

		// Clear the contents of outdated section
		projectTrieNode.Data.This, err = projectTrieNode.Data.This.Clear(name)
		if err != nil {
			helpers.SendCtx(ctx, errCh, err)
		}

		// Register all not deleted child sections
		for childName, childSection := range projectTrieNode.GetAll() {
			if childSection.Data.Status != info.StatusDeleted {
				childSection.Data.This.SetName(childName)
				projectTrieNode.Data.This.RegisterEntry(childSection.Data.This)
			}
		}
	}

	switch projectTrieNode.Data.Status {
	case info.StatusDeleted:
		return true
	case info.StatusCreated, info.StatusDocsOutdated, info.StatusCodeOutdated, info.StatusModified:
		isMySectionOutdated.Store(true)
	case info.StatusProbablyModified, info.StatusUnmodified:
		// Save code and docs of unmodified section
		if projectTrieNode.Data.CodeHash != "" {
			if err := h.vcs.SaveOldEntry(ctx, projectTrieNode.Data.CodeHash); err != nil {
				helpers.SendCtx(ctx, errCh, fmt.Errorf("error saving code in section %s: %w", name, err))
				return false
			}
		}
		if projectTrieNode.Data.DocsHash != "" {
			if err := h.vcs.SaveOldEntry(ctx, projectTrieNode.Data.DocsHash); err != nil {
				helpers.SendCtx(ctx, errCh, fmt.Errorf("error saving docs in section %s: %w", name, err))
				return false
			}
		}
	}

	//fmt.Println("saving", name, hasher.ConvertToHex(projectTrieNode.Data.This.GetHash()), isMySectionOutdated.Load())

	// Save current section
	if isMySectionOutdated.Load() {
		if err := h.vcs.SaveNewEntry(ctx, projectTrieNode.Data.This); err != nil {
			helpers.SendCtx(ctx, errCh, err)
			return false
		}
	} else {
		if err := h.vcs.SaveOldEntry(ctx, projectTrieNode.Data.Hash); err != nil {
			helpers.SendCtx(ctx, errCh, err)
			return false
		}
	}

	return isMySectionOutdated.Load()
}
