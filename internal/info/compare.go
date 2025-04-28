package info

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func Compare(
	ctx context.Context, files *trie.Node[File], sections *trie.Node[Section], rootNode *trie.Node[TrieSection],
	rootPath string,
) *SectionsStatuses {
	sectionsStatuses := newSectionsStatuses()
	if sections == nil {
		rootNode.Data.Status = StatusCreated
		sectionsStatuses.Add(StatusCreated, SectionStatus{Path: rootPath, Trie: nil, FullTrieNode: rootNode})
		return sectionsStatuses
	}

	var wg sync.WaitGroup

	done := make(chan struct{})

	wg.Add(1)
	go compare(ctx, files, sections, rootNode, rootPath, sectionsStatuses, &wg)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
	}

	return sectionsStatuses
}

func compare(
	ctx context.Context, files *trie.Node[File], sections *trie.Node[Section], curNode *trie.Node[TrieSection],
	path string, sectionsStatuses *SectionsStatuses, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	curNode.Data.Section = sections.Data.This

	fileInfo := files.Data
	sectionInfo := sections.Data
	if !fileInfo.IsDir {
		status := compareFileAndSection(fileInfo, sectionInfo)
		curNode.Data.Status = status
		sectionsStatuses.Add(status, SectionStatus{Path: path, Trie: sections, FullTrieNode: curNode})
		return
	}

	childrenFiles := files.GetAll()
	childrenSections := sections.GetAll()

	seen := make(map[string]struct{})
	for name, file := range childrenFiles {
		filePath := filepath.Join(path, name)
		nextNode := curNode.Insert(name)

		if section, ok := childrenSections[name]; ok {
			wg.Add(1)
			go compare(ctx, file, section, nextNode, filePath, sectionsStatuses, wg)
		} else {
			nextNode.Data.Status = StatusCreated
			sectionsStatuses.Add(StatusCreated, SectionStatus{Path: filePath, Trie: nil, FullTrieNode: nextNode})
		}
		seen[name] = struct{}{}
	}

	for name, section := range childrenSections {
		filePath := filepath.Join(path, name)
		if _, ok := seen[name]; !ok {
			nextNode := curNode.Insert(name)
			nextNode.Data.Section = section.Data.This
			nextNode.Data.Status = StatusCreated
			sectionsStatuses.Add(StatusDeleted, SectionStatus{Path: filePath, Trie: nil, FullTrieNode: nextNode})
		}
	}
}
