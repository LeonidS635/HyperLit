package parser

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) parseDir(
	ctx context.Context, path string, section *tree.Tree, sectionsTrieNode *trie.Node[info.Section],
) {
	defer p.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	sectionWg := sync.WaitGroup{}
	for _, file := range getDirEntries(ctx, path, p.sema, p.errCh) {
		filePath := filepath.Join(path, file.Name())

		match := false
		for _, glob := range p.ignoreGlobs {
			if match, _ = filepath.Match(glob, filePath); match {
				break
			}
		}
		if match {
			continue
		}

		subSection, err := tree.Prepare(file.Name())
		if err != nil {
			helpers.SendCtx(ctx, p.errCh, err)
			return
		}
		nextNode := sectionsTrieNode.Insert(file.Name())

		sectionWg.Add(1)
		if file.IsDir() {
			p.wg.Add(1)
			go func() {
				defer sectionWg.Done()
				p.parseDir(ctx, filePath, subSection, nextNode)
				section.RegisterEntry(subSection)
			}()
		} else {
			p.wg.Add(1)
			go func() {
				defer sectionWg.Done()
				p.parseFile(ctx, filePath, subSection, nextNode)
				section.RegisterEntry(subSection)
			}()
		}
	}

	_ = helpers.WaitCtx(ctx, &sectionWg, nil)

	// Fill data in a global sections tree
	sectionsTrieNode.Data = info.Section{
		Path: path,
		Hash: hasher.ConvertToHex(section.GetHash()),

		MTime: time.Now(),
		This:  section,
	}
}

func getDirEntries(
	ctx context.Context, path string, sema *resourceslimiter.Semaphore, errCh chan<- error,
) []os.DirEntry {
	ok := sema.Acquire(ctx)
	if !ok {
		return nil
	}
	defer sema.Release()

	entries, err := os.ReadDir(path)
	if err != nil {
		helpers.SendCtx(ctx, errCh, err)
		return nil
	}
	return entries
}
