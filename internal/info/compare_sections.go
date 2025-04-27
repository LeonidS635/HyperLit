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
	rootPath string, sectionsStatuses *SectionsStatuses,
) {
	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go compareSectionsInOneFile(ctx, newSections, prevSections, rootPath, sectionsStatuses, &wg)
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
	path string, sectionsStatuses *SectionsStatuses, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	newSectionInfo := newSections.Data
	prevSectionInfo := prevSections.Data

	status := StatusUnmodified
	if !areSectionsEqual(newSectionInfo, prevSectionInfo) {
		status = StatusModified
	}
	sectionsStatuses.Add(status, SectionStatus{Path: path, Trie: prevSections})

	newChildrenSections := newSections.GetAll()
	prevChildrenSections := prevSections.GetAll()

	seen := make(map[string]struct{})
	for name, newS := range newChildrenSections {
		sectionPath := filepath.Join(path, name)

		if prevS, ok := newChildrenSections[name]; ok {
			wg.Add(1)
			go compareSectionsInOneFile(ctx, newS, prevS, sectionPath, sectionsStatuses, wg)
		} else {
			sectionsStatuses.Add(StatusCreated, SectionStatus{Path: sectionPath, Trie: nil})
		}
		seen[name] = struct{}{}
	}

	for name := range prevChildrenSections {
		if _, ok := seen[name]; !ok {
			sectionsStatuses.Add(StatusDeleted, SectionStatus{Path: filepath.Join(path, name), Trie: nil})
		}
	}
}
