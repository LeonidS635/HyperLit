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
	h := &HyperLit{
		hlPath:      path,
		projectPath: projectPath,
		rootSection: trie.NewNode[info.TrieSection](),

		parser: parser.NewParser(),
		vcs:    vcs.NewVCS(path),
	}
	h.docsGenerator = docsgenerator.NewGenerator(path, h.vcs.GetDocsAndCodeFromTree)
	return h
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

func (h *HyperLit) Diff() {}

func (h *HyperLit) removeUnused() {
	return
}
