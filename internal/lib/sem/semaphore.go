package sem

// Semaphore is the interface that wraps basic Acquire and Release
// methods to control access to a limited resource pool.
type Semaphore interface {
	// Acquire blocks until a slot is available in the semaphore.
	Acquire()

	// Release frees a slot in the semaphore.
	Release()
}

// semaphore is a concrete implementation of the Semaphore interface
// using a buffered channel.
type semaphore struct {
	ch chan struct{}
}

// New returns a Semaphore that allows up to n concurrent holders.
func New(n int) Semaphore {
	return &semaphore{
		ch: make(chan struct{}, n),
	}
}

// Acquire takes a slot in the semaphore, blocking if none are available.
func (s *semaphore) Acquire() { s.ch <- struct{}{} }

// Release frees a previously acquired slot in the semaphore.
func (s *semaphore) Release() { <-s.ch }
