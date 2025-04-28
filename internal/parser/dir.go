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

func (p *Parser) parseDir(ctx context.Context, path string, rootSection *tree.Tree, rootNode *trie.Node[info.Section]) {
	p.wg.Add(1)
	go p.parseDirSection(ctx, path, rootSection, rootNode)
	p.wg.Wait()
}

func (p *Parser) parseDirSection(
	ctx context.Context, path string, section *tree.Tree, curNode *trie.Node[info.Section],
) {
	defer p.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	sectionWg := sync.WaitGroup{}
	done := make(chan struct{})

	for _, file := range getDirEntries(ctx, path, p.sema, p.errCh) {
		filePath := filepath.Join(path, file.Name())

		subSection, err := tree.Prepare(file.Name())
		if err != nil {
			helpers.SendCtx(ctx, p.errCh, err)
			return
		}
		nextNode := curNode.Insert(file.Name())

		sectionWg.Add(1)
		if file.IsDir() {
			p.wg.Add(1)
			go func() {
				defer sectionWg.Done()
				p.parseDirSection(ctx, filePath, subSection, nextNode)
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
	go func() {
		sectionWg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return
	case <-done:
	}

	helpers.SendCtx(ctx, p.sectionsCh, Section(section))

	curNode.Data = info.Section{
		Hash:     hasher.ConvertToHex(section.GetHash()),
		CodeHash: "",
		DocsHash: "",

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
