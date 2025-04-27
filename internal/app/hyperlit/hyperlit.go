package hyperlit

import (
	"context"
	"fmt"
	"os"

	"github.com/LeonidS635/HyperLit/internal/docsgenerator"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser"
	"github.com/LeonidS635/HyperLit/internal/vcs"
)

type HyperLit struct {
	hlPath      string
	projectPath string

	rootSection      *trie.Node[info.TrieSection]
	sectionsStatuses *info.SectionsStatuses

	docsGenerator docsgenerator.Generator
	parser        *parser.Parser
	vcs           vcs.VCS
}

func New(path, projectPath string) *HyperLit {
	return &HyperLit{
		hlPath:      path,
		projectPath: projectPath,
		rootSection: trie.NewNode[info.TrieSection](),

		docsGenerator: docsgenerator.NewGenerator(path),
		parser:        parser.NewParser(),
		vcs:           vcs.NewVCS(path),
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

func (h *HyperLit) Status(ctx context.Context) {
	err := h.getSectionsStatuses(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//for status, paths := range h.sectionsStatuses.Get(info.StatusProbablyModified) {
	//	fmt.Println(paths, status)
	//}
	fmt.Println("Before saving:")
	h.sectionsStatuses.Print()
	fmt.Println("After saving:")
	fmt.Println(h.commitSections(ctx))
	h.sectionsStatuses.Print()
}

func (h *HyperLit) Diff() {}

func (h *HyperLit) removeUnused() {
	return
}
