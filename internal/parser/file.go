package parser

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser/sections"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) parseFile(
	ctx context.Context, path string, section *tree.Tree, sectionsTrieNode *trie.Node[info.Section],
) {
	defer p.wg.Done()
	if helpers.IsCtxCancelled(ctx) {
		return
	}

	for _, glob := range p.ignoreGlobs {
		if match, _ := filepath.Match(glob, path); match {
			return
		}
	}

	ok := p.sema.Acquire(ctx)
	if !ok {
		return
	}
	defer p.sema.Release()

	file, err := os.Open(path)
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			helpers.SendCtx(ctx, p.errCh, err)
		}
	}()

	fileStat, err := file.Stat()
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
		return
	}

	sectionsParser, err := sections.NewParser(path, fileStat.ModTime(), bufio.NewScanner(file), p.entriesCh)
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
	}

	if err = sectionsParser.Parse(ctx, path, section, sectionsTrieNode); err != nil {
		helpers.SendCtx(ctx, p.errCh, fmt.Errorf("error parsing %s: %w", path, err))
	}
}
