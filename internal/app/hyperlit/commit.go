package hyperlit

import (
	"context"
	"fmt"
	"os"
)

func (h *HyperLit) CommitFirstStep(ctx context.Context) {
	if err := h.getSectionsStatuses(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := h.commitSections(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := h.printSectionsInfo(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (h *HyperLit) CommitSecondStep(ctx context.Context) {
	h.saveSections(ctx)

	if err := h.docsGenerator.Generate(h.rootSection, h.projectPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := h.vcs.Dump(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h.removeUnused()
}
