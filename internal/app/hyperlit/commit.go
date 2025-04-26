package hyperlit

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

func (h *HyperLit) Commit(ctx context.Context) {
	changedFiles := h.Status(ctx)

	parseCtx, parseCancel := context.WithCancel(ctx)
	defer parseCancel()

	parseWg := sync.WaitGroup{}
	parseWg.Add(1)

	sectionsCh := h.parser.Sections()

	go func() {
		defer parseWg.Done()

		err := h.vcs.Save(parseCtx, sectionsCh)
		if err != nil {
			log.Println(err)
		}
		parseCancel()
	}()

	for _, file := range changedFiles {
		if err := h.parser.Parse(parseCtx, file.path); err != nil {
			parseCancel()
			fmt.Println(err)
			os.Exit(1)
		}
	}

	parseWg.Wait()

	h.removeUnused()
}
