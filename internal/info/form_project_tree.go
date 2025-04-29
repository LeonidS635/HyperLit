package info

import (
	"context"
	"fmt"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
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
		fmt.Printf(
			"Unmodified section with hash %q and name %q\n", hasher.ConvertToHex(sections.Data.This.GetHash()),
			sections.Data.This.GetName(),
		)
		fmt.Printf(
			"Project node %q has hash %q\n", project.Data.Section.GetName(),
			hasher.ConvertToHex(project.Data.Section.GetHash()),
		)
		for childName, child := range sections.GetAll() {
			fmt.Printf(
				"Child of %q: %q has hash %q\n", hasher.ConvertToHex(sections.Data.This.GetHash()), childName,
				hasher.ConvertToHex(child.Data.This.GetHash()),
			)
			next := project.Insert(childName)
			next.Data.Section = child.Data.This
			next.Data.Status = StatusUnmodified

			wg.Add(1)
			go formProjectTree(ctx, child, next, wg)
		}
	default:
	}
}
