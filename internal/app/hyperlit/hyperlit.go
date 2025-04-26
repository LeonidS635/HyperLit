package hyperlit

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser"
	"github.com/LeonidS635/HyperLit/internal/parser/dirtraverser"
	"github.com/LeonidS635/HyperLit/internal/vcs"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

type HyperLit struct {
	hlPath      string
	projectPath string

	parser *parser.Parser
	vcs    vcs.VCS
}

func New(path, projectPath string) *HyperLit {
	return &HyperLit{
		hlPath:      path,
		projectPath: projectPath,
		parser:      parser.NewParser(),
		vcs:         vcs.NewVCS(path),
	}
}

func (h *HyperLit) Init(ctx context.Context) {
	if err := os.Mkdir(h.hlPath, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := h.vcs.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (h *HyperLit) Status(ctx context.Context) []info.FileStatus {
	//root, err := h.vcs.Read(ctx, "bde668440151c9daf9f43555e1cfa49558c3b78bdeacb37ad76162cf726a3181")
	//root, err := h.parser.Traverse(ctx, "testdata")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for path, node := range root.GetAll() {
	//	fmt.Println(path, node.GetAll(), node.Data)
	//	for path, node := range node.GetAll() {
	//		fmt.Println(path, node.GetAll(), node.Data)
	//	}
	//}
	//return nil

	statusCtx, statusCtxCancel := context.WithCancel(ctx)
	defer statusCtxCancel()

	filesRoot, sectionsRoot := trie.NewNode[info.File](), trie.NewNode[info.Section]()
	var filesErr, sectionsErr error

	wg := sync.WaitGroup{}
	done := make(chan struct{})

	wg.Add(2)
	go func() {
		filesRoot, filesErr = h.parser.Traverse(statusCtx, "testdata")
		if filesErr != nil {
			fmt.Println(filesErr)
			statusCtxCancel()
		}
	}()
	go func() {
		sectionsRoot, sectionsErr = h.vcs.Read(
			statusCtx, "bde668440151c9daf9f43555e1cfa49558c3b78bdeacb37ad76162cf726a3181",
		)
		if sectionsErr != nil {
			fmt.Println(sectionsErr)
			statusCtxCancel()
		}
	}()
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
	}

	if filesErr != nil || sectionsErr != nil {
		return nil
	}

	return info.Compare(ctx, filesRoot, sectionsRoot)
}

func (h *HyperLit) compareFileAndEntry(ctx context.Context, fileInfo dirtraverser.FileInfo, entryInfo entry.Info) (
	bool, error,
) {
	parseCtx, parseCancel := context.WithCancel(ctx)
	defer parseCancel()

	sections := make(map[string]parser.Section)
	go func() {
		sectionsCh := h.parser.Sections()

		for {
			select {
			case <-parseCtx.Done():
				return
			case section := <-sectionsCh:
				sections[section.GetPath()] = section // Section Names can be equal in different files but not in one file
			}
		}
	}()

	err := h.parser.ParseFile(ctx, fileInfo.Path)
	if err != nil {
		parseCancel()
		return false, err
	}

	//fileEntry := h.vcs.LoadEntry(ctx, entryInfo.Hash)
	//if fileEntry.Type != format.TreeType {
	//	return false, fmt.Errorf("expected tree entry, found %v", fileEntry.Type)
	//}
	//
	//subEntries, err := tree.Parse(fileEntry.Data)
	//if err != nil {
	//	return false, err
	//}

	return false, nil
}

func (h *HyperLit) Diff() {}

func (h *HyperLit) removeUnused() {
	return
}
