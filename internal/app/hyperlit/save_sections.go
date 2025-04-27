package hyperlit

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (h *HyperLit) saveSections(ctx context.Context) {
	if h.saveSection(ctx, h.rootSection, h.projectPath) {
		h.vcs.SaveRootHash(h.rootSection.Data.Section.GetHash())
	}
}

func (h *HyperLit) saveSection(ctx context.Context, curNode *trie.Node[info.TrieSection], name string) bool {
	switch curNode.Data.Status {
	case info.StatusCreated, info.StatusModified:
		h.vcs.SaveEntry(ctx, curNode.Data.Section)
		return true
	case info.StatusUnmodified:
		modified := false
		for childName, child := range curNode.GetAll() {
			modified = modified || h.saveSection(ctx, child, childName)
		}
		if modified {
			curNode.Data.Section, _ = tree.Prepare(name)
			for childName, child := range curNode.GetAll() {
				child.Data.Section.SetName(childName)
				curNode.Data.Section.RegisterEntry(child.Data.Section)
			}
			h.vcs.SaveEntry(ctx, curNode.Data.Section)
		}
		return modified
	}
	return false
}
