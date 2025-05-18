package hyperlit

import (
	"context"
)

func (h *HyperLit) Docs(ctx context.Context) error {
	return h.docsGenerator.StartServer()
}
