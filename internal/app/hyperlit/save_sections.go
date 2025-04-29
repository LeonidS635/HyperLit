package hyperlit

import (
	"context"
	"fmt"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
)

func (h *HyperLit) saveSections(ctx context.Context) {
	fmt.Println(
		"Root section", hasher.ConvertToHex(h.rootSection.Data.Section.GetHash()), "has status",
		h.rootSection.Data.Status,
	)
	if h.saveSection(ctx, h.rootSection, h.projectPath) {
		h.vcs.SaveRootHash(h.rootSection.Data.Section.GetHash())
	}
}

func (h *HyperLit) saveSection(ctx context.Context, curNode *trie.Node[info.TrieSection], name string) bool {
	fmt.Println(
		"Node", hasher.ConvertToHex(curNode.Data.Section.GetHash()), "has children:",
		curNode.GetAll(), "and status", curNode.Data.Status,
	)

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
	}

	if isMySectionOutdated {
		fmt.Printf("Saving modified %q\n", curNode.Data.Section.GetName())
		h.vcs.SaveNewEntry(ctx, curNode.Data.Section)
	} else {
		fmt.Printf("Saving unmodified %q\n", curNode.Data.Section.GetName())
		h.vcs.SaveOldEntry(ctx, curNode.Data.Section)
	}

	return isMySectionOutdated
}
