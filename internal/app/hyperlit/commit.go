//package hyperlit
//
//import (
//	"context"
//	"fmt"
//	"log"
//	"os"
//	"sync"
//
//	"github.com/LeonidS635/HyperLit/internal/info"
//	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
//)
//
//func (h *HyperLit) Commit(ctx context.Context) {
//	if err := h.getSectionsStatuses(ctx); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	changedFiles := h.sectionsStatuses
//	fmt.Println(changedFiles)
//
//	parseCtx, parseCancel := context.WithCancel(ctx)
//	defer parseCancel()
//
//	sectionsCh, errCh := h.parser.InitChannels()
//	done := make(chan struct{})
//
//	go func() {
//		defer close(done)
//
//		if err := h.vcs.Save(parseCtx, sectionsCh); err != nil {
//			fmt.Println(err)
//		}
//		parseCancel()
//	}()
//
//	wg := sync.WaitGroup{}
//	for _, sectionStatus := range changedFiles.Get(info.StatusCreated) {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			rootHash, _ := h.parser.Parse(parseCtx, sectionStatus.Path)
//			if rootHash != nil {
//				fmt.Println(hasher.ConvertToHex(rootHash))
//				err := h.vcs.SaveRootHash(rootHash)
//				if err != nil {
//					log.Println(err)
//				}
//				parseCancel()
//			}
//		}()
//	}
//
//	go func() {
//		wg.Wait()
//		h.parser.CloseChannels()
//	}()
//
//	select {
//	case <-ctx.Done():
//		return
//	case err, ok := <-errCh:
//		if ok {
//			log.Println(err)
//		}
//	case <-done:
//	}
//
//	h.removeUnused()
//}

package hyperlit

import (
	"context"
	"fmt"
	"os"
)

func (h *HyperLit) Commit(ctx context.Context) {
	if err := h.getSectionsStatuses(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := h.commitSections(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h.sectionsStatuses.Print()

	h.removeUnused()
}
