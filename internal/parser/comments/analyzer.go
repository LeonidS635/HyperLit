package comments

import (
	"fmt"
	"path/filepath"
)

type Analyzer struct {
	syntax Syntax

	isInMultiLineSection bool
}

func NewAnalyzer(filename string) (*Analyzer, error) {
	ext := filepath.Ext(filename)
	if s, ok := commentsSyntax[ext]; ok {
		return &Analyzer{syntax: s, isInMultiLineSection: false}, nil
	}
	return nil, fmt.Errorf("unsupported file extension: %q", ext)
}
