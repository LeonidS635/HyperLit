package info

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func Compare(
	ctx context.Context, files *trie.Node[File], sections *trie.Node[Section], rootPath string,
) *SectionsStatuses {
	sectionsStatuses := newSectionsStatuses()
	if sections == nil {
		sectionsStatuses.Add(StatusCreated, SectionStatus{Path: rootPath, Trie: nil})
		return sectionsStatuses
	}

	var wg sync.WaitGroup

	done := make(chan struct{})

	wg.Add(1)
	go compare(ctx, files, sections, rootPath, sectionsStatuses, &wg)
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
	ctx context.Context, files *trie.Node[File], sections *trie.Node[Section],
	path string, sectionsStatuses *SectionsStatuses, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	fileInfo := files.Data
	sectionInfo := sections.Data
	if !fileInfo.IsDir {
		status := StatusUnmodified
		if !areEqual(fileInfo, sectionInfo) {
			status = StatusProbablyModified
		}

		sectionsStatuses.Add(status, SectionStatus{Path: path, Trie: sections})
		return
	}

	childrenFiles := files.GetAll()
	childrenSections := sections.GetAll()

	fmt.Println(path, childrenFiles, childrenSections)

	seen := make(map[string]struct{})
	for name, file := range childrenFiles {
		filePath := filepath.Join(path, name)

		if section, ok := childrenSections[name]; ok {
			wg.Add(1)
			go compare(ctx, file, section, filePath, sectionsStatuses, wg)
		} else {
			sectionsStatuses.Add(StatusCreated, SectionStatus{Path: filePath, Trie: nil})
		}
		seen[name] = struct{}{}
	}

	for name := range childrenSections {
		if _, ok := seen[name]; !ok {
			sectionsStatuses.Add(StatusDeleted, SectionStatus{filepath.Join(path, name), nil})
		}
	}
}
