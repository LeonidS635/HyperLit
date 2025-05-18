package parser

import (
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/parser/globs"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

type Parser struct {
	ignoreGlobs []string

	sema *resourceslimiter.Semaphore
	wg   sync.WaitGroup

	entriesCh chan entry.Interface
	errCh     chan error
}

func NewParser(projectPath, hlPath string) *Parser {
	p := &Parser{
		ignoreGlobs: globs.GetGlobsToIgnore(projectPath),
		sema:        resourceslimiter.NewSemaphore(),
		wg:          sync.WaitGroup{},
	}
	p.ignoreGlobs = append(p.ignoreGlobs, hlPath)
	return p
}

func (p *Parser) initChannels() {
	p.entriesCh = make(chan entry.Interface)
	p.errCh = make(chan error)
}

func (p *Parser) closeChannels() {
	close(p.entriesCh)
	close(p.errCh)
}
