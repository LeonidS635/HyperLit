package parser

import (
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
	"github.com/LeonidS635/HyperLit/internal/parser/globs"
)

type Parser struct {
	ignoreGlobs []string

	sema *resourceslimiter.Semaphore
	wg   sync.WaitGroup
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
