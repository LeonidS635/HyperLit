package info

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func propagateStatus(
	ctx context.Context, status int, sectionsTrieNode *trie.Node[Section], sectionsStates *SectionsStates,
) {
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	sectionsTrieNode.Data.Status = status
	sectionsStates.Add(status, sectionsTrieNode.Data.Path, sectionsTrieNode, sectionsTrieNode)
	for _, childSection := range sectionsTrieNode.GetAll() {
		propagateStatus(ctx, status, childSection, sectionsStates)
	}
}
