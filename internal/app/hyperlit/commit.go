package hyperlit

import (
	"context"
)

func (h *HyperLit) CommitFirstStep(ctx context.Context) (bool, error) {
	if err := h.Init(ctx); err != nil {
		return false, err
	}

	if err := h.getSectionsStatuses(ctx); err != nil {
		return false, err
	}

	if err := h.commitSections(ctx); err != nil {
		return false, err
	}

	return h.printSectionsInfo()
}

func (h *HyperLit) CommitSecondStep(ctx context.Context) error {
	if err := h.saveSections(ctx); err != nil {
		return err
	}

	if err := h.docsGenerator.Generate(h.projectTrie, h.projectName); err != nil {
		return err
	}

	return h.vcs.Dump(ctx)
}

func (h *HyperLit) Clear() error {
	return h.vcs.Clear()
}
