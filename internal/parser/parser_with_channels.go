package parser

import (
	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

type parserWithChannels struct {
	*Parser

	entriesCh chan entry.Interface
	errCh     chan error
}

func newParserWithChannels(p *Parser) *parserWithChannels {
	return &parserWithChannels{
		Parser:    p,
		entriesCh: make(chan entry.Interface, resourceslimiter.MaxOpenedEntries),
		errCh:     make(chan error, 1),
	}
}
