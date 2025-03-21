package concurrency

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(maxTickets int) *Semaphore {
	return &Semaphore{
		tickets: make(chan struct{}, maxTickets),
	}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}
