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
	rootNode *trie.Node[TrieSection],
	rootPath string, sectionsStatuses *SectionsStatuses,
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
	curNode *trie.Node[TrieSection],
	path string, sectionsStatuses *SectionsStatuses, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	newSectionInfo := newSections.Data
	prevSectionInfo := prevSections.Data
	curNode.Data.Section = newSectionInfo.This

	status := StatusUnmodified
	if !areSectionsEqual(newSectionInfo, prevSectionInfo) {
		status = StatusModified
	}
	curNode.Data.Status = status
	sectionsStatuses.Add(status, SectionStatus{Path: path, Trie: prevSections, FullTrieNode: curNode})

	newChildrenSections := newSections.GetAll()
	prevChildrenSections := prevSections.GetAll()

	seen := make(map[string]struct{})
	for name, newS := range newChildrenSections {
		sectionPath := filepath.Join(path, name)
		nextNode := curNode.Insert(name)

		if prevS, ok := prevChildrenSections[name]; ok {
			wg.Add(1)
			go compareSectionsInOneFile(ctx, newS, prevS, nextNode, sectionPath, sectionsStatuses, wg)
		} else {
			nextNode.Data.Status = StatusCreated
			sectionsStatuses.Add(StatusCreated, SectionStatus{Path: sectionPath, Trie: nil, FullTrieNode: nextNode})
		}
		seen[name] = struct{}{}
	}

	for name, prevS := range prevChildrenSections {
		sectionPath := filepath.Join(path, name)
		if _, ok := seen[name]; !ok {
			nextNode := curNode.Insert(name)
			nextNode.Data.Status = StatusDeleted
			nextNode.Data.Section = prevS.Data.This
			sectionsStatuses.Add(StatusDeleted, SectionStatus{Path: sectionPath, Trie: nil, FullTrieNode: nextNode})
		}
	}
}
