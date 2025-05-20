package hyperlit

import (
	"context"
)

func (h *HyperLit) Docs(ctx context.Context, port int, md bool) error {
	return h.docsGenerator.StartServer(ctx, port, md)
}
