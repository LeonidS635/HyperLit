package hyperlit

import (
	"context"
	"fmt"
	"os"

	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser"
	"github.com/LeonidS635/HyperLit/internal/vcs"
)

type HyperLit struct {
	hlPath      string
	projectPath string

	sectionsStatuses *info.SectionsStatuses

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
