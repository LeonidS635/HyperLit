package parser

import (
	"context"
	"os"
)

func (p *Parser) Parse(ctx context.Context, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	defer p.Close()

	if info.IsDir() {
		return p.ParseFile(ctx, path)
	}
	return p.ParseFile(ctx, path)
}
