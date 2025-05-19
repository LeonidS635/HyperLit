package parser

import (
	"context"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

func (p *Parser) HandleParsedSections(
	ctx context.Context, path string, handler func(ctx context.Context, section entry.Interface) error,
) (*trie.Node[info.Section], error) {
	parseCtx, parseCancel := context.WithCancel(ctx)
	defer parseCancel()

	localParser := newParserWithChannels(p)
	var rootSectionsTrieNode *trie.Node[info.Section]

	var parseWg sync.WaitGroup
	errCh := make(chan error)

	parseWg.Add(1)
	go func() {
		defer parseWg.Done()

		var err error
		rootSectionsTrieNode, err = localParser.parse(parseCtx, path)
		if err != nil {
			helpers.SendCtx(parseCtx, errCh, err)
			return
		}

		close(localParser.entriesCh)
	}()

	parseWg.Add(1)
	go func() {
		defer parseWg.Done()

		for {
			select {
			case <-parseCtx.Done():
				return
			case section, ok := <-localParser.entriesCh:
				if !ok {
					return
				}
				if err := handler(parseCtx, section); err != nil {
					helpers.SendCtx(parseCtx, errCh, err)
					return
				}
			}
		}
	}()

	if err := helpers.WaitCtx(ctx, &parseWg, errCh); err != nil {
		return nil, err
	}
	return rootSectionsTrieNode, nil
}
