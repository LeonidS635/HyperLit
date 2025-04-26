package info

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func Compare(ctx context.Context, files *trie.Node[File], sections *trie.Node[Section]) []FileStatus {
	var filesStatus []FileStatus
	var wg sync.WaitGroup
	var mu sync.Mutex
	path := ""

	done := make(chan struct{})

	wg.Add(1)
	go compare(ctx, files, sections, path, &filesStatus, &mu, &wg)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
	}

	return filesStatus
}

func compare(
	ctx context.Context, files *trie.Node[File], sections *trie.Node[Section],
	path string, filesStatus *[]FileStatus, mu *sync.Mutex, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	fileInfo := files.Data
	sectionInfo := sections.Data
	if !fileInfo.IsDir && areEqual(fileInfo, sectionInfo) {
		mu.Lock()
		*filesStatus = append(
			*filesStatus, FileStatus{
				Path:   path,
				Status: StatusUnmodified,
			},
		)
		mu.Unlock()
		return
	}

	childrenFiles := files.GetAll()
	childrenSections := sections.GetAll()

	seen := make(map[string]struct{})
	for name, file := range childrenFiles {
		if section, ok := childrenSections[name]; ok {
			wg.Add(1)
			go compare(ctx, file, section, filepath.Join(path, name), filesStatus, mu, wg)
		} else {
			mu.Lock()
			*filesStatus = append(
				*filesStatus, FileStatus{
					Path:   path,
					Status: StatusCreated,
				},
			)
			mu.Unlock()
		}
		seen[name] = struct{}{}
	}

	for name, section := range childrenSections {
		if _, ok := seen[name]; !ok {
			if file, ok := childrenFiles[name]; ok {
				wg.Add(1)
				go compare(ctx, file, section, filepath.Join(path, name), filesStatus, mu, wg)
			} else {
				mu.Lock()
				*filesStatus = append(
					*filesStatus, FileStatus{
						Path:   path,
						Status: StatusDeleted,
					},
				)
				mu.Unlock()
			}
			seen[name] = struct{}{}
		}
	}
}
