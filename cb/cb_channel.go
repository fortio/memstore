// FIFO Queue with fixed capacity.
// Circular Buffer implementation in Go with both
// pub/sub thread safe blocking API and pure FIFO queue with set capacity
// unsynchronized base.
// Channel / go idiomatic version.
package cb

type CircularBufferChan[T any] struct {
	buffer chan T
}

func NewC[T any](capacity int) Queue[T] {
	cb := &CircularBufferChan[T]{
		buffer: make(chan T, capacity),
	}
	return cb
}

func (cb *CircularBufferChan[T]) Empty() bool {
	return len(cb.buffer) == 0
}

func (cb *CircularBufferChan[T]) Full() bool {
	return len(cb.buffer) == cap(cb.buffer)
}

func (cb *CircularBufferChan[T]) Size() int {
	return len(cb.buffer)
}

func (cb *CircularBufferChan[T]) Capacity() int {
	return cap(cb.buffer)
}

// Push adds an item to the queue. returns false if queue is full.
// Note: might block and not return false at times. Use PushBlocking for
// correct version.
func (cb *CircularBufferChan[T]) Push(item T) bool {
	// Note: this is for equivalent with the array/conditional variable version
	// but isn't correct, as in Full can return false and yet the push can block
	// if another producer enqueued in the meanwhile.
	if cb.Full() {
		return false
	}
	cb.buffer <- item
	return true
}

// Pop removes an item from the queue. returns false if queue is empty.
// Note: might block and not return false at times. Use PopBlocking for
// correct version.
func (cb *CircularBufferChan[T]) Pop() (value T, ok bool) {
	if cb.Empty() {
		return
	}
	ok = true
	value = <-cb.buffer
	return
}

// Thread safe blocking versions:

// Push adds an item to the queue. blocks if queue is full.
func (cb *CircularBufferChan[T]) PushBlocking(item T) {
	cb.buffer <- item
}

// Pop removes an item from the queue. blocks if queue is empty.
func (cb *CircularBufferChan[T]) PopBlocking() T {
	return <-cb.buffer
}
