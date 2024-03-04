package semaphore

type Semaphore struct {
	ch chan struct{}
}

func New(n int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, n),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}
