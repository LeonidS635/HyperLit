package info

import (
	"context"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

func FormProjectTree(ctx context.Context, sections *trie.Node[Section], project *trie.Node[TrieSection]) {
	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go formProjectTree(ctx, sections, project, &wg)

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}
}

func formProjectTree(
	ctx context.Context, sections *trie.Node[Section], project *trie.Node[TrieSection], wg *sync.WaitGroup,
) {
	defer wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	switch project.Data.Status {
	case StatusUnmodified:
		for childName, child := range sections.GetAll() {
			next := project.Insert(childName)
			next.Data.Section = child.Data.This
			next.Data.Status = StatusUnmodified

			wg.Add(1)
			go formProjectTree(ctx, child, next, wg)
		}
	default:
	}
}
