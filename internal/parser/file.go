package parser

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) ParseFile(ctx context.Context, path string) error {
	if helpers.IsCtxCancelled(ctx) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	section, err := tree.Prepare(filepath.Base(path), path)
	if err != nil {
		return err
	}

	lineNumber := 0
	if err = p.parseSection(ctx, bufio.NewScanner(file), &lineNumber, section); err != nil {
		return fmt.Errorf("error parsing %s: %w", path, err)
	}
	return nil
}
