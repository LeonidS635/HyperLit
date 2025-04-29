package hyperlit

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

func (h *HyperLit) saveSections(ctx context.Context) {
	if h.saveSection(ctx, h.rootSection, h.projectPath) {
		h.vcs.SaveRootHash(h.rootSection.Data.Section.GetHash())
	}
}

func (h *HyperLit) saveSection(ctx context.Context, curNode *trie.Node[info.TrieSection], name string) bool {
	isMySectionOutdated := false
	for childName, childSection := range curNode.GetAll() {
		if childSection.Data.Status == info.StatusDeleted {
			isMySectionOutdated = true
		} else {
			isMySectionOutdated = h.saveSection(ctx, childSection, childName) || isMySectionOutdated
		}
	}

	if isMySectionOutdated {
		curNode.Data.Section, _ = curNode.Data.Section.Clear(name)
		for childName, child := range curNode.GetAll() {
			if child.Data.Status != info.StatusDeleted {
				child.Data.Section.SetName(childName)
				curNode.Data.Section.RegisterEntry(child.Data.Section)
			}
		}
	}

	switch curNode.Data.Status {
	case info.StatusDeleted:
		return true
	case info.StatusCreated, info.StatusDocsOutdated, info.StatusCodeOutdated, info.StatusModified:
		isMySectionOutdated = true
	case info.StatusProbablyModified, info.StatusUnmodified:
		// Quick fix, need to rewrite
		entries, _ := curNode.Data.Section.GetEntries()
		for _, e := range entries {
			switch e.Type {
			case format.DocsType, format.CodeType:
				h.vcs.SaveOldEntry(ctx, e)
			default:
			}
		}
	}

	if isMySectionOutdated {
		h.vcs.SaveNewEntry(ctx, curNode.Data.Section)
	} else {
		h.vcs.SaveOldEntry(ctx, curNode.Data.Section)
	}

	return isMySectionOutdated
}
