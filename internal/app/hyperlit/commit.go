package hyperlit

import (
	"context"
	"fmt"
	"os"
)

func (h *HyperLit) Commit(ctx context.Context) {
	if err := h.getSectionsStatuses(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := h.commitSections(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h.saveSections(ctx)

	h.vcs.Dump(ctx)

	h.removeUnused()
}
