package hyperlit

import (
	"context"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/docsgenerator"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser"
	"github.com/LeonidS635/HyperLit/internal/vcs"
)

const hlDirName = "hl"

type HyperLit struct {
	projectName string
	projectPath string
	hlPath      string

	projectTrie    *trie.Node[info.Section]
	sectionsStates *info.SectionsStates

	docsGenerator docsgenerator.Generator
	parser        *parser.Parser
	vcs           vcs.VCS
}

func New(projectPath string) *HyperLit {
	hlPath := filepath.Join(projectPath, hlDirName)
	h := &HyperLit{
		projectName: filepath.Base(projectPath),
		projectPath: projectPath,
		hlPath:      hlPath,

		sectionsStates: info.NewSectionsStates(),

		parser: parser.NewParser(projectPath, hlPath),
		vcs:    vcs.NewVCS(projectPath, hlPath),
	}
	h.docsGenerator = docsgenerator.NewGenerator(hlPath, h.vcs.LoadEntryData)
	return h
}

func (h *HyperLit) Init(ctx context.Context) error {
	if err := os.Mkdir(h.hlPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return h.vcs.Init()
}
