package helpers

import "context"

func SendCtx[ChValue any](ctx context.Context, ch chan<- ChValue, val ChValue) {
	select {
	case <-ctx.Done():
	case ch <- val:
	}
}
