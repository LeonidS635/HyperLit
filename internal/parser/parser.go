package parser

import (
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

type Syntax struct {
	DocsStartSeq []byte
	DocsEndSeq   []byte
	CodeEndSeq   []byte
}

type Section = entry.Interface

type Parser struct {
	syntax Syntax

	sema *resourceslimiter.Semaphore
	wg   *sync.WaitGroup

	sectionsCh chan Section
	errCh      chan error
}

func NewParser() *Parser {
	return &Parser{
		syntax: Syntax{
			DocsStartSeq: []byte("@@docs"),
			DocsEndSeq:   []byte("@@/docs"),
			CodeEndSeq:   []byte("@@/code"),
		},
		sema: resourceslimiter.NewSemaphore(),
		wg:   &sync.WaitGroup{},
	}
}

func (p *Parser) InitChannels() (<-chan Section, <-chan error) {
	p.sectionsCh = make(chan Section)
	p.errCh = make(chan error)
	return p.sectionsCh, p.errCh
}

func (p *Parser) CloseChannels() {
	close(p.sectionsCh)
	close(p.errCh)
}
