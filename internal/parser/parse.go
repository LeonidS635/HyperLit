package parser

import (
	"context"
	"os"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *parserWithChannels) parse(ctx context.Context, path string) (*trie.Node[info.Section], error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	section, err := tree.Prepare(fileInfo.Name())
	if err != nil {
		return nil, err
	}
	sectionsTrieRootNode := trie.NewNode[info.Section]()

	if fileInfo.IsDir() {
		p.wg.Add(1)
		go p.parseDir(ctx, path, section, sectionsTrieRootNode)
	} else {
		p.wg.Add(1)
		go p.parseFile(ctx, path, section, sectionsTrieRootNode)
	}

	if err = helpers.WaitCtx(ctx, &p.wg, p.errCh); err != nil {
		return nil, err
	}
	return sectionsTrieRootNode, nil
}
