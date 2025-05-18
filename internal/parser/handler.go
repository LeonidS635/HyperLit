package parser

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

// TODO: Think about possible send on close error channel problem
// TODO: Make channels buffered
// TODO: replace global parser channels with local ones

func (p *Parser) HandleParsedSections(
	ctx context.Context, path string, handler func(ctx context.Context, section entry.Interface) error,
) (*trie.Node[info.Section], error) {
	parseCtx, parseCancel := context.WithCancel(ctx)
	defer parseCancel()

	var rootSectionsTrieNode *trie.Node[info.Section]
	var err error

	p.initChannels()
	go func() {
		rootSectionsTrieNode, err = p.parse(parseCtx, path)
		p.closeChannels()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-p.errCh:
			if err != nil {
				return nil, err
			}
		case section, ok := <-p.entriesCh:
			if !ok {
				return rootSectionsTrieNode, err
			}
			if err := handler(parseCtx, section); err != nil {
				return nil, err
			}
		}
	}
}
