package hyperlit

import (
	"context"
	"fmt"
	"os"
)

func (h *HyperLit) Docs(ctx context.Context) {
	if err := h.docsGenerator.StartServer(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
