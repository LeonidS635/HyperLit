package parser

import (
	"context"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
)

func (p *Parser) Traverse(ctx context.Context, path string) (*trie.Node[info.File], error) {
	// Reopen channels
	p.initChannels()
	defer p.closeChannels()

	filesTrieRootNode := trie.NewNode[info.File]()

	// Context for the traversal with cancel in case of an error
	traverseCtx, traverseCancel := context.WithCancel(ctx)
	defer traverseCancel()

	// Start traversing
	p.wg.Add(1)
	go p.traverse(traverseCtx, path, filesTrieRootNode)

	if err := helpers.WaitCtx(traverseCtx, &p.wg, p.errCh); err != nil {
		return nil, err
	}
	return filesTrieRootNode, nil
}

func (p *Parser) traverse(ctx context.Context, path string, curNode *trie.Node[info.File]) {
	defer p.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
		return
	}

	curNode.Data = info.File{
		IsDir: fileInfo.IsDir(),
		Path:  path,
		Size:  fileInfo.Size(),
		MTime: fileInfo.ModTime(),
	}

	if fileInfo.IsDir() {
		for _, entry := range p.getDirEntries(ctx, path) {
			match := false
			for _, glob := range p.ignoreGlobs {
				if match, _ = filepath.Match(glob, filepath.Join(path, entry.Name())); match {
					break
				}
			}
			if match {
				continue
			}

			next := curNode.Insert(entry.Name())
			p.wg.Add(1)
			go p.traverse(ctx, filepath.Join(path, entry.Name()), next)
		}
	}
}

func (p *Parser) getDirEntries(ctx context.Context, dir string) []os.DirEntry {
	ok := p.sema.Acquire(ctx)
	if !ok {
		return nil
	}
	defer p.sema.Release()

	entries, err := os.ReadDir(dir)
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
		return nil
	}
	return entries
}
