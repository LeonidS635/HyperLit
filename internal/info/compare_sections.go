package info

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func CompareSectionsInOneFile(
	ctx context.Context, newSections *trie.Node[Section], prevSections *trie.Node[Section],
	rootNode *trie.Node[TrieSection], rootPath string, sectionsStatuses *SectionsStatuses,
) {
	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go compareSectionsInOneFile(ctx, newSections, prevSections, rootNode, rootPath, sectionsStatuses, &wg)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}
}

func compareSectionsInOneFile(
	ctx context.Context, newSections *trie.Node[Section], prevSections *trie.Node[Section],
	curNode *trie.Node[TrieSection], path string, sectionsStatuses *SectionsStatuses, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	var newSectionInfo, prevSectionInfo Section
	var newSectionChildren, prevSectionChildren map[string]*trie.Node[Section]

	if prevSections == nil {
		curNode.Data.Status = StatusCreated
		sectionsStatuses.Add(StatusCreated, SectionStatus{Path: path, Trie: prevSections, FullTrieNode: curNode})
	} else {
		prevSectionInfo = prevSections.Data
		prevSectionChildren = prevSections.GetAll()
	}

	if newSections == nil {
		curNode.Data.Status = StatusDeleted
		sectionsStatuses.Add(StatusDeleted, SectionStatus{Path: path, Trie: prevSections, FullTrieNode: curNode})
	} else {
		newSectionInfo = newSections.Data
		newSectionChildren = newSections.GetAll()
		curNode.Data.Section = newSectionInfo.This
	}

	if prevSections != nil && newSections != nil {
		status := compareTwoSections(newSectionInfo, prevSectionInfo)
		curNode.Data.Status = status
		sectionsStatuses.Add(status, SectionStatus{Path: path, Trie: prevSections, FullTrieNode: curNode})
	}

	seen := make(map[string]struct{})
	for name, newS := range newSectionChildren {
		sectionPath := filepath.Join(path, name)
		nextNode := curNode.Insert(name)

		wg.Add(1)
		go compareSectionsInOneFile(ctx, newS, prevSectionChildren[name], nextNode, sectionPath, sectionsStatuses, wg)

		seen[name] = struct{}{}
	}

	for name, prevS := range prevSectionChildren {
		if _, ok := seen[name]; !ok {
			sectionPath := filepath.Join(path, name)
			nextNode := curNode.Insert(name)

			wg.Add(1)
			go compareSectionsInOneFile(ctx, nil, prevS, nextNode, sectionPath, sectionsStatuses, wg)
		}
	}
}
