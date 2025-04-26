package parser

//import (
//	"context"
//	"os"
//	"path/filepath"
//	"sync"
//
//	"github.com/LeonidS635/HyperLit/internal/helpers"
//	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
//)
//
//func (p *Parser) ParseDir(ctx context.Context, path string) {
//	if helpers.IsCtxCancelled(ctx) {
//		return
//	}
//
//	files, err := os.ReadDir(path)
//	if err != nil {
//		return nil, err
//	}
//
//	section, err := tree.Prepare(filepath.Base(path), path)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, file := range files {
//		if file.IsDir() {
//			dirSection, err := p.ParseDir(ctx, filepath.Join(path, file.Name()))
//			if err != nil {
//				return nil, err
//			}
//			section.RegisterEntry(dirSection)
//		} else {
//			fileSection, err := p.ParseFile(ctx, filepath.Join(path, file.Name()))
//			if err != nil {
//				return nil, err
//			}
//			section.RegisterEntry(fileSection)
//		}
//	}
//	return section, nil
//}
//
//func (p *Parser) parseDir(ctx context.Context, path, name string, wg *sync.WaitGroup) (Section, error) {
//	if helpers.IsCtxCancelled(ctx) {
//		return nil, nil
//	}
//
//	files, err := os.ReadDir(path)
//	if err != nil {
//		return nil, err
//	}
//
//	section, err := tree.Prepare(name)
//	if err != nil {
//		return nil, err
//	}
//
//	errs := make(chan error)
//	for _, file := range files {
//		if file.IsDir() {
//			//wg.Add(1)
//			//go func() {
//			//	defer wg.Done()
//			dirSection, err := p.parseDir(ctx, filepath.Join(path, file.Name()), file.Name(), wg)
//			if err != nil {
//				//errs <- err
//				return nil, err
//			}
//			section.RegisterEntry(dirSection)
//			//}()
//		} else {
//			fileSection, err := p.ParseFile(ctx, filepath.Join(path, file.Name()), file.Name())
//			if err != nil {
//				return nil, err
//			}
//			section.RegisterEntry(fileSection)
//		}
//	}
//
//	wg.Wait()
//
//	select {
//	case <-ctx.Done():
//		return nil, nil
//	case err = <-errs:
//		return nil, err
//	default:
//		return section, nil
//	}
//}
