package helpers

import "context"

// TODO: think about sending in closed channel panic issue

func SendCtx[ChValue any](ctx context.Context, ch chan<- ChValue, val ChValue) {
	select {
	case <-ctx.Done():
	case ch <- val:
	}
}
