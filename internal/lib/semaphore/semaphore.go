package semaphore

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, size),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}
