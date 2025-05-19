package hyperlit

import (
	"context"
)

func (h *HyperLit) Docs(ctx context.Context, port int) error {
	return h.docsGenerator.StartServer(ctx, port)
}
