package helpers

import (
	"context"
	"sync"
)

func WaitCtx(ctx context.Context, wg *sync.WaitGroup, errCh <-chan error) error {
	done := make(chan struct{})

	go func() {
		if wg != nil {
			wg.Wait()
			close(done)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-done:
		return nil
	}
}
