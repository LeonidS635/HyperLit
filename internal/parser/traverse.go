package parser

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser/dirtraverser"
)

func (p *Parser) Traverse(ctx context.Context, rootPath string) (*trie.Node[info.File], error) {
	t := dirtraverser.NewDirTraverser()
	root, errCh := t.GetOutputs()

	readCtx, readCancel := context.WithCancel(ctx)
	defer readCancel()

	go t.Traverse(readCtx, rootPath)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err, ok := <-errCh:
		if ok {
			return nil, err
		}
		return root, nil
	}
}
