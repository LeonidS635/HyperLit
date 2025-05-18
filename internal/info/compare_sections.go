package info

import (
	"context"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func CompareSectionsTries(
	ctx context.Context, newSectionsTrieNode *trie.Node[Section], oldSectionsTrieNode *trie.Node[Section],
	sectionsStates *SectionsStates,
) (*trie.Node[Section], error) {
	// newSectionsTrieNode sections can't be nil
	if oldSectionsTrieNode == nil {
		propagateStatus(ctx, StatusCreated, newSectionsTrieNode, sectionsStates)
		return newSectionsTrieNode, nil
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go compareSectionsTries(ctx, newSectionsTrieNode, oldSectionsTrieNode, sectionsStates, &wg)

	_ = helpers.WaitCtx(ctx, &wg, nil)
	return newSectionsTrieNode, ctx.Err()
}

func compareSectionsTries(
	ctx context.Context, newSectionsTrieNode *trie.Node[Section], oldSectionsTrieNode *trie.Node[Section],
	sectionsStates *SectionsStates, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	status := compareTwoSections(newSectionsTrieNode.Data, oldSectionsTrieNode.Data)
	newSectionsTrieNode.Data.Status = status
	sectionsStates.Add(status, newSectionsTrieNode.Data.Path, newSectionsTrieNode, nil)

	newSectionChildren, oldSectionChildren := newSectionsTrieNode.GetAll(), oldSectionsTrieNode.GetAll()

	seen := make(map[string]struct{})
	for name, newChild := range newSectionChildren {
		if oldChild, ok := oldSectionChildren[name]; ok {
			wg.Add(1)
			go compareSectionsTries(ctx, newChild, oldChild, sectionsStates, wg)
		} else {
			propagateStatus(ctx, StatusCreated, newChild, sectionsStates)
		}
		seen[name] = struct{}{}
	}

	for name, oldChild := range oldSectionChildren {
		if _, ok := seen[name]; !ok {
			nextNode := newSectionsTrieNode.Insert(name)
			nextNode.Data.Path = oldChild.Data.Path
			propagateStatus(ctx, StatusDeleted, nextNode, sectionsStates)
		}
	}
}
