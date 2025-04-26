package resourceslimiter

import "context"

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore() *Semaphore {
	return &Semaphore{make(chan struct{}, MaxOpenedEntries)}
}

func (s *Semaphore) Acquire(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case s.ch <- struct{}{}:
		return true
	}
}

func (s *Semaphore) Release() {
	<-s.ch
}
