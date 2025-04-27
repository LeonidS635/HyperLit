//package storage
//
//import (
//	"context"
//	"fmt"
//	"sync"
//
//	"github.com/LeonidS635/HyperLit/internal/helpers"
//	"github.com/LeonidS635/HyperLit/internal/helpers/resourceslimiter"
//	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
//	"github.com/LeonidS635/HyperLit/internal/vcs/objects"
//	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
//	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
//	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
//)
//
//type hashTraverser struct {
//	ObjectsStorage
//
//	sema *resourceslimiter.Semaphore
//	wg   *sync.WaitGroup
//
//	sectionsCh chan objects.Section
//	errCh      chan error
//}
//
//func newHashTraverser() *hashTraverser {
//	return &hashTraverser{
//		sema: resourceslimiter.NewSemaphore(),
//		wg:   &sync.WaitGroup{},
//	}
//}
//
//func (t *hashTraverser) getChannels() (<-chan objects.Section, <-chan error) {
//	t.sectionsCh = make(chan objects.Section, resourceslimiter.MaxOpenedEntries)
//	t.errCh = make(chan error)
//
//	return t.sectionsCh, t.errCh
//}
//
//func (t *hashTraverser) Traverse(ctx context.Context, path string) {
//	t.wg.Add(1)
//	t.traverse(ctx, path)
//	t.wg.Wait()
//
//	close(t.sectionsCh)
//}
//
//func (t *hashTraverser) traverse(ctx context.Context, hash string) {
//	defer t.wg.Done()
//	if helpers.IsCtxCancelled(ctx) {
//		return
//	}
//
//	e, err := t.LoadEntry(hash)
//	if err != nil {
//		helpers.SendCtx(ctx, t.errCh, err)
//		return
//	}
//	if e.Type != format.TreeType {
//		return
//	}
//
//	section := objects.Section{Path: e.Path}
//	needToSave := false
//
//	childEntries, err := tree.Parse(e.Data)
//	if err != nil {
//		helpers.SendCtx(ctx, t.errCh, err)
//		return
//	}
//
//	for _, childEntry := range childEntries {
//		switch childEntry.Type {
//		case format.TreeType:
//			childEntries, err := tree.Parse(e.Data)
//			if err != nil {
//				helpers.SendCtx(ctx, t.errCh, err)
//				return
//			}
//
//			for _, childEntry := range childEntries {
//				t.wg.Add(1)
//				go t.traverse(ctx, hasher.ConvertToHex(childEntry.GetHash()))
//			}
//		case format.CodeType:
//			section.Code = e
//		case format.DocsType:
//			section.Docs = e
//			needToSave = true
//		default:
//			helpers.SendCtx(ctx, t.errCh, fmt.Errorf("unknown entry type: %v", e.Type))
//			return
//		}
//	}
//
//	if needToSave {
//		helpers.SendCtx(ctx, t.sectionsCh, section)
//	}
//}
//
//func (t *hashTraverser) getChildren(ctx context.Context, content []byte) []entry.Entry {
//	ok := t.sema.Acquire(ctx)
//	if !ok {
//		return nil
//	}
//	defer t.sema.Release()
//
//	entries, err := tree.Parse(content)
//	if err != nil {
//		helpers.SendCtx(ctx, t.errCh, err)
//		return nil
//	}
//	return entries
//}
//
//func (s ObjectsStorage) Traverse(ctx context.Context, rootHash string) {
//	s.traverse(ctx, rootHash, make(chan<- objects.Section))
//}
//
//func (s ObjectsStorage) traverse(ctx context.Context, hash string, sectionsCh chan<- objects.Section) {
//	if helpers.IsCtxCancelled(ctx) {
//		return
//	}
//
//	e, err := s.LoadEntry(hash)
//	if err != nil {
//		//helpers.SendCtx(ctx, t.errCh, err)
//		return
//	}
//	if e.Type != format.TreeType {
//		return
//	}
//
//	section := objects.Section{Path: e.Path}
//	needToSave := false
//
//	childEntries, err := tree.Parse(e.Data)
//	if err != nil {
//		//helpers.SendCtx(ctx, t.errCh, err)
//		return
//	}
//
//	for _, childEntry := range childEntries {
//		switch childEntry.Type {
//		case format.TreeType:
//			childEntries, err := tree.Parse(e.Data)
//			if err != nil {
//				//helpers.SendCtx(ctx, t.errCh, err)
//				return
//			}
//
//			for _, childEntry := range childEntries {
//				s.traverse(ctx, hasher.ConvertToHex(childEntry.GetHash()), sectionsCh)
//			}
//		case format.CodeType:
//			section.Code = e
//		case format.DocsType:
//			section.Docs = e
//			needToSave = true
//		default:
//			//helpers.SendCtx(ctx, t.errCh, fmt.Errorf("unknown entry type: %v", e.Type))
//		}
//	}
//
//	if needToSave {
//		helpers.SendCtx(ctx, sectionsCh, section)
//	}
//}

package storage

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/storage/hashtraverser"
)

func (s ObjectsStorage) Traverse(ctx context.Context, rootHash string) (*trie.Node[info.Section], <-chan error) {
	t := hashtraverser.NewHashTraverser(s.LoadEntry)
	go t.Traverse(ctx, rootHash)
	return t.GetOutputs()
}
