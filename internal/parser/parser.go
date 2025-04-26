package parser

import (
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

	sectionsCh chan Section
	errCh      chan error

	resourcesControlSema *resourceslimiter.Semaphore
}

func NewParser() *Parser {
	return &Parser{
		syntax: Syntax{
			DocsStartSeq: []byte("@@docs"),
			DocsEndSeq:   []byte("@@/docs"),
			CodeEndSeq:   []byte("@@/code"),
		},
		resourcesControlSema: resourceslimiter.NewSemaphore(),
	}
}

// Sections must be called before Parse
func (p *Parser) Sections() <-chan Section {
	p.sectionsCh = make(chan Section)
	return p.sectionsCh
}

// Close must be called after Parse finished
func (p *Parser) Close() {
	close(p.sectionsCh)
}
