package parser

import (
	"context"
	"os"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) Parse(ctx context.Context, path string) ([]byte, *trie.Node[info.Section]) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
		return nil, nil
	}

	section, err := tree.Prepare(fileInfo.Name(), path)
	if err != nil {
		helpers.SendCtx(ctx, p.errCh, err)
		return nil, nil
	}
	rootNode := trie.NewNode[info.Section]()

	if fileInfo.IsDir() {
		p.parseDir(ctx, path, section, rootNode)
	} else {
		p.wg.Add(1)
		p.parseFile(ctx, path, section, rootNode)
	}

	return section.GetHash(), rootNode
}
