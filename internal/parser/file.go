package parser

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) parseFile(
	ctx context.Context, path string, fileSection *tree.Tree, fileNode *trie.Node[info.Section],
) {
	defer p.wg.Done()

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
	defer file.Close()

	lineNumber := 0
	if err = p.parseSection(ctx, bufio.NewScanner(file), &lineNumber, fileSection, fileNode); err != nil {
		helpers.SendCtx(ctx, p.errCh, fmt.Errorf("error parsing %s: %w", path, err))
	}
}
