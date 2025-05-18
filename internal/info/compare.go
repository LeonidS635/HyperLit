package info

import (
	"context"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func Compare(
	ctx context.Context, filesTrieNode *trie.Node[File], sectionsTrieNode *trie.Node[Section],
	sectionsStates *SectionsStates,
) (*trie.Node[Section], error) {
	sectionsTrieRootNode := sectionsTrieNode
	if sectionsTrieNode == nil {
		sectionsTrieRootNode = trie.NewNode[Section]()
		sectionsTrieRootNode.Data.Path = filesTrieNode.Data.Path
		sectionsStates.Add(StatusCreated, sectionsTrieRootNode.Data.Path, sectionsTrieRootNode, nil)
		return sectionsTrieRootNode, nil
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go compare(ctx, filesTrieNode, sectionsTrieRootNode, sectionsStates, &wg)

	_ = helpers.WaitCtx(ctx, &wg, nil)
	return sectionsTrieRootNode, ctx.Err()
}

func compare(
	ctx context.Context, filesTrieNode *trie.Node[File], sectionsTrieNode *trie.Node[Section],
	sectionsStates *SectionsStates, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	fileInfo := filesTrieNode.Data
	sectionInfo := sectionsTrieNode.Data

	if !fileInfo.IsDir {
		status := compareFileAndSection(fileInfo, sectionInfo)

		if status == StatusUnmodified {
			propagateStatus(ctx, status, sectionsTrieNode, sectionsStates)
		} else {
			sectionsTrieNode.Data.Status = status
			sectionsStates.Add(status, sectionsTrieNode.Data.Path, sectionsTrieNode, sectionsTrieNode)
		}

		return
	}

	childrenFiles := filesTrieNode.GetAll()
	childrenSections := sectionsTrieNode.GetAll()

	seen := make(map[string]struct{})
	for name, childFile := range childrenFiles {
		if childSection, ok := childrenSections[name]; ok {
			wg.Add(1)
			go compare(ctx, childFile, childSection, sectionsStates, wg)
		} else {
			nextNode := sectionsTrieNode.Insert(name)
			nextNode.Data.Path = childFile.Data.Path
			sectionsStates.Add(StatusCreated, nextNode.Data.Path, nextNode, nil)
		}
		seen[name] = struct{}{}
	}

	for name, childSection := range childrenSections {
		if _, ok := seen[name]; !ok {
			propagateStatus(ctx, StatusDeleted, childSection, sectionsStates)
		}
	}
}
